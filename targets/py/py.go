package python

import uct "github.com/qlova/uct/assembler"

func init() {
	uct.RegisterTarget(uct.Target{
		Name: "Python",
		FileExtension: "py",
		Runtime: Runtime,
		
		Header: "#!/usr/bin/env python3",
		
		Assembler: func(uc *uct.Assembler) {
			if uc.String == "ERROR" {
				uc.String = "runtime.Error"
			}
			
			if uc.String == "GLOBAL" {
				uc.String = "runtime.Global"
			}
			
			switch uc.Instruction {
				case uct.MAIN:
					uc.WriteLine("if __name__ == \"__main__\":")
					uc.IncreaseIndentation()
					uc.WriteLine("runtime = Runtime()")

				case uct.EXIT:
					if uc.Indentation() > 1 {
						uc.WriteLine("sys.exit(runtime.Error)")
					} else {
						uc.DecreaseIndentation()
					}

				case uct.CODE:
					uc.WriteLine("def "+uc.String+"(runtime):")
					
					if data := OptimiseCode(uc.String); data != "" {
						uc.WriteLine(data)
					}
					
					uc.IncreaseIndentation()

				case uct.BACK:
					if uc.Indentation() > 1 {
						uc.WriteLine("return")
					} else {
						uc.DecreaseIndentation()
					}

				case uct.CALL:
					uc.WriteLine(uc.String+"(runtime)")
				
				//case uct.FORK:
					//uc.WriteLine("go "+uc.String+"(runtime)")

				case uct.LOOP:
					uc.WriteLine("while True:")
					uc.IncreaseIndentation()

				case uct.DONE:
					uc.WriteLine("break")
				
				case uct.REDO:
					uc.DecreaseIndentation()
					
				case uct.IF:
					uc.WriteLine("if runtime.Stack.pop() != 0:")
					uc.IncreaseIndentation()
				
				case uct.OR:
					uc.DecreaseIndentation()
					uc.WriteLine("else:")
					uc.IncreaseIndentation()
				
				case uct.NO:
					uc.DecreaseIndentation()

				case uct.DATA:
					uc.WriteString(uc.String+" = bytearray([")
					for i, value := range uc.Data {
						uc.WriteString(value.String())
						if i < len(uc.Data)-1 {
							uc.WriteString(",")
						}
					}
					uc.WriteLine("])")

				case uct.LIST:
					uc.WriteLine("runtime.Lists.append([])")
				
				case uct.LOAD:
					uc.WriteLine("runtime.load()")
				
				case uct.MAKE:
					uc.WriteLine("runtime.Lists.append([0]*runtime.Stack.pop())")

				case uct.PUT:
					uc.WriteLine("runtime.Lists[-1].append(runtime.Stack.pop())")
				
				case uct.POP:
					uc.WriteLine("runtime.Stack.append(runtime.Lists[-1].pop())")
					
				case uct.GET:
					uc.WriteLine("runtime.Stack.append(runtime.Lists[-1][runtime.Stack.pop()])")
				
				case uct.SET:
					uc.WriteLine("runtime.Lists[-1][runtime.Stack[-2]] = runtime.Stack[-1]")
					uc.WriteLine("runtime.Stack.pop()")
					uc.WriteLine("runtime.Stack.pop()")
				
				case uct.SIZE:
					uc.WriteLine("runtime.Stack.append(len(runtime.Lists[-1]))")
				
				case uct.USED:
					uc.WriteLine("runtime.Lists.pop()")
					
				case uct.PUSH:
					uc.WriteLine("runtime.Stack.append("+uc.String+")")
				
				case uct.PUSH_LIST:
					uc.WriteLine("runtime.Lists.append("+uc.String+")")
				
				case uct.PUSH_PIPE:
					uc.WriteLine("runtime.Pipes.append("+uc.String+")")
					
				case uct.DROP, uct.FREE:
					uc.WriteLine("runtime.Stack.pop()")
				
				case uct.DROP_LIST, uct.FREE_LIST:
					uc.WriteLine("runtime.Lists.pop()")
				
				case uct.DROP_PIPE, uct.FREE_PIPE:
					uc.WriteLine("runtime.Pipes.pop()")
				
				case uct.PULL:
					uc.WriteLine(uc.String+" = runtime.Stack.pop()")
				
				case uct.PULL_LIST:
					uc.WriteLine(uc.String+" = runtime.Lists.pop()")
				
				case uct.PULL_PIPE:
					uc.WriteLine(uc.String+" = runtime.Pipes.pop()")
				
				case uct.NAME:
					uc.WriteLine(uc.String+" = runtime.Stack.pop()")

				case uct.NAME_LIST:
					uc.WriteLine(uc.String+" = runtime.Lists.pop()")
					
				case uct.NAME_PIPE:
					uc.WriteLine(uc.String+" = runtime.Pipes.pop()")
			
				case uct.HEAP_LIST:
					uc.WriteLine("runtime.heaplist()")

				case uct.COPY:
					uc.WriteLine("runtime.Stack.append(runtime.Stack[-1])")

				case uct.COPY_LIST:
					uc.WriteLine("runtime.Lists.append(runtime.Lists[-1])")
					
				case uct.COPY_PIPE:
					uc.WriteLine("runtime.Pipes.append(runtime.Pipes[-1])")
				
				case uct.SWAP:
					uc.WriteLine("runtime.Stack[-1], runtime.Stack[-2] = runtime.Stack[-2], runtime.Stack[-1]")

				case uct.SWAP_LIST:
					uc.WriteLine("runtime.Lists[-1], runtime.Lists[-2] = runtime.Lists[-2], runtime.Lists[-1]")
					
				case uct.SWAP_PIPE:
					uc.WriteLine("runtime.Pipes[-1], runtime.Pipes[-2] = runtime.Pipes[-2], runtime.Pipes[-1]")
					
				case uct.WRAP:
					uc.WriteLine("runtime.Pipes.append(WrappedFunction("+uc.String+"))")
					
				case uct.OPEN:
					uc.WriteLine("runtime.open()")
				
				case uct.SEND:
					uc.WriteLine("runtime.send()")
				
				case uct.READ:
					uc.WriteLine("runtime.read()")
				
					
				case uct.INT:
					uc.WriteLine("runtime.Stack.append("+uc.Int.String()+")")
				
				case uct.ADD:
					uc.WriteLine("runtime.Stack.append(runtime.Stack.pop()+runtime.Stack.pop())")
				
				case uct.SUB:
					uc.WriteLine("runtime.Stack.append(-runtime.Stack.pop()+runtime.Stack.pop())")
				
				case uct.MUL:
					uc.WriteLine("runtime.Stack.append(runtime.Stack.pop()*runtime.Stack.pop())")
					
				case uct.DIV:
					uc.WriteLine("runtime.div()")
					
				case uct.MOD:
					uc.WriteLine("runtime.mod()")
				
				case uct.POW:
					uc.WriteLine("runtime.pow()")
				
				case uct.LESS:
					uc.WriteLine("runtime.Stack.append(int(runtime.Stack.pop() < runtime.Stack.pop()))")
					
				case uct.MORE:
					uc.WriteLine("runtime.Stack.append(int(runtime.Stack.pop() > runtime.Stack.pop()))")
					
				case uct.SAME:
					uc.WriteLine("runtime.Stack.append(int(runtime.Stack.pop() == runtime.Stack.pop()))")
					
				case uct.FLIP:
					uc.WriteLine("runtime.Stack.append(-runtime.Stack.pop())")
					
				case uct.NATIVE:
					uc.WriteLine(uc.String)
					
				
				default:
					println("WARNING "+uct.InstructionString(uc.Instruction)+" Not implemented")
			}
		},
	})
}
