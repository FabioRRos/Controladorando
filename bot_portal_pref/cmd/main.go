package main

import (
	scraper "bot_portal_pref/Internal/Scraper"
	"bot_portal_pref/Internal/database"
	services "bot_portal_pref/Services"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"sync" // Importação necessária para o WaitGroup
)

func main() {
	fmt.Println("🚀 Iniciando o robô concorrente (Receitas e Despesas simultâneas)...")

	pastaDestino := "fileTemp"

	// O WaitGroup é o nosso maestro: ele garante que o programa não termine
	// antes das duas goroutines terminarem seus trabalhos.
	var wg sync.WaitGroup

	wg.Add(8)

	go processarReceitas(&wg, pastaDestino)
	go processarDespesas(&wg, pastaDestino)
	go processarLicitacoes(&wg, pastaDestino)
	go processarContratos(&wg, pastaDestino)
	go processarCargosSalarios(&wg, pastaDestino)
	go processarFolhas(&wg, pastaDestino)
	go processarAPI_TCE(&wg)
	go processarAPI_TCE_Receitas(&wg)
	wg.Wait()

	fmt.Println("✅ Processo 100% finalizado com sucesso nas duas frentes!")
}

// Separamos o fluxo de Receitas em uma função isolada
func processarReceitas(wg *sync.WaitGroup, pastaDestino string) {
	// O defer garante que, independente de dar erro ou sucesso,
	// ele avisa o maestro que essa tarefa terminou.
	defer wg.Done()

	prefixo := "[RECEITAS]"
	fmt.Printf("%s Iniciando rotina...\n", prefixo)

	err := scraper.BaixarCSV(pastaDestino)
	if err != nil {
		log.Printf("%s ❌ Erro ao baixar o CSV: %v\n", prefixo, err)
		return // Retorna para abortar só essa rotina, sem matar o programa principal
	}

	caminhoCompleto := filepath.Join(pastaDestino, "Portal Transparencia Receitas Acumuladas - Exercício 2026.csv")

	fmt.Printf("%s Lendo o CSV e convertendo para entidades...\n", prefixo)
	listaReceitas, err := services.ConverterCSVParaEntidade(caminhoCompleto)
	if err != nil {
		log.Printf("%s ❌ Erro ao converter o CSV: %v\n", prefixo, err)
		return
	}

	fmt.Printf("%s Salvando no PostgreSQL...\n", prefixo)
	err = database.SalvarReceitas(listaReceitas)
	if err != nil {
		log.Printf("%s ❌ Erro ao importar para o banco: %v\n", prefixo, err)
		return
	}

	// Faxina
	err = os.Remove(caminhoCompleto)
	if err != nil {
		fmt.Printf("%s Aviso: Não foi possível deletar o arquivo: %v\n", prefixo, err)
	}

	fmt.Printf("%s Sucesso! Concluído com êxito.\n", prefixo)
}

// Separamos o fluxo de Despesas na outra função isolada
func processarDespesas(wg *sync.WaitGroup, pastaDestino string) {
	defer wg.Done()

	prefixo := "[DESPESAS]"
	fmt.Printf("%s Iniciando rotina...\n", prefixo)

	nomeDoArquivo := "exportacao_despesas.csv"
	caminhoCompleto := filepath.Join(pastaDestino, nomeDoArquivo)

	_, err := scraper.BaixarCSVDespesas(pastaDestino, nomeDoArquivo)
	if err != nil {
		log.Printf("%s ❌ Erro ao baixar Despesas: %v\n", prefixo, err)
		return
	}

	fmt.Printf("%s Lendo CSV e limpando encoding...\n", prefixo)
	listaDespesas, err := services.ConverterCSVDespesasParaEntidade(caminhoCompleto)
	if err != nil {
		log.Printf("%s ❌ Erro ao converter o CSV: %v\n", prefixo, err)
		return
	}

	fmt.Printf("%s Salvando no PostgreSQL...\n", prefixo)
	err = database.SalvarDespesas(listaDespesas)
	if err != nil {
		log.Printf("%s ❌ Erro ao importar para o banco: %v\n", prefixo, err)
		return
	}

	// Faxina
	os.Remove(caminhoCompleto)

	fmt.Printf("%s Sucesso! %d despesas importadas.\n", prefixo, len(listaDespesas))
}

