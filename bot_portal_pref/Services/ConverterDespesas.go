package services

import (
	entity "bot_portal_pref/Models/Entity"
	"encoding/csv"
	"fmt"
	"os"
	"strconv"
	"strings"

	"golang.org/x/text/encoding/charmap"
)

func ConverterCSVDespesasParaEntidade(caminhoArquivo string) ([]entity.Despesas, error) {
	file, err := os.Open(caminhoArquivo)
	if err != nil {
		return nil, fmt.Errorf("erro ao abrir arquivo: %w", err)
	}
	defer file.Close()

	// A mágica anti-erro de encoding (traduz de ISO-8859-1 para UTF-8)
	decodificador := charmap.ISO8859_1.NewDecoder().Reader(file)
	reader := csv.NewReader(decodificador)
	reader.Comma = ';'
	reader.LazyQuotes = true

	linhas, err := reader.ReadAll()
	if err != nil {
		return nil, fmt.Errorf("erro na leitura do CSV: %w", err)
	}

	var lista []entity.Despesas

	for i, linha := range linhas {
		if i == 0 {
			continue // Pula cabeçalho
		}

		if len(linha) < 36 {
			continue // Ignora linhas em branco ou quebradas no fim do arquivo
		}

		despesa := entity.Despesas{
			Empenho:               converterInt(linha[0]),
			Tipo:                  strings.TrimSpace(linha[1]),
			NoFicha:               converterInt(linha[2]),
			Data:                  converterDataParaSQL(linha[3]), // Converte DD/MM/YYYY para YYYY-MM-DD
			CodForn:               converterInt(linha[4]),
			NomeFornecedor:        strings.TrimSpace(linha[5]),
			CpfCnpj:               strings.TrimSpace(linha[6]),
			Dotacao:               converterMoeda(linha[7]),
			AlteracaoDotacao:      converterMoeda(linha[8]),
			DotacaoAtual:          converterMoeda(linha[9]),
			ValorAnulado:          converterMoeda(linha[10]),
			Reforco:               converterMoeda(linha[11]),
			ValorEmpenhado:        converterMoeda(linha[12]),
			ValorLiquidado:        converterMoeda(linha[13]),
			ValorPago:             converterMoeda(linha[14]),
			EmpenhadoAteHoje:      converterMoeda(linha[15]),
			LiquidadoAteHoje:      converterMoeda(linha[16]),
			PagoAteHoje:           converterMoeda(linha[17]),
			Local:                 strings.TrimSpace(linha[18]),
			Funcional:             strings.TrimSpace(linha[19]),
			Funcao:                converterInt(linha[20]),
			NomeFuncao:            strings.TrimSpace(linha[21]),
			Subfuncao:             converterInt(linha[22]),
			NomeSubfuncao:         strings.TrimSpace(linha[23]),
			CodAplicacao:          strings.TrimSpace(linha[24]),
			DescricaoCodAplicacao: strings.TrimSpace(linha[25]),
			Natureza:              strings.TrimSpace(linha[26]),
			NomeNatureza:          strings.TrimSpace(linha[27]),
			Fonte:                 converterInt(linha[28]),
			FonteRecurso:          strings.TrimSpace(linha[29]),
			CodFonte:              converterInt(linha[30]),
			CodigoFonte:           strings.TrimSpace(linha[31]),
			FonteSTN:              strings.TrimSpace(linha[32]),
			NomeFonteSTN:          strings.TrimSpace(linha[33]),
			ProcLicitatorio:       strings.TrimSpace(linha[34]),
			Modalidade:            strings.TrimSpace(linha[35]),
		}
		lista = append(lista, despesa)
	}

	return lista, nil
}

// converterInt transforma texto vazio em 0 de forma segura
func converterInt(valor string) int {
	v := strings.TrimSpace(valor)
	if v == "" || v == "-" {
		return 0
	}
	resultado, _ := strconv.Atoi(v)
	return resultado
}

// converterDataParaSQL converte "05/01/2026" para "2026-01-05" que o PostgreSQL exige
func converterDataParaSQL(dataBR string) string {
	partes := strings.Split(strings.TrimSpace(dataBR), "/")
	if len(partes) == 3 {
		return fmt.Sprintf("%s-%s-%s", partes[2], partes[1], partes[0])
	}
	return ""
}

// Se você já tem a converterMoeda neste arquivo, não precisa copiar de novo!
