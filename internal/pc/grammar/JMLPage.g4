grammar JMLPage;



program
    : doctype imports? page EOF
;

doctype
    : '_doctype' 'page' IDENTIFIER
;

imports
    : importStatement+
;

importStatement
    : 'import' component 'from' STRING
;

component
    : IDENTIFIER
;

page
    : 'Page' '{' pageBody '}'
;

pageBody
    : pageProperty* componentInvocation?
;

pageProperty
    : IDENTIFIER ':' value
;

componentInvocation
    : IDENTIFIER '{' componentBody '}'
;

componentBody
    : componentProperty*
;

componentProperty
    : IDENTIFIER ':' value
;

value
    : STRING
    | IDENTIFIER
    | componentInvocation
;




IDENTIFIER
    : [a-zA-Z_][a-zA-Z0-9_]*
;

STRING
    : '"' (~["\r\n] | '\\"')* '"'
;

COMMENT
    : '//' ~[\r\n]* -> skip
;

MULTILINE_COMMENT
    : '/*' .*? '*/' -> skip
;

WS
    : [ \t\r\n]+ -> skip
;



