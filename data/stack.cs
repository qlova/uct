using System;
using System.Collections;
using System.Collections.Generic;
using System.Numerics;
using System.Reflection;
using System.IO;
using System.Security.Cryptography;

public class stack {
	
	public Array Numbers = new Array();
	
	public ArrayArray Arrays = new ArrayArray();
	public PipeArray Pipes = new PipeArray();
	
	public Array ActiveArray = new Array();

	public static String[] ARGS; //Stores the programs arguments.
	
	public BigInteger ERROR = new BigInteger(0);

	public class PipeArray {
		List<Pipe> Value;
		
		public void push(Pipe n) {
			Value.Add(n);
		}
		
		public PipeArray() {
			Value = new List<Pipe>();
		}
		
		public Pipe pop() {
			Pipe temp;
			temp = Value[Value.Count-1];
			Value.RemoveAt(Value.Count-1);
			return temp;
		}
	}
	
	public class Pipe {
		public String Name;
		public StreamReader FileRead;
		public StreamWriter FileWrite;
		public MethodInfo Func;
		
		public Pipe() {}
		
		public Pipe(MethodInfo m) {
			Func = m;
		}
	}

	public class ArrayArray {
		List<Array> Value;
		
		public void push(Array n) {
			Value.Add(n);
		}
		
		public ArrayArray() {
			Value = new List<Array>();
		}
		
		public Array pop() {
			Array temp;
			temp = Value[Value.Count-1];
			Value.RemoveAt(Value.Count-1);
			return temp;
		}
		
		public BigInteger size() {
			return new BigInteger(Value.Count);
		}
		
		public Array index(BigInteger n) {
			return Value[(int)n];
		}
	}
	
	public class Array {
		List<BigInteger> Value;
		//new ArrayList<Integer>();
		
		public void push(BigInteger n) {
			Value.Add(n);
		}
		
		public void set(BigInteger index, BigInteger n) {
			Value[(int)index] = n;
		}
		
		public Array join(Array s) {
			Array newList = new Array(); 
			newList.Value.AddRange(Value);
			newList.Value.AddRange(s.Value);
			return newList;
		}
		
		public Array() {
			Value = new List<BigInteger>();
		}
		
		public Array(int size) {
			Value = new List<BigInteger>(new BigInteger[size]);
		}
		
		public Array(string n) {
			Value = new List<BigInteger>();
			for (int i = 0; i < n.Length; ++i) {
				Value.Add(n[i]);
			}
		}
		
