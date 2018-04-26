package golang

import uct "github.com/qlova/uct/assembler"

var Runtime = uct.Runtime{
	Name: "runtime.go",
	Data: `
import "math/big"
import "io"
import "os"
import "bufio"
import "strconv"
import "errors"
	
type Runtime struct {
	Error Int
	Global List
	Channel Pipe

	//Stack
	Stack []Int
	Lists []*List
	Pipes []Pipe
	
	//Heap
	TheHeap []Int
	TheHeapRoom []int
	
	TheListHeap []*List
	TheListHeapRoom []int
	
	ThePipeHeap []Pipe
	ThePipeHeapRoom []int
}
	
func (runtime *Runtime) Init() {}

func (runtime *Runtime) Pull() Int {
	var result = runtime.Stack[len(runtime.Stack)-1]
	runtime.Stack = runtime.Stack[:len(runtime.Stack)-1]
	return result
}

func (runtime *Runtime) Sub() {
	runtime.Stack[len(runtime.Stack)-1] = runtime.Stack[len(runtime.Stack)-1].Flip()
	runtime.Add()
}

func (runtime *Runtime) Add() {
	var a, b = runtime.Stack[len(runtime.Stack)-2], runtime.Stack[len(runtime.Stack)-1]
	
	runtime.Stack = runtime.Stack[:len(runtime.Stack)-2]
	
	if ca, cb := a.Int.Bits() == nil, b.Int.Bits() == nil; ca || cb {
		if ca && cb {
			
			if (a.Small > 0 && b.Small > (1<<63 - 1) - a.Small) || (a.Small < 0 && b.Small < (-1 << 63) - a.Small) {
				var r big.Int
				r.Add(big.NewInt(a.Small), big.NewInt(b.Small))
				runtime.Stack = append(runtime.Stack, Int{Int:r})
				return
			}
			
			runtime.Stack = append(runtime.Stack, Int{Small:a.Small+b.Small})
			
		} else if !ca {
			
			var r big.Int
			r.Add(&a.Int, big.NewInt(b.Small))
			runtime.Stack = append(runtime.Stack, Int{Int:r})
			
		} else if !cb {
			
			var r big.Int
			r.Add(big.NewInt(a.Small), &b.Int)
			runtime.Stack = append(runtime.Stack, Int{Int:r})
			
		}
	} else {
		var r big.Int
		r.Add(&a.Int, &b.Int)
		runtime.Stack = append(runtime.Stack, Int{Int:r})
	}
}

func (runtime *Runtime) Mul() {
	var a, b = runtime.Stack[len(runtime.Stack)-2], runtime.Stack[len(runtime.Stack)-1]
	
	runtime.Stack = runtime.Stack[:len(runtime.Stack)-2]
	
	if ca, cb := a.Int.Bits() == nil, b.Int.Bits() == nil; ca || cb {
		if ca && cb {
			
			var z = Int{Small:a.Small*b.Small}
			runtime.Stack = append(runtime.Stack, z)
			
			if (a.Small != 0 && z.Small / a.Small != b.Small) {
				var r big.Int
				r.Mul(big.NewInt(a.Small), big.NewInt(b.Small))
				runtime.Stack[len(runtime.Stack)-1] = Int{Int:r}
				return
			}

		} else if !ca {
			
			var r big.Int
			r.Mul(&a.Int, big.NewInt(b.Small))
			runtime.Stack = append(runtime.Stack, Int{Int:r})
			
		} else if !cb {
			
			var r big.Int
			r.Mul(big.NewInt(a.Small), &b.Int)
			runtime.Stack = append(runtime.Stack, Int{Int:r})
			
		}
	} else {
		var r big.Int
		r.Mul(&a.Int, &b.Int)
		runtime.Stack = append(runtime.Stack, Int{Int:r})
	}
}

var Zero_go = big.NewInt(0)

func (runtime *Runtime) Div() {
	var a, b = runtime.Stack[len(runtime.Stack)-2], runtime.Stack[len(runtime.Stack)-1]
	runtime.Stack = runtime.Stack[:len(runtime.Stack)-2]
	
	defer func(a Int) {
        if r := recover(); r != nil {
		    if (len(a.Int.Bytes()) == 0 && a.Small == 0) || (len(a.Int.Bytes()) != 0 && a.Int.Cmp(Zero_go) == 0) {
		    	runtime.Stack = append(runtime.Stack, Int{Small:1})
		    	return
		    }
		   	runtime.Stack = append(runtime.Stack, Int{Small:0})
		}
    }(a)
	
	if ca, cb := a.Int.Bits() == nil, b.Int.Bits() == nil; ca || cb {
		if ca && cb {
			
			if b.Small < 0 && a.Small == (-1 << 63) {
				var r big.Int
				r.Div(big.NewInt(a.Small), big.NewInt(b.Small))
				runtime.Stack = append(runtime.Stack, Int{Int:r})
				return
			}
			
			runtime.Stack = append(runtime.Stack, Int{Small:a.Small/b.Small})

		} else if !ca {
			
			var r big.Int
			r.Div(&a.Int, big.NewInt(b.Small))
			runtime.Stack = append(runtime.Stack, Int{Int:r})
			
		} else if !cb {
			
			var r big.Int
			r.Div(big.NewInt(a.Small), &b.Int)
			runtime.Stack = append(runtime.Stack, Int{Int:r})
			
		}
	} else {
		var r big.Int
		r.Div(&a.Int, &b.Int)
		runtime.Stack = append(runtime.Stack, Int{Int:r})
	}
}

func (runtime *Runtime) Mod() {
	var a, b = runtime.Stack[len(runtime.Stack)-2], runtime.Stack[len(runtime.Stack)-1]
	
	runtime.Stack = runtime.Stack[:len(runtime.Stack)-2]
	
	if ca, cb := a.Int.Bits() == nil, b.Int.Bits() == nil; ca || cb {
		if ca && cb {
			
			runtime.Stack = append(runtime.Stack, Int{Small:a.Small%b.Small})

		} else if !ca {
			
			var r big.Int
			r.Mod(&a.Int, big.NewInt(b.Small))
			runtime.Stack = append(runtime.Stack, Int{Int:r})
			
		} else if !cb {
			
			var r big.Int
			r.Mod(big.NewInt(a.Small), &b.Int)
			runtime.Stack = append(runtime.Stack, Int{Int:r})
			
		}
	} else {
		var r big.Int
		r.Mod(&a.Int, &b.Int)
		runtime.Stack = append(runtime.Stack, Int{Int:r})
	}
}


func (runtime *Runtime) Pow() {
	var a, b = runtime.Stack[len(runtime.Stack)-2], runtime.Stack[len(runtime.Stack)-1]
	
	runtime.Stack = runtime.Stack[:len(runtime.Stack)-2]

	if ca, cb := a.Int.Bits() == nil, b.Int.Bits() == nil; ca || cb {
		if ca && cb {
			
			var za = Int{Int:*big.NewInt(a.Small)}
			var zb = Int{Int:*big.NewInt(b.Small)}
			
			za.Int.Exp(&za.Int, &zb.Int, nil)
			
			runtime.Stack = append(runtime.Stack, za)

		} else if !ca {
			
			var r big.Int
			r.Exp(&a.Int, big.NewInt(b.Small), nil)
			runtime.Stack = append(runtime.Stack, Int{Int:r})
			
		} else if !cb {
			
			var r big.Int
			r.Exp(big.NewInt(a.Small), &b.Int, nil)
			runtime.Stack = append(runtime.Stack, Int{Int:r})
			
		}
	} else {
		var r big.Int
		r.Exp(&a.Int, &b.Int, nil)
		runtime.Stack = append(runtime.Stack, Int{Int:r})
	}
} 

func (runtime *Runtime) Same() {
	var a, b = runtime.Stack[len(runtime.Stack)-1], runtime.Stack[len(runtime.Stack)-2]
	
	runtime.Stack = runtime.Stack[:len(runtime.Stack)-2]
	
	if ca, cb := a.Int.Bits() == nil, b.Int.Bits() == nil; ca || cb {
		if ca && cb {
			
			if a.Small == b.Small {
				runtime.Stack = append(runtime.Stack, Int{Small:1})
			} else {
				runtime.Stack = append(runtime.Stack, Int{Small:0})
			}
			
		} else if !ca {

			if a.Int.Cmp(big.NewInt(b.Small)) == 0 {
				runtime.Stack = append(runtime.Stack, Int{Small:1})
			} else {
				runtime.Stack = append(runtime.Stack, Int{Small:0})
			}
			
		} else if !cb {

			if big.NewInt(a.Small).Cmp(&b.Int) == 0 {
				runtime.Stack = append(runtime.Stack, Int{Small:1})
			} else {
				runtime.Stack = append(runtime.Stack, Int{Small:0})
			}
			
		}
	} else {
		if a.Int.Cmp(&b.Int) == 0 {
			runtime.Stack = append(runtime.Stack, Int{Small:1})
		} else {
			runtime.Stack = append(runtime.Stack, Int{Small:0})
		}
	}
}

func (runtime *Runtime) More() {
	var a, b = runtime.Stack[len(runtime.Stack)-1], runtime.Stack[len(runtime.Stack)-2]
	
	runtime.Stack = runtime.Stack[:len(runtime.Stack)-2]
	
	if ca, cb := a.Int.Bits() == nil, b.Int.Bits() == nil; ca || cb {
		if ca && cb {
			
			if a.Small > b.Small {
				runtime.Stack = append(runtime.Stack, Int{Small:1})
			} else {
				runtime.Stack = append(runtime.Stack, Int{Small:0})
			}
			
		} else if !ca {

			if a.Int.Cmp(big.NewInt(b.Small)) == 1 {
				runtime.Stack = append(runtime.Stack, Int{Small:1})
			} else {
				runtime.Stack = append(runtime.Stack, Int{Small:0})
			}
			
		} else if !cb {

			if big.NewInt(a.Small).Cmp(&b.Int) == 1 {
				runtime.Stack = append(runtime.Stack, Int{Small:1})
			} else {
				runtime.Stack = append(runtime.Stack, Int{Small:0})
			}
			
		}
	} else {
		if a.Int.Cmp(&b.Int) == 1 {
			runtime.Stack = append(runtime.Stack, Int{Small:1})
		} else {
			runtime.Stack = append(runtime.Stack, Int{Small:0})
		}
	}
}

func (runtime *Runtime) Less() {
	var a, b = runtime.Stack[len(runtime.Stack)-1], runtime.Stack[len(runtime.Stack)-2]
	
	runtime.Stack = runtime.Stack[:len(runtime.Stack)-2]
	
	if ca, cb := a.Int.Bits() == nil, b.Int.Bits() == nil; ca || cb {
		if ca && cb {
			
			if a.Small < b.Small {
				runtime.Stack = append(runtime.Stack, Int{Small:1})
			} else {
				runtime.Stack = append(runtime.Stack, Int{Small:0})
			}
			
		} else if !ca {

			if a.Int.Cmp(big.NewInt(b.Small)) == -1 {
				runtime.Stack = append(runtime.Stack, Int{Small:1})
			} else {
				runtime.Stack = append(runtime.Stack, Int{Small:0})
			}
			
		} else if !cb {

			if big.NewInt(a.Small).Cmp(&b.Int) == -1 {
				runtime.Stack = append(runtime.Stack, Int{Small:1})
			} else {
				runtime.Stack = append(runtime.Stack, Int{Small:0})
			}
			
		}
	} else {
		if a.Int.Cmp(&b.Int) == -1 {
			runtime.Stack = append(runtime.Stack, Int{Small:1})
		} else {
			runtime.Stack = append(runtime.Stack, Int{Small:0})
		}
	}
}

func (runtime *Runtime) Open() {
	var uri = runtime.Lists[len(runtime.Lists)-1]
	runtime.Lists = runtime.Lists[:len(runtime.Lists)-1]
	
	if len(uri.Bytes) == 0 {
		
		runtime.Pipes = append(runtime.Pipes, StdPipe{})
		
	} else {
		runtime.Error.Small = 1
		runtime.Error.SetBits(nil)
	}
}

func (runtime *Runtime) Load() {
	var name string
	var variable string
	
	var result = &List{}
	
	var uri = runtime.Lists[len(runtime.Lists)-1]
	runtime.Lists = runtime.Lists[:len(runtime.Lists)-1]
	
	if len(uri.Bytes) > 1 && uri.Bytes[0] == '$' {
	
		for i := 0; i < len(uri.Bytes); i++ {
			if i == 0 {
				continue
			}
			name += string(rune(uri.Bytes[i]))
		}
		variable = os.Getenv(name)

	} else if len(uri.Bytes) == 1 {
		
		if len(os.Args) > int(uri.Bytes[0]) {
			if char := int(uri.Bytes[0]); char < 0 {
				
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
					runtime.Error = Int{Small: 1}
				}
			} else {
				//This is not supposed to crash!
				variable = os.Args[uri.Mixed[0].Small]
			}
		} else {
			runtime.Error = Int{Small: 1}
		}
	}
	
	
	result.Bytes = []byte(variable)
	runtime.Lists = append(runtime.Lists, result)
}

func (runtime *Runtime) Send() {
	var pipe = runtime.Pipes[len(runtime.Pipes)-1]
	runtime.Pipes = runtime.Pipes[:len(runtime.Pipes)-1]
	
	var bytes = runtime.Lists[len(runtime.Lists)-1]
	runtime.Lists = runtime.Lists[:len(runtime.Lists)-1]
	
	_, err := pipe.Write(bytes.Bytes)
	if err != nil {
		if err.Error() == "function" {
			pipe.(*WrappedFunction).Function(runtime)
			return
		}
		runtime.Error.Small = 1
		runtime.Error.SetBits(nil)
	}
}

func (runtime *Runtime) Read() {
	var size = runtime.Stack[len(runtime.Stack)-1]
	runtime.Stack = runtime.Stack[:len(runtime.Stack)-1]
	
	var pipe = runtime.Pipes[len(runtime.Pipes)-1]
	runtime.Pipes = runtime.Pipes[:len(runtime.Pipes)-1]
	
	switch {
		case size.Small == 0:
			
			panic("unimplemented!")
			
		case size.Small > 0:
			
			var list = &List{
				Bytes: make([]byte, size.Small),
			}
			
			pipe.Read(list.Bytes)
			
			runtime.Lists = append(runtime.Lists, list)
			
		case size.Small < 0:
			
			var l = new(List)
			
			for {
				var b [1]byte
				
				pipe.Read(b[:])
				
				if b[0] == byte(-size.Small) {
					runtime.Lists = append(runtime.Lists, l)
					break
				}
				
				l.Put(Int{Small:int64(b[0])})
			}
	}
}

func (runtime *Runtime) HeapList() {
	address := runtime.Stack[len(runtime.Stack)-1]
	runtime.Stack = runtime.Stack[:len(runtime.Stack)-1]

	switch {
		//Put the value on the heap.
		case address.Small == 0:
			
			var list = runtime.Lists[len(runtime.Lists)-1]
			runtime.Lists = runtime.Lists[:len(runtime.Lists)-1]
			
			if len(runtime.TheListHeapRoom) > 0 {
				address := runtime.TheListHeapRoom[len(runtime.TheListHeapRoom)-1]
				runtime.TheListHeapRoom = runtime.TheListHeapRoom[:len(runtime.TheListHeapRoom)-1]
				
				runtime.TheListHeap[address] = list
				runtime.Stack = append(runtime.Stack, Int{Small:int64(address-1)})
				
			} else {
				runtime.TheListHeap = append(runtime.TheListHeap, list)
				runtime.Stack = append(runtime.Stack, Int{Small:int64(len(runtime.TheListHeap))})
			}
			
		//Get from the heap.
		case address.Small > 0:
			runtime.Lists = append(runtime.Lists, runtime.TheListHeap[int64(address.Small)-1])
		
		//Free heap space.
		case address.Small < 0:
			address.Small = -address.Small
			runtime.TheListHeap[int(address.Small-1)] = nil
			runtime.TheListHeapRoom = append(runtime.TheListHeapRoom, int(address.Small)-1)
	}
}

type Int struct {
	big.Int
	Small int64
}

func (value Int) Init() {}

func (value Int) True() bool {
	return !(value.Bits() == nil && value.Small == 0)
}

func (value Int) Flip() Int {
	var result big.Int
	if value.Bits() == nil {
		return Int{Small: int64(^uint64(int64(value.Small)-1))}
	}
	result.Neg(&value.Int)
	return Int{Int:result}
}

type List struct {
	Mixed []Int
	Bytes []byte
}

func (list *List) Init() {}

//Maybe optimise memory usage here?
//Bug here.
func (list *List) Put(value Int) {
	list.Bytes = append(list.Bytes, byte(value.Small))
	list.Mixed = append(list.Mixed, value)
}

//Maybe optimise memory usage here?
func (list *List) Pop() Int {
	var result = list.Mixed[len(list.Mixed)-1]
	list.Bytes = list.Bytes[:len(list.Bytes)-1]
	list.Mixed = list.Mixed[:len(list.Mixed)-1]
	return result
}

type Pipe interface {
	io.ReadWriter
	io.Closer
	io.Seeker
	Info(List) List
	Init()
}

var Stdin = bufio.NewReader(os.Stdin)
var Stdout = bufio.NewWriter(os.Stdout)
var IntSize = strconv.IntSize

type StdPipe struct {}

func (StdPipe) Init() {}

func (StdPipe) Read(p []byte) (n int, err error) {
	Stdout.Flush()
	return Stdin.Read(p)
}

func (StdPipe) Write(p []byte) (n int, err error) {
	return Stdout.Write(p)
}

func (StdPipe) Seek(offset int64, whence int) (int64, error) {
	return 0, nil
}

func (StdPipe) Close() (error) {
	return nil
}

func (StdPipe) Info(data List) (List) {
	return List{}
}

type WrappedFunction struct {
	Function func(*Runtime)
	Data List
}

func (WrappedFunction) Init() {}

func (wrapped *WrappedFunction) Read(p []byte) (n int, err error) {
	return 0, nil
}

func (wrapped *WrappedFunction) Write(p []byte) (n int, err error) {
	return 0, errors.New("function")
}

func (wrapped *WrappedFunction) Seek(offset int64, whence int) (int64, error) {
	return 0, nil
}

func (wrapped *WrappedFunction) Close() (error) {
	return nil
}

func (wrapped *WrappedFunction) Info(data List) (List) {
	if data.Mixed == nil {
		return wrapped.Data
	} else {
		wrapped.Data = data
	}
	return List{}
}

`}
