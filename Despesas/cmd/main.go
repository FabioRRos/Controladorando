package main

import (
	"despesas/automacao"
	"despesas/repository"
	"despesas/services"
	"despesas/utils"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"
)

var strConn string = "postgres://pirajui:vKEP82XuP@ssw0rdMoreka@2129@localhost:5432/DespesaService?sslmode=disable"

func main() {
	if err := utils.InitLogger(); err != nil {
		log.Fatal(err)
	}

	utils.Logger.Println("Processo iniciado")

	DespesasInicio()

	EmpenhoInicio()

}

func EmpenhoInicio() {
	utils.Logger.Println("############# Chamando Função Empenho Inicio ############# ")

	services.ProcessarDespesas(strConn)
}

func DespesasInicio() {

	automacao.BaixarDespesas()
	utils.Logger.Println("############# Chamando Função DespensaInicio ############# ")

	diretorioatual, _ := os.Getwd()
	var caminho string
	if strings.HasSuffix(diretorioatual, "cmd") {
		caminho = filepath.Join(diretorioatual, "temp", "despesas.csv")
	} else {
		caminho = filepath.Join(diretorioatual, "cmd", "temp", "despesas.csv")
	}

	listaDespesas, err := services.ConverterCSVDespesasParaEntidade(caminho)
	if err != nil {
		log.Fatalf("Erro ao converter CSV: %v", err)
	}
	err = repository.SalvarDespesas(listaDespesas, strConn)
	if err != nil {
		log.Fatalf("Erro ao salvar no banco: %v", err)
	}

	err = os.Remove(caminho)
	if err != nil {
		fmt.Printf("Aviso: Não foi possível excluir o arquivo temporário: %v\n", err)
	} else {
		fmt.Println("Arquivo temporário removido com sucesso!")
	}
	fmt.Println("Processo concluído com sucesso!")

}
