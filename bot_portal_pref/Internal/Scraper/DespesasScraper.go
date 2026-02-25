package scraper

import (
	"fmt"
	"os"
	"path/filepath"
	"time"

	"github.com/playwright-community/playwright-go"
)

func BaixarCSVDespesas(pastaDestino string, nomeArquivoCustomizado string) (string, error) {

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

	// 1. Criamos a página direto, sem usar o context (igualzinho nas Receitas)
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

	fmt.Println("Clicando em Despesas...")
	if err = frame.Locator("#btnAtalhoDespesa").Click(); err != nil {
		return "", fmt.Errorf("erro ao clicar em Despesas: %w", err)
	}

	time.Sleep(3 * time.Second)

	fmt.Println("Preenchendo as datas...")
	if err = frame.Locator("#datDataInicial_I").Fill("01/01/2026"); err != nil {
		return "", fmt.Errorf("erro ao preencher data inicial: %w", err)
	}
	if err = frame.Locator("#datDataFinal_I").Fill("31/01/2026"); err != nil {
		return "", fmt.Errorf("erro ao preencher data final: %w", err)
	}

	if err := os.MkdirAll(pastaDestino, os.ModePerm); err != nil {
		return "", fmt.Errorf("erro ao criar a pasta: %w", err)
	}

	fmt.Println("Aguardando a geração do CSV pelo servidor (com paciência de 60s)...")
	download, err := page.ExpectDownload(func() error {
		return frame.Locator("#btnExportarCSV").Click()
	})

	if err != nil {
		return "", fmt.Errorf("erro ao interceptar o download: %w", err)
	}

	caminhoCompleto := filepath.Join(pastaDestino, nomeArquivoCustomizado)
	if err = download.SaveAs(caminhoCompleto); err != nil {
		return "", fmt.Errorf("erro ao salvar o arquivo CSV: %w", err)
	}

	fmt.Printf("Sucesso! Arquivo de despesas salvo em: %s\n", caminhoCompleto)

	// Mais um sleep pra você ver o final antes de fechar o Chrome
	time.Sleep(2 * time.Second)

	return caminhoCompleto, nil
}
