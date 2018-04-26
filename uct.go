package main

import "github.com/qlova/uct/assembler"
import "fmt"
import "os"
import "path"

import _ "github.com/qlova/uct/targets/go"

func main() {
	
	if len(os.Args) != 3 {
		fmt.Println("[usage] uct file.u file.ext")
		return
	}
	
	var input = os.Args[1]
	var output = os.Args[2]
	
	var asm = new(assembler.Assembler)
	
	file, err := os.Open(input)
	if err != nil {
		fmt.Println("Could not open ", asm.Input, "! ", err.Error())
		return
	}
	asm.Input = append(asm.Input, file)
	
	asm.Output, err = os.Create(output)
	if err != nil {
		fmt.Println("Could not open ", asm.Input, "! ", err.Error())
		return
	}
	
	err = asm.Assemble(path.Ext(output)[1:])
	if err != nil {
		fmt.Println("Could not compile ", input, "! ", err.Error())
		return
	}
}
