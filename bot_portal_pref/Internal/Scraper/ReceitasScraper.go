package scraper

import (
	"fmt"
	"os"
	"path/filepath"
	"time"

	"github.com/playwright-community/playwright-go"
)

// Agora ela recebe a pasta e retorna um erro (se der ruim)
func BaixarCSV(pastaDestino string) error {

	// Inicia o Playwright
	pw, err := playwright.Run()
	if err != nil {
		return fmt.Errorf("erro ao iniciar o Playwright: %w", err)
	}
	defer pw.Stop()

	// Abre o Chromium. Headless = false deixa você ver a mágica acontecendo
	browser, err := pw.Chromium.Launch(playwright.BrowserTypeLaunchOptions{
		Headless: playwright.Bool(true),
	})
	if err != nil {
		return fmt.Errorf("erro ao lançar o navegador: %w", err)
	}
	defer browser.Close()

	page, err := browser.NewPage()
	if err != nil {
		return fmt.Errorf("erro ao criar a página: %w", err)
	}

	// 1. Acessa o site

	fmt.Println("Acessando o portal de transparência...")
	if _, err = page.Goto("http://prefeiturapirajui1.ddns.net:8079/Transparencia/"); err != nil {
		return fmt.Errorf("erro ao acessar a página: %w", err)
	}

	time.Sleep(2 * time.Second)
	fmt.Println("Clicando em Receitas...")
	frame := page.FrameLocator("iframe, frame").First()
	// 2. Clica em "Receitas"
	if err = frame.Locator("#btnAtalhoReceita").Click(); err != nil {
		return fmt.Errorf("erro ao clicar em Receitas: %w", err)
	}
	time.Sleep(2 * time.Second)
	// 3. Seleciona a data inicial e data final

	fmt.Println("Preenchendo as datas...")
	if err = frame.Locator("#datDataInicial_I").Fill("01/01/2026"); err != nil {
		return fmt.Errorf("erro ao preencher data inicial: %w", err)
	}
	if err = frame.Locator("#datDataFinal_I").Fill("31/01/2026"); err != nil {
		return fmt.Errorf("erro ao preencher data final: %w", err)
	}

	// Cria a pasta fileTemp
	if err := os.MkdirAll(pastaDestino, os.ModePerm); err != nil {
		return fmt.Errorf("erro ao criar a pasta: %w", err)
	}

	// 4 e 5. Clica em CSV e Salva na pasta
	fmt.Println("Aguardando o download do CSV...")
	download, err := page.ExpectDownload(func() error {
		return frame.Locator("#btnExportarCSV").Click()
	})
	if err != nil {
		return fmt.Errorf("erro ao interceptar o download: %w", err)
	}

	caminhoCompleto := filepath.Join(pastaDestino, download.SuggestedFilename())
	if err = download.SaveAs(caminhoCompleto); err != nil {
		return fmt.Errorf("erro ao salvar o arquivo CSV: %w", err)
	}

	fmt.Printf("Sucesso ABSOLUTO! Arquivo salvo em: %s\n", caminhoCompleto)

	// Um sleepzinho no final só pra você comemorar olhando pra tela antes de fechar
	time.Sleep(2 * time.Second)

	return nil
}