func processarLicitacoes(wg *sync.WaitGroup, pastaDestino string) {
	defer wg.Done()
	prefixo := "[LICITAÇÕES]"
	nomeDoArquivo := "exportacao_licitacoes.csv"

	fmt.Printf("%s Iniciando rotina...\n", prefixo)
	caminho, err := scraper.BaixarCSVLicitacoes(pastaDestino, nomeDoArquivo)
	if err != nil {
		fmt.Printf("%s ❌ Erro no download: %v\n", prefixo, err)
		return
	}

	lista, err := services.ConverterCSVLicitacoesParaEntidade(caminho)
	if err != nil {
		fmt.Printf("%s ❌ Erro na conversão: %v\n", prefixo, err)
		return
	}

	if err = database.SalvarLicitacoes(lista); err != nil {
		fmt.Printf("%s ❌ Erro no banco: %v\n", prefixo, err)
		return
	}

	os.Remove(caminho)
	fmt.Printf("%s Sucesso! %d importadas.\n", prefixo, len(lista))
}

func processarContratos(wg *sync.WaitGroup, pastaDestino string) {
	defer wg.Done()
	prefixo := "[CONTRATOS]"
	nomeDoArquivo := "exportacao_contratos.csv"

	fmt.Printf("%s Iniciando rotina...\n", prefixo)
	caminho, err := scraper.BaixarCSVContratos(pastaDestino, nomeDoArquivo)
	if err != nil {
		fmt.Printf("%s ❌ Erro no download: %v\n", prefixo, err)
		return
	}

	lista, err := services.ConverterCSVContratosParaEntidade(caminho)
	if err != nil {
		fmt.Printf("%s ❌ Erro na conversão: %v\n", prefixo, err)
		return
	}

	if err = database.SalvarContratos(lista); err != nil {
		fmt.Printf("%s ❌ Erro no banco: %v\n", prefixo, err)
		return
	}

	os.Remove(caminho)
	fmt.Printf("%s Sucesso! Tabela limpa e %d contratos importados.\n", prefixo, len(lista))
}

func processarCargosSalarios(wg *sync.WaitGroup, pastaDestino string) {
	defer wg.Done()
	prefixo := "[CARGOS]"
	nomeDoArquivo := "exportacao_cargos_salarios.csv"

	fmt.Printf("%s Iniciando rotina...\n", prefixo)
	caminho, err := scraper.BaixarCSVCargosSalarios(pastaDestino, nomeDoArquivo)
	if err != nil {
		fmt.Printf("%s ❌ Erro no download: %v\n", prefixo, err)
		return
	}

	lista, err := services.ConverterCSVCargosSalariosParaEntidade(caminho)
	if err != nil {
		fmt.Printf("%s ❌ Erro na conversão: %v\n", prefixo, err)
		return
	}

	if err = database.SalvarCargosSalarios(lista); err != nil {
		fmt.Printf("%s ❌ Erro no banco: %v\n", prefixo, err)
		return
	}

	os.Remove(caminho)
	fmt.Printf("%s Sucesso! Tabela limpa e %d cargos importados.\n", prefixo, len(lista))
}

