//Compiled to Go with UCT (Universal Code Translator)
package main

import "math/big"
import "os"
import "os/exec"
import "io"
import "crypto/rand"
import "net"
import "strings"
import "strconv"
import "bufio"
import "net/http"
import "io/ioutil"
import "sync"

var Networks_In = make(map[string]net.Listener)

//This is the Go stack implementation.
// It holds arrays for the 3 types:
//		Numbers
//		Arrays
//		Pipes
//
// It also holds the ERROR variable for the current thread.
// The currently active array is stored as ActiveArray.
type Stack struct {
	sync.Mutex

	Numbers *Array
	Arrays 	[]*Array
	Pipes	[]*Pipe
	
	ERROR Number
	ActiveArray *Array
	TheHeap []*Array;
	TheITHeap []*Pipe;
	HeapRoom []int
	ITHeapRoom []int
	Map map[string]Number
	
	Inbox chan *Array
	Outbox chan *Array
	
	//Experimental.
	Eval map[string]func(*Stack)
}

func (stack *Stack) Copy() (n *Stack) {
	n = new(Stack)
	n.Eval = stack.Eval
	
	n.Numbers = &Array{}
	n.Numbers.Small = make([]byte, len(stack.Numbers.Small))
	copy(n.Numbers.Small, stack.Numbers.Small)
	n.Numbers.Big = make([]Number, len(stack.Numbers.Big))
	copy(n.Numbers.Big, stack.Numbers.Big)
	
	n.Arrays = make([]*Array, len(stack.Arrays))
	for i := range stack.Arrays {
		n.Arrays[i] = &Array{}
		n.Arrays[i].Small = make([]byte, len(stack.Arrays[i].Small))
		copy(n.Arrays[i].Small, stack.Arrays[i].Small)
		
		n.Arrays[i].Big = make([]Number, len(stack.Arrays[i].Big))
		copy(n.Arrays[i].Big, stack.Arrays[i].Big)
	}
	
	n.Outbox = make(chan *Array, 1)
	stack.Inbox = n.Outbox
	
	n.Pipes = make([]*Pipe, len(stack.Pipes))
	copy(n.Pipes, stack.Pipes)
	
	return
}

func (stack *Stack) Pipe() *Pipe {
	var pipe Pipe
	pipe.Channel = make(chan *Array)
	pipe.Data = &Array{}
	return &pipe
}

func (stack *Stack) Array() *Array {
	var array Array
	stack.ActiveArray = &array
	return &array
}

func (z *Stack) Init() {
	z.Map = make(map[string]Number)
	z.Numbers = &Array{}
	z.Eval = make(map[string]func(*Stack))
}

func (z *Stack) Exit(n int) {
	os.Exit(n)
}

func (stack *Stack) Execute() {
	text := stack.Grab().String()
	
	out, err := exec.Command("sh", "-c", text).Output()
    if err != nil {
        stack.ERROR = NewNumber(1)
        stack.Share(&Array{})
		return
    }
    
	stack.Share(&Array{Small:out})
}

func (stack *Stack) Delete() {
	text := stack.Grab().String()
	
	if text == "" {
		stack.ERROR = NewNumber(1)
		return
	}
	
	if os.Remove(text) != nil {
		stack.ERROR = NewNumber(1)
	}
	
}

func (stack *Stack) Move() {
	dest := stack.Grab().String()
	src := stack.Grab().String()
	
	if dest == "" || src == "" {
		stack.ERROR = NewNumber(1)
		return
	}
	
	if os.Rename(src, dest) != nil {
		stack.ERROR = NewNumber(1)
	}
	
}

func (stack *Stack) Share(array *Array) {
	stack.Arrays = append(stack.Arrays, array)
}
func (stack *Stack) Grab() (array *Array) {
	array = stack.Arrays[len(stack.Arrays)-1]
	stack.Arrays = stack.Arrays[:len(stack.Arrays)-1]
	return
}

