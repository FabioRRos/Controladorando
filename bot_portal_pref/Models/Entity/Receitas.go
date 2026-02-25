package entity

type Receitas struct {
	Id                 int
	Codigo             string  // Classificação orçamentária da receita.
	Especificacao      string  // O nome da conta (ou da arrecadação)
	CodAplicacao       string  // Origem ou destino da receita
	FonteSTN           string  // Padrão nacional de codificação para que o Tesouro consiga consolidar as contas
	FonteRecurco       string  // Indica se o dinheiro é livre ou se é destinado
	PrevisaoInicial    float64 // Previsão da arrecadação
	PrevisaoAtualizada float64 // Previsão atualizada da arrecadação
	ArrecadacaoPeriodo float64 // O quanto de grana entrou no cofre
	ArrecadacaoTotal   float64 // A arrecadação total deste cara
}
