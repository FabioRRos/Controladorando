package services

import (
	entity "bot_portal_pref/Models/Entity"
	"encoding/csv"
	"os"
	"strings"

	"golang.org/x/text/encoding/charmap"
)

func ConverterCSVContratosParaEntidade(caminho string) ([]entity.Contratos, error) {
	file, err := os.Open(caminho)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	decodificador := charmap.ISO8859_1.NewDecoder().Reader(file)
	reader := csv.NewReader(decodificador)
	reader.Comma = ';'
	reader.LazyQuotes = true

	linhas, err := reader.ReadAll()
	if err != nil {
		return nil, err
	}

	var lista []entity.Contratos

	for i, linha := range linhas {
		if i == 0 || len(linha) < 15 {
			continue
		}

		c := entity.Contratos{
			NumContrato:          strings.TrimSpace(linha[0]),
			NumDetalhadoContrato: strings.TrimSpace(linha[1]),
			NumModalidade:        strings.TrimSpace(linha[2]),
			Modalidade:           strings.TrimSpace(linha[3]),
			Exercicio:            converterInt(linha[4]),
			FundamentoLegal:      strings.TrimSpace(linha[5]),
			ProcLicitatorio:      strings.TrimSpace(linha[6]),
			CpfCnpjFornecedor:    strings.TrimSpace(linha[7]),
			Fornecedor:           strings.TrimSpace(linha[8]),
			Valor:                converterMoeda(linha[9]),
			VigenciaInicial:      converterDataParaSQL(linha[10]),
			VencimentoAtual:      converterDataParaSQL(linha[11]),
			Objeto:               strings.TrimSpace(linha[12]),
			Tipo:                 converterInt(linha[13]),
			ContratoRateio:       strings.TrimSpace(linha[14]),
		}
		lista = append(lista, c)
	}
	return lista, nil
}
