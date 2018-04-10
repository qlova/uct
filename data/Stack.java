//Compiled to Java with UCT (Universal Code Translator)

//Import java libraries.
import java.math.BigInteger; 	//Support numbers of any size.
import java.util.Hashtable;  	//Hashtable is a useful utilPipey.
import java.util.ArrayList;  	//ArrayLists are helpful.
import java.util.List;
import java.io.*;				//Deal with files.
import java.net.*;				//Deal with the network.
import java.lang.reflect.*;		//Reflection for methods.
import java.util.Scanner;

//Random numbers.
import java.security.SecureRandom;

//This is for threading.
import java.util.concurrent.Executor;
import java.util.concurrent.Executors;
import java.util.concurrent.ArrayBlockingQueue;

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
    Array 				Numbers;
    ArrayArray 			Arrays;
    PipeArray 			Pipes;

    Number 				ERROR;

    Array				ActiveArray;
    ArrayArray			Heap;
    List<Integer>	HeapRoom;
    
    PipeArray			HeapIt;
    List<Integer>	HeapItRoom;
    
    Hashtable<String, Number> Map;

    static Opener              opener;
    static Loader               loader;
    
    ArrayBlockingQueue<Array> inbox;
    ArrayBlockingQueue<Array> outbox;
    

    //This hashtable keeps track of Servers currently listening on the specified port.
    static Hashtable<String, ServerSocket> Networks_In = new Hashtable<String, ServerSocket>();

    //This will store the system arguments.
    public static String[] Arguments;

    //This is the threading pool.
    public static Executor ThreadPool = Executors.newCachedThreadPool();


    interface Opener
    {
        Pipe open(String uri);
    }
    static public void SetOpener(Opener o) {
        opener = o;
    }

    interface Loader {
        Array load(String uri);
    }
    static public void SetLoader(Loader o) {
        loader = o;
    }

    //This creates an empty stack. Ready for use.
    public Stack() {
        Numbers 		= new Array();
        Arrays 			= new ArrayArray();
        Pipes 		    = new PipeArray();
        ERROR 			= new Number();

        Heap			= new ArrayArray();
        HeapRoom		= new ArrayList<Integer>();
        
		HeapIt			= new PipeArray();
        HeapItRoom		= new ArrayList<Integer>();
        Map				= new Hashtable<String, Number>();
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
        
        n.outbox = new ArrayBlockingQueue<Array>(1);
        inbox = n.outbox;

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
    
    //CONNECT
    void connect() {
        Number num = Map.get(grab().String());
        if (num == null) {
            push(new Number(0));
        	ERROR = new Number(1);
        } else {
        	push(num);
        }
    }

    //LINK 
    void link() {
        Map.put(grab().String(), pull());
    }
    
    //SLICE 
    void slice() {
        Array result = new Array();
        result.List = grab().List.subList(pull().intValue(), pull().intValue());
        share(result);
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
    Number get() {
    	if (ActiveArray.List.size() == 0) {
    		ERROR = new Number(4);
    		return new Number(0);
    	}
        return ActiveArray.index(pull().mod(ActiveArray.size()));
    }

    //SET number
    void set(Number b) {
    	if (ActiveArray.List.size() == 0) {
    		ERROR = new Number(4);
    		return;
    	}
        ActiveArray.set(pull().mod(ActiveArray.size()), b);
    }

    //POP number
    Number pop() {
        return ActiveArray.pop();
    }

    //PULL number
    Number pull() {
        return Numbers.pop();
    }
    
    void heap() {
    	Number address = pull();
	
		if (address.intValue() == 0) {
			if (HeapRoom.size() > 0) {
				Integer address2 = HeapRoom.remove(HeapRoom.size()-1);
				
				Heap.List.set(((address2)%(Heap.List.size()+1)-1),  grab());
				push(new Number(address2));
			} else {
				Heap.push(grab());
				push(Heap.size());
			}
			
		} else if (address.intValue() > 0) {
			share(Heap.List.get((address.intValue()%(Heap.List.size()+1))-1));
			
		} else if (address.intValue() < 0) {
			Heap.List.set(((-address.intValue())%(Heap.List.size()+1)-1),  null);
			HeapRoom.add(-address.intValue());
		}
    }
    
    void heapit() {
    	Number address = pull();
	
		if (address.intValue() == 0) {
			if (HeapItRoom.size() > 0) {
				Integer address2 = HeapItRoom.remove(HeapItRoom.size()-1);
				
				HeapIt.List.set(((address2)%(HeapIt.List.size()+1)-1),  take());
				push(new Number(address2));
			} else {
				HeapIt.push(take());
				push(HeapIt.size());
			}
			
		} else if (address.intValue() > 0) {
			relay(HeapIt.List.get((address.intValue()%(HeapIt.List.size()+1))-1));
			
		} else if (address.intValue() < 0) {
			HeapIt.List.set(((-address.intValue())%(HeapIt.List.size()+1)-1),  null);
			HeapItRoom.add(-address.intValue());
		}
    }

    void stdout() {
        Array text = grab();
        for (int i = 0; i < text.size().intValue(); i++) {
            if (text.index(new Number(i)) != null) {
                int c = text.index(new Number(i)).intValue();
                System.out.print((char)(c));
            }
        }
    }

    void stdin() {
        Number length = pull();
        
        //This is the mode we use.
        // >0 is number of bytes to read.
        // <0 is character to read.
        // 0 is read a line.
        if (length.compareTo(new Number(0)) == 0) {
        
        	byte[] b = new byte[1];
        	
        	String input = "";
        	
        	while (true) {
        		try {
        			int n = System.in.read(b);
        		} catch (Exception e) {
        			ERROR = new Number(-1);
        			share(new Array(input));
        			return;
        		}
        		if (b[0] == '\n') {
        			share(new Array(input));
        			return;
        		}
        		input += ((char)(b[0]));
        	}
        	
        	
        } else if (length.compareTo(new Number(0)) == -1)  {
        
       		byte[] b = new byte[1];
        	
        	String input = "";
        	
        	while (true) {
        		try {
        			int n = System.in.read(b);
        		} catch (Exception e) {
        			ERROR = new Number(-1);
        			share(new Array(input));
        			return;
        		}
        		if (b[0] == -length.intValue()) {
        			share(new Array(input));
        			return;
        		}
        		input += ((char)(b[0]));
        	}
        
        
        } else { //length is > 0
        	
        	
        	byte[] b = new byte[length.intValue()];
        	
    		try {
    			int n = System.in.read(b);
    		} catch (Exception e) {
    			ERROR = new Number(-1);
    			share(new Array(b));
    			return;
    		}
        	
			 share(new Array(b));
			 return;
        }
       
    }
    
    public void scope(String name, Class<?> c) {
		Class[] cArg = new Class[1]; 
		cArg[0] = Stack.class; 
		try { 
			relay(new Stack.Pipe(c.getDeclaredMethod(name, cArg))); 
		} catch (NoSuchMethodException e) { 
			throw new RuntimeException(e); 
		}
    }

    void in() {
        Pipe file = take();
        Number length = pull();

        if (file.external != null) {
            share(file.external.In(length.longValue()));
            return;
        }
        
        if (file.udp && file.udpsock != null) {
        	//TODO store this as a static member of stack to avoid the creation of huge buffers all the time.
        	byte[] buffer = new byte[65535];
        	DatagramPacket packet = new DatagramPacket(buffer, 65535);
        	try {
        		((MulticastSocket)file.udpsock).receive(packet);
        	} catch (Exception e) {
        		ERROR = new Number(404);
        		return;
        	}
        	
        	String message = "";
        	try {
        		message = new String(buffer, "UTF-8");
        	} catch (Exception e) {
        		ERROR = new Number(404);
        		return;
        	}
        	
        	file.ip = packet.getAddress().getHostAddress();
        	share(new Array(message));
        	
        	return;
        }
       
       	//This is the mode we use.
        // >0 is number of bytes to read.
        // <0 is character to read.
        // 0 is read a line.
        if (length.compareTo(new Number(0)) == 0) {
        
        	byte[] b = new byte[1];
        	
        	String input = "";
        	
        	while (true) {
        		try {
        			int n = file.input.read(b);
        			
        			if (n < 0) {
		    			ERROR = new Number(-1);
		    			share(new Array(input));
		    			return;
		    		}
        			
        		} catch (Exception e) {
        			ERROR = new Number(-1);
        			share(new Array(input));
        			return;
        		}
        		
        		
        		if (b[0] == '\n') {
        			share(new Array(input));
        			return;
        		}
        		input += ((char)( b[0] & 0xff));
        	}
        	
        	
        } else if (length.compareTo(new Number(0)) == -1)  {
        
       		byte[] b = new byte[1];
        	
        	String input = "";
        	
        	while (true) {
        		try {
        			int n = file.input.read(b);
        			
        			if (n < 0) {
		    			ERROR = new Number(-1);
		    			share(new Array(input));
		    			return;
		    		}
        			
        		} catch (Exception e) {
        			ERROR = new Number(-1);
        			share(new Array(input));
        			return;
        		}
        		
        		if (b[0] == -length.intValue()) {
        			share(new Array(input));
        			return;
        		}
        		input += ((char)(b[0] & 0xff));
        	}
        
        
        } else { //length is > 0
        	
        	
        	byte[] b = new byte[length.intValue()];
        	
    		try {
    			int n = file.input.read(b);
    		} catch (Exception e) {
    			ERROR = new Number(-1);
    			share(new Array(b));
    			return;
    		}
        	
			 share(new Array(b));
			 return;
        }
    }

    void out() {
        Pipe file = take();
        Array text = grab();

        if (file.external != null) {
            file.external.Out(text);
            return;
        }
        
        if (file.udp && file.udpsock!= null) {
        	try {
		    	DatagramPacket packet = new DatagramPacket(text.String().getBytes(), text.size().intValue(), file.group, file.port);
		    	((MulticastSocket)file.udpsock).send(packet);
		    	
		    } catch (Exception e) {
		    	ERROR = new Number(404);
		    }
		    return;
        } 

        if (file.output != null) {
            //TODO optimise to send in a single packet.
            for (int i = 0; i < text.size().intValue(); i++) {
                if (text.index(new Number(i)) != null) {
                    int c = text.index(new Number(i)).intValue();
                    try {
                        file.output.write((char)(c));
                    }catch(Exception e){
                        push(new Number(-1));
                    }
                }
            }
            push(new Number(0));
            return;
        }

        if (file.Name.length() > 0 && (text.size().intValue() == 0 || file.output == null )) {
            if (file.Name.charAt(file.Name.length()-1) == '/') {
                if (new File(file.Name).exists()) {

                } else {
                    try {
                        File f = new File(file.Name);
                        if (!f.mkdir()) {
                            push(new Number(-1));
                            return;
                        }
                        push(new Number(0));
                        return;
                    } catch (Exception e) {
                        push(new Number(-1));
                        return;
                    }
                }
            } else if (new File(file.Name).exists()) {

            } else {
                try {
                    File newfile = new File(file.Name);
                    newfile.createNewFile();
                    try {
					    file.input = new FileInputStream(newfile);
					    file.output = new FileOutputStream(newfile, true);
					}catch(Exception e){
					    push(new Number(-1));
						return;
					}
                } catch (Exception e)  {
                    push(new Number(-1));
                    return;
                }
            }
        }

        for (int i = 0; i < text.size().intValue(); i++) {
            if (text.index(new Number(i)) != null) {
                int c = text.index(new Number(i)).intValue();
                try {
                    file.output.write((char)(c));
                }catch(Exception e){
                    push(new Number(-1));
                }
            }
        }
        push(new Number(0));
    }
    
    void execute() {
    	String command = grab().String();
    	String variable = "";
    	try {
			
			String[] complete = {"sh","-c", command};
			Scanner s = new Scanner(Runtime.getRuntime().exec(complete).getInputStream());
			s.useDelimiter("\\A");
    	
			variable = s.next();
		} catch (Exception e) {
			ERROR = new Number(1);
		}
		share(new Array(variable));
    }
    
    void delete() {
    	String name = grab().String();
    	
    	if (name.equals("")) {
    		ERROR = new Number(1);
    		return;
    	}
    	
    	File file = new File(name);
    	
    	if(!file.delete()) { 
    		ERROR = new Number(1);
    		return;
		}
    }


    //LOAD
    void load() {
        String name = "";
        String variable = "";

        //This request is what we need to load.
        Array request = grab();

        if (loader != null) {
            Array arr = loader.load(request.String());
            if (arr != null) {
                share(arr);
                return;
            }
        }

        Array result = new Array();

        //The request is an enviromental variable.
        if (request.index(0) == '$' && request.length() > 1) {

            //We parse the rest of the string.
            for (int i = 1; i < request.length(); i++) {
                try {
                    name += (char)(request.index(i));
                } catch (Exception e) {

                }
            }
            share(new Array(System.getenv(name)));
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
                        share(new Array(variable));
                        return;
                    } catch (Exception e) {
                        ERROR = new Number(-1);
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
                        share(new Array(variable));
                        return;
                    } catch (Exception e) {
                        try {
                            InetAddress[] addresses = InetAddress.getAllByName(protocol[1]);
                            for (int i=0; i<addresses.length-1; i++) {
                                variable += addresses[i].getHostAddress();
                                variable += " ";
                            }
                            share(new Array(variable));
                        return;
                        } catch (Exception e2) {
                            ERROR = new Number(-1);
                        }
                    }
            }
        }

        if (Arguments.length > request.index(new Number(0)).intValue()-1) {
            share(new Array(Arguments[request.index(new Number(0)).intValue()-1]));
            return;
        } else {
        	ERROR = new Number(1);
        	share(result);
        	return;
        }
    }


    void open() {

        Array request = grab();
        String path = request.String();

        if (opener != null) {
            Pipe pipe = opener.open(path);
            if (pipe != null) {
                relay(pipe);
                return;
            }
        }

        Pipe pipe = new Pipe();
        pipe.Name = path;
        
        if (path.equals("")) {
        	relay(pipe);
        	ERROR = new Number(1);
        	return;
        }

        //Load various protocols.
        String[] protocol = pipe.Name.split("://", 2);
        if (protocol.length > 1) {
            switch (protocol[0]) {
            	case "multicast":
            		String[] hostport = protocol[1].split(":", 2);
            		if (hostport.length > 1) {
            			try {
		        			MulticastSocket client = new MulticastSocket(Integer.valueOf(hostport[1]));
		        			InetAddress group = InetAddress.getByName(hostport[0]);
		        			
		        			client.joinGroup(group);
		        			
		        			pipe.udpsock = client;
		        			pipe.udp = true;
		        			pipe.group = group;
		        			pipe.port = Integer.valueOf(hostport[1]);
		        			
		        			relay(pipe);
		        			return;
		        			
		        		} catch (Exception e) {
		        			ERROR = new Number(404);
		        			relay(new Pipe());
		        		}
            		}
            		break;
                case "tcp":
                    ServerSocket server = Networks_In.get(protocol[1]);
                    if (server != null) {
                        try {
                            Socket client = server.accept();
                            pipe.socket = client;
                            pipe.input = client.getInputStream();
                            pipe.output = client.getOutputStream();
                        } catch( Exception e) {
                            ERROR = new Number(404);
                            relay(pipe);
                            return;
                        }
                        relay(pipe);
                        return;
                    } else {
                        hostport = protocol[1].split(":", 2);
                        if (hostport.length > 1) {
                            //TODO
                            try {
                                Socket req = new Socket(hostport[0], (int)Integer.valueOf(hostport[1]));
                                pipe.socket = req;
                                pipe.input = req.getInputStream();
                                pipe.output = req.getOutputStream();
                            } catch (Exception e) {
                                ERROR = new Number(404);
                                relay(pipe);
                                return;
                            }
                           	relay(pipe);
                           	return;
                        }
                    }
            }
        }

        File file = new File(path);

        /*if (!file.exists()) {
        	try {
        		file.createNewFile();
			}catch(Exception e){
                push(new Number(-1));
    			relay(pipe);
            }
        }*/
        if (!file.exists()) {
        	ERROR = new Number(404);
			relay(pipe);
			return;
        }
            try {
                pipe.input = new FileInputStream(file);
                pipe.output = new FileOutputStream(file, true);
            }catch(Exception e){
    			relay(pipe);
    			return;
            }
            relay(pipe);
            return;
        
       
    }

    void info () {
        String request = "";
        String variable = "";
        Array result = new Array();

        Pipe file = take();
        Array text = grab();

        for (int i = 0; i < text.size().intValue(); i++) {
            if (text.index(new Number(i)) != null) {
                int c = text.index(new Number(i)).intValue();
                request += (char)(c);
            }
        }

        switch (request) {
            case "ip":
                if (file.socket != null || file.udpsock != null) {
                	if (file.udp) {
                		variable = file.ip;
                	} else {
                    	variable = file.socket.getLocalAddress().getHostAddress();
                    }
                }
        }

        for (int i = 0; i < variable.length(); i++) {
            result.push(new Number(variable.charAt(i)));
        }
        share(result);
    }

    public static class Number {
        BigInteger a;

        public Number() {
            a = BigInteger.ZERO;
        }

        public Number(int n) {
            a = BigInteger.valueOf(n);
        }

        public Number(BigInteger n) {
            a = n;
        }
        
        public Number(String s) {
            a = new BigInteger(s);
        }

        public int compareTo(Number b) {
            return a.compareTo(b.a);
        }

        public int intValue() {
            return a.intValue();
        }

        public long longValue() { return a.longValue(); }

        //These functions can be inlined. Only here for reference.
        Number slt(Number b) {
            return new Number(compareTo(b) == -1 ? BigInteger.ONE : BigInteger.ZERO);
        }

        Number seq(Number b) {
            return new Number(compareTo(b) == 0 ? BigInteger.ONE : BigInteger.ZERO);
        }

        Number sge(Number b) {
            return new Number(compareTo(b) >= 0 ? BigInteger.ONE : BigInteger.ZERO);
        }

        Number sgt(Number b) {
            return new Number(compareTo(b) == 1 ? BigInteger.ONE : BigInteger.ZERO);
        }

        Number sne(Number b) {
            return new Number(compareTo(b) != 0 ? BigInteger.ONE : BigInteger.ZERO);
        }

        Number sle(Number b) {
            return new Number(compareTo(b) <= 0 ? BigInteger.ONE : BigInteger.ZERO);
        }

        Number add(Number b) {
            return new Number(a.add(b.a));
        }

        Number sub(Number b) {
            return new Number(a.subtract(b.a));
        }

        Number mod(Number b) {
            try {
                return new Number(a.mod(b.a));
            } catch (Exception e) {
                return new Number(BigInteger.ZERO);
            }
        }

        Number div(Number b) {
            try {
                return new Number(a.divide(b.a));
            } catch (Exception e) {
                if (compareTo(new Number(0)) == 0) {
                    SecureRandom srand = new SecureRandom();
                    return new Number(srand.nextInt(255)+1);
                } else {
                    return new Number(BigInteger.ZERO);
                }
            }
        }

        Number mul(Number b) {
            return new Number(a.multiply(b.a));
        }

        Number pow(Number b) {
            return new Number(a.pow(b.intValue()));
        }
    }

    //Pipe implementation.
    public static class Pipe {
		Array data;
		Method method;
		
		//Make these into interface.
			String Name;

			//Input and Output.
			InputStream input;
			OutputStream output;
			
			//UDP stuff.
		MulticastSocket udpsock;
			boolean udp;
			InetAddress group;
			int port;
			String ip;

			//Types of pipes.
			File file;
			Socket socket;

        Pipeable external;

        public interface Pipeable {
            Array In(long length);
            void Out(Array data);
        }
        
        public Pipe() {
			data = new Array();
        }

        void Set(Pipeable p) {
            external = p;
        }
        
        
        public Pipe(Method m) {
			data = new Array();
        	method = m;
        }
        
        void exe(Stack stack) {
        	try {
        		method.invoke(null, stack); 
        	} catch (Exception e) {  
        		throw new RuntimeException(e);
        	}
        }

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
        List<Pipe> List;

        public PipeArray() {
            List = new ArrayList<Pipe>();
        }
        
        public Number size() {
            return new Number(List.size());
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
        List<Array> List;

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
        List<Number> List;
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
        
        public Array(int n) {
            List = java.util.Arrays.asList(new Number[n]);
        }
        
        public Array(byte[] s) {
            List = new ArrayList<Number>();
            for (int i = 0; i < s.length; ++i) {
                List.add(new Number(s[i] & 0xff));
            }
        }

        public String String() {
            String name = "";
            for (int i = 0; i < size().intValue(); i++) {
                if (index(new Number(i)) != null) {
                    int c = index(new Number(i)).intValue();
                    name += (char)(c);
                }
            }
            return name;
        }

        Number pop() {
            return List.remove(List.size()-1);
        }

        Number size() {
            return new Number(List.size());
        }

        Number index(Number n) {
        	if (List.get(n.intValue()) == null) {
        		return new Number(0);
        	}
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