func (stack *Stack) Relay(pipe *Pipe) {
	stack.Pipes = append(stack.Pipes, pipe)
}
func (stack *Stack) Take() (pipe *Pipe) {
	pipe = stack.Pipes[len(stack.Pipes)-1]
	stack.Pipes = stack.Pipes[:len(stack.Pipes)-1]
	return
}

func (stack *Stack) Link() {
	array := stack.Grab()
	if array.Big == nil {
		stack.Map[string(array.Small)] = stack.Pull()
	} else {
		stack.Map[array.String()] = stack.Pull()
	}
}
func (stack *Stack) Connect() {
	array := stack.Grab()
	if array.Big == nil {
		n, ok := stack.Map[string(array.Small)]
		stack.Push(n)
		if !ok {
			stack.ERROR = NewNumber(1)
		}
	} else {
		n, ok := stack.Map[array.String()]
		stack.Push(n)
		if !ok {
			stack.ERROR = NewNumber(1)
		}
	} 
}


func (stack *Stack) Slice() {
	array := stack.Grab()
	if array.Big == nil {
		stack.Share(&Array{Small:array.Small[stack.Pull().ToInt():stack.Pull().ToInt()]})
	} else {
		stack.Share(&Array{Big:array.Big[stack.Pull().ToInt():stack.Pull().ToInt()]})
	} 
}

func (stack *Stack) Push(number Number) {
	if number.Int == nil && number.Small < 256 && number.Small >= 0 && stack.Numbers.Big == nil {
		stack.Numbers.Small = append(stack.Numbers.Small, byte(number.Small))
	} else {
		stack.Numbers.Grow()
		stack.Numbers.Big = append(stack.Numbers.Big, number)
	}
}
func (stack *Stack) Pull() (number Number) {
	if stack.Numbers.Big == nil {
		number.Small = int64(stack.Numbers.Small[len(stack.Numbers.Small)-1])
		stack.Numbers.Small = stack.Numbers.Small[:len(stack.Numbers.Small)-1]
	} else {
		stack.Numbers.Grow()
		number = stack.Numbers.Big[len(stack.Numbers.Big)-1]
		stack.Numbers.Big = stack.Numbers.Big[:len(stack.Numbers.Big)-1]
	}
	return
}

func (stack *Stack) Heap() {
	address := stack.Pull()
	
	switch {
		case address.Small == 0:
			if len(stack.HeapRoom) > 0 {
				address := stack.HeapRoom[len(stack.HeapRoom)-1]
				stack.HeapRoom = stack.HeapRoom[:len(stack.HeapRoom)-1]
				
				stack.TheHeap[(address%(len(stack.TheHeap)+1))-1] = stack.Grab()
				stack.Push(NewNumber(address))
				
			} else {
				stack.TheHeap = append(stack.TheHeap, stack.Grab())
				stack.Push(NewNumber(len(stack.TheHeap)))
			}
			
		case address.Small > 0:
			stack.Share(stack.TheHeap[(int(address.Small)%(len(stack.TheHeap)+1))-1])
			
		case address.Small < 0:
			stack.TheHeap[(-int(address.Small))%(len(stack.TheHeap)+1)-1] = &Array{}
			stack.HeapRoom = append(stack.HeapRoom, -int(address.Small))
	}
}

func (stack *Stack) HeapIt() {
	address := stack.Pull()
	
	switch {
		case address.Small == 0:
			if len(stack.ITHeapRoom) > 0 {
				address := stack.ITHeapRoom[len(stack.HeapRoom)-1]
				stack.ITHeapRoom = stack.ITHeapRoom[:len(stack.ITHeapRoom)-1]
				
				stack.TheITHeap[(address%(len(stack.TheITHeap)+1))-1] = stack.Take()
				stack.Push(NewNumber(address))
				
			} else {
				stack.TheITHeap = append(stack.TheITHeap, stack.Take())
				stack.Push(NewNumber(len(stack.TheITHeap)))
			}
			
		case address.Small > 0:
			stack.Relay(stack.TheITHeap[(int(address.Small)%(len(stack.TheITHeap)+1))-1])
			
		case address.Small < 0:
			stack.TheITHeap[(-int(address.Small))%(len(stack.TheITHeap)+1)-1] = &Pipe{}
			stack.ITHeapRoom = append(stack.ITHeapRoom, -int(address.Small))
	}
}


