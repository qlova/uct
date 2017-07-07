package uct

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
	Pass, Check string
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
		if asm["FILE"].Path != "" {
		//if _, err := os.Stat(path.Dir(name)+"/"+asm["FILE"].Path); os.IsNotExist(err) {
			data, err := Asset("data/"+asm["FILE"].Path)
			if err == nil {
				ioutil.WriteFile(path.Dir(name)+"/"+asm["FILE"].Path, data, 0600)
			} else {
				fmt.Println(err.Error())
				os.Exit(1)
			}
		//}
		}
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
	"RUBY":true,
	"LUA":true,
	"JAVASCRIPT":true,
	"QML":true,
	"RUST": true,
	"BASH": true,
}

var Functions []string
var LastFunction string

func (asm Assemblable) Assemble(command string, args []string) ([]byte, error) {

	instruction, ok := asm[command]
	if !ok {
		if Languages[command] {
			return []byte(""), nil
		}
		return nil, errors.New("Unrecognised command! "+command)
	}
	
	defer func() {
		asm["last"] = Instruction{Data:command}
	}()
	
	if instruction.All {
		return []byte(asm.Indentation(instruction.Indentation)+strings.Join(args, " ")+"\n"), nil
	}
	if instruction.Data  == "" {
		return nil, errors.New("Bad command! "+command)
	} 
	if command == "NUMBER" || command == "SIZE" || command == "STRING" || command == "BIG" {
		if instruction.Function != nil {
			result := instruction.Function(args)
			return []byte(result), nil
		}
		if asm["BASH"].All && command == "SIZE" {
			return []byte(fmt.Sprintf(instruction.Data, LastFunction+args[0])), nil
		} else {
			return []byte(fmt.Sprintf(instruction.Data, args[0])), nil
		}
	}
	if command == "ERRORS" {
		return []byte(fmt.Sprintf(instruction.Data)), nil
	}
	
	//This is to fix the BASH circular reference bug.
	if command == "FUNCTION" {
		if len(args) > 0 {
			LastFunction = args[0]
		}
	}
	if command == "SOFTWARE" {
		LastFunction = "software"
	}
	
	//Especially for python, need to pass blank functions
	if instruction.Check != "" {
		//println(asm["last"].Data)
		if asm["last"].Data == instruction.Check {
			asm["last"] = Instruction{}
			a, err := asm.Assemble(command, args)
			return append([]byte(asm.Indentation(1)+instruction.Pass), a...), err
		}
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
			if instruct, ok := asm[arg[1:]]; ok {
				if instruct.Global {
					if asm["PREFIXGLOBALS"].Global {
						arg = "#$"+arg[1:]
					} else {
						//This is for rust. Inline data.
						//Not sure if it needs to be here, so commented it out.
						//Could cause a bug with indexing in rust.
						
						//args[i] = asm[arg].Data
					}
				} else {
					arg = "#u_"+arg[1:] //Fix for reserved words in languages being used as a length.
				}
			}
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
				} else if asm["RUST"].All {
					//This is for rust. Inline data.
					args[i] = asm[arg].Data
				}
			} else {
				args[i] = "u_"+args[i]
			}
		} else {
			//Hardcoded fix for BASH's "circular name reference" bug.
			if command != "FUNCTION" && command != "DATA" && command != "RUN" && command != "SCOPE" && asm["BASH"].All {
				args[i] = LastFunction+args[i]
			}
		}
	}
	
		
	//Experimental.
	if command == "FUNCTION" {
		Functions = append(Functions, args[0])
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
		asm[args[0]] = Instruction{Global:true, Data: args[1]}
	}
	
	var postpend string
	//Experimental.
	if _, ok := asm["EVALUATION"]; ok && command == "SOFTWARE"  {
		for _,name := range Functions {
			postpend += fmt.Sprintf(asm["EVALUATION"].Data+"\n", name, name)
		}
	}
	
	data := strings.Replace(instruction.Data, "\n", "\n"+asm.Indentation(instruction.Indentation), -1)
	
	if len(varaidic) > 0 {
		var value = fmt.Sprint(varaidic[0])
		//RUST workaraound.
		if asm["RUST"].All && (strings.Count(value, "NewStringArray") >= 1 || strings.Count(value, ".len().to_bigint().unwrap()") >= 1) {
			tmp++
		
			var name = "rust_workaround"+fmt.Sprint(tmp)
		
			return []byte(asm.Indentation(instruction.Indentation)+
			"let mut "+name+" = "+value+";\n"+
			asm.Indentation(instruction.Indentation)+
				fmt.Sprintf(data+"\n"+postpend, name)), nil
		}
	}
	
	if strings.Count(data, "%s") >= 1 {
	
		return []byte(asm.Indentation(instruction.Indentation)+
			fmt.Sprintf(data+"\n"+postpend, varaidic...)), nil
	} else {
		return []byte(asm.Indentation(instruction.Indentation)+
			fmt.Sprint(data+"\n"+postpend)), nil
	}
}

var tmp = 0
