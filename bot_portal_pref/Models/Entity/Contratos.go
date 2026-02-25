package entity

type Contratos struct {
	Id                   int
	NumContrato          string  // N° Contrato
	NumDetalhadoContrato string  // Nº Detalhado do Contrato
	NumModalidade        string  // N° Modalidade
	Modalidade           string  // Modalidade
	Exercicio            int     // Exercício
	FundamentoLegal      string  // Fundamento Legal
	ProcLicitatorio      string  // Proc. Licitatório
	CpfCnpjFornecedor    string  // CPF/CNPJ Fornecedor
	Fornecedor           string  // Fornecedor
	Valor                float64 // Valor
	VigenciaInicial      string  // Vigência Inicial
	VencimentoAtual      string  // Vencimento Atual
	Objeto               string  // Objeto
	Tipo                 int     // Tipo
	ContratoRateio       string  // Contrato de Rateio
}
