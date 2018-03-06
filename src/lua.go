package uct

import "flag"

var LuaReserved = []string{
	"and",       "break",     "do",        "else",      "elseif",
    "end",       "false",     "for",       "function",  "if",
    "in" ,       "local",     "nil",       "not",       "or",
    "repeat",    "return",    "then",      "true",      "until",     "while",
}

//This is the Java compiler for uct.
var Lua bool

func init() {
	flag.BoolVar(&Lua, "lua", false, "Target Lua")
	
	RegisterAssembler(LuaAssembly, &Lua, "lua", "--")

	for _, word := range LuaReserved {
		LuaAssembly[word] = Reserved()
	}
}

var LuaAssembly = Assemblable{
	//Special commands.
	"HEADER": Instruction{
		Data:   `require "stack"`,
		Args: 1,
	},

	"FOOTER": Instruction{ Data: "main()"},

	"FILE": Instruction{
		Path: "stack.lua",
	},
	
	"LUA": Instruction{All:true},

	"NUMBER": is("bigint(%s)", 1),
	"BIG": 	is("bigint(\"%s\")", 1),
	"SIZE":   is("bigint(#%s)", 1),
	"STRING": is("tobytes(%s)", 1),
	"ERRORS":  is("stack.ERROR", 1),
	
	//"LINK":  is("stack:link()"),
	//"CONNECT":  is("stack:connect()"),
	//"SLICE":  is("stack:slice()"),

	"SOFTWARE": Instruction{
		Data:   "function main()\n\tlocal stack = NewStack()",
		Indent: 1,
	},

	"FUNCTION": is("function %s(stack)", 1, 1),
	
	"EXIT": Instruction{
		Indented:    1,
		Data:        "end",
		Indent:      -1,
		Indentation: -1,
		Else: &Instruction{
			Data: "os.exit(stack.ERROR)",
		},
	},
	
	"RETURN": Instruction{
		Indented:    1,
		Data:        "end\n",
		Indent:      -1,
		Indentation: -1,
		Else: &Instruction{
			Data: "return",
		},
	},
	
	"SCOPE": is(`stack:relay(NewPipe(%s))`, 1),
	
	"EXE": is("%s:exe(stack)", 1),

	//Optimised
	"PUSH": is("stack:push(%s)", 1),
	"PULL": is("local %s = stack:pull()", 1),
	"SHARE": is("stack:share(%s)", 1),
	"GRAB":  is("local %s = stack:grab()", 1),
	"PUT":   is("stack:put(%s)", 1),
	"POP":   is("local %s = stack:pop()", 1),
	"PLACE": is("stack.activearray = %s", 1),
	"ARRAY":  is("local %s = stack:array()", 1),
	"RENAME": is("%s = stack.grab()", 1),
	"EVAL": is("dostring(stack:grabstring()+'(stack)')"),
	
	"RELOAD": is("%s = stack:take()", 1),

	"RELAY": is("stack:relay(%s)", 1),
	"TAKE":  is("local %s = stack:take()", 1),

	"GET": is("local %s = stack:get()", 1),
	"SET": is("stack:set(%s)", 1),

	"VAR": is("local %s = bigint(0)", 1),

	"OPEN":   is("stack:open()"),
	"EXECUTE": is("stack:stdout()"),
	"DELETE": is("stack:delete()"),
	"LOAD":   is("stack:load()"),
	"OUT":    is("stack:out()"),
	"STAT":   is("stack:info()"),
	"IN":     is("stack:in()"),
	"STDOUT": is("stack:stdout()"),
	"STDIN":  is("stack:stdin()"),
	"HEAP":   is("stack:heap()"),
	"HEAPIT":   is("stack:heapit()"),
	
	"MAKE":   is("stack:share(make(tonumber(tostring(stack:pull()))))"),

	"CLOSE": is("%s:close()", 1),

	"LOOP":   is("while true do", 0, 1),
	"BREAK":  is("break"),
	"REPEAT": is("end", 0, -1, -1),

	"IF":   is("if %s ~= bigint(0) then", 1, 1),
	"ELSE": is("else", 0, 0, -1),
	"END":  is("end", 0, -1, -1),

	"RUN":  is("%s(stack)", 1),
	"DATA": is("local %s = %s", 2),

	//Threading.
	//"PIPE": is("%s = stack.pipe(stack.channel); stack.channel = stack.channel + 1", 1),
	//"FORK": is("stack.thread('%s')", 1),
	
	
	
	//"INBOX": is("while (stack.inbox.length <= 0) {} stack.share(stack.inbox.shift())", 0),
	//"OUTBOX": is("stack.outbox()", 0),

	"ADD": is("%s = %s + %s", 3),
	"SUB": is("%s = %s - %s", 3),
	"MUL": is("%s = %s * %s", 3),
	"DIV": is("%s = %s / %s", 3),
	"MOD": is("%s = %s %% %s", 3),
	"POW": is("%s = %s ^ %s", 3),

	"SLT": is("%s = %s < %s and bigint(1) or bigint(0)", 3),
	"SEQ": is("%s = %s == %s and bigint(1) or bigint(0)", 3),
	"SGE": is("%s = %s >= %s and bigint(1) or bigint(0)", 3),
	"SGT": is("%s = %s > %s and bigint(1) or bigint(0)", 3),
	"SNE": is("%s = %s ~= %s and bigint(1) or bigint(0)", 3),
	"SLE": is("%s = %s <= %s and bigint(1) or bigint(0)", 3),

	"JOIN": is("%s = join(%s, %s)", 3),
	"ERROR": is("stack.ERROR = %s", 1),
}
