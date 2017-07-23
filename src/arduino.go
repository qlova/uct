package uct

import "flag"
import "strconv"
import "fmt"

var ArduinoReserved = []string{
	"char",
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


struct ExtendedArray { 
    int data[14];
    ExtendedArray *next;
};

struct Array { 
    int length;
    int refcount; //reference counter.
    int nesting;
    int data[11];
    ExtendedArray *next;
};


int ERROR;
Array *activearray;

Array *arrays[5];
int array_pointer;
int numbers[5];
int number_pointer;

int (*pipes[5])();
int pipe_pointer;

//Garbage collection.
int nesting;
Array *garbagetruck[48];
int garbage_pointer;

void extended_array_collect(ExtendedArray *arr) {
	if (!arr->next) {
		extended_array_collect(arr->next);
	}
	free(arr);
}

void array_collect(Array *arr) {
	if (!arr->next) {
		extended_array_collect(arr->next);
	}
	free(arr);
}

void stack_collect() {
  nesting--;
  for (int i = 0; i < 64; i++) {
    if (!garbagetruck[i]) {
      //Check the reference counter and the nesting value.
      //If the reference counter is zero and the nesting value is greater than the current nesting value.
      //Free the memory.
      Array *arr = garbagetruck[i];
      if (arr->refcount == 0 && arr->nesting > nesting) {
        array_collect(arr);
        garbagetruck[i] = 0;
      }
    } 
  }
}

Array *stack_data(int data[]) {
  int length = sizeof(data) / sizeof(int);
  Array *a = stack_array(length);
  for (int i = 0; i < length; i++) {
    stack_put(data[i]);
  }
  a->refcount = 1;
  return a;
}
	
Array *stack_array(int l) {
  Array *arr = new Array();
  
  //Because ardunio is buggy.
  //Don't remove these next two lines because, well, It breaks stdin.
  //(I have no clue)
  Serial.print(F(""));
  Serial.flush();
  
  arr->length = 0; //Length
  arr->refcount = 0; //Reference Count.
  arr->nesting = nesting; //Nesting value.
  activearray = arr;
  for (int i = 0; i < 64; i++) {
    if (!garbagetruck[i]) {
      garbagetruck[i] = arr;
    } 
  }
  return arr;
}

int *stack_get_array_segment(Array *arr, int length) {
	//Linked list time!
  	//I hope this works!
  	
   ExtendedArray *next;
  
  	if (!(arr->next)) {
  	  arr->next = new ExtendedArray();
  	  next = arr->next;
  	} else {
  	  next = arr->next;
  	}
  	
  	length -= 11;
  	while (length > 14) {
  		if (!(next->next)) {
	  		next->next = new ExtendedArray();
	  	}
	  	next = next->next;
  		length -= 14;
  	}
  	
  	return &(next->data[length]);
}

void stack_put(int value) {
  int length = activearray->length;
  
  if (length >= 11) {
  
  	int *pointer = stack_get_array_segment(activearray, length);
  	*pointer = value;
  	activearray->length = length + 1;

  	
  } else {
	  activearray->data[length] = value;
	  activearray->length = length + 1;
  }
}

Array *stack_join(Array *a, Array *b) {
	Array *c = stack_array(0);
	for (int i = 0; i < a->length; i++) {
		activearray = a;
		stack_push(i);
		int value = stack_get();
		
		activearray = c;
		stack_put(value);
	}
	for (int i = 0; i < b->length; i++) {
		activearray = b;
		stack_push(i);
		int value = stack_get();
		
		activearray = c;
		stack_put(value);
	}
	
	return c;
}

int  stack_pop() {
  int length = activearray->length;
	
	if (length >= 11) {
	
		int value = *stack_get_array_segment(activearray, length);
		activearray->length = length - 1;
		return value;

	} else { 
	  int value = activearray->data[length];
	  activearray->length = length - 1;
	  return value;
	}
}


int stack_get() {
  int index = stack_pull();
  int length = activearray->length;
  if (length == 0) {
    return 0;
  }
  
  index = index % length;
  
  if (index < 11) {
    return activearray->data[index];
  } else {
  	int *value = stack_get_array_segment(activearray, index);
  	return *value;
  }
}

int stack_set(int value) {
  int index = stack_pull();
  int length = activearray->length;
  index = index % length;
  
  if (index < 11) {
    activearray->data[index] = value;
  }
 int *pointer = stack_get_array_segment(activearray, index);
  *pointer = value;
}

void stack_push(int num) {
  numbers[number_pointer] = num;
  number_pointer++;
}

int stack_pull() {
  number_pointer--;
  return numbers[number_pointer]; 
}

void stack_relay(int (*pipe)()) {
  pipes[pipe_pointer] = pipe;
  pipe_pointer++;
}

int (*stack_take())() {
  pipe_pointer--;
  return pipes[pipe_pointer]; 
}

void stack_share(Array *arr) {
  arrays[array_pointer] = arr;
  arr->refcount++;
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

Array *stack_grab() {
  array_pointer--;
  Array *arr = arrays[array_pointer];
  arr->refcount--;
  arr->nesting = nesting;
  return arr;
}

void stack_stdout() {
  nesting++;
  Array *arr = stack_grab();
  
  for (int i = 0; i < arr->length; i++) {
    if (i < 11) {
      Serial.write(arr->data[i]);
    } else {
       Serial.write(*(stack_get_array_segment(arr, i)));
    }
    Serial.flush();
  }
  
  stack_collect();
}

//Toaster-level stdin.
void stack_stdin() {
	int mode = stack_pull();
	if (!mode) {
		mode = -10;
	}
	
	//Magic hacks. I don't think Serial wants to run in setu?!s
	Serial.println(F(" "));
	Serial.flush();
	
	Array *result = stack_array(0);
	
	if (mode > 0) {
		while (mode > 0) {
			if (Serial.available() > 0) {
				stack_put(Serial.read());
				mode--;
			}
		}
	} else {
		while (1) {
		  Serial.flush();
  		while (Serial.available() <= 0) { asm("nop"); }
  		int value = Serial.read();
  		
  		if (value == -mode) {
  			break;
  		}
  		Serial.flush();
  		stack_put(value);
  		Serial.flush();
		}
	}

	stack_share(result);	
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
	
	"SIZE":   is("%s->length", 1),
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

	"ARRAY":  is("Array *%s = stack_array(0);", 1),
	"MAKE":  is("stack_share(stack_array(stack_pull()));"),
	"RENAME": is("%s = activearray;", 1),
	
	"RELOAD": is("%s = stack.take();", 1),

	"SHARE": is("stack_share(%s);", 1),
	"GRAB":  is("Array *%s = stack_grab();", 1),

	"RELAY": is("stack_relay(%s);", 1),
	"TAKE":  is("int (*%s)() = stack_take();", 1),

	"GET": is("int %s = stack_get();", 1),
	"SET": is("stack_set(%s);", 1),

	"VAR": is("int %s = 0;", 1),

	"OPEN":   is("stack.open();"),
	"DELETE":   is("stack.delete();"),
	"EXECUTE":   is("stack.execute();"),
	"LOAD":   is("stack.load();"),
	"OUT":    is("stack.out();"),
	"STAT":   is("stack.info();"),
	"IN":     is("stack.in();"),
	"STDOUT": is("stack_stdout();"),
	"STDIN":  is("stack_stdin();"),
	"HEAP":   is(" "),
	"LINK":   is("stack.link();"),
	"CONNECT":   is("stack.connect();"),
	"SLICE":   is("stack.slice();"),

	"CLOSE": is("%s.close();", 1),

	"LOOP":   is("while (1) { nesting++;", 0, 1),
	"BREAK":  is("stack_collect(); break;"),
	"REPEAT": is("stack_collect(); }", 0, -1, -1),

	"IF":   is("if (%s) { nesting++;", 1, 1),
	"ELSE": is("stack_collect(); } else { nesting++;", 0, 0, -1),
	"END":  is("stack_collect(); }", 0, -1, -1),

	"RUN":  is("%s();", 1),
	
	"DATA":  is("Array %s[]=%s;", 2),
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
