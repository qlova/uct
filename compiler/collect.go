package compiler

func (c *Compiler) CollectAll() {
	for _, scope := range c.Scope {
		for name, variable := range scope.Variables {
			if !variable.Protected {
				c.PushType(variable.Type, name)
				
				if variable.Type.Collect != nil {
					variable.Type.Collect(c)
				}
				
				c.FreeType(variable.Type)
			}
		}
	}
}
