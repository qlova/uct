package compiler

import "github.com/qlova/uct/assembler"
import "io"
import "math/big"
import "bytes"


type Base interface{
	Push(c *Compiler, name string)
	Pull(c *Compiler, name string)
	Drop(c *Compiler)
	Free(c *Compiler)
	
	//Attach/Detach from a list.
	Attach(c *Compiler)
	Detach(c *Compiler)
}

func WriteString(output io.Writer, data string) {
	
	var length = big.NewInt(int64(len(data))).Bytes()
	
	if len(length) > 255 {
		panic("NUMBER TOO BIG")
	}
	
	output.Write([]byte{byte(len(length))})
	output.Write(length)
	
	//Should panic if string is longer than an int64...
	//Will this ever happen though?? na!
	
	for _, char := range []byte(data) {
		output.Write([]byte{1, char})	
	}
}

func Int(value int64) *big.Int {
	return big.NewInt(value)
}

func (c *Compiler) Main() {
	c.Output.Write([]byte{1, assembler.MAIN})
}

func (c *Compiler) Link(data string) {
	c.Output.Write([]byte{1, assembler.LINK})
	WriteString(c.Output, data)
}

func (c *Compiler) Exit() {
	c.Output.Write([]byte{1, assembler.EXIT})
}

func (c *Compiler) Code(data string) {
	c.Output.Write([]byte{1, assembler.CODE})
	WriteString(c.Output, data)
}

func (c *Compiler) Back() {
	c.Output.Write([]byte{1, assembler.BACK})
}

func (c *Compiler) Call(f *Function) {
	
	//Load the function if possible.
	if !f.compiled {
		f.compiled = true
		c.SwapOutput()
		
		var header, output, scope = c.Header, c.Output, c.Scope
		
		c.Header = c.Output
		var b bytes.Buffer
		c.Output = &b
		c.Scope = nil
		
		c.GainScope()
		f.Compile(c)
		c.LoseScope()
		c.Scope = scope
		
		c.Header, c.Output = header, output
		
		c.Code(f.Name[c.Language])
		
		c.Output.Write(b.Bytes())
	
		c.Back()
		c.SwapOutput()
		
		
	}
	
	c.CallRaw(f.Name[c.Language])
}

func (c *Compiler) CallRaw(data string) {
	c.Output.Write([]byte{1, assembler.CALL})
	WriteString(c.Output, data)
}

func (c *Compiler) Fork() {
	c.Output.Write([]byte{1, assembler.FORK})
}


func (c *Compiler) Error() string {
	return "ERROR"
}

func (c *Compiler) Global() string {
	return "GLOBAL"
}

func (c *Compiler) Channel() string {
	return "CHANNEL"
}

func (c *Compiler) Loop() {
	c.Output.Write([]byte{1, assembler.LOOP})
}

func (c *Compiler) Done() {
	c.Output.Write([]byte{1, assembler.DONE})
}

func (c *Compiler) Redo() {
	c.Output.Write([]byte{1, assembler.REDO})
}

func (c *Compiler) If() {
	c.Output.Write([]byte{1, assembler.IF})
}
func (c *Compiler) Or() {
	c.Output.Write([]byte{1, assembler.OR})
}
func (c *Compiler) No() {
	c.Output.Write([]byte{1, assembler.NO})
}

func (c *Compiler) BigData(name string, data ...*big.Int) {
	c.Output.Write([]byte{1, assembler.DATA})
	WriteString(c.Output, name)
	
	var length = big.NewInt(int64(len(data))).Bytes()
	
	if len(length) > 255 {
		panic("NUMBER TOO BIG")
	}
	
	c.Output.Write([]byte{byte(len(length))})
	c.Output.Write(length)
	
	for _, value := range data {
		var bytes = value.Bytes()
		
		if len(bytes) > 255 {
			panic("NUMBER TOO BIG")
		}
		
		c.Output.Write([]byte{byte(len(bytes))})
		c.Output.Write(bytes)
	}
}

func (c *Compiler) Data(name string, data []byte) {
	c.Output.Write([]byte{1, assembler.DATA})
	WriteString(c.Output, name)
	
	var length = big.NewInt(int64(len(data))).Bytes()
	
	if len(length) > 255 {
		panic("NUMBER TOO BIG")
	}
	
	c.Output.Write([]byte{byte(len(length))})
	c.Output.Write(length)
	
	for _, value := range data {		
		c.Output.Write([]byte{1, value})
	}
}

func (c *Compiler) List() {
	c.Output.Write([]byte{1, assembler.LIST})
}

func (c *Compiler) Load() {
	c.Output.Write([]byte{1, assembler.LOAD})
}

func (c *Compiler) Make() {
	c.Output.Write([]byte{1, assembler.MAKE})
}

func (c *Compiler) Put() {
	c.Output.Write([]byte{1, assembler.PUT})
}

func (c *Compiler) Pop() {
	c.Output.Write([]byte{1, assembler.POP})
}
func (c *Compiler) Get() {
	c.Output.Write([]byte{1, assembler.GET})
}
func (c *Compiler) Set() {
	c.Output.Write([]byte{1, assembler.SET})
}
func (c *Compiler) Size() {
	c.Output.Write([]byte{1, assembler.SIZE})
}
func (c *Compiler) Used() {
	c.Output.Write([]byte{1, assembler.USED})
}

func (c *Compiler) Free() {
	c.Output.Write([]byte{1, assembler.FREE})
}
func (c *Compiler) FreeList() {
	c.Output.Write([]byte{1, assembler.FREE_LIST})
}
func (c *Compiler) FreePipe() {
	c.Output.Write([]byte{1, assembler.FREE_PIPE})
}

