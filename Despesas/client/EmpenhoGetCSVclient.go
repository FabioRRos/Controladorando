package client

import (
	"bytes"
	"despesas/model/entity"
	"encoding/csv"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"time"

	"golang.org/x/text/encoding/charmap"
)

// FetchEmpenho faz o POST, pega o CSV, converte e retorna o PRIMEIRO registro como entity.Empenho.
// cookieValue: só o valor do ASP.NET_SessionId (ex: "brpt0wkxbl1frgchuc4obxey")
func FetchEmpenho(cookieValue string) (entity.Empenho, error) {
	endpoint := "http://prefeiturapirajui1.ddns.net:8079/Transparencia/DadosEmpenho.aspx"

	// Form data do clique do menu
	form := url.Values{}
	form.Set("__EVENTTARGET", "ASPxPopupMenu2")
	form.Set("__EVENTARGUMENT", "CLICK:0")

	req, err := http.NewRequest("POST", endpoint, strings.NewReader(form.Encode()))
	if err != nil {
		return entity.Empenho{}, err
	}

	// Headers mínimos e estáveis
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Set("Accept", "*/*")
	req.Header.Set("User-Agent", "Mozilla/5.0")
	req.Header.Set("Referer", endpoint)

	// Cookie (monta igual ao Postman)
	req.Header.Set("Cookie", "ASP.NET_SessionId="+cookieValue)

	client := &http.Client{Timeout: 30 * time.Second}

	resp, err := client.Do(req)
	if err != nil {
		return entity.Empenho{}, err
	}
	defer resp.Body.Close()

	bodyRaw, err := io.ReadAll(resp.Body)
	if err != nil {
		return entity.Empenho{}, err
	}

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		// devolve um pedaço do body pra debug sem explodir o console
		snippet := string(bodyRaw)
		if len(snippet) > 500 {
			snippet = snippet[:500]
		}
		return entity.Empenho{}, fmt.Errorf("POST retornou HTTP %d. Body (início): %s", resp.StatusCode, snippet)
	}

	// Converte Windows-1252 -> UTF-8 (resolve "Educa��o" etc)
	bodyUTF8, err := charmap.Windows1252.NewDecoder().Bytes(bodyRaw)
	if err != nil {
		// se falhar, tenta seguir com raw mesmo
		bodyUTF8 = bodyRaw
	}

	emp, err := parseFirstEmpenhoFromCSV(bodyUTF8)
	if err != nil {
		return entity.Empenho{}, err
	}

	return emp, nil
}

func parseFirstEmpenhoFromCSV(csvBytes []byte) (entity.Empenho, error) {
	r := csv.NewReader(bytes.NewReader(csvBytes))
	r.Comma = ';'
	r.LazyQuotes = true

	// 1ª linha: "Dados Empenho;;;;"
	_, err := r.Read()
	if err != nil {
		return entity.Empenho{}, fmt.Errorf("erro lendo linha 1 do csv: %w", err)
	}

	// 2ª linha: header
	_, err = r.Read()
	if err != nil {
		return entity.Empenho{}, fmt.Errorf("erro lendo header do csv: %w", err)
	}

	// Próxima(s) linha(s): dados (podem ter \n dentro de campos com aspas)
	rec, err := r.Read()
	if err == io.EOF {
		return entity.Empenho{}, errors.New("csv veio sem registros (apenas cabeçalho)")
	}
	if err != nil {
		return entity.Empenho{}, fmt.Errorf("erro lendo registro do csv: %w", err)
	}

	// Mapeamento por índice (batendo com o header que você mostrou)
	// 0 Ex.; 1 Data; 2 Emp.; 3 Tipo; 4 Favorecido; 5 CPF/CNPJ; 6 Valor; ...
	get := func(i int) string {
		if i < 0 || i >= len(rec) {
			return ""
		}
		return strings.TrimSpace(rec[i])
	}

	emp := entity.Empenho{
		Exercicio:               parseInt(get(0)),
		Data:                    get(1),
		NumeroEmpenho:           parseInt(get(2)),
		Tipo:                    get(3),
		Favorecido:              get(4),
		CpfCnpj:                 get(5),
		Valor:                   parseBRFloat(get(6)),
		ProcessoContratacao:     get(7),
		TipoLicitacao:           get(8),
		NumeroLicitacao:         get(9),
		Orgao:                   get(10),
		UnidadeOrcamentaria:     get(11),
		ProjetoAtividade:        get(12),
		VinculoOrcamentario:     get(13),
		GrupoFonte:              get(14),
		CodigoFonte:             get(15),
		Elemento:                get(16),
		Historico:               get(17),
		Poder:                   get(18),
		Funcao:                  get(19),
		Subfuncao:               get(20),
		Programa:                get(21),
		Fonro:                   get(22),
		FonteSTN:                get(23),
		CategoriaEconomica:      get(24),
		GrupoNatureza:           get(25),
		ModalidadeAplicacao:     get(26),
		Desdobro:                get(27),
		Natureza:                get(28),
		NumeroContrato:          get(29),
		ContratoNumDetalhado:    get(30),
		VigenciaInicio:          get(31),
		VigenciaFim:             get(32),
		Termo:                   get(33),
		ContratoAditId:          get(34),
		TermoResgatado:          get(35),
		NumeroConvenio:          get(36),
		AnoConvenio:             parseInt(get(37)),
		TipoFundamento:          get(38),
		Inciso:                  get(39),
		ContratoVencimentoAtual: get(40),
		Ficha:                   parseInt(get(41)),
	}

	return emp, nil
}

func parseInt(s string) int {
	s = strings.TrimSpace(s)
	if s == "" {
		return 0
	}
	i, _ := strconv.Atoi(s)
	return i
}

// parseBRFloat converte "9.719,48" -> 9719.48
func parseBRFloat(s string) float64 {
	s = strings.TrimSpace(s)
	if s == "" {
		return 0
	}
	s = strings.ReplaceAll(s, ".", "")
	s = strings.ReplaceAll(s, ",", ".")
	f, _ := strconv.ParseFloat(s, 64)
	return f
}
