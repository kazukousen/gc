# gc

## Grammars

```
Type = int

Declaration   = VarDecl .
TopLevelDecl = FunctionDecl .

VarDecl = "var" ( VarSpec | "(" { VarSpec ";" } ")" ) .
VarSpec = IdentifierList ( Type [ "=" ExpressionList ] | "=" ExpressionList ) .

FunctionDecl = "func" FunctionName Signature [ FunctionBody ] .
FunctionName = identifier .
FunctionBody = Block .

StatementList  = { Statement ";" } .
Statement      = Declaration | ReturnStmt | Block | IfStmt | ForStmt | SimpleStmt .
ReturnStmt     = "return" [ ExpressionList ] .
Block          = "{" StatementList "}" .
ExpressionStmt = Expression .

SimpleStmt     = ExpressionStmt | ShortVarDecl | Assignment .
ShortVarDecl   = IdentifierList ":=" ExpressionList .
Assignment     = ExpressionList assign_op ExpressionList .

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
