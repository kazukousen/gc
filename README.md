# gc

## Grammars

```
StatementList  = { Statement ";" } .
Statement      = ReturnStmt | ExpressionStmt .
ReturnStmt     = "return" [ ExpressionList ] .
ExpressionList = Expression { "," Expression } .
ExpressionStmt = Expression .

Expression = UnaryExpr | Expression binary_op Expression .
UnaryExpr  = PrimaryExpr | unary_op UnaryExpr .
binary_op  = rel_op | add_op | mul_op .
rel_op     = "==" | "!=" | "<" | "<=" | ">" | ">=" .
add_op     = "+" | "-" .
mul_op     = "*" | "/" .

unary_op   = "+" | "-" .
```
