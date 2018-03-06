package uct

import "flag"

var GoReserved = []string{
	"break",        "default",      "func",         "interface",    "select",
	"case",         "defer",        "go",           "map",          "struct",
	"chan",         "else",         "goto",         "package",      "switch",
	"const",        "fallthrough",  "if",           "range",        "type",
	"continue",     "for",          "import",       "return",       "var",
	"bool",			"byte", 		"len", 			"open", 		"file", 
	"close", 		"load", 		"copy",			"new",
}

//This is the Java compiler for uct.
var Go bool

func init() {
	flag.BoolVar(&Go, "go", false, "Target Go")

	RegisterAssembler(GoAssembly, &Go, "go", "//")

	for _, word := range GoReserved {
		GoAssembly[word] = Reserved()
	}
}

var GoAssembly = Assemblable{
	//Special commands.
	"HEADER": Instruction{
		Data:   "package main",
		Args: 1,
	},

	"FOOTER": Instruction{},

	"FILE": Instruction{
		Path: "stack.go",
	},

	"GO": Instruction{All:true},

	"NUMBER": is("NewNumber(%s)", 1),
	"BIG": 	is("NewNumberFromString(`%s`)", 1),
	"SIZE":   is("%s.Len()", 1),
	//"READY":   is(),
	"STRING": is("NewStringArray(%s)", 1),
	"ERRORS":  is("stack.ERROR", 1),
	
	//Experimental.
	"EVAL": is(`stack.Eval[stack.Grab().String()](stack)`),
	"EVALUATION": is(`stack.Eval["%s"] = %s`, 2),

	"SOFTWARE": Instruction{
		Data:   "func main() { stack := &Stack{}; stack.Init();",
		Indent: 1,
	},
	"EXIT": Instruction{
		Indented:    1,
		Data:        "stack.Exit(stack.ERROR.ToInt())}\n",
		Indent:      -1,
		Indentation: -1,
		Else: &Instruction{
			Data: "stack.Exit(stack.ERROR.ToInt())",
		},
	},

	"FUNCTION": is("func %s(stack *Stack) {", 1, 1),
	"RETURN": Instruction{
		Indented:    1,
		Data:        "}\n",
		Indent:      -1,
		Indentation: -1,
		Else: &Instruction{
			Data: "return",
		},
	},
	
	"SCOPE": is(`stack.Relay(&Pipe{Function:%s})`, 1),
	
	"EXE": is("%s.Exe(stack)", 1),

	"PUSH": is("stack.Push(%s)", 1),
	"PULL": is("%s := stack.Pull(); %s.Init()", 1),

	"PUT":   is("stack.Put(%s)", 1),
	"POP":   Instruction {
		Data: " ",
		Args: 1,
		Function: func(args []string) string {
			return args[0]+" := stack.Pop(); "+args[0]+".Add("+args[0]+", NewNumber(0))\n"
		},
	},
	"PLACE": is("stack.ActiveArray = %s", 1),
	"ARRAY":  is("var %s = &Array{}; stack.ActiveArray = %s", 1),
	"MAKE":  is("stack.Share(&Array{Small:make([]byte, stack.Pull().ToInt())})"),
	"RENAME": is("%s = stack.Grab()", 1),
	"RELOAD": is("%s = stack.Take()", 1),

	"SHARE": is("stack.Share(%s)", 1),
	"GRAB":  is("%s := stack.Grab(); %s.Init()", 1),

	"PIPE": is("var %s = stack.Pipe()", 1),
	"RELAY": is("stack.Relay(%s)", 1),
	"TAKE":  is("%s := stack.Take(); %s.Init()", 1),

	"GET": is("%s := stack.Get()", 1),
	"SET": is("stack.Set(%s)", 1),

	"VAR": is("var %s Number", 1),

	"INBOX":   is("stack.Share(<- stack.Inbox)"),
	"READYBOX":   is("stack.Push(NewNumber(len(stack.Inbox)))"),
	"OUTBOX":   is("stack.Outbox <- stack.Grab()"),

	"OPEN":   is("stack.Open()"),
	"DELETE": is("stack.Delete()"),
	"MOVE":   is("stack.Move()"),
	"EXECUTE":is("stack.Execute()"),
	"LINK":   is("stack.Link()"),
	"CONNECT":is("stack.Connect()"),
	"SLICE":  is("stack.Slice()"),
	"LOAD":   is("stack.Load()"),
	"OUT":    is("stack.Out()"),
	"STAT":   is("stack.Info()"),
	"IN":     is("stack.In()"),
	"STDOUT": is("stack.Stdout()"),
	"STDIN":  is("stack.Stdin()"),
	"HEAP":  is("stack.Heap()"),
	"HEAPIT":  is("stack.HeapIt()"),

	"CLOSE": is("%s.Close()", 1),

	"LOOP":   is("for {", 0, 1),
	"BREAK":  is("break"),
	"REPEAT": is("}", 0, -1, -1),

	"IF":   is("if %s.True()  {", 1, 1),
	"ELSE": is("} else {", 0, 0, -1),
	"END":  is("}", 0, -1, -1),

	"RUN":  is("%s(stack)", 1),
	"DATA": is("var %s *Array = %s;", 2),

	"FORK": is("go %s(stack.Copy())\n", 1),

	"ADD": is("%s.Add(%s, %s)", 3),
	"SUB": is("%s.Sub(%s, %s)", 3),
	"MUL": is("%s.Mul(%s, %s)", 3),
	"DIV": is("%s.Div(%s, %s)", 3),
	"MOD": is("%s.Mod(%s, %s)", 3),
	"POW": is("%s.Pow(%s, %s)", 3),

	"SLT": is("%s = %s.Slt(%s)", 3),
	"SEQ": is("%s = %s.Seq(%s)", 3),
	"SGE": is("%s = %s.Sge(%s)", 3),
	"SGT": is("%s = %s.Sgt(%s)", 3),
	"SNE": is("%s = %s.Sne(%s)", 3),
	"SLE": is("%s = %s.Sle(%s)", 3),

	"JOIN": is("%s = %s.Join(%s);", 3),
	"ERROR": is("stack.ERROR = %s;", 1),
}