		public Array(BigInteger[] n) {
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
	
	public Array array() {
		Array a = new Array();
		
		place(a);
		
		return a;
	}
	
	public void heap() {
		//Do this laterz.
	}
	
	public void share(Array n) {
		Arrays.push(n);
	}
	
	public void relay(Pipe n) {
		Pipes.push(n);
	}
	
	public void put(BigInteger a) {
		ActiveArray.push(a);
	}
	
	public void sets(BigInteger a) {
		ActiveArray.set(mod(pull(), ActiveArray.size()), a);
	}
	
	public BigInteger gets() {
		return ActiveArray.index(mod(pull(), ActiveArray.size()));
	}
	
	public void place(Array a) {
		ActiveArray = a;
	}
	
	public BigInteger pop() {
		return ActiveArray.pop();
	}
	
	
	public Array grab() {
		return Arrays.pop();
	}
	
	public Pipe take() {
		return Pipes.pop();
	}

	public void push(BigInteger n) {
		Numbers.push(n);
	}
	
	public BigInteger pull() {
		return Numbers.pop();
	}
	
	public void load() {
		String name = "";
		String variable = "";
		Array result = new Array();
	
		Array text = grab();
	
		if ((int)text.index(0) == 36 && (int)text.size() > 1) {
	
			for (int i = 1; i < (int)text.size(); i++) {
				int c = (int)text.index(new BigInteger(i));
				name += (char)(c);
			} 
			
			variable = Environment.GetEnvironmentVariable(name);
		} else {
			if (ARGS.Length > (int)text.index(0)) {
				variable = ARGS[(int)text.index(0)];
			} 
		}
		
		if (variable == null) {
			share(result);
			return;
		}
	
		for (int i = 0; i < variable.Length; i++) {
		    result.push(new BigInteger(variable[i]));
		}
		share(result);
	}
	
	public Pipe openit() {
		String filename = "";
		Array text = grab();
		for (int i = 0; i < (int)text.size(); i++) {
			int c = (int)text.index(new BigInteger(i));
			filename += (char)(c);
		} 
		
		Pipe it = new Pipe();
		it.Name = filename;
		
		if (File.Exists(filename)) {
			push(new BigInteger(0));
			try {
				var oStream = new FileStream(filename, FileMode.Append, FileAccess.Write, FileShare.Read); 
				var iStream = new FileStream(filename, FileMode.Open, FileAccess.Read, FileShare.ReadWrite);
				it.FileRead = new StreamReader(oStream);
				it.FileWrite = new StreamWriter(iStream);
			}catch(System.ArgumentException){
			}
			return it;
		}
		if (Directory.Exists(filename)) {
			push(new BigInteger(0));
			return it;
		}
		push(new BigInteger(-1));
		return it;
	}
	
	public void inn(Pipe file) {
		BigInteger length = pop();
		for (int i = 0; i < (int)length; i++) {
			/*try {*/
				push(new BigInteger(file.FileRead.Read()));
			/*}catch(IOException e){
				push(new BigInteger(-1));
			}*/
		}
	}
	
	public void outy(Pipe file) {
		Array text = grab();
		
		if ((int)text.size() == 0 || file.FileWrite == null ) {
			if (file.Name[file.Name.Length-1] == '/') {
				if (Directory.Exists(file.Name)) {
					
				} else {
					try {
 						Directory.CreateDirectory(file.Name);
						push(new BigInteger(0));
						return;
					} catch {
						push(new BigInteger(-1));
						return;
					}
				}
			} else if (File.Exists(file.Name)) {
				
			} else {
				try {
					File.Create(file.Name);
					push(new BigInteger(0));
					return;
				} catch  {
					push(new BigInteger(-1));
					return;
				}
			}
		}

		
		for (int i = 0; i < (int)text.size(); i++) {
			int c = (int)text.index(new BigInteger(i));
			file.FileWrite.Write((char)(c));
		} 
		push(new BigInteger(0));
	}
	
	public void close(Pipe file) {
		try {
			file.FileRead.Close();
			file.FileWrite.Close();
		}catch{
		}
	}

	public void stdout() {
		Array text = grab();
		for (int i = 0; i < (int)text.size(); i++) {
			int c = (int)text.index(new BigInteger(i));
			System.Console.Write((char)(c));
		} 
	}
	
	public void stdin() {
		BigInteger mode = pull();
		
		Array result = new Array();
		if (mode > 0) {
			for (int i = 0; i < (int)mode; i++) {
				int c = Console.Read();
				if (c == -1) {
					ERROR = 1;
					break;
				}
				result.push(new BigInteger(c));
			}
		} else if (mode == 0) {
			while (true) {
				int c = Console.Read();
				if (c == -1) {
					ERROR = 1;
					break;
				}
				if (c == (int)'\n') {
					break;
				}
				result.push(new BigInteger(c));
			}
		} else {
			while (true) {
				int c = Console.Read();
				if (c == -1) {
					ERROR = 1;
					break;
				}
				if (c == -((int)mode)) {
					break;
				}
				result.push(new BigInteger(c));
			}
		}
		
		share(result);
	}
	
	public static BigInteger slt(BigInteger a, BigInteger b) {
		if (a.CompareTo(b) == -1) {
			return new BigInteger(1);
		}
		return new BigInteger(0);
	}
	
	public static BigInteger seq(BigInteger a, BigInteger b) {
		if (a.CompareTo(b) == 0) {
			return new BigInteger(1);
		}
		return new BigInteger(0);
	}
	
	public static BigInteger sge(BigInteger a, BigInteger b) {
		if (a.CompareTo(b) >= 0) {
			return new BigInteger(1);
		}
		return new BigInteger(0);
	}
	
	public static BigInteger sgt(BigInteger a, BigInteger b) {
		if (a.CompareTo(b) == 1) {
			return new BigInteger(1);
		}
		return new BigInteger(0);
	}
	
	public static BigInteger sne(BigInteger a, BigInteger b) {
		if (a.CompareTo(b) != 0) {
			return new BigInteger(1);
		}
		return new BigInteger(0);
	}
	
	public static BigInteger sle(BigInteger a, BigInteger b) {
		if (a.CompareTo(b) <= 0) {
			return new BigInteger(1);
		}
		return new BigInteger(0);
	}
	
	public static BigInteger div(BigInteger a, BigInteger b) {
		try {
			return a/b;
		} catch {
			if (a.CompareTo(new BigInteger(0)) == 0) {
				RandomNumberGenerator rng = RNGCryptoServiceProvider.Create();
				 byte[] randomNumber = new byte[1];
				 rng.GetBytes(randomNumber);
				return new BigInteger((int)randomNumber[0]);
			} else {
				return new BigInteger(0);
			}
		}
	}
	
	public static BigInteger mod(BigInteger a, BigInteger b) {
		try {
			return ((a%b) + b) % b ;
		} catch {
			if (a.CompareTo(new BigInteger(0)) == 0) {
				RandomNumberGenerator rng = RNGCryptoServiceProvider.Create();
				 byte[] randomNumber = new byte[1];
				 rng.GetBytes(randomNumber);
				return new BigInteger((int)randomNumber[0]);
			} else {
				return new BigInteger(0);
			}
		}
	}
	
	public static BigInteger mul(BigInteger a, BigInteger b) {
		return a*b;
	}
	
	public static BigInteger pow(BigInteger a, BigInteger b) {
		return BigInteger.Pow(a, (int)b);
	}
}
