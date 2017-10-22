# UCT Commands

This document will go through the available commands in uct assembly.

SOFTWARE/EXIT
=======
	This is the entrypoint of your application. The close command will exit the software.
	SOFTWARE must be followed at some point by EXIT.
	
	Your writing software, so all software needs a beginning and an end.
	The Software is a beginning and the Exit is an end. An Exit is where the software will close.
	
	example:
		
		#Basic application.
		SOFTWARE
			
		EXIT

FUNCTION
==========
	The function command allows you to define functions.
	FUNCTION must be followed at some point by RETURN.
	
	A Function is a role or purpose of your software, with functions you can organise your software into different parts.
	
	example:
		
		FUNCTION name
			#Do something here.
		RETURN

ERROR
=======
	Error is a global value for a thread that is set to a value when an instruction has failed.
	It is mainly useful for checking whether or not IO operations succeed. (OPEN and LOAD)
	
	example:
		
		SOFTWARE
			ARRAY a
			SHARE a
			OPEN
			
			IF ERROR
				#open failed!
			END
		EXIT
		
HEAP
======
	The heap is a heap of arrays.
	
	example:
	
		DATA array "12345"
	
		SOFTWARE
		
			SHARE array
			PUSH 0
			HEAP
			POP address
			
			MUL address address -1
			PUSH address
			HEAP
			GRAB samearray
		EXIT
		
		
PUT/PLACE/POP
===========
	These allow you to add values to arrays.
	
	example:
		
		SOFTWARE
			ARRAY a
				PUT 1
			ARRAY b
				PUT 1
			PLACE a
				PUT 2
			PLACE b
				PUT 2
				POP variable
				
			# a is now [1, 2]
			# b is now [1]
			# variable is now 2
		EXIT

PUSH/PULL
==========
	These two commands can be used to push values to functions.
	
	example:
		FUNCTION echo
			PULL number
			
			ARRAY text
				PUT number
			
			PASS text
			STDOUT
		RETURN
		
		#This will print 'a'
		SOFTWARE
			PUSH 97
			RUN echo
		EXIT

SHARE/GRAB
===========
	These can be used to pass and take arrays.
	
	example:
		FUNCTION echo
			GRAB stack
			SHARE stack
			STDOUT
		RETURN
	
		#This will print 'a'
		SOFTWARE
		
SCOPE/RUN
=======
	Transforms a function defined with the function instruction into a pipe.
	
	example:
		FUNCTION f
			
		RETURN
		
		SOFTWARE
			SCOPE f
			TAKE pipe
			
			#This runs f()
			RUN pipe
		EXIT
			

RELAY/TAKE
============
	These are used for passing I/O's
	
	example:
		FUNCTION message
			TAKE file
			GRAB text
			
			
			OUT
		RETURN
		
		SOFTWARE
		
			OPEN file
			STACK name
				PUT 97
			SHARE name
			
			RUN echo
		EXIT
	
		
		
GET
=====
	The get command will index an array.

	example:
	
		SOFTWARE
			ARRAY items
				PUT 22
				PUT 11
				PUSH 0
				GET a
			#a now holds the value 22
		EXIT
		
SET
=====
	The set command will modify an array.

	example:
	
		SOFTWARE
			ARRAY ages
				PUT 35
				PUT 16
				PUT 24
				
				PUSH 1
				SET 18
			#String will now be: [35, 18, 24]
		EXIT


VAR
===
	Initialise a value.
	
	example:
		
		SOFTWARE
			VAR a
			VAR b
			
			#a=0 and b=10
		EXIT

ARRAY/RENAME
===============
	Initialise an array.
	
	example:
		
		SOFTWARE
			ARRAY A
			ARRAY B
			SHARE B
			RENAME A
			# A = B
		END

OPEN
====
	Opens inputs and outputs.
	Must check status.
	
	example:
	
		STRINGDATA input "file.name"
		
		SOFTWARE
			PUSHSTRING input
			OPEN it
			POP failed
			
			IF failed
				ERROR 1
				RETURN
			END
			
			#it contains the input/output scope.
		END

