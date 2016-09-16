package main

import "io/ioutil"
import "path"
import "path/filepath"
import "errors"
import "fmt"
import "strings"
import "strconv"
import "os"

type Assemblable map[string]Instruction

type Instruction struct {
	Data, Path string
	Indent, Indented, Args int
	
	All bool
	Global bool
	
	Indentation int
	Else *Instruction
	Function func(args []string) string
}

func is(s string, ns ...int) Instruction {
	var args, ind, indent int
	if len(ns) > 0 {
		args = ns[0]
	}
	if len(ns) > 1 {
		indent = ns[1]
	}
	if len(ns) > 2 {
		ind = ns[2]
	}
	return Instruction{Data:s, Args:args, Indent:indent, Indentation:ind}
}

func Reserved() Instruction {
	return Instruction{}
}

func (asm Assemblable) Header() []byte {
	b, err := asm.Assemble("HEADER", []string{asm["NAME"].Data})
	if err != nil {
		panic(err.Error())
	}
	return b
}

func (asm Assemblable) Indentation(n ...int) string {
	if _, ok := asm["INDENT"]; !ok {
		return ""
	}
	if len(n) > 0 {
		if int(asm["INDENT"].Indent)+n[0] < 0 {
			return ""
		}
		return strings.Repeat("\t", int(asm["INDENT"].Indent)+n[0])
	} else {
		return strings.Repeat("\t", int(asm["INDENT"].Indent))
	}
}

func (asm Assemblable) SetFileName(name string) {
	if asm["FILE"].Data == "" {
		//if _, err := os.Stat(path.Dir(name)+"/"+asm["FILE"].Path); os.IsNotExist(err) {
			data, err := Asset("data/"+asm["FILE"].Path)
			if err == nil {
				ioutil.WriteFile(path.Dir(name)+"/"+asm["FILE"].Path, data, 0600)
			} else {
				fmt.Println(err.Error())
				os.Exit(1)
			}
		//}
	} else {
		println("[depreciated]")
		if _, err := os.Stat(path.Dir(name)+asm["FILE"].Path); os.IsNotExist(err) {
			ioutil.WriteFile(path.Dir(name)+asm["FILE"].Path, []byte(asm["FILE"].Data), 0600)
		}
	}
	var extension = filepath.Ext(name)
	name = filepath.Base(name[0:len(name)-len(extension)])
	asm["NAME"] = Instruction{Data:path.Base(name)}
}

func (asm Assemblable) Footer() []byte {
	b, _ := asm.Assemble("FOOTER", nil)
	return b
}

var Languages = map[string]bool {
	"JAVA":true,
	"GO":true,
	"PYTHON":true,
}

func (asm Assemblable) Assemble(command string, args []string) ([]byte, error) {
	instruction, ok := asm[command]
	if !ok {
		if Languages[command] {
			return []byte(""), nil
		}
		return nil, errors.New("Unrecognised command! "+command)
	}
	if instruction.All {
		return []byte(strings.Join(args, " ")+"\n"), nil
	}
	if instruction.Data  == "" {
		return nil, errors.New("Bad command! "+command)
	} 
	if command == "NUMBER" || command == "SIZE" || command == "STRING" {
		return []byte(fmt.Sprintf(instruction.Data, args[0])), nil
	}
	if command == "ERRORS" {
		return []byte(fmt.Sprintf(instruction.Data)), nil
	}
	for i, arg := range args {
		if len(arg) == 0 {
			break
		}
	
		var b []byte
		if _, err := strconv.Atoi(arg); err == nil {
			b, err = asm.Assemble("NUMBER", []string{arg})
			if err != nil {
				panic(err.Error())
			}
			args[i] = string(b)
			continue
		}
		
		if _, err := strconv.Atoi(arg); err == nil {
			b, err = asm.Assemble("NUMBER", []string{arg})
			if err != nil {
				panic(err.Error())
			}
			args[i] = string(b)
			continue
		} else if _, err := strconv.Atoi(string(arg[0])); err == nil || arg[0] == '-' {
			b, err = asm.Assemble("BIG", []string{arg})
			if err != nil {
				panic(err.Error())
			}
			args[i] = string(b)
			continue
		}
	
		if arg[0] == '#' {
			b, _ = asm.Assemble("SIZE", []string{arg[1:]})
			args[i] = string(b)
			continue
		}
		if arg[0] == '"' {
			b, _ = asm.Assemble("STRING", []string{strings.Join(args[i:], " ")})
			args[i] = string(b)
			args = args[:i+1]
			break
		}
		if arg == "ERROR" {
			b, _ = asm.Assemble("ERRORS", []string{arg[1:]})
			args[i] = string(b)
			continue
		}
		
		
		if instruct, ok := asm[arg]; ok {
			if instruct.Global {
				if asm["PREFIXGLOBALS"].Global {
					args[i] = "$"+args[i]
				}
			} else {
				args[i] = "u_"+args[i]
			}
		}
	}
	
	if indent, ok := asm["INDENT"]; ok {
		if (instruction.Indented == -1 && 0 != indent.Indent) ||
		(instruction.Indented != 0 && instruction.Indented != indent.Indent) {
			instruction = *instruction.Else
		}
	}
	
	if len(args) != instruction.Args {
		return nil, errors.New(command+" Argument count mismatch! "+fmt.Sprintf("%v != %v args:%v", len(args), instruction.Args,args))
	}
	
	if instruction.Function != nil {
		result := instruction.Function(args)
		return []byte(result), nil
	}
	
	varaidic := make([]interface {}, len(args))
	for i, v := range args {
		varaidic[i] = v
	}
	if len(args) > 0 && strings.Count(instruction.Data, "%s") > instruction.Args {
		varaidic = append(varaidic, varaidic[len(varaidic)-1])
	}
	
	if instruction.Indent != 0 {
		if _, ok := asm["INDENT"]; !ok {
			asm["INDENT"] = Instruction{}
		}
		defer func() {
			asm["INDENT"] = Instruction{Indent: asm["INDENT"].Indent+instruction.Indent}
			if asm["INDENT"].Indent < 0 {
				asm["INDENT"] = Instruction{Indent: 0}
			}
		}()
		
	}
	
	//Keep a record of globals.
	if command == "DATA" {
		asm[args[0]] = Instruction{Global:true}
	}
	
	if strings.Count(instruction.Data, "%s") >= 1 {
	
		return []byte(asm.Indentation(instruction.Indentation)+
			fmt.Sprintf(instruction.Data+"\n", varaidic...)), nil
	} else {
		return []byte(asm.Indentation(instruction.Indentation)+
			fmt.Sprint(instruction.Data+"\n")), nil
	}
}
