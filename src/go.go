package main

import "errors"
import "flag"
import "strings"
import "strconv"

//Go flag.
var Go bool
func init() {
	flag.BoolVar(&Go, "go", false, "Target Go")
	
	RegisterAssembler(new(GoAssembler), &Go, "go", "//")
}

type GoAssembler struct {
	Indentation int
}

func (g *GoAssembler) SetFileName(s string) {

}

func (g *GoAssembler) Header() []byte {
	return []byte(
	`
package main

import "math/big"
import "os"
import "fmt"

type AnyInt struct {
	*big.Int
	Small int64
}

func (z *AnyInt) Add(a, b AnyInt) {
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
			*z = AnyInt{Small:a.Small+b.Small}
		}
	} else {
		z.Int.Add(a.Int, b.Int)
	}
} 

func (z *AnyInt) Sub(a, b AnyInt) {
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
			*z = AnyInt{Small:a.Small-b.Small}
		}
	} else {
		z.Int.Sub(a.Int, b.Int)
	}
} 

func (z *AnyInt) Mul(a, b AnyInt) {
	z.Int = big.NewInt(0)
	if ca, cb := a.Int == nil, b.Int == nil; ca || cb {
		if !ca {
			z.Int.Mul(a.Int, big.NewInt(b.Small))
		} else if !cb {
			z.Int.Mul(big.NewInt(a.Small), b.Int)
		} else {
			*z = AnyInt{Small:a.Small*b.Small}
			if (a.Small != 0 && z.Small / a.Small != b.Small) {
				z.Int = big.NewInt(0)
				z.Int.Mul(big.NewInt(a.Small), big.NewInt(b.Small))
			}
		}
	} else {
		z.Int.Mul(a.Int, b.Int)
	}
} 

func (z *AnyInt) Div(a, b AnyInt) {
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
			*z = AnyInt{Small:a.Small/b.Small}
		}
	} else {
		z.Int.Div(a.Int, b.Int)
	}
} 

func (z AnyInt) If() bool {
	if z.Int == nil {
		if z.Small != 0 {
			return true
		}
		return false
	}
	return z.Int.Cmp(big.NewInt(0)) != 0
} 

func (z AnyInt) Int64() int64 {
	if z.Int == nil {
		return z.Small
	}
	return z.Int.Int64()
} 


func BigInt(text string) AnyInt {
	n := AnyInt{Int:new(big.Int)}
	fmt.Sscan(text, n)
	return n
}

type String []AnyInt
type StringString [][]AnyInt

var N String
var N2 StringString

func (s *String) push(n AnyInt) {
	*s = append(*s, n)
}

func (p *String) pop() (n AnyInt) {
	s := *p
	n = s[len(s)-1]
	s = s[:len(s)-1]
	*p = s
	return
}

func (s *StringString) push(n String) {
	*s = append(*s, n)
}

func (p *StringString) pop() (n String) {
	s := *p
	n = s[len(s)-1]
	s = s[:len(s)-1]
	*p = s
	return
}

func pushstring(n String) {
	N2.push(n)
}

func popstring() (n String) {
	return N2.pop()
}

func push(n AnyInt) {
	if n.Int == nil {
		N.push(n)
	} else {
		N.push(AnyInt{Int:big.NewInt(0).Set(n.Int)})
	}
}

func pop() (n AnyInt) {
	return N.pop()
}

func stdout() {
	text := popstring()
	for _, v := range text {
		if v.Int != nil {
			os.Stdout.Write(v.Bytes())
		} else {
			os.Stdout.Write([]byte{byte(v.Small)})
		}
	} 
}

func stdin() {
	length := pop()
	var b []byte = make([]byte, 1)
	if length.Int != nil {
		for i := big.NewInt(0); i.Cmp(length.Int) < 0; i.Add(i, big.NewInt(1)) {	
			var n int
			for n == 0 {
				n, _ = os.Stdin.Read(b)
			}
			push(AnyInt{ Small:int64(b[0]) })
		}
	} else {
		for i := int64(0); i < length.Small; i++ {	
			var n int
			for n == 0 {
				n, _ = os.Stdin.Read(b)
			}
			push(AnyInt{Small:int64(b[0])})
		}
	}
}

func __bool2int(b bool) int64 {
	if b {
		return 1
	}
	return 0
}

func slt(a AnyInt, b AnyInt) AnyInt {
	if ca, cb := a.Int == nil, b.Int == nil; ca || cb {
		if !ca {
			return AnyInt{Small:__bool2int( a.Cmp(big.NewInt(b.Small)) == -1 )}
		} else if !cb {
			return AnyInt{Small:__bool2int( big.NewInt(a.Small).Cmp(b.Int) == -1 )}
		} else {
			return AnyInt{Small:__bool2int( a.Small < b.Small )}
		}
	} else {
		return AnyInt{Small:__bool2int( a.Cmp(b.Int) == -1 )}
	}
}

func seq(a AnyInt, b AnyInt) AnyInt {
	if ca, cb := a.Int == nil, b.Int == nil; ca || cb {
		if !ca {
			return AnyInt{Small:__bool2int( a.Cmp(big.NewInt(b.Small)) == 0 )}
		} else if !cb {
			return AnyInt{Small:__bool2int( big.NewInt(a.Small).Cmp(b.Int) == 0 )}
		} else {
			return AnyInt{Small:__bool2int( a.Small == b.Small )}
		}
	} else {
		return AnyInt{Small:__bool2int( a.Cmp(b.Int) == 0 )}
	}
}

func sge(a AnyInt, b AnyInt) AnyInt {
	if ca, cb := a.Int == nil, b.Int == nil; ca || cb {
		if !ca {
			return AnyInt{Small:__bool2int( a.Cmp(big.NewInt(b.Small)) >= 0 )}
		} else if !cb {
			return AnyInt{Small:__bool2int( big.NewInt(a.Small).Cmp(b.Int) >= 0 )}
		} else {
			return AnyInt{Small:__bool2int( a.Small >= b.Small )}
		}
	} else {
		return AnyInt{Small:__bool2int( a.Cmp(b.Int) >= 0 )}
	}
}

func sle(a AnyInt, b AnyInt) AnyInt {
	if ca, cb := a.Int == nil, b.Int == nil; ca || cb {
		if !ca {
			return AnyInt{Small:__bool2int( a.Cmp(big.NewInt(b.Small)) <= 0 )}
		} else if !cb {
			return AnyInt{Small:__bool2int( big.NewInt(a.Small).Cmp(b.Int) <= 0 )}
		} else {
			return AnyInt{Small:__bool2int( a.Small <= b.Small )}
		}
	} else {
		return AnyInt{Small:__bool2int( a.Cmp(b.Int) <= 0 )}
	}
}

func sgt(a AnyInt, b AnyInt) AnyInt {
	if ca, cb := a.Int == nil, b.Int == nil; ca || cb {
		if !ca {
			return AnyInt{Small:__bool2int( a.Cmp(big.NewInt(b.Small)) == 1 )}
		} else if !cb {
			return AnyInt{Small:__bool2int( big.NewInt(a.Small).Cmp(b.Int) == 1 )}
		} else {
			return AnyInt{Small:__bool2int( a.Small > b.Small )}
		}
	} else {
		return AnyInt{Small:__bool2int( a.Cmp(b.Int) == 1 )}
	}
}

func sne(a AnyInt, b AnyInt) AnyInt {
	if ca, cb := a.Int == nil, b.Int == nil; ca || cb {
		if !ca {
			return AnyInt{Small:__bool2int( a.Cmp(big.NewInt(b.Small)) != 0 )}
		} else if !cb {
			return AnyInt{Small:__bool2int( big.NewInt(a.Small).Cmp(b.Int) != 0 )}
		} else {
			return AnyInt{Small:__bool2int( a.Small != b.Small )}
		}
	} else {
		return AnyInt{Small:__bool2int( a.Cmp(b.Int) != 0 )}
	}
}
`)
}

