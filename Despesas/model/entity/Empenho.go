package entity

type Empenho struct {
	Id                      int
	Exercicio               int     // Ano do exercício
	Data                    string  // Data do empenho
	NumeroEmpenho           int     // Número do empenho
	Tipo                    string  // Tipo (AD, OR, DA...)
	Favorecido              string  // Nome do Favorecido
	CpfCnpj                 string  // CPF/CNPJ
	Valor                   float64 // Valor do empenho
	ProcessoContratacao     string  // Processo de contratação
	TipoLicitacao           string  // Modalidade (DISPENSA, PREGÃO...)
	NumeroLicitacao         string  // Número da licitação
	Orgao                   string  // Órgão
	UnidadeOrcamentaria     string  // Unidade Orçamentária
	ProjetoAtividade        string  // Projeto/Atividade
	VinculoOrcamentario     string  // Vínculo Orçamentário
	GrupoFonte              string  // Grupo da Fonte
	CodigoFonte             string  // Código da Fonte
	Elemento                string  // Elemento (ex: 30 - Material de Consumo)
	Historico               string  // Histórico do empenho
	Poder                   string  // Poder (Executivo/Legislativo)
	Funcao                  string  // Código/Nome da Função
	Subfuncao               string  // Código/Nome da Subfunção
	Programa                string  // Programa
	Fonro                   string  // Código interno da Fonte
	FonteSTN                string  // Código Fonte STN
	CategoriaEconomica      string  // Categoria Econômica
	GrupoNatureza           string  // Grupo Natureza
	ModalidadeAplicacao     string  // Modalidade de Aplicação
	Desdobro                string  // Desdobro
	Natureza                string  // Natureza completa (3.3.90.30.07)
	NumeroContrato          string  // Número do contrato
	ContratoNumDetalhado    string  // Número detalhado do contrato
	VigenciaInicio          string  // Vigência início
	VigenciaFim             string  // Vigência fim
	Termo                   string  // Termo
	ContratoAditId          string  // ID do aditivo
	TermoResgatado          string  // Termo resgatado
	NumeroConvenio          string  // Número do convênio
	AnoConvenio             int     // Ano do convênio
	TipoFundamento          string  // Tipo de fundamento legal
	Inciso                  string  // Inciso legal
	ContratoVencimentoAtual string  // Data vencimento atual do contrato
	Ficha                   int     // Número da ficha

}
