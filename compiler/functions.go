package compiler

type Function struct {
	Name Translatable
	
	Arguments []Type
	Tokens []string
	
	Returns []Type
	
	Flags []Flag
	
	Inline func(*Compiler)
	Compile func(*Compiler)
	
	Variadic bool
	
	compiled bool
	
	Data interface{}
}

func (c *Compiler) GetFunction(name string) *Function {
	for _, t := range c.Functions {
		if t.Name[c.Language] == name {
			return t
		}
	}
	return nil
}
