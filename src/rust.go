package main

import "flag"

var RustReserved = []string{
	/*"break",        "default",      "func",         "interface",    "select",
	"case",         "defer",        "go",           "map",          "struct",
	"chan",         "else",         "goto",         "package",      "switch",
	"const",        "fallthrough",  "if",           "range",        "type",
	"continue",     "for",          "import",       "return",       "var",
	"bool",			"byte", 		"len", 			"open", 		"file", 
	"close", 		"load", 		"copy",			"new",*/
}

//This is the Java compiler for uct.
var Rust bool

func init() {
	flag.BoolVar(&Rust, "rs", false, "Target Rust")

	RegisterAssembler(RustAssembly, &Rust, "rs", "//")

	for _, word := range RustReserved {
		RustAssembly[word] = Reserved()
	}
}

var RustAssembly = Assemblable{
	//Special commands.
	"HEADER": Instruction{
		Data: `
#![allow(warnings)]
extern crate num;
extern crate rand;
mod stack;

use num::bigint::BigInt;
use num::bigint::{ToBigInt, Sign};
use num::ToPrimitive;

use stack::Stack;
use stack::Div;
use stack::Mod;
use stack::Pipe;
use stack::Pointer;
`,
		Args: 1,
	},

	"FOOTER": Instruction{},

	"FILE": Instruction{
		Path: "stack.rs",
	},

	"RUST": Instruction{All:true},

	"NUMBER": is("%s.to_bigint().unwrap()", 1),
	"BIG": 	is("BigInt(`%s`)", 1),
	"SIZE":   is("stack.Borrow[%s.address].len().to_bigint().unwrap()", 1),
	"STRING": is("stack.NewStringArray(%s)", 1),
	"ERRORS":  is("stack.ERROR", 1),
	
	//Experimental.
	//"EVAL": is(`stack.Eval[stack.Grab().String()](stack)`),
	//"EVALUATION": is(`stack.Eval["%s"] = %s`, 2),

	"SOFTWARE": Instruction{
		Data:   "fn main() { let mut stack = Stack::new();",
		Indent: 1,
	},
	"EXIT": Instruction{
		Indented:    1,
		Data:        "}\n",
		Indent:      -1,
		Indentation: -1,
		Else: &Instruction{
			Data: "exit",
		},
	},

	"FUNCTION": is("fn %s(mut stack: &mut Stack) {", 1, 1),
	"RETURN": Instruction{
		Indented:    1,
		Data:        "}\n",
		Indent:      -1,
		Indentation: -1,
		Else: &Instruction{
			Data: "return",
		},
	},
	
	"SCOPE": is(`stack.Relay(Pipe{Function:%s});`, 1),
	
	"EXE": is("%s.Exe(&mut stack);", 1),

	"PUSH": is("stack.Push(%s.clone());", 1),
	"PULL": is("let mut %s = stack.Pull();", 1),

	"PUT":   is("stack.Put(%s);", 1),
	"POP":   is("let mut %s = stack.Pop();", 1),
	"PLACE": is("stack.ActiveArray.address = %s.address;", 1),
	"ARRAY":  is("let mut %s = stack.Array();", 1),
	"MAKE":  is("stack.Make();"),
	"RENAME": is("%s = stack.ActiveArray.clone(); stack.Library[stack.ActiveArray.address].set(stack.Library[stack.ActiveArray.address].get()+1);", 1),
	"RELOAD": is("%s = stack.Take();", 1),

	"SHARE": is("stack.Share(%s.clone());", 1),
	"GRAB":  is("let mut %s = stack.Grab();", 1),

	"PIPE": is("let mut %s = stack.Pipe();", 1),
	"RELAY": is("stack.Relay(%s);", 1),
	"TAKE":  is("let mut %s = stack.Take();", 1),

	"GET": is("let mut %s = stack.Get();", 1),
	"SET": is("stack.Set(%s);", 1),

	"VAR": is("let mut %s = 0.to_bigint().unwrap();", 1),

	//"INBOX":   is("stack.Share(<- stack.Inbox)"),
	//"OUTBOX":   is("stack.Outbox <- stack.Grab()"),

	"OPEN":   is("stack.Open();"),
	"DELETE": is("stack.Delete();"),
	"MOVE":   is("stack.Move();"),
	"EXECUTE":is("stack.Execute();"),
	"LINK":   is("stack.Link();"),
	"CONNECT":is("stack.Connect();"),
	"SLICE":  is("stack.Slice();"),
	"LOAD":   is("stack.Load();"),
	"OUT":    is("stack.Out();"),
	"STAT":   is("stack.Info();"),
	"IN":     is("stack.In();"),
	"STDOUT": is("stack.Stdout();"),
	"STDIN":  is("stack.Stdin();"),
	"HEAP":  is("stack.Heap();"),
	"HEAPIT":  is("stack.HeapIt();"),

	"CLOSE": is("%s.Close();", 1),

	"LOOP":   is("loop {", 0, 1),
	"BREAK":  is("break;"),
	"REPEAT": is("}", 0, -1, -1),

	"IF":   is("if %s != 0.to_bigint().unwrap() {", 1, 1),
	"ELSE": is("} else {", 0, 0, -1),
	"END":  is("}", 0, -1, -1),

	"RUN":  is("%s(&mut stack);", 1),
	
	"DATA": is(" ", 2),

	"FORK": is("go %s(stack.Copy())\n", 1),

	"ADD": is("%s = &%s + &%s;", 3),
	"SUB": is("%s = &%s - &%s;", 3),
	"MUL": is("%s = &%s * &%s;", 3),
	"DIV": is("%s = Div(&%s, &%s);", 3),
	"MOD": is("%s = Mod(&%s, &%s);", 3),
	"POW": is("%s = num::pow(%s, %s.to_usize().unwrap());", 3),

	"SLT": is("%s = if %s < %s { 1.to_bigint().unwrap() } else { 0.to_bigint().unwrap() };", 3),
	"SEQ": is("%s = if %s == %s { 1.to_bigint().unwrap() } else { 0.to_bigint().unwrap() };", 3),
	"SGE": is("%s = if %s >= %s { 1.to_bigint().unwrap() } else { 0.to_bigint().unwrap() };", 3),
	"SGT": is("%s = if %s > %s { 1.to_bigint().unwrap() } else { 0.to_bigint().unwrap() };", 3),
	"SNE": is("%s = if %s != %s { 1.to_bigint().unwrap() } else { 0.to_bigint().unwrap() };", 3),
	"SLE": is("%s = if %s <= %s { 1.to_bigint().unwrap() } else { 0.to_bigint().unwrap() };", 3),

	"JOIN": is("%s = stack.Join(&%s, &%s);", 3),
	"ERROR": is("stack.ERROR = %s;", 1),
}
