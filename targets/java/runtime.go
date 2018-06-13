package java

import uct "github.com/qlova/uct/assembler"

var Runtime = uct.Runtime{
	Name: "Runtime.java",
	Data: `import java.math.BigInteger;
import java.lang.reflect.*;
	
class Runtime { 
	Int Error = new Int();
	List Global = new List();
	Pipe Channel;
	
	Int[] Stack = new Int[20];
	List[] Lists = new List[20];;
	Pipe[] Pipes = new Pipe[20];;
	
	int StackPointer = 0;
	int ListsPointer = 0;
	int PipesPointer = 0;
	
	Int[] TheHeap = new Int[1];
	int[] TheHeapRoom = new int[1];
	
	List[] TheListHeap = new List[1];
	int[] TheListHeapRoom = new int[1];
	
	int TheListHeapPointer = 0;
	int TheListHeapRoomPointer = 0;
	
	Pipe[] ThePipeHeap = new Pipe[1];
	int[] ThePipeHeapRoom = new int[1];
	
	void push(Int i) {
		StackPointer++;
		if (Stack.length == StackPointer) {
			Int[] tmp = new Int[Stack.length*2];
			System.arraycopy(Stack, 0, tmp, 0, Stack.length);
			Stack = tmp;
		}
		Stack[StackPointer] = i;
	}
	
	Int pull() {
		Int result = Stack[StackPointer];
		StackPointer--;
		return result;
	}
	
	void pushList(List l) {
		ListsPointer++;
		if (Lists.length == ListsPointer) {
			List[] tmp = new List[Lists.length*2];
			System.arraycopy(Lists, 0, tmp, 0, Lists.length);
			Lists = tmp;
		}
		Lists[ListsPointer] = l;
	}
	
	void heapList() {
    	Int address = pull();
	
		if (address.Small == 0) {
			if (TheListHeapRoomPointer > 0) {
				
				int free_spot = TheListHeapRoom[TheListHeapRoomPointer];
				TheListHeapRoomPointer--;
				
				TheListHeap[free_spot] = Lists[ListsPointer];
				ListsPointer--;
				
				push(new Int((long)free_spot));
			} else {
				
				TheListHeapPointer++;
				if (TheListHeap.length == TheListHeapPointer) {
					List[] tmp = new List[TheListHeap.length*2];
					System.arraycopy(TheListHeap, 0, tmp, 0, TheListHeap.length);
					TheListHeap = tmp;
				}
				TheListHeap[TheListHeapPointer] = Lists[ListsPointer];
				ListsPointer--;

				push(new Int((long)TheListHeapPointer));
			}
			
		} else if (address.Small > 0) {
			pushList(TheListHeap[(int)address.Small]);
			
		} else if (address.Small < 0) {
			
			TheListHeapRoomPointer++;
			if (TheListHeapRoom.length == TheListHeapRoomPointer) {
				int[] tmp = new int[TheListHeapRoom.length*2];
				System.arraycopy(TheListHeapRoom, 0, tmp, 0, TheListHeapRoom.length);
				TheListHeapRoom = tmp;
			}
			TheListHeapRoom[TheListHeapRoomPointer] = (int)-address.Small;
			
			TheListHeap[TheListHeapPointer] = null;
			TheListHeapPointer--;
		}
    }
	
	void pushPipe(Pipe p) {
		PipesPointer++;
		if (Pipes.length == PipesPointer) {
			Pipe[] tmp = new Pipe[Pipes.length*2];
			System.arraycopy(Pipes, 0, tmp, 0, Pipes.length);
			Pipes = tmp;
		}
		Pipes[PipesPointer] = p;
	}
	
	void sub() {
		Int a = Stack[StackPointer-1];
		Int b = Stack[StackPointer];
		StackPointer -= 2;
		
		Int result  = new Int();
		
		if (a.Big == null || b.Big == null) {
			if (a.Big == null && b.Big == null) {
				try {
					result.Small = Math.subtractExact(a.Small, b.Small);
				} catch (ArithmeticException e) {
					result.Big = BigInteger.valueOf(a.Small).subtract(BigInteger.valueOf(b.Small));
				}
			}
		} else {
			result.Big = a.Big.subtract(b.Big);
		}
		
		push(result);
	}
	
	void add() {
		Int a = Stack[StackPointer];
		Int b = Stack[StackPointer-1];
		StackPointer -= 2;
		
		Int result = new Int();
		
		if (a.Big == null || b.Big == null) {
			if (a.Big == null && b.Big == null) {
				try {
					result.Small = Math.addExact(a.Small, b.Small);
				} catch (ArithmeticException e) {
					result.Big = BigInteger.valueOf(a.Small).add(BigInteger.valueOf(b.Small));
				}
			}
		} else {
			result.Big = a.Big.add(b.Big);
		}
		
		push(result);
	}
	
	void div() {
		Int a = Stack[StackPointer-1];
		Int b = Stack[StackPointer];
		StackPointer -= 2;
		
		Int result = new Int();
		
		if (a.Big == null || b.Big == null) {
			if (a.Big == null && b.Big == null) {
				//TODO check for case when there is overflow.
				
				try {
					result.Small = a.Small/b.Small;
				} catch (Exception e) {
					if (a.Small == 0) {
						result.Small = 1;
					} else {
						result.Small = 0;
					}
				}
			}
		} else {
			result.Big = a.Big.divide(b.Big);
		}
		
		push(result);
	}
	
	void mul() {
		Int a = Stack[StackPointer];
		Int b = Stack[StackPointer-1];
		StackPointer -= 2;
		
		Int result = new Int();
		
		if (a.Big == null || b.Big == null) {
			if (a.Big == null && b.Big == null) {
				try {
					result.Small = Math.multiplyExact(a.Small, b.Small);
				} catch (ArithmeticException e) {
					result.Big = BigInteger.valueOf(a.Small).multiply(BigInteger.valueOf(b.Small));
				}
			} else if (a.Big == null) {
				result.Big = BigInteger.valueOf(a.Small).multiply(b.Big);
			} else {
				result.Big = a.Big.multiply(BigInteger.valueOf(b.Small));
			}
		} else {
			result.Big = a.Big.multiply(b.Big);
		}
		
		push(result);
	}
	
	void mod() {
		Int a = Stack[StackPointer-1];
		Int b = Stack[StackPointer];
		StackPointer -= 2;
		
		Int result = new Int();
		
		if (a.Big == null || b.Big == null) {
			if (a.Big == null && b.Big == null) {
				//TODO check for case when there is overflow.
				result.Small = a.Small%b.Small;
			}
		} else {
			result.Big = a.Big.mod(b.Big);
		}
		
		push(result);
	}
	
	void pow() {
		Int a = Stack[StackPointer-1];
		Int b = Stack[StackPointer];
		StackPointer -= 2;
		
		Int result = new Int();
		
		if (a.Big == null || b.Big == null) {
			if (a.Big == null && b.Big == null) {
				//TODO check for case when there is overflow.
				result.Big = BigInteger.valueOf(a.Small).pow((int)b.Small);
			}
		} else {
			result.Big = a.Big.mod(b.Big);
		}
		
		push(result);
	}
	
	void same() {
		Int a = Stack[StackPointer];
		Int b = Stack[StackPointer-1];
		StackPointer -= 2;
		
		Int result = new Int();
		
		if (a.Big == null || b.Big == null) {
			if (a.Big == null && b.Big == null) {
				if (a.Small == b.Small) {
					result.Small = 1;
				}
			}
		} else {
			if (a.Big.compareTo(b.Big) == 0) {
				result.Small = 1;
			}
		}
		
		push(result);
	}
	
	void more() {
		Int a = Stack[StackPointer];
		Int b = Stack[StackPointer-1];
		StackPointer -= 2;
		
		Int result = new Int();
		
		if (a.Big == null || b.Big == null) {
			if (a.Big == null && b.Big == null) {
				if (a.Small > b.Small) {
					result.Small = 1;
				}
			}
		} else {
			if (a.Big.compareTo(b.Big) == 1) {
				result.Small = 1;
			}
		}
		
		push(result);
	}
	
	void less() {
		Int a = Stack[StackPointer];
		Int b = Stack[StackPointer-1];
		StackPointer -= 2;
		
		Int result = new Int();
		
		if (a.Big == null || b.Big == null) {
			if (a.Big == null && b.Big == null) {
				if (a.Small < b.Small) {
					result.Small = 1;
				}
			}
		} else {
			if (a.Big.compareTo(b.Big) == -1) {
				result.Small = 1;
			}
		}
		
		push(result);
	}
	
	void swap() {
		Int tmp = Stack[StackPointer];
		Stack[StackPointer] = Stack[StackPointer-1];
		Stack[StackPointer-1] = tmp;
	}
	
	void swapList() {
		List tmp = Lists[ListsPointer];
		Lists[ListsPointer] = Lists[ListsPointer-1];
		Lists[ListsPointer-1] = tmp;
	}
	
	void swapPipe() {
		Pipe tmp = Pipes[PipesPointer];
		Pipes[PipesPointer] = Pipes[PipesPointer-1];
		Pipes[PipesPointer-1] = tmp;
	}
	
	void load() {
		List uri = Lists[ListsPointer];
		ListsPointer--;
		
		List result = new List(null, null);
		String name = new String();
		String variable = new String();
		
		if (uri.Bytes.length > 1 && (uri.Bytes[0] & 0xff) == '$') {
			for (int i = 0; i < uri.Bytes.length; i++) {
				if (i == 0) {
					continue;
				}
				name += (char)uri.Bytes[i];
			}
			variable = System.getenv(name);
		}
		
		if (variable == null) {
			result.Bytes = new byte[1];
			result.BytesSize = 0;
		} else {
			result.Bytes = variable.getBytes();
			result.BytesSize = result.Bytes.length;
		}
		
		pushList(result);
	}
	
	void open() {
		List uri = Lists[ListsPointer];
		ListsPointer--;
		
		if (uri.BytesSize == 0) {
			pushPipe(new StdPipe());
		} else {
			pushPipe(null);
			Error.Small = 1;
			Error.Big = null;
		}
	}
	
	void read() {
		Pipe pipe = Pipes[PipesPointer];
		PipesPointer--;
	
		Int size = Stack[StackPointer];
		StackPointer--;
		
		if (size.Small == 0) {
			
		} else if (size.Small < 0) {
			
			byte[] b = new byte[1];
			List l = new List(new byte[1], null);
			l.BytesSize = 0;
			
			while (true) {
				int amount = pipe.read(this, b);
				//TODO check errors here.
				
				if ((b[0] & 0xff) == -size.Small) {
					pushList(l);
					break;
				}
				
				if (l.Bytes.length == l.BytesSize) {
					byte[] tmp = new byte[l.Bytes.length*2];
					System.arraycopy(l.Bytes, 0, tmp, 0, l.Bytes.length);
					l.Bytes = tmp;
				}
				l.Bytes[l.BytesSize] = b[0];
				l.BytesSize++;
			}
			
		} else {
			byte[] data = new byte[(int)size.Small];
			
			int amount = pipe.read(this, data);
			
			List result = new List();
			result.Bytes = data;
			
			pushList(result);
		}
	}
	
	void send() {
		Pipe pipe = Pipes[PipesPointer];
		PipesPointer--;
		
		List data = Lists[ListsPointer];
		ListsPointer--;
		
		if (pipe == null) {
			return;
		}
		
		try {
			pipe.write(this, data.Bytes);
		} catch (Exception e) {
			Error.Small = 1;
			Error.Big = null;
		}
	}

static class Int {
	long Small = 0;
	BigInteger Big;
	
	Int() {}
	
	Int(long n) {
		Small = n;
	}
	
	boolean isTrue() {
		return !((Big == null || Big.intValue() == 0) && Small == 0);
	}
	
	Int flip() {
		Int result = new Int();
		if (Big == null) {
			result.Small = -Small;
			return result;
		}
		result.Big = Big.negate();
		return result;
	}
}

static class List {
	Int[] Mixed;
	byte[] Bytes;
	
	int MixedSize = 0;
	int BytesSize = 0;
	
	List() {
		Mixed = new Int[1];
		Bytes = new byte[1];
	}
	
	List(Int i) {
		Mixed = new Int[(int)i.Small];
		Bytes = new byte[(int)i.Small];
		MixedSize = (int)i.Small;
		BytesSize = (int)i.Small;
	}
	
	List(byte[] b, Int[] m) {
		Mixed = m;
		if (m != null) {
			MixedSize = m.length;
		}
		Bytes = b;
		if (b != null) {
			BytesSize = b.length;
		}
	}
	
	//Maybe optimise memory usage here?
	//Bug here.
	void put(Int i) {
		if (Mixed.length == MixedSize) {
			Int[] tmp = new Int[Mixed.length*2];
			System.arraycopy(Mixed, 0, tmp, 0, Mixed.length);
			Mixed = tmp;
		}
		Mixed[MixedSize] = i;
		MixedSize++;
		
		if (Bytes.length == BytesSize) {
			byte[] tmp = new byte[Bytes.length*2];
			System.arraycopy(Bytes, 0, tmp, 0, Bytes.length);
			Bytes = tmp;
		}
		Bytes[BytesSize] = (byte)i.Small;
		BytesSize++;
	}
	
	Int get(Int i) {
		if (Mixed == null) {
			return new Int(Bytes[(int)i.Small] & 0xff);
		}
		return Mixed[(int)i.Small];
	}
	
	void set(long i, Int v) {
		Mixed[(int)i] = v;
	}
	
	Int size() {
		if (Mixed == null) {
			return new Int(BytesSize);
		}
		return new Int(MixedSize);
	}
	
	Int pop() {
		Int result = Mixed[MixedSize];
		MixedSize--;
		BytesSize--;
		return result;
	}
}

static interface Pipe {
	public void write(Runtime r, byte[] b);
	public int read(Runtime r, byte[] b);
}

static class StdPipe implements Pipe {
	public int read(Runtime r, byte[] b) {
		try {
			return System.in.read(b);
		} catch (Exception e) {
			r.Error.Small = 1;
			r.Error.Big = null;
			return 0;
		}
	}

	public void write(Runtime r, byte[] b) {
		try {
			System.out.write(b);
		} catch (Exception e) {
			r.Error.Small = 1;
			r.Error.Big = null;
		}
	}
}

static class MethodPipe implements Pipe {
	Method method;
	MethodPipe(Method m) {
		method = m;
	}
	
	public int read(Runtime r, byte[] b) { return 0; }

	public void write(Runtime r, byte[] b) {
		try {
			method.invoke(null, r);
		} catch (Exception e) {
			r.Error.Small = 1;
			r.Error.Big = null;
		}
	}
}

static Pipe wrap(String name, Class<?> c) {
	Class[] cArg = new Class[1]; 
	cArg[0] = Runtime.class; 
	try { 
		return new MethodPipe(c.getDeclaredMethod(name, cArg)); 
	} catch (NoSuchMethodException e) { 
		throw new RuntimeException(e); 
	}
}
`}
