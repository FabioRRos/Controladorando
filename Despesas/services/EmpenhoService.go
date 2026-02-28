package services

import (
	"despesas/automacao"
	"despesas/model/entity"
	"despesas/repository"
	"encoding/csv"
	"fmt"
	"os"
	"path/filepath"

	"golang.org/x/text/encoding/charmap"
)

var listEmpenho []entity.Empenho

func ProcessarDespesas(strConn string) error {

	Unicas(strConn)
	//Multiplas(strConn)

	return nil
}

func Unicas(strConn string) error {

	lista, err := repository.BuscarDespesasPendentesUnicas(strConn)
	if err != nil {
		return err
	}

	// Canal que transporta o caminho do arquivo baixado
	arquivosChan := make(chan string)

	// ===============================
	// Goroutine PRODUTORA (Download)
	// ===============================
	go func() {
		for _, k := range lista {

			// Faz download
			err := automacao.BaixarEmpenhos(k.Empenho, k.Tipo, false)
			if err != nil {
				fmt.Println("Erro no download:", err)
				continue
			}

			// Monta nome do arquivo
			nomeArquivo := fmt.Sprintf("%d_%s.csv", k.Empenho, k.Tipo)

			dir := `C:\Pirajui\Controladorando\Despesas\cmd\temp`
			caminhoFinal := filepath.Join(dir, nomeArquivo)

			// Envia para processamento
			arquivosChan <- caminhoFinal
		}

		// Fecha canal quando terminar
		close(arquivosChan)
	}()

	// ===============================
	// CONSUMIDOR (Processamento)
	// ===============================
	for caminho := range arquivosChan {

		empRet, err := ProcessarDetalheEmpenho(caminho)
		if err != nil {
			fmt.Println("Erro ao processar:", err)
			continue
		}

		listEmpenho = append(listEmpenho, empRet)
		fmt.Println("Estou com uma lista de empenho com total de", len(listEmpenho))
	}

	return nil
}

/*
func Multiplas(strConn string) {
	lista, err := repository.BuscarDespesasPendentesMultiplas(strConn)
	if err != nil {
		return err
	}

	return nill
}*/

func ProcessarDetalheEmpenho(caminho string) (entity.Empenho, error) {

	fmt.Println("Iniciando a leitura do arquivo....")

	file, err := os.Open(caminho)
	if err != nil {
		return entity.Empenho{}, fmt.Errorf("não foi possível abrir o arquivo: %w", err)
	}
	defer file.Close()

	decodificador := charmap.ISO8859_1.NewDecoder().Reader(file)

	reader := csv.NewReader(decodificador)
	reader.Comma = ';'
	reader.LazyQuotes = true

	linhas, err := reader.ReadAll()
	if err != nil {
		return entity.Empenho{}, fmt.Errorf("erro na leitura do CSV: %w", err)
	}

	if len(linhas) < 3 {
		return entity.Empenho{}, fmt.Errorf(
			"arquivo inválido: esperado ao menos 3 linhas, obtido %d",
			len(linhas),
		)
	}

	dados := linhas[2]

	if len(dados) < 42 {
		return entity.Empenho{}, fmt.Errorf(
			"linha inválida: esperado ao menos 42 colunas, obtido %d",
			len(dados),
		)
	}

	empenho := entity.Empenho{
		Exercicio:               converterInt(dados[0]),
		Data:                    converterDataParaSQL(dados[1]),
		NumeroEmpenho:           converterInt(dados[2]),
		Tipo:                    dados[3],
		Favorecido:              dados[4],
		CpfCnpj:                 dados[5],
		Valor:                   converterMoeda(dados[6]),
		ProcessoContratacao:     dados[7],
		TipoLicitacao:           dados[8],
		NumeroLicitacao:         dados[9],
		Orgao:                   dados[10],
		UnidadeOrcamentaria:     dados[11],
		ProjetoAtividade:        dados[12],
		VinculoOrcamentario:     dados[13],
		GrupoFonte:              dados[14],
		CodigoFonte:             dados[15],
		Elemento:                dados[16],
		Historico:               dados[17],
		Poder:                   dados[18],
		Funcao:                  dados[19],
		Subfuncao:               dados[20],
		Programa:                dados[21],
		Fonro:                   dados[22],
		FonteSTN:                dados[23],
		CategoriaEconomica:      dados[24],
		GrupoNatureza:           dados[25],
		ModalidadeAplicacao:     dados[26],
		Desdobro:                dados[27],
		Natureza:                dados[28],
		NumeroContrato:          dados[29],
		ContratoNumDetalhado:    dados[30],
		VigenciaInicio:          dados[31],
		VigenciaFim:             dados[32],
		Termo:                   dados[33],
		ContratoAditId:          dados[34],
		TermoResgatado:          dados[35],
		NumeroConvenio:          dados[36],
		AnoConvenio:             converterInt(dados[37]),
		TipoFundamento:          dados[38],
		Inciso:                  dados[39],
		ContratoVencimentoAtual: dados[40],
		Ficha:                   converterInt(dados[41]),
	}

	return empenho, nil
}
