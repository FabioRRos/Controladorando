package services

import (
	entity "bot_portal_pref/Models/Entity"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

// BuscarReceitasTCE bate na rota de receitas do tribunal
func BuscarReceitasTCE(municipio string, ano int, mes int) ([]entity.ReceitasTce, error) {
	url := fmt.Sprintf("https://transparencia.tce.sp.gov.br/api/json/receitas/%s/%d/%d", municipio, ano, mes)

	resp, err := http.Get(url)
	if err != nil {
		return nil, fmt.Errorf("erro ao conectar na API do TCE (Receitas): %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("API do TCE retornou status: %d", resp.StatusCode)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("erro ao ler o corpo da resposta: %w", err)
	}

	var receitas []entity.ReceitasTce

	if err := json.Unmarshal(body, &receitas); err != nil {
		return nil, fmt.Errorf("erro ao decodificar JSON de receitas: %w", err)
	}

	return receitas, nil
}
