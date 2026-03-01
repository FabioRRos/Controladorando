package repository

import (
	"database/sql"
	"despesas/model/dto"
	"fmt"
)

func BuscarDespesasPendentesUnicas(connStr string) ([]dto.DespesasDTO, error) {
	// 1. Abrir conexão
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		return nil, fmt.Errorf("erro ao conectar no banco.")
	}
	defer db.Close()

	// 2. Executar a query
	query := `SELECT d.id ,d.empenho, d.tipo
FROM despesas d
WHERE d.empenho IN (
    SELECT empenho
    FROM despesas
    GROUP BY empenho
    HAVING COUNT(*) = 1
)order by id asc;`
	rows, err := db.Query(query)
	if err != nil {
		return nil, fmt.Errorf("erro ao buscar despesas pendentes: %w", err)
	}
	defer rows.Close()

	var lista []dto.DespesasDTO

	// 3. Iterar pelos resultados
	for rows.Next() {
		var d dto.DespesasDTO
		err := rows.Scan(&d.Id, &d.Empenho, &d.Tipo)
		if err != nil {
			return nil, fmt.Errorf("erro ao escanear registro: %w", err)
		}
		lista = append(lista, d)
	}

	// 4. Verificar erros após o loop
	if err = rows.Err(); err != nil {
		return nil, err
	}

	return lista, nil
}

func BuscarDespesasPendentesMultiplas(connStr string) ([]dto.DespesasDTO, error) {
	// 1. Abrir conexão
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		return nil, fmt.Errorf("erro ao conectar no banco: %w", err)
	}
	defer db.Close()

	// 2. Executar a query
	query := `SELECT d.id ,d.empenho, d.tipo
FROM despesas d
WHERE d.empenho IN (
    SELECT empenho
    FROM despesas
    GROUP BY empenho
    HAVING COUNT(*) > 1
)
order by id asc;;`
	rows, err := db.Query(query)
	if err != nil {
		return nil, fmt.Errorf("erro ao buscar despesas pendentes: %w", err)
	}
	defer rows.Close()

	var lista []dto.DespesasDTO

	// 3. Iterar pelos resultados
	for rows.Next() {
		var d dto.DespesasDTO
		err := rows.Scan(&d.Id, &d.Empenho, &d.Tipo)
		if err != nil {
			return nil, fmt.Errorf("erro ao escanear registro: %w", err)
		}
		lista = append(lista, d)
	}

	// 4. Verificar erros após o loop
	if err = rows.Err(); err != nil {
		return nil, err
	}

	return lista, nil
}
func MarcarDespesaComoProcessada(connStr string, id int) error {

	fmt.Println("Baixar ID:", id)

	// 1. Abrir conexão
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		return fmt.Errorf("erro ao conectar no banco: %w", err)
	}
	defer db.Close()

	// 2. Executar o Update
	query := `UPDATE despesas SET processado = true WHERE id = $1`

	result, err := db.Exec(query, id)
	if err != nil {
		return fmt.Errorf("erro ao atualizar status da despesa %d: %w", id, err)
	}

	// 3. (Opcional) Verificar se alguma linha foi realmente afetada
	linhasAfetadas, _ := result.RowsAffected()
	if linhasAfetadas == 0 {
		return fmt.Errorf("nenhuma despesa encontrada com o ID %d", id)
	}

	return nil
}
