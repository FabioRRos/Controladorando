package entity

type ReceitasTce struct {
	Orgao                string
	Mes                  string
	Ds_fonte_recurso     string
	Ds_cd_aplicacao_fixo string
	Ds_alinea            string
	Ds_subalinea         string
	Vl_arrecadacao       string // Recebemos como string da API e convertemos no banco
}
