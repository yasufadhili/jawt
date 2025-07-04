grammar JML;


source
    : jmlDocument
    | tsSource
    ;

jmlDocument
    : jmlDoctypeDeclaration
    ;

jmlDoctypeDeclaration
    : '_doctype' ('page' | 'component')
    ;

jmlImports
    :  jmlImportDeclaration+
    ;

jmlImportDeclaration
    : 'import' jmlImportType identifier 'from' StringLiteral
    ;

jmlImportType
    : 'component'
    | 'script'
    ;

jmlContent
    : jmlPageContent
    | jmlComponentContent
    ;

jmlPageContent
    : jmlElement
    ;

jmlComponentContent
    :  jmlElement
    ;

jmlElement
    : identifier '{' jmlElementBody '}'
    ;

jmlElementBody
    :  (jmlProperty | jmlElement | jmlControlStructure)*
    ;

jmlProperty
    :   identifier ':' jmlValue
    ;

jmlValue
    : StringLiteral
    | NumericLiteral
    | BooleanLiteral
    | identifier
    ;

jmlControlStructure
    :   jmlForLoop
    |   jmlIfStatement
    |   jmlSwitchStatement
    |   jmlPropsAccess
    ;

jmlForLoop
    : 'for' '(' identifier 'in' expression ')' '{' jmlElementBody '}'
    ;

jmlIfStatement
    : 'if' '(' expression ')' '{' jmlElementBody '}'
    ;

jmlElseClause
    : 'else' '{' jmlElementBody '}'
    | 'else' jmlIfStatement
    ;

jmlSwitchStatement
    : 'switch' '(' expression ')' '{' jmlCaseClause* jmlDefaultClause? '}'
    ;

jmlCaseClause
    :  'case' expression ':' jmlElementBody
    ;

jmlDefaultClause
    : 'default' ':' jmlElementBody
    ;

jmlPropsAccess
    : 'props' '.' identifier ('.' identifier)*
    ;






// TypeScript

tsSource
    : tsSourceElements? EOF
    ;

tsSourceElements
    : tsSourceElement+
    ;

tsSourceElement
    : importStatement
    | exportStatement
    | statement
    | declaration
    ;

importStatement
    : 'import' importClause 'from' StringLiteral eos
    | 'import' StringLiteral eos
    | 'import' '=' identifier StringLiteral eos
    ;

importClause
    : importedDefaultBinding (',' namedImports)?
    | namedImports
    | namespaceImport
    ;

importedDefaultBinding
    : identifier
    ;

namedImports
    : '{' (importSpecifier (',' importSpecifier)*)? '}'
    ;

importSpecifier
    : identifier ('as' identifier)?
    ;

namespaceImport
    : '*' 'as' identifier
    ;

exportStatement
    : 'export' '*' 'from' StringLiteral eos
    | 'export' exportClause ('from' StringLiteral)? eos
    | 'export' declaration
    | 'export' 'default' (assignmentExpression | declaration) eos
    ;

exportClause
    : '{' (exportSpecifier (',' exportSpecifier)*)? '}'
    ;

exportSpecifier
    : identifier ('as' identifier)?
    ;


declaration
    :  variableDeclaration
    ;

variableDeclaration
    : bindingPattern typeAnnotation? initializer? eos
    | 'var' variableDeclarationList eos
    | 'let' variableDeclarationList eos
    | 'const' variableDeclarationList eos
    ;

variableDeclarationList
    : variableBinding (',' variableBinding)*
    ;

variableBinding
    : bindingPattern typeAnnotation? initializer?
    ;

bindingPattern
    : identifier
    | arrayBindingPattern
    | objectBindingPattern
    ;

arrayBindingPattern
    : '[' (bindingElement (',' bindingElement)*)? ']'
    ;

bindingElement
    : bindingPattern initializer?
    | '...' bindingPattern
    ;

objectBindingPattern
    : '{' (bindingProperty (',' bindingProperty)*)? '}'
    ;

bindingProperty
    : identifier typeAnnotation? initializer?
    | propertyName ':' bindingPattern initializer?
    ;

initializer
    : '=' assignmentExpression
    ;

// Functions

functionDeclaration
    : 'async'? 'function' identifier callSignature ('{' functionBody '}' | eos)
    ;

