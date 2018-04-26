 package compiler
 
 type Syntax struct {
	Name string
	 
	Statements []Statement
	Expressions []Expression
	Types []Type
	Functions []*Function
	
	Operators []Operator
	
	Aliases map[string]string
 }
 
 func NewSyntax(name string) Syntax {
	 var s Syntax
	 s.Name = name
	 s.Aliases = make(map[string]string)
	 return s
 }

 func (s *Syntax) RegisterStatement(statement Statement) {
	s.Statements = append(s.Statements, statement)
 }
 
 func (s *Syntax) RegisterExpression(expression Expression) {
	s.Expressions = append(s.Expressions, expression)
 }
 
  func (s *Syntax) RegisterOperator(symbol string, p int) {
	s.Operators = append(s.Operators, Operator{
		Symbol: symbol,
		Precedence: p,
	})
 }
 
 func (s *Syntax) RegisterType(t Type) {
	s.Types = append(s.Types, t)
 }
 
  func (s *Syntax) RegisterFunction(f *Function) {
	s.Functions = append(s.Functions, f)
 }
 
   func (s *Syntax) RegisterAlias(from string, to string) {
	s.Aliases[from] = to
 }

 func (c *Compiler) SetSyntax(s Syntax) {
	c.Syntax = s
 }
