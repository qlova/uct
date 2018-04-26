 package compiler
 
type Scope struct {
	Variables map[string]Variable
	Flags map[string]Flag
}

func NewScope() Scope {
	var s Scope
	s.Variables = make(map[string]Variable)
	s.Flags = make(map[string]Flag)
	return s
}

func (c *Compiler) GainScope() {
	c.Scope = append(c.Scope, NewScope())
}

func (c *Compiler) GetScope() Scope {
	return c.Scope[len(c.Scope)-1]
}

func (c *Compiler) LoseScope() {
	
	var scope = c.Scope[len(c.Scope)-1]
	
	for name, variable := range scope.Variables {
		if !variable.Protected {
			c.PushType(variable.Type, name)
			
			if variable.Type.Collect != nil {
				variable.Type.Collect(c)
			}
			
			c.FreeType(variable.Type)
		}
	}
	
	for _, flag := range scope.Flags {
		if flag.OnLost != nil {
			flag.OnLost(c)
		}
	}
	
	c.Scope = c.Scope[:len(c.Scope)-1]
}
