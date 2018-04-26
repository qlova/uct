package compiler

type Statement struct {
	Name Translatable
	
	OnScan func(*Compiler)
	Detect func(*Compiler) bool
}

func (c *Compiler) ScanStatement() {
	var token = c.Scan()
	
	if token == "\n" || token == "" {
		return
	}
	
	var found = false
	for _, statement := range c.Statements {
		if statement.Name[c.Language] == token {
			found = true
			statement.OnScan(c)
			return
		}
	}
	
	for i := len(c.Statements)-1; i>0; i-- {
		statement := c.Statements[i]
		
		if statement.Detect != nil && statement.Detect(c) {
			found = true
			return
		}
	}
	
	if !found {
		c.RaiseError(Translatable{
				English:"Unknown Statement: "+token,
		})
	}
}
