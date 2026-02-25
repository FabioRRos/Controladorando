package services

import (
	entity "bot_portal_pref/Models/Entity"
	"encoding/csv"
	"os"
	"strings"

	"golang.org/x/text/encoding/charmap"
)

func ConverterCSVLicitacoesParaEntidade(caminho string) ([]entity.Licitacoes, error) {
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

	var lista []entity.Licitacoes

	for i, linha := range linhas {
		if i == 0 || len(linha) < 19 {
			continue
		}

		lic := entity.Licitacoes{
			ProcLicitatorio:     strings.TrimSpace(linha[0]),
			ProcAdministrativo:  strings.TrimSpace(linha[1]),
			Modalidade:          strings.TrimSpace(linha[2]),
			Exercicio:           converterInt(linha[3]),
			NumMod:              converterInt(linha[4]),
			Situacao:            strings.TrimSpace(linha[5]),
			DataAbertPropost:    converterDataParaSQL(linha[6]),
			HoraAbertPropost:    strings.TrimSpace(linha[7]),
			ValorPrevisto:       converterMoeda(linha[8]),
			ValorTotalLicitacao: converterMoeda(linha[9]),
			Objeto:              strings.TrimSpace(linha[10]),
			DataEdital:          converterDataParaSQL(linha[11]),
			DataEncerramento:    converterDataParaSQL(linha[12]),
			Carona:              strings.TrimSpace(linha[13]),
			RegPreco:            strings.TrimSpace(linha[14]),
			PrazoEntregaInicio:  strings.TrimSpace(linha[15]),
			ArtigoInciso:        strings.TrimSpace(linha[16]),
			DataInicioProposta:  converterDataParaSQL(linha[17]),
			DataFimProposta:     converterDataParaSQL(linha[18]),
		}
		lista = append(lista, lic)
	}
	return lista, nil
}
