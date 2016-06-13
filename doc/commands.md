# UCT Commands

This document will go through the available commands in uct assembly.

ROUTINE
=======
	The routine command is needed in all uct software. It is the entrypoint of your application.
	
	example:
		
		#Basic application.
		ROUTINE
			
		END

SUBROUTINE
==========
	The subroutine command allows you to define subroutines.
	
	example:
		
		SUBROUTINE myfunction
			#Do something here.
		END

FUNC
====
	The func command creates a new scoped function.
	
	example:
		
		SUBROUTINE anotherfunction
		
		END
		
		ROUTINE
			FUNC anotherfunction
			#you can now use anotherfunction as if it were a variable.
		END

EXE
====
	exe will run a scoped function.
	
	example:
		
		SUBROUTINE needstorun
		
		END
		
		ROUTINE
			FUNC needstorun
			#This will run the scoped function.
			EXE needstorun
		END

PUSH, PUSHSTRING, PUSHFUNC, PUSHIT
==================================
	The push functions push the given scoped variable onto the respective stack.
	The default push can push onto strings too.
	
	example:
		ROUTINE
			STRING mystring
			PUSH 97 mystring
			PUSH 98 mystring
			PUSHSTRING mystring
			STDOUT #Will print "ab"
		END

POP, POPSTRING POPFUNC, POPIT
=============================
	The pop functions pop a value off the stack and bring the values into scope.
	You do no declare a value which you pop.
	The plain pop function can be used to pop off strings.
	
	example:
		ROUTINE
			PUSH 20
			POP a
			#a now has the value '20'
			
			STRING anotherstring
			PUSH 10 anotherstring
			
			POP b anotherstring
			#b now has the value 10
		END
		
INDEX
=====
	The index command will index a string.

	example:
	
		ROUTINE
			STRING items
			PUSH 22 items
			PUSH 11 items
			INDEX items 0 a
			#a now holds the value 22
		END

SET
===
	Set can modify a string by index.
	
	example:
		
		ROUTINE
			STRING ages
			PUSH 35 ages
			PUSH 16 ages
			PUSH 24 ages
			
			SET ages 1 18
			
			#String will now be: [35, 18, 24]
		END

VAR
===
	Initialise a value.
	
	example:
		
		ROUTINE
			VAR a
			VAR b 10
			
			#a=0 and b=10
		END

STRING
======
	Initialise a string.
	
	example:
		
		ROUTINE
			STRING A
		END

OPEN
====
	Opens inputs and outputs.
	Must check status.
	
	example:
	
		STRINGDATA input "file.name"
		
		ROUTINE
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
		
		ROUTINE
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
		
		ROUTINE
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
		
		ROUTINE
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
		
		ROUTINE
			PUSHSTRING home
			LOAD
			STDOUT
		END
			
STDOUT
======
	Sends string on the stack to stdout.
	
	example:
		
		STRINGDATA helloworld "Hello World"
		
		ROUTINE
			PUSHSTRING helloworld
			STDOUT
		END

STDIN
=====
	Gets values from stdin.
	Must push to stack how many values.
	Blocks on recieve.
	
	example:
	
		ROUTINE
			PUSH 2
			STDIN 
			
			POP a
			POP b
		END

LOOP, REPEAT, BREAK
====
	Not so Ininite loop.
	
	example:
		
		ROUTINE
			LOOP
				BREAK
			REPEAT
		END

IF
==
	If will branch to end if equal to 0
	
	example:
		
		ROUTINE
			VAR a
			VAR b 1
			IF a
				#This will not run.
			ELSEIF b
				#This will run.
			ELSE
				#This will not run.
			END
		END

RUN
===
	Jump to a subroutine.
	
	example:
		
		SUBROUTINE next
			#This runs!
		END
		
		ROUTINE
			RUN next
		END

RETURN
======
	Return from current routine or subroutine.
	
	example:
		
		ROUTINE
			RETURN
		END

STRINGDATA
==========
	Literal string. Must be outside routine.
	
	example:
		STRINGDATA text "How are you?"
		
		ROUTINE
			PUSHSTRING text
			STDOUT
		END

JOIN
====
	Join two strings together.
	
	example:
		STRINGDATA a "First "
		STRINGDATA b "Last"
		
		ROUTINE
			JOIN c a b
			PUSHSTRING c
			STDOUT # Will output "First Last"
		END

ADD, SUB, MUL, DIV, MOD, POW
============================
	Maths...
	
	example:
		
		ROUTINE
			VAR a = 22
			VAR b = 2
			VAR c
			
			ADD c a b
		END

SLT, SEQ, SGE, SGT, SNE, SLE
============================
	Conditions! (Set less then, set equal... etc)
	
	example:
		
		ROUTINE
			VAR a 50
			VAR b 60
			VAR c
			
			SGE c a b
			
			# c will still be 0
			
			SNE c a b
			
			# c will be 1
		END
