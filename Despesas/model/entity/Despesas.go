package entity

type Despesas struct {
	Id                    int
	Empenho               int     // Número do empenho
	Tipo                  string  // Tipo (AD, OR, DA...)
	NoFicha               int     // N° da Ficha
	Data                  string  // Data da transação
	CodForn               int     // Código do Fornecedor
	NomeFornecedor        string  // Nome do Fornecedor
	CpfCnpj               string  // CPF/CNPJ
	Dotacao               float64 // Dotação inicial
	AlteracaoDotacao      float64 // Alteração de dotação
	DotacaoAtual          float64 // Dotação atualizada
	ValorAnulado          float64 // Valor anulado
	Reforco               float64 // Reforço
	ValorEmpenhado        float64 // Valor empenhado
	ValorLiquidado        float64 // Valor liquidado
	ValorPago             float64 // Valor pago
	EmpenhadoAteHoje      float64 // Acumulado empenhado
	LiquidadoAteHoje      float64 // Acumulado liquidado
	PagoAteHoje           float64 // Acumulado pago
	Local                 string  // Código do Local
	Funcional             string  // Código Funcional
	Funcao                int     // Código da Função
	NomeFuncao            string  // Nome da Função (Saúde, Administração...)
	Subfuncao             int     // Código da Subfunção
	NomeSubfuncao         string  // Nome da Subfunção
	CodAplicacao          string  // Código de aplicação (ex: 310.000)
	DescricaoCodAplicacao string  // Descrição do Cód. de aplicação
	Natureza              string  // Código da Natureza
	NomeNatureza          string  // Descrição da Natureza
	Fonte                 int     // Código da Fonte
	FonteRecurso          string  // Nome da Fonte (TESOURO)
	CodFonte              int     // Código da Fonte detalhado
	CodigoFonte           string  // Descrição da Fonte (Recursos Ordinarios)
	FonteSTN              string  // Código Fonte STN (1.500)
	NomeFonteSTN          string  // Descrição Fonte STN
	ProcLicitatorio       string  // Processo Licitatório
	Modalidade            string  // Modalidade (DISPENSA, INEXIGIBILIDADE...)
	processado            bool    //Pra saber se já processei o empenho referente a este cara
}
