package main

import "errors"
import "flag"
import "strings"
import "strconv"

//Java flag.
var Java bool
func init() {
	flag.BoolVar(&Java, "java", false, "Target Java")
	
	RegisterAssembler(new(JavaAssembler), &Java, "java", "//")
}

type JavaAssembler struct {
	Indentation int
	FileName string
}

func (g *JavaAssembler) SetFileName(s string) {
	g.FileName = s
}

func (g *JavaAssembler) Header() []byte {
	defer func() { g.Indentation++ }()
	return []byte(
	`
import java.math.BigInteger;
import java.util.ArrayList; 
import java.io.IOException; 
public class `+g.FileName+` {
	
	static NString N = new NString();
	static NStringString N2 = new NStringString();

	static class NStringString {
		ArrayList<NString> List;
		
		void push(NString n) {
			List.add(n);
		}
		
		public NStringString() {
			List = new ArrayList<NString>();
		}
		
		NString pop() {
			return List.remove(List.size()-1);
		}
		
		BigInteger size() {
			return BigInteger.valueOf(List.size());
		}
		
		NString index(BigInteger n) {
			return List.get(n.intValue());
		}
	}
	
	static class NString {
		ArrayList<BigInteger> List;
		//new ArrayList<Integer>();
		
		void push(BigInteger n) {
			List.add(n);
		}
		
		NString join(NString s) {
			NString newList = new NString(); 
			newList.List.addAll(List);
			newList.List.addAll(s.List);
			return newList;
		}
		
		public NString(BigInteger... n) {
			List = new ArrayList<BigInteger>();
			for (int i = 0; i < n.length; ++i) {
				List.add(n[i]);
			}
		}
		
		BigInteger pop() {
			return List.remove(List.size()-1);
		}
		
		BigInteger size() {
			return BigInteger.valueOf(List.size());
		}
		
		BigInteger index(BigInteger n) {
			return List.get(n.intValue());
		}
	}
	
	static void pushstring(NString n) {
		N2.push(n);
	}
	
	static NString popstring() {
		return N2.pop();
	}

	static void push(BigInteger n) {
		N.push(n);
	}
	
	static BigInteger pop() {
		return N.pop();
	}

	static void stdout() {
		NString text = popstring();
		for (int i = 0; i < text.size().intValue(); i++) {
			if (text.index(BigInteger.valueOf(i)) != null) {
				int c = text.index(BigInteger.valueOf(i)).intValue();
				System.out.print((char)(c));
			}
		} 
	}
	
	static void stdin() {
		BigInteger length = pop();
		for (int i = 0; i < length.intValue(); i++) {
			try {
				push(BigInteger.valueOf(System.in.read()));
			}catch(IOException e){
				push(BigInteger.valueOf(-1));
			}
		}
	}
	
	static BigInteger slt(BigInteger a, BigInteger b) {
		if (a.compareTo(b) == -1) {
			return BigInteger.valueOf(1);
		}
		return BigInteger.valueOf(0);
	}
	
	static BigInteger seq(BigInteger a, BigInteger b) {
		if (a.compareTo(b) == 0) {
			return BigInteger.valueOf(1);
		}
		return BigInteger.valueOf(0);
	}
	
	static BigInteger sge(BigInteger a, BigInteger b) {
		if (a.compareTo(b) >= 0) {
			return BigInteger.valueOf(1);
		}
		return BigInteger.valueOf(0);
	}
	
	static BigInteger sgt(BigInteger a, BigInteger b) {
		if (a.compareTo(b) == 1) {
			return BigInteger.valueOf(1);
		}
		return BigInteger.valueOf(0);
	}
	
	static BigInteger sne(BigInteger a, BigInteger b) {
		if (a.compareTo(b) != 0) {
			return BigInteger.valueOf(1);
		}
		return BigInteger.valueOf(0);
	}
	
	static BigInteger sle(BigInteger a, BigInteger b) {
		if (a.compareTo(b) <= 0) {
			return BigInteger.valueOf(1);
		}
		return BigInteger.valueOf(0);
	}
`)
}