func (stack *Stack) Put(number Number) {
	var array = stack.ActiveArray
	if number.Int == nil && number.Small < 256  && number.Small >= 0 && array.Big == nil {
		array.Small = append(array.Small, byte(number.Small))
	} else {
		array.Grow()
		array.Big = append(array.Big, number)
	}
}
func (stack *Stack) Pop() (number Number) {
	var array = stack.ActiveArray
	if array.Big == nil {
		number.Small = int64(array.Small[len(array.Small)-1])
		array.Small = array.Small[:len(array.Small)-1]
	} else {
		array.Grow()
		number = array.Big[len(array.Big)-1]
		array.Big = array.Big[:len(array.Big)-1]
	}
	return
}

func (stack *Stack) Place(array *Array) {
	stack.ActiveArray = array
}
func (stack *Stack) Get() (number Number) {
	var array = stack.ActiveArray
	var index Number
	
	if array.Big == nil {
		if len(array.Small) == 0 {
			number = NewNumber(0)
			stack.ERROR = NewNumber(4)
			return
		}
		index.Mod(stack.Pull(), NewNumber(len(array.Small)))
		number = NewNumber(int(array.Small[index.ToInt()]))
	} else {	
		if len(array.Big) == 0 {
			number = NewNumber(0)
			stack.ERROR = NewNumber(4)
			return
		}
		index.Mod(stack.Pull(), NewNumber(len(array.Big)))
		number = array.Big[index.ToInt()]
	}
	return
}
func (stack *Stack) Set(number Number) {
	var array = stack.ActiveArray
	var index Number
	
	if number.Int == nil && array.Big == nil && number.Small < 256  && number.Small >= 0 {
		if len(array.Small) == 0 {
			number = NewNumber(0)
			stack.ERROR = NewNumber(4)
			return
		}
		index.Mod(stack.Pull(), NewNumber(len((*stack.ActiveArray).Small)))
		array.Small[index.ToInt()] = byte(number.Small)
	} else {
		array.Grow()
		if len(array.Big) == 0 {
			number = NewNumber(0)
			stack.ERROR = NewNumber(4)
			return
		}
		index.Mod(stack.Pull(), NewNumber(len((*stack.ActiveArray).Big)))
		array.Big[index.ToInt()] = number
	}
}

func (stack *Stack) Load() {
	var name string
	var variable string
	var result = stack.Array()
	var err error
	
	text := stack.Grab()
	
	if text.Index(0).ToInt() == '$' && text.Len().Small > 1 {
	
		for i := 0; i < int(text.Len().Small); i++ {
			if i == 0 {
				continue
			}
			name += string(rune(text.Index(i).ToInt()))
		}
		variable = os.Getenv(name)
	} else {
	
		name = text.String()
		
		protocol := strings.SplitN(name, "://", 2)
		if len(protocol) > 1 {
			switch protocol[0] {
				case "tcp":
					listener, err := net.Listen("tcp", ":"+protocol[1])
					
					if err != nil {
					
						stack.ERROR = NewNumber(1)
						
					} else {
						_, variable, _ = net.SplitHostPort(listener.Addr().String())
						if protocol[1] == "0" {
							Networks_In[variable] = listener
						} else {
							Networks_In[protocol[1]] = listener
						}
					}
				case "dns":
					//This can be optimised. Check the string.
					hosts, err := net.LookupAddr(protocol[1])
					if err != nil {
						hosts, err = net.LookupHost(protocol[1])
						if err != nil {
							stack.ERROR = NewNumber(1)
						}
					}
					variable = strings.Join(hosts, " ")
				default:
					if err != nil {
						stack.ERROR = NewNumber(1)
					}
			}
		} else {
	
			if len(os.Args) > int(text.Index(0).ToInt()) {
				if char := text.Index(0).ToInt(); char < 0 {
					var found bool
					for i, v := range os.Args {
						if len(v) == 2 && v[0] == '-' && v[1] == byte(-char) {
							found = true
							if i+1 < len(os.Args) {
								variable = os.Args[i+1]
							}
						}
					}
					if !found {
						stack.ERROR = NewNumber(1)
					}
				} else {
					variable = os.Args[char]
				}
			} else {
				stack.ERROR = NewNumber(1)
			}
		}
	}
	
	result.Small = []byte(variable)
	stack.Share(result)
}

