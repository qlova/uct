package main

import "flag"
import "github.com/qlova/uct/src"
import (
	"os"
	"path/filepath"
)

func main() {
	flag.Parse()

	//Set the assembler.
	for _, asm := range uct.Assemblers {
		if *asm.Flag {
			uct.SetAssembler(asm)
			
			var extension = filepath.Ext(flag.Arg(0))
			var err error
			uct.Output, err = os.Create(flag.Arg(0)[0:len(flag.Arg(0))-len(extension)]+"."+asm.Ext)
			if err != nil {
				uct.Output = os.Stdout
			} else {
				os.Chmod(flag.Arg(0)[0:len(flag.Arg(0))-len(extension)]+"."+asm.Ext, 0755)
			}
		}
	}
	
	if !uct.AssemblerReady() {
		os.Stderr.Write([]byte("Please provide an assembler!"))
		os.Exit(1)
	}
	
	{
		

		uct.SetFileName(flag.Arg(0))
	}

	//Write any necessary headers.
	uct.Output.Write(uct.Header())
	
	//Reset the assembler's instruction count. 
	err := uct.Assemble(flag.Arg(0))
	if err != nil {
		os.Stderr.Write([]byte(err.Error()+"\n"))
		os.Exit(1)
	}
	//fmt.Println(aliases)
	
	//Write any necessary footers.
	uct.Output.Write(uct.Footer())
}
