# gc

## Grammars

```
Type = int

TopLevelDecl = FunctionDecl .
FunctionDecl = "func" FunctionName Signature [ FunctionBody ] .
FunctionName = identifier .
FunctionBody = Block .

StatementList  = { Statement ";" } .
Statement      = ReturnStmt | Block | IfStmt | ForStmt | SimpleStmt .
ReturnStmt     = "return" [ ExpressionList ] .
Block          = "{" StatementList "}" .
SimpleStmt     = ExpressionStmt | Assignment .
ExpressionStmt = Expression .
Assignment     = Expression assign_op Expression .

Signature      = Parameters [ Type ] .
Parameters     = "(" [ ParameterList ] ")" .
ParameterList  = ParameterDecl { "," ParameterDecl } .
ParameterDecl  = [ IdentifierList ] Type .

IdentifierList = identifier { "," identifier } .
ExpressionList = Expression { "," Expression } .

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
unary_op   = "+" | "-" | "*" | "&" .

PrimaryExpr = Operand .

Operand     = Literal | OperandName [ Arguments ] | "(" Expression ")" .
Literal     = BasicLit .
BasicLit    = int_lit .

OperandName = identifier .
Arguments   = "(" [ ExpressionList [ "..." ] [ "," ] ] ")" .
```
