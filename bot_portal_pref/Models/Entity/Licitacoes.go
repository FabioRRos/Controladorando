package entity

type Licitacoes struct {
	Id                  int
	ProcLicitatorio     string  // Proc. Licitatório
	ProcAdministrativo  string  // Proc. Administrativo
	Modalidade          string  // Modalidade
	Exercicio           int     // Exercício
	NumMod              int     // N° Mod.
	Situacao            string  // Situação
	DataAbertPropost    string  // Data Abert. Propost.
	HoraAbertPropost    string  // Hora Abert. Propost.
	ValorPrevisto       float64 // Valor Previsto
	ValorTotalLicitacao float64 // Valor Total Licitação
	Objeto              string  // Objeto
	DataEdital          string  // Data do Edital
	DataEncerramento    string  // Data Encerramento
	Carona              string  // Carona
	RegPreco            string  // Reg. Preço
	PrazoEntregaInicio  string  // Prazo de Entrega/Início
	ArtigoInciso        string  // Artigo/Inciso
	DataInicioProposta  string  // Data Inicio Proposta
	DataFimProposta     string  // Data Fim Proposta
}
