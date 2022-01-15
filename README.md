# gc

## Grammars

```
StatementList  = { Statement ";" } .
Statement      = ReturnStmt | Block | ForStmt | SimpleStmt .
ReturnStmt     = "return" [ ExpressionList ] .
Block          = "{" StatementList "}" .
ExpressionList = Expression { "," Expression } .
SimpleStmt     = ExpressionStmt | Assignment .
ExpressionStmt = Expression .
Assignment     = Expression assign_op Expression .

ForStmt    = "for" [ Condition | ForClause ] Block .
Condition  = Expression .
ForClause  = [ InitStmt ] ";" [ Condition ] ";" [ PostStmt ] .
InitStmt   = SimpleStmt .
PostStmt   = SimpleStmt .

Expression = UnaryExpr | Expression binary_op Expression .
UnaryExpr  = PrimaryExpr | unary_op UnaryExpr .
assign_op  = "=" .
binary_op  = rel_op | add_op | mul_op .
rel_op     = "==" | "!=" | "<" | "<=" | ">" | ">=" .
add_op     = "+" | "-" .
mul_op     = "*" | "/" .

unary_op   = "+" | "-" .
```
