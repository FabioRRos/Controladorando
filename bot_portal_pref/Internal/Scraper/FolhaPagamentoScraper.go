package scraper

import (
	"os"
	"path/filepath"
	"time"

	"github.com/playwright-community/playwright-go"
)

// BaixarCSVFolhas baixa tanto a mensal quanto a rescisão na mesma viagem
func BaixarCSVFolhas(pastaDestino string, nomeArquivoMensal string, nomeArquivoRescisao string) (string, string, error) {
	pw, err := playwright.Run()
	if err != nil {
		return "", "", err
	}
	defer pw.Stop()

	browser, err := pw.Chromium.Launch(playwright.BrowserTypeLaunchOptions{
		Headless: playwright.Bool(true),
	})
	if err != nil {
		return "", "", err
	}
	defer browser.Close()

	page, err := browser.NewPage()
	if err != nil {
		return "", "", err
	}

	page.SetDefaultTimeout(60000)

	if _, err = page.Goto("http://prefeiturapirajui1.ddns.net:8079/Transparencia/"); err != nil {
		return "", "", err
	}

	frame := page.FrameLocator("iframe, frame").First()

	// 1. Entra na tela de Servidores
	if err = frame.Locator("#btnAtalhoServidores").Click(); err != nil {
		return "", "", err
	}
	time.Sleep(3 * time.Second)

	os.MkdirAll(pastaDestino, os.ModePerm)

	// ================= FOLHA MENSAL =================
	// O rbListagemServidoresAtivos já vem checado, então vamos direto para o tipo:
	if err = frame.Locator("#rbTipoRefMensal").Click(); err != nil {
		return "", "", err
	}
	time.Sleep(4 * time.Second) // PostBack

	downloadMensal, err := page.ExpectDownload(func() error {
		return frame.Locator("#btnExportarCSV").Click()
	})
	if err != nil {
		return "", "", err
	}

	caminhoMensal := filepath.Join(pastaDestino, nomeArquivoMensal)
	if err = downloadMensal.SaveAs(caminhoMensal); err != nil {
		return "", "", err
	}

	// ================= FOLHA DE RESCISÃO =================
	if err = frame.Locator("#rbRecisao").Click(); err != nil {
		return "", "", err
	}
	time.Sleep(4 * time.Second) // PostBack

	downloadRescisao, err := page.ExpectDownload(func() error {
		return frame.Locator("#btnExportarCSV").Click()
	})
	if err != nil {
		return "", "", err
	}

	caminhoRescisao := filepath.Join(pastaDestino, nomeArquivoRescisao)
	if err = downloadRescisao.SaveAs(caminhoRescisao); err != nil {
		return "", "", err
	}

	return caminhoMensal, caminhoRescisao, nil
}
