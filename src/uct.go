//This was originally an assembler written by me, I am now bootstrapping Universal Code Translator with this.
package uct

import (
	"os"
	"bufio"
	"fmt"
	"strings"
	"io"
	"path/filepath"
	"errors"
	"strconv"
)

//Map of aliases, we store all replacements in this map.
var aliases map[string]string = make(map[string]string)
//Map of inline functions.
var inlines map[string]*InlineFunction = make(map[string]*InlineFunction)

var imported map[string]bool = make(map[string]bool)

//Alias hook, if we are creating an inline function then prepend the inline name.
func alias(name, value string) {
	if inlining != "" {
		aliases[inlining+"."+name] = value
	} else {
		aliases[name] = value
	}
}

//This is the inmemory structure of an inline function.
type InlineFunction struct {
	//These are the parameters the function takes.
	Aliases []string
	//This is the string of instructions.
	Instructions string
}

type NestedInlineFunction struct {
	*InlineFunction
	*bufio.Reader
	Name string
}

//An interface which assembles into machinecode.
type Assembler interface {

	//This should write any necessary headers to the binary.
	Header() []byte
	
	//This will assemble a instruction from text to binary.
	//All arguments are decimal numbers. $1 will be passed as 1.
	Assemble(string, []string) ([]byte, error)
	
	//This should write any necessary footers to the binary.
	Footer() []byte
	
	SetFileName(string)
}

var number = 0 				//The line number.
var instruction uint = 0 	//The instruction counter.
var assembler Assembler		//What assembler are we using?

//Output file.
var Output io.Writer
var comment string

var inlining string //Are we defining an inline function?

//TODO: make nested inlines.

var running int //Are we "running" an inline? LOL
var runningstack []NestedInlineFunction = make([]NestedInlineFunction, 20)

var inlinereader *bufio.Reader

type registeredAssembler struct {
	Flag *bool
	Ext string
	Comment string
	Assembler
}

var Assemblers []registeredAssembler

func RegisterAssembler(asm Assembler, flag *bool, ext, comment string) {
	Assemblers = append(Assemblers, registeredAssembler{Flag:flag, Ext:ext, Assembler: asm, Comment: comment})
}

func SetAssembler(asm registeredAssembler) {
	assembler = asm.Assembler
	comment = asm.Comment
}

func AssemblerReady() bool {
	return !(assembler == nil)
}

func SetFileName(name string) {
	assembler.SetFileName(name)
}

func Header() []byte {
	return assembler.Header()
}

func Footer() []byte {
	return assembler.Footer()
}

