package main

import "flag"

var JavaReserved = []string{
	"abstract", "continue", "for", "new", "switch", "assert", "default",
	"goto", "package", "synchronized", "boolean", "do", "if", "private",
	"this", "break", "double", "implements", "protected", "throw",
	"byte", "else", "import", "public", "throws", "case", "enum", "instanceof",
	"return", "transient", "catch", "extends", "int", "short", "try",
	"char", "final", "interface", "static", "void", "class", "finally",
	"long", "strictfp", "volatile", "const", "float", "native", "super", "while",
}

//This is the Java compiler for uct.
var Java bool

func init() {
	flag.BoolVar(&Java, "java", false, "Target Java")

	RegisterAssembler(JavaAssembly, &Java, "java", "//")

	for _, word := range JavaReserved {
		JavaAssembly[word] = Reserved()
	}
}

var JavaAssembly = Assemblable{
	//Special commands.
	"HEADER": Instruction{
		Data:   "public class %s {",
		Indent: 1,
		Args:   1,
	},

	"FOOTER": Instruction{
		Data:        "}",
		Indent:      -1,
		Indentation: -1,
	},

	"FILE": Instruction{
		Path: "Stack.java",
	},
	
	"JAVA": Instruction{All:true},
	
	"EVAL": is("Class[] cArg = new Class[1]; cArg[0] = Stack.class; try { new Object() { }.getClass().getEnclosingClass().getDeclaredMethod(stack.grab().String(), cArg).invoke(null, stack); } catch (Exception e) { throw new RuntimeException(e); }"),

	"NUMBER": is("new Stack.Number(%s)", 1),
	"BIG": is("new Stack.Number(%s)", 1),
	"SIZE":   is("%s.size()", 1),
	"STRING": is("new Stack.Array(%s)", 1),
	"ERRORS":  is("stack.ERROR", 1),

	"SOFTWARE": Instruction{
		Data:   "public static void main(String[] args) { Stack stack = new Stack(); stack.Arguments = args;",
		Indent: 1,
	},
	"EXIT": Instruction{
		Indented:    2,
		Data:        "}\n",
		Indent:      -1,
		Indentation: -1,
		Else: &Instruction{
			Data: "System.exit(stack.ERROR.intValue());",
		},
	},

	"FUNCTION": is("static void %s(Stack stack) {", 1, 1),
	"RETURN": Instruction{
		Indented:    2,
		Data:        "}\n",
		Indent:      -1,
		Indentation: -1,
		Else: &Instruction{
			Data: "return;",
		},
	},
	
	"SCOPE": is(`Class[] cArg = new Class[1]; cArg[0] = Stack.class; try { stack.relay(new Stack.Pipe((new Object() { }.getClass().getEnclosingClass().getDeclaredMethod("%s", cArg)))); } catch (NoSuchMethodException e) { throw new RuntimeException(e); }`, 1),
	
	"EXE": is("%s.exe(stack);", 1),

	"PUSH": is("stack.push(%s);", 1),
	"PULL": is("Stack.Number %s = stack.pull();", 1),

	"PUT":   is("stack.put(%s);", 1),
	"POP":   is("Stack.Number %s = stack.pop();", 1),
	"PLACE": is("stack.place(%s);", 1),

	"ARRAY":  is("Stack.Array %s = stack.array();", 1),
	"MAKE":  is("stack.share(new Stack.Array(stack.pull().intValue()));"),
	"RENAME": is("%s = stack.ActiveArray;", 1),
	
	"RELOAD": is("%s = stack.take();", 1),

	"SHARE": is("stack.share(%s);", 1),
	"GRAB":  is("Stack.Array %s = stack.grab();", 1),

	"RELAY": is("stack.relay(%s);", 1),
	"TAKE":  is("Stack.Pipe %s = stack.take();", 1),

	"GET": is("Stack.Number %s = stack.get();", 1),
	"SET": is("stack.set(%s);", 1),

	"VAR": is("Stack.Number %s = new Stack.Number();", 1),

	"OPEN":   is("stack.open();"),
	"DELETE":   is("stack.delete();"),
	"EXECUTE":   is("stack.execute();"),
	"LOAD":   is("stack.load();"),
	"OUT":    is("stack.out();"),
	"STAT":   is("stack.info();"),
	"IN":     is("stack.in();"),
	"STDOUT": is("stack.stdout();"),
	"STDIN":  is("stack.stdin();"),
	"HEAP":   is("stack.heap();"),
	"LINK":   is("stack.link();"),
	"CONNECT":   is("stack.connect();"),
	"SLICE":   is("stack.slice();"),

	"CLOSE": is("%s.close();", 1),

	"LOOP":   is("while (true) {", 0, 1),
	"BREAK":  is("break;"),
	"REPEAT": is("}", 0, -1, -1),

	"IF":   is("if (%s.compareTo(new Stack.Number(0)) != 0 ) {", 1, 1),
	"ELSE": is("} else {", 0, 0, -1),
	"END":  is("}", 0, -1, -1),

	"RUN":  is("%s(stack);", 1),
	"DATA": is("static Stack.Array %s = %s;", 2),

	"FORK": is("{ Stack s = stack.copy(); Stack.ThreadPool.execute(() -> %s(s)); }\n", 1),

	"ADD": is("%s = %s.add(%s);", 3),
	"SUB": is("%s = %s.sub(%s);", 3),
	"MUL": is("%s = %s.mul(%s);", 3),
	"DIV": is("%s = %s.div(%s);", 3),
	"MOD": is("%s = %s.mod(%s);", 3),
	"POW": is("%s = %s.pow(%s);", 3),

	"SLT": is("%s = %s.compareTo(%s) == -1 ? new Stack.Number(1) : new Stack.Number(0);", 3),
	"SEQ": is("%s = %s.compareTo(%s) == 0 ? new Stack.Number(1) : new Stack.Number(0);", 3),
	"SGE": is("%s = %s.compareTo(%s) >= 0 ? new Stack.Number(1) : new Stack.Number(0);", 3),
	"SGT": is("%s = %s.compareTo(%s) == 1 ? new Stack.Number(1) : new Stack.Number(0);", 3),
	"SNE": is("%s = %s.compareTo(%s) != 0 ? new Stack.Number(1) : new Stack.Number(0);", 3),
	"SLE": is("%s = %s.compareTo(%s) <= 0 ? new Stack.Number(1) : new Stack.Number(0);", 3),

	"JOIN": is("%s = %s.join(%s);", 3),
	"ERROR": is("stack.ERROR = %s;", 1),
}
