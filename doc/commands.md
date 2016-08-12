# UCT Commands

This document will go through the available commands in uct assembly.

SOFTWARE/EXIT
=======
	This is the entrypoint of your application. The close command will exit the software.
	SOFTWARE must be followed at some point by EXIT
	
	example:
		
		#Basic application.
		SOFTWARE
			
		EXIT

FUNCTION
==========
	The function command allows you to define functions.
	FUNCTION must be followed at some point by RETURN.
	
	example:
		
		FUNCTION name
			#Do something here.
		RETURN
		
HEAP
======
	The heap is a heap of arrays.
	
	example:
		
		
PUT/PLACE/POP
===========
	These allow you to add values to stacks.
	
	example:
		
		SOFTWARE
			STACK a
				PUT 1
			STACK b
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
			
			STACK text
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
	These can be used to pass and take stacks, stacks are threadsafe.
	
	example:
		FUNCTION echo
			GRAB stack
			SHARE stack
			STDOUT
		RETURN
	
		#This will print 'a'
		SOFTWARE
			

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

# WIP

RING/LINK/NEXT/SELECT
=======================
	This is a data structure for IO devices.
	Select will select a device which is ready to recieve. (Blocking)
	The call to next after this will give you that file.
	
	
	example:
	
		SOFTWARE
			RING files
				LINK file
				SELECT file
				
			RELAY file
			IN
			
			
			NEXT file
			
		EXIT
	
		
		
GET
=====
	The get command will index a stack.

	example:
	
		SOFTWARE
			STACK items
				PUT 22
				PUT 11
				PUSH 0
				GET a
			#a now holds the value 22
		EXIT
		
SET
=====
	The set command will modify a string.

	example:
	
		SOFTWARE
			STACK ages
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

STACK/RESTACK
===============
	Initialise a stack.
	
	example:
		
		SOFTWARE
			STACK A
			STACK B
			SHARE B
			RESTACK A
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
	
	example:
		
		STRINGDATA home "$HOME"
		
		SOFTWARE
			PUSHSTRING home
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

FORK
====
	Executes a function with a copy of the stack, that function will run independantly.
	
	example:
		FUNCTION myfunction
			#Do something here.
		RETURN
		
		SOFTWARE
			FORK myfunction
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
