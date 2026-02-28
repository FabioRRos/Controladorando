package database

import (
	entity "bot_portal_pref/Models/Entity"
	"database/sql"
)

func SalvarReceitasTCE(receitas []entity.ReceitasTce, mes string) error {
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

	// A documentação do TCE diz que 'mes' vem por extenso no JSON (ex: "JANEIRO"),
	// mas para garantir, limpamos usando exatamente o dado que a API nos devolveu.
	if len(receitas) > 0 {
		mesExtensoApi := receitas[0].Mes
		queryLimpa := "DELETE FROM receitas_tce WHERE mes = $1"
		if _, err = tx.Exec(queryLimpa, mesExtensoApi); err != nil {
			tx.Rollback()
			return err
		}
	}

	query := `
		INSERT INTO receitas_tce (
			orgao, mes, ds_fonte_recurso, ds_cd_aplicacao_fixo, 
			ds_alinea, ds_subalinea, vl_arrecadacao
		) VALUES ($1, $2, $3, $4, $5, $6, $7)`

	stmt, err := tx.Prepare(query)
	if err != nil {
		tx.Rollback()
		return err
	}
	defer stmt.Close()

	for _, r := range receitas {
		_, err = stmt.Exec(
			r.Orgao,
			r.Mes,
			r.Ds_fonte_recurso,
			r.Ds_cd_aplicacao_fixo,
			r.Ds_alinea,
			r.Ds_subalinea,
			converterMoedaTce(r.Vl_arrecadacao), // Usa a mesma função auxiliar que já criamos nas despesas
		)
		if err != nil {
			tx.Rollback()
			return err
		}
	}

	return tx.Commit()
}