callSignature
    : typeParameters? '(' parameterList? ')' typeAnnotation?
    ;

parameterList
    : restParameter
    | parameter (',' parameter)* (',' restParameter)?
    ;

parameter
    : accessibilityModifier? bindingPattern ('?' | initializer)? typeAnnotation?
    ;

restParameter
    : '...' bindingPattern typeAnnotation?
    ;

functionBody
    : tsSourceElements?
    ;


// Class

classDeclaration
    : 'abstract'? 'class' identifier typeParameters? classHeritage? '{' classBody '}'
    ;

classHeritage
    : extendsClause? implementsClause?
    ;

extendsClause
    : 'extends' typeReference
    ;

implementsClause
    : 'implements' typeReference (',' typeReference)*
    ;

classBody
    : classMember*
    ;

classMember
    : constructorDeclaration
    | propertyMemberDeclaration
    | methodDeclaration
    | accessorDeclaration
    | indexMemberDeclaration
    ;

constructorDeclaration
    : accessibilityModifier? 'constructor' '(' parameterList? ')' '{' functionBody '}'
    ;

propertyMemberDeclaration
    : accessibilityModifier? 'static'? 'readonly'? propertyName ('?' | '!')? typeAnnotation? initializer? eos
    ;

methodDeclaration
    : accessibilityModifier? 'static'? 'async'? propertyName '?'? callSignature ('{' functionBody '}' | eos)
    ;

accessorDeclaration
    : accessibilityModifier? 'static'? ('get' | 'set') propertyName '(' parameterList? ')' typeAnnotation? '{' functionBody '}'
    ;

indexMemberDeclaration
    : '[' bindingPattern ':' type ']' typeAnnotation eos
    ;


// Interface Declaration
interfaceDeclaration
    : 'interface' identifier typeParameters? interfaceExtendsClause? objectType
    ;

interfaceExtendsClause
    : 'extends' typeReference (',' typeReference)*
    ;

// Type Alias Declaration
typeAliasDeclaration
    : 'type' identifier typeParameters? '=' type eos
    ;

// Enum Declaration
enumDeclaration
    : 'const'? 'enum' identifier '{' enumBody? '}'
    ;

enumBody
    : enumMember (',' enumMember)* ','?
    ;

enumMember
    : propertyName ('=' assignmentExpression)?
    ;


// Namespace Declaration
namespaceDeclaration
    : 'namespace' identifier '{' tsSourceElements? '}'
    | 'module' identifier '{' tsSourceElements? '}'
    ;

// Ambient Declarations
ambientDeclaration
    : 'declare' ambientDeclarationElement
    ;

ambientDeclarationElement
    : variableDeclaration
    | functionDeclaration
    | classDeclaration
    | interfaceDeclaration
    | typeAliasDeclaration
    | enumDeclaration
    | namespaceDeclaration
    ;



// Statements
statement
    : expressionStatement
    | ifStatement
    | iterationStatement
    | continueStatement
    | breakStatement
    | returnStatement
    | withStatement
    | labelledStatement
    | switchStatement
    | throwStatement
    | tryStatement
    | debuggerStatement
    | emptyStatement
    | block
    ;

block
    : '{' tsSourceElements? '}'
    ;

expressionStatement
    : expression eos
    ;

ifStatement
    : 'if' '(' expression ')' statement ('else' statement)?
    ;

iterationStatement
    : 'do' statement 'while' '(' expression ')' eos
    | 'while' '(' expression ')' statement
    | 'for' '(' (variableDeclarationList | expression)? ';' expression? ';' expression? ')' statement
    | 'for' '(' (variableDeclaration | assignmentExpression) ('in' | 'of') expression ')' statement
    ;

continueStatement
    : 'continue' identifier? eos
    ;

breakStatement
    : 'break' identifier? eos
    ;

returnStatement
    : 'return' expression? eos
    ;

withStatement
    : 'with' '(' expression ')' statement
    ;

labelledStatement
    : identifier ':' statement
    ;

switchStatement
    : 'switch' '(' expression ')' caseBlock
    ;

caseBlock
    : '{' caseClauses? (defaultClause caseClauses?)? '}'
    ;

caseClauses
    : caseClause+
    ;

