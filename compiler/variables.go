package compiler

type Variable struct {
	Type
	Defined, Protected, Modified, Embedded bool
	
	Index int
	List string
	
	DefinedAtLineNumber int
	DefinedAtLine string
}

func (c *Compiler) SetVariable(name string, t Type) {
	if _, ok := c.GetScope().Variables[name]; ok {
		c.RaiseError(Translatable{
			English: name+" already defined!",
		})
	}
	c.GetScope().Variables[name] = Variable{ 
		Type: t,
		Defined: true,
		DefinedAtLineNumber: c.Scanners[len(c.Scanners)-1].Line,
		DefinedAtLine: c.CurrentLines[len(c.Scanners)-1],
	}
}

func (c *Compiler) UpdateVariable(name string, t Type) {
	for i:=len(c.Scope)-1; i>=0; i-- {
		if v, ok := c.Scope[i].Variables[name]; ok {
			v.Type = t
			c.Scope[i].Variables[name] = v
		}
	} 
}

func (c *Compiler) GetVariable(name string) Variable {
	for i:=len(c.Scope)-1; i>=0; i-- {
		if v, ok := c.Scope[i].Variables[name]; ok {
			return v
		}
	} 

	return Variable{}
}
