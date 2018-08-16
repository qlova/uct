package compiler

import "fmt"
import "strings"
import "runtime"

func (c *Compiler) RaiseError(message Translatable) {
	c.Language = English
	
	fmt.Fprint(c.StdErr[len(c.StdErr)-1], c.Scanners[len(c.Scanners)-1].Line+c.LineOffset, ": ")
	var x = c.Scanners[len(c.Scanners)-1].Column+1
	
	if len(c.Token()) > 1 {
		x -= len(c.Token())/2
	}

	if c.NextToken != "" && !c.DoneToken {
		x--
	}
	
	//var line = c.CurrentLines[len(c.Scanners)-1]
	for tok := c.Scan(); tok != "\n"; tok = c.Scan() {
		if tok != "\n" {
			//line = c.CurrentLines[len(c.Scanners)-1]
		}
	}
	fmt.Fprintln(c.StdErr[len(c.StdErr)-1], c.LastLine)
	
	fmt.Fprintln(c.StdErr[len(c.StdErr)-1], strings.Repeat(" ", x), "^")
	fmt.Fprintln(c.StdErr[len(c.StdErr)-1], message[c.Language])
	panic("error")
}

func (c *Compiler) Expecting(token string) {
	if c.Scan() != token {
		c.RaiseError(Translatable{
			English: "Expecting "+token,
		})
	}
}

func (c *Compiler) ExpectingType(t Type) {
	if b := c.ScanExpression(); !b.Equals(t) {
		c.RaiseError(Translatable{
			English: "Expecting this to be of type "+t.Name[c.Language]+", however, it is of type "+b.Name[c.Language],
		})
	}
}

func (c *Compiler) Unexpected(token string) {
	c.RaiseError(Translatable{
		English: "Unexpected "+token,
	})
}

func (c *Compiler) Unimplemented() {
	
	pc, _, _, _ := runtime.Caller(1)
	name := runtime.FuncForPC(pc).Name()
	_, line := runtime.FuncForPC(pc).FileLine(pc)
	
	c.RaiseError(Translatable{
		English: "What you are trying to do is valid, however, \n it is not yet implemented by this compiler!\n("+name+":"+fmt.Sprint(line)+")",
	})
}

func (c *Compiler) UndefinedError(token string) {
	c.RaiseError(Translatable{
		English: "Undefined: "+token,
	})
}


func (c *Compiler) Expected(tokens ...string) {
	
	var symbols string 
	
	for i, symbol := range tokens {
		symbols += "'"+symbol+"'"
		if i < len(symbol)-1 {
			symbols += ","
		}
	}
	
	c.RaiseError(Translatable{
		English: "Unexpected "+c.Token()+", Expected "+symbols,
	})
}

