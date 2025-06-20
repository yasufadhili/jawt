grammar JMLPage;

import JML;


program
    : doctypeSpecifier imports? page EOF
;


page
    : 'Page' '{' '}'
;


pageBody
    : pageProperty*
;


pageProperty
    : IDENTIFIER ':' STRING
;

