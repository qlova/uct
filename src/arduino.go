package uct

import "flag"
import "strconv"
import "fmt"

var ArduinoReserved = []string{
	"char", "on", "off",
}

//This is the Java compiler for uct.
var Arduino bool

func init() {
	flag.BoolVar(&Arduino, "arduino", false, "Target the Arduino")

	RegisterAssembler(ArduinoAssembly, &Arduino, "ino", "//")

	for _, word := range ArduinoReserved {
		ArduinoAssembly[word] = Reserved()
	}
}

var ArduinoAssembly = Assemblable{
	//Special commands.
	"HEADER": Instruction{
		Data:   `


//System configiration.
	//Tweak these values based upon the amount of RAM your microcontroller has.	
	const int HEAP_SIZE = 48; //This is how many blocks the system can allocate before running out of memory.
	const int BLOCK_SIZE = 15; //This is the how many integers a block can hold.
//End configuriation.

enum {
	length_pos, refcount_pos, nesting_pos, data_pos
};

struct Block {
	int index[BLOCK_SIZE];
	Block *next;
};

int ERROR;
Block *activearray;

Block *arrays[5];
int array_pointer;
int numbers[5];
int number_pointer;

int (*pipes[5])();
int pipe_pointer;

//Garbage collection.
volatile int nesting;
Block *heap[HEAP_SIZE];

void collect_block(Block *b) {
	if (b->index[length_pos] > (BLOCK_SIZE-data_pos)) {
		collect_block(b->next);
	}
	b->index[length_pos] = -2;
}

void stack_collect() {
  nesting--;
  
  for (int i = 0; i < HEAP_SIZE; i++) {
		if (heap[i] && heap[i]->index[length_pos] >= 0) {
			int nest = heap[i]->index[nesting_pos];
			int refcount = heap[i]->index[refcount_pos];
			
			if (refcount == 0 && nest > nesting) {
				collect_block(heap[i]);
			}
		}
   }
}

void stack_push(int num) {
  numbers[number_pointer] = num;
  number_pointer++;
}

int stack_pull() {
  number_pointer--;
  return numbers[number_pointer]; 
}

Block *get_free_block() {

	for (int i = 0; i < HEAP_SIZE; i++) {
		if (heap[i] && heap[i]->index[length_pos] == -2) {
		  	return heap[i];
		} 
 	}

	return new Block();
}

//Create a new array.
Block *stack_array(int l) {
  //Need to check if there is a block which we can recycle!
  Block *arr = get_free_block();
  
  if (arr->index[length_pos] == -2) {
    
  	arr->index[length_pos] = l;
  	arr->index[refcount_pos] = 0;
  	arr->index[nesting_pos] = nesting;
  	activearray = arr;
  	return arr;
  	
  } else {
    
    int done = 0;
    for (int i = 0; i < HEAP_SIZE; i++) {
      if (heap[i] == 0) {
        heap[i] = arr;
        activearray = arr;
        done = 1;
        break;
      } 
    }
    
    if (done) {
    	arr->index[length_pos] = l;
      return arr;
    } else {
      Serial.println(F("[ERROR] Out of memory!"));
      exit(1);
    }
  }
}

int *stack_get_array_segment(Block *arr, int length) {
	//Linked list time!
  	//I hope this works!
  	
	Block *next;
  
  	if (!(arr->next)) {
  	  arr->next = get_free_block();
  	  next = arr->next;
  	  next->index[0] = -1;
  	} else {
  	  next = arr->next;
  	}
  	
  	length -= (BLOCK_SIZE - data_pos);
  	while (length > BLOCK_SIZE-1) {
  		if (!(next->next)) {
	  		next->next = get_free_block();
	  		next = next->next;
	  		next->index[0] = -1;
	  	} else {
	  		next = next->next;
	  	}
  		length -= (BLOCK_SIZE-1);
  	}
  	
  	return &(next->index[length+1]);
}

void stack_put(int value) {
  int length = activearray->index[length_pos];
  
  if (length >= (BLOCK_SIZE - data_pos)) {
  
  	int *pointer = stack_get_array_segment(activearray, length);
  	*pointer = value;
  	activearray->index[length_pos] = length + 1;
  	
  } else {
	  activearray->index[length+data_pos] = value;
	  activearray->index[length_pos] = length + 1;
  }
}

int stack_pop() {
  int length = activearray->index[length_pos];
	
	if (length >= (BLOCK_SIZE - data_pos)) {
	
		int value = *stack_get_array_segment(activearray, length);
		activearray->index[length_pos] = length - 1;
		return value;

	} else { 
		int value = activearray->index[length+data_pos];
		activearray->index[length_pos] = length - 1;
		return value;
	}
}

void stack_heap() {
	int mode = stack_pull();
	
	if (mode == 0) {
	  int address = (int)stack_array(0);
		stack_push(address);
	}
	if (mode > 0) {
	  
	  Block *tmp = mode;
		stack_share(tmp);
	}
	if (mode < 0) {
	 Block *tmp = mode;
		tmp->index[length_pos] = -2;
	}
}

int stack_get() {
  int index = stack_pull();
  int length = activearray->index[length_pos];
  if (length == 0) {
    return 0;
  }
  
  index = index % length;
  
  if (index < (BLOCK_SIZE - data_pos)) {
    return activearray->index[index+data_pos];
  } else {
  	return *(stack_get_array_segment(activearray, index));
  }
}

int stack_set(int value) {
  int index = stack_pull();
  int length = activearray->index[length_pos];
  index = index % length;
  
  if (index < (BLOCK_SIZE - data_pos)) {
    activearray->index[index+data_pos] = value;
  }
 int *pointer = stack_get_array_segment(activearray, index);
  *pointer = value;
}

void stack_relay(int (*pipe)()) {
  pipes[pipe_pointer] = pipe;
  pipe_pointer++;
}

int (*stack_take())() {
  pipe_pointer--;
  return pipes[pipe_pointer]; 
}

void stack_share(Block *arr) {
  arrays[array_pointer] = arr;
  arr->index[refcount_pos]++;
  array_pointer++;
}

int powint(int x, int y)
{
 int val=x;
 for(int z=0;z<=y;z++)
 {
   if(z==0)
     val=1;
   else
     val=val*x;
 }
 return val;
}

Block *stack_grab() {
  array_pointer--;
  Block *arr = arrays[array_pointer];
  arr->index[refcount_pos]--;
  arr->index[nesting_pos] = nesting;
  return arr;
}

Block *stack_join(Block *a, Block *b) {
	Block *c = stack_array(0);
	for (int i = 0; i < a->index[length_pos]; i++) {
		activearray = a;
		stack_push(i);
		int value = stack_get();
		
		activearray = c;
		stack_put(value);
	}
	for (int i = 0; i < b->index[length_pos]; i++) {
		activearray = b;
		stack_push(i);
		int value = stack_get();
		
		activearray = c;
		stack_put(value);
	}
	
	return c;
}

Block *stack_data(int data[]) {
  int length = sizeof(data) / sizeof(int);
  Block *a = stack_array(length);
  for (int i = 0; i < length; i++) {
    stack_put(data[i]);
  }
  a->index[refcount_pos] = 1;
  return a;
}

void stack_stdout() {
  nesting++;
  Block *arr = stack_grab();
  

  
  for (int i = 0; i < arr->index[length_pos]; i++) {
    if (i < (BLOCK_SIZE - data_pos)) {
      Serial.write(arr->index[i+data_pos]);
    } else {
      Serial.write(*(stack_get_array_segment(arr, i)));
    }
    Serial.flush();
  }

  stack_collect();
}

//Toaster-level stdin.
void stack_stdin() {
	
}
`,
		Args:   1,
	},

	"FOOTER": Instruction{
		Data:        "",
	},

	//Don't need a stack file.
	"FILE": Instruction{},
	
	"ARDUINO": Instruction{All:true},
	
	//"EVAL": is("Class[] cArg = new Class[1]; cArg[0] = Stack.class; try { new Object() { }.getClass().getEnclosingClass().getDeclaredMethod(stack.grab().String(), cArg).invoke(null, stack); } catch (Exception e) { throw new RuntimeException(e); }"),

	//BigInts are not supported yet.. there is a BC library ported to the Arduino, I don't know if it should be added or not, arduino's are just hardware restricted.
	"NUMBER": is("%s", 1),
	"BIG": is("%s", 1),
	
	"SIZE":   is("%s->index[length_pos]", 1),
	"ERRORS":  is("ERROR", 1),

	"SOFTWARE": Instruction{
		Data:   "void setup() { Serial.begin(9600);",
		Indent: 1,
	},
	"EXIT": Instruction{
		Indented:    1,
		Data:        "stack_collect(); } void loop(){ exit(0); }\n",
		Indent:      -1,
		Indentation: -1,
		Else: &Instruction{
			Data: "exit(ERROR);",
		},
	},

	"FUNCTION": is("void %s() { nesting++;", 1, 1),
	"RETURN": Instruction{
		Indented:    1,
		Data:        "stack_collect(); }\n",
		Indent:      -1,
		Indentation: -1,
		Else: &Instruction{
			Data: "stack_collect(); return;",
		},
	},
	
	"SCOPE": is(`stack_relay(&%s);`, 1),
	
	"EXE": is("(*%s)();", 1),

	"PUSH": is("stack_push(%s);", 1),
	"PULL": is("int %s = stack_pull();", 1),

	"PUT":   is("stack_put(%s);", 1),
	"POP":   is("int %s = stack_pop();", 1),
	"PLACE": is("activearray = %s;", 1),

	"ARRAY":  is("Block *%s = stack_array(0);", 1),
	"MAKE":  is("stack_share(stack_array(stack_pull()));"),
	"RENAME": is("%s = activearray;", 1),
	
	"RELOAD": is("%s = stack.take();", 1),

	"SHARE": is("stack_share(%s);", 1),
	"GRAB":  is("Block *%s = stack_grab();", 1),

	"RELAY": is("stack_relay(%s);", 1),
	"TAKE":  is("int (*%s)() = stack_take();", 1),

	"GET": is("int %s = stack_get();", 1),
	"SET": is("stack_set(%s);", 1),

	"VAR": is("int %s = 0;", 1),

	"OPEN":   is(""),
	"DELETE":   is(""),
	"EXECUTE":   is(""),
	"LOAD":   is(""),
	"OUT":    is(""),
	"STAT":   is(""),
	"IN":     is(""),
	"STDOUT": is("stack_stdout();"),
	"STDIN":  is("stack_stdin();"),
	"HEAP":   is("stack_heap();"),
	"LINK":   is(""),
	"CONNECT":   is(""),
	"SLICE":   is(""),

	"CLOSE": is("", 1),

	"LOOP":   is("while (1) { nesting++;", 0, 1),
	"BREAK":  is("stack_collect(); break;"),
	"REPEAT": is("stack_collect(); }", 0, -1, -1),

	"IF":   is("if (%s) { nesting++;", 1, 1),
	"ELSE": is("stack_collect(); } else { nesting++;", 0, 0, -1),
	"END":  is("stack_collect(); }", 0, -1, -1),

	"RUN":  is("%s();", 1),
	
	"DATA":  is("Block %s[]=%s;", 2),
	"STRING": Instruction{
		Args: 1,
		Data: " ",
		Function: func(args []string) string {
			str, _ := strconv.Unquote(args[0])
			
			var result string = "{"+fmt.Sprint(len(str))+", 0, 0,"
			for _, char := range str {
				result += fmt.Sprint(int(char))+", "
			}
			result = result[:len(result)-2]
			result += "}"
			return result
		},
	},

	"FORK": is("{ Stack s = stack.copy(); Stack.ThreadPool.execute(() -> %s(s)); }\n", 1),

	"ADD": is("%s = %s + %s;", 3),
	"SUB": is("%s = %s - %s;", 3),
	"MUL": is("%s = %s * %s;", 3),
	"DIV": is("%s = %s / %s;", 3),
	"MOD": is("%s = %s %% %s;", 3),
	"POW": is("%s = powint(%s, %s);", 3),

	"SLT": is("%s = (%s < %s) ? 1 : 0;", 3),
	"SEQ": is("%s = (%s == %s) ? 1 : 0;", 3),
	"SGE": is("%s = (%s >= %s)  ? 1 : 0;", 3),
	"SGT": is("%s = (%s > %s) ? 1 : 0;", 3),
	"SNE": is("%s = (%s != %s) ? 1 : 0;", 3),
	"SLE": is("%s = (%s <= %s) ? 1 : 0;", 3),

	"JOIN": is("%s = stack_join(%s, %s);", 3),
	"ERROR": is("ERROR = %s;", 1),
}
