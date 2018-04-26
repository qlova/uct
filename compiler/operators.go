package compiler 

type Operator struct {
	Symbol string
	Precedence int
}

func (c *Compiler) GetOperator(symbol string) Operator {
	for _, op := range c.Operators {
		if op.Symbol == symbol {
			return op
		}
	}

	return Operator{Precedence:-1}
}
