 package java

import uct "github.com/qlova/uct/assembler"
import "math/big"

func init() {
	uct.RegisterTarget(uct.Target{
		Name: "Java",
		FileExtension: "java",
		Runtime: Runtime,
		
		Header: "",
		
		Assembler: func(uc *uct.Assembler) {

			if uc.String == "ERROR" {
				uc.String = "runtime.Error"
			}
			
			switch uc.Instruction {
				case uct.MAIN:
					uc.WriteLine("public static void main(String[] args) {")
					uc.IncreaseIndentation()
					uc.WriteLine("Runtime runtime = new Runtime();")

				case uct.EXIT:
					if uc.Indentation() > 1 {
						uc.WriteLine("System.exit(runtime.Error.Small);")
					} else {
						uc.DecreaseIndentation()
						uc.WriteLine("}}")
					}

				case uct.CODE:
					uc.WriteLine("static void "+uc.String+"(Runtime runtime) {")
					
					if data := OptimiseCode(uc.String); data != "" {
						uc.WriteLine(data)
						uc.DisableOutput()
					}
					
					uc.IncreaseIndentation()

				case uct.BACK:
					if uc.Indentation() > 1 {
						uc.WriteLine("if (true) return;")
					} else {
						uc.EnableOutput()
						uc.DecreaseIndentation()
						uc.WriteLine("}")
					}

				case uct.CALL:
					uc.WriteLine(uc.String+"(runtime);")
				
				//case uct.FORK:
					//uc.WriteLine("go "+uc.String+"(runtime)")

				case uct.LOOP:
					uc.WriteLine("while(true) {")
					uc.IncreaseIndentation()
					uc.WriteLine("if (false) break;")

				case uct.DONE:
					uc.WriteLine("if (true) break;")
				
				case uct.REDO:
					uc.DecreaseIndentation()
					uc.WriteLine("}")
					
				case uct.IF:
					uc.WriteLine("if (runtime.pull().isTrue()) {")
					uc.IncreaseIndentation()
				
				case uct.OR:
					uc.DecreaseIndentation()
					uc.WriteLine("} else {")
					uc.IncreaseIndentation()
				
				case uct.NO:
					uc.DecreaseIndentation()
					uc.WriteLine("}")

				case uct.DATA:
					uc.WriteString("static List "+uc.String+" = new List(new byte[]{")
					for i, value := range uc.Data {
						var r big.Int
						if value.Int64() > 127 {
							r.Sub(&value, big.NewInt(256))
						} else {
							r = value
						}
						uc.WriteString(r.String())
						if i < len(uc.Data)-1 {
							uc.WriteString(",")
						}
					}
					uc.WriteString("}, new Int[]{")
					for i, value := range uc.Data {
						uc.WriteString("new Int("+value.String()+")")
						if i < len(uc.Data)-1 {
							uc.WriteString(",")
						}
					}
					uc.WriteLine("});")

				case uct.LIST:
					uc.WriteLine("runtime.pushList(new List());")
				
				case uct.LOAD:
					uc.WriteLine("runtime.load();")
				
				case uct.MAKE:
					uc.WriteLine("runtime.pushList(new List(runtime.pull()));")

				case uct.PUT:
					uc.WriteLine("runtime.Lists[runtime.ListsPointer].put(runtime.pull());")
					
				case uct.POP:
					uc.WriteLine("runtime.push(runtime.Lists[runtime.ListsPointer].pop());")
					
				case uct.GET:
					uc.WriteLine("runtime.push(runtime.Lists[runtime.ListsPointer].get(runtime.pull()));")
				
				case uct.SET:
					uc.WriteLine("runtime.Lists[runtime.ListsPointer].set(runtime.Stack[runtime.StackPointer-1].Small, runtime.Stack[runtime.StackPointer]);")
					uc.WriteLine("runtime.StackPointer -= 2;")
				
				case uct.SIZE:
					uc.WriteLine("runtime.push(runtime.Lists[runtime.ListsPointer].size());")
				
				case uct.USED:
					uc.WriteLine("runtime.ListsPointer--;")
					
				case uct.PUSH:
					uc.WriteLine("runtime.push("+uc.String+");")
				
				case uct.PUSH_LIST:
					uc.WriteLine("runtime.pushList("+uc.String+");")
				
				case uct.PUSH_PIPE:
					uc.WriteLine("runtime.pushPipe("+uc.String+");")
					
				case uct.DROP, uct.FREE:
					uc.WriteLine("runtime.StackPointer--;")
				
				case uct.DROP_LIST, uct.FREE_LIST:
					uc.WriteLine("runtime.ListsPointer--;")
				
				case uct.DROP_PIPE, uct.FREE_PIPE:
					uc.WriteLine("runtime.PipesPointer--;")
				
				case uct.PULL:
					uc.WriteLine("Int "+uc.String+" = runtime.Stack[runtime.StackPointer];")
					uc.WriteLine("runtime.StackPointer--;")
				
				case uct.PULL_LIST:
					uc.WriteLine("List "+uc.String+" = runtime.Lists[runtime.ListsPointer];")
					uc.WriteLine("runtime.ListsPointer--;")
				
				case uct.PULL_PIPE:
					uc.WriteLine("Pipe "+uc.String+" = runtime.Pipes[runtime.PipesPointer];")
					uc.WriteLine("runtime.PipesPointer--;")
				
				case uct.NAME:
					uc.WriteLine(uc.String+" = runtime.Stack[runtime.StackPointer];")
					uc.WriteLine("runtime.StackPointer--;")
				
				case uct.NAME_LIST:
					uc.WriteLine(uc.String+" = runtime.Lists[runtime.ListsPointer];")
					uc.WriteLine("runtime.ListsPointer--;")
					
				case uct.NAME_PIPE:
					uc.WriteLine(uc.String+" = runtime.Pipes[runtime.PipesPointer];")
					uc.WriteLine("runtime.PipesPointer--;")
			
				case uct.HEAP_LIST:
					uc.WriteLine("runtime.heapList();")

				case uct.COPY:
					uc.WriteLine("runtime.push(runtime.Stack[runtime.StackPointer]);")

				case uct.COPY_LIST:
					uc.WriteLine("runtime.pushList(runtime.Lists[runtime.ListsPointer]);")
					
				case uct.COPY_PIPE:
					uc.WriteLine("runtime.pushPipe(runtime.Pipes[runtime.PipesPointer]);")
					
				
				case uct.SWAP:
					uc.WriteLine("runtime.swap();")

				case uct.SWAP_LIST:
					uc.WriteLine("runtime.swapList();")
					
				case uct.SWAP_PIPE:
					uc.WriteLine("runtime.swapPipe();")

					
				case uct.WRAP:
					uc.WriteLine("runtime.pushPipe(runtime.wrap(\""+uc.String+"\", new Object() {}.getClass().getEnclosingClass()));")
					
				case uct.OPEN:
					uc.WriteLine("runtime.open();")
				
				case uct.SEND:
					uc.WriteLine("runtime.send();")
				
				case uct.READ:
					uc.WriteLine("runtime.read();")
				
					
				case uct.INT:
					uc.WriteLine("runtime.push(new Int("+uc.Int.String()+"));")
				
				case uct.ADD:
					uc.WriteLine("runtime.add();")
				
				case uct.SUB:
					uc.WriteLine("runtime.sub();")	
				
				case uct.MUL:
					uc.WriteLine("runtime.mul();")
					
				case uct.DIV:
					uc.WriteLine("runtime.div();")
					
				case uct.MOD:
					uc.WriteLine("runtime.mod();")
				
				case uct.POW:
					uc.WriteLine("runtime.pow();")
				
				case uct.LESS:
					uc.WriteLine("runtime.less();")
					
				case uct.MORE:
					uc.WriteLine("runtime.more();")
					
				case uct.SAME:
					uc.WriteLine("runtime.same();")
					
				case uct.FLIP:
					uc.WriteLine("runtime.Stack[runtime.StackPointer] = runtime.Stack[runtime.StackPointer].flip();")
					
				case uct.NATIVE:
					uc.WriteLine(uc.String)
				
				default:
					println("WARNING "+uct.InstructionString(uc.Instruction)+" Not implemented")
			}
		},
	})
}
