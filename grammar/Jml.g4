grammar Jml;

// Root rule - a JML document must start with a doctype declaration
document
    : doctypeDeclaration imports? sourceElements? EOF
    ;

// Doctype declaration - required at the beginning of every JML document
doctypeDeclaration
    : '_doctype' doctypeKind doctypeName
    ;

doctypeKind
    : 'page'
    | 'component'
    ;

doctypeName
    : Identifier
    ;

// Import declarations - specific to JML with different import types
imports
    : importDeclaration+
    ;

importDeclaration
    : componentImport
    | scriptImport
    | browserImport
    | moduleImport
    ;

componentImport
    : 'import' 'component' Identifier 'from' StringLiteral ';'?
    ;

scriptImport
    : 'import' 'script' Identifier 'from' StringLiteral ';'?
    ;

browserImport
    : 'import' 'browser' ';'?
    ;

// Traditional module imports (for compatibility)
moduleImport
    : 'import' importClause 'from' StringLiteral ';'?
    | 'import' StringLiteral ';'?
    ;

importClause
    : Identifier
    | namedImports
    | Identifier ',' namedImports
    ;

namedImports
    : '{' importsList? '}'
    ;

importsList
    : importSpecifier (',' importSpecifier)*
    ;

importSpecifier
    : Identifier
    | Identifier 'as' Identifier
    ;

// Source elements - the main content after doctype and imports
sourceElements
    : pageContent      // For pages - single root element
    | componentContent // For components - can have multiple elements
    ;

// Page content - must have exactly one root element
pageContent
    : componentElement
    ;

// Component content - can have multiple top-level elements
componentContent
    : sourceElement+
    ;

sourceElement
    : statement
    | declaration
    | exportDeclaration
    | componentElement
    ;

// Component element - JML's declarative syntax for UI components
componentElement
    : Identifier '{' componentBody? '}'
    ;

componentBody
    : componentBodyElement+
    ;

componentBodyElement
    : componentProperty
    | componentElement
    | statement
    | forLoop
    | ifCondition
    ;

componentProperty
    : Identifier ':' expression ';'?
    ;

// Control structures within components
forLoop
    : 'for' '(' Identifier 'in' expression ')' '{' componentBody? '}'
    ;

ifCondition
    : 'if' '(' expression ')' '{' componentBody? '}' ('else' '{' componentBody? '}')?
    ;

// Export declarations
exportDeclaration
    : 'export' declaration
    | 'export' 'default' (declaration | expression) ';'?
    | 'export' '{' exportsList '}' ';'?
    | 'export' '{' exportsList '}' 'from' StringLiteral ';'?
    ;

exportsList
    : exportSpecifier (',' exportSpecifier)*
    ;

exportSpecifier
    : Identifier
    | Identifier 'as' Identifier
    ;

// Declarations
declaration
    : variableDeclaration
    | functionDeclaration
    | classDeclaration
    | interfaceDeclaration
    | typeAliasDeclaration
    | enumDeclaration
    | propertyDeclaration
    | stateDeclaration
    ;

propertyDeclaration
    : 'property' typeAnnotation? Identifier (':' expression)? ';'?
    ;

stateDeclaration
    : 'state' Identifier '{' stateBody? '}' ';'?
    ;

stateBody
    : stateBodyElement+
    ;

stateBodyElement
    : propertyDeclaration
    | Identifier ':' expression ';'?
    ;

// Variable declarations
variableDeclaration
    : variableDeclarationList ';'?
    ;

variableDeclarationList
    : variableDeclarationKind variableDeclarator (',' variableDeclarator)*
    ;

variableDeclarationKind
    : 'var'
    | 'let'
    | 'const'
    ;

variableDeclarator
    : Identifier typeAnnotation? ('=' expression)?
    ;

// Function declarations
functionDeclaration
    : 'function' Identifier '(' parameterList? ')' typeAnnotation? functionBody
    ;

parameterList
    : parameter (',' parameter)*
    ;

parameter
    : Identifier typeAnnotation? ('=' expression)?
    | '...' Identifier typeAnnotation?
    ;

functionBody
    : '{' statementList? '}'
    ;

// Class declarations
classDeclaration
    : 'class' Identifier typeParameters? classHeritage? '{' classBody? '}'
    ;

classHeritage
    : 'extends' typeReference ('implements' typeReference (',' typeReference)*)?
    | 'implements' typeReference (',' typeReference)*
    ;

classBody
    : classMember+
    ;

classMember
    : methodDefinition
    | propertyDefinition
    | constructorDefinition
    | propertyDeclaration
    | stateDeclaration
    ;

constructorDefinition
    : 'constructor' '(' parameterList? ')' functionBody
    ;

methodDefinition
    : accessibilityModifier? 'static'? Identifier '(' parameterList? ')' typeAnnotation? functionBody
    ;

