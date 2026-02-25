package scraper

import (
	"os"
	"path/filepath"
	"time"

	"github.com/playwright-community/playwright-go"
)

func BaixarCSVContratos(pastaDestino string, nomeArquivo string) (string, error) {
	pw, err := playwright.Run()
	if err != nil {
		return "", err
	}
	defer pw.Stop()

	browser, err := pw.Chromium.Launch(playwright.BrowserTypeLaunchOptions{
		Headless: playwright.Bool(true),
	})
	if err != nil {
		return "", err
	}
	defer browser.Close()

	page, err := browser.NewPage()
	if err != nil {
		return "", err
	}

	page.SetDefaultTimeout(60000)

	if _, err = page.Goto("http://prefeiturapirajui1.ddns.net:8079/Transparencia/"); err != nil {
		return "", err
	}

	frame := page.FrameLocator("iframe, frame").First()

	// 1. Clica em Contratos
	if err = frame.Locator("#btnAtalhoContrato").Click(); err != nil {
		return "", err
	}

	// Aguarda a tabela carregar completamente antes de tentar baixar
	time.Sleep(5 * time.Second)

	os.MkdirAll(pastaDestino, os.ModePerm)

	// 2. Exporta direto
	download, err := page.ExpectDownload(func() error {
		return frame.Locator("#ASPxPageControl1_btnExportarCSV").Click()
	})
	if err != nil {
		return "", err
	}

	caminho := filepath.Join(pastaDestino, nomeArquivo)
	if err = download.SaveAs(caminho); err != nil {
		return "", err
	}

	return caminho, nil
}
