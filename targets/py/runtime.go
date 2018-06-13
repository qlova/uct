package python

import uct "github.com/qlova/uct/assembler"

var Runtime = uct.Runtime{
	Name: "runtime.py",
	Data: `#!/usr/bin/env python3
import os
import sys
import collections

class Empty:
	name = ""
	file = None
	
	def __init__(self, name):
		self.name = name
		
	def write(self, data):
		if self.file == None:
			if self.name[-1] == "/":
				os.makedirs(self.name)
			else:
				self.file = open(self.name, "ba+")

		return self.file.write(data)
		

class Std:
	def write(self, bytes):
		sys.stdout.buffer.write(bytes)
		return True
	def read(self, size):
		sys.stdout.buffer.flush()
		return sys.stdin.buffer.read(size)
	
class WrappedFunction:
	function = None

	def __init__(self, f):
		self.function = f
		
	def write(self, bytes):
		return False
	
class Runtime:
	def __init__(self):
		self.Error = 0
		self.Global = []
		self.Channel = None
		
		self.Stack = collections.deque()
		self.Lists = collections.deque()
		self.Pipes = collections.deque()
		
		self.TheHeap = []
		self.TheHeapRoom = collections.deque()

		self.TheListHeap = []
		self.TheListHeapRoom = collections.deque()
		
		self.ThePipeHeap = []
		self.ThePipeHeapRoom = collections.deque()
	
	def div(self):
		b = self.Stack.pop()
		a = self.Stack.pop()
		
		if b == 0:
			if a == 0:
				runtime.Stack.append(1)
				return
			else:
				runtime.Stack.append(0)
				return
				
		d, r = divmod(a, b)
			
		if r != 0 and b < 0:
			d += 1
		runtime.Stack.append(d)
	
	def div(self):
		b = self.Stack.pop()
		a = self.Stack.pop()
		
		if b == 0:
			if a == 0:
				runtime.Stack.append(1)
				return
			else:
				runtime.Stack.append(0)
				return
		
		runtime.Stack.append(a//b)
	
	def mod(self):
		b = self.Stack.pop()
		a = self.Stack.pop()
		runtime.Stack.append(a%b)
		
	def pow(self):
		b = self.Stack.pop()
		a = self.Stack.pop()
		runtime.Stack.append(a**b)

	def send(self):
		pipe = self.Pipes.pop()
		bytes = self.Lists.pop()

		try:
			if not pipe.write(bytearray(bytes)):
				if hasattr(pipe, "function"):
					pipe.function(runtime)
					return
				self.Error = 1
		except PermissionError:
			self.Error = 13
		except:
			self.Error = 1
			
	def read(self):
		size = self.Stack.pop()
		pipe = self.Pipes.pop()
		
		if size == 0:
			raise ValueError('Read 0 unimplemented!')
		elif size > 0:
			self.Lists.append(pipe.read(size))
		else:
			text = []
			try:
				while 1:
					text += pipe.read(1)
					if text[-1] == -size:
						text = text[:-1]
						break
			except:
				self.Error = 1
			self.Lists.append(text)
	
	def open(self):
		uri = bytearray(self.Lists.pop()).decode("utf8")
		if len(uri) == 0:
			self.Pipes.append(Std())
			return
		
		try:
			if os.path.isfile(uri) or os.path.isdir(uri):
				self.Pipes.append(Empty(uri))
			else:
				self.Error = 2
				self.Pipes.append(Empty(uri))
		except:
			#Stupid unicode errors..
			self.Error = 1
			self.Pipes.append(Empty(uri))
	
	def load(self):
		name = ""
		variable = ""
		
		uri = self.Lists.pop()
		
		if len(uri) > 1 and uri[0] == ord("$"):
			uri.pop(0)
			name = uri.decode("utf8")
			variable = os.environ[name]
		elif len(uri) > 0:
			try:
				file = open(uri.decode("utf8"), "rb")
				self.Lists.append(file.read())
				file.close()
				return
			except:
				self.Error = 1
		
		self.Lists.append(bytearray(variable, "utf8"))
		
	def heaplist(self):
		address = self.Stack.pop()
		
		#Get an object off the heap.
		if address > 0:
			self.Lists.append(self.TheListHeap[address-1])
		
		#Delete the object.
		elif address < 0:
			address = -address
			self.TheListHeap[address-1] = []
			self.TheListHeapRoom.append(address-1)
		
		#Add an object.
		elif address == 0:
			if len(self.TheListHeapRoom) > 0:
				address = self.TheListHeapRoom.pop()
				self.TheListHeap[address] = self.Lists.pop()
				self.Stack.append(address-1)
			else:
				self.TheListHeap.append(self.Lists.pop())
				self.Stack.append(len(self.TheListHeap))
`}
