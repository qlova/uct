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

var RubyData = make(map[string]bool)

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
require 'securerandom'
	
N = []
N2 = []
F = []
F2 = []

$error = 0

def push(n)
	N << n
end

def pushit(n)
	F2 << n
end

def pushstring(n)
	N2 << n
end

def pushfunc(n)
	F << n
end

def pop()
	return N.pop
end

def popit()
	return F2.pop
end

def popstring()
	return N2.pop
end

def popfunc()
	return F.pop
end

def loadit() 
	text = popstring()
	result = []
	name = ""
	
	variable = ""	
	if text[0] == 36 and text.length > 0
		for i in 1..text.length-1
			name += text[i].chr
		end

		begin
			variable = ENV[name]
		rescue
			pushstring(result)
			return
		end
	else
		begin
			variable = ARGV[text[0]]
		rescue
			pushstring(result)
			return
		end
	end
	
	if !variable
		pushstring(result)
		return
	end
	
	
	for c in variable.split("")
		result << c.ord
	end
	pushstring(result)
end


def openit()

	filename = ""
	text = popstring()
	for i in 0..text.length-1
		filename += text[i].chr
	end
	
	file = [nil, nil]
	file[0] = filename
	begin
		file[1] = File.open(filename, "a+")
	rescue
		if File.directory?(filename)
			push(0)
			return file
		end
		push(-1)
		return file
	end
	push(0)
	return file
end

def out(file)
	text = popstring()
	
	if file[1] == nil
		if file[0][file[0].length-1] == "/"
			if File.directory?(file[0])
			
			else
				begin
					Dir.mkdir(file[0])
				rescue
					push(-1)
					return 
				end
			end
			push(0)
			return
		else
			if File.file?(file[0])
			
			else
				begin
					file[1] = File.open(file[0], 'w')
				rescue
					push(-1)
					return
				end
			end
			push(0)
			return
		end
	end
	
	for i in 0..text.length-1
		file[1].puts text[i].chr
	end
	push(0)
end

def inn(file)
	length = pop()
	for i in 1..length
		v = file[1].read(1)
		if v == ""
			push(-1000)
			return
		end
		push(v.ord)
	end
end

def close(file)
	if file[1]
		file[1].close
	end
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
		txt = STDIN.read(1)
		if not txt
			push(-1000)
			return
		end
		push(txt.ord)
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

def div(a, b)
	begin
		return a / b
	rescue
		if a == 0
			return (SecureRandom.random_number 255) + 1
		end
		return 0
	end
end

def mul(a, b)
	if a == 0 and b == 0
		return (SecureRandom.random_number 255) + 1
	end
	return a*b
end

def pow(a,b)
	if a == 0
		if b % 2 != 0
			return (SecureRandom.random_number 255) + 1
		end
		return 0
	end
	return a**b
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
		if RubyData[arg] {
			args[i] = "$"+arg
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
			case "byte", "end", "open", "close", "self":
				args[i] = "u_"+args[i]
			case "ERROR":
				args[i] = "$error"
		}
	} 
	switch command {
		case "ROUTINE":
			return []byte(""), nil
		case "SUBROUTINE":
			defer func() { g.Indentation++ }()
			return []byte("def "+args[0]+"()\n"), nil
		case "PUSH", "PUSHSTRING", "PUSHFUNC", "PUSHIT":
			var name string
			if command == "PUSHSTRING" {
				name = "string"
			}
			if command == "PUSHFUNC" {
				name = "func"
			}
			if command == "PUSHIT" {
				name = "it"
			}
			if len(args) == 1 {
				return []byte(g.indt()+"push"+name+"("+args[0]+")\n"), nil
			} else {
				return []byte(g.indt()+args[1]+".push("+args[0]+")\n"), nil
			}
		case "POP", "POPSTRING", "POPFUNC", "POPIT":
			var name string
			if command == "POPSTRING" {
				name = "string"
			}
			if command == "POPFUNC" {
				name = "func"
			}
			if command == "POPIT" {
				name = "it"
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
			
		case "FUNC":
			return []byte(g.indt()+args[0]+" = method(:"+args[1]+") \n"), nil
		case "EXE":
			return []byte(g.indt()+args[0]+".call() \n"), nil
			
		//IT stuff.
		case "OPEN":
			return []byte(g.indt()+""+args[0]+" = openit()\n"), nil
		case "OUT":
			return []byte(g.indt()+"out("+args[0]+")\n"), nil
		case "IN":
			return []byte(g.indt()+"inn("+args[0]+")\n"), nil
		case "CLOSE":
			return []byte(g.indt()+"close("+args[0]+")\n"), nil
		
		case "ERROR":
			return []byte(g.indt()+"$error ="+args[0]+"\n"), nil
			
		case "VAR":
			if len(args) == 1 {
				return []byte(g.indt()+args[0]+" = 0 \n"), nil
			} else {
				return []byte(g.indt()+args[0]+" = "+args[1]+" \n"), nil
			}
		case "STRING":
			return []byte(g.indt()+args[0]+" = [] \n"), nil
		case "STDOUT", "STDIN":
			return []byte(g.indt()+strings.ToLower(command)+"()\n"), nil
		case "LOAD":
			return []byte(g.indt()+"loadit()\n"), nil
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
			RubyData[args[0]] = true
			return []byte(g.indt()+"$"+args[0]+" = ["+args[1]+"] \n"), nil
		case "JOIN":
			return []byte(g.indt()+args[0]+" = "+args[1]+" + "+args[2]+"\n"), nil
		
		//Maths.
		case "ADD":
			return []byte(g.indt()+args[0]+" = "+args[1]+" + "+args[2]+"\n"), nil
		case "SUB":
			return []byte(g.indt()+args[0]+" = "+args[1]+" - "+args[2]+"\n"), nil
		case "MUL":
			return []byte(g.indt()+args[0]+" = mul("+args[1]+", "+args[2]+")\n"), nil
		case "DIV":
			return []byte(g.indt()+args[0]+" = div("+args[1]+", "+args[2]+")\n"), nil
		case "MOD":
			return []byte(g.indt()+args[0]+" = "+args[1]+" % "+args[2]+"\n"), nil
		case "POW":
			return []byte(g.indt()+args[0]+" = pow("+args[1]+", "+args[2]+")\n"), nil
			
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
