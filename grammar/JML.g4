grammar JML;


literal
    : INTEGER
;


INTEGER: [0-9]+ ;

IDENTIFIER: [a-zA-Z_][a-zA-Z0-9_]* ;

WHITESPACE  : [ \t\r\n]+ -> skip ;
