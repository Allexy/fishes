package tokenizer

// Possible token types
const (
	TT_BOF               TokenType = 0
	TT_DEFAULT           TokenType = 1  // ...
	TT_WORD              TokenType = 2  // abc
	TT_NUMBER            TokenType = 3  // 123.0
	TT_STRING            TokenType = 4  // "..."
	TT_LOGIC             TokenType = 5  // true / false
	TT_OPERATOR          TokenType = 6  // +-*/...
	TT_O_PAREN           TokenType = 7  // (
	TT_C_PAREN           TokenType = 8  // )
	TT_O_BRACKET         TokenType = 9  // [
	TT_C_BRACKET         TokenType = 10 // ]
	TT_O_BRACE           TokenType = 11 // {
	TT_C_BRACE           TokenType = 12 // }
	TT_COLON             TokenType = 13 // :
	TT_SEMICOLON         TokenType = 14 // ;
	TT_COMA              TokenType = 15 // ,
	TT_POINT             TokenType = 16 // .
	TT_AT                TokenType = 17 // @
	TT_ASSIGNMENT        TokenType = 18 // =
	TT_ARROW             TokenType = 19 // =>
	TT_VARIABLE          TokenType = 20 // $varname
	TT_COMMENT           TokenType = 21 // #.... or //...
	TT_MULTILINE_COMMENT TokenType = 22 // /*...*/
	TT_WHITE_SPACE       TokenType = 23 // any white space
	TT_EOF               TokenType = 24
)

// Key words
const (
	KW_TRUE  = "true"
	KW_FALSE = "false"
	KW_BIND  = "=>"
)

// Operators
const (
	OP_PLUS         = "+"
	OP_MINUS        = "-"
	OP_DIV          = "/"
	OP_MULT         = "*"
	OP_MOD          = "%"
	OP_ASSIGN       = "="
	OP_EQUALS       = "=="
	OP_NOT_EQUALS   = "!="
	OP_GT           = ">"
	OP_GTE          = ">="
	OP_LT           = "<"
	OP_LTE          = "<="
	OP_NOT          = "!"
	OP_AND          = "&&"
	OP_OR           = "||"
	OP_INC          = "++"
	OP_DEC          = "--"
	OP_PLUS_ASSIGN  = "+="
	OP_MINUS_ASSIGN = "-="
	OP_DIV_ASSIGN   = "/="
	OP_MULT_ASSIGN  = "*="
	OP_MOD_ASSIGN   = "%="
)
