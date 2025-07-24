package calculator

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/yoru0/goodi/internal/ports"
)

type CalculatorBot struct {
	Store ports.ResultScore
}

func (cb *CalculatorBot) Process(input string) string {
	tokens := strings.Fields(input)
	if len(tokens) != 3 || tokens[0] != "!add" {
		return "Usage: !add <num1> <num2>"
	}

	a, err1 := strconv.Atoi(tokens[1])
	b, err2 := strconv.Atoi(tokens[2])
	if err1 != nil || err2 != nil {
		return "Both arguments must be integers."
	}

	result := a + b
	err := cb.Store.SaveResult(a, b, result)
	if err != nil {
		return "Failed to save to database."
	}

	return fmt.Sprintf("Result: %d", result)
}
