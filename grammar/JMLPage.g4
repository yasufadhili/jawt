grammar JMLPage;

import JML;


program
    : doctypeSpecifier imports? page EOF
;


page
    : 'Page' '{' pageBody '}'
;


pageBody
    : pageProperty*
;


pageProperty
    : IDENTIFIER ':' propertyValue
;


propertyValue
    : literal
;

