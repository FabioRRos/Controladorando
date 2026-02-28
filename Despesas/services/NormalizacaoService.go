package services

import (
	"fmt"
	"strconv"
	"strings"
)

// converterInt transforma texto vazio em 0 de forma segura
func converterInt(valor string) int {
	v := strings.TrimSpace(valor)
	if v == "" || v == "-" {
		return 0
	}
	resultado, _ := strconv.Atoi(v)
	return resultado
}

func converterDataParaSQL(dataBR string) string {
	partes := strings.Split(strings.TrimSpace(dataBR), "/")
	if len(partes) == 3 {
		return fmt.Sprintf("%s-%s-%s", partes[2], partes[1], partes[0])
	}
	return ""
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