func (g *GoAssembler) indt(n ...int) string {
	if len(n) > 0 {
		return strings.Repeat("\t", int(g.Indentation)+n[0])
	} else {
		return strings.Repeat("\t", int(g.Indentation))
	}
}

func (g *GoAssembler) Assemble(command string, args []string) ([]byte, error) {
	//Format argument expressions.
	//var expression bool
	for i, arg := range args {
		if arg == "=" {
			//expression = true
			continue
		}
		if _, err := strconv.Atoi(arg); err == nil {
			args[i] = "AnyInt{Small:"+arg+"}"
			continue
		} else if _, err := strconv.Atoi(string(arg[0])); err == nil {
			args[i] = "BigInt(\""+arg+"\")"
			continue
		} else if arg[0] == '-' {
			args[i] = "BigInt(\""+arg+"\")"
			continue
		}
		if arg[0] == '#' {
			args[i] = "AnyInt{Small:int64(len("+arg[1:]+"))}"
		}
		if arg[0] == '"' {
			var newarg string
			var j = i
			arg = arg[1:]
			
			stringloop:
			arg = strings.Replace(arg, "\\n", "\n", -1)
			for _, v := range arg {
				if v == '"' {
					goto end
				}
				newarg += "AnyInt{Small:"+strconv.Itoa(int(v))+"},"
			}
			if len(arg) == 0 {
				goto end
			}
			newarg += "AnyInt{Small:"+strconv.Itoa(int(' '))+"},"
			j++
			//println(arg)
			arg = args[j]
			goto stringloop
			end:
			//println(newarg)
			args[i] = newarg
		}
		switch arg {
			case "byte":
				args[i] = "u_"+args[i]
		}
	} 
	switch command {
		case "ROUTINE":
			defer func() { g.Indentation++ }()
			return []byte("func main() {\n"), nil
		case "SUBROUTINE":
			defer func() { g.Indentation++ }()
			return []byte("func "+args[0]+"() {\n"), nil
		case "PUSH", "PUSHSTRING":
			var name string
			if command == "PUSHSTRING" {
				name = "string"
			}
			if len(args) == 1 {
				return []byte(g.indt()+"push"+name+"("+args[0]+")\n"), nil
			} else {
				return []byte(g.indt()+args[1]+".push("+args[0]+")\n"), nil
			}
		case "POP", "POPSTRING":
			var name string
			if command == "POPSTRING" {
				name = "string"
			}
			if len(args) == 1 {
				return []byte(g.indt()+args[0]+" := pop"+name+"()\n"), nil
			} else {
				return []byte(g.indt()+args[0]+" := "+args[1]+".pop()\n"), nil
			}
		case "INDEX":
			return []byte(g.indt()+args[2]+" := "+args[0]+"["+args[1]+".Int64()]\n"), nil
		case "VAR":
			if len(args) == 1 {
				return []byte(g.indt()+"var "+args[0]+" AnyInt \n"), nil
			} else {
				return []byte(g.indt()+"var "+args[0]+" = "+args[1]+" \n"), nil
			}
		case "STRING":
			return []byte(g.indt()+"var "+args[0]+" String\n"), nil
		case "STDOUT":
			return []byte(g.indt()+"stdout()\n"), nil
		case "STDIN":
			return []byte(g.indt()+"stdin()\n"), nil
		case "LOOP":
			defer func() { g.Indentation++ }()
			return []byte(g.indt()+"for {\n"), nil
		case "REPEAT", "END", "DONE":
			g.Indentation--
			return []byte(g.indt()+"}\n"), nil
		case "IF":
			defer func() { g.Indentation++ }()
			return []byte(g.indt()+"if "+args[0]+".If() {\n"), nil
		case "RUN":
			return []byte(g.indt()+args[0]+"()\n"), nil
		case "ELSE":
			return []byte(g.indt(-1)+"} else {\n"), nil
		case "ELSEIF":
			return []byte(g.indt(-1)+"} else if "+args[0]+".If() {\n"), nil
		case "BREAK":
			return []byte(g.indt()+"break\n"), nil
		case "RETURN":
			return []byte(g.indt()+"return\n"), nil
		case "STRINGDATA":
			return []byte(g.indt()+"var "+args[0]+" String = String{"+args[1]+"} \n"), nil
		case "JOIN":
			return []byte(g.indt()+args[0]+" = append("+args[1]+", "+args[2]+"...) \n"), nil
		
		//Maths.
		case "ADD":
			return []byte(g.indt()+args[0]+".Add("+args[1]+","+args[2]+")\n"), nil
		case "SUB":
			return []byte(g.indt()+args[0]+".Sub("+args[1]+","+args[2]+")\n"), nil
		case "MUL":
			return []byte(g.indt()+args[0]+".Mul("+args[1]+","+args[2]+")\n"), nil
		case "DIV":
			return []byte(g.indt()+args[0]+".Div("+args[1]+","+args[2]+")\n"), nil
			
		case "SLT", "SEQ", "SGE", "SGT", "SNE", "SLE":
			return []byte(g.indt()+args[0]+" = "+strings.ToLower(command)+"("+args[1]+","+args[2]+")\n"), nil
		default:
			/*if expression {
				return []byte(g.indt()+command+" "+strings.Join(args, " ")+"\n"), nil
			}*/
	}
	return nil, errors.New("Unrecognised expression: "+command+" "+strings.Join(args, " "))
}

func (g *GoAssembler) Footer() []byte {
	return []byte("")
}
