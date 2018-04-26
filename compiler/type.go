package compiler

type DefaultBase byte

func (d DefaultBase) Push(c *Compiler, name string) {
	switch d {
		case INT:
			c.Push(name)
		case LIST:
			c.PushList(name)
		case PIPE:
			c.PushPipe(name)
		
		default:
			panic("Invalid default base!")
	}
}

func (d DefaultBase) Pull(c *Compiler, name string) {
	switch d {
		case INT:
			c.Pull(name)
		case LIST:
			c.PullList(name)
		case PIPE:
			c.PullPipe(name)
		
		default:
			panic("Invalid default base!")
	}
}

func (d DefaultBase) Free(c *Compiler) {
	switch d {
		case INT:
			c.Free()
		case LIST:
			c.FreeList()
		case PIPE:
			c.FreePipe()
		
		default:
			panic("Invalid default base!")
	}
}

func (d DefaultBase) Drop(c *Compiler) {
	switch d {
		case INT:
			c.Drop()
		case LIST:
			c.DropList()
		case PIPE:
			c.DropPipe()
		
		default:
			panic("Invalid default base!")
	}
}

const (
	INT DefaultBase = iota
	LIST
	PIPE
)

type Data interface {
	Equals(Data) bool
	Name(Language) string
}

type Type struct {
	Data
	
	Name Translatable
	Base Base
	
	Immutable bool
	Defined bool
	
	Casts []func(c *Compiler, t Type) bool
	Shunts map[string]func(*Compiler, Type) Type
	
	Collect func(*Compiler)
	
	EmbeddedStatement func(*Compiler, Type) 
}

func (t Type) String() string {
	if data := t.Data; data != nil && data.Name(0) != "" {
		return t.Name[0]+"."+data.Name(0)
	}
	return t.Name[0]
}

func (t Type) Equals(b Type) bool {
	if t.Name[English] == b.Name[English] {
		return true
	}
	return false
}
 
func (c *Compiler) GetType(name string) *Type {
	for i, t := range c.Types {
		if t.Name[c.Language] == name {
			return &c.Types[i]
		}
	}
	return nil
}

func (c *Compiler) Cast(a, b Type) bool {
	if a.Equals(b) {
		return true
	}
	
	for _, t := range a.Casts {
		if t(c, b) {
			return true
		}
	}
	c.RaiseError(Translatable{
		English: "Cannot cast "+a.String()+" to "+b.String(),
	})
	return false
}
