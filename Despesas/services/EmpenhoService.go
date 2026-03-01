package services

import (
	"despesas/automacao"
	"despesas/client"
	"despesas/model/entity"
	"despesas/repository"
	"despesas/utils"
	"fmt"
	"time"
)

var listEmpenho []entity.Empenho

func ProcessarDespesas(strConn string) error {

	Unicas(strConn)
	//Multiplas(strConn)

	return nil
}

var hora = time.Now().Format("15:04:05")

func Unicas(strConn string) error {
	utils.Logger.Println("Buscando lista de despesas unicas")
	lista, err := repository.BuscarDespesasPendentesUnicas(strConn)

	fmt.Println("Iniciando empenho")
	if err != nil {
		return err
	}
	utils.Logger.Println("Iniciando varredura de itens não baixados.")
	for _, k := range lista {
		utils.Logger.Println("ID", k.Id, " - Despensa ", k.Empenho)

		start := hora

		utils.Logger.Println("Inicio:", start)
		cookie, err := automacao.BaixarEmpenhos(k.Empenho, k.Tipo, false)

		if err != nil {
			utils.Logger.Println("Erro no processamento:", err)
			continue
		}

		empenho, err := client.FetchEmpenho(cookie)
		if err != nil {
			utils.Logger.Println("Erro na requisição:", err)
			continue
		}

		err = repository.InsertEmpenho(strConn, empenho)
		if err != nil {
			utils.Logger.Println("Erro ao salvar empenho no banco:", err)
			continue
		}
		fmt.Println("ID:", k.Id)
		err = repository.MarcarDespesaComoProcessada(strConn, k.Id)
		if err != nil {
			utils.Logger.Println("Erro ao baixar despensa no banco:", err)
			continue
		}

		end := hora
		utils.Logger.Println("Termino:", end)
	}
	return nil
}

func Multiplas(strConn string) error {
	utils.Logger.Println("Buscando lista de despesas multiplas")
	lista, err := repository.BuscarDespesasPendentesMultiplas(strConn)

	fmt.Println("Iniciando empenho")
	if err != nil {
		return err
	}
	utils.Logger.Println("Iniciando varredura de itens não baixados.")
	for _, k := range lista {
		utils.Logger.Println("ID", k.Id, " - Despensa ", k.Empenho)

		start := hora

		utils.Logger.Println("Inicio:", start)
		cookie, err := automacao.BaixarEmpenhos(k.Empenho, k.Tipo, true)

		if err != nil {
			utils.Logger.Println("Erro no processamento:", err)
			continue
		}

		empenho, err := client.FetchEmpenho(cookie)
		if err != nil {
			utils.Logger.Println("Erro na requisição:", err)
			continue
		}

		err = repository.InsertEmpenho(strConn, empenho)
		if err != nil {
			utils.Logger.Println("Erro ao salvar empenho no banco:", err)
			continue
		}
		fmt.Println("ID:", k.Id)
		err = repository.MarcarDespesaComoProcessada(strConn, k.Id)
		if err != nil {
			utils.Logger.Println("Erro ao baixar despensa no banco:", err)
			continue
		}

		end := hora
		utils.Logger.Println("Termino:", end)
	}
	return nil
}
