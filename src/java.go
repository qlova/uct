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
		Data: JavaFile,
		Path: "/Stack.java",
	},

	"NUMBER": is("new Stack.Number(%s)", 1),
	"SIZE":   is("%s.size()", 1),
	"STRING": is("new Stack.Array(%s)", 1),
	"ERROR":  is("stack.ERROR", 1),

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
			Data: "System.exit(stack.ERROR);",
		},
	},

	"FUNCTION": is("static void %s(Stack stack) {", 1, 1),
	"RETURN":   is("}", 0, -1, -1),

	"PUSH": is("stack.push(%s);", 1),
	"PULL": is("Stack.Number %s = stack.pull();", 1),

	"PUT":   is("stack.put(%s);", 1),
	"POP":   is("Number %s = stack.pop();", 1),
	"PLACE": is("stack.place(%s);", 1),

	"ARRAY":  is("Stack.Array %s = stack.array();", 1),
	"RENAME": is("%s = stack.ActiveArray", 1),

	"SHARE": is("stack.share(%s);", 1),
	"GRAB":  is("ARRAY %s = stack.grab();", 1),

	"RELAY": is("stack.relay(%s);", 1),
	"TAKE":  is("PIPE %s = stack.take();", 1),

	"GET": is("Stack.Number %s = stack.get();", 1),
	"SET": is("stack.set(%s);", 1),

	"VAR": is("Stack.Number %s = new Number();", 1),

	"OPEN":   is("stack.open();"),
	"LOAD":   is("stack.load();"),
	"OUT":    is("stack.out();"),
	"IN":     is("stack.in();"),
	"STDOUT": is("stack.stdout();"),
	"STDIN":  is("stack.stdin();"),

	"CLOSE": is("%s.close();", 1),

	"LOOP":   is("while (true) {"),
	"BREAK":  is("break"),
	"REPEAT": is("}"),

	"IF":   is("if (%s.compareTo(new Number(1)) == 0 ) {", 1, 1),
	"ELSE": is("} else {", 0, 0, -1),
	"END":  is("}", 0, -1, -1),

	"RUN":  is("%s();", 1),
	"DATA": is("static Stack.Array %s = %s;", 2),

	"FORK": is("{ Stack s = stack.copy(); pool__java.execute(() -> %s(s)); }\n"),

	"ADD": is("%s = %s.add(%s);"),
	"SUB": is("%s = %s.sub(%s);"),
	"MUL": is("%s = %s.mul(%s);"),
	"DIV": is("%s = %s.div(%s);"),
	"MOD": is("%s = %s.mod(%s);"),
	"POW": is("%s = %s.pow(%s);"),

	"SLT": is("%s = %s.compareTo(%s) == -1 ? Stack.One : Stack.Zero;"),
	"SEQ": is("%s = %s.compareTo(%s) == 0 ? Stack.One : Stack.Zero;"),
	"SGE": is("%s = %s.compareTo(%s) >= 0 ? Stack.One : Stack.Zero;"),
	"SGT": is("%s = %s.compareTo(%s) == 1 ? Stack.One : Stack.Zero;"),
	"SNE": is("%s = %s.compareTo(%s) != 0 ? Stack.One : Stack.Zero;"),
	"SLE": is("%s = %s.compareTo(%s) <= 0 ? Stack.One : Stack.Zero;"),

	"JOIN": is("%s = %s.join(%s);"),
}