caseClause
    : 'case' expression ':' tsSourceElements?
    ;

defaultClause
    : 'default' ':' tsSourceElements?
    ;

throwStatement
    : 'throw' expression eos
    ;

tryStatement
    : 'try' block (catchClause finallyClause? | finallyClause)
    ;

catchClause
    : 'catch' ('(' bindingPattern typeAnnotation? ')')? block
    ;

finallyClause
    : 'finally' block
    ;

debuggerStatement
    : 'debugger' eos
    ;

emptyStatement
    : eos
    ;

// Expressions

expression
    : assignmentExpression (',' assignmentExpression)*
    ;

assignmentExpression
    :
    ;

assignmentOperator
    : '=' | '*=' | '/=' | '%=' | '+=' | '-=' | '<<=' | '>>=' | '>>>=' | '&=' | '^=' | '|='
    ;

conditionalExpression
    : logicalOrExpression ('?' assignmentExpression ':' assignmentExpression)?
    ;

logicalOrExpression
    : logicalAndExpression ('||' logicalAndExpression)*
    ;

logicalAndExpression
    : bitwiseOrExpression ('&&' bitwiseOrExpression)*
    ;

bitwiseOrExpression
    : bitwiseXorExpression ('|' bitwiseXorExpression)*
    ;

bitwiseXorExpression
    : bitwiseAndExpression ('^' bitwiseAndExpression)*
    ;

bitwiseAndExpression
    : equalityExpression ('&' equalityExpression)*
    ;

equalityExpression
    : relationalExpression (('==' | '!=' | '===' | '!==') relationalExpression)*
    ;

relationalExpression
    : shiftExpression (('<' | '>' | '<=' | '>=' | 'instanceof' | 'in') shiftExpression)*
    ;

shiftExpression
    : additiveExpression (('<<' | '>>' | '>>>') additiveExpression)*
    ;

additiveExpression
    : multiplicativeExpression (('+' | '-') multiplicativeExpression)*
    ;

multiplicativeExpression
    : exponentiationExpression (('*' | '/' | '%') exponentiationExpression)*
    ;

exponentiationExpression
    : unaryExpression ('**' exponentiationExpression)?
    ;

unaryExpression
    : postfixExpression
    | 'delete' unaryExpression
    | 'void' unaryExpression
    | 'typeof' unaryExpression
    | '++' unaryExpression
    | '--' unaryExpression
    | '+' unaryExpression
    | '-' unaryExpression
    | '~' unaryExpression
    | '!' unaryExpression
    | 'await' unaryExpression
    ;

postfixExpression
    : leftHandSideExpression
    | leftHandSideExpression '++'
    | leftHandSideExpression '--'
    ;

leftHandSideExpression
    : newExpression
    | callExpression
    ;

newExpression
    : memberExpression
    | 'new' newExpression
    ;

memberExpression
    : primaryExpression
    | memberExpression '[' expression ']'
    | memberExpression '.' identifier
    //| memberExpression templateLiteral
    | 'new' memberExpression arguments
    | 'super' '[' expression ']'
    | 'super' '.' identifier
    ;

callExpression
    : memberExpression arguments
    | 'super' arguments
    | callExpression arguments
    | callExpression '[' expression ']'
    | callExpression '.' identifier
    // | callExpression templateLiteral
    ;

arguments
    : '(' (argumentList ','?)? ')'
    ;

argumentList
    : assignmentExpression (',' assignmentExpression)*
    ;

primaryExpression
    : 'this'
    | identifier
    | literal
    | arrayLiteral
    | objectLiteral
    | functionExpression
    | arrowFunctionExpression
    | classExpression
    | '(' expression ')'
    | typeAssertion
    ;

// Literals
literal
    : BooleanLiteral
    | NumericLiteral
    | StringLiteral
    | NoSubstitutionTemplateLiteral
    | templateLiteral
    | RegularExpressionLiteral
    | 'null'
    ;

templateLiteral
    : TemplateHead (expression TemplateMiddle)* expression TemplateTail
    ;

arrayLiteral
    : '[' (elementList ','?)? ']'
    ;

elementList
    : assignmentExpression (',' assignmentExpression)*
    | '...' assignmentExpression
    ;

