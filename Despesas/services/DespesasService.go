package services

import (
	"despesas/model/entity"
	"encoding/csv"
	"fmt"
	"os"
	"strings"

	"golang.org/x/text/encoding/charmap"
)

func ConverterCSVDespesasParaEntidade(caminhoArquivo string) ([]entity.Despesas, error) {

	file, err := os.Open(caminhoArquivo)

	if err != nil {

		return nil, fmt.Errorf("erro ao abrir arquivo: %w", err)

	}

	defer file.Close()

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

			continue

		}

		if strings.TrimSpace(linha[0]) == "" {
			continue
		}

		if len(linha) < 36 {

			continue // Ignora linhas em branco ou quebradas no fim do arquivo

		}

		despesa := entity.Despesas{

			Empenho: converterInt(linha[0]),

			Tipo: strings.TrimSpace(linha[1]),

			NoFicha: converterInt(linha[2]),

			Data: converterDataParaSQL(linha[3]), // Converte DD/MM/YYYY para YYYY-MM-DD

			CodForn: converterInt(linha[4]),

			NomeFornecedor: strings.TrimSpace(linha[5]),

			CpfCnpj: strings.TrimSpace(linha[6]),

			Dotacao: converterMoeda(linha[7]),

			AlteracaoDotacao: converterMoeda(linha[8]),

			DotacaoAtual: converterMoeda(linha[9]),

			ValorAnulado: converterMoeda(linha[10]),

			Reforco: converterMoeda(linha[11]),

			ValorEmpenhado: converterMoeda(linha[12]),

			ValorLiquidado: converterMoeda(linha[13]),

			ValorPago: converterMoeda(linha[14]),

			EmpenhadoAteHoje: converterMoeda(linha[15]),

			LiquidadoAteHoje: converterMoeda(linha[16]),

			PagoAteHoje: converterMoeda(linha[17]),

			Local: strings.TrimSpace(linha[18]),

			Funcional: strings.TrimSpace(linha[19]),

			Funcao: converterInt(linha[20]),

			NomeFuncao: strings.TrimSpace(linha[21]),

			Subfuncao: converterInt(linha[22]),

			NomeSubfuncao: strings.TrimSpace(linha[23]),

			CodAplicacao: strings.TrimSpace(linha[24]),

			DescricaoCodAplicacao: strings.TrimSpace(linha[25]),

			Natureza: strings.TrimSpace(linha[26]),

			NomeNatureza: strings.TrimSpace(linha[27]),

			Fonte: converterInt(linha[28]),

			FonteRecurso: strings.TrimSpace(linha[29]),

			CodFonte: converterInt(linha[30]),

			CodigoFonte: strings.TrimSpace(linha[31]),

			FonteSTN: strings.TrimSpace(linha[32]),

			NomeFonteSTN: strings.TrimSpace(linha[33]),

			ProcLicitatorio: strings.TrimSpace(linha[34]),

			Modalidade: strings.TrimSpace(linha[35]),
		}

		lista = append(lista, despesa)

	}

	return lista, nil

}
