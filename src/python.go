package main

import "flag"

var PythonReserved = []string{
	"and",       "del",       "from",      "not",       "while",   
	"as",        "elif",      "global",    "or",        "with",    
	"assert",    "else",      "if",        "pass",      "yield",    
	"break",     "except",    "import",    "print",     "len",          
	"class",     "exec",      "in",        "raise", 	"open",             
	"continue",  "finally",   "is",        "return",    "bool",   
	"def",       "for",       "lambda",    "try",		"copy",

}

//This is the Java compiler for uct.
var Python bool

func init() {
	flag.BoolVar(&Python, "py", false, "Target Python")

	RegisterAssembler(PythonAssembly, &Python, "py", "#")

	for _, word := range PythonReserved {
		PythonAssembly[word] = Reserved()
	}
}

var PythonAssembly = Assemblable{
	//Special commands.
	"HEADER": Instruction{
		Data:   "#! /bin/python3\nimport stack\nimport sys\nimport threading",
		Args: 1,
	},

	"FOOTER": Instruction{},

	"FILE": Instruction{
		Data: PythonFile,
		Path: "/stack.py",
	},

	"NUMBER": is("%s", 1),
	"BIG": 	is("%s", 1),
	"SIZE":   is("len(%s)", 1),
	"STRING": is("bytes(%s, 'utf-8')", 1),
	"ERRORS":  is("stack.ERROR", 1),
	
	"LINK":  is("stack.link()"),
	"CONNECT":  is("stack.connect()"),
	"SLICE":  is("stack.slice()"),

	"SOFTWARE": Instruction{
		Data:   "stack = stack.Stack()\n",
	},
	"EXIT": Instruction{
		Data:        "sys.exit(stack.ERROR)",
	},

	"FUNCTION": is("def %s(stack):", 1, 1),
	"RETURN": Instruction{
		Indented:    1,
		Data:        "\n",
		Indent:      -1,
		Else: &Instruction{
			Data: "return",
		},
	},
	
	"SCOPE": is(`stack.relay(stack.pipe(%s))`, 1),
	
	"EXE": is("%s.exe(stack)", 1),

	"PUSH": is("stack.push(%s)", 1),
	"PULL": is("%s = stack.pull()", 1),

	"PUT":   is("stack.put(%s)", 1),
	"POP":   is("%s = stack.pop()", 1),
	"PLACE": is("stack.activearray = %s", 1),
	"ARRAY":  is("%s = stack.array()", 1),
	"RENAME": is("%s = stack.activearray", 1),
	"RELOAD": is("%s = stack.take()", 1),

	"SHARE": is("stack.share(%s)", 1),
	"GRAB":  is("%s = stack.grab()", 1),

	"RELAY": is("stack.relay(%s)", 1),
	"TAKE":  is("%s = stack.take()", 1),

	"GET": is("%s = stack.get()", 1),
	"SET": is("stack.set(%s)", 1),

	"VAR": is("%s = 0", 1),

	"OPEN":   is("stack.open()"),
	"LOAD":   is("stack.load()"),
	"OUT":    is("stack.out()"),
	"STAT":   is("stack.info()"),
	"IN":     is("stack.inn()"),
	"STDOUT": is("stack.stdout()"),
	"STDIN":  is("stack.stdin()"),
	"HEAP":   is("stack.heap()"),
	"HEAPIT":   is("stack.heapit()"),

	"CLOSE": is("%s.close()", 1),

	"LOOP":   is("while 1:", 0, 1),
	"BREAK":  is("break"),
	"REPEAT": is("\n", 0, -1, -1),

	"IF":   is("if %s != 0:", 1, 1),
	"ELSE": is("else:", 0, 0, -1),
	"END":  is("\n", 0, -1, -1),

	"RUN":  is("%s(stack)", 1),
	"DATA": is("%s = %s;", 2),

	"FORK": is("threading.Thread(target=%s, args=(stack.copy(),)).start()\n", 1),

	"ADD": is("%s = %s + %s", 3),
	"SUB": is("%s = %s - %s", 3),
	"MUL": is("%s = stack.mul(%s, %s)", 3),
	"DIV": is("%s = stack.div(%s, %s)", 3),
	"MOD": is("%s = %s %% %s", 3),
	"POW": is("%s = stack.pow(%s, %s)", 3),

	"SLT": is("%s = int(%s <  %s)", 3),
	"SEQ": is("%s = int(%s == %s)", 3),
	"SGE": is("%s = int(%s >= %s)", 3),
	"SGT": is("%s = int(%s >  %s)", 3),
	"SNE": is("%s = int(%s != %s)", 3),
	"SLE": is("%s = int(%s <= %s)", 3),

	"JOIN": is("%s = %s + %s", 3),
	"ERROR": is("stack.ERROR = %s", 1),
}

