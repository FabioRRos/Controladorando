package entity

type FolhaPagamento struct {
	Id                          int
	Detalhe                     string  // Detalhe (Vazio no CSV geralmente)
	Referencia                  string  // Referência (Ex: Folha Mensal - Janeiro)
	ReferenciaSalarial          string  // Referência Salarial
	Nome                        string  // Nome
	Divisao                     string  // Divisão
	Cargo                       string  // Cargo
	Matricula                   int     // Matrícula
	Proventos                   float64 // PROVENTOS
	Descontos                   float64 // DESCONTOS
	Liquido                     float64 // LIQUIDO
	DataAdmissao                string  // Data Admissão
	DataDesligamento            string  // Data Desligamento
	TipoRegime                  string  // Tipo de Regime
	SituacaoFuncional           string  // Situação Funcional
	TipoContrato                string  // Tipo de Contrato
	DataPrevistaTerminoContrato string  // Data Prevista Termino Contrato
}
