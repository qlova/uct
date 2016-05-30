package main

import "errors"
import "flag"
import "strings"
import "strconv"

//CSharp flag.
var CSharp bool
func init() {
	flag.BoolVar(&CSharp, "cs", false, "Target C#")
	
	RegisterAssembler(new(CSharpAssembler), &CSharp, "cs", "//")
}

type CSharpAssembler struct {
	Indentation int
	FileName string
}

func (g *CSharpAssembler) SetFileName(s string) {
	g.FileName = s
}

func (g *CSharpAssembler) Header() []byte {
	defer func() { g.Indentation++ }()
	return []byte(
	`
using System;
using System.Collections;
using System.Collections.Generic;
using System.Numerics;
using System.Reflection;

public class `+g.FileName+` {
	
	static NString N = new NString();
	static NStringString N2 = new NStringString();
	static Funcs F = new Funcs();

	public class Funcs {
		List<MethodInfo> Value;
		
		public void push(MethodInfo n) {
			Value.Add(n);
		}
		
		public Funcs() {
			Value = new List<MethodInfo>();
		}
		
		public MethodInfo pop() {
			MethodInfo temp;
			temp = Value[Value.Count-1];
			Value.RemoveAt(Value.Count-1);
			return temp;
		}
	}

	public class NStringString {
		List<NString> Value;
		
		public void push(NString n) {
			Value.Add(n);
		}
		
		public NStringString() {
			Value = new List<NString>();
		}
		
		public NString pop() {
			NString temp;
			temp = Value[Value.Count-1];
			Value.RemoveAt(Value.Count-1);
			return temp;
		}
		
		public BigInteger size() {
			return new BigInteger(Value.Count);
		}
		
		public NString index(BigInteger n) {
			return Value[(int)n];
		}
	}
	
	public class NString {
		List<BigInteger> Value;
		//new ArrayList<Integer>();
		
		public void push(BigInteger n) {
			Value.Add(n);
		}
		
		public void set(BigInteger index, BigInteger n) {
			Value.Insert((int)index, n);
		}
		
		public NString join(NString s) {
			NString newList = new NString(); 
			newList.Value.AddRange(Value);
			newList.Value.AddRange(s.Value);
			return newList;
		}
		
		public NString(params BigInteger[] n) {
			Value = new List<BigInteger>();
			for (int i = 0; i < n.Length; ++i) {
				Value.Add(n[i]);
			}
		}
		
		public BigInteger pop() {
			BigInteger temp;
			temp = Value[Value.Count-1];
			Value.RemoveAt(Value.Count-1);
			return temp;
		}
		
		public BigInteger size() {
			return new BigInteger(Value.Count);
		}
		
		public BigInteger index(BigInteger n) {
			return Value[(int)n];
		}
	}
	
	static void pushstring(NString n) {
		N2.push(n);
	}
	
	static NString popstring() {
		return N2.pop();
	}
	
	static void pushfunc(MethodInfo n) {
		F.push(n);
	}
	
	static MethodInfo popfunc() {
		return F.pop();
	}


	static void push(BigInteger n) {
		N.push(n);
	}
	
	static BigInteger pop() {
		return N.pop();
	}

	static void stdout() {
		NString text = popstring();
		for (int i = 0; i < (int)text.size(); i++) {
			int c = (int)text.index(new BigInteger(i));
			System.Console.Write((char)(c));
		} 
	}
	
	static void stdin() {
		BigInteger length = pop();
		for (int i = 0; i < (int)length; i++) {
			/*try {*/
				push(new BigInteger(Console.Read()));
			/*}catch(IOException e){
				push(new BigInteger(-1));
			}*/
		}
	}
	
	static BigInteger slt(BigInteger a, BigInteger b) {
		if (a.CompareTo(b) == -1) {
			return new BigInteger(1);
		}
		return new BigInteger(0);
	}
	
	static BigInteger seq(BigInteger a, BigInteger b) {
		if (a.CompareTo(b) == 0) {
			return new BigInteger(1);
		}
		return new BigInteger(0);
	}
	
	static BigInteger sge(BigInteger a, BigInteger b) {
		if (a.CompareTo(b) >= 0) {
			return new BigInteger(1);
		}
		return new BigInteger(0);
	}
	
	static BigInteger sgt(BigInteger a, BigInteger b) {
		if (a.CompareTo(b) == 1) {
			return new BigInteger(1);
		}
		return new BigInteger(0);
	}
	
	static BigInteger sne(BigInteger a, BigInteger b) {
		if (a.CompareTo(b) != 0) {
			return new BigInteger(1);
		}
		return new BigInteger(0);
	}
	
	static BigInteger sle(BigInteger a, BigInteger b) {
		if (a.CompareTo(b) <= 0) {
			return new BigInteger(1);
		}
		return new BigInteger(0);
	}
`)
}

func (g *CSharpAssembler) indt(n ...int) string {
	if len(n) > 0 {
		return strings.Repeat("\t", int(g.Indentation)+n[0])
	} else {
		return strings.Repeat("\t", int(g.Indentation))
	}
}

