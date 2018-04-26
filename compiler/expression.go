package compiler

type Expression struct {
	Name Translatable
	OnScan func(*Compiler) Type
	
	Detect func(*Compiler) *Type
}

func (c *Compiler) scanExpression() Type {
	var token = c.Scan()
	
	for _, expression := range c.Expressions {
		if expression.Name[c.Language] == token {
			return expression.OnScan(c)
		}
	}
	
	for _, expression := range c.Expressions {
		if expression.Detect != nil {
			if t := expression.Detect(c); t != nil {
				return *t
			}
		}
	}
	
	/*c.RaiseError(Translatable{
			English: "Unknown Expression: "+c.Token(),
	})*/
	return Type{Name: NoTranslation(c.Token())}
}

func (c *Compiler) ScanExpression() Type {
	return c.Shunt(c.scanExpression(), 0)
}
