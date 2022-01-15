# gc

## Grammars

```
StatementList  = { Statement ";" } .
Statement      = ReturnStmt | Block | IfStmt | ForStmt | SimpleStmt .
ReturnStmt     = "return" [ ExpressionList ] .
Block          = "{" StatementList "}" .
ExpressionList = Expression { "," Expression } .
SimpleStmt     = ExpressionStmt | Assignment .
ExpressionStmt = Expression .
Assignment     = Expression assign_op Expression .

IfStmt = "if" [ SimpleStmt ";" ] Expression Block [ "else" ( IfStmt | Block ) ] .

ForStmt    = "for" [ Condition | ForClause ] Block .
Condition  = Expression .
ForClause  = [ InitStmt ] ";" [ Condition ] ";" [ PostStmt ] .
InitStmt   = SimpleStmt .
PostStmt   = SimpleStmt .

Expression  = UnaryExpr | Expression binary_op Expression .
UnaryExpr   = PrimaryExpr | unary_op UnaryExpr .

assign_op   = "=" .
binary_op   = rel_op | add_op | mul_op .
rel_op      = "==" | "!=" | "<" | "<=" | ">" | ">=" .
add_op      = "+" | "-" .
mul_op      = "*" | "/" .
unary_op   = "+" | "-" .

PrimaryExpr = Operand .

Operand     = Literal | OperandName [ Arguments ] | "(" Expression ")" .
Literal     = BasicLit .
BasicLit    = int_lit .

OperandName = identifier .
Arguments   = "(" [ ExpressionList [ "..." ] [ "," ] ] ")" .
```