objectLiteral
    : '{' (propertyDefinitionList ','?)? '}'
    ;

propertyDefinitionList
    : propertyDefinition (',' propertyDefinition)*
    ;

propertyDefinition
    : propertyName ':' assignmentExpression
    | 'get' propertyName '(' ')' typeAnnotation? '{' functionBody '}'
    | 'set' propertyName '(' bindingPattern typeAnnotation? ')' '{' functionBody '}'
    | propertyName '?' '(' parameterList? ')' typeAnnotation? '{' functionBody '}'
    | identifier
    | '...' assignmentExpression
    ;

propertyName
    : identifier
    | StringLiteral
    | NumericLiteral
    | '[' assignmentExpression ']'
    ;

functionExpression
    : 'async'? 'function' identifier? callSignature '{' functionBody '}'
    ;

arrowFunctionExpression
    : 'async'? arrowFormalParameters '=>' (assignmentExpression | '{' functionBody '}')
    ;

arrowFormalParameters
    : identifier
    | callSignature
    ;

classExpression
    : 'class' identifier? typeParameters? classHeritage? '{' classBody '}'
    ;


// Type System
typeAssertion
    : '<' type '>' unaryExpression
    | unaryExpression 'as' type
    ;

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
    : parenthesizedType
    | predefinedType
    | typeReference
    | objectType
    | arrayType
    | tupleType
    | functionType
    | constructorType
    | typeQuery
    | thisType
    | literalType
    ;

parenthesizedType
    : '(' type ')'
    ;

predefinedType
    : 'any'
    | 'number'
    | 'boolean'
    | 'string'
    | 'symbol'
    | 'void'
    | 'undefined'
    | 'null'
    | 'never'
    | 'object'
    | 'unknown'
    ;

typeReference
    : typeName typeArguments?
    ;

typeName
    : identifier
    | namespaceName '.' identifier
    ;

namespaceName
    : identifier
    | namespaceName '.' identifier
    ;

typeArguments
    : '<' typeArgumentList '>'
    ;

typeArgumentList
    : type (',' type)*
    ;

objectType
    : '{' typeBody? '}'
    ;

typeBody
    : typeMember (';' | ',') (typeMember (';' | ','))* (';' | ',')?
    | typeMember (';' | ',')?
    ;

typeMember
    : propertySignature
    | callSignature
    | constructSignature
    | indexSignature
    | methodSignature
    ;

propertySignature
    : propertyName '?'? typeAnnotation?
    ;

callSignature
    : typeParameters? '(' parameterList? ')' typeAnnotation?
    ;

constructSignature
    : 'new' typeParameters? '(' parameterList? ')' typeAnnotation?
    ;

indexSignature
    : '[' bindingPattern ':' ('string' | 'number') ']' typeAnnotation
    ;

methodSignature
    : propertyName '?'? callSignature
    ;

arrayType
    : primaryType '[' ']'
    ;

tupleType
    : '[' tupleElementTypes? ']'
    ;

tupleElementTypes
    : type (',' type)*
    ;

functionType
    : typeParameters? '(' parameterList? ')' '=>' type
    ;

constructorType
    : 'new' typeParameters? '(' parameterList? ')' '=>' type
    ;

typeQuery
    : 'typeof' typeQueryExpression
    ;

typeQueryExpression
    : identifier
    | typeQueryExpression '.' identifier
    ;

thisType
    : 'this'
    ;

literalType
    : BooleanLiteral
    | NumericLiteral
    | StringLiteral
    ;

typeParameters
    : '<' typeParameterList '>'
    ;

typeParameterList
    : typeParameter (',' typeParameter)*
    ;

typeParameter
    : identifier constraint? defaultType?
    ;

constraint
    : 'extends' type
    ;

defaultType
    : '=' type
    ;

// Modifiers
accessibilityModifier
    : 'public'
    | 'private'
    | 'protected'
    ;

// End of statement
eos
    : ';'
    | EOF
    | {this.lineTerminatorAhead()}?
    | {this.closeBrace()}?
    ;