func (stack *Stack) Open() {
	var err error
	
	text := stack.Grab()
	
	var name string
	for i := 0; i < int(text.Len().Small); i++ {
		name += string(rune(text.Index(i).ToInt()))
	}

	var filename string = text.String()
	
	
	var it Pipe
	it.Name = filename
	
	protocol := strings.SplitN(filename, "://", 2)
	if len(protocol) > 1 {
		switch protocol[0] {
			case "dir":
				it.Files, err = ioutil.ReadDir(protocol[1])
				if err != nil {
					stack.ERROR = NewNumber(404)
					stack.Relay(&it)
					return
				}
				stack.Relay(&it)
				return
		
			case "http":
				resp, err := http.Get(filename)
				it.In = resp.Body
				if err != nil {
					stack.ERROR = NewNumber(404)
					stack.Relay(&it)
					return
				}
				stack.Relay(&it)
				return
		
			case "tcp":
				if listener, ok := Networks_In[protocol[1]]; ok {
					
					it.Connection, err = listener.Accept()
					it.In = it.Connection
					it.Out = it.Connection
					if err != nil {
						stack.ERROR = NewNumber(404)
						stack.Relay(&it)
						return
					}
					stack.Relay(&it)
					return
					
				} else {
					it.Connection, err = net.Dial("tcp", protocol[1])
					it.In = it.Connection
					it.Out = it.Connection
					if err != nil {
						stack.ERROR = NewNumber(404)
						stack.Relay(&it)
						return
					}
					stack.Relay(&it)
					return
				}
		}
	}

	f, err := os.OpenFile(filename, os.O_RDWR|os.O_APPEND, 0666)
	if err != nil {
		f, err = os.Open(filename)
	}
	if err == nil {
		stat, _ := f.Stat()
		it.In = f
		it.Out = f
		it.Size = int(stat.Size())
		stack.Relay(&it)
		return
	}
	it.In = nil
	it.Out = nil
	
	/*it.Pipe, err = os.Create(filename)
	if err == nil {
		stack.Push(NewNumber(0))
		stack.Relay(&it)
		return
	}*/
	
	if stat, err := os.Stat(filename); err == nil {
		it.Size = int(stat.Size()) 
		stack.Relay(&it)
		return
	}
	stack.ERROR = NewNumber(404)
	stack.Relay(&it)
	return
}

func (stack *Stack) Info() {
	var request string
	var variable string
	
	var result = stack.Array()

	text := stack.Grab()
	it := stack.Take()
	
	request = text.String()
	
	switch request {
		case "address":
			if it.Connection != nil {
				variable = it.Connection.RemoteAddr().String()
			}
		case "ip":
			if it.Connection != nil {
				variable = it.Connection.RemoteAddr().(*net.TCPAddr).IP.String()
			}
		case "port":
			if it.Connection != nil {
				variable = strconv.Itoa(it.Connection.RemoteAddr().(*net.TCPAddr).Port)
			}
		default:
			stack.Put(NewNumber(it.Size))
			stack.Share(result)
			return
	}
	
	result.Small = []byte(variable)
	stack.Share(result)
}