propertyDefinition
    : accessibilityModifier? 'static'? Identifier typeAnnotation? ('=' expression)? ';'?
    ;

accessibilityModifier
    : 'public'
    | 'private'
    | 'protected'
    ;

// Interface declarations
interfaceDeclaration
    : 'interface' Identifier typeParameters? interfaceHeritage? '{' interfaceBody? '}'
    ;

interfaceHeritage
    : 'extends' typeReference (',' typeReference)*
    ;

interfaceBody
    : interfaceMember+
    ;

interfaceMember
    : propertySignature
    | methodSignature
    | indexSignature
    ;

propertySignature
    : Identifier '?'? typeAnnotation ';'?
    ;

methodSignature
    : Identifier '(' parameterList? ')' typeAnnotation? ';'?
    ;

indexSignature
    : '[' Identifier ':' type ']' ':' type ';'?
    ;

// Type alias declarations
typeAliasDeclaration
    : 'type' Identifier typeParameters? '=' type ';'?
    ;

// Enum declarations
enumDeclaration
    : 'enum' Identifier '{' enumBody? '}'
    ;

enumBody
    : enumMember (',' enumMember)* ','?
    ;

enumMember
    : Identifier ('=' expression)?
    ;

// Type annotations and types
typeAnnotation
    : ':' type
    ;

type
    : unionType
    ;

unionType
    : intersectionType ('|' intersectionType)*
    ;

intersectionType
    : primaryType ('&' primaryType)*
    ;

primaryType
    : baseType ('[' ']')*
    ;

baseType
    : typeReference
    | objectType
    | tupleType
    | primitiveType
    | literalType
    | '(' type ')'
    ;

primitiveType
    : 'string'
    | 'number'
    | 'boolean'
    | 'void'
    | 'any'
    | 'unknown'
    | 'never'
    | 'undefined'
    | 'null'
    ;

literalType
    : StringLiteral
    | NumericLiteral
    | 'true'
    | 'false'
    ;

typeReference
    : Identifier typeArguments?
    ;

typeArguments
    : '<' type (',' type)* '>'
    ;

typeParameters
    : '<' typeParameter (',' typeParameter)* '>'
    ;

typeParameter
    : Identifier ('extends' type)?
    ;

objectType
    : '{' objectTypeBody? '}'
    ;

objectTypeBody
    : objectTypeMember (',' objectTypeMember)* ','?
    ;

objectTypeMember
    : propertySignature
    | methodSignature
    | indexSignature
    ;

tupleType
    : '[' type (',' type)* ']'
    ;

// Statements
statement
    : block
    | expressionStatement
    | ifStatement
    | iterationStatement
    | returnStatement
    | breakStatement
    | continueStatement
    | throwStatement
    | tryStatement
    | ';'
    ;

statementList
    : statement+
    ;

block
    : '{' statementList? '}'
    ;

expressionStatement
    : expression ';'?
    ;

ifStatement
    : 'if' '(' expression ')' statement ('else' statement)?
    ;

iterationStatement
    : 'for' '(' (variableDeclarationList | expression)? ';' expression? ';' expression? ')' statement
    | 'for' '(' (variableDeclarationList | leftHandSideExpression) 'in' expression ')' statement
    | 'for' '(' (variableDeclarationList | leftHandSideExpression) 'of' expression ')' statement
    | 'while' '(' expression ')' statement
    | 'do' statement 'while' '(' expression ')' ';'?
    ;

returnStatement
    : 'return' expression? ';'?
    ;

breakStatement
    : 'break' ';'?
    ;

continueStatement
    : 'continue' ';'?
    ;

throwStatement
    : 'throw' expression ';'?
    ;

tryStatement
    : 'try' block catchClause? finallyClause?
    ;

catchClause
    : 'catch' '(' Identifier typeAnnotation? ')' block
    ;

finallyClause
    : 'finally' block
    ;

// Expressions
expression
    : conditionalExpression
    ;

conditionalExpression
    : logicalOrExpression ('?' expression ':' expression)?
    ;

logicalOrExpression
    : logicalAndExpression ('||' logicalAndExpression)*
    ;

logicalAndExpression
    : equalityExpression ('&&' equalityExpression)*
    ;

equalityExpression
    : relationalExpression (('==' | '!=' | '===' | '!==') relationalExpression)*
    ;

relationalExpression
    : additiveExpression (('<' | '>' | '<=' | '>=' | 'instanceof' | 'in') additiveExpression)*
    ;

additiveExpression
    : multiplicativeExpression (('+' | '-') multiplicativeExpression)*
    ;

multiplicativeExpression
    : unaryExpression (('*' | '/' | '%') unaryExpression)*
    ;

