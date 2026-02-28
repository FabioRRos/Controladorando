package automacao

import (
	"despesas/model/entity"
	"fmt"
	"log"
	"path/filepath"
	"time"

	"github.com/playwright-community/playwright-go"
)

func BaixarDespesas() {
	// 1. Iniciar Playwright
	pw, err := playwright.Run()
	if err != nil {
		log.Fatalf("Não foi possível iniciar o Playwright: %v", err)
	}
	defer pw.Stop()

	// 2. Lançar o Browser
	browser, err := pw.Chromium.Launch(playwright.BrowserTypeLaunchOptions{
		Headless: playwright.Bool(false), // Mude para true depois
	})
	if err != nil {
		log.Fatalf("Não foi possível lançar o browser: %v", err)
	}
	defer browser.Close()

	page, err := browser.NewPage()
	if err != nil {
		log.Fatalf("Não foi possível criar a página: %v", err)
	}

	// Variáveis
	varlink := entity.LinkDespesas()
	datainicial := "01/01/2026"
	datafinal := "26/02/2026"

	// 3. Navegação
	fmt.Println("Navegando para o portal...")
	if _, err = page.Goto(varlink); err != nil {
		log.Fatalf("Erro ao navegar: %v", err)
	}
	frame := page.FrameLocator("iframe, frame").First()

	// 4. Interação (O Playwright espera o elemento aparecer automaticamente)
	fmt.Println("Clicando no atalho de despesas...")
	if err = frame.Locator("#btnAtalhoDespesa").Click(); err != nil {
		log.Fatalf("Erro ao clicar no atalho: %v", err)
	}

	// Preencher datas (usando os IDs que você mapeou)
	// O Fill limpa o campo e digita automaticamente
	if err = frame.Locator("#datDataInicial_I").Fill(datainicial); err != nil {
		log.Fatalf("Erro ao preencher data inicial: %v", err)
	}

	if err = frame.Locator("#datDataFinal_I").Fill(datafinal); err != nil {
		log.Fatalf("Erro ao preencher data final: %v", err)
	}

	// 5. Capturar o Download
	// Em vez de Sleep, o Playwright "espera" pelo evento de download
	fmt.Println("Iniciando exportação...")

	download, err := page.ExpectDownload(func() error {
		return frame.Locator("#btnExportarCSV").Click()
	})
	if err != nil {
		log.Fatalf("Erro ao aguardar download: %v", err)
	}

	caminhoCompleto := filepath.Join("temp", "despesas.csv")
	if err = download.SaveAs(caminhoCompleto); err != nil {
		log.Fatalf("erro ao salvar o arquivo CSV: %v", err)
	}

	fmt.Printf("Sucesso! Arquivo de despesas salvo em: %s\n", caminhoCompleto)

	// Mais um sleep pra você ver o final antes de fechar o Chrome
	time.Sleep(2 * time.Second)

}
