package services

import (
	entity "bot_portal_pref/Models/Entity"
	"encoding/csv"
	"os"
	"strings"

	"golang.org/x/text/encoding/charmap"
)

func ConverterCSVCargosSalariosParaEntidade(caminho string) ([]entity.CargosSalarios, error) {
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

	var lista []entity.CargosSalarios

	for i, linha := range linhas {
		// Pula o cabeçalho ou linhas muito curtas
		if i == 0 || len(linha) < 6 {
			continue
		}

		//Pula as "linhas fantasmas" onde a coluna ID (índice 1) está vazia
		if strings.TrimSpace(linha[1]) == "" {
			continue
		}

		c := entity.CargosSalarios{
			PlanoCargo: strings.TrimSpace(linha[0]),
			CargoId:    converterInt(linha[1]),
			Cargo:      strings.TrimSpace(linha[2]),
			Referencia: strings.TrimSpace(linha[3]),
			Valor:      converterMoeda(linha[4]),
			Codigo:     strings.TrimSpace(linha[5]),
		}
		lista = append(lista, c)
	}
	return lista, nil
}
