grammar JML;


doctypeSpecifier
    : '_doctype' doctype IDENTIFIER NEWLINE
;


doctype
    : 'page'
    | 'component'
    | 'module'
;


imports
    : importStatement+
;


importStatement
    : 'import' doctype IDENTIFIER 'from' STRING NEWLINE
;


literal
    : INTEGER
    | STRING
    | IDENTIFIER
;


INTEGER
    : [0-9]+
;

COMP_ID
    : [A-Z][a-zA-Z0-9_]*
;

IDENTIFIER
    : [a-zA-Z_][a-zA-Z0-9_]*
;

STRING
    : '"' (~["\r\n] | '\\"')* '"'
;

NEWLINE
    : ('\r'? '\n')+
;

WHITESPACE
    : [ \t]+ -> skip
;

COMMENT
    : '//' ~[\r\n]* -> skip
;

MULTILINE_COMMENT
    : '/*' .*? '*/' -> skip
;
