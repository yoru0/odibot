package domain

import "time"

type Calculation struct {
	ID        string    `json:"id"         db:"id"`
	UserID    string    `json:"user_id"    db:"user_id"`
	Operation string    `json:"operation"  db:"operation"`
	Operand1  float64   `json:"operand1"   db:"operand1"`
	Operand2  float64   `json:"operand2"   db:"operand2"`
	Result    float64   `json:"result"     db:"result"`
	CreatedAt time.Time `json:"created_at" db:"created_at"`
}

type CalculationResult struct {
	Result      float64      `json:"result"`
	Calculation *Calculation `json:"calculation"`
}
