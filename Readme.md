# Universal compiler and translator

The technology consists of Universal Assembly and a Universal Compiler.
With UCT, it is possible to target other programming languages. What this means, is that you can write a compiler with UCT and then compile it into a language such as Python or Java. For a list of targets, see below:

## Philosophy

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
	* Arrays
	* Pipes

The numbers have a strange property of being unbinded.
This means they can contain any finite +/- value or zero.


## High level UCT
There is a programming language in development called "I" which follows the UCT philosophy.
You can find it at http://github.com/qlova/ilang
It compiles to UCT and I encourage you to use it.

## Experimental features
In order to be more useful as a language Universal Assembly has some experimental features.

**I/O**
This will be a officially supported feature of Universal Assembly.
I/O allows file system access and manipulating files through pipes.

**Threading**
With the fork instruction, new threads can be created, the thread communication protocol is still in design.

**Networking**
An extension of I/O, networking will enable reading and writing of various internet protocols.


## Targets

These are the languages available as a target.
The table also shows experimental feature support.
These are the official supported lanuages for development and testing.

| Language |  I/O  | Networking | Threading |
|----------|-------|------------|-----------|
|Go		   |  NO   |    NO      |    SOME    |
|Python	   |  NO   |    NO      |    NO     |


These targets are currently not enabled, they have been in the past and will be added again in a future release.

| Language |  I/O  | Networking | Threading |
|----------|-------|------------|-----------|
|Java	   |  YES  |    YES     |    YES    |
|Rust      |  SOME |    NO      |    NO     |
|Javascript|  NO   |    NO      |    YES    |
|Arduino   |  NO   |    NO      |    NO     |
|C++       |  NO   |    NO      |    NO     |
|Bash	   |  NO   |    NO      |    NO     |
|Ruby	   |  SOME |     NO     |    NO     |
|C#		   |  SOME |     NO     |    NO     |
|Lua       |  SOME |     NO     |    NO     |
