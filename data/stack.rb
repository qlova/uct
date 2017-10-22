require 'securerandom'

STDIN.binmode

class Pipe
	attr_accessor :name
	attr_accessor :file
	
	def initialize(f)
		@name = ""
		@file = nil
		@connection = nil
		@info = nil
		
		@function = f
	end
	
	def exe(stack)
		@function.call(stack)
	end
	
	def close()
		if @file
			@file.close
		end
	end
end

class Stack
	attr_accessor :activearray
	attr_accessor :error
	attr_accessor :map
	attr_accessor :numbers
	attr_accessor :arrays
	attr_accessor :pipes

	def initialize()
		@numbers = []
		@arrays  = []
		@pipes   = []
		@error   = 0
	
		@activearray = []
	
		@theheap = []
		@heaproom= []
	
		@theitheap = []
		@heapitroom = []
	
		@map = Hash.new
	end
	
	def copy()
		n = Stack.new
		n.numbers = self.numbers.dup
		n.arrays  = self.arrays.dup
		n.pipes   = self.pipes.dup
		
		return n
	end
	
	def execute()
		result = []
		s = self.grab
		name = ""
		for i in 0..s.length-1
			name += s[i].chr
		end
		
		variable = `#{name}`
		
		for c in variable.split("")
			result << c.ord
		end
		self.share(result)
	end
	
	def link()
		s = self.grab
		n = self.pull
		
		name = ""
		for i in 0..s.length-1
			name += s[i].chr
		end
		
		self.map[name] = n
	end
	
	def connect()
		s = self.grab
		name = ""
		for i in 0..s.length-1
			name += s[i].chr
		end
		
		v = self.map[name]
		if v == nil
			@error = 1
			self.push 0
		else
			self.push v
		end
	end
	
	def slice()
		s = self.grab
		self.share s[self.pull..self.pull-1]
	end
	
	def array()
		@activearray = []
		return @activearray
	end
	
	def share(array)
		@arrays << array
	end
	
	def grab()
		return @arrays.pop
	end
	
	def relay(pipe)
		@pipes << pipe
	end
	
	def take()
		return @pipes.pop
	end
	
	def push(num)
		@numbers << num
	end
	
	def pull()
		return @numbers.pop
	end
	
	def put(num)
		@activearray << num
	end
	
	def pop()
		return @activearray.pop
	end
	
	def place(array)
		@activearray = array
	end
	
	def get()
		result = @activearray[self.pull%(@activearray.length)]
		if result == nil 
			return 0
		end
		return result
	end
	
	def set(num)
		@activearray[self.pull%(@activearray.length)] = num
	end
	
	def stdout()
		text = self.grab
		for i in 0..text.length-1
			print(text[i].chr)
		end
	end
	def stdin()
		length = self.pull
		
		if length == 0
			self.share(STDIN.gets)
		elsif length > 0
			self.share(STDIN.read(length))
		else
			text = []
			while true
				b = STDIN.read(1).ord
				if b == -length
					break
				end
				text << b
			end
			self.share(text)
		end
	end
	def div(a, b)
		begin
			return a / b
		rescue
			if a == 0
				return (SecureRandom.random_number 255) + 1
			end
			return 0
		end
	end
	
	
	def openit()
		filename = ""
		text = self.grab
		for i in 0..text.length-1
			filename += text[i].chr
		end
	
		file = Pipe.new(nil)
		file.name = filename
		begin
			file.file = File.open(filename, "a+")
		rescue
			if File.directory?(filename)
				self.relay file
				return
			end
			self.error = 404
			self.relay file
		end
		self.relay file
	end
	
	def inn()
		length = self.pull
		file   = self.take
		
		if length == 0
			self.share(file.file.gets)
		elsif length > 0
			self.share(file.file.read(length))
		else
			text = []
			while true
				b = file.file.read(1).ord
				if b == -length
					break
				end
				text << b
			end
			self.share(text)
		end
	end
	
	def out()
		text = self.grab
		file = self.take
	
		if file.file == nil
			if file.name[file.name.length-1] == "/"
				if File.directory?(file.name)
			
				else
					begin
						Dir.mkdir(file.name)
					rescue
						self.error = 401
					end
				end
				return
			else
				if File.file?(file.name)
			
				else
					begin
						file.file = File.open(file.name, 'w')
					rescue
						self.error = 401
					end
				end
				return
			end
		end
	
		for i in 0..text.length-1
			file.file.puts text[i].chr
		end
	end
	
	def loadit() 
		text = self.grab
		result = []
		name = ""
	
		variable = ""	
		if text[0] == 36 and text.length > 0
			for i in 1..text.length-1
				name += text[i].chr
			end
			begin
				variable = ENV[name]
			rescue
				self.share(result)
				self.error = 404
				return
			end
		else
			begin
				variable = ARGV[text[0]]
			rescue
				self.share(result)
				self.error = 404
				return
			end
		end
	
		if !variable
			self.share(result)
			self.error = 404
			return
		end
	
	
		for c in variable.split("")
			result << c.ord
		end
		self.share(result)
	end
end
