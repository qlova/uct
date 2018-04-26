package compiler

import "io"
import "os"
import "text/scanner"
import "fmt"
import "runtime/debug"

type Compiler struct {
	Syntax
	
	CurrentFunction *Function
	
	Language Language
	
	Scope []Scope
	
	Header io.Writer
	Output io.Writer

	StdErr []io.Writer
	
	DisableOutput bool
	
	
	Scanners []*scanner.Scanner

	CurrentLines []string
	LastLine string

	LineOffset int
	NextToken string
	DoneToken bool
	token string
	
	Errors bool
}

func (c *Compiler) SwapOutput() {
	c.Header, c.Output = c.Output, c.Header
}

func (c *Compiler) Token() string {
	return c.token
}

func (c *Compiler) Peek() string {
	if c.NextToken != "" && !c.DoneToken {
		return c.NextToken
	}
	
	var token = c.Scanners[len(c.Scanners)-1].TokenText()
	c.NextToken = c.Scan()
	c.token = token
	return c.NextToken
}

func (c *Compiler) scan() string {
	if c.NextToken != "" && !c.DoneToken  {
		c.DoneToken = true
		c.token = c.NextToken
		return c.NextToken
	}
	if c.DoneToken {
		c.NextToken = ""
		c.DoneToken = false
	}
	
	tok := c.Scanners[len(c.Scanners)-1].Scan()
	
	for c.Scanners[len(c.Scanners)-1].TokenText() == " " || c.Scanners[len(c.Scanners)-1].TokenText() == "\t" {
		if c.Scanners[len(c.Scanners)-1].TokenText() != "\t" {
			c.CurrentLines[len(c.Scanners)-1] += c.Scanners[len(c.Scanners)-1].TokenText()
		}
		tok = c.Scanners[len(c.Scanners)-1].Scan()
	}
	
	if tok == scanner.EOF {
		c.CurrentLines = c.CurrentLines[:len(c.Scanners)-1]
		c.Scanners = c.Scanners[:len(c.Scanners)-1]
		return ""
	}
	
	c.CurrentLines[len(c.Scanners)-1] += c.Scanners[len(c.Scanners)-1].TokenText()
	
	if c.Scanners[len(c.Scanners)-1].TokenText() == "\n" {
		c.LastLine = c.CurrentLines[len(c.Scanners)-1][:len(c.CurrentLines[len(c.Scanners)-1])-1]
		c.CurrentLines[len(c.Scanners)-1] = ""
	}
	
	c.token = c.Scanners[len(c.Scanners)-1].TokenText()
	return c.Scanners[len(c.Scanners)-1].TokenText()
}

func (c *Compiler) Scan() string {
	var token = c.scan()
	
	if alias, ok := c.Aliases[token]; ok {
		return alias
	}
	
	return token
}

func (c *Compiler) ScanIf(test string) bool {
	if c.Peek() == test {
		c.Scan()
		return true
	}
	return false
}

func (c *Compiler) AddInput(input io.Reader) {
	var s = new(scanner.Scanner)
	s.Init(input)
	s.Whitespace = 0
	c.Scanners = append(c.Scanners, s)
	c.CurrentLines = append(c.CurrentLines, "")
}

func (c *Compiler) Compile() {
	
	defer func() {
        if r := recover(); r != nil {
			if r == "error" {
				c.Errors = true
				
				if os.Getenv("PANIC") == "1" { 
					panic("PANIC=1")
				}
				
			} else if r != "done" {
				fmt.Println(r, string(debug.Stack()))
			}
        }
    }()
	
	if len(c.Scanners) == 0 {
		return
	}
	
	for {
		c.ScanStatement()
		if len(c.Scanners) == 0 {
			return
		}
	}
}
