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
		case NULL:
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
		case NULL:
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
		case NULL:
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
		case NULL:
		default:
			panic("Invalid default base!")
	}
}

func (d DefaultBase) Attach(c *Compiler) {
	switch d {
		case INT:
			c.Set()
		case LIST:
			c.Int(0)
			c.HeapList()
			c.Set()
		case PIPE:
			c.Int(0)
			c.HeapPipe()
			c.Set()
		case NULL:
		default:
			panic("Invalid default base!")
	}
}

func (d DefaultBase) Detach(c *Compiler) {
	switch d {
		case INT:
			c.Get()
		case LIST:
			c.Get()
			c.HeapList()
		case PIPE:
			c.Get()
			c.HeapPipe()
		case NULL:
		default:
			panic("Invalid default base!")
	}
}

const (
	NULL DefaultBase = iota
	INT 
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
	Fake bool
	Constant bool
	
	Cast func(c *Compiler, a Type, b Type) bool
	Shunt func(c *Compiler, symbol string, a Type, b Type) *Type
	
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
		
		if t.Data != nil && b.Data != nil {
			if !t.Data.Equals(b.Data) {
				return false
			}
		}
		
		return true
	}
	return false
}
 
func (t Type) With(d Data) Type {
	t.Data = d
	return t
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
	
	if a.Cast != nil &&  a.Cast(c, a, b) {
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
