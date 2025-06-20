grammar JMLComponent;

import JML;

program
    :
;


component
    :   COMP_ID '{' componentBody '}' NEWLINE
;


componentInvocation
    :   COMP_ID '{' componentBody '}' NEWLINE
;


componentBody
    :   componentProperty*
;


componentProperty
    :   IDENTIFIER ':' propertyValue NEWLINE
;


propertyValue
    : STRING
    | IDENTIFIER
    | component
;


componentContent
    :
;


