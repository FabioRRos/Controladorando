package scraper

import (
	"fmt"
	"os"
	"path/filepath"
	"time"

	"github.com/playwright-community/playwright-go"
)

func BaixarCSVLicitacoes(pastaDestino string, nomeArquivo string) (string, error) {
	pw, err := playwright.Run()
	if err != nil {
		return "", fmt.Errorf("erro ao iniciar o Playwright: %w", err)
	}
	defer pw.Stop()

	browser, err := pw.Chromium.Launch(playwright.BrowserTypeLaunchOptions{
		Headless: playwright.Bool(true),
	})
	if err != nil {
		return "", fmt.Errorf("erro ao lançar o navegador: %w", err)
	}
	defer browser.Close()

	page, err := browser.NewPage()
	if err != nil {
		return "", fmt.Errorf("erro ao criar a página: %w", err)
	}

	page.SetDefaultTimeout(60000)

	fmt.Println("Acessando o portal de transparência...")
	if _, err = page.Goto("http://prefeiturapirajui1.ddns.net:8079/Transparencia/"); err != nil {
		return "", fmt.Errorf("erro ao acessar a página: %w", err)
	}

	frame := page.FrameLocator("iframe, frame").First()

	// 1. Clica no botão de Licitações
	fmt.Println("Clicando em Licitações...")
	if err = frame.Locator("#btnAtalhoLicitacao").Click(); err != nil {
		return "", fmt.Errorf("erro ao clicar em Licitações: %w", err)
	}

	// Como é um grid que atualiza via AJAX quando você digita,
	// um pequeno sleep aqui garante que a tabela filtrou antes de você mandar exportar.
	time.Sleep(2 * time.Second)

	if err := os.MkdirAll(pastaDestino, os.ModePerm); err != nil {
		return "", fmt.Errorf("erro ao criar a pasta: %w", err)
	}

	// 3. Exporta o CSV
	fmt.Println("Aguardando a geração do CSV de Licitações...")
	download, err := page.ExpectDownload(func() error {
		return frame.Locator("#ASPxCallbackPanel1_ASPxPageControl1_btnExportarCSV").Click()
	})
	if err != nil {
		return "", fmt.Errorf("erro ao interceptar o download: %w", err)
	}

	caminhoCompleto := filepath.Join(pastaDestino, nomeArquivo)
	if err = download.SaveAs(caminhoCompleto); err != nil {
		return "", fmt.Errorf("erro ao salvar o arquivo CSV: %w", err)
	}

	return caminhoCompleto, nil
}
