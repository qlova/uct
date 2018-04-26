package compiler

import "bytes"
import "text/scanner"

type Cache struct {
	bytes.Buffer
}

func (c *Compiler) NewCache(open, close string) Cache {
	var cache Cache
	
	var depth = 1
	for {
		c.Scanners[len(c.Scanners)-1].Scan()
		tok := c.Scanners[len(c.Scanners)-1].TokenText()
		
		if tok == close {
			depth--
			if depth == 0 {
				break
			}
		} else if tok == open {
			depth++
		}
		
		cache.Write([]byte(tok))
	}
	
	return cache
}


func (c *Compiler) CompileCache(cache Cache, name string, line int) {
	
	
	var s scanner.Scanner
	s.Init(&cache)
	s.Filename = name
	s.Whitespace = 0

	c.Scanners = append(c.Scanners, &s)
	c.CurrentLines = append(c.CurrentLines, "")
	c.LineOffset = line
	
	var length = len(c.Scanners)
	
	for {
		c.ScanStatement()
		
		if len(c.Scanners) < length {
			c.LineOffset = 0
			return
		}
	}
}

func (c *Compiler) LoadCache(cache Cache, name string, line int) {
	var s scanner.Scanner
	s.Init(&cache)
	s.Filename = name
	s.Whitespace = 0

	c.Scanners = append(c.Scanners, &s)
	c.CurrentLines = append(c.CurrentLines, "")
	c.LineOffset = line
}
