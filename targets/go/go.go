package golang

import uct "github.com/qlova/uct/assembler"

func init() {
	uct.RegisterTarget(uct.Target{
		Name: "Go",
		FileExtension: "go",
		Runtime: Runtime,
		
		Header: "package main",
		
		Assembler: func(uc *uct.Assembler) {
			if uc.String == "ERROR" {
				uc.String = "runtime.Error"
			}
			
			switch uc.Instruction {
				case uct.MAIN:
					uc.WriteLine("func main() {")
					uc.IncreaseIndentation()
					uc.WriteLine("var runtime = new(Runtime)")
					uc.WriteLine("runtime.Init()")

				case uct.EXIT:
					if uc.Indentation() > 1 {
						uc.WriteLine("os.Exit(runtime.Error)")
					} else {
						uc.DecreaseIndentation()
						uc.WriteLine("Stdout.Flush()}")
					}

				case uct.CODE:
					uc.WriteLine("func "+uc.String+"(runtime *Runtime) {")
					
					if data := OptimiseCode(uc.String); data != "" {
						uc.WriteLine(data)
					}
					
					uc.IncreaseIndentation()

				case uct.BACK:
					if uc.Indentation() > 1 {
						uc.WriteLine("return")
					} else {
						uc.DecreaseIndentation()
						uc.WriteLine("}")
					}

				case uct.CALL:
					uc.WriteLine(uc.String+"(runtime)")
				
				case uct.FORK:
					uc.WriteLine("go "+uc.String+"(runtime)")

				case uct.LOOP:
					uc.WriteLine("for {")
					uc.IncreaseIndentation()

				case uct.DONE:
					uc.WriteLine("break")
				
				case uct.REDO:
					uc.DecreaseIndentation()
					uc.WriteLine("}")
					
				case uct.IF:
					uc.WriteLine("if runtime.Pull().True() {")
					uc.IncreaseIndentation()
				
				case uct.OR:
					uc.DecreaseIndentation()
					uc.WriteLine("} else {")
					uc.IncreaseIndentation()
				
				case uct.NO:
					uc.DecreaseIndentation()
					uc.WriteLine("}")

				case uct.DATA:
					uc.WriteString("var "+uc.String+" = &List{Bytes:[]byte{")
					for i, value := range uc.Data {
						uc.WriteString(value.String())
						if i < len(uc.Data)-1 {
							uc.WriteString(",")
						}
					}
					uc.WriteString("}, Mixed: []Int{")
					for i, value := range uc.Data {
						uc.WriteString("Int{Small:"+value.String()+"}")
						if i < len(uc.Data)-1 {
							uc.WriteString(",")
						}
					}
					uc.WriteLine("}}")

				case uct.LIST:
					uc.WriteLine("runtime.Lists = append(runtime.Lists, &List{})")
				
				case uct.LOAD:
					uc.WriteLine("runtime.Load()")
				
				case uct.MAKE:
					uc.WriteLine("runtime.Lists = append(runtime.Lists, &List{Mixed: make([]Int, runtime.Stack[len(runtime.Stack)-1].Small), Bytes: make([]byte, runtime.Stack[len(runtime.Stack)-1].Small)})")
					uc.WriteLine("runtime.Stack = runtime.Stack[:len(runtime.Stack)-1]")

				case uct.PUT:
					uc.WriteLine("runtime.Lists[len(runtime.Lists)-1].Put(runtime.Stack[len(runtime.Stack)-1])")
					uc.WriteLine("runtime.Stack = runtime.Stack[:len(runtime.Stack)-1]")
				
				case uct.POP:
					uc.WriteLine("runtime.Stack = append(runtime.Stack, runtime.Lists[len(runtime.Lists)-1].Pop())")
					
				case uct.GET:
					uc.WriteLine("if runtime.Lists[len(runtime.Lists)-1].Mixed == nil {")
						uc.WriteLine("runtime.Stack = append(runtime.Stack, Int{Small:int64(runtime.Lists[len(runtime.Lists)-1].Bytes[runtime.Pull().Small])})")
					uc.WriteLine("} else {")
						uc.WriteLine("runtime.Stack = append(runtime.Stack, runtime.Lists[len(runtime.Lists)-1].Mixed[runtime.Pull().Small])")
					uc.WriteLine("}")
				
				case uct.SET:
					uc.WriteLine("runtime.Lists[len(runtime.Lists)-1].Mixed[runtime.Stack[len(runtime.Stack)-2].Small] = runtime.Stack[len(runtime.Stack)-1]")
					uc.WriteLine("runtime.Stack = runtime.Stack[:len(runtime.Stack)-2]")
				
				case uct.SIZE:
					uc.WriteLine("runtime.Stack = append(runtime.Stack, Int{Small:int64(len(runtime.Lists[len(runtime.Lists)-1].Mixed))})")
				
				case uct.USED:
					uc.WriteLine("runtime.Lists = runtime.Lists[:len(runtime.Lists)-1]")
					
				case uct.PUSH:
					uc.WriteLine("runtime.Stack = append(runtime.Stack,"+uc.String+")")
				
				case uct.PUSH_LIST:
					uc.WriteLine("runtime.Lists = append(runtime.Lists, "+uc.String+")")
				
				case uct.PUSH_PIPE:
					uc.WriteLine("runtime.Pipes = append(runtime.Pipes, "+uc.String+")")
					
				case uct.DROP, uct.FREE:
					uc.WriteLine("runtime.Stack = runtime.Stack[:len(runtime.Stack)-1]")
				
				case uct.DROP_LIST, uct.FREE_LIST:
					uc.WriteLine("runtime.Lists = runtime.Lists[:len(runtime.Lists)-1]")
				
				case uct.DROP_PIPE, uct.FREE_PIPE:
					uc.WriteLine("runtime.Pipes = runtime.Pipes[:len(runtime.Pipes)-1]")
				
				case uct.PULL:
					uc.WriteLine("var "+uc.String+" = runtime.Stack[len(runtime.Stack)-1]")
					uc.WriteLine("runtime.Stack = runtime.Stack[:len(runtime.Stack)-1]")
					uc.WriteLine(uc.String+".Init()")
				
				case uct.PULL_LIST:
					uc.WriteLine("var "+uc.String+" = runtime.Lists[len(runtime.Lists)-1]")
					uc.WriteLine("runtime.Lists = runtime.Lists[:len(runtime.Lists)-1]")
					uc.WriteLine(uc.String+".Init()")
				
				case uct.PULL_PIPE:
					uc.WriteLine("var "+uc.String+" = runtime.Pipes[len(runtime.Pipes)-1]")
					uc.WriteLine("runtime.Pipes = runtime.Pipes[:len(runtime.Pipes)-1]")
					uc.WriteLine(uc.String+".Init()")
				
				case uct.NAME:
					uc.WriteLine(uc.String+" = runtime.Stack[len(runtime.Stack)-1]")
					uc.WriteLine("runtime.Stack = runtime.Stack[:len(runtime.Stack)-1]")
				
				case uct.NAME_LIST:
					uc.WriteLine(uc.String+" = runtime.Lists[len(runtime.Lists)-1]")
					uc.WriteLine("runtime.Lists = runtime.Lists[:len(runtime.Lists)-1]")
					
				case uct.NAME_PIPE:
					uc.WriteLine(uc.String+" = runtime.Pipes[len(runtime.Pipes)-1]")
					uc.WriteLine("runtime.Pipes = runtime.Pipes[:len(runtime.Pipes)-1]")
			
				case uct.HEAP_LIST:
					uc.WriteLine("runtime.HeapList()")

				case uct.COPY:
					uc.WriteLine("runtime.Stack = append(runtime.Stack, runtime.Stack[len(runtime.Stack)-1])")

				case uct.COPY_LIST:
					uc.WriteLine("runtime.Lists = append(runtime.Lists, runtime.Lists[len(runtime.Lists)-1])")
					
				case uct.COPY_PIPE:
					uc.WriteLine("runtime.Pipes = append(runtime.Pipes, runtime.Pipes[len(runtime.Pipes)-1])")
					
				
				case uct.SWAP:
					uc.WriteLine("runtime.Stack[len(runtime.Stack)-1], runtime.Stack[len(runtime.Stack)-2] = runtime.Stack[len(runtime.Stack)-2], runtime.Stack[len(runtime.Stack)-1]")

				case uct.SWAP_LIST:
					uc.WriteLine("runtime.Lists[len(runtime.Lists)-1], runtime.Lists[len(runtime.Lists)-2] = runtime.Lists[len(runtime.Lists)-2], runtime.Lists[len(runtime.Lists)-1]")
					
				case uct.SWAP_PIPE:
					uc.WriteLine("runtime.Pipes[len(runtime.Pipes)-1], runtime.Pipes[len(runtime.Pipes)-2] = runtime.Pipes[len(runtime.Pipes)-2], runtime.Pipes[len(runtime.Pipes)-1]")

					
				case uct.WRAP:
					uc.WriteLine("runtime.Pipes = append(runtime.Pipes, &WrappedFunction{Function:"+uc.String+"})")
					
				case uct.OPEN:
					uc.WriteLine("runtime.Open()")
				
				case uct.SEND:
					uc.WriteLine("runtime.Send()")
				
				case uct.READ:
					uc.WriteLine("runtime.Read()")
				
					
				case uct.INT:
					uc.WriteLine("runtime.Stack = append(runtime.Stack, Int{Small:"+uc.Int.String()+"})")
				
				case uct.ADD:
					uc.WriteLine("runtime.Add()")
				
				case uct.SUB:
					uc.WriteLine("runtime.Sub()")	
				
				case uct.MUL:
					uc.WriteLine("runtime.Mul()")
					
				case uct.DIV:
					uc.WriteLine("runtime.Div()")
					
				case uct.MOD:
					uc.WriteLine("runtime.Mod()")
				
				case uct.POW:
					uc.WriteLine("runtime.Pow()")
				
				case uct.LESS:
					uc.WriteLine("runtime.Less()")
					
				case uct.MORE:
					uc.WriteLine("runtime.More()")
					
				case uct.SAME:
					uc.WriteLine("runtime.Same()")
					
				case uct.FLIP:
					uc.WriteLine("runtime.Stack[len(runtime.Stack)-1] = runtime.Stack[len(runtime.Stack)-1].Flip()")
					
				case uct.NATIVE:
					uc.WriteLine(uc.String)
					
				
				default:
					println("WARNING "+uct.InstructionString(uc.Instruction)+" Not implemented")
			}
		},
	})
}
