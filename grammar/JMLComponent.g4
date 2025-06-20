grammar JMLComponent;

import JML;


component
    :   COMP_ID '{' componentBody '}' NEWLINE
;


componentInvocation
    :   COMP_ID '{' componentBody '}' NEWLINE
    |   COMP_ID NEWLINE
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
    | componentInvocation
;

