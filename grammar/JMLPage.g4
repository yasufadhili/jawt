grammar JMLPage;

import JML;


page
    : doctypeSpecifier imports? pageDefinition NEWLINE* EOF
;


pageDefinition
    : 'Page' '{' NEWLINE pageBody NEWLINE? '}'
;


pageBody
    : (pageProperty | componentInvocation)*
;


pageProperty
    : IDENTIFIER ':' propertyValue NEWLINE
;


componentInvocation
    : COMP_ID '{' componentBody '}' NEWLINE*
    | COMP_ID NEWLINE
;


componentBody
    : (componentProperty | componentInvocation)*
;


componentProperty
    : IDENTIFIER ':' propertyValue NEWLINE
;


propertyValue
    : literal
    | componentInvocation
;


literal
    : INTEGER
    | STRING
    | IDENTIFIER
;