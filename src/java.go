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
import java.io.*;
import java.lang.reflect.*; 
public class `+g.FileName+` {
	
	static NString N = new NString();
	static NStringString N2 = new NStringString();
	static Funcs F = new Funcs();
	static ITs F2 = new ITs();
	
	static BigInteger ERROR = BigInteger.valueOf(0);
	
	static class Funcs {
		ArrayList<Method> List;
		
		public void push(Method n) {
			List.add(n);
		}
		
		public Funcs() {
			List = new ArrayList<Method>();
		}
		
		public Method pop() {
			return List.remove(List.size()-1);
		}
	}
	
	static class IT {
		String Name;
		FileReader FileRead;
		FileWriter FileWrite;
	}
	
	static class ITs {
		ArrayList<IT> List;
		
		public void push(IT n) {
			List.add(n);
		}
		
		public ITs() {
			List = new ArrayList<IT>();
		}
		
		public IT pop() {
			return List.remove(List.size()-1);
		}
	}

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
		
		void set(BigInteger index, BigInteger n) {
			List.set(index.intValue(), n);
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
	
	static void pushit(IT n) {
		F2.push(n);
	}
	
	static NString popstring() {
		return N2.pop();
	}
	
	static IT popit() {
		return F2.pop();
	}
	
	
	static void pushfunc(Method n) {
		F.push(n);
	}
	
	static Method popfunc() {
		return F.pop();
	}

	static void push(BigInteger n) {
		N.push(n);
	}
	
	static BigInteger pop() {
		return N.pop();
	}
	
	static IT open() {
		String filename = "";
		NString text = popstring();
		for (int i = 0; i < text.size().intValue(); i++) {
			if (text.index(BigInteger.valueOf(i)) != null) {
				int c = text.index(BigInteger.valueOf(i)).intValue();
				filename += (char)(c);
			}
		} 
		
		File file = new File(filename);
		IT it = new IT();
		it.Name = filename;
		
		if (file.exists()) {
			push(BigInteger.valueOf(0));
			try {
				it.FileRead = new FileReader(file);
				it.FileWrite = new FileWriter(file, true);
			}catch(FileNotFoundException e){
				
			}catch(IOException e){
			
			}
			return it;
		}
		push(BigInteger.valueOf(-1));
		return it;
	}
	
	static void in(IT file) {
		BigInteger length = pop();
		for (int i = 0; i < length.intValue(); i++) {
			try {
				push(BigInteger.valueOf(file.FileRead.read()));
			}catch(IOException e){
				push(BigInteger.valueOf(-1000));
			}
		}
	}
	
	static void out(IT file) {
		NString text = popstring();
		for (int i = 0; i < text.size().intValue(); i++) {
			if (text.index(BigInteger.valueOf(i)) != null) {
				int c = text.index(BigInteger.valueOf(i)).intValue();
				try {
					file.FileWrite.write((char)(c));
				}catch(IOException e){
					push(BigInteger.valueOf(-1));
				}
			}
		}
	}
	
	static void close(IT file) {
		try {
			file.FileRead.close();
			file.FileWrite.close();
		}catch(IOException e){
		}catch(NullPointerException e){
		}
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
			case "char", "byte", "open":
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
		case "FUNC":
			return []byte(g.indt()+"Method "+args[0]+" = null; try { "+args[0]+" = " +
				g.FileName+".class.getDeclaredMethod(\""+args[1]+
				"\", (Class<?>[])null); } catch (NoSuchMethodException e) { throw new RuntimeException(e); }\n"), nil
		case "EXE":
			return []byte(g.indt()+"try { "+args[0]+
				".invoke(null); } catch (Exception e) {  throw new RuntimeException(e);}\n"), nil
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
				return []byte(g.indt()+"push"+name+"("+args[0]+");\n"), nil
			} else {
				return []byte(g.indt()+args[1]+".push("+args[0]+");\n"), nil
			}
		case "POP", "POPSTRING", "POPFUNC", "POPIT":
			var name string
			var typ string = "BigInteger"
			if command == "POPSTRING" {
				name = "string"
				typ  = "NString"
			}
			if command == "POPFUNC" {
				name = "func"
				typ  = "Method"
			}
			if command == "POPIT" {
				name = "it"
				typ  = "IT"
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
				return []byte(g.indt()+"BigInteger "+args[0]+" = BigInteger.valueOf(0); \n"), nil
			} else {
				return []byte(g.indt()+"BigInteger "+args[0]+" = "+args[1]+"; \n"), nil
			}
		
		//IT stuff.
		case "OPEN":
			return []byte(g.indt()+"IT "+args[0]+" = open();\n"), nil
		case "OUT":
			return []byte(g.indt()+"out("+args[0]+");\n"), nil
		case "IN":
			return []byte(g.indt()+"in("+args[0]+");\n"), nil
		case "CLOSE":
			return []byte(g.indt()+"close("+args[0]+");\n"), nil
			
		case "ERROR":
			return []byte(g.indt()+"ERROR = "+args[0]+";\n"), nil	
	
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
		case "MOD":
			return []byte(g.indt()+args[0]+" = "+args[1]+".mod("+args[2]+");\n"), nil
		case "POW":
			return []byte(g.indt()+args[0]+" = "+args[1]+".pow("+args[2]+".intValue());\n"), nil
			
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