func (stack *Stack) Out() {
	var err error
	
	text := stack.Grab()
	f := stack.Take()
	
	if f.Channel != nil {
		f.Channel <- text
		return
	}
	
	if f.Out == nil {
		if f.Name[len(f.Name)-1] == '/' {
			i, err := os.Stat(f.Name)
			if err == nil && i.IsDir() {
				
			} else {
				err := os.Mkdir(f.Name, 0666)
				if err != nil {
					stack.Push(NewNumber(-1))
					return
				}
			}
		} else {
			f.Out, err = os.Create(f.Name)
			if err != nil {
				stack.Push(NewNumber(-1))
				return
			}
		}
	}
	if int(text.Len().Small) == 0 {
		stack.Push(NewNumber(0))
		return
	}
	for i := 0; i < int(text.Len().Small); i++ {
		v := text.Index(i)
		if v.Int != nil {
			_, err := f.Out.Write(v.Bytes())
			if err != nil {
				stack.Push(NewNumber(-1))
				return
			}
		} else {
			_, err := f.Out.Write([]byte{byte(v.Small)})
			if err != nil {
				stack.Push(NewNumber(-1))
				return
			}
		}
	} 
	
	stack.Push(NewNumber(0))
}

func (stack *Stack) Stdout() {
	text := stack.Grab()
	
	if text.Big == nil {
		os.Stdout.Write(text.Small)
		return
	}
	
	for i := 0; i < int(text.Len().Small); i++ {
		v := text.Index(i)
		if v.Int != nil {
			os.Stdout.Write(v.Bytes())
		} else {
			os.Stdout.Write([]byte{byte(v.Small)})
		}
	} 
}

func (stack *Stack) In() {
	length := stack.Pull()
	f := stack.Take()
	
	if f.Channel != nil {
		stack.Share(<-f.Channel)
		return
	}
	
	//Null pipe.
	if f.In == nil {
	
			
		if f.Files != nil {
			if f.Index < len(f.Files) {
				stack.Share(NewStringArray(f.Files[f.Index].Name()))
				f.Index++
				return
			}
		}
	
		stack.ERROR = NewNumber(1)
		stack.Share(&Array{})
		return
	}

	switch {
		case length.Small > 0, length.Int != nil:
		
			var b []byte = make([]byte, int(length.ToInt()))
			n, err := f.In.Read(b)
			if err != nil || n <= 0 {
				//println(err.Error())
				if n <= 0 {
					b = []byte{}
				}
				stack.Share(&Array{Small:b})
				stack.ERROR = NewNumber(1)
				return
			}
			stack.Share(&Array{Small:b[:n]})
		case length.Small == 0:
			length.Small = -10
			fallthrough
			
		case length.Small < 0:
			var result []byte = make([]byte, 0, 32)
			var b []byte = make([]byte, 1)
			for {
				n, err := f.In.Read(b)
				
				if err != nil || n <= 0 {
					stack.Share(&Array{Small:result})
					stack.ERROR = NewNumber(1)
					return
				}
				
				if b[0] == byte(-length.Small) {
					stack.Share(&Array{Small:result})
					return
				}
				
				result = append(result, b[0])
			}
			
	}

}

var stdin_reader = bufio.NewReader(os.Stdin)

func  (stack *Stack) Stdin() {
	length := stack.Pull()
	
	switch {
			
		case length.Small > 0, length.Int != nil:
		
			var b []byte = make([]byte, int(length.ToInt()))
			n, err := os.Stdin.Read(b)
			if err != nil || n <= 0 {
				stack.Share(&Array{})
				return
			}
			stack.Share(&Array{Small:b})
		
		case length.Small == 0:
			
			b, err := stdin_reader.ReadBytes('\n')
			if err != nil || len(b) == 0 {
				stack.Share(&Array{})
				return
			}
			if len(b) > 2 && b[len(b)-2] == '\r' {
				b = b[:len(b)-1]
			}
			stack.Share(&Array{Small:b[:len(b)-1]})
			
		case length.Small < 0:
			b, err := stdin_reader.ReadBytes(byte(-length.Small))
			if err != nil || len(b) == 0 {
				stack.Share(&Array{})
				return
			}
			stack.Share(&Array{Small:b[:len(b)-1]})
	}
}