func processarFolhas(wg *sync.WaitGroup, pastaDestino string) {
	defer wg.Done()
	prefixo := "[FOLHAS]"
	nomeMensal := "exportacao_folha_mensal.csv"
	nomeRescisao := "exportacao_folha_rescisao.csv"

	fmt.Printf("%s Iniciando rotina dupla (Mensal e Rescisão)...\n", prefixo)

	caminhoMensal, caminhoRescisao, err := scraper.BaixarCSVFolhas(pastaDestino, nomeMensal, nomeRescisao)
	if err != nil {
		fmt.Printf("%s ❌ Erro no download: %v\n", prefixo, err)
		return
	}

	// === Processa a Folha Mensal ===
	listaMensal, err := services.ConverterCSVFolhaParaEntidade(caminhoMensal)
	if err == nil {
		if err = database.SalvarFolha(listaMensal, "folha_mensal"); err != nil {
			fmt.Printf("%s ❌ Erro banco Mensal: %v\n", prefixo, err)
		} else {
			fmt.Printf("%s ✔️ %d registros salvos na folha_mensal.\n", prefixo, len(listaMensal))
		}
	}
	os.Remove(caminhoMensal)

	// === Processa a Folha de Rescisão ===
	listaRescisao, err := services.ConverterCSVFolhaParaEntidade(caminhoRescisao)
	if err == nil {
		if err = database.SalvarFolha(listaRescisao, "folha_rescisao"); err != nil {
			fmt.Printf("%s ❌ Erro banco Rescisão: %v\n", prefixo, err)
		} else {
			fmt.Printf("%s ✔️ %d registros salvos na folha_rescisao.\n", prefixo, len(listaRescisao))
		}
	}
	os.Remove(caminhoRescisao)
}
func processarAPI_TCE(wg *sync.WaitGroup) {
	defer wg.Done()
	prefixo := "[API TCE]"

	// Configurações da busca
	municipio := "pirajui"
	ano := 2025
	mes := 12

	fmt.Printf("%s Consultando API (Município: %s | %d/%d)...\n", prefixo, municipio, mes, ano)

	listaTce, err := services.BuscarDespesasTCE(municipio, ano, mes)
	if err != nil {
		fmt.Printf("%s ❌ Erro na consulta da API: %v\n", prefixo, err)
		return
	}

	if len(listaTce) == 0 {
		fmt.Printf("%s ⚠️ Nenhum dado retornado do TCE para este período.\n", prefixo)
		return
	}

	fmt.Printf("%s Salvando %d registros no PostgreSQL...\n", prefixo, len(listaTce))

	// Passamos o mês como string para fazer o DELETE condicional no banco
	mesStr := fmt.Sprintf("%02d", mes)
	if err = database.SalvarDespesasTCE(listaTce, mesStr); err != nil {
		fmt.Printf("%s ❌ Erro ao salvar dados da API: %v\n", prefixo, err)
		return
	}

	fmt.Printf("%s ✔️ Sucesso Absoluto! Despesas do TCE importadas.\n", prefixo)
}

func processarAPI_TCE_Receitas(wg *sync.WaitGroup) {
	defer wg.Done()
	prefixo := "[API TCE - RECEITAS]"

	municipio := "pirajui"
	ano := 2025
	mes := 12

	fmt.Printf("%s Consultando API...\n", prefixo)

	listaTce, err := services.BuscarReceitasTCE(municipio, ano, mes)
	if err != nil {
		fmt.Printf("%s ❌ Erro na consulta da API: %v\n", prefixo, err)
		return
	}

	if len(listaTce) == 0 {
		fmt.Printf("%s ⚠️ Nenhum dado retornado do TCE para este período.\n", prefixo)
		return
	}

	fmt.Printf("%s Salvando %d registros no PostgreSQL...\n", prefixo, len(listaTce))

	if err = database.SalvarReceitasTCE(listaTce, fmt.Sprintf("%d", mes)); err != nil {
		fmt.Printf("%s ❌ Erro ao salvar dados da API: %v\n", prefixo, err)
		return
	}

	fmt.Printf("%s ✔️ Sucesso Absoluto! Receitas do TCE importadas.\n", prefixo)
}
