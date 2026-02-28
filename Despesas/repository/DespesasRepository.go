package repository

import (
	"database/sql"
	"despesas/model/entity"
	"fmt"

	_ "github.com/lib/pq"
)

func SalvarDespesas(lista []entity.Despesas, connStr string) error {
	// 1. Abrir conexão
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		return fmt.Errorf("erro ao conectar no banco: %w", err)
	}
	defer db.Close()

	tx, err := db.Begin()
	if err != nil {
		return err
	}

	query := `INSERT INTO despesas (
		empenho, tipo, no_ficha, data, cod_forn, nome_fornecedor, cpf_cnpj,
		dotacao, alteracao_dotacao, dotacao_atual, valor_anulado, reforco,
		valor_empenhado, valor_liquidado, valor_pago, empenhado_ate_hoje,
		liquidado_ate_hoje, pago_ate_hoje, local, funcional, funcao,
		nome_funcao, subfuncao, nome_subfuncao, cod_aplicacao,
		descricao_cod_aplicacao, natureza, nome_natureza, fonte,
		fonte_recurso, cod_fonte, codigo_fonte, fonte_stn,
		nome_fonte_stn, proc_licitatorio, modalidade
	) VALUES (
		$1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, $15,
		$16, $17, $18, $19, $20, $21, $22, $23, $24, $25, $26, $27, $28,
		$29, $30, $31, $32, $33, $34, $35, $36
	)`

	stmt, err := tx.Prepare(query)
	if err != nil {
		tx.Rollback()
		return err
	}
	defer stmt.Close()

	for _, d := range lista {
		var dataParaBanco interface{}
		dataParaBanco = d.Data
		if d.Data == "" {
			dataParaBanco = nil
		}
		_, err := stmt.Exec(
			d.Empenho, d.Tipo, d.NoFicha, dataParaBanco, d.CodForn, d.NomeFornecedor, d.CpfCnpj,
			d.Dotacao, d.AlteracaoDotacao, d.DotacaoAtual, d.ValorAnulado, d.Reforco,
			d.ValorEmpenhado, d.ValorLiquidado, d.ValorPago, d.EmpenhadoAteHoje,
			d.LiquidadoAteHoje, d.PagoAteHoje, d.Local, d.Funcional, d.Funcao,
			d.NomeFuncao, d.Subfuncao, d.NomeSubfuncao, d.CodAplicacao,
			d.DescricaoCodAplicacao, d.Natureza, d.NomeNatureza, d.Fonte,
			d.FonteRecurso, d.CodFonte, d.CodigoFonte, d.FonteSTN,
			d.NomeFonteSTN, d.ProcLicitatorio, d.Modalidade,
		)
		if err != nil {
			tx.Rollback() // Se um der erro, cancela tudo
			return fmt.Errorf("erro ao inserir empenho %d: %w", d.Empenho, err)
		}
	}

	// 5. Finalizar transação
	return tx.Commit()
}
