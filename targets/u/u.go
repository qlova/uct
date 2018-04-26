package u

import uct "github.com/qlova/uct/assembler"
import "fmt"

func init() {
	uct.RegisterTarget(uct.Target{
		Name: "Universal Assembly",
		FileExtension: "uct",
		
		Header: "",
		
		Assembler: func(uc *uct.Assembler) {
			switch uc.Instruction {
				case uct.LINK, uct.CALL, uct.PUSH, uct.PUSH_LIST, 
				uct.PUSH_PIPE, uct.PULL, uct.PULL_LIST, uct.PULL_PIPE, 
				uct.NAME, uct.NAME_LIST, uct.NAME_PIPE, uct.WRAP:
				
					fmt.Println(uct.InstructionString(uc.Instruction), " ", uc.String)
				
				case uct.INT:
					fmt.Println("INT ", uc.Int.String())
				
				case uct.DATA:
					fmt.Print("DATA ", uc.Int.String())
					for _, value := range uc.Data {
						fmt.Print(value.String(), " ")
					}
					fmt.Println()
				
				default:
					fmt.Println(uct.InstructionString(uc.Instruction))
			}
		},
	})
}

