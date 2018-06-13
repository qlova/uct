package compiler

type Expression struct {
	Name Translatable
	OnScan func(*Compiler) Type
	
	Detect func(*Compiler) *Type
}

func (c *Compiler) Expression() Type {
	var token = c.Scan()
	for token == "\n" {
		token = c.Scan()
	}
	
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
	
	
	return Type{Name: NoTranslation(c.Token()), Fake: true}
}


func (c *Compiler) ScanExpression() Type {
	var result = c.Shunt(c.Expression(), 0)
	
	if result.Fake {
		c.RaiseError(Translatable{
				English: "Unknown Expression: "+result.Name[c.Language],
		})
	}
	
	return result
}