// Lexer Rules
// Keywords
Abstract: 'abstract';
Any: 'any';
As: 'as';
Async: 'async';
Await: 'await';
Boolean: 'boolean';
Break: 'break';
Case: 'case';
Catch: 'catch';
Class: 'class';
Const: 'const';
Constructor: 'constructor';
Continue: 'continue';
Debugger: 'debugger';
Declare: 'declare';
Default: 'default';
Delete: 'delete';
Do: 'do';
Else: 'else';
Enum: 'enum';
Export: 'export';
Extends: 'extends';
False: 'false';
Finally: 'finally';
For: 'for';
From: 'from';
Function: 'function';
Get: 'get';
If: 'if';
Implements: 'implements';
Import: 'import';
In: 'in';
Instanceof: 'instanceof';
Interface: 'interface';
Is: 'is';
Keyof: 'keyof';
Let: 'let';
Module: 'module';
Namespace: 'namespace';
Never: 'never';
New: 'new';
Null: 'null';
Number: 'number';
Object: 'object';
Of: 'of';
Package: 'package';
Private: 'private';
Protected: 'protected';
Public: 'public';
Readonly: 'readonly';
Return: 'return';
Set: 'set';
Static: 'static';
String: 'string';
Super: 'super';
Switch: 'switch';
Symbol: 'symbol';
This: 'this';
Throw: 'throw';
True: 'true';
Try: 'try';
Type: 'type';
Typeof: 'typeof';
Undefined: 'undefined';
Unknown: 'unknown';
Var: 'var';
Void: 'void';
While: 'while';
With: 'with';
Yield: 'yield';

// Literals
BooleanLiteral
    : False
    | True
    ;

NumericLiteral
    : DecimalLiteral
    | HexIntegerLiteral
    | OctalIntegerLiteral
    | OctalIntegerLiteral2
    | BinaryIntegerLiteral
    ;

DecimalLiteral
    : DecimalIntegerLiteral '.' DecimalDigits? ExponentPart?
    | '.' DecimalDigits ExponentPart?
    | DecimalIntegerLiteral ExponentPart?
    ;

DecimalIntegerLiteral
    : '0'
    | NonZeroDigit DecimalDigits?
    ;

DecimalDigits
    : DecimalDigit+
    ;

DecimalDigit
    : [0-9]
    ;

NonZeroDigit
    : [1-9]
    ;

ExponentPart
    : ExponentIndicator SignedInteger
    ;

ExponentIndicator
    : [eE]
    ;

SignedInteger
    : DecimalDigits
    | '+' DecimalDigits
    | '-' DecimalDigits
    ;

HexIntegerLiteral
    : '0' [xX] HexDigit+
    ;

HexDigit
    : [0-9a-fA-F]
    ;

OctalIntegerLiteral
    : '0' OctalDigit+
    ;

OctalIntegerLiteral2
    : '0' [oO] OctalDigit+
    ;

OctalDigit
    : [0-7]
    ;

BinaryIntegerLiteral
    : '0' [bB] BinaryDigit+
    ;

BinaryDigit
    : [01]
    ;

StringLiteral
    : '"' DoubleStringCharacter* '"'
    | '\'' SingleStringCharacter* '\''
    ;