unaryExpression
    : postfixExpression
    | ('++' | '--' | '+' | '-' | '!' | '~' | 'typeof' | 'void' | 'delete') unaryExpression
    ;

postfixExpression
    : leftHandSideExpression ('++' | '--')?
    ;

leftHandSideExpression
    : newExpression
    | callExpression
    ;

newExpression
    : memberExpression
    | 'new' newExpression
    ;

callExpression
    : memberExpression arguments
    | callExpression arguments
    | callExpression '[' expression ']'
    | callExpression '.' Identifier
    ;

memberExpression
    : primaryExpression
    | memberExpression '[' expression ']'
    | memberExpression '.' Identifier
    | 'new' memberExpression arguments
    ;

arguments
    : '(' argumentList? ')'
    ;

argumentList
    : expression (',' expression)*
    ;

primaryExpression
    : 'this'
    | 'super'
    | Identifier
    | literal
    | arrayLiteral
    | objectLiteral
    | functionExpression
    | arrowFunction
    | '(' expression ')'
    | templateLiteral
    ;

literal
    : 'null'
    | 'undefined'
    | 'true'
    | 'false'
    | NumericLiteral
    | StringLiteral
    | RegexLiteral
    ;

arrayLiteral
    : '[' elementList? ']'
    ;

elementList
    : expression (',' expression)* ','?
    ;

objectLiteral
    : '{' propertyNameAndValueList? '}'
    ;

propertyNameAndValueList
    : propertyAssignment (',' propertyAssignment)* ','?
    ;

propertyAssignment
    : propertyName ':' expression
    | Identifier
    | '...' expression
    ;

propertyName
    : Identifier
    | StringLiteral
    | NumericLiteral
    | '[' expression ']'
    ;

functionExpression
    : 'function' Identifier? '(' parameterList? ')' typeAnnotation? functionBody
    ;

arrowFunction
    : arrowFunctionParameters '=>' (functionBody | expression)
    ;

arrowFunctionParameters
    : Identifier
    | '(' parameterList? ')'
    ;

templateLiteral
    : TemplateStringLiteral
    ;

// Lexer rules
Identifier
    : [a-zA-Z_$][a-zA-Z0-9_$]*
    ;

NumericLiteral
    : DecimalLiteral
    | HexIntegerLiteral
    | OctalIntegerLiteral
    | BinaryIntegerLiteral
    ;

DecimalLiteral
    : DecimalIntegerLiteral '.' [0-9]* ExponentPart?
    | '.' [0-9]+ ExponentPart?
    | DecimalIntegerLiteral ExponentPart?
    ;

DecimalIntegerLiteral
    : '0'
    | [1-9][0-9]*
    ;

ExponentPart
    : [eE] [+-]? [0-9]+
    ;

HexIntegerLiteral
    : '0' [xX] [0-9a-fA-F]+
    ;

OctalIntegerLiteral
    : '0' [0-7]+
    ;

BinaryIntegerLiteral
    : '0' [bB] [01]+
    ;

StringLiteral
    : '"' DoubleStringCharacter* '"'
    | '\'' SingleStringCharacter* '\''
    ;

TemplateStringLiteral
    : '`' TemplateStringCharacter* '`'
    ;

RegexLiteral
    : '/' RegexBody '/' RegexFlags?
    ;

fragment DoubleStringCharacter
    : ~["\\\r\n]
    | EscapeSequence
    ;

fragment SingleStringCharacter
    : ~['\\\r\n]
    | EscapeSequence
    ;

fragment TemplateStringCharacter
    : ~[`\\]
    | EscapeSequence
    | '${' .*? '}'
    ;

fragment EscapeSequence
    : '\\' ['"\\bfnrtv]
    | '\\' [0-3]? [0-7]? [0-7]
    | '\\' 'u' [0-9a-fA-F] [0-9a-fA-F] [0-9a-fA-F] [0-9a-fA-F]
    | '\\' 'x' [0-9a-fA-F] [0-9a-fA-F]
    ;

fragment RegexBody
    : RegexFirstChar RegexChar*
    ;

fragment RegexFirstChar
    : ~[*\r\n\\/[]
    | RegexBackslashSequence
    | RegexClass
    ;

fragment RegexChar
    : ~[\r\n\\/[]
    | RegexBackslashSequence
    | RegexClass
    ;

fragment RegexBackslashSequence
    : '\\' ~[\r\n]
    ;

fragment RegexClass
    : '[' RegexClassChar* ']'
    ;

fragment RegexClassChar
    : ~[\r\n\]\\]
    | RegexBackslashSequence
    ;

fragment RegexFlags
    : [gimuy]*
    ;

// Skip whitespace and comments
WS
    : [ \t\r\n\u000C]+ -> skip
    ;

LineComment
    : '//' ~[\r\n]* -> skip
    ;

BlockComment
    : '/*' .*? '*/' -> skip
    ;