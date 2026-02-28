package database

import (
	entity "bot_portal_pref/Models/Entity"
	"database/sql"
	"fmt"
	"strconv"
	"strings"
)

func SalvarDespesasTCE(despesas []entity.DespesasTce, mes string) error {
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

	// Limpa o mês que estamos importando para não duplicar dados
	queryLimpa := "DELETE FROM despesas_tce WHERE mes = $1"
	if _, err = tx.Exec(queryLimpa, mes); err != nil {
		tx.Rollback()
		return err
	}

	query := `
		INSERT INTO despesas_tce (
			orgao, mes, evento, nr_empenho, id_fornecedor, 
			nm_fornecedor, dt_emissao_despesa, vl_despesa
		) VALUES ($1, $2, $3, $4, $5, $6, $7, $8)`

	stmt, err := tx.Prepare(query)
	if err != nil {
		tx.Rollback()
		return err
	}
	defer stmt.Close()

	for _, d := range despesas {
		_, err = stmt.Exec(
			d.Orgao,
			d.Mes,
			d.Evento,
			d.Nr_empenho,
			d.Id_fornecedor,
			d.Nm_fornecedor,
			converterDataTceParaSQL(d.Dt_emissao_despesa),
			converterMoedaTce(d.Vl_despesa),
		)
		if err != nil {
			tx.Rollback()
			return err
		}
	}

	return tx.Commit()
}

// Funções locais rápidas de conversão
func converterMoedaTce(valor string) float64 {
	v := strings.TrimSpace(valor)
	if v == "" {
		return 0.0
	}
	v = strings.ReplaceAll(v, ".", "")
	v = strings.ReplaceAll(v, ",", ".")
	resultado, _ := strconv.ParseFloat(v, 64)
	return resultado
}

func converterDataTceParaSQL(dataBR string) interface{} {
	partes := strings.Split(strings.TrimSpace(dataBR), "/")
	if len(partes) == 3 {
		return fmt.Sprintf("%s-%s-%s", partes[2], partes[1], partes[0])
	}
	return nil
}
