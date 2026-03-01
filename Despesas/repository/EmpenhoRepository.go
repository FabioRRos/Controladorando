package repository

import (
	"context"
	"database/sql"
	"despesas/model/entity"
	"fmt"
	"time"

	_ "github.com/lib/pq"
)

func parseDateBR(s string) interface{} {
	if s == "" {
		return nil
	}
	t, err := time.Parse("02/01/2006", s)
	if err != nil {
		// se vier algo inválido, melhor NULL do que quebrar tudo
		return nil
	}
	return t
}

func InsertEmpenho(connString string, emp entity.Empenho) error {
	db, err := sql.Open("postgres", connString)
	if err != nil {
		return fmt.Errorf("erro abrindo conexão: %w", err)
	}
	defer db.Close()

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// Datas: converte "" -> NULL e "dd/mm/yyyy" -> time.Time
	dataDB := parseDateBR(emp.Data)
	vigIniDB := parseDateBR(emp.VigenciaInicio)
	vigFimDB := parseDateBR(emp.VigenciaFim)
	vencDB := parseDateBR(emp.ContratoVencimentoAtual)

	query := `
	INSERT INTO empenhos (
		exercicio, data, numero_empenho, tipo, favorecido, cpf_cnpj, valor,
		processo_contratacao, tipo_licitacao, numero_licitacao,
		orgao, unidade_orcamentaria, projeto_atividade, vinculo_orcamentario,
		grupo_fonte, codigo_fonte, elemento, historico, poder, funcao,
		subfuncao, programa, fonro, fonte_stn, categoria_economica,
		grupo_natureza, modalidade_aplicacao, desdobro, natureza,
		numero_contrato, contrato_num_detalhado, vigencia_inicio,
		vigencia_fim, termo, contrato_adit_id, termo_resgatado,
		numero_convenio, ano_convenio, tipo_fundamento, inciso,
		contrato_vencimento_atual, ficha
	)
	VALUES (
		$1,$2,$3,$4,$5,$6,$7,
		$8,$9,$10,
		$11,$12,$13,$14,
		$15,$16,$17,$18,$19,$20,
		$21,$22,$23,$24,$25,
		$26,$27,$28,$29,
		$30,$31,$32,
		$33,$34,$35,$36,
		$37,$38,$39,$40,
		$41,$42
	)
	`

	_, err = db.ExecContext(ctx, query,
		emp.Exercicio,
		dataDB, // <- aqui
		emp.NumeroEmpenho,
		emp.Tipo,
		emp.Favorecido,
		emp.CpfCnpj,
		emp.Valor,
		emp.ProcessoContratacao,
		emp.TipoLicitacao,
		emp.NumeroLicitacao,
		emp.Orgao,
		emp.UnidadeOrcamentaria,
		emp.ProjetoAtividade,
		emp.VinculoOrcamentario,
		emp.GrupoFonte,
		emp.CodigoFonte,
		emp.Elemento,
		emp.Historico,
		emp.Poder,
		emp.Funcao,
		emp.Subfuncao,
		emp.Programa,
		emp.Fonro,
		emp.FonteSTN,
		emp.CategoriaEconomica,
		emp.GrupoNatureza,
		emp.ModalidadeAplicacao,
		emp.Desdobro,
		emp.Natureza,
		emp.NumeroContrato,
		emp.ContratoNumDetalhado,
		vigIniDB, // <- aqui
		vigFimDB, // <- aqui
		emp.Termo,
		emp.ContratoAditId,
		emp.TermoResgatado,
		emp.NumeroConvenio,
		emp.AnoConvenio,
		emp.TipoFundamento,
		emp.Inciso,
		vencDB, // <- aqui
		emp.Ficha,
	)
	if err != nil {
		return fmt.Errorf("erro no insert: %w", err)
	}

	return nil
}
