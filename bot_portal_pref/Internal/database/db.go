package database

import (
	entity "bot_portal_pref/Models/Entity"
	"database/sql"
	"fmt"

	_ "github.com/lib/pq"
)

var strConexao string = "postgres://pirajui:vKEP82XuP@ssw0rdMoreka@2129@localhost:5432/Pirajui?sslmode=disable"

func SalvarReceitas(receitas []entity.Receitas) error {
	//	strConexao := "postgres://pirajui:vKEP82XuP@ssw0rdMoreka@2129@localhost:5432/Pirajui?sslmode=disable"

	db, err := sql.Open("postgres", strConexao)
	if err != nil {
		return fmt.Errorf("erro ao conectar no banco: %w", err)
	}
	defer db.Close()

	tx, err := db.Begin()
	if err != nil {
		return err
	}

	query := `
		INSERT INTO receitas (
			codigo, especificacao, cod_aplicacao, fonte_stn, fonte_recurso, 
			previsao_inicial, previsao_atualizada, arrecadacao_periodo, arrecadacao_total
		) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)`

	stmt, err := tx.Prepare(query)
	if err != nil {
		return err
	}
	defer stmt.Close()

	for _, r := range receitas {
		_, err = stmt.Exec(
			r.Codigo,
			r.Especificacao,
			r.CodAplicacao,
			r.FonteSTN,
			r.FonteRecurco,
			r.PrevisaoInicial,
			r.PrevisaoAtualizada,
			r.ArrecadacaoPeriodo,
			r.ArrecadacaoTotal,
		)

		if err != nil {
			tx.Rollback()
			return fmt.Errorf("erro ao inserir receita %s: %w", r.Codigo, err)
		}
	}

	if err = tx.Commit(); err != nil {
		return fmt.Errorf("erro ao salvar a transação no banco: %w", err)
	}

	return nil
}

func SalvarDespesas(despesas []entity.Despesas) error {
	//strConexao := "postgres://pirajui:vKEP82XuP@ssw0rdMoreka@2129@localhost:5432/Pirajui?sslmode=disable"

	db, err := sql.Open("postgres", strConexao)
	if err != nil {
		return fmt.Errorf("erro ao conectar no banco: %w", err)
	}
	defer db.Close()

	tx, err := db.Begin()
	if err != nil {
		return err
	}

	query := `
		INSERT INTO despesas (
			empenho, tipo, no_ficha, data_transacao, cod_forn, nome_fornecedor, cpf_cnpj,
			dotacao, alteracao_dotacao, dotacao_atual, valor_anulado, reforco, valor_empenhado,
			valor_liquidado, valor_pago, empenhado_ate_hoje, liquidado_ate_hoje, pago_ate_hoje,
			local_codigo, funcional, funcao_codigo, nome_funcao, subfuncao_codigo, nome_subfuncao,
			cod_aplicacao, descricao_cod_aplicacao, natureza, nome_natureza, fonte_codigo,
			fonte_recurso, cod_fonte_detalhado, codigo_fonte_desc, fonte_stn, nome_fonte_stn,
			proc_licitatorio, modalidade
		) VALUES (
			$1, $2, $3, $4, $5, $6, $7, $8, $9, $10,
			$11, $12, $13, $14, $15, $16, $17, $18, $19, $20,
			$21, $22, $23, $24, $25, $26, $27, $28, $29, $30,
			$31, $32, $33, $34, $35, $36
		)`

	stmt, err := tx.Prepare(query)
	if err != nil {
		return err
	}
	defer stmt.Close()

	for _, d := range despesas {
		// O banco espera um formato de data válido. Se vier vazio, passamos NULL.
		var data interface{} = d.Data
		if d.Data == "" {
			data = nil
		}

		_, err = stmt.Exec(
			d.Empenho, d.Tipo, d.NoFicha, data, d.CodForn, d.NomeFornecedor, d.CpfCnpj,
			d.Dotacao, d.AlteracaoDotacao, d.DotacaoAtual, d.ValorAnulado, d.Reforco, d.ValorEmpenhado,
			d.ValorLiquidado, d.ValorPago, d.EmpenhadoAteHoje, d.LiquidadoAteHoje, d.PagoAteHoje,
			d.Local, d.Funcional, d.Funcao, d.NomeFuncao, d.Subfuncao, d.NomeSubfuncao,
			d.CodAplicacao, d.DescricaoCodAplicacao, d.Natureza, d.NomeNatureza, d.Fonte,
			d.FonteRecurso, d.CodFonte, d.CodigoFonte, d.FonteSTN, d.NomeFonteSTN,
			d.ProcLicitatorio, d.Modalidade,
		)

		if err != nil {
			tx.Rollback()
			return fmt.Errorf("erro ao inserir empenho %d: %w", d.Empenho, err)
		}
	}

	if err = tx.Commit(); err != nil {
		return fmt.Errorf("erro ao salvar a transação no banco: %w", err)
	}

	return nil
}

