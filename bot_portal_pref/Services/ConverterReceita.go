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

func ConverterCSVParaEntidade(caminhoArquivo string) ([]entity.Receitas, error) {
	file, err := os.Open(caminhoArquivo)
	if err != nil {
		return nil, fmt.Errorf("erro ao abrir o arquivo CSV: %w", err)
	}
	defer file.Close()

	// A MÁGICA ACONTECE AQUI:
	// Traduzimos o arquivo de ISO-8859-1 (padrão Brasil antigo) para UTF-8
	decodificador := charmap.ISO8859_1.NewDecoder().Reader(file)

	// Agora passamos o decodificador para o CSV Reader, em vez do arquivo cru
	reader := csv.NewReader(decodificador)
	reader.Comma = ';'
	reader.LazyQuotes = true

	linhas, err := reader.ReadAll()
	if err != nil {
		return nil, fmt.Errorf("erro ao ler linhas do CSV: %w", err)
	}

	var listaReceitas []entity.Receitas

	for i, linha := range linhas {
		if i == 0 {
			continue // Pula a linha de cabeçalho
		}

		// Segurança extra: ignora linhas em branco no final do arquivo
		if len(linha) < 9 {
			continue
		}

		receita := entity.Receitas{
			Codigo:             strings.TrimSpace(linha[0]),
			Especificacao:      strings.TrimSpace(linha[1]),
			CodAplicacao:       strings.TrimSpace(linha[2]),
			FonteSTN:           strings.TrimSpace(linha[3]),
			FonteRecurco:       strings.TrimSpace(linha[4]), // FonteRecurco conforme sua struct
			PrevisaoInicial:    converterMoeda(linha[5]),
			PrevisaoAtualizada: converterMoeda(linha[6]),
			ArrecadacaoPeriodo: converterMoeda(linha[7]),
			ArrecadacaoTotal:   converterMoeda(linha[8]),
		}

		listaReceitas = append(listaReceitas, receita)
	}

	return listaReceitas, nil
}

func converterMoeda(valor string) float64 {
	valorLimpo := strings.TrimSpace(valor)
	if valorLimpo == "" || valorLimpo == "-" {
		return 0.0
	}
	valorLimpo = strings.ReplaceAll(valorLimpo, ".", "")
	valorLimpo = strings.ReplaceAll(valorLimpo, ",", ".")

	resultado, _ := strconv.ParseFloat(valorLimpo, 64)
	return resultado
}
