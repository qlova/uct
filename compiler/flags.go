package compiler

type Flag struct {
	Name Translatable
	Value int
	Data string
	Bool bool
	
	Defined bool
	
	OnLost func(*Compiler)
}

func (c *Compiler) SetFlag(f Flag) {
	var scope = c.Scope[len(c.Scope)-1]
	
	f.Defined = true
	
	scope.Flags[f.Name[c.Language]] = f
}


func (c *Compiler) SetGlobalFlag(f Flag) {
	var scope = c.GlobalScope
	
	f.Defined = true
	
	scope.Flags[f.Name[c.Language]] = f
}

func (c *Compiler) GlobalFlagExists(f Flag) bool {
	var name = f.Name[c.Language]
	if _, ok := c.GlobalScope.Flags[name]; ok {
		return true
	}

	return false
}


func (c *Compiler) GetFlag(flag Flag) (Flag, int) {
	var name = flag.Name[c.Language]
	for i:=len(c.Scope)-1; i>=0; i-- {
		if v, ok := c.Scope[i].Flags[name]; ok {
			return v, i
		}
	} 

	return Flag{}, -1
}

func (c *Compiler) UpdateFlag(flag Flag) {
	var name = flag.Name[c.Language]
	for i:=len(c.Scope)-1; i>=0; i-- {
		if _, ok := c.Scope[i].Flags[name]; ok {
			c.Scope[i].Flags[name] = flag
		}
	}
}

func (c *Compiler) DeleteFlag(flag Flag) {
	var name = flag.Name[c.Language]
	for i:=len(c.Scope)-1; i>=0; i-- {
		if _, ok := c.Scope[i].Flags[name]; ok {
			delete(c.Scope[i].Flags, name)
		}
	}
}
