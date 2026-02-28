package services

import (
	entity "bot_portal_pref/Models/Entity"
	"encoding/csv"
	"os"
	"strings"

	"golang.org/x/text/encoding/charmap"
)

func ConverterCSVFolhaParaEntidade(caminho string) ([]entity.FolhaPagamento, error) {
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

	var lista []entity.FolhaPagamento

	for i, linha := range linhas {
		if i == 0 || len(linha) < 16 {
			continue
		}

		// Ignora as linhas de cabeçalho agrupadora que vêm sem matrícula e sem nome
		if strings.TrimSpace(linha[6]) == "" && strings.TrimSpace(linha[3]) == "" {
			continue
		}

		f := entity.FolhaPagamento{
			Detalhe:                     strings.TrimSpace(linha[0]),
			Referencia:                  strings.TrimSpace(linha[1]),
			ReferenciaSalarial:          strings.TrimSpace(linha[2]),
			Nome:                        strings.TrimSpace(linha[3]),
			Divisao:                     strings.TrimSpace(linha[4]),
			Cargo:                       strings.TrimSpace(linha[5]),
			Matricula:                   converterInt(linha[6]),
			Proventos:                   converterMoeda(linha[7]),
			Descontos:                   converterMoeda(linha[8]),
			Liquido:                     converterMoeda(linha[9]),
			DataAdmissao:                converterDataParaSQL(linha[10]),
			DataDesligamento:            converterDataParaSQL(linha[11]),
			TipoRegime:                  strings.TrimSpace(linha[12]),
			SituacaoFuncional:           strings.TrimSpace(linha[13]),
			TipoContrato:                strings.TrimSpace(linha[14]),
			DataPrevistaTerminoContrato: converterDataParaSQL(linha[15]),
		}
		lista = append(lista, f)
	}
	return lista, nil
}