func (g *CSharpAssembler) Assemble(command string, args []string) ([]byte, error) {
	//Format argument expressions.
	//var expression bool
	for i, arg := range args {
		if arg == "=" {
			//expression = true
			continue
		}
		if _, err := strconv.Atoi(arg); err == nil {
			args[i] = "new BigInteger("+arg+")"
			continue
		}
		if arg[0] == '#' {
			args[i] = ""+arg[1:]+".size()"
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
				newarg += "new BigInteger("+strconv.Itoa(int(v))+"),"
			}
			if len(arg) == 0 {
				goto end
			}
			newarg += "new BigInteger("+strconv.Itoa(int(' '))+"),"
			j++
			//println(arg)
			arg = args[j]
			goto stringloop
			end:
			//println(newarg)
			if len(newarg) > 0 {
				newarg = newarg[:len(newarg)-1]
			}
			args[i] = newarg
		}
		switch arg {
			case "char", "byte":
				args[i] = "u_"+args[i]
		}
	} 
	switch command {
		case "ROUTINE":
			defer func() { g.Indentation++ }()
			return []byte(g.indt()+"public static void Main(String[] args) {\n"), nil
		case "SUBROUTINE":
			defer func() { g.Indentation++ }()
			return []byte(g.indt()+"static void "+args[0]+"() {\n"), nil
		case "FUNC":
			return []byte(g.indt()+"MethodInfo "+args[0]+
				" = typeof(Functions).GetTypeInfo().GetDeclaredMethod(\""+args[1]+"\"); "), nil
		case "EXE":
			return []byte(g.indt()+args[0]+".Invoke(null, (Object[])null);"), nil
		case "PUSH", "PUSHSTRING", "PUSHFUNC":
			var name string
			if command == "PUSHSTRING" {
				name = "string"
			}
			if command == "PUSHFUNC" {
				name = "func"
			}
			if len(args) == 1 {
				return []byte(g.indt()+"push"+name+"("+args[0]+");\n"), nil
			} else {
				return []byte(g.indt()+args[1]+".push("+args[0]+");\n"), nil
			}
		case "POP", "POPSTRING", "POPFUNC":
			var name string
			var typ string = "BigInteger"
			if command == "POPSTRING" {
				name = "string"
				typ  = "NString"
			}
			if command == "POPFUNC" {
				name = "func"
				typ  = "MethodInfo"
			}
			if len(args) == 1 {
				return []byte(g.indt()+typ+" "+args[0]+" = pop"+name+"();\n"), nil
			} else {
				return []byte(g.indt()+typ+" "+args[0]+" = "+args[1]+".pop();\n"), nil
			}
		case "INDEX":
			return []byte(g.indt()+"BigInteger "+args[2]+" = "+args[0]+".index("+args[1]+");\n"), nil
		case "SET":
			return []byte(g.indt()+args[0]+".set("+args[1]+", "+args[2]+");\n"), nil
		case "VAR":
			if len(args) == 1 {
				return []byte(g.indt()+"BigInteger "+args[0]+" = 0; \n"), nil
			} else {
				return []byte(g.indt()+"BigInteger "+args[0]+" = "+args[1]+"; \n"), nil
			}
		case "STRING":
			return []byte(g.indt()+"NString "+args[0]+" = new NString();\n"), nil
		case "STDOUT":
			return []byte(g.indt()+"stdout();\n"), nil
		case "STDIN":
			return []byte(g.indt()+"stdin();\n"), nil
		case "LOOP":
			defer func() { g.Indentation++ }()
			return []byte(g.indt()+"while (true) {\n"), nil
		case "REPEAT", "END", "DONE":
			g.Indentation--
			return []byte(g.indt()+"}\n"), nil
		case "IF":
			defer func() { g.Indentation++ }()
			return []byte(g.indt()+"if ("+args[0]+".CompareTo(0) != 0) {\n"), nil
		case "RUN":
			return []byte(g.indt()+args[0]+"();\n"), nil
		case "ELSE":
			return []byte(g.indt(-1)+"} else {\n"), nil
		case "ELSEIF":
			return []byte(g.indt(-1)+"} else if ("+args[0]+".CompareTo(0) != 0) {\n"), nil
		case "BREAK":
			return []byte(g.indt()+"break;\n"), nil
		case "RETURN":
			return []byte(g.indt()+"return;\n"), nil
		case "STRINGDATA":
			return []byte(g.indt()+"static NString "+args[0]+" = new NString("+args[1]+"); \n"), nil
		case "JOIN":
			return []byte(g.indt()+args[0]+" = "+args[1]+".join("+args[2]+"); \n"), nil
		
		//Maths.
		case "ADD":
			return []byte(g.indt()+args[0]+" = "+args[1]+" + "+args[2]+";\n"), nil
		case "SUB":
			return []byte(g.indt()+args[0]+" = "+args[1]+" - "+args[2]+";\n"), nil
		case "MUL":
			return []byte(g.indt()+args[0]+" = "+args[1]+" * "+args[2]+";\n"), nil
		case "DIV":
			return []byte(g.indt()+args[0]+" = "+args[1]+" / "+args[2]+";\n"), nil
		case "MOD":
			return []byte(g.indt()+args[0]+" = (("+args[1]+"%"+args[2]+") + "+args[2]+") % "+args[2]+" ;\n"), nil
		case "POW":
			return []byte(g.indt()+args[0]+" = BigInteger.Pow("+args[1]+","+args[2]+");\n"), nil
		
		case "SLT", "SEQ", "SGE", "SGT", "SNE", "SLE":
			return []byte(g.indt()+args[0]+" = "+strings.ToLower(command)+"("+args[1]+","+args[2]+");\n"), nil
		default:
			/*if expression {
				return []byte(g.indt()+command+" "+strings.Join(args, " ")+"\n"), nil
			}*/
	}
	return nil, errors.New("Unrecognised expression: "+command+" "+strings.Join(args, " "))
}

func (g *CSharpAssembler) Footer() []byte {
	return []byte("}")
}
