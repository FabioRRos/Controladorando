package automacao

import (
	"despesas/model/entity"
	"fmt"
	"os"
	"strconv"

	"github.com/playwright-community/playwright-go"
)

// headlessFromEnv lê HEADLESS=1/true para rodar headless (default: false para debug).
func headlessFromEnv() bool {
	v := os.Getenv("HEADLESS")
	if v == "" {
		return false
	}
	b, err := strconv.ParseBool(v)
	if err != nil {
		return false
	}
	return b
}

func toInt(v any) (int, bool) {
	switch x := v.(type) {
	case int:
		return x, true
	case int32:
		return int(x), true
	case int64:
		return int(x), true
	case float32:
		return int(x), true
	case float64:
		return int(x), true
	default:
		return 0, false
	}
}

// newPW cria Playwright + Browser + Context + Page, com timeouts e opções úteis.
func newPW(headless bool) (*playwright.Playwright, playwright.Browser, playwright.BrowserContext, playwright.Page, error) {
	// Se você já instala os browsers no build/CI, pode remover isso.
	if err := playwright.Install(); err != nil {
		return nil, nil, nil, nil, fmt.Errorf("playwright install: %w", err)
	}

	pw, err := playwright.Run()
	if err != nil {
		return nil, nil, nil, nil, fmt.Errorf("playwright run: %w", err)
	}

	browser, err := pw.Chromium.Launch(playwright.BrowserTypeLaunchOptions{
		Headless: playwright.Bool(headless),
	})
	if err != nil {
		pw.Stop()
		return nil, nil, nil, nil, fmt.Errorf("chromium launch: %w", err)
	}

	context, err := browser.NewContext(playwright.BrowserNewContextOptions{
		IgnoreHttpsErrors: playwright.Bool(true),
		AcceptDownloads:   playwright.Bool(true), // <<< ESSENCIAL
	})
	if err != nil {
		_ = browser.Close()
		pw.Stop()
		return nil, nil, nil, nil, fmt.Errorf("new context: %w", err)
	}

	page, err := context.NewPage()
	if err != nil {
		_ = context.Close()
		_ = browser.Close()
		pw.Stop()
		return nil, nil, nil, nil, fmt.Errorf("new page: %w", err)
	}

	// Timeouts padrão
	page.SetDefaultTimeout(60_000)
	page.SetDefaultNavigationTimeout(60_000)

	return pw, browser, context, page, nil
}

func BaixarEmpenhos(idEmpenho int, tipo string, caso2 bool) error {
	headless := headlessFromEnv()

	fmt.Println("Iniciando Playwright...")
	pw, browser, context, page, err := newPW(headless)
	if err != nil {
		return err
	}
	defer func() {
		fmt.Println("Encerrando navegador...")
		_ = page.Close()
		_ = context.Close()
		_ = browser.Close()
		pw.Stop()
	}()

	varlink := entity.LinkDespesas()

	fmt.Println("Navegando para:", varlink)
	if _, err = page.Goto(varlink, playwright.PageGotoOptions{
		WaitUntil: playwright.WaitUntilStateNetworkidle,
		Timeout:   playwright.Float(60_000),
	}); err != nil {
		return fmt.Errorf("navegar para portal: %w", err)
	}
	fmt.Println("Página carregada")

	frame := page.FrameLocator("iframe").First()

	// ===============================
	// Clique 1 - Despesas
	// ===============================
	fmt.Println("Aguardando botão Despesas...")
	btn := frame.Locator("#btnAtalhoDespesa")

	if err := btn.WaitFor(playwright.LocatorWaitForOptions{
		State:   playwright.WaitForSelectorStateVisible,
		Timeout: playwright.Float(60_000),
	}); err != nil {
		return fmt.Errorf("aguardar botão atalho despesas: %w", err)
	}

	fmt.Println("Clicando em Despesas...")
	if err := btn.Click(); err != nil {
		return fmt.Errorf("clicar no atalho despesas: %w", err)
	}
	fmt.Println("Clique em Despesas realizado")

	page.WaitForTimeout(200)

	// ===============================
	// Filtro
	// ===============================
	fmt.Println("Aguardando campo de filtro...")
	filtro := frame.Locator("#gridDespesas_DXFREditorcol3_I")

	if err := filtro.WaitFor(playwright.LocatorWaitForOptions{
		State:   playwright.WaitForSelectorStateVisible,
		Timeout: playwright.Float(60_000),
	}); err != nil {
		return fmt.Errorf("aguardar filtro de id do empenho: %w", err)
	}

	fmt.Println("Preenchendo filtro com ID:", idEmpenho)
	if err := filtro.Fill(strconv.Itoa(idEmpenho)); err != nil {
		return fmt.Errorf("preencher filtro idEmpenho: %w", err)
	}

	fmt.Println("Confirmando filtro (Enter)...")
	if err := filtro.Press("Enter"); err != nil {
		return fmt.Errorf("confirmar filtro (Enter): %w", err)
	}

	page.WaitForTimeout(3000)

	// ===============================
	// Clique 2 - Detalhes
	// ===============================
	fmt.Println("Aguardando botão Detalhes...")
	btnDetalhes := frame.Locator("img[title='Detalhes do Empenho']").First()

	if err := btnDetalhes.WaitFor(playwright.LocatorWaitForOptions{
		State:   playwright.WaitForSelectorStateVisible,
		Timeout: playwright.Float(60_000),
	}); err != nil {
		return fmt.Errorf("aguardar botão Detalhes do Empenho: %w", err)
	}

	fmt.Println("Clicando em Detalhes...")
	if err := btnDetalhes.Click(); err != nil {
		return fmt.Errorf("clicar em Detalhes do Empenho: %w", err)
	}
	fmt.Println("Clique em Detalhes realizado")

	page.WaitForTimeout(4000)

	// ===============================
	// Captura de Cookies de Sessão
	// ===============================
	fmt.Println("Extraindo cookies da sessão...")

	// Captura todos os cookies do contexto atual
	cookies, err := context.Cookies()
	if err != nil {
		return fmt.Errorf("erro ao capturar cookies: %w", err)
	}

	// Criando um mapa para você salvar e usar as variáveis como quiser
	cookieMap := make(map[string]string)
	for _, c := range cookies {
		cookieMap[c.Name] = c.Value
		// Opcional: Logar para debug
		// fmt.Printf("Cookie encontrado: %s = %s\n", c.Name, c.Value)
	}

	// Identificando o cookie de sessão comum em portais DevExpress (ASP.NET)
	if sessao, ok := cookieMap["ASP.NET_SessionId"]; ok {
		fmt.Println("Cookie é> ", sessao)
	} else {
		fmt.Println("⚠️ Cookie ASP.NET_SessionId não encontrado. Verifique os nomes disponíveis no mapa.")
	}

	// Se você precisar de outros cookies (como os da Assessor ou Pronim),
	// eles estarão todos dentro do cookieMap.

	return nil
}