//Edit this in a Java IDE.
const JavaFile = `
//Compiled to Java with UCT (Universal Code Translator)

//Import java libraries.
import java.math.BigInteger; 	//Support numbers of any size.
import java.util.Hashtable;  	//Hashtable is a useful utilPipey.
import java.util.ArrayList;  	//ArrayLists are helpful.
import java.io.*;				//Deal with files.
import java.net.*;				//Deal with the network.
import java.lang.reflect.*;		//Reflection for methods.

//Random numbers.
import java.security.SecureRandom;

//This is for threading.
import java.util.concurrent.Executor;
import java.util.concurrent.Executors;

//This is the Java stack implementation.
// It holds arrays for the 4 types:
//		Numbers
//		Arrays
//		Functions
//		Pipes
//
// It also holds the ERROR variable for the current thread.
// The currently active array is stored as ActiveArray.
public class Stack {
    Array 			Numbers;
    ArrayArray 		Arrays;
    PipeArray 		Pipes;

    Number 			ERROR;

    Array			ActiveArray;
    Array			Heap;

    //This hashtable keeps track of Servers currently listening on the specified port.
    static Hashtable<String, ServerSocket> Networks_In = new Hashtable<String, ServerSocket>();

    //This will store the system arguments.
    public static String[] Arguments;

  	static final Number Zero = (Number)Number.ZERO;
    static final Number One = (Number)Number.ONE;

    //This is the threading pool.
    public static Executor ThreadPool = Executors.newCachedThreadPool();

    //This creates an empty stack. Ready for use.
    public Stack() {
        Numbers 		= new Array();
        Arrays 			= new ArrayArray();
        Pipes 		    = new PipeArray();
        ERROR 			= new Number();

        Heap			= new Array();
    }

    //This returns a copy of a stack which can be used by another thread.
    public Stack copy() {
        Stack n = new Stack();
        n.Numbers = new Array();

        n.Numbers.List = new ArrayList<Number>(Numbers.List);

        n.Arrays.List = new ArrayList<Array>(Arrays.List.size());
        for (int i = 0; i < Arrays.List.size()-1; i++) {
            n.Arrays.List.set(i, new Array());
            n.Arrays.List.get(i).List = new ArrayList<>(Arrays.List.get(i).List);
        }

        n.Pipes = new PipeArray();
        n.Pipes.List = new ArrayList<>(Pipes.List);
        return n;
    }

    //This stuff can be inlined in the compiler, they are only here for reference.

    //SHARE array
    void share(Array a) {
        Arrays.push(a);
    }

    //GRAB array
    Array grab() {
        return Arrays.pop();
    }

    //RELAY pipe
    void relay(Pipe p) {
        Pipes.push(p);
    }

    //TAKE pipe
    Pipe take() {
        return Pipes.pop();
    }

    //PUSH number
    void push(Number n) {
        Numbers.push(n);
    }

    //PUT number
    void put(Number b) {
        ActiveArray.push(b);
    }

    //PLACE array
    void place(Array a) {
        ActiveArray = a;
    }
    
    //ARRAY name
    Array array() {
    	Array a = new Array();
    	ActiveArray = a;
    	return a;
    }

    //GET number
    Number get(Array a) {
        return (Number)ActiveArray.index(pull());
    }

    //SET number
    void set(Number b) {
        ActiveArray.set(pull(), b);
    }

    //POP number
    Number pop(Array a) {
        return (Number)ActiveArray.pop();
    }

    //PULL number
    Number pull() {
        return (Number)Numbers.pop();
    }

    void stdout() {
        Array text = grab();
        for (int i = 0; i < text.size().intValue(); i++) {
            if (text.index(BigInteger.valueOf(i)) != null) {
                int c = text.index(BigInteger.valueOf(i)).intValue();
                System.out.print((char)(c));
            }
        }
    }

    void stdin() {
        Number length = pull();
        for (int i = 0; i < length.intValue(); i++) {
            try {
                int c = System.in.read();
                if (c == -1) {
                    push((Number)BigInteger.valueOf(-1000));
                    return;
                }
                push((Number)BigInteger.valueOf(c));
            }catch(Exception e){
                push((Number)BigInteger.valueOf(-1000));
                return;
            }
        }
    }

    void in() {
        Pipe file = take();
        BigInteger length = pull();
        int n = 0;
        byte[] b = new byte[length.intValue()];

        if (file.input != null) {
            try {
                n = file.input.read(b);
            }catch(Exception e){
                push((Number)BigInteger.valueOf(-1000));
            }
        }

        if ((b.length > 1) || (n <= 0)) {
            push((Number)BigInteger.valueOf(-1000));
        }

        for (int i = n-1; i >= 0; i--) {
            push((Number)BigInteger.valueOf(b[i]));
        }

        return;
    }

    void out() {
        Pipe file = take();
        Array text = grab();

        if (file.output != null) {
            //TODO optimise to send in a single packet.
            for (int i = 0; i < text.size().intValue(); i++) {
                if (text.index(BigInteger.valueOf(i)) != null) {
                    int c = text.index(BigInteger.valueOf(i)).intValue();
                    try {
                        file.output.write((char)(c));
                    }catch(Exception e){
                        push((Number)BigInteger.valueOf(-1));
                    }
                }
            }
            push((Number)BigInteger.valueOf(0));
            return;
        }

        if (text.size().intValue() == 0 || file.output == null ) {
            if (file.Name.charAt(file.Name.length()-1) == '/') {
                if (new File(file.Name).exists()) {

                } else {
                    try {
                        File f = new File(file.Name);
                        if (!f.mkdir()) {
                            push((Number)BigInteger.valueOf(-1));
                            return;
                        }
                        push((Number)BigInteger.valueOf(0));
                        return;
                    } catch (Exception e) {
                        push((Number)BigInteger.valueOf(-1));
                        return;
                    }
                }
            } else if (new File(file.Name).exists()) {

            } else {
                try {
                    new File(file.Name).createNewFile();
                    push((Number)BigInteger.valueOf(0));
                    return;
                } catch (Exception e)  {
                    push((Number)BigInteger.valueOf(-1));
                    return;
                }
            }
        }

        for (int i = 0; i < text.size().intValue(); i++) {
            if (text.index(BigInteger.valueOf(i)) != null) {
                int c = text.index(BigInteger.valueOf(i)).intValue();
                try {
                    file.output.write((char)(c));
                }catch(Exception e){
                    push((Number)BigInteger.valueOf(-1));
                }
            }
        }
        push((Number)BigInteger.valueOf(0));
    }


    //LOAD
    static void load(Stack stack) {
        String name = "";
        String variable = "";

        Array result = new Array();

        //This request is what we need to load.
        Array request = stack.grab();

        //The request is an enviromental variable.
        if (request.index(0) == '$' && request.length() > 1) {

            //We parse the rest of the string.
            for (int i = 1; i < request.length(); i++) {
                try {
                    name += (char)(request.index(i));
                } catch (Exception e) {

                }
            }
            stack.share(new Array(System.getenv(name)));
            return;
        }


        name = request.String();

        //Load various protocols.
        //Protocols are seperated by ://
        //For example http://, dns://, tcp://
        String[] protocol = name.split("://", 2);
        if (protocol.length > 1) {
            switch (protocol[0]) {
                case "tcp":
                    try {
                        ServerSocket ss = new ServerSocket(Integer.parseInt(protocol[1]));
                        String port = String.valueOf(ss.getLocalPort());
                        if (protocol[1].equals("0")) {
                            Networks_In.put(port, ss);
                            variable = String.valueOf(port);
                        } else {
                            Networks_In.put(protocol[1], ss);
                            variable = protocol[1];
                        }
                    } catch (Exception e) {
                        stack.ERROR = new Number(-1);
                    }
                    break;

                case "dns":
                    try {
                        InetAddress address = InetAddress.getByName(protocol[1]);
                        variable = address.getHostName();
                        if (variable == protocol[1]) {
                            variable = "";
                            throw null;
                        }
                    } catch (Exception e) {
                        try {
                            InetAddress[] addresses = InetAddress.getAllByName(protocol[1]);
                            for (int i=0; i<addresses.length-1; i++) {
                                variable += addresses[i].getHostAddress();
                                variable += " ";
                            }
                        } catch (Exception e2) {
                            stack.ERROR = new Number(-1);
                        }
                    }
            }
        }

        if (Arguments.length > request.index(BigInteger.valueOf(0)).intValue()-1) {
            stack.share(new Array(Arguments[request.index(BigInteger.valueOf(0)).intValue()-1]));
            return;
        }

        stack.share(result);
        stack.ERROR = new Number(404);
        return;
    }


    static Pipe open(Stack stack) {

        Array request = stack.grab();
        String path = request.String();

        Pipe Pipe = new Pipe();
        Pipe.Name = path;

        //Load various protocols.
        String[] protocol = Pipe.Name.split("://", 2);
        if (protocol.length > 1) {
            switch (protocol[0]) {
                case "tcp":
                    ServerSocket server = Networks_In.get(protocol[1]);
                    if (server != null) {
                        try {
                            Socket client = server.accept();
                            Pipe.socket = client;
                            Pipe.input = client.getInputStream();
                            Pipe.output = client.getOutputStream();
                        } catch( Exception e) {
                            stack.push((Number)BigInteger.valueOf(-1));
                            return Pipe;
                        }
                        stack.push((Number)BigInteger.valueOf(0));
                        return Pipe;
                    } else {
                        String[] hostport = protocol[1].split(":", 2);
                        if (hostport.length > 1) {
                            //TODO
                            try {
                                Socket req = new Socket(hostport[0], (int)Integer.valueOf(hostport[1]));
                                Pipe.socket = req;
                                Pipe.input = req.getInputStream();
                                Pipe.output = req.getOutputStream();
                            } catch (Exception e) {
                                stack.push((Number)BigInteger.valueOf(-1));
                                return Pipe;
                            }
                            stack.push((Number)BigInteger.valueOf(0));
                            return Pipe;
                        }
                    }
            }
        }

        File file = new File(path);

        if (file.exists()) {
            stack.push((Number)BigInteger.valueOf(0));
            try {
                Pipe.input = new FileInputStream(file);
                Pipe.output = new FileOutputStream(file);
            }catch(Exception e){

            }
            return Pipe;
        }
        stack.push((Number)BigInteger.valueOf(-1));
        return Pipe;
    }

    static void info (Stack s) {
        String request = "";
        String variable = "";
        Array result = new Array();

        Pipe file = s.take();
        Array text = s.grab();

        for (int i = 0; i < text.size().intValue(); i++) {
            if (text.index(BigInteger.valueOf(i)) != null) {
                int c = text.index(BigInteger.valueOf(i)).intValue();
                request += (char)(c);
            }
        }

        switch (request) {
            case "ip":
                if (file.socket != null) {
                    variable = file.socket.getLocalAddress().getHostAddress();
                }
        }

        for (int i = 0; i < variable.length(); i++) {
            result.push((Number)BigInteger.valueOf(variable.charAt(i)));
        }
        s.share(result);
    }

    public static class Number extends BigInteger {
        public Number() {
            super("0");
        }

        public Number(int n) {
            super(String.valueOf(n));
        }

        //These functions can be inlined. Only here for reference.
        Number slt(Number b) {
            return compareTo(b) == -1 ? (Number)ONE : (Number)ZERO;
        }

        Number seq(Number b) {
            return compareTo(b) == 0 ? (Number)ONE : (Number)ZERO;
        }

        Number sge(Number b) {
            return compareTo(b) >= 0 ? (Number)ONE : (Number)ZERO;
        }

        Number sgt(Number b) {
            return compareTo(b) == 1 ? (Number)ONE : (Number)ZERO;
        }

        Number sne(Number b) {
            return compareTo(b) != 0 ? (Number)ONE : (Number)ZERO;
        }

        Number sle(Number b) {
            return compareTo(b) <= 0 ? (Number)ONE : (Number)ZERO;
        }

        Number div(Number b) {
            try {
                return (Number)divide(b);
            } catch (Exception e) {
                if (compareTo(Number.valueOf(0)) == 0) {
                    SecureRandom srand = new SecureRandom();
                    return new Number(srand.nextInt(255)+1);
                } else {
                    return (Number)ZERO;
                }
            }
        }

        Number mul(BigInteger b) {
            if (intValue() == 0 && b.intValue() == 0) {
                SecureRandom srand = new SecureRandom();
                return new Number(srand.nextInt(255)+1);
            }
            return (Number)multiply(b);
        }

        Number pow(Number b) {
            if (intValue() == 0) {
                if (b.mod(Number.valueOf(2)).intValue() != 0) {
                    SecureRandom srand = new SecureRandom();
                    return new Number(srand.nextInt(255)+1);
                }
                return (Number)ZERO;
            }
            return (Number)pow(b.intValue());
        }
    }

    //Pipe implementation.
    public static class Pipe {
        String Name;

        //Input and Output.
        InputStream input;
        OutputStream output;

        //Types of pipes.
        File file;
        Socket socket;

        void close() {
            try {
                input.close();
                output.close();
            }catch (Exception e)  {
            }
            try {
                socket.close();
            }catch (Exception e)  {
            }
        }
    }

    public static class PipeArray {
        ArrayList<Pipe> List;

        public PipeArray() {
            List = new ArrayList<Pipe>();
        }

        public Pipe pop() {
            return List.remove(List.size()-1);
        }
        public void push(Pipe n) {
            List.add(n);
        }
    }

    //An array of arrays, or a "Heap".
    public static class ArrayArray {
        ArrayList<Array> List;

        void push(Array n) {
            List.add(n);
        }

        public ArrayArray() {
            List = new ArrayList<Array>();
        }

        Array pop() {
            return List.remove(List.size()-1);
        }


        Number size() {
            return new Number(List.size());
        }

        Array index(Number n) {
            return List.get(n.intValue());
        }
    }

    public static class Array {
        ArrayList<Number> List;
        //new ArrayList<Integer>();

        void push(Number n) {
            List.add(n);
        }

        void set(Number index, Number n) {
            List.set(index.intValue(), n);
        }

        Array join(Array s) {
            Array newList = new Array();
            newList.List.addAll(List);
            newList.List.addAll(s.List);
            return newList;
        }

        public Array(Number... n) {
            List = new ArrayList<Number>();
            for (int i = 0; i < n.length; ++i) {
                List.add(n[i]);
            }
        }

        public Array(String s) {
            List = new ArrayList<Number>();
            for (int i = 0; i < s.length(); ++i) {
                List.add(new Number(s.charAt(i)));
            }
        }

        public String String() {
            String name = "";
            for (int i = 0; i < size().intValue(); i++) {
                if (index(BigInteger.valueOf(i)) != null) {
                    int c = index(BigInteger.valueOf(i)).intValue();
                    name += (char)(c);
                }
            }
            return name;
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

        int index(int n) {
            return List.get(n).intValue();
        }
        int length() {
            return List.size();
        }
    }
}
`
