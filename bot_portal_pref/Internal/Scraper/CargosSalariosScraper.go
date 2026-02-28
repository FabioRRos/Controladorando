package scraper

import (
	"fmt"
	"os"
	"path/filepath"
	"time"

	"github.com/playwright-community/playwright-go"
)

func BaixarCSVCargosSalarios(pastaDestino string, nomeArquivo string) (string, error) {

	hoje := time.Now()
	dia := hoje.Day()
	if dia != 15 {
		return "", fmt.Errorf("Ainda não é dia 15")
	}

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

	// 1. Clica no menu de Servidores
	if err = frame.Locator("#btnAtalhoServidores").Click(); err != nil {
		return "", err
	}
	time.Sleep(3 * time.Second) // Aguarda a tela de servidores abrir

	// 2. Clica no RadioBox de Cargos e Salários
	if err = frame.Locator("#rbListagemCargoSalario").Click(); err != nil {
		return "", err
	}

	os.MkdirAll(pastaDestino, os.ModePerm)

	// 3. Exporta o CSV
	download, err := page.ExpectDownload(func() error {
		return frame.Locator("#btnExportarCSVCargoSalario").Click()
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