type Number struct {
	*big.Int
	Small int64
}

func (z *Number) Init() {

}

func NewNumber(n int) Number {
	return Number{Small:int64(n)}
}

func NewNumberFromString(n string) Number {
	var num = Number{}
	num.Int = big.NewInt(0)
	num.SetString(n, 10)

	return num
}

func (z *Number) Add(a, b Number) {
	z.Int = big.NewInt(0)
	if ca, cb := a.Int == nil, b.Int == nil; ca || cb {
		if !ca {
			z.Int.Add(a.Int, big.NewInt(b.Small))
		} else if !cb {
			z.Int.Add(big.NewInt(a.Small), b.Int)
		} else {
			if (a.Small > 0 && b.Small > (1<<63 - 1) - a.Small) {
				z.Int.Add(big.NewInt(a.Small), big.NewInt(b.Small))
				return
			} else if (a.Small < 0 && b.Small < (-1 << 63) - a.Small) {
				z.Int.Add(big.NewInt(a.Small), big.NewInt(b.Small))
				return
			}
			*z = Number{Small:a.Small+b.Small}
		}
	} else {
		z.Int.Add(a.Int, b.Int)
	}
} 

func (z *Number) Sub(a, b Number) {
	z.Int = big.NewInt(0)
	if ca, cb := a.Int == nil, b.Int == nil; ca || cb {
		if !ca {
			z.Int.Sub(a.Int, big.NewInt(b.Small))
		} else if !cb {
			z.Int.Sub(big.NewInt(a.Small), b.Int)
		} else {
			if (a.Small > 0 && -b.Small > (1<<63 - 1) - a.Small) {
				z.Int.Sub(big.NewInt(a.Small), big.NewInt(b.Small))
				return
			} else if (a.Small < 0 && -b.Small < (-1 << 63) - a.Small) {
				z.Int.Sub(big.NewInt(a.Small), big.NewInt(b.Small))
				return
			}
			*z = Number{Small:a.Small-b.Small}
		}
	} else {
		z.Int.Sub(a.Int, b.Int)
	}
} 

func (z *Number) Mul(a, b Number) {
	
	z.Int = big.NewInt(0)
	if ca, cb := a.Int == nil, b.Int == nil; ca || cb {
		if !ca {
			z.Int.Mul(a.Int, big.NewInt(b.Small))
		} else if !cb {
			z.Int.Mul(big.NewInt(a.Small), b.Int)
		} else {
			*z = Number{Small:a.Small*b.Small}
			if (a.Small != 0 && z.Small / a.Small != b.Small) {
				z.Int = big.NewInt(0)
				z.Int.Mul(big.NewInt(a.Small), big.NewInt(b.Small))
			}
		}
	} else {
		z.Int.Mul(a.Int, b.Int)
	}
} 

var Zero_go = big.NewInt(0)

func (z *Number) Div(a, b Number) {
	defer func() {
        if r := recover(); r != nil {
		    if (a.Int == nil && a.Small == 0) || (a.Int != nil && a.Int.Cmp(Zero_go) == 0) {
		    	var b []byte = []byte{1}
		    	rand.Read(b)
		    	*z =NewNumber(int(b[0]+1))
		    	return
		    }
		   	*z = NewNumber(0)
		}
    }()
	z.Int = big.NewInt(0)
	if ca, cb := a.Int == nil, b.Int == nil; ca || cb {
		if !ca {
			z.Int.Div(a.Int, big.NewInt(b.Small))
		} else if !cb {
			z.Int.Div(big.NewInt(a.Small), b.Int)
		} else {
			if b.Small < 0 && a.Small == (-1 << 63) {
				z.Int = big.NewInt(0)
				z.Int.Div(big.NewInt(a.Small), big.NewInt(b.Small))
				return
			}
			*z = Number{Small:a.Small/b.Small}
		}
	} else {
		z.Int.Div(a.Int, b.Int)
	}
} 