OUT
===
	Sends a value to an output.
	Must pop a status after.

	example:
		
		STRINGDATA input "file.name"
		STRINGDATA data "1234567890"
		
		SOFTWARE
			PUSHSTRING input
			OPEN it
			POP failed
			
			IF failed
				ERROR 1
				RETURN
			END
			
			PUSHSTRING data
			OUT it
			
			POP failed2
			IF failed2
				ERROR 1
				RETURN
			END
		END

IN
==
	Read N many values from input.
	Need to check for early return. (-1000 value)
	
	example:
	
		STRINGDATA input "file.name"
		
		SOFTWARE
			PUSHSTRING input
			OPEN it
			POP failed
			
			IF failed
				ERROR 1
				RETURN
			END
			
			PUSH 2 #Read 2 values.
			IN it
			
			POP a
			SNE condition a -1000
			IF condition
				POP b
			END
		END

CLOSE
=====
	Close an input/output value.
	
	example:
		
		STRINGDATA input "file.name"
		
		SOFTWARE
			PUSHSTRING input
			OPEN it
			POP failed
			
			IF failed
				ERROR 1
				RETURN
			END
			CLOSE it
		END

LOAD
====
	Loads a string from the system.
	Also loads command line arguments if the array passed to it consists of one element (the nth command line argument).
	
	example:
		
		DATA home "$HOME"
		
		SOFTWARE
			SHARE home
			LOAD
			STDOUT
		END
			
STDOUT
======
	Sends string on the stack to stdout.
	
	example:
		
		STRINGDATA helloworld "Hello World"
		
		SOFTWARE
			PUSHSTRING helloworld
			STDOUT
		END

STDIN
=====
	Gets values from stdin.
	Must push to stack how many values.
	Blocks on recieve.
	
	example:
	
		SOFTWARE
			PUSH 2
			STDIN 
			
			POP a
			POP b
		END

LOOP, REPEAT, BREAK
====
	Not so Ininite loop.
	
	example:
		
		SOFTWARE
			LOOP
				BREAK
			REPEAT
		END

IF, ELSE
==
	If will branch to else or end if equal to 0
	
	example:
		
		SOFTWARE
			VAR a
			IF a
				#This will not run.
			ELSE
				#This will not run.
			END
		END

RUN
===
	Jump to a function.
	
	example:
		
		FUNCTION next
			#This runs!
		RETURN
		
		SOFTWARE
			RUN next
		EXIT

RETURN
======
	Return from current SOFTWARE or subSOFTWARE.
	
	example:
		
		SOFTWARE
			RETURN
		EXIT

DATA
==========
	Literal string. Must be outside SOFTWARE.
	
	example:
		DATA question "How are you?"
		
		SOFTWARE
			SHARE question
			STDOUT
		EXIT

JOIN
====
	Join two strings together.
	
	example:
		DATA a "First "
		DATA b "Last"
		
		SOFTWARE
			JOIN c a b
			PUSHSTRING c
			STDOUT # Will output "First Last"
		EXIT

FORK, INBOX, OUTBOX
====
	Threading instructions.
	Fork Executes a function with a copy of the stack, that function will run independantly.|
	Inbox blocks and recieves messages from threads, outbox is used within a thread and sends messages to the parent thread.
	
	example:
		FUNCTION myfunction
			#Do something here.
			ARRAY result
				PUT '!'
			SHARE result
			OUTBOX
		RETURN
		
		SOFTWARE
			FORK myfunction

			INBOX
			STDOUT
			#Prints '!'
		EXIT


ADD, SUB, MUL, DIV, MOD, POW
============================
	Maths...
	
	example:
		
		SOFTWARE
			VAR a
			ADD a 0 22
			
			VAR b
			ADD b 0 2
			
			ADD c a b
		EXIT

SLT, SEQ, SGE, SGT, SNE, SLE
============================
	Conditions! (Set less then, set equal... etc)
	
	example:
		
		SOFTWARE
			VAR a 50
			VAR b 60
			VAR c
			
			SGE c a b
			
			# c will still be 0
			
			SNE c a b
			
			# c will be 1
		EXIT