func SalvarLicitacoes(licitacoes []entity.Licitacoes) error {
	strConexao := "postgres://pirajui:vKEP82XuP@ssw0rdMoreka@2129@localhost:5432/Pirajui?sslmode=disable"
	db, err := sql.Open("postgres", strConexao)
	if err != nil {
		return err
	}
	defer db.Close()

	tx, err := db.Begin()
	if err != nil {
		return err
	}

	queryLimpaBase := `truncate table licitacoes restart identity;`

	stmt1, err := tx.Prepare(queryLimpaBase)
	if err != nil {
		return err
	}
	defer stmt1.Close()

	query := `
		INSERT INTO licitacoes (
			proc_licitatorio, proc_administrativo, modalidade, exercicio, num_mod,
			situacao, data_abert_propost, hora_abert_propost, valor_previsto, valor_total_licitacao,
			objeto, data_edital, data_encerramento, carona, reg_preco, prazo_entrega_inicio,
			artigo_inciso, data_inicio_proposta, data_fim_proposta
		) VALUES (
			$1, $2, $3, $4, $5, $6, $7, $8, $9, $10,
			$11, $12, $13, $14, $15, $16, $17, $18, $19
		)`

	stmt, err := tx.Prepare(query)
	if err != nil {
		return err
	}
	defer stmt.Close()

	for _, l := range licitacoes {
		_, err = stmt.Exec(
			l.ProcLicitatorio, l.ProcAdministrativo, l.Modalidade, l.Exercicio, l.NumMod,
			l.Situacao,
			validaData(l.DataAbertPropost), l.HoraAbertPropost, l.ValorPrevisto, l.ValorTotalLicitacao,
			l.Objeto,
			validaData(l.DataEdital), validaData(l.DataEncerramento), l.Carona, l.RegPreco,
			l.PrazoEntregaInicio, l.ArtigoInciso, validaData(l.DataInicioProposta), validaData(l.DataFimProposta),
		)
		if err != nil {
			tx.Rollback()
			return err
		}
	}
	return tx.Commit()
}

// validaData garante que strings vazias virem NULL no banco
func validaData(d string) interface{} {
	if d == "" {
		return nil
	}
	return d
}

