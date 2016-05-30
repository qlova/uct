package main

import "errors"
import "flag"
import "strconv"
import "strings"

//Go flag.
var Python bool
func init() {
	flag.BoolVar(&Python, "py", false, "Target Python")
	
	RegisterAssembler(new(PythonAssembler), &Python, "py", "#")
}

type PythonAssembler struct {
	Indentation int
}


func (g *PythonAssembler) indt(n ...int) string {
	if len(n) > 0 {
		return strings.Repeat("\t", int(g.Indentation)+n[0])
	} else {
		return strings.Repeat("\t", int(g.Indentation))
	}
}

func (g *PythonAssembler) Header() []byte {
	return []byte(
	`
import sys
	
N = []
N2 = []

def push(n):
	N.append(n)

def pushstring(n):
	N2.append(n)

def pop():
	return N.pop()

def popstring():
	return N2.pop()

def stdout():
	text = popstring()
	for i in range(0, len(text)):
		print(chr(text[i]), end="")

def stdin():
	length = pop()
	for i in range(0, length):
		push(ord(sys.stdin.read(1)))
		
def seq(a, b):
	if a == b:
		return 1
	return 0

def sge(a, b):
	if a >= b:
		return 1
	return 0

def sgt(a, b):
	if a > b:
		return 1
	return 0
	
def sne(a, b):
	if a != b:
		return 1
	return 0

def sle(a, b):
	if a <= b:
		return 1
	return 0

def slt(a, b):
	if a < b:
		return 1
	return 0
`)
}

func (g *PythonAssembler) SetFileName(s string) {
}

func (g *PythonAssembler) Assemble(command string, args []string) ([]byte, error) {
	//Format argument expressions.
	//var expression bool
	for i, arg := range args {
		if arg == "=" {
			//expression = true
			continue
		}
		if _, err := strconv.Atoi(arg); err == nil {
			args[i] = arg
			continue
		}
		if arg[0] == '#' {
			args[i] = "len("+arg[1:]+")"
		}
		if arg[0] == '"' {
			var newarg string
			var j = i
			arg = arg[1:]
			
			stringloop:
			arg = strings.Replace(arg, "\\n", "\n", -1)
			for _, v := range arg {
				if v == '"' {
					goto end
				}
				newarg += strconv.Itoa(int(v))+","
			}
			if len(arg) == 0 {
				goto end
			}
			newarg += strconv.Itoa(int(' '))+","
			j++
			//println(arg)
			arg = args[j]
			goto stringloop
			end:
			//println(newarg)
			args[i] = newarg
		}
		//RESERVED names in the language.
		switch arg {
			case "byte", "len":
				args[i] = "u_"+args[i]
		}
	} 
	switch command {
		case "ROUTINE":
			return []byte(""), nil
		case "SUBROUTINE":
			defer func() { g.Indentation++ }()
			return []byte("def "+args[0]+"():\n"), nil
		case "PUSH", "PUSHSTRING":
			var name string
			if command == "PUSHSTRING" {
				name = "string"
			}
			if len(args) == 1 {
				return []byte(g.indt()+"push"+name+"("+args[0]+")\n"), nil
			} else {
				return []byte(g.indt()+args[1]+".append("+args[0]+")\n"), nil
			}
		case "POP", "POPSTRING":
			var name string
			if command == "POPSTRING" {
				name = "string"
			}
			if len(args) == 1 {
				return []byte(g.indt()+args[0]+" = pop"+name+"()\n"), nil
			} else {
				return []byte(g.indt()+args[0]+" = "+args[1]+".pop()\n"), nil
			}
		case "INDEX":
			return []byte(g.indt()+args[2]+" = "+args[0]+"["+args[1]+"]\n"), nil
		case "SET":
			return []byte(g.indt()+args[0]+"["+args[1]+"] = "+args[2]+";\n"), nil
		case "VAR":
			if len(args) == 1 {
				return []byte(g.indt()+args[0]+" = 0 \n"), nil
			} else {
				return []byte(g.indt()+args[0]+" = "+args[1]+" \n"), nil
			}
		case "STRING":
			return []byte(g.indt()+args[0]+" = [] \n"), nil
		case "STDOUT":
			return []byte(g.indt()+"stdout()\n"), nil
		case "STDIN":
			return []byte(g.indt()+"stdin()\n"), nil
		case "LOOP":
			defer func() { g.Indentation++ }()
			return []byte(g.indt()+"while True:\n"), nil
		case "REPEAT", "END", "DONE":
			if g.Indentation > 0 {
				g.Indentation--
			}
			return []byte(g.indt()+"\n"), nil
		case "IF":
			defer func() { g.Indentation++ }()
			return []byte(g.indt()+"if "+args[0]+" != 0:\n"), nil
		case "RUN":
			return []byte(g.indt()+args[0]+"()\n"), nil
		case "ELSE":
			return []byte(g.indt(-1)+"else:\n"), nil
		case "ELSEIF":
			return []byte(g.indt(-1)+"elif "+args[0]+" != 0:\n"), nil
		case "BREAK":
			return []byte(g.indt()+"break\n"), nil
		case "RETURN":
			return []byte(g.indt()+"return\n"), nil
		case "STRINGDATA":
			return []byte(g.indt()+args[0]+" = ["+args[1]+"] \n"), nil
		case "JOIN":
			return []byte(g.indt()+args[0]+" = "+args[1]+" + "+args[2]+"\n"), nil
		
		//Maths.
		case "ADD":
			return []byte(g.indt()+args[0]+" = "+args[1]+" + "+args[2]+"\n"), nil
		case "SUB":
			return []byte(g.indt()+args[0]+" = "+args[1]+" - "+args[2]+"\n"), nil
		case "MUL":
			return []byte(g.indt()+args[0]+" = "+args[1]+" * "+args[2]+"\n"), nil
		case "DIV":
			return []byte(g.indt()+args[0]+" = "+args[1]+" // "+args[2]+"\n"), nil
		case "MOD":
			return []byte(g.indt()+args[0]+" = "+args[1]+" % "+args[2]+"\n"), nil
		case "POW":
			return []byte(g.indt()+args[0]+" = "+args[1]+" ** "+args[2]+"\n"), nil
			
		case "SLT", "SEQ", "SGE", "SGT", "SNE", "SLE":
			return []byte(g.indt()+args[0]+" = "+strings.ToLower(command)+"("+args[1]+","+args[2]+")\n"), nil
		default:
			/*if expression {
				return []byte(g.indt()+command+" "+strings.Join(args, " ")+"\n"), nil
			}*/
	}
	return nil, errors.New("Unrecognised expression: "+command+" "+strings.Join(args, " "))

}

func (g *PythonAssembler) Footer() []byte {
	return []byte("")
}
