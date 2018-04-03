package uct

import "flag"

var JavascriptReserved = []string{
	"and",       "del",       "from",      "not",       "while",   
	"as",        "elif",      "global",    "or",        "with",    
	"assert",    "else",      "if",        "pass",      "yield",    
	"break",     "except",    "import",    "print",     "len",          
	"class",     "exec",      "in",        "raise", 	"open",             
	"continue",  "finally",   "is",        "return",    "bool",   
	"def",       "for",       "lambda",    "try",		"copy",

}

//This is the Java compiler for uct.
var Javascript bool

func init() {
	flag.BoolVar(&Javascript, "js", false, "Target Javascript")
	
	JavascriptAssembly["INBOX"] = Instruction{
		Args: 0,
		Data: " ",
		Function: func(args []string) (result string) {
			
			result = "_unwind = setInterval(function() { if (stack.inbox <= 0) { return; } clearInterval(_unwind); stack.share(stack.inbox.shift())\n"
			
			JavascriptAssembly["js_Unwind"] = Instruction {
				Args: JavascriptAssembly["js_Unwind"].Args+1,
			}
			
			return result
		},
	}
	
	JavascriptAssembly["RETURN"] = Instruction{
		Indented:    1,
		Data:        "}\n",
		Indent:      -1,
		Else: &Instruction{
			Data: "return",
		},
		Function: func(args []string) (result string) {
			
			var times int = JavascriptAssembly["js_Unwind"].Args
			
			JavascriptAssembly["js_Unwind"] = Instruction{}
			
			
			
			for i:=0; i < times; i++ {
				result += `}, 100)`
			}
			
			return result+"}\n"
		},
	}
	
	JavascriptAssembly["EXIT"] =  Instruction{
		Indented:    1,
		Data:        "}\n",
		Indent:      -1,
		Indentation: -1,
		Else: &Instruction{
			Data: "return;",
		},
		Function: func(args []string) (result string) {
			
			var times int = JavascriptAssembly["js_Unwind"].Args
			
			JavascriptAssembly["js_Unwind"] = Instruction{}
			
			println(times)
			
			for i:=0; i < times; i++ {
				result += `}, 100)`
			}
			
			return result+"}\n"
		},
	}
	
	RegisterAssembler(JavascriptAssembly, &Javascript, "js", "//")

	for _, word := range JavascriptReserved {
		JavascriptAssembly[word] = Reserved()
	}
}

var JavascriptAssembly = Assemblable{
	//Special commands.
	"HEADER": Instruction{
		Data:   `
var bigInt = bigInt
if (typeof require != 'undefined') { require('./stack.js'); bigInt = global.bigInt } 
//Helper functions
String.prototype.getBytes = function () {
var bytes = [];
for (var i = 0; i < this.length; ++i) {
bytes.push(this.charCodeAt(i));
}
return bytes;
};
`,
		Args: 1,
	},

	"FOOTER": Instruction{ Data: "if (typeof WorkerGlobalScope !== 'undefined' && self instanceof WorkerGlobalScope) {} else { main(); }"},

	"FILE": Instruction{
		Path: "stack.js",
	},
	
	"JAVASCRIPT": Instruction{All:true},

	"NUMBER": is("bigInt(%s)", 1),
	"BIG": 	is("bigInt(\"%s\")", 1),
	"SIZE":   is("bigInt(%s.length)", 1),
	"STRING": is("%s.getBytes()", 1),
	"ERRORS":  is("stack.ERROR", 1),
	
	"LINK":  is("stack.link();"),
	"CONNECT":  is("stack.connect();"),
	"SLICE":  is("stack.slice();"),

	"SOFTWARE": Instruction{
		Data:   "function main() {\n\tstack = new Stack();\nbigInt = stack.bigInt;",
		Indent: 1,
	},

	"FUNCTION": is("function %s(stack) {", 1, 1),
	
	"SCOPE": is(`stack.relay(stack.pipe(%s));`, 1),
	
	"EXE": is("%s.exe(stack);", 1),

	//Optimised
	"PUSH": is("stack.numbers.push(%s);", 1),
	"PULL": is("var %s = stack.numbers.pop();", 1),
	"SHARE": is("stack.arrays.push(%s);", 1),
	"GRAB":  is("var %s = stack.arrays.pop();", 1),
	"PUT":   is("stack.activearray.push(%s);", 1),
	"POP":   is("var %s = stack.activearray.pop();", 1),
	"PLACE": is("stack.activearray = %s;", 1),
	"ARRAY":  is("stack.activearray = []; var %s = stack.activearray;", 1),
	"RENAME": is("%s = stack.grab();", 1),
	"EVAL": is("eval(stack.grabstring()+'(stack)');"),
	
	"RELOAD": is("%s = stack.take();", 1),

	"RELAY": is("stack.relay(%s);", 1),
	"TAKE":  is("var %s = stack.take();", 1),

	"GET": is("var %s = stack.get();", 1),
	"SET": is("stack.set(%s);", 1),

	"VAR": is("var %s = bigInt();", 1),

	"OPEN":   is("stack.open();"),
	"EXECUTE": is("stack.stdout();"),
	"DELETE": is("stack.delete();"),
	"LOAD":   is("stack.load();"),
	"OUT":    is("stack.out();"),
	"STAT":   is("stack.info();"),
	"IN":     is("stack.in();"),
	"STDOUT": is("stack.stdout();"),
	"STDIN":  is("stack.stdin();"),
	"HEAP":   is("stack.heap();"),
	"HEAPIT":   is("stack.heapit();"),
	"MAKE":   is("stack.share(new Array(stack.pull().toJSNumber()).fill(bigInt.zero));"),

	"CLOSE": is("%s.close();", 1),

	"LOOP":   is("while (1) {", 0, 1),
	"BREAK":  is("break;"),
	"REPEAT": is("}", 0, -1, -1),

	"IF":   is("if (%s != 0) {", 1, 1),
	"ELSE": is("} else {", 0, 0, -1),
	"END":  is("}", 0, -1, -1),

	"RUN":  is("%s(stack);", 1),
	"DATA": is_data("var %s = [", "%s", ",",  "];"),

	//Threading.
	"PIPE": is("%s = stack.pipe(stack.channel); stack.channel = stack.channel + 1;", 1),
	"FORK": is("stack.thread('%s');", 1),
	
	
	
	//"INBOX": is("while (stack.inbox.length <= 0) {} stack.share(stack.inbox.shift())", 0),
	"OUTBOX": is("stack.outbox()", 0),

	"ADD": is("%s = %s.add(%s)", 3),
	"SUB": is("%s = %s.subtract(%s)", 3),
	"MUL": is("%s = %s.multiply(%s)", 3),
	"DIV": is("%s = %s.divide(%s)", 3),
	"MOD": is("%s = stack.mod(%s, %s)", 3),
	"POW": is("%s = %s.pow(%s)", 3),

	"SLT": is("%s = %s.lt(%s) ? bigInt.one : bigInt.zero;", 3),
	"SEQ": is("%s = %s.equals(%s) ? bigInt.one : bigInt.zero;", 3),
	"SGE": is("%s = %s.geq(%s) ? bigInt.one : bigInt.zero;", 3),
	"SGT": is("%s = %s.gt(%s) ? bigInt.one : bigInt.zero;", 3),
	"SNE": is("%s = %s.neq(%s) ? bigInt.one : bigInt.zero;", 3),
	"SLE": is("%s = %s.leq(%s)? bigInt.one : bigInt.zero;", 3),

	"JOIN": is("%s = %s.concat(%s)", 3),
	"ERROR": is("stack.ERROR = %s", 1),
}