func SalvarContratos(contratos []entity.Contratos) error {
	strConexao := "postgres://pirajui:vKEP82XuP@ssw0rdMoreka@2129@localhost:5432/Pirajui?sslmode=disable"
	db, err := sql.Open("postgres", strConexao)
	if err != nil {
		return err
	}
	defer db.Close()

	tx, err := db.Begin()
	if err != nil {
		return err
	}

	// 1. LIMPA A TABELA ANTES DE INSERIR
	// Usamos tx.Exec direto pois a query não tem parâmetros ($1, etc), é mais limpo.
	_, err = tx.Exec("TRUNCATE TABLE contratos RESTART IDENTITY;")
	if err != nil {
		tx.Rollback()
		return err
	}

	// 2. PREPARA O INSERT
	query := `
		INSERT INTO contratos (
			num_contrato, num_detalhado_contrato, num_modalidade, modalidade, exercicio,
			fundamento_legal, proc_licitatorio, cpf_cnpj_fornecedor, fornecedor, valor,
			vigencia_inicial, vencimento_atual, objeto, tipo, contrato_rateio
		) VALUES (
			$1, $2, $3, $4, $5, $6, $7, $8, $9, $10,
			$11, $12, $13, $14, $15
		)`

	stmt, err := tx.Prepare(query)
	if err != nil {
		tx.Rollback()
		return err
	}
	defer stmt.Close()

	// 3. EXECUTA AS INSERÇÕES
	for _, c := range contratos {
		_, err = stmt.Exec(
			c.NumContrato, c.NumDetalhadoContrato, c.NumModalidade, c.Modalidade, c.Exercicio,
			c.FundamentoLegal, c.ProcLicitatorio, c.CpfCnpjFornecedor, c.Fornecedor, c.Valor,
			validaData(c.VigenciaInicial), validaData(c.VencimentoAtual), c.Objeto, c.Tipo, c.ContratoRateio,
		)
		if err != nil {
			tx.Rollback()
			return err
		}
	}

	// Salva a limpeza e a inserção de uma vez só!
	return tx.Commit()
}

func SalvarCargosSalarios(cargos []entity.CargosSalarios) error {
	strConexao := "postgres://pirajui:vKEP82XuP@ssw0rdMoreka@2129@localhost:5432/Pirajui?sslmode=disable"
	db, err := sql.Open("postgres", strConexao)
	if err != nil {
		return err
	}
	defer db.Close()

	tx, err := db.Begin()
	if err != nil {
		return err
	}

	// Limpa a tabela
	_, err = tx.Exec("TRUNCATE TABLE cargos_salarios RESTART IDENTITY;")
	if err != nil {
		tx.Rollback()
		return err
	}

	query := `
		INSERT INTO cargos_salarios (
			plano_cargo, cargo_id, cargo, referencia, valor, codigo
		) VALUES ($1, $2, $3, $4, $5, $6)`

	stmt, err := tx.Prepare(query)
	if err != nil {
		tx.Rollback()
		return err
	}
	defer stmt.Close()

	for _, c := range cargos {
		_, err = stmt.Exec(
			c.PlanoCargo, c.CargoId, c.Cargo, c.Referencia, c.Valor, c.Codigo,
		)
		if err != nil {
			tx.Rollback()
			return err
		}
	}

	return tx.Commit()
}

func SalvarFolha(folhas []entity.FolhaPagamento, nomeTabela string) error {
	strConexao := "postgres://pirajui:vKEP82XuP@ssw0rdMoreka@2129@localhost:5432/Pirajui?sslmode=disable"
	db, err := sql.Open("postgres", strConexao)
	if err != nil {
		return err
	}
	defer db.Close()

	tx, err := db.Begin()
	if err != nil {
		return err
	}

	// A consulta usa o nomeTabela dinâmico (folha_mensal ou folha_rescisao)
	query := fmt.Sprintf(`
		INSERT INTO %s (
			detalhe, referencia, referencia_salarial, nome, divisao, cargo, matricula,
			proventos, descontos, liquido, data_admissao, data_desligamento, tipo_regime,
			situacao_funcional, tipo_contrato, data_prevista_termino_contrato
		) VALUES (
			$1, $2, $3, $4, $5, $6, $7, $8, $9, $10,
			$11, $12, $13, $14, $15, $16
		)`, nomeTabela)

	stmt, err := tx.Prepare(query)
	if err != nil {
		tx.Rollback()
		return err
	}
	defer stmt.Close()

	for _, f := range folhas {
		_, err = stmt.Exec(
			f.Detalhe, f.Referencia, f.ReferenciaSalarial, f.Nome, f.Divisao, f.Cargo, f.Matricula,
			f.Proventos, f.Descontos, f.Liquido,
			validaData(f.DataAdmissao), validaData(f.DataDesligamento), f.TipoRegime,
			f.SituacaoFuncional, f.TipoContrato, validaData(f.DataPrevistaTerminoContrato),
		)
		if err != nil {
			tx.Rollback()
			return err
		}
	}

	return tx.Commit()
}
