#Universal compiler and translator

UCT is an innovative new technology for computer science.

The technology consists of Universal Assembly and a Universal Compiler.
With UCT, it is possible to target other programming languages. What this means, is that you can write a program in Universal Assembly and then compile it into a language such as Python or Java. For a list of targets, see below:

##Philosophy/Rant (Trigger warning!)
Most Programming languages are stupid.  
We as a species are too attached to the limitations of computer hardware. For some reason programming is still chained to concepts such as bits and bytes and obscure things which have little meaning to anything else in the world. If I want to program, I don't care about what hardware I have! I just want my instructions to be followed.  
Java is not a solution to this at all. Don't program in Java.  
Software should be written in a more mathmatical manner.. not with strange symbols I might add! But with the abstraction and elegance.
You should not write any programs in UCT but you should accept it as the future assembly language which will empower the people of humanity.  
There is a core principle of Universal Assembly:

**Hardware Agnostic**
Universal Assembly does not care about the hardware it runs on, it could be running on a fridge, a toaster.. nearly anything in your kitchen.  
I have prepared a few words you should erase out of your mind when thinking about Universal Assembly:
	
	* bit
	* byte
	* int
	* float32
	* bool
	* RAM
	* CPU
	* etc...

You may be thinking how in the world you could program without such terms?  
Well you can. Universal Assembly assumes you are running on a perfect calculator.
It may not support quantum computers yet but it can still suit your non-quantum needs.
The only data types in Universal Assembly are:

	* Numbers
	* Strings (of Numbers)
	* Input & Output
	* Functions

The numbers have a strange property of being unbinded.
This means they can contain any finite +/- value or zero.


##High level UCT
There is a programming language in development called "I" which follows the UCT philosophy.
You can find it at http://github.com/qlova/ilang
It compiles to UCT and I encourage you to use it.

##Experimental features
In order to be more useful as a language Universal Assembly has some experimental features.

**I/O**
This will be a officially supported feature of Universal Assembly.
I/O allows file system access and manipulating file registers.

**Threading**
With the fork instruction, new threads can be created, the thread communication protocol is still in design.

**Networking**
An extension of I/O, networking will enable reading and writing of various internet protocols.


##Targets

These are the languages available as a target.
The table also shows experimental feature support.

| Language |  I/O  | Networking | Threading |
|----------|-------|------------|-----------|
|Go		   |  YES  |    YES     |    YES    |
|Bash	   |  SOME |     NO     |    NO     |
|Python	   |  YES  |     NO     |    NO     |
|Java	   |  SOME |    YES     |    YES    |
|Ruby	   |  SOME |     NO     |    NO     |
|C#		   |  SOME |     NO     |    NO     |
|Lua       |  SOME |     NO     |    NO     |

#Install

```bash
#Make sure $GOPATH/bin is in $PATH then:
go get github.com/qlova/uct
```

#Using

```bash
		uct -ext input.u
		#eg for go:
		uct -go input.u

		#you will then find a file (input.go)
		#Which you can compile with your favourite go compiler.
```


#Hello World
```u
		STRINGDATA helloworld "Hello World\n"

		# Prints "Hello World"
		ROUTINE
			PUSHSTRING helloworld
			STDOUT
		END
```