func (z *Number) Mod(a, b Number) {
	defer func() {
        if r := recover(); r != nil {
		    if a.Small == 0 || a.Int != nil && a.Int.Cmp(Zero_go) == 0 {
		    	var b []byte = []byte{1}
		    	rand.Read(b)
		    	*z =NewNumber(int(b[0]+1))
		    	return
		    }
		   	*z = NewNumber(0)
		}
    }()
	z.Int = big.NewInt(0)
	if ca, cb := a.Int == nil, b.Int == nil; ca || cb {
		if !ca {
			z.Int.Mod(a.Int, big.NewInt(b.Small))
		} else if !cb {
			z.Int.Mod(big.NewInt(a.Small), b.Int)
		} else {
			z.Int = big.NewInt(0)
			z.Int.Mod(big.NewInt(a.Small), big.NewInt(b.Small))
			return
		}
	} else {
		z.Int.Mod(a.Int, b.Int)
	}
} 

func (z *Number) Pow(a, b Number) {
	z.Int = big.NewInt(0)
	if ca, cb := a.Int == nil, b.Int == nil; ca || cb {
		if !ca {
			z.Int.Exp(a.Int, big.NewInt(b.Small), nil)
		} else if !cb {
			z.Int.Exp(big.NewInt(a.Small), b.Int, nil)
		} else {
			z.Int = big.NewInt(0)
			z.Int.Exp(big.NewInt(a.Small), big.NewInt(b.Small), nil)
			return
		}
	} else {
		z.Int.Exp(a.Int, b.Int, nil)
	}
} 

func (a Number) Slt(b Number) Number {
	if ca, cb := a.Int == nil, b.Int == nil; ca || cb {
		if !ca {
			return NewNumber(__bool2int( a.Cmp(big.NewInt(b.Small)) == -1 ))
		} else if !cb {
			return NewNumber(__bool2int( big.NewInt(a.Small).Cmp(b.Int) == -1 ))
		} else {
			return NewNumber(__bool2int( a.Small < b.Small ))
		}
	} else {
		return NewNumber(__bool2int( a.Cmp(b.Int) == -1 ))
	}
}

func (a Number) Seq(b Number) Number {
	if ca, cb := a.Int == nil, b.Int == nil; ca || cb {
		if !ca {
			return NewNumber(__bool2int( a.Cmp(big.NewInt(b.Small)) == 0 ))
		} else if !cb {
			return NewNumber(__bool2int( big.NewInt(a.Small).Cmp(b.Int) == 0 ))
		} else {
			return NewNumber(__bool2int( a.Small == b.Small ))
		}
	} else {
		return NewNumber(__bool2int( a.Cmp(b.Int) == 0 ))
	}
}

func (a Number) Sge(b Number) Number {
	if ca, cb := a.Int == nil, b.Int == nil; ca || cb {
		if !ca {
			return NewNumber(__bool2int( a.Cmp(big.NewInt(b.Small)) >= 0 ))
		} else if !cb {
			return NewNumber(__bool2int( big.NewInt(a.Small).Cmp(b.Int) >= 0 ))
		} else {
			return NewNumber(__bool2int( a.Small >= b.Small ))
		}
	} else {
		return NewNumber(__bool2int( a.Cmp(b.Int) >= 0 ))
	}
}

func (a Number) Sle(b Number) Number {
	if ca, cb := a.Int == nil, b.Int == nil; ca || cb {
		if !ca {
			return NewNumber(__bool2int( a.Cmp(big.NewInt(b.Small)) <= 0 ))
		} else if !cb {
			return NewNumber(__bool2int( big.NewInt(a.Small).Cmp(b.Int) <= 0 ))
		} else {
			return NewNumber(__bool2int( a.Small <= b.Small ))
		}
	} else {
		return NewNumber(__bool2int( a.Cmp(b.Int) <= 0 ))
	}
}

