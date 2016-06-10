package main

import "errors"
import "flag"
import "strconv"
import "strings"

//Go flag.
var Lua bool
func init() {
	flag.BoolVar(&Lua, "lua", false, "Target Lua")
	
	RegisterAssembler(new(LuaAssembler), &Lua, "lua", "--")
}

type LuaAssembler struct {
	Indentation int
}


func (g *LuaAssembler) indt(n ...int) string {
	if len(n) > 0 {
		return strings.Repeat("\t", int(g.Indentation)+n[0])
	} else {
		return strings.Repeat("\t", int(g.Indentation))
	}
}

func (g *LuaAssembler) Header() []byte {
	return []byte(
	`
	
N = {}
N2 = {}
F = {}
F2 = {}

ERROR = 0

function push(n)
	table.insert(N, n)
end

function pushit(n)
	table.insert(F2, n)
end

function pushstring(n)
	table.insert(N2, n)
end

function pushfunc(n)
	table.insert(F, n)
end

function pop()
	return table.remove(N)
end

function popit()
	return table.remove(F2)
end

function popstring()
	return table.remove(N2)
end

function popfunc()
	return table.remove(F)
end

function open()
	local filename = ""
	local text = popstring()
	for i = 1, #text do
		filename = filename..string.char(text[i])
	end
	
	file = io.open(filename, "a+")
	if file == nil then	
		local _, _, code = os.rename(filename, filename)
		if code == 2 then
			push(-1)
			return
		end
		push(0)
		return file
	else
		push(0)
		return file
	end
end

function out(file)
	local text = popstring()
	for i = 1, #text do
		out:write(string.char(text[i]))	
	end
end

function inn(file)
	local length = pop()
	for i = 1, length do
		push(string.byte(file:read(1)))
	end
end

function close(file)
	file:close()
end

function stdout()
	local text = popstring()
	for i = 1, #text do
		io.stdout:write(string.char(text[i]))	
	end
end

function stdin()
	local length = pop()
	for i = 1, length do
		push(string.byte(io.stdin:read(1)))
	end
end
		
function seq(a, b)
	if a == b then
		return 1
	end
	return 0
end

function sge(a, b)
	if a >= b then
		return 1
	end
	return 0
end

function sgt(a, b)
	if a > b then
		return 1
	end
	return 0
end

function sne(a, b)
	if a ~= b then
		return 1
	end
	return 0
end

function sle(a, b)
	if a <= b then
		return 1
	end
	return 0
end


function slt(a, b)
	if a < b then
		return 1
	end
	return 0
end

function join(a,b)
	local c = {}
	table.foreach(a,function(i,v)table.insert(c,v)end)
	table.foreach(b,function(i,v)table.insert(c,v)end)
	return c
end
`)
}

func (g *LuaAssembler) SetFileName(s string) {
}

func (g *LuaAssembler) Assemble(command string, args []string) ([]byte, error) {
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
			args[i] = arg
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
			case "end", "close", "open":
				args[i] = "u_"+args[i]
		}
	} 
	switch command {
		case "ROUTINE":
			return []byte(""), nil
		case "SUBROUTINE":
			defer func() { g.Indentation++ }()
			return []byte("function "+args[0]+"()\n"), nil
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
				return []byte(g.indt()+"table.insert("+args[1]+","+args[0]+")\n"), nil
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
				return []byte(g.indt()+"local "+args[0]+" = pop"+name+"()\n"), nil
			} else {
				return []byte(g.indt()+"local "+args[0]+" = table.remove("+args[1]+")\n"), nil
			}
		case "INDEX":
			return []byte(g.indt()+args[2]+" = "+args[0]+"["+args[1]+"+1]\n"), nil
		case "FUNC":
			return []byte(g.indt()+args[0]+" = "+args[1]+" \n"), nil
		case "EXE":
			return []byte(g.indt()+args[0]+"() \n"), nil
			
		//IT stuff.
		case "OPEN":
			return []byte(g.indt()+""+args[0]+" = open()\n"), nil
		case "OUT":
			return []byte(g.indt()+"out("+args[0]+")\n"), nil
		case "IN":
			return []byte(g.indt()+"inn("+args[0]+")\n"), nil
		case "CLOSE":
			return []byte(g.indt()+"close("+args[0]+")\n"), nil
			
		case "ERROR":
			return []byte(g.indt()+"ERROR = "+args[0]+"\n"), nil
		
		case "SET":
			return []byte(g.indt()+args[0]+"["+args[1]+"+1] = "+args[2]+"\n"), nil
		case "VAR":
			if len(args) == 1 {
				return []byte(g.indt()+"local "+args[0]+" = 0 \n"), nil
			} else {
				return []byte(g.indt()+"local "+args[0]+" = "+args[1]+" \n"), nil
			}
		case "STRING":
			return []byte(g.indt()+args[0]+" = {} \n"), nil
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
			return []byte(g.indt()+"if "+args[0]+" ~= 0 then\n"), nil
		case "RUN":
			return []byte(g.indt()+args[0]+"()\n"), nil
		case "ELSE":
			return []byte(g.indt(-1)+"else\n"), nil
		case "ELSEIF":
			return []byte(g.indt(-1)+"elsif "+args[0]+" ~= 0 then\n"), nil
		case "BREAK":
			return []byte(g.indt()+"break\n"), nil
		case "RETURN":
			return []byte(g.indt()+"return\n"), nil
		case "STRINGDATA":
			return []byte(g.indt()+args[0]+" = {"+args[1]+"} \n"), nil
		case "JOIN":
			return []byte(g.indt()+args[0]+" = join("+args[1]+","+args[2]+")\n"), nil
		
		//Maths.
		case "ADD":
			return []byte(g.indt()+args[0]+" = "+args[1]+" + "+args[2]+"\n"), nil
		case "SUB":
			return []byte(g.indt()+args[0]+" = "+args[1]+" - "+args[2]+"\n"), nil
		case "MUL":
			return []byte(g.indt()+args[0]+" = "+args[1]+" * "+args[2]+"\n"), nil
		case "DIV":
			return []byte(g.indt()+args[0]+" = math.floor("+args[1]+" / "+args[2]+")\n"), nil
		case "MOD":
			return []byte(g.indt()+args[0]+" = "+args[1]+" % "+args[2]+"\n"), nil
		case "POW":
			return []byte(g.indt()+args[0]+" = "+args[1]+" ^ "+args[2]+"\n"), nil
			
		case "SLT", "SEQ", "SGE", "SGT", "SNE", "SLE":
			return []byte(g.indt()+args[0]+" = "+strings.ToLower(command)+"("+args[1]+","+args[2]+")\n"), nil
		default:
			/*if expression {
				return []byte(g.indt()+command+" "+strings.Join(args, " ")+"\n"), nil
			}*/
	}
	return nil, errors.New("Unrecognised expression: "+command+" "+strings.Join(args, " "))

}

func (g *LuaAssembler) Footer() []byte {
	return []byte("")
}
