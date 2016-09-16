package main

import "flag"

var RubyReserved = []string{
	"BEGIN",	"do",		"next",		"then",
    /*"END",*/	"else",		"nil",		"true",
	"alias",	"elsif",	"not",		"undef",
	"and",		"end",		"or",		"unless",
	"begin",	"ensure",	"redo",		"until",
	"break",	"false",	"rescue",	"when",
	"case",		"for",		"retry",	"while",
	"class",	"if",		"return",
	"def",		"in",		"self",		"__FILE__",
	"defined?",	"module",	"super",	"__LINE__",
	
	//Functions we use.
	"print", 	"open",

}

//This is the Java compiler for uct.
var Ruby bool

func init() {
	flag.BoolVar(&Ruby, "rb", false, "Target Ruby")

	RegisterAssembler(RubyAssembly, &Ruby, "rb", "#")

	for _, word := range RubyReserved {
		RubyAssembly[word] = Reserved()
	}
}

var RubyAssembly = Assemblable{
	//Special commands.
	"HEADER": Instruction{
		Data:   "load 'stack.rb'",
		Args: 1,
	},

	"FOOTER": Instruction{},

	"FILE": Instruction{
		Path: "stack.rb",
	},
	
	"RUBY": Instruction{All:true},

	"NUMBER": is("%s", 1),
	"BIG": 	is("%s", 1),
	"SIZE":   is("%s.length", 1),
	"STRING": is("%s.bytes.to_a", 1),
	
	"ERRORS":  is("stack.error", 1),
	
	"LINK":  is("stack.link"),
	"CONNECT":  is("stack.connect"),
	"SLICE":  is("stack.slice"),

	"SOFTWARE": Instruction{
		Data:   "stack = Stack.new\n",
	},
	"EXIT": Instruction{
		Data:   "exit(stack.error)",
	},
	
	"PREFIXGLOBALS": Instruction{
		Global:   true,
	},

	"FUNCTION": is("def %s(stack)", 1, 1),
	"RETURN": Instruction{
		Indented:    1,
		Data:        "end",
		Indent:      -1,
		Indentation: -1,
		Else: &Instruction{
			Data: "return",
		},
	},
	
	"SCOPE": is(`stack.relay(Pipe.new(method(:%s)))`, 1),
	
	"EXE": is("%s.exe(stack)", 1),

	"PUSH": is("stack.push %s", 1),
	"PULL": is("%s = stack.pull", 1),

	"PUT":   is("stack.put %s", 1),
	"POP":   is("%s = stack.pop", 1),
	"PLACE": is("stack.activearray = %s", 1),
	"ARRAY":  is("%s = stack.array", 1),
	"RENAME": is("%s = stack.activearray", 1),
	"RELOAD": is("%s = stack.take", 1),

	"SHARE": is("stack.share %s", 1),
	"GRAB":  is("%s = stack.grab", 1),

	"RELAY": is("stack.relay %s", 1),
	"TAKE":  is("%s = stack.take", 1),

	"GET": is("%s = stack.get", 1),
	"SET": is("stack.set %s", 1),

	"VAR": is("%s = 0", 1),

	"OPEN":   is("stack.openit"),
	"EXECUTE": is("stack.execute"),
	"DELETE": is("stack.delete"),
	"LOAD":   is("stack.loadit"),
	"OUT":    is("stack.out"),
	"STAT":   is("stack.info"),
	"IN":     is("stack.inn"),
	"STDOUT": is("stack.stdout"),
	"STDIN":  is("stack.stdin"),
	"HEAP":   is("stack.heap"),
	"HEAPIT":   is("stack.heapit"),

	"CLOSE": is("%s.close", 1),

	"LOOP":   is("while true", 0, 1),
	"BREAK":  is("break"),
	"REPEAT": is("end", 0, -1, -1),

	"IF":   is("if %s != 0", 1, 1),
	"ELSE": is("else", 0, 0, -1),
	"END":  is("end", 0, -1, -1),

	"RUN":  is("%s(stack)", 1),
	"DATA": is("$%s = %s", 2),

	"FORK": is("threading.Thread(target=%s, args=(stack.copy(),)).start()\n", 1),

	"ADD": is("%s = %s + %s", 3),
	"SUB": is("%s = %s - %s", 3),
	"MUL": is("%s = %s * %s", 3),
	"DIV": is("%s = stack.div(%s, %s)", 3),
	"MOD": is("%s = %s %% %s", 3),
	"POW": is("%s = %s**%s", 3),

	"SLT": is("%s = (%s <  %s) ? 1: 0", 3),
	"SEQ": is("%s = (%s == %s) ? 1: 0", 3),
	"SGE": is("%s = (%s >= %s) ? 1: 0", 3),
	"SGT": is("%s = (%s >  %s) ? 1: 0", 3),
	"SNE": is("%s = (%s != %s) ? 1: 0", 3),
	"SLE": is("%s = (%s <= %s) ? 1: 0", 3),

	"JOIN": is("%s = %s + %s", 3),
	"ERROR": is("stack.error = %s", 1),
}