func (a Number) Sgt(b Number) Number {
	if ca, cb := a.Int == nil, b.Int == nil; ca || cb {
		if !ca {
			return NewNumber(__bool2int( a.Cmp(big.NewInt(b.Small)) == 1 ))
		} else if !cb {
			return NewNumber(__bool2int( big.NewInt(a.Small).Cmp(b.Int) == 1 ))
		} else {
			return NewNumber(__bool2int( a.Small > b.Small ))
		}
	} else {
		return NewNumber(__bool2int( a.Cmp(b.Int) == 1 ))
	}
}

func (a Number) Sne(b Number) Number {
	if ca, cb := a.Int == nil, b.Int == nil; ca || cb {
		if !ca {
			return NewNumber(__bool2int( a.Cmp(big.NewInt(b.Small)) != 0 ))
		} else if !cb {
			return NewNumber(__bool2int( big.NewInt(a.Small).Cmp(b.Int) != 0 ))
		} else {
			return NewNumber(__bool2int( a.Small != b.Small ))
		}
	} else {
		return NewNumber(__bool2int( a.Cmp(b.Int) != 0 ))
	}
}

func (z Number) True() bool {
	if z.Int == nil {
		if z.Small != 0 {
			return true
		}
		return false
	}
	return z.Int.Cmp(big.NewInt(0)) != 0
} 

func (z Number) ToInt() int {
	if z.Int == nil {
		return int(z.Small)
	}
	return int(z.Int.Int64())
} 


type Pipe struct {
	Function func(*Stack)
	Data *Array
	
	//This stuff needs to be wrapped up in an interface.
		Name string
		
		Size int
		
		Pipe io.ReadWriteCloser
		In io.ReadCloser
		Out io.WriteCloser
		
		Channel chan *Array
		
		Connection net.Conn
		Reader *bufio.Reader
		
		//Other protocols.
		Files []os.FileInfo
		Index int
	

}

func (z *Pipe) Init() {

}

func (pipe *Pipe) Exe(stack *Stack) {
	if pipe.Function != nil {
		pipe.Function(stack)
	}
}

func (pipe *Pipe) Close() {
	if pipe.In != nil {
		pipe.In.Close()
	}
	if pipe.Out != nil {
		pipe.Out.Close()
	}
}

type Array struct {
	Big []Number
	Small []byte
}

func (z *Array) Init() {

}

func (z *Array) Grow() {
	if z.Big == nil {
		z.Big = make([]Number, len(z.Small))
		for i := range z.Small {
			z.Big[i] = Number{Small:int64(z.Small[i])}
		}
	}
}

func (z *Array) Index(n int) Number {
	if z.Big == nil {
		return Number{Small:int64(z.Small[n])}
	} else {
		return z.Big[n]
	}
}

func (z *Array) Len() Number {
	if z.Big == nil {
		return Number{Small:int64(len(z.Small))}
	} else {
		return Number{Small:int64(len(z.Big))}
	}
}

func (z *Array) String() string {
	if z.Big == nil {
		return string(z.Small)
	} else {
		var name string
		for i := 0; i < int(z.Len().Small); i++ {
			name += string(rune(z.Big[i].ToInt()))
		}
		return name
	}
}

func NewStringArray(s string) *Array {
	var result = Array{Small:[]byte(s)}
	return &result
}

func (array *Array) Join(b *Array) *Array {
	switch {
		case array.Big == nil && b.Big == nil:
			return &Array{Small:append(array.Small, b.Small...)}
			
		case (array.Small != nil && b.Big != nil) || (array.Big != nil && b.Small != nil):
			array.Grow()
			b.Grow()
			fallthrough
		case array.Big != nil && b.Big != nil:
			return &Array{Big:append(array.Big, b.Big...)}
	}
	//TODO, copy b array
	return b
}

func __bool2int(b bool) int {
	if b {
		return 1
	}
	return 0
}
