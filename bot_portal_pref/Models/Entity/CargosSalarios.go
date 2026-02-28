package entity

type CargosSalarios struct {
	Id         int
	PlanoCargo string  // Plano Cargo
	CargoId    int     // ID (do CSV)
	Cargo      string  // Cargo
	Referencia string  // Referência
	Valor      float64 // Valor
	Codigo     string  // Código
}
