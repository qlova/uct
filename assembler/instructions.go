 package assembler
 
 const (
	MAIN = iota
	EXIT
	LINK
	
	CODE
	BACK
	CALL
	FORK
	
	//Global Variables
	ERROR
	GLOBAL
	CHANNEL
	
	LOOP
	DONE
	REDO
	
	IF
	OR
	NO
	
	DATA
	
	LIST
	LOAD
	MAKE

	PUT
	POP
	GET
	SET
	
	USED
	SIZE
	
	FREE
	FREE_LIST
	FREE_PIPE

	HEAP
	HEAP_LIST
	HEAP_PIPE	

	PUSH
	PUSH_LIST
	PUSH_PIPE
	
	PULL
	PULL_LIST
	PULL_PIPE
	
	NAME
	NAME_LIST
	NAME_PIPE
	
	DROP
	DROP_LIST
	DROP_PIPE
	
	COPY
	COPY_LIST
	COPY_PIPE
	
	SWAP
	SWAP_LIST
	SWAP_PIPE
	
	PIPE
	WRAP
	OPEN
	
	READ
	SEND
	STOP
	SEEK
	
	INFO
	MOVE
	
	INT
	
	ADD
	SUB
	MUL
	DIV
	MOD
	POW
	
	LESS
	MORE
	SAME
	FLIP
	
	//Instructions for targets
	NATIVE
 )

func InstructionString(instruction byte) string {
	switch instruction {
		case MAIN: return "MAIN"
		case EXIT: return "EXIT"
		case LINK: return "LINK"
		case CODE: return "CODE"
		case BACK: return "BACK"
		case CALL: return "CALL"
		case FORK: return "FORK"
		case ERROR: return "ERROR"
		case GLOBAL: return "GLOBAL"
		case CHANNEL: return "CHANNEL"
		case LOOP: return "LOOP"
		case DONE: return "DONE"
		case REDO: return "REDO"
		case IF: return "IF"
		case OR: return "OR"
		case NO: return "NO"
		case DATA: return "DATA"
		case LIST: return "LIST"
		case LOAD: return "LOAD"
		case MAKE: return "MAKE"
		case PUT: return "PUT"
		case POP: return "POP"
		case GET: return "GET"
		case SET: return "SET"
		case USED: return "USED"
		case SIZE: return "SIZE"
		case FREE: return "FREE"
			case FREE_LIST: return "FREE_LIST"
				case FREE_PIPE: return "FREE_PIPE"
		case HEAP: return "HEAP"
		case HEAP_LIST: return "HEAP_LIST"
		case HEAP_PIPE: return "HEAP_PIPE"
		case PUSH: return "PUSH"
		case PUSH_LIST: return "PUSH_LIST"
		case PUSH_PIPE: return "PUSH_PIPE"
		case PULL: return "PULL"
		case PULL_LIST: return "PULL_LIST"
		case PULL_PIPE: return "PULL_PIPE"
		case NAME: return "NAME"
		case NAME_LIST: return "NAME_LIST"
		case NAME_PIPE: return "NAME_PIPE"
		case COPY: return "COPY"
		case COPY_LIST: return "COPY_LIST"
		case COPY_PIPE: return "COPY_PIPE"
		case PIPE: return "PIPE"
		case WRAP: return "WRAP"
		case OPEN: return "OPEN"
		case READ: return "READ"
		case SEND: return "SEND"
		case STOP: return "STOP"
		case SEEK: return "SEEK"
		case INFO: return "INFO"
		case MOVE: return "MOVE"
		case INT: return "INT"
		case ADD: return "ADD"
		case SUB: return "SUB"
		case MUL: return "MUL"
		case DIV: return "DIV"
		case MOD: return "MOD"
		case POW: return "POW"
		case LESS: return "LESS"
		case MORE: return "MORE"
		case SAME: return "SAME"
		case FLIP: return "FLIP"
		
		//Native targets.
		case NATIVE: return "NATIVE"
	}
	return "UNKNOWN"
}
