package uct

import "flag"

var CSReserved = []string{
	//Figure these out.
	"abstract",	"as",		"base",			"bool",		"break",		"byte",
	"case",		"catch",	"char",			"checked",	"class",		"const",
	"continue",	"decimal",	"default",		"delegate",	"do",			"double",
	"else",		"enum",		"event",		"explicit",	"extern",		"finally",
	"fixed",	"float",	"for",			"foreach",	"goto",			"if",
	"implicit",	"in",		"int",			"interface","internal",		"is",
	"lock",		"long",		"namespace",	"new",		"null",			"object",
	"operator",	"out",		"override",		"params",	"private",		"s",
	"public",	"readonly",	"ref",			"return",	"sbyte",		"sealed",
	"short",	"sizeof",	"salloc",		"static",	"string",		"struct",
	"switch",	"this",		"throw",		"try",		"typeof",		"unit",
	"ulong",	"unchecked","unsafe",		"ushort",	"using",		"virtual",
	"void",		"volatile",	"while",		"FALSE",	"TRUE",
}

//This is the Java compiler for uct.
var CS bool

func init() {
	flag.BoolVar(&CS, "cs", false, "Target C#")

	RegisterAssembler(CSAssembly, &CS, "cs", "//")

	for _, word := range CSReserved {
		CSAssembly[word] = Reserved()
	}
}

var CSAssembly = Assemblable{
	//Special commands.
	"HEADER": Instruction{
		Data:   "using System;\nusing System.Reflection;\nusing System.Numerics;\nclass %s {",
		Indent: 1,
		Args:   1,
	},

	"FOOTER": Instruction{
		Data:        "}",
		Indent:      -1,
		Indentation: -1,
	},

	"FILE": Instruction{
		Path: "stack.cs",
	},
	
	"CSHARP": Instruction{All:true},
	
	//"EVAL": is("Class[] cArg = new Class[1]; cArg[0] = Stack.class; try { new Object() { }.getClass().getEnclosingClass().getDeclaredMethod(s.grab().String(), cArg).invoke(null, s); } catch (Exception e) { throw new RuntimeException(e); }"),

	"NUMBER": is("new BigInteger(%s)", 1),
	"BIG": is("new BigInteger(\"%s\")", 1),
	"SIZE":   is("%s.size()", 1),
	"STRING": is("new stack.Array(%s)", 1),
	"ERRORS":  is("s.ERROR", 1),

	"SOFTWARE": Instruction{
		Data:   "public static void Main(string[] args) {\n stack s = new stack(); stack.ARGS=args;",
		Indent: 1,
	},
	"EXIT": Instruction{
		Indented:    2,
		Data:        "}\n",
		Indent:      -1,
		Indentation: -1,
		Else: &Instruction{
			Data: "System.Exit(((int)s.ERROR));",
		},
	},

	"FUNCTION": is("static void %s(stack s) {", 1, 1),
	"RETURN": Instruction{
		Indented:    2,
		Data:        "}\n",
		Indent:      -1,
		Indentation: -1,
		Else: &Instruction{
			Data: "return;",
		},
	},
	
	"SCOPE": is(`s.relay(new stack.Pipe(typeof(Functions).GetTypeInfo().GetDeclaredMethod("%s")));`, 1),
	
	"EXE": is("%s.Func.Invoke(null, new Object[]{s});", 1),

	"PUSH": is("s.push(%s);", 1),
	"PULL": is("BigInteger %s = s.pull();", 1),

	"PUT":   is("s.put(%s);", 1),
	"POP":   is("BigInteger %s = s.pop();", 1),
	"PLACE": is("s.place(%s);", 1),

	"ARRAY":  is("stack.Array %s = s.array();", 1),
	"MAKE":  is("s.share(new stack.Array(((int)s.pull())));"),
	"RENAME": is("%s = s.grab();", 1),
	
	"RELOAD": is("%s = s.take();", 1),

	"SHARE": is("s.share(%s);", 1),
	"GRAB":  is("stack.Array %s = s.grab();", 1),

	"RELAY": is("s.relay(%s);", 1),
	"TAKE":  is("stack.Pipe %s = s.take();", 1),

	"GET": is("BigInteger %s = s.gets();", 1),
	"SET": is("s.sets(%s);", 1),

	"VAR": is("BigInteger %s = new BigInteger();", 1),

	"OPEN":   is("s.open();"),
	"DELETE":   is("s.delete();"),
	"EXECUTE":   is("s.execute();"),
	"LOAD":   is("s.load();"),
	"OUT":    is("s.out();"),
	"STAT":   is("s.info();"),
	"IN":     is("s.in();"),
	"STDOUT": is("s.stdout();"),
	"STDIN":  is("s.stdin();"),
	"HEAP":   is("s.heap();"),
	"LINK":   is("s.link();"),
	"CONNECT":   is("s.connect();"),
	"SLICE":   is("s.slice();"),

	"CLOSE": is("%s.close();", 1),

	"LOOP":   is("while (true) {", 0, 1),
	"BREAK":  is("break;"),
	"REPEAT": is("}", 0, -1, -1),

	"IF":   is("if (%s.CompareTo(new BigInteger(0)) != 0 ) {", 1, 1),
	"ELSE": is("} else {", 0, 0, -1),
	"END":  is("}", 0, -1, -1),

	"RUN":  is("%s(s);", 1),
	"DATA": is_data("static stack.Array %s = new stack.Array(", "%s", ",", ")"),

	"FORK": is("{ Stack s = s.copy(); Stack.ThreadPool.execute(() -> %s(s)); }\n", 1),

	"ADD": is("%s = %s + %s;", 3),
	"SUB": is("%s = %s - %s;", 3),
	"MUL": is("%s = stack.mul(%s, %s);", 3),
	"DIV": is("%s = stack.div(%s, %s);", 3),
	"MOD": is("%s = stack.mod(%s, %s);", 3),
	"POW": is("%s = stack.pow(%s, %s);", 3),

	"SLT": is("%s = stack.slt(%s, %s);", 3),
	"SEQ": is("%s = stack.seq(%s, %s);", 3),
	"SGE": is("%s = stack.sge(%s, %s);", 3),
	"SGT": is("%s = stack.sgt(%s, %s);", 3),
	"SNE": is("%s = stack.sne(%s, %s);", 3),
	"SLE": is("%s = stack.sle(%s, %s);", 3),

	"JOIN": is("%s = %s.join(%s);", 3),
	"ERROR": is("s.ERROR = %s;", 1),
}
