package services

import (
	entity "bot_portal_pref/Models/Entity"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

// BuscarDespesasTCE consome a API oficial do Tribunal de Contas
func BuscarDespesasTCE(municipio string, ano int, mes int) ([]entity.DespesasTce, error) {
	url := fmt.Sprintf("https://transparencia.tce.sp.gov.br/api/json/despesas/%s/%d/%d", municipio, ano, mes)

	resp, err := http.Get(url)
	if err != nil {
		return nil, fmt.Errorf("erro ao conectar na API do TCE: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("API do TCE retornou status: %d", resp.StatusCode)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("erro ao ler o corpo da resposta: %w", err)
	}

	var despesas []entity.DespesasTce

	// A mágica acontece aqui: o JSON entra e o Go mapeia para a struct automaticamente
	if err := json.Unmarshal(body, &despesas); err != nil {
		return nil, fmt.Errorf("erro ao decodificar JSON: %w", err)
	}

	return despesas, nil
}
