package automacao

import (
	"despesas/model/entity"
	"despesas/utils"
	"fmt"
	"os"
	"strconv"
	"strings"

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
		AcceptDownloads:   playwright.Bool(true),
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

	page.SetDefaultTimeout(60_000)
	page.SetDefaultNavigationTimeout(60_000)

	return pw, browser, context, page, nil
}

func BaixarEmpenhos(idEmpenho int, tipo string, caso2 bool) (string, error) {
	headless := headlessFromEnv()

	utils.Logger.Println("Iniciando Playwright...")
	pw, browser, context, page, err := newPW(headless)
	if err != nil {
		return "", err
	}
	defer func() {
		utils.Logger.Println("Encerrando navegador...")
		_ = page.Close()
		_ = context.Close()
		_ = browser.Close()
		pw.Stop()
	}()

	varlink := entity.LinkDespesas()

	utils.Logger.Println("Navegando para:", varlink)

	if _, err = page.Goto(varlink, playwright.PageGotoOptions{
		WaitUntil: playwright.WaitUntilStateNetworkidle,
		Timeout:   playwright.Float(60_000),
	}); err != nil {
		return "", fmt.Errorf("navegar para portal: %w", err)
	}
	utils.Logger.Println("Página carregada")

	frame := page.FrameLocator("iframe").First()

	// ===============================
	// Clique 1 - Despesas
	// ===============================
	utils.Logger.Println("Aguardando botão Despesas...")
	btn := frame.Locator("#btnAtalhoDespesa")

	if err := btn.WaitFor(playwright.LocatorWaitForOptions{
		State:   playwright.WaitForSelectorStateVisible,
		Timeout: playwright.Float(60_000),
	}); err != nil {
		return "", fmt.Errorf("aguardar botão atalho despesas: %w", err)
	}

	utils.Logger.Println("Clicando em Despesas...")
	if err := btn.Click(); err != nil {
		return "", fmt.Errorf("clicar no atalho despesas: %w", err)
	}
	utils.Logger.Println("Clique em Despesas realizado")

	page.WaitForTimeout(200)

	// ===============================
	// Filtro
	// ===============================
	utils.Logger.Println("Aguardando campo de filtro...")
	filtro := frame.Locator("#gridDespesas_DXFREditorcol3_I")

	if err := filtro.WaitFor(playwright.LocatorWaitForOptions{
		State:   playwright.WaitForSelectorStateVisible,
		Timeout: playwright.Float(60_000),
	}); err != nil {
		return "", fmt.Errorf("aguardar filtro de id do empenho: %w", err)
	}

	utils.Logger.Println("Preenchendo filtro com ID:", idEmpenho)
	if err := filtro.Fill(strconv.Itoa(idEmpenho)); err != nil {
		return "", fmt.Errorf("preencher filtro idEmpenho: %w", err)
	}

	utils.Logger.Println("Confirmando filtro (Enter)...")
	if err := filtro.Press("Enter"); err != nil {
		return "", fmt.Errorf("confirmar filtro (Enter): %w", err)
	}

	page.WaitForTimeout(3000)

	// ===============================
	// Clique 2 - Detalhes
	// ===============================
	utils.Logger.Println("Aguardando botão Detalhes...")
	btnDetalhes := frame.Locator("img[title='Detalhes do Empenho']").First()

	if err := btnDetalhes.WaitFor(playwright.LocatorWaitForOptions{
		State:   playwright.WaitForSelectorStateVisible,
		Timeout: playwright.Float(60_000),
	}); err != nil {
		return "", fmt.Errorf("aguardar botão Detalhes do Empenho: %w", err)
	}

	utils.Logger.Println("Clicando em Detalhes...")
	if err := btnDetalhes.Click(); err != nil {
		return "", fmt.Errorf("clicar em Detalhes do Empenho: %w", err)
	}
	utils.Logger.Println("Clique em Detalhes realizado")

	// ===============================
	// CASO 2 - MULTIPLOS
	// ===============================

	if caso2 {
		utils.Logger.Println("CASO 2 - Múltiplos")
		tipoParam := strings.ToUpper(strings.TrimSpace(tipo))

		// 1. Definir o locator das linhas
		rowsLocator := frame.Locator("tr[id^='gridDespesasEmpenhos_DXDataRow']")

		// 2. FUNDAMENTAL: Aguardar que pelo menos uma linha apareça (o Count() sozinho não espera)
		if err := rowsLocator.First().WaitFor(playwright.LocatorWaitForOptions{
			State:   playwright.WaitForSelectorStateVisible,
			Timeout: playwright.Float(30_000),
		}); err != nil {
			return "", fmt.Errorf("o grid de empenhos não carregou a tempo: %w", err)
		}

		// 3. Pegar todas as linhas de uma vez para processar em memória (mais rápido)
		rows, err := rowsLocator.All()
		if err != nil {
			return "", fmt.Errorf("erro ao listar linhas: %w", err)
		}

		utils.Logger.Printf("Total de linhas encontradas: %d", len(rows))

		var targetRow playwright.Locator
		for i, row := range rows {
			// Buscamos o texto da 4ª coluna (índice 3)
			// Usamos TextContent para pegar o texto bruto, lidando melhor com as <div> internas
			text, err := row.Locator("td").Nth(3).TextContent()
			if err != nil {
				continue
			}

			tipoLinha := strings.ToUpper(strings.TrimSpace(text))

			// Comparação flexível: o seu 'DA' estava limpo, mas o 'AD' tinha um div antes.
			// O TrimSpace do Go resolve a maioria, mas se falhar, o EqualFold é mais robusto.
			if tipoLinha == tipoParam {
				utils.Logger.Printf("Linha %d corresponde ao tipo: %s", i, tipoLinha)
				targetRow = row
				break
			}
		}

		if targetRow == nil {
			return "", fmt.Errorf("não encontrou linha com tipo %s entre as %d disponíveis", tipoParam, len(rows))
		}

		// 4. Clique no botão de detalhes dentro da linha encontrada
		btnDetalhes := targetRow.Locator("img[title='Detalhes do Empenho']")

		utils.Logger.Println("Clicando em Detalhes...")
		if err := btnDetalhes.Click(playwright.LocatorClickOptions{
			Timeout: playwright.Float(15_000),
		}); err != nil {
			return "", fmt.Errorf("falha ao clicar no botão da linha: %w", err)
		}

		utils.Logger.Println("Clique realizado com sucesso")
	}
	page.WaitForTimeout(4000)

	// ===============================
	// Captura de Cookies de Sessão
	// ===============================
	utils.Logger.Println("Extraindo cookies da sessão...")

	// Captura todos os cookies do contexto atual
	cookies, err := context.Cookies()
	if err != nil {
		return "", fmt.Errorf("erro ao capturar cookies: %w", err)
	}

	// Criando um mapa para você salvar e usar as variáveis como quiser
	cookieMap := make(map[string]string)
	for _, c := range cookies {
		cookieMap[c.Name] = c.Value
	}

	// Identificando o cookie de sessão comum em portais DevExpress (ASP.NET)
	if sessao, ok := cookieMap["ASP.NET_SessionId"]; ok {
		return sessao, nil

	} else {
		err = fmt.Errorf("⚠️ Cookie ASP.NET_SessionId não encontrado. Verifique os nomes disponíveis no mapa.")
		return "", err
	}

}