//Edit this in a Java IDE.
const PythonFile = `
#This is the python stack implementation for UCT.
import socket
import sys
import os
import time

Networks_In = {}

class Pipe:
	def __init__(self, f=None):
		self.name = ""
		self.file = None
		self.connection = None
		self.info = None
		self.function = f

	def exe(self, stack):
		self.function(stack)

	def close(self):
		if self.file:
			self.file.close()
		if self.connection:
			self.connection.close()

class Stack:
	def __init__(self):
		self.numbers = []
		self.arrays = []
		self.pipes = []
		self.ERROR = 0

		self.activearray = []
		
		self.theheap = []
		self.heaproom = []
		
		self.theitheap = []
		self.heapitroom = []
		
		self.map = {}

	def copy(self):
		n = Stack()
		n.numbers = self.numbers.copy()

		n.arrays = self.arrays.copy()

		n.pipes = self.pipes.copy()
		return n
	
	def link(self):
		s = self.grab()
		n = self.pull()
		name = ""
		for i in s:
			name += chr(i)
		self.map[name] = n
		
	def connect(self):
		s = self.grab()
		name = ""
		for i in s:
			name += chr(i)
		try:
			self.push(self.map[name])
		except:
			self.push(0)
			self.ERROR=1
	
	def slice(self):
		s = self.grab()
		self.share(s[self.pull():self.pull()])
		
	def pipe(self, f):
		return Pipe(f)

	def array(self):
		self.activearray = []
		return self.activearray

	def share(self, array):
		self.arrays.append(array)

	def grab(self):
		return self.arrays.pop()

	def relay(self, array):
		self.pipes.append(array)

	def take(self):
		return self.pipes.pop()

	def push(self, array):
		self.numbers.append(array)

	def pull(self):
		return self.numbers.pop()

	def put(self, number):
		self.activearray.append(number)

	def pop(self):
		return self.activearray.pop()

	def place(self, array):
		self.activearray = array

	def get(self):
		return self.activearray[self.pull() % len(self.activearray)]

	def set(self, number):
		self.activearray[self.pull() % len(self.activearray)] = number
		
	def heap(self):
		address = self.pull()
		
		#Get an object off the heap.
		if address > 0:
			self.share(self.theheap[(address%len(self.theheap))-1])
		
		#Delete the object.
		elif address < 0:
			self.theheap[(-address)%(len(self.theheap)+1)-1] = None
			self.heaproom.append(-address)
		
		#Add an object.
		elif address == 0:
			if len(self.heaproom) > 0:
				address = self.heaproom.pop()
				self.theheap[(address%len(self.theheap))-1] = self.grab()
				self.push(address)
			else:
				self.theheap.append(self.grab())
				self.push(len(self.theheap))
	
	def heapit(self):
		address = self.pull()
		
		#Get an object off the heap.
		if address > 0:
			self.share(self.theitheap[(address%len(self.theitheap))-1])
		
		#Delete the object.
		elif address < 0:
			self.theitheap[(-address)%(len(self.theitheap)+1)-1] = None
			self.heapitroom.append(-address)
		
		#Add an object.
		elif address == 0:
			if len(self.heapitroom) > 0:
				address = self.heapitroom.pop()
				self.theitheap[(address%len(self.theitheap))-1] = self.grab()
				self.push(address)
			else:
				self.theitheap.append(self.grab())
				self.push(len(self.theitheap))

	def load(self):
		text = self.grab()
		result = []
		name = ""

		variable = ""
		if text[0] == 36 and len(text) > 0:
			for i in range(1, len(text)):
				name += chr(text[i])
			try:
				variable = os.environ[name]
			except:
				self.share(result)
				return
		else:
			for i in text:
				name += chr(i)

			protocol = name.split("://", 2)
			if len(protocol) > 1:
				if protocol[0] == "tcp":
					try:
						address = ('localhost', int(protocol[1]))
						listener = socket.socket(socket.AF_INET, socket.SOCK_STREAM)
						listener.setsockopt(socket.SOL_SOCKET, socket.SO_REUSEADDR, 1)
						listener.bind(address)
						listener.listen(5)
						variable = str(listener.getsockname()[1])

						Networks_In[variable] = listener
					except:
						self.ERROR = 1
				elif protocol[0] == "dns":
					hosts = None
					try:
						tmp = int(protocol[1][-1])
						hosts = socket.gethostbyaddr(protocol[1])
						variable = hosts[0]
					except:
						try:
							hosts = socket.gethostbyname(protocol[1])
							variable = hosts[2].join(" ")
						except:
							self.ERROR = 1


				else:
					self.ERROR = 1
			else:
				try:
					variable = sys.argv[text[0]]
				except:
					self.ERROR = 1
					self.share(result)
					return

		for char in variable:
			result.append(ord(char))
		self.share(result)

	def open(self):
		filename = ""
		text = self.grab()
		for i in range(0, len(text)):
			filename += chr(text[i])

		file = Pipe()
		file.name = filename

		#This is for protocols such as tcp.
		protocol = filename.split("://", 2)
		if len(protocol) > 1:
			if protocol[0] == "tcp":
				try:
					listener = Networks_In[protocol[1]]
					try:
						conn = listener.accept()
						file.connection = conn[0]
						file.info = conn[1]
						self.push(0)
						self.relay(file)
						return
					except:
						raise 
						self.push(-1)
						self.relay(file)
						return
				except:
					try:
						s = socket.socket(socket.AF_INET, socket.SOCK_STREAM)
						address = protocol[1].split(":")[0]
						port = int(protocol[1].split(":")[1])
						s.connect((address, port))
						file.connection = s
						self.push(0)
					except:
						raise 
						self.push(-1)
					self.relay(file)
					return

		#This is for files.
		try:
			file.file = open(filename, "r+b")
		except:
			if os.path.isdir(filename):
				self.push(0)
			else:
				self.push(-1)
		else:
			self.push(0)
		self.relay(file)

	def info(self):
		request = ""
		variable = ""

		result = []
		text = self.grab()
		pipe = self.take()

		for i in range(0, len(text)):
			request += chr(text[i])

		if request == "ip":
			variable = pipe.info[0]

		for char in variable:
			result.append(ord(char))
		self.share(result)

	def out(self):
		pipe = self.take()

		if pipe.file is None and pipe.connection is None:
			if pipe.name[-1] == "/":
				if not os.path.isdir(pipe.name):
					try:
						os.mkdir(pipe.name)
						self.push(0)
					except:
						self.push(-1)
				else:
					self.push(0)
				return
			else:
				try:
					pipe.file = open(pipe.name, "r+b")
				except:
					self.push(-1)
					return

		text = self.grab()

		if pipe.connection:
			pipe.connection.send(bytes(text))
		else:
			for i in range(0, len(text)):
				pipe.file.write(chr(text[i]))

		self.push(0)

	def inn(self):
		length = self.pull()
		pipe = self.take()
		text = ""

		if length == 0:
			length = -ord("\n")
			
		if length > 0:
			try:
				if pipe.file:
					self.share(pipe.file.read(length))
				elif pipe.connection:
					self.share(pipe.connection.recv(length))
			except:
				self.ERROR=1
				self.share([])
			return
		else:
			bytes = []
			while 1:
				try:
					if pipe.file:
						bytes += pipe.file.read(1)
					elif pipe.connection:
						bytes += pipe.connection.recv(1)
					if bytes[-1] == -length:
						bytes = bytes[:-1]
						break
				except:
					self.ERROR=1
					break

			self.share(bytes)

	def stdout(self):
		text = self.grab()
		for i in range(0, len(text)):
			print(chr(text[i]), end="")

	def stdin(self):
		length = self.pull()
		text = ""

		if length == 0:
			text = sys.stdin.readln()
		elif length > 0:
			text = sys.stdin.read(length)
		else:
			while 1:
				text += sys.stdin.read(1)
				if ord(text[-1]) == -length:
					text = text[:-1]
					break

		self.share(bytes(text, "utf-8"))

	def div(self, a, b):
		try:
			return a // b
		except ZeroDivisionError:
			if a == 0:
				return ord(os.urandom(1)) + 1
			return 0

	def mul(self, a, b):
		return a * b

	def pow(self, a, b):
		if a == 0:
			if b % 2:
				return ord(os.urandom(1)) + 1
			return 0
		return a ** b
`
