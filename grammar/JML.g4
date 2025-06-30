grammar JML;


document
    : doctypeSpecifier imports? documentContent EOF
;


documentContent
    : pageDefinition           // For pages
    | componentDefinition      // For components
    | moduleDefinition         // For modules (future extension)
;


doctypeSpecifier
    : '_doctype' doctype IDENTIFIER
;


doctype
    : 'page'
    | 'component'
    | 'module'
;


// Import system
imports
    : importStatement+
;


importStatement
    : 'import' doctype IDENTIFIER 'from' STRING
    | 'import' IDENTIFIER 'from' STRING
    | 'import' IDENTIFIER
;


// Page-specific content (single child component constraint)
pageDefinition
    : 'Page' '{' pageBody '}'
;


pageBody
    : pageProperties? singleComponentChild?
;


pageProperties
    : pageProperty+
;


pageProperty
    : IDENTIFIER ':' propertyValue
;


// Pages can only have one child component
singleComponentChild
    : componentElement
;


// Component-specific content (can have multiple children)
componentDefinition
    : componentElement
;


componentElement
    : IDENTIFIER componentBlock?
;


componentBlock
    : '{' componentBody '}'
;


componentBody
    : (componentProperty | componentElement | scriptFunction)*
;


componentProperty
    : IDENTIFIER ':' propertyValue
;


// Property values can be literals, expressions, or nested components
propertyValue
    : literal
    | expression
    | componentElement
    | arrayLiteral
    | objectLiteral
;


expression
    : IDENTIFIER                                                 #IdExpr
    | functionCall                                               #FuncCallExpr
    | memberAccess                                               #MemberAccessExpr
    | literal                                                    #LiteralExpr
    | arrayLiteral                                               #ArrayExpr
    | objectLiteral                                              #ObjectExpr
    | '(' expression ')'                                         #ParenExpr
    | expression binaryOperator expression                       #BinaryExpr
    | expression '?' expression ':' expression                   #ConditionalExpr
;


functionCall
    : IDENTIFIER '(' argumentList? ')'
;


argumentList
    : expression (',' expression)*
;


memberAccess
    : IDENTIFIER ('.' IDENTIFIER)+
;


binaryOperator
    : '+'
    | '-'
    | '*'
    | '/'
    | '=='
    | '!='
    | '<'
    | '>'
    | '<='
    | '>='
    | '&&'
    | '||'
;


// Array and object literals for complex properties
arrayLiteral
    : '[' (propertyValue (',' propertyValue)*)? ']'
;


objectLiteral
    : '{' (objectProperty (',' objectProperty)*)? '}'
;


objectProperty
    : (IDENTIFIER | STRING) ':' propertyValue
;


// Script functions (for components)
scriptFunction
    : functionDeclaration
;


functionDeclaration
    : 'function' IDENTIFIER '(' parameterList? ')' ':' typeAnnotation? '{' functionBody '}'
;


parameterList
    : parameter (',' parameter)*
;


parameter
    : IDENTIFIER ':' typeAnnotation
;


typeAnnotation
    : 'void'
    | 'string'
    | 'number'
    | 'boolean'
    | 'any'
    | IDENTIFIER
;


functionBody
    : statement*
;


statement
    : expressionStatement
    | returnStatement
    | ifStatement
    | variableDeclaration
;


expressionStatement
    : expression
;


returnStatement
    : 'return' expression?
;


ifStatement
    : 'if' '(' expression ')' '{' statement* '}' ('else' '{' statement* '}')?
;


variableDeclaration
    : ('let' | 'const' | 'var') IDENTIFIER ':' typeAnnotation? '=' expression
;


// Module content (placeholder for future extension)
moduleDefinition
    : moduleFunction+
;


moduleFunction
    : 'export' functionDeclaration
;


literal
    : INTEGER
    | FLOAT
    | STRING
    | BOOLEAN
    | NULL
;


// Lexer rules
INTEGER
    : [0-9]+
;


FLOAT
    : [0-9]+ '.' [0-9]+
;


BOOLEAN
    : 'true'
    | 'false'
;


NULL
    : 'null'
;


IDENTIFIER
    : [a-zA-Z_][a-zA-Z0-9_]*
;


STRING
    : '"' (~["\r\n\\] | '\\' .)* '"'
    | '\'' (~['\r\n\\] | '\\' .)* '\''
;


// Special handling for template literals (for dynamic content)
TEMPLATE_LITERAL
    : '`' (~[`\\] | '\\' . | '${' ~[}]* '}')* '`'
;


WHITESPACE
    : [ \r\n\t]+ -> skip
;


COMMENT
    : '//' ~[\r\n]* -> skip
;


MULTILINE_COMMENT
    : '/*' .*? '*/' -> skip
;


// Operators and punctuation
LPAREN : '(' ;
RPAREN : ')' ;
LBRACE : '{' ;
RBRACE : '}' ;
LBRACKET : '[' ;
RBRACKET : ']' ;
SEMICOLON : ';' ;
COMMA : ',' ;
DOT : '.' ;
COLON : ':' ;
QUESTION : '?' ;
PLUS : '+' ;
MINUS : '-' ;
MULTIPLY : '*' ;
DIVIDE : '/' ;
EQUALS : '=' ;
DOUBLE_EQUALS : '==' ;
NOT_EQUALS : '!=' ;
LESS_THAN : '<' ;
GREATER_THAN : '>' ;
LESS_EQUALS : '<=' ;
GREATER_EQUALS : '>=' ;
LOGICAL_AND : '&&' ;
LOGICAL_OR : '||' ;
LOGICAL_NOT : '!' ;