func (g *JavaAssembler) indt(n ...int) string {
	if len(n) > 0 {
		return strings.Repeat("\t", int(g.Indentation)+n[0])
	} else {
		return strings.Repeat("\t", int(g.Indentation))
	}
}

func (g *JavaAssembler) Assemble(command string, args []string) ([]byte, error) {
	//Format argument expressions.
	//var expression bool
	for i, arg := range args {
		if arg == "=" {
			//expression = true
			continue
		}
		if _, err := strconv.Atoi(arg); err == nil {
			args[i] = "BigInteger.valueOf("+arg+")"
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
				newarg += "BigInteger.valueOf("+strconv.Itoa(int(v))+"),"
			}
			if len(arg) == 0 {
				goto end
			}
			newarg += "BigInteger.valueOf("+strconv.Itoa(int(' '))+"),"
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
			return []byte(g.indt()+"public static void main(String[] args) {\n"), nil
		case "SUBROUTINE":
			defer func() { g.Indentation++ }()
			return []byte(g.indt()+"static void "+args[0]+"() {\n"), nil
		case "PUSH", "PUSHSTRING":
			var name string
			if command == "PUSHSTRING" {
				name = "string"
			}
			if len(args) == 1 {
				return []byte(g.indt()+"push"+name+"("+args[0]+");\n"), nil
			} else {
				return []byte(g.indt()+args[1]+".push("+args[0]+");\n"), nil
			}
		case "POP", "POPSTRING":
			var name string
			var typ string = "BigInteger"
			if command == "POPSTRING" {
				name = "string"
				typ  = "NString"
			}
			if len(args) == 1 {
				return []byte(g.indt()+typ+" "+args[0]+" = pop"+name+"();\n"), nil
			} else {
				return []byte(g.indt()+typ+" "+args[0]+" = "+args[1]+".pop();\n"), nil
			}
		case "INDEX":
			return []byte(g.indt()+"BigInteger "+args[2]+" = "+args[0]+".index("+args[1]+");\n"), nil
		case "VAR":
			if len(args) == 1 {
				return []byte(g.indt()+"BigInteger "+args[0]+" = BigInteger.valueOf(0); \n"), nil
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
			return []byte(g.indt()+"if ("+args[0]+".compareTo(BigInteger.valueOf(0)) != 0) {\n"), nil
		case "RUN":
			return []byte(g.indt()+args[0]+"();\n"), nil
		case "ELSE":
			return []byte(g.indt(-1)+"} else {\n"), nil
		case "ELSEIF":
			return []byte(g.indt(-1)+"} else if ("+args[0]+".compareTo(BigInteger.valueOf(0)) != 0) {\n"), nil
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
			return []byte(g.indt()+args[0]+" = "+args[1]+".add("+args[2]+");\n"), nil
		case "SUB":
			return []byte(g.indt()+args[0]+" = "+args[1]+".subtract("+args[2]+");\n"), nil
		case "MUL":
			return []byte(g.indt()+args[0]+" = "+args[1]+".multiply("+args[2]+");\n"), nil
		case "DIV":
			return []byte(g.indt()+args[0]+" = "+args[1]+".divide("+args[2]+");\n"), nil
			
		case "SLT", "SEQ", "SGE", "SGT", "SNE", "SLE":
			return []byte(g.indt()+args[0]+" = "+strings.ToLower(command)+"("+args[1]+","+args[2]+");\n"), nil
		default:
			/*if expression {
				return []byte(g.indt()+command+" "+strings.Join(args, " ")+"\n"), nil
			}*/
	}
	return nil, errors.New("Unrecognised expression: "+command+" "+strings.Join(args, " "))
}

func (g *JavaAssembler) Footer() []byte {
	return []byte("}")
}