func (c *Compiler) FreeType(t Type) {
	t.Base.Free(c)
}


func (c *Compiler) Heap() {
	c.Output.Write([]byte{1, assembler.HEAP})
}
func (c *Compiler) HeapList() {
	c.Output.Write([]byte{1, assembler.HEAP_LIST})
}
func (c *Compiler) HeapPipe() {
	c.Output.Write([]byte{1, assembler.HEAP_PIPE})
}

func (c *Compiler) Push(data string) {
	c.Output.Write([]byte{1, assembler.PUSH})
	WriteString(c.Output, data)
}

func (c *Compiler) PushList(data string) {
	c.Output.Write([]byte{1, assembler.PUSH_LIST})
	WriteString(c.Output, data)
}

func (c *Compiler) PushType(t Type, data string) {
	t.Base.Push(c, data)
}

func (c *Compiler) PushPipe(data string) {
	c.Output.Write([]byte{1, assembler.PUSH_PIPE})
	WriteString(c.Output, data)
}

func (c *Compiler) Pull(data string) {
	c.Output.Write([]byte{1, assembler.PULL})
	WriteString(c.Output, data)
}

func (c *Compiler) PullList(data string) {
	c.Output.Write([]byte{1, assembler.PULL_LIST})
	WriteString(c.Output, data)
}

func (c *Compiler) PullPipe(data string) {
	c.Output.Write([]byte{1, assembler.PULL_PIPE})
	WriteString(c.Output, data)
}

func (c *Compiler) PullType(t Type, data string) {
	t.Base.Pull(c, data)
}

func (c *Compiler) Drop() {
	c.Output.Write([]byte{1, assembler.DROP})
}

func (c *Compiler) DropList() {
	c.Output.Write([]byte{1, assembler.DROP_LIST})
}

func (c *Compiler) DropPipe() {
	c.Output.Write([]byte{1, assembler.DROP_PIPE})
}

func (c *Compiler) DropType(t Type) {
	t.Base.Drop(c)
}

func (c *Compiler) Name(data string) {
	c.Output.Write([]byte{1, assembler.NAME})
	WriteString(c.Output, data)
}

func (c *Compiler) NameList(data string) {
	c.Output.Write([]byte{1, assembler.NAME_LIST})
	WriteString(c.Output, data)
}

func (c *Compiler) NamePipe(data string) {
	c.Output.Write([]byte{1, assembler.NAME_PIPE})
	WriteString(c.Output, data)
}

func (c *Compiler) Copy() {
	c.Output.Write([]byte{1, assembler.COPY})
}

func (c *Compiler) CopyList() {
	c.Output.Write([]byte{1, assembler.COPY_LIST})
}

func (c *Compiler) CopyPipe() {
	c.Output.Write([]byte{1, assembler.COPY_PIPE})
}

func (c *Compiler) Swap() {
	c.Output.Write([]byte{1, assembler.SWAP})
}

func (c *Compiler) SwapList() {
	c.Output.Write([]byte{1, assembler.SWAP_LIST})
}

func (c *Compiler) SwapPipe() {
	c.Output.Write([]byte{1, assembler.SWAP_PIPE})
}

func (c *Compiler) Pipe() {
	c.Output.Write([]byte{1, assembler.PIPE})
}

func (c *Compiler) Wrap(data string) {
	c.Output.Write([]byte{1, assembler.WRAP})
	WriteString(c.Output, data)
}

func (c *Compiler) Open() {
	c.Output.Write([]byte{1, assembler.OPEN})
}

func (c *Compiler) Read() {
	c.Output.Write([]byte{1, assembler.READ})
}

func (c *Compiler) Send() {
	c.Output.Write([]byte{1, assembler.SEND})
}

func (c *Compiler) Stop() {
	c.Output.Write([]byte{1, assembler.STOP})
}

func (c *Compiler) Seek() {
	c.Output.Write([]byte{1, assembler.SEEK})
}

func (c *Compiler) Info() {
	c.Output.Write([]byte{1, assembler.INFO})
}

func (c *Compiler) Move() {
	c.Output.Write([]byte{1, assembler.MOVE})
}

func (c *Compiler) Int(value int64) {
	c.BigInt(Int(value))
}

func (c *Compiler) BigInt(value *big.Int) {
	c.Output.Write([]byte{1, assembler.INT})
	
	var bytes = value.Bytes()
	
	if len(bytes) > 255 {
		panic("NUMBER TOO BIG")
	}
	
	c.Output.Write([]byte{byte(len(bytes))})
	c.Output.Write(bytes)
}

func (c *Compiler) Add() {
	c.Output.Write([]byte{1, assembler.ADD})
}

func (c *Compiler) Sub() {
	c.Output.Write([]byte{1, assembler.SUB})
}

func (c *Compiler) Mul() {
	c.Output.Write([]byte{1, assembler.MUL})
}

func (c *Compiler) Div() {
	c.Output.Write([]byte{1, assembler.DIV})
}

func (c *Compiler) Mod() {
	c.Output.Write([]byte{1, assembler.MOD})
}

func (c *Compiler) Pow() {
	c.Output.Write([]byte{1, assembler.POW})
}

func (c *Compiler) Less() {
	c.Output.Write([]byte{1, assembler.LESS})
}
func (c *Compiler) More() {
	c.Output.Write([]byte{1, assembler.MORE})
}

func (c *Compiler) Same() {
	c.Output.Write([]byte{1, assembler.SAME})
}

func (c *Compiler) Flip() {
	c.Output.Write([]byte{1, assembler.FLIP})
}


func (c *Compiler) Native(target, data string) {
	c.Output.Write([]byte{1, assembler.NATIVE})
	WriteString(c.Output, target)
	WriteString(c.Output, data)
}