DoubleStringCharacter
    : ~["\\\r\n]
    | '\\' EscapeSequence
    | LineContinuation
    ;

SingleStringCharacter
    : ~['\\\r\n]
    | '\\' EscapeSequence
    | LineContinuation
    ;

EscapeSequence
    : CharacterEscapeSequence
    | OctalEscapeSequence
    | HexEscapeSequence
    | UnicodeEscapeSequence
    ;

CharacterEscapeSequence
    : SingleEscapeCharacter
    | NonEscapeCharacter
    ;

SingleEscapeCharacter
    : ['"\\bfnrtv]
    ;

NonEscapeCharacter
    : ~['"\\bfnrtv0-9xu\r\n]
    ;

HexEscapeSequence
    : 'x' HexDigit HexDigit
    ;

OctalEscapeSequence
    : '0' OctalDigit? OctalDigit? // To handle \0, \07, \077
    ;

UnicodeEscapeSequence
    : 'u' HexDigit HexDigit HexDigit HexDigit
    | 'u{' HexDigit+ '}'
    ;

LineContinuation
    : '\\' [\r\n\u2028\u2029]
    ;

// Template Literals
TemplateStringLiteral
    : '`' ('\\`' | ~'`')* '`'
    ;

NoSubstitutionTemplateLiteral
    : '`' TemplateCharacters? '`'
    ;

TemplateHead
    : '`' TemplateCharacters? '${'
    ;

TemplateMiddle
    : '}' TemplateCharacters? '${'
    ;

TemplateTail
    : '}' TemplateCharacters? '`'
    ;

TemplateCharacters
    : TemplateCharacter+
    ;

TemplateCharacter
    : '$' ~'{'
    | '\\' EscapeSequence
    | '\\' ~[`\\$]
    | LineContinuation
    | LineTerminatorSequence
    | ~[`\\$\r\n\u2028\u2029]
    ;


// Regular Expression
RegularExpressionLiteral
    : '/' RegularExpressionBody '/' RegularExpressionFlags?
    ;

RegularExpressionBody
    : RegularExpressionFirstChar RegularExpressionChar*
    ;

RegularExpressionFirstChar
    : ~[*/[\r\n\u2028\u2029]
    | RegularExpressionBackslashSequence
    | RegularExpressionClass
    ;

RegularExpressionChar
    : ~[/[\r\n\u2028\u2029]
    | RegularExpressionBackslashSequence
    | RegularExpressionClass
    ;

RegularExpressionBackslashSequence
    : '\\' ~[\r\n\u2028\u2029]
    ;

RegularExpressionClass
    : '[' RegularExpressionClassChar* ']'
    ;

RegularExpressionClassChar
    : ~[\]\\r\n\u2028\u2029]
    | RegularExpressionBackslashSequence
    ;

RegularExpressionFlags
    : [gimsuxy]+
    ;

// Identifiers
identifier
    : Identifier
    | 'async'
    ;

Identifier
    : IdentifierStart IdentifierPart*
    ;

IdentifierStart
    : UnicodeLetter
    | '$'
    | '_'
    | '\\' UnicodeEscapeSequence
    ;

IdentifierPart
    : IdentifierStart
    | UnicodeDigit
    | UnicodeConnectorPunctuation
    | '\\' UnicodeEscapeSequence
    ;

UnicodeLetter
    : [a-zA-Z]
    ;

UnicodeDigit
    : [0-9]
    ;

fragment UnicodeConnectorPunctuation
    : '_'
    ;

// Whitespace and Comments
LineTerminator
    : [\r\n\u2028\u2029] -> skip
    ;

LineTerminatorSequence
    : '\r\n'
    | [\r\n\u2028\u2029]
    ;

WhiteSpaces
    : [\t\u000B\u000C\u0020\u00A0]+ -> skip
    ;

MultiLineComment
    : '/*' .*? '*/' -> skip
    ;

SingleLineComment
    : '//' ~[\r\n\u2028\u2029]* -> skip
    ;

// Punctuation
OpenBracket: '[';
CloseBracket: ']';
OpenParen: '(';
CloseParen: ')';
OpenBrace: '{';
CloseBrace: '}';
SemiColon: ';';
Comma: ',';
Assign: '=';
QuestionMark: '?';
Colon: ':';
Ellipsis: '...';
Dot: '.';
PlusPlus: '++';
MinusMinus: '--';
Plus: '+';
Minus: '-';
BitNot: '~';
Not: '!';
Multiply: '*';
Divide: '/';
Modulus: '%';
RightShiftArithmetic: '>>';
LeftShiftArithmetic: '<<';
RightShiftLogical: '>>>';
LessThan: '<';
MoreThan: '>';
LessThanEquals: '<=';
GreaterThanEquals: '>=';
Equals_: '==';
NotEquals: '!=';
IdentityEquals: '===';
IdentityNotEquals: '!==';
BitAnd: '&';
BitXOr: '^';
BitOr: '|';
And: '&&';
Or: '||';
MultiplyAssign: '*=';
DivideAssign: '/=';
ModulusAssign: '%=';
PlusAssign: '+=';
MinusAssign: '-=';
LeftShiftArithmeticAssign: '<<=';
RightShiftArithmeticAssign: '>>=';
RightShiftLogicalAssign: '>>>=';
BitAndAssign: '&=';
BitXorAssign: '^=';
BitOrAssign: '|=';
PowerAssign: '**=';
Arrow: '=>';
NullCoalescing: '??';