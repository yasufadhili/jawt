grammar JMLPage;

import JML;


program
    : doctypeSpecifier imports? page NEWLINE* EOF
;


page
    : 'Page' '{' NEWLINE pageBody NEWLINE? '}'
;


pageBody
    : pageProperty*
;


pageProperty
    : IDENTIFIER ':' propertyValue NEWLINE
;


propertyValue
    : literal
;

