import socket
import sys
import os
import time
import queue
import multiprocessing

#This is a table mapping port strings to serversockets.
Networks_In = {}

#This is the UCT Pipe implementation.
class Pipe:
	def __init__(self, f=None):
		self.name 		= ""
		self.file 		= None
		self.connection = None
		self.info 		= None
		self.function 	= f
		self.q = None
		self.size = 0


	def exe(self, stack):
		self.function(stack)

	def close(self):
		if self.file:
			self.file.close()
		if self.connection:
			self.connection.close()

#This is the UCT stack implementation.
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
		
		self.inbox = None
		self.outbox = None

	def copy(self):
		n = Stack()
		n.numbers = self.numbers.copy()

		n.arrays = self.arrays.copy()

		n.pipes = self.pipes.copy()
		
		parent_conn, child_conn = multiprocessing.Pipe()
		n.outbox = child_conn
		self.inbox = parent_conn
		
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
	
	def queue(self):
		p = Pipe()
		p.q = multiprocessing.Queue()
		return p

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
		if len(self.activearray) == 0:
			self.ERROR = 4
			return 0
		return self.activearray[self.mod(self.pull(), len(self.activearray))]

	def set(self, number):
		if len(self.activearray) == 0:
			self.ERROR = 4
			return 0
		self.activearray[self.mod(self.pull(), len(self.activearray))] = number
		
	def heap(self):
		address = self.pull()
		
		#Get an object off the heap.
		if address > 0:
			self.share(self.theheap[(address%len(self.theheap))-1])
		
		#Delete the object.
		elif address < 0:
			self.theheap[(-address)%(len(self.theheap)+1)-1] = []
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
			self.relay(self.theitheap[(address%len(self.theitheap))-1])
		
		#Delete the object.
		elif address < 0:
			self.theitheap[(-address)%(len(self.theitheap)+1)-1] = None
			self.heapitroom.append(-address)
		
		#Add an object.
		elif address == 0:
			if len(self.heapitroom) > 0:
				address = self.heapitroom.pop()
				self.theitheap[(address%len(self.theitheap))-1] = self.take()
				self.push(address)
			else:
				self.theitheap.append(self.take())
				self.push(len(self.theitheap))

	def load(self):
		text = self.grab()
		result = []
		name = ""

		variable = ""
		if len(text) > 0 and text[0] == 36:
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
	
	def execute(self):
		result = []
		command = ""
		for i in range(0, len(command)):
			command += chr(text[i])
			
		variable = os.popen(command).read()
		
		for char in variable:
			result.append(ord(char))
		self.share(result)

	def delete(self):
		filename = ""
		text = self.grab()
		for i in range(0, len(text)):
			filename += chr(text[i])
		
		if filename == "":
			self.ERROR = 1
			return
		
		try:
			os.remove(filename)
		except:
			self.ERROR = 1
			return
		

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
						self.relay(file)
						return
					except:
						self.ERROR = 404
						self.relay(file)
						return
				except:
					try:
						s = socket.socket(socket.AF_INET, socket.SOCK_STREAM)
						address = protocol[1].split(":")[0]
						port = int(protocol[1].split(":")[1])
						s.connect((address, port))
						file.connection = s
					except:
						self.ERROR = 404
						
					self.relay(file)
					return

		#This is for files.
		try:
			file.file = open(filename, "r+b")
			file.size = os.path.getsize(filename)
		except:
			try:
				file.file = open(filename, "rb")
				file.size = os.path.getsize(filename)
			except:
				if not os.path.isdir(filename):
					self.ERROR = 404
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
		else:
			result.append(pipe.size)
			self.share(result)
			return

		for char in variable:
			result.append(ord(char))
		self.share(result)

	def out(self):
		pipe = self.take()
		
		if pipe.q:
			pipe.q.put(self.grab())
			return

		if pipe.file is None and pipe.connection is None:
			if len(pipe.name) > 0 and pipe.name[-1] == "/":
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
					pipe.file = open(pipe.name, "w+b")
				except:
					self.push(-1)
					return

		text = self.grab()
		
		try:
			if pipe.connection:
				pipe.connection.send(bytes(text))
			else:
				pipe.file.write(bytes(text))
		except:
			self.ERROR = 1
			self.push(-1)
			return

		self.push(0)

	def inn(self):
		length = self.pull()
		pipe = self.take()
		text = ""
		
		if pipe.q:
			self.share(pipe.q.get())
			return

		if length == 0:
			length = -ord("\n")
			
		if length > 0:
			bytes = []
			try:
				if pipe.file:
					bytes=pipe.file.read(length)
					if len(b) < length:
						self.ERROR=1
					self.share(list(bytes))
				elif pipe.connection:
					self.share(list(pipe.connection.recv(length)))
				else:
					self.ERROR=1
					self.share([])
			except:
				self.ERROR=1
				self.share(list(bytes))
			return
		else:
			bytes = []
			while 1:
				try:
					if pipe.file:
						b = pipe.file.read(1)
						if len(b) == 0:
							self.ERROR=1
							break
						bytes += b
					elif pipe.connection:
						bytes += pipe.connection.recv(1)
					if bytes[-1] == -length:
						bytes = bytes[:-1]
						break
				except:
					self.ERROR=1
					break

			self.share(list(bytes))

	def stdout(self):
		text = self.grab()
		for i in range(0, len(text)):
			sys.stdout.buffer.write(bytes([text[i]]))

	def stdin(self):
		length = self.pull()
		text = ""

		if length == 0:
			length = -10
		
		if length > 0:
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
	
	def mod(self, a, b):
		try:
			return a % b
		except ZeroDivisionError:
			return 0

	def mul(self, a, b):
		return a * b

	def pow(self, a, b):
		return a ** b
