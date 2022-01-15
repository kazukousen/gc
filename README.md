# gc

## Grammars

```
Block          = "{" StatementList "}" .
StatementList  = { Statement ";" } .
Statement      = ReturnStmt | SimpleStmt .
ReturnStmt     = "return" [ ExpressionList ] .
ExpressionList = Expression { "," Expression } .
SimpleStmt     = ExpressionStmt | Assignment .
ExpressionStmt = Expression .
Assignment     = Expression assign_op Expression .

Expression = UnaryExpr | Expression binary_op Expression .
UnaryExpr  = PrimaryExpr | unary_op UnaryExpr .
assign_op  = "=" .
binary_op  = rel_op | add_op | mul_op .
rel_op     = "==" | "!=" | "<" | "<=" | ">" | ">=" .
add_op     = "+" | "-" .
mul_op     = "*" | "/" .

unary_op   = "+" | "-" .
```