func Assemble(filename string) error {

	//Open main.s in the current directory and begin compilation.
	file, err := os.Open(filename)
	if err != nil {
		return errors.New("Could not find "+filename+" file!"+err.Error())
	}

	//Read the first line of the file and check the architecture.
	reader := bufio.NewReader(file)
	
	//Set our line number to 0.
	var number = 0

	//Loop through the lines of the file.
	for  {
	
		var line string
	
		if running == 0 {
			line, err = reader.ReadString('\n')
			if err == io.EOF {
				break
			} else if err != nil {
				return errors.New("Error reading file.. Is it courrupted? D: "+err.Error())
			}

			number++
				
		} else {
			
			line, err = runningstack[running].ReadString('\n')
			if err == io.EOF {
				
				running--
				continue
				
			} else if err != nil {
				return errors.New("Error reading file.. Is it courrupted? D: "+err.Error())
			}

		}
		
		//Figure out what the line is doing.
		//If the line is a comment, skip the line.
		if trim := strings.TrimSpace(line); len(trim) > 0 {
			if trim[0] == '#' {
				if comment != "" {
					Output.Write([]byte(strings.Replace(line, "#", comment+" ", 1)))
				}
				continue
			}
			line = trim
		} else {
			continue
		}

		//Split the line into space-sperated tokens.
		tokens := strings.Split(line, " ")
		if len(tokens) == 0 {
			continue
		}
		
		if inlining != "" {
			//If we are inlining a function and we hit .return then stop inlining and continue.
			if len(tokens[0]) > 0 && tokens[0] == ".return" {
				inlining = ""
				continue
			}
			inlines[inlining].Instructions += line+"\n"
			continue
		}

		//Detect labels.
		if len(tokens[0]) > 0 && tokens[0][len(tokens[0])-1] == ':' {
			
			var name string
			if filename != "main.s" {
				var extension = filepath.Ext(filename)
				name = filename[0:len(filename)-len(extension)]+"."
			}
			
			alias(name+tokens[0][:len(tokens[0])-1], fmt.Sprint(instruction))
			continue
		}

		//Resolve aliases.
		resolve:
		for i, token := range tokens {
			
			if token == "\\" && len(tokens) > i+1 && tokens[i+1] == "t" {
				token = "\t"
				tokens[i] = token
				 tokens[i+1] = ""
			}
			
			if token == "=" && len(tokens) > i+1 {
				tokens[i] = token+tokens[i+1]
				 tokens[i+1] = ""
			}
		
			if _, ok := Languages[tokens[0]]; tokens[0] == "DATA" || ok {
				continue
			}
			if len(token) == 0 {
				continue
			}
			if token != "," && token[len(token)-1] == ',' {
				token = token[:len(token)-1]
				tokens[i] = token
			}
			
			var extension = filepath.Ext(filename)
			var name = filepath.Base(filename[0:len(filename)-len(extension)])+"."
			
			var found bool
			
			if alias, ok := aliases[name+token]; ok {
				token = alias
				tokens[i] = alias
				found = true
				//output.Writeln("aliasing", token, alias)
			}
			if alias, ok := aliases[token]; ok {
				token = alias
				tokens[i] = alias
				found = true
				//output.Writeln("aliasing", token, alias)
			}
			
			if running > 0 {
				if alias, ok := aliases[runningstack[running].Name+"."+token]; ok {
					token = alias
					tokens[i] = alias
					found = true
					//output.Writeln("aliasing", token, alias)
				}
			}
			
			if inlining != "" {
				if alias, ok := aliases[inlining+"."+token]; ok {
					token = alias
					tokens[i] = alias
					found = true
					//output.Writeln("aliasing", token, alias)
				}
			}
			if len(token) > 0 && token[0] == '$' {
				token = token[1:]
				tokens[i] = token
			}
			
			if token == "=" && len(tokens) > i+1 && tokens[i+1] == "=" {
				token = "=="
				tokens[i] = token
			 	tokens[i+1] = ""
			 	continue
			}
			
			if len(token) > 2 && token[0] == '0' && token[1] == 'x' {
				var hex uint
				fmt.Sscanf(token[2:], "%x", &hex)
				//output.Writeln("converted hex", token ,"value to", hex)
				token = fmt.Sprint(hex)
				tokens[i] = token
			} else if len(token) > 2 && token[0] == '0' {
				var binary uint
				fmt.Sscanf(token[1:], "%b", &binary)
				//output.Writeln("converted hex", token ,"value to", hex)
				token = fmt.Sprint(binary)
				tokens[i] = token
			}
			
			//Blank aliases are ignored.
			if token == "_" {
				tokens = append(tokens[:i], tokens[i+1:]...)
				goto resolve
			}
			
			var v bool
			
			_, ok := Languages[tokens[0]]
			
			if !ok && len(token) > 1 && !found && tokens[0] != ".alias"  {
				for ii, c := range token {
					switch c {
						case '=':
							token = strings.Replace(token, string(c), "_eq_", -1)
							tokens[i] = token
							v = true
						case '+':
							token = strings.Replace(token, string(c), "_plus_", -1)
							tokens[i] = token
							v = true
						case '-':
							if _, err := strconv.Atoi(token); err != nil {
								token = strings.Replace(token, string(c), "_minus_", -1)
								tokens[i] = token
								v = true
							}
						case '/':
							token = strings.Replace(token, string(c), "_over_", -1)
							tokens[i] = token
							v = true
						case '*':
							token = strings.Replace(token, string(c), "_times_", -1)
							tokens[i] = token
							v = true
						case ')', '(':
							if ii != len(token)-1 {
								token = strings.Replace(token, string(c), "_l_", -1)
								tokens[i] = token
								v = true
							}
						case '<':
							 token = strings.Replace(token, string(c), "_lt_", -1)
							tokens[i] = token
							v = true
						case '>':
							token = strings.Replace(token, string(c), "_gt_", -1)
							tokens[i] = token
							v = true
						case '!':
							token = strings.Replace(token, string(c), "_not_", -1)
							tokens[i] = token
							v = true
						case '.':
							if i > 0 {		
								token = strings.Replace(token, ".", "_", 1)
								tokens[i] = token
							}
					}
				}
				if v {
					if token[0] == '#' {
						//token = "#v_"+token[1:]
						tokens[i] = token
					} else {
						//token = "v_"+token
						tokens[i] = token
					}
				}
				
				//Something to do with keeping things lowercase
				//I don't know why this is important??
							
				/*if i > 0 && token != "ERROR" && token != strings.ToLower(token) {
					if token[0] == '#' {
						token = "#l_"+strings.ToLower(token)[1:]
						tokens[i] = token
					} else {
						token = "l_"+strings.ToLower(token)
						tokens[i] = token
					}
				}*/
			}
		}

		switch tokens[0] {
		
			//Import another file and assemble it.
			case ".import":
				if len(tokens) != 2 {
					return  errors.New(fmt.Sprint(number)+": Import needs a filename.")
				}
				if !imported[tokens[1]+".u"] {
					err := Assemble(tokens[1]+".u")
					if err != nil {
						return errors.New(filename+":"+err.Error())
					}
					imported[tokens[1]+".u"] = true
				}
			
			//This is why quasm is so powerful! EVERYTHING IS AN ALIAS!
			case ".var", ".alias", ".const", ".global":
				if len(tokens) != 3 {
					return  errors.New(fmt.Sprint(number)+": Alias decleration needs a name and a value padded with single spaces.")
				}
				var name string
				if filename != "main.s" && tokens[0] != ".global" {
					var extension = filepath.Ext(filename)
					name = filepath.Base(filename[0:len(filename)-len(extension)])+"."
				}
				alias(name+tokens[1], tokens[2])
			//Create placeholder aliases!
			case ".blank":
				if len(tokens) != 2 {
					return  errors.New(fmt.Sprint(number)+": Blank decleration needs a name")
				}
				var name string
				if filename != "main.s" {
					var extension = filepath.Ext(filename)
					name = filepath.Base(filename[0:len(filename)-len(extension)])+"."
				}
				alias(name+tokens[1], "_")
			//RUN INLINE FUNCTIONS WITH PARAMETERS 8)
			case ".run", ".":
				if len(tokens) < 2 {
					return  errors.New(fmt.Sprint(number)+": Run requires a label.")
				}
				if _, ok := inlines[tokens[1]]; ok {
				
					runningstack[running+1]=NestedInlineFunction{
						Name: tokens[1],
						InlineFunction: inlines[tokens[1]], 
						Reader: bufio.NewReader(strings.NewReader(inlines[tokens[1]].Instructions)),
					}
					running++
					
					for i, aliasname := range inlines[tokens[1]].Aliases {
						if len(tokens)-1 < i {
							break
						}
						alias(tokens[1]+"."+aliasname, tokens[i])
					}
				} else {
					return  errors.New(fmt.Sprint(number)+": Could not find inline definition for "+tokens[1])
				}
			//Create inline functions!!! <3 <3 <3
			case ".inline":
				if len(tokens) < 2 {
					return  errors.New(fmt.Sprint(number)+": Inline decleration needs a name.")
				}
				inlining = tokens[1]
				inlines[inlining] = &InlineFunction{Aliases:tokens}
				
			//Use assembler to finally assemble commands to machine code.
			default:
				assembly, err := assembler.Assemble(tokens[0], tokens[1:])
				if err != nil {
					return  errors.New(fmt.Sprint(number)+": "+err.Error())
				}
				Output.Write(assembly)
				instruction++
		}
 	}
 	return nil
}
