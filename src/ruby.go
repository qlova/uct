package main

import "errors"
import "flag"
import "strconv"
import "strings"

//Go flag.
var Ruby bool
func init() {
	flag.BoolVar(&Ruby, "rb", false, "Target Ruby")
	
	RegisterAssembler(new(RubyAssembler), &Ruby, "rb", "#")
}

type RubyAssembler struct {
	Indentation int
}


func (g *RubyAssembler) indt(n ...int) string {
	if len(n) > 0 {
		return strings.Repeat("\t", int(g.Indentation)+n[0])
	} else {
		return strings.Repeat("\t", int(g.Indentation))
	}
}

func (g *RubyAssembler) Header() []byte {
	return []byte(
	`
	
N = []
N2 = []

def push(n)
	N << n
end

def pushstring(n)
	N2 << n
end

def pop()
	return N.pop
end

def popstring()
	return N2.pop
end

def stdout()
	text = popstring()
	for i in 0..text.length-1
		print(text[i].chr)
	end
end

def stdin()
	length = pop()
	for i in 1..length
		push(STDIN.read(1).ord)
	end
end
		
def seq(a, b)
	if a == b
		return 1
	end
	return 0
end

def sge(a, b)
	if a >= b
		return 1
	end
	return 0
end

def sgt(a, b)
	if a > b
		return 1
	end
	return 0
end

def sne(a, b)
	if a != b
		return 1
	end
	return 0
end

def sle(a, b)
	if a <= b
		return 1
	end
	return 0
end


def slt(a, b)
	if a < b
		return 1
	end
	return 0
end
`)
}

func (g *RubyAssembler) SetFileName(s string) {
}

func (g *RubyAssembler) Assemble(command string, args []string) ([]byte, error) {
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
			args[i] = arg[1:]+".length"
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
			case "byte", "end":
				args[i] = "u_"+args[i]
		}
	} 
	switch command {
		case "ROUTINE":
			return []byte(""), nil
		case "SUBROUTINE":
			defer func() { g.Indentation++ }()
			return []byte("def "+args[0]+"()\n"), nil
		case "PUSH", "PUSHSTRING":
			var name string
			if command == "PUSHSTRING" {
				name = "string"
			}
			if len(args) == 1 {
				return []byte(g.indt()+"push"+name+"("+args[0]+")\n"), nil
			} else {
				return []byte(g.indt()+args[1]+".push("+args[0]+")\n"), nil
			}
		case "POP", "POPSTRING":
			var name string
			if command == "POPSTRING" {
				name = "string"
			}
			if len(args) == 1 {
				return []byte(g.indt()+args[0]+" = pop"+name+"()\n"), nil
			} else {
				return []byte(g.indt()+args[0]+" = "+args[1]+".pop\n"), nil
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
			return []byte(g.indt()+"while true do\n"), nil
		case "REPEAT", "END", "DONE":
			if g.Indentation > 0 {
				g.Indentation--
				return []byte(g.indt()+"end\n"), nil
			}
			return []byte(""), nil
		case "IF":
			defer func() { g.Indentation++ }()
			return []byte(g.indt()+"if "+args[0]+" != 0\n"), nil
		case "RUN":
			return []byte(g.indt()+args[0]+"()\n"), nil
		case "ELSE":
			return []byte(g.indt(-1)+"else\n"), nil
		case "ELSEIF":
			return []byte(g.indt(-1)+"elsif "+args[0]+" != 0:\n"), nil
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
			return []byte(g.indt()+args[0]+" = "+args[1]+" / "+args[2]+"\n"), nil
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

func (g *RubyAssembler) Footer() []byte {
	return []byte("")
}
