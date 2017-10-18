package parser

import (
	"go/token"
	"testing"

	"github.com/nordsieck/defect"
)

type Tmap map[string][][]*Token
type Output struct {
	tree []string
	rest [][]*Token
}
type Omap map[string]Output

func GetStateOutput(ss []State) []StateOutput {
	var result []StateOutput
	for _, s := range ss {
		var r []string
		for _, rend := range s.r {
			r = append(r, string(rend.Render()))
		}
		result = append(result, StateOutput{r, s.t})
	}
	return result
}

type StateOutput struct {
	s []string
	t []*Token
}

var (
	empty = [][]*Token(nil)

	_else        = &Token{token.ELSE, `else`}
	_fallthrough = &Token{token.FALLTHROUGH, `fallthrough`}
	_if          = &Token{token.IF, `if`}
	_int         = &Token{token.IDENT, `int`}
	_var         = &Token{token.VAR, `var`}
	a            = &Token{token.IDENT, `a`}
	add          = &Token{tok: token.ADD}
	arrow        = &Token{tok: token.ARROW}
	assign       = &Token{tok: token.ASSIGN}
	b            = &Token{token.IDENT, `b`}
	c            = &Token{token.IDENT, `c`}
	colon        = &Token{tok: token.COLON}
	comma        = &Token{tok: token.COMMA}
	d            = &Token{token.IDENT, `d`}
	define       = &Token{tok: token.DEFINE}
	dot          = &Token{tok: token.PERIOD}
	ellipsis     = &Token{tok: token.ELLIPSIS}
	inc          = &Token{tok: token.INC}
	lbrace       = &Token{tok: token.LBRACE}
	lbrack       = &Token{tok: token.LBRACK}
	lparen       = &Token{tok: token.LPAREN}
	one          = &Token{token.INT, `1`}
	rbrace       = &Token{tok: token.RBRACE}
	rbrack       = &Token{tok: token.RBRACK}
	ret          = &Token{token.SEMICOLON, "\n"}
	rparen       = &Token{tok: token.RPAREN}
	semi         = &Token{token.SEMICOLON, `;`}
	two          = &Token{token.INT, `2`}
	zero         = &Token{token.INT, `0`}
)

func TestAddOp(t *testing.T) {
	resultState(t, AddOp, map[string][]StateOutput{
		`+`: {{[]string{``, `+`}, []*Token{}}},
		`-`: {{[]string{``, `-`}, []*Token{}}},
		`|`: {{[]string{``, `|`}, []*Token{}}},
		`&`: {{[]string{``, `&`}, []*Token{}}},
		`1`: nil,
	})
}

func TestAliasDeclState(t *testing.T) {
	resultState(t, AliasDeclState, map[string][]StateOutput{`a=b`: {min(`a=b`)}})
}

func TestAnonymousField(t *testing.T) {
	resultState(t, AnonymousField, map[string][]StateOutput{
		`a.a`:  {{[]string{``, `a.a`}, []*Token{ret}}, {[]string{``, `a`}, []*Token{ret, a, dot}}},
		`*a.a`: {{[]string{``, `*a.a`}, []*Token{ret}}, {[]string{``, `*a`}, []*Token{ret, a, dot}}},
	})
}

func TestAnonymousField_Render(t *testing.T) {
	defect.Equal(t, string(anonymousField{true, a}.Render()), `*a`)
	defect.Equal(t, string(anonymousField{false, &qualifiedIdent{a, b}}.Render()), `a.b`)
}

func TestArguments(t *testing.T) {
	resultState(t, Arguments, map[string][]StateOutput{
		`()`:        {min(`()`)},
		`(a)`:       {min(`(a)`)},
		`(a,a)`:     {min(`(a,a)`)},
		`(a,a,)`:    {min(`(a,a,)`)},
		`(a,a...,)`: {min(`(a,a...,)`)},
	})
}

func TestArguments_Render(t *testing.T) {
	defect.Equal(t, string(arguments{}.Render()), `()`)
	defect.Equal(t, string(arguments{expressionList: expressionList{a, b}}.Render()), `(a,b)`)
	defect.Equal(t, string(arguments{expressionList: expressionList{a, b}, ellipsis: true, comma: true}.Render()), `(a,b...,)`)
}

func TestArrayType(t *testing.T) {
	resultState(t, ArrayType, map[string][]StateOutput{
		`[1]int`:   {min(`[1]int`)},
		`[a]int`:   {min(`[a]int`)},
		`[a.a]int`: {min(`[a.a]int`)},
	})
}

func TestArrayType_Render(t *testing.T) {
	defect.Equal(t, string(arrayType{one, _int}.Render()), `[1]int`)
}

func TestAssignment(t *testing.T) {
	remaining(t, Assignment, Tmap{`a = 1`: {{ret}}})
}

func TestAssignmentState(t *testing.T) {
	resultState(t, AssignmentState, map[string][]StateOutput{`a=1`: {min(`a=1`)}})
}

func TestAssignOp(t *testing.T) {
	resultState(t, AssignOp, map[string][]StateOutput{
		`+=`:  {{[]string{``, `+=`}, []*Token{}}},
		`-=`:  {{[]string{``, `-=`}, []*Token{}}},
		`*=`:  {{[]string{``, `*=`}, []*Token{}}},
		`/=`:  {{[]string{``, `/=`}, []*Token{}}},
		`%=`:  {{[]string{``, `%=`}, []*Token{}}},
		`&=`:  {{[]string{``, `&=`}, []*Token{}}},
		`|=`:  {{[]string{``, `|=`}, []*Token{}}},
		`^=`:  {{[]string{``, `^=`}, []*Token{}}},
		`<<=`: {{[]string{``, `<<=`}, []*Token{}}},
		`>>=`: {{[]string{``, `>>=`}, []*Token{}}},
		`&^=`: {{[]string{``, `&^=`}, []*Token{}}},
		`=`:   {{[]string{``, `=`}, []*Token{}}},
	})
}

func TestBasicLit(t *testing.T) {
	resultState(t, BasicLit, map[string][]StateOutput{
		``:    nil,
		`1`:   {{[]string{``, `1`}, []*Token{ret}}},
		`1.1`: {{[]string{``, `1.1`}, []*Token{ret}}},
		`1i`:  {{[]string{``, `1i`}, []*Token{ret}}},
		`'a'`: {{[]string{``, `'a'`}, []*Token{ret}}},
		`"a"`: {{[]string{``, `"a"`}, []*Token{ret}}},
		`a`:   nil,
		`_`:   nil,
	})
}

func TestBinaryOp(t *testing.T) {
	resultState(t, BinaryOp, map[string][]StateOutput{
		`==`: {{[]string{``, `==`}, []*Token{}}},
		`+`:  {{[]string{``, `+`}, []*Token{}}},
		`*`:  {{[]string{``, `*`}, []*Token{}}},
		`||`: {{[]string{``, `||`}, []*Token{}}},
		`&&`: {{[]string{``, `&&`}, []*Token{}}},
		`1`:  nil,
	})
}

func TestBlock(t *testing.T) {
	remaining(t, Block, Tmap{
		`{}`:        {{ret}},
		`{a()}`:     {{ret}},
		`{a();}`:    {{ret}},
		`{a();b()}`: {{ret}},
	})
}

func TestBlockState(t *testing.T) {
	resultState(t, BlockState, map[string][]StateOutput{
		`{}`:        {min(`{}`)},
		`{;}`:       {min(`{;}`)},
		`{a++}`:     {min(`{a++}`)},
		`{a++;}`:    {min(`{a++;}`)},
		`{a++;b++}`: {min(`{a++;b++}`)},
	})
}

func TestBreakStmt(t *testing.T) {
	remaining(t, BreakStmt, Tmap{
		`break a`: {{ret, a}, {ret}},
		`a`:       empty,
	})
}

func TestBreakStmtState(t *testing.T) {
	resultState(t, BreakStmtState, map[string][]StateOutput{
		`a`:       nil,
		`break`:   {min(`break`)},
		`break a`: {{[]string{``, `break`}, []*Token{ret, a}}, min(`break a`)},
	})
}

func TestChannelType(t *testing.T) {
	resultState(t, ChannelType, map[string][]StateOutput{
		`chan int`:   {{[]string{``, `chan int`}, []*Token{ret}}},
		`<-chan int`: {{[]string{``, `<-chan int`}, []*Token{ret}}},
		`chan<- int`: {{[]string{``, `chan<- int`}, []*Token{ret}}},
		`int`:        nil,
	})
}

func TestChannelType_Render(t *testing.T) {
	defect.Equal(t, string(channelType{nil, &Token{token.IDENT, `int`}}.Render()), `chan int`)
	defect.Equal(t, string(channelType{&_true, &Token{token.IDENT, `int`}}.Render()), `<-chan int`)
	defect.Equal(t, string(channelType{&_false, &Token{token.IDENT, `int`}}.Render()), `chan<- int`)
}

func TestCompositeLit(t *testing.T) {
	remaining(t, CompositeLit, Tmap{
		`T{1}`: {{ret}},
		`T{foo: "bar", baz: "quux",}`: {{ret}, {ret}, {ret}, {ret}},
	})
}

func TestCompositeLitState(t *testing.T) {
	resultState(t, CompositeLitState, map[string][]StateOutput{
		`a{1}`: {min(`a{1}`)},
		`a{foo:"bar",baz:"quux",}`: {min(`a{foo:"bar",baz:"quux",}`)},
	})
}

func TestCommCase(t *testing.T) {
	remaining(t, CommCase, Tmap{
		`default`:   {{}},
		`case <-a`:  {{ret}},
		`case a<-5`: {{ret}, {ret, {token.INT, `5`}, arrow}},
	})
}

func TestCommClause(t *testing.T) {
	remaining(t, CommClause, Tmap{
		`default:`:     {{}},
		`case <-a:`:    {{}},
		`case a<-5:`:   {{}},
		`default: a()`: {{ret, rparen, lparen, a}, {ret, rparen, lparen}, {ret}, {}},
		`case <-a: b(); c()`: {
			{ret, rparen, lparen, c, semi, rparen, lparen, b},
			{ret, rparen, lparen, c, semi, rparen, lparen},
			{ret, rparen, lparen, c, semi},
			{ret, rparen, lparen, c},
			{ret, rparen, lparen}, {ret}, {}},
	})
}

func TestConstDecl(t *testing.T) {
	remaining(t, ConstDecl, Tmap{
		`const a = 1`:                         {{ret}},
		`const a int = 1`:                     {{ret}},
		`const (a = 1)`:                       {{ret}},
		`const (a int = 1)`:                   {{ret}},
		`const (a, b int = 1, 2; c int = 3;)`: {{ret}},
	})
}

func TestConstDeclState(t *testing.T) {
	resultState(t, ConstDeclState, map[string][]StateOutput{
		`const a=1`:                    {min(`const a=1`)},
		`const a int=1`:                {min(`const a int=1`)},
		`const ()`:                     {min(`const ()`)},
		`const (a=1)`:                  {min(`const (a=1)`)},
		`const (a=1;)`:                 {min(`const (a=1;)`)},
		`const (a=1;b=2)`:              {min(`const (a=1;b=2)`)},
		`const (a,b int=1,2;c int=3;)`: {min(`const (a,b int=1,2;c int=3;)`)},
	})
}

func TestConstDecl_Render(t *testing.T) {
	defect.Equal(t, string(constDecl{[]Renderer{constSpec{a, nil, one}}, false, false}.Render()), `const a=1`)
	defect.Equal(t, string(constDecl{nil, true, false}.Render()), `const ()`)
	defect.Equal(t, string(constDecl{[]Renderer{constSpec{a, nil, one}}, true, true}.Render()), `const (a=1;)`)
	defect.Equal(t, string(constDecl{[]Renderer{constSpec{a, nil, one}, constSpec{b, nil, two}}, true, false}.Render()), `const (a=1;b=2)`)
}

func TestConstSpec(t *testing.T) {
	remaining(t, ConstSpec, Tmap{
		`a = 1`:           {{ret}},
		`a int = 1`:       {{ret}},
		`a, b int = 1, 2`: {{ret, two, comma}, {ret}},
	})
}

func TestConstSpecState(t *testing.T) {
	resultState(t, ConstSpecState, map[string][]StateOutput{
		`a=1`:         {min(`a=1`)},
		`a int=1`:     {min(`a int=1`)},
		`a,b int=1,2`: {{[]string{``, `a,b int=1`}, []*Token{ret, two, comma}}, min(`a,b int=1,2`)},
	})
}

func TestConstSpec_Render(t *testing.T) {
	defect.Equal(t, string(constSpec{a, _int, one}.Render()), `a int=1`)
	defect.Equal(t, string(constSpec{a, nil, one}.Render()), `a=1`)
}

func TestContinueStmt(t *testing.T) {
	remaining(t, ContinueStmt, Tmap{
		`continue a`: {{ret, a}, {ret}},
		`a`:          empty,
	})
}

func TestContinueStmtState(t *testing.T) {
	resultState(t, ContinueStmtState, map[string][]StateOutput{
		`a`:          nil,
		`continue`:   {min(`continue`)},
		`continue a`: {{[]string{``, `continue`}, []*Token{ret, a}}, min(`continue a`)},
	})
}

func TestConversion(t *testing.T) {
	remaining(t, Conversion, Tmap{
		`float(1)`:   {{ret}},
		`(int)(5,)`:  {{ret}},
		`a.a("foo")`: {{ret}},
	})
}

func TestConversionState(t *testing.T) {
	resultState(t, ConversionState, map[string][]StateOutput{
		`float(1)`:   {min(`float(1)`)},
		`(int)(5,)`:  {min(`(int)(5,)`)},
		`a.a("foo")`: {min(`a.a("foo")`)},
	})
}

func TestConversion_Render(t *testing.T) {
	defect.Equal(t, string(conversion{_int, a, false}.Render()), `int(a)`)
	defect.Equal(t, string(conversion{_int, a, true}.Render()), `int(a,)`)
}

func TestDeclaration(t *testing.T) {
	remaining(t, Declaration, Tmap{
		`const a = 1`: {{ret}},
		`type a int`:  {{ret}},
		`var a int`:   {{ret}},
	})
}

func TestDeclarationState(t *testing.T) {
	resultState(t, DeclarationState, map[string][]StateOutput{
		`const a=1`:  {min(`const a=1`)},
		`type a int`: {min(`type a int`)},
		`var a int`:  {min(`var a int`)},
	})
}

func TestDeferStmt(t *testing.T) {
	remaining(t, DeferStmt, Tmap{
		`defer a()`: {{ret, rparen, lparen}, {ret}},
		`defer`:     empty,
		`a()`:       empty,
	})
}

func TestElement(t *testing.T) {
	remaining(t, Element, Tmap{
		`1+1`:   {{ret, one, add}, {ret}},
		`{1+1}`: {{ret}},
	})
}

func TestElementState(t *testing.T) {
	resultState(t, ElementState, map[string][]StateOutput{
		`1+1`:   {{[]string{``, `1`}, []*Token{ret, one, add}}, min(`1+1`)},
		`{1+1}`: {min(`{1+1}`)},
	})
}

func TestElementList(t *testing.T) {
	remaining(t, ElementList, Tmap{
		`1:1`: {{ret, one, colon}, {ret}},
		`1:1, 1:1`: {
			{ret, one, colon, one, comma, one, colon},
			{ret, one, colon, one, comma},
			{ret, one, colon}, {ret}},
	})
}

func TestElementListState(t *testing.T) {
	resultState(t, ElementListState, map[string][]StateOutput{
		`1:1`: {{[]string{``, `1`}, []*Token{ret, one, colon}}, min(`1:1`)},
		`1:1,1:1`: {
			{[]string{``, `1`}, []*Token{ret, one, colon, one, comma, one, colon}},
			{[]string{``, `1:1`}, []*Token{ret, one, colon, one, comma}},
			{[]string{``, `1:1,1`}, []*Token{ret, one, colon}},
			min(`1:1,1:1`)},
	})
}

func TestElementList_Render(t *testing.T) {
	defect.Equal(t, string(elementList{a}.Render()), `a`)
	defect.Equal(t, string(elementList{a, b}.Render()), `a,b`)
}

func TestEllipsisArrayType(t *testing.T) {
	resultState(t, EllipsisArrayType, map[string][]StateOutput{`[...]int`: {{[]string{``, `[...]int`}, []*Token{ret}}}})
}

func TestEllipsisArrayType_Render(t *testing.T) {
	defect.Equal(t, string(ellipsisArrayType{typ{_int, 0}}.Render()), `[...]int`)
}

func TestEmptyStmt(t *testing.T) {
	remaining(t, EmptyStmt, Tmap{`1`: {{ret, one}}})
}

func TestEmptyStmtState(t *testing.T) {
	resultState(t, EmptyStmtState, map[string][]StateOutput{`1`: {{[]string{``}, []*Token{ret, one}}}})
}

func TestExprCaseClause(t *testing.T) {
	remaining(t, ExprCaseClause, Tmap{
		`default: a()`: {{ret, rparen, lparen, a}, {ret, rparen, lparen}, {ret}, {}},
	})
}

func TestExpression(t *testing.T) {
	remaining(t, Expression, Tmap{
		`1`:    {{ret}},
		`1+1`:  {{ret, one, add}, {ret}},
		`1+-1`: {{ret, one, {tok: token.SUB}, add}, {ret}},
	})
}

func TestExpressionState(t *testing.T) {
	resultState(t, ExpressionState, map[string][]StateOutput{
		`1`:    {{[]string{``, `1`}, []*Token{ret}}},
		`1+1`:  {{[]string{``, `1`}, []*Token{ret, one, add}}, min(`1+1`)},
		`1+-1`: {{[]string{``, `1`}, []*Token{ret, one, {tok: token.SUB}, add}}, min(`1+-1`)},
	})
}

func TestExpression_Render(t *testing.T) {
	defect.Equal(t, string(expression{a}.Render()), `a`)
	defect.Equal(t, string(expression{a, add, b}.Render()), `a+b`)
}

func TestExpressionList(t *testing.T) {
	remaining(t, ExpressionList, Tmap{
		`1`:   {{ret}},
		`1,1`: {{ret, one, comma}, {ret}},
	})
}

func TestExpressionListState(t *testing.T) {
	resultState(t, ExpressionListState, map[string][]StateOutput{
		`1`:   {min(`1`)},
		`a,b`: {{[]string{``, `a`}, []*Token{ret, b, comma}}, min(`a,b`)},
	})
}

func TestExpressionList_Render(t *testing.T) {
	defect.Equal(t, string(expressionList{}.Render()), ``)
	defect.Equal(t, string(expressionList{a}.Render()), `a`)
	defect.Equal(t, string(expressionList{a, b, c}.Render()), `a,b,c`)
}

func TestExpressionStmt(t *testing.T) {
	remaining(t, ExpressionStmt, Tmap{`1`: {{ret}}})
}

func TestExpressionStmtState(t *testing.T) {
	resultState(t, ExpressionStmtState, map[string][]StateOutput{`1`: {min(`1`)}})
}

func TestExprSwitchCase(t *testing.T) {
	remaining(t, ExprSwitchCase, Tmap{
		`case 1`:    {{ret}},
		`case 1, 1`: {{ret, one, comma}, {ret}},
		`default`:   {{}},
	})
}

func TestExprSwitchStmt(t *testing.T) {
	remaining(t, ExprSwitchStmt, Tmap{
		`switch{}`:                                               {{ret}},
		`switch{default:}`:                                       {{ret}},
		`switch a := 1; {}`:                                      {{ret}},
		`switch a {}`:                                            {{ret}},
		`switch a := 1; a {}`:                                    {{ret}},
		`switch a := 1; a {default:}`:                            {{ret}},
		`switch {case true:}`:                                    {{ret}},
		`switch {case true: a()}`:                                {{ret}},
		`switch{ case true: a(); case false: b(); default: c()}`: {{ret}},
	})
}

func TestFallthroughStmt(t *testing.T) {
	resultState(t, FallthroughStmt, map[string][]StateOutput{
		`fallthrough`: {{[]string{``, `fallthrough`}, []*Token{ret}}},
		`a`:           nil,
	})
}

func TestFieldDecl(t *testing.T) {
	resultState(t, FieldDecl, map[string][]StateOutput{
		`a,a int`: {{[]string{``, `a,a int`}, []*Token{ret}}, {[]string{``, `a`}, []*Token{ret, _int, a, comma}}},
		`a,a int "foo"`: {{[]string{``, `a,a int`}, []*Token{ret, {token.STRING, `"foo"`}}}, {[]string{``, `a,a int "foo"`}, []*Token{ret}},
			{[]string{``, `a`}, []*Token{ret, {token.STRING, `"foo"`}, _int, a, comma}}},
		`*int`:       {{[]string{``, `*int`}, []*Token{ret}}},
		`*int "foo"`: {{[]string{``, `*int`}, []*Token{ret, {token.STRING, `"foo"`}}}, {[]string{``, `*int "foo"`}, []*Token{ret}}},
	})
}

func TestFieldDecl_Render(t *testing.T) {
	defect.Equal(t, string(fieldDecl{anonField: &Token{token.IDENT, `int`}, tag: &Token{token.STRING, `a`}}.Render()), `int a`)
	defect.Equal(t, string(fieldDecl{idList: identifierList{&Token{token.IDENT, `a`}, &Token{token.IDENT, `b`}}, typ: &Token{token.IDENT, `int`}}.Render()), `a,b int`)
}

func TestForClause(t *testing.T) {
	remaining(t, ForClause, Tmap{
		`;;`:       {{}, {}, {}, {}},
		`a := 1;;`: {{}, {}},
		`; a < 5;`: {{}, {}, {}, {}},
		`;; a++`: {{ret, inc, a}, {ret, inc, a}, {ret, inc, a},
			{ret, inc, a}, {ret, inc}, {ret, inc}, {ret}, {ret}},
		`a := 1; a < 5; a++`: {{ret, inc, a}, {ret, inc, a}, {ret, inc}, {ret}},
	})
}

func TestForStmt(t *testing.T) {
	remaining(t, ForStmt, Tmap{
		`for {}`:                               {{ret}},
		`for true {}`:                          {{ret}},
		`for a := 0; a < 5; a++ {}`:            {{ret}},
		`for range a {}`:                       {{ret}},
		`for { a() }`:                          {{ret}},
		`for { a(); }`:                         {{ret}},
		`for { a(); b() }`:                     {{ret}},
		`for a := 0; a < 5; a++ { a(); b(); }`: {{ret}},
	})
}

func TestFunction(t *testing.T) {
	remaining(t, Function, Tmap{
		`(){}`:                 {{ret}},
		`()(){}`:               {{ret}},
		`(int)int{ return 0 }`: {{ret}},
		`(){ a() }`:            {{ret}},
	})
}

func TestFunctionState(t *testing.T) {
	resultState(t, FunctionState, map[string][]StateOutput{
		`(){}`:             {min(`(){}`)},
		`()(){}`:           {min(`()(){}`)},
		`(int)int{return}`: {min(`(int)int{return}`)},
		`(){a()}`:          {min(`(){a()}`)},
	})
}

func TestFunctionDecl(t *testing.T) {
	remaining(t, FunctionDecl, Tmap{
		`func f(){}`:                        {{ret}},
		`func f()(){}`:                      {{ret}},
		`func f(int)int{ a(); return b() }`: {{ret}},
	})
}

func TestFunctionLit(t *testing.T) {
	remaining(t, FunctionLit, Tmap{`func(){}`: {{ret}}})
}

func TestFunctionLitState(t *testing.T) {
	resultState(t, FunctionLitState, map[string][]StateOutput{`func(){}`: {min(`func(){}`)}})
}

func TestFunctionType(t *testing.T) {
	resultState(t, FunctionType, map[string][]StateOutput{
		`func()`:             {{[]string{``, `func()`}, []*Token{ret}}},
		`func()()`:           {{[]string{``, `func()`}, []*Token{ret, rparen, lparen}}, {[]string{``, `func()()`}, []*Token{ret}}},
		`func(int, int) int`: {{[]string{``, `func(int,int)`}, []*Token{ret, _int}}, {[]string{``, `func(int,int)int`}, []*Token{ret}}},
	})
}

func TestFunctionType_Render(t *testing.T) {
	defect.Equal(t, string(functionType{signature{parameters: parameters{r: parameterList{}}}}.Render()), `func()`)
}

func TestGoStmt(t *testing.T) {
	remaining(t, GoStmt, Tmap{`go a()`: {{ret, rparen, lparen}, {ret}}})
}

func TestGoStmtState(t *testing.T) {
	resultState(t, GoStmtState, map[string][]StateOutput{
		`go a()`: {{[]string{``, `go a`}, []*Token{ret, rparen, lparen}}, min(`go a()`)},
	})
}

func TestGotoStmt(t *testing.T) {
	remaining(t, GotoStmt, Tmap{`goto a`: {{ret}}, `goto`: empty, `a`: empty})
}

func TestIdentifierList(t *testing.T) {
	resultState(t, IdentifierList, map[string][]StateOutput{
		`a`:   {{[]string{``, `a`}, []*Token{ret}}},
		`a,b`: {{[]string{``, `a`}, []*Token{ret, b, comma}}, {[]string{``, `a,b`}, []*Token{ret}}},
		`1`:   nil,
		`_`:   {{[]string{``, `_`}, []*Token{ret}}},
	})
}

func TestIdentifierList_Render(t *testing.T) {
	defect.Equal(t, string(identifierList{a, one}.Render()), `a,1`)
}

func TestIfStmt(t *testing.T) {
	remaining(t, IfStmt, Tmap{
		`if a {}`:              {{ret}},
		`if a := false; a {}`:  {{ret}},
		`if a {} else {}`:      {{ret, rbrace, lbrace, _else}, {ret}},
		`if a {} else if b {}`: {{ret, rbrace, lbrace, b, _if, _else}, {ret}},
		`if a := false; a {} else if b {} else {}`: {{ret, rbrace, lbrace, _else, rbrace, lbrace, b, _if, _else},
			{ret, rbrace, lbrace, _else}, {ret}},
		`if a { b() }`: {{ret}},
	})
}

func TestImportDecl(t *testing.T) {
	remaining(t, ImportDecl, Tmap{
		`import "a"`:                     {{ret}},
		`import ()`:                      {{ret}},
		`import ("a")`:                   {{ret}},
		`import ("a";)`:                  {{ret}},
		`import ("a"; "b")`:              {{ret}},
		`import ("a";. "b";_ "c";d "d")`: {{ret}},
	})
}

func TestImportSpec(t *testing.T) {
	remaining(t, ImportSpec, Tmap{
		`"a"`:   {{ret}},
		`. "a"`: {{ret}},
		`_ "a"`: {{ret}},
		`a "a"`: {{ret}},
	})
}

func TestIncDecStmt(t *testing.T) {
	remaining(t, IncDecStmt, Tmap{
		`a++`: {{ret}},
		`a--`: {{ret}},
	})
}

func TestIncDecStmtState(t *testing.T) {
	resultState(t, IncDecStmtState, map[string][]StateOutput{
		`a++`: {min(`a++`)},
		`a--`: {min(`a--`)},
	})
}

func TestIndex(t *testing.T) {
	resultState(t, Index, map[string][]StateOutput{
		`a`:   nil,
		`[a]`: {min(`[a]`)},
		`[1]`: {min(`[1]`)},
	})
}

func TestIndex_Render(t *testing.T) {
	defect.Equal(t, string(index{one}.Render()), `[1]`)
}

func TestInterfaceType(t *testing.T) {
	resultState(t, InterfaceType, map[string][]StateOutput{
		`interface{}`:       {{[]string{``, `interface{}`}, []*Token{ret}}},
		`interface{;}`:      nil,
		`interface{a}`:      {{[]string{``, `interface{a;}`}, []*Token{ret}}},
		`interface{a;}`:     {{[]string{``, `interface{a;}`}, []*Token{ret}}},
		`interface{a()}`:    {{[]string{``, `interface{a();}`}, []*Token{ret}}},
		`interface{a();a;}`: {{[]string{``, `interface{a();a;}`}, []*Token{ret}}},
	})
}

func TestInterfaceType_Render(t *testing.T) {
	defect.Equal(t, string(interfaceType{}.Render()), `interface{}`)
	defect.Equal(t, string(interfaceType{methodSpec{iTypeName: a}}.Render()), `interface{a;}`)
	defect.Equal(t, string(interfaceType{methodSpec{iTypeName: a}, methodSpec{iTypeName: b}}.Render()), `interface{a;b;}`)
}

func TestKey(t *testing.T) {
	remaining(t, Key, Tmap{
		`a`:     {{ret}, {ret}},
		`1+a`:   {{ret, a, add}, {ret}},
		`"foo"`: {{ret}},
		`{a}`:   {{ret}},
	})
}

func TestKeyState(t *testing.T) {
	resultState(t, KeyState, map[string][]StateOutput{
		`a`:     {{[]string{``, `a`}, []*Token{ret}}},
		`1+a`:   {{[]string{``, `1`}, []*Token{ret, a, add}}, min(`1+a`)},
		`"foo"`: {min(`"foo"`)},
		`{a}`:   {min(`{a}`)},
	})
}

func TestKeyedElement(t *testing.T) {
	remaining(t, KeyedElement, Tmap{
		`1`:   {{ret}},
		`1:1`: {{ret, one, colon}, {ret}},
	})
}

func TestKeyedElementState(t *testing.T) {
	resultState(t, KeyedElementState, map[string][]StateOutput{
		`1`:   {min(`1`)},
		`1:1`: {{[]string{``, `1`}, []*Token{ret, one, colon}}, min(`1:1`)},
	})
}

func TestKeyedElement_Render(t *testing.T) {
	defect.Equal(t, string(keyedElement{element: b}.Render()), `b`)
	defect.Equal(t, string(keyedElement{a, b}.Render()), `a:b`)
}

func TestLabeledStmt(t *testing.T) {
	remaining(t, LabeledStmt, Tmap{
		`a: var b int`: {{ret}, {ret, _int, b, _var}},
	})
}

func TestLabeledStmtState(t *testing.T) {
	resultState(t, LabeledStmtState, map[string][]StateOutput{
		`a:var b int`: {min(`a:var b int`), {[]string{``, `a:`}, []*Token{ret, _int, b, _var}}},
	})
}

func TestLiteral(t *testing.T) {
	remaining(t, Literal, Tmap{
		`1`:        {{ret}},
		`T{1}`:     {{ret}},
		`a`:        empty,
		`_`:        empty,
		`func(){}`: {{ret}},
	})
}

func TestLiteralState(t *testing.T) {
	resultState(t, LiteralState, map[string][]StateOutput{
		`1`:        {min(`1`)},
		`a{1}`:     {min(`a{1}`)},
		`a`:        nil,
		`_`:        nil,
		`func(){}`: {min(`func(){}`)},
	})
}

func TestLiteralType(t *testing.T) {
	resultState(t, LiteralType, map[string][]StateOutput{
		`struct{}`:    {{[]string{``, `struct{}`}, []*Token{ret}}},
		`[1]int`:      {{[]string{``, `[1]int`}, []*Token{ret}}},
		`[...]int`:    {{[]string{``, `[...]int`}, []*Token{ret}}},
		`[]int`:       {{[]string{``, `[]int`}, []*Token{ret}}},
		`map[int]int`: {{[]string{``, `map[int]int`}, []*Token{ret}}},
		`a.a`:         {{[]string{``, `a.a`}, []*Token{ret}}, {[]string{``, `a`}, []*Token{ret, a, dot}}},
	})
}

func TestLiteralValue(t *testing.T) {
	remaining(t, LiteralValue, Tmap{
		`{1}`:           {{ret}},
		`{0: 1, 1: 2,}`: {{ret}},
	})
}

func TestLiteralValueState(t *testing.T) {
	resultState(t, LiteralValueState, map[string][]StateOutput{
		`{1}`:        {min(`{1}`)},
		`{0:1,1:2,}`: {min(`{0:1,1:2,}`)},
	})
}

func TestMapType(t *testing.T) {
	resultState(t, MapType, map[string][]StateOutput{
		`map[int]int`: {{[]string{``, `map[int]int`}, []*Token{ret}}},
		`int`:         nil,
	})
}

func TestMapType_Render(t *testing.T) {
	defect.Equal(t, string(mapType{&Token{token.IDENT, `a`}, &Token{token.IDENT, `b`}}.Render()), `map[a]b`)
}

func TestMethodDecl(t *testing.T) {
	remaining(t, MethodDecl, Tmap{
		`func (m M) f(){}`:                 {{ret}},
		`func (m M) f()(){}`:               {{ret}},
		`func (m M) f(int)int{ return 0 }`: {{ret}},
		`func (m *M) f() { a(); b(); }`:    {{ret}},
	})
}

func TestMethodExpr(t *testing.T) {
	remaining(t, MethodExpr, Tmap{
		`1`:        empty,
		`a.a`:      {{ret}},
		`a.a.a`:    {{ret}, {ret, a, dot}},
		`(a).a`:    {{ret}},
		`(a.a).a`:  {{ret}},
		`(*a.a).a`: {{ret}},
	})
}

func TestMethodExprState(t *testing.T) {
	resultState(t, MethodExprState, map[string][]StateOutput{
		`1`:        nil,
		`a.a`:      {min(`a.a`)},
		`a.a.a`:    {min(`a.a.a`), {[]string{``, `a.a`}, []*Token{ret, a, dot}}},
		`(a).a`:    {min(`(a).a`)},
		`(a.a).a`:  {min(`(a.a).a`)},
		`(*a.a).a`: {min(`(*a.a).a`)},
	})
}

func TestMethodExpr_Render(t *testing.T) {
	defect.Equal(t, string(methodExpr{receiverType{r: qualifiedIdent{a, a}}, a}.Render()), `a.a.a`)
	defect.Equal(t, string(methodExpr{receiverType{r: qualifiedIdent{a, a}, parens: 2}, a}.Render()), `((a.a)).a`)
	defect.Equal(t, string(methodExpr{receiverType{r: qualifiedIdent{a, a}, parens: 1, pointer: true}, a}.Render()), `(*a.a).a`)
}

func TestMethodSpec(t *testing.T) {
	resultState(t, MethodSpec, map[string][]StateOutput{
		`a()`:   {{[]string{``, `a()`}, []*Token{ret}}, {[]string{``, `a`}, []*Token{ret, rparen, lparen}}},
		`a`:     {{[]string{``, `a`}, []*Token{ret}}},
		`a.a()`: {{[]string{``, `a.a`}, []*Token{ret, rparen, lparen}}, {[]string{``, `a`}, []*Token{ret, rparen, lparen, a, dot}}},
	})
}

func TestMethodSpec_Render(t *testing.T) {
	defect.Equal(t, string(methodSpec{name: a, signature: signature{parameters: parameters{r: parameterList{}}}}.Render()), `a()`)
	defect.Equal(t, string(methodSpec{iTypeName: a}.Render()), `a`)
}

func TestMulOp(t *testing.T) {
	resultState(t, MulOp, map[string][]StateOutput{
		`*`:  {{[]string{``, `*`}, []*Token{}}},
		`/`:  {{[]string{``, `/`}, []*Token{}}},
		`%`:  {{[]string{``, `%`}, []*Token{}}},
		`<<`: {{[]string{``, `<<`}, []*Token{}}},
		`>>`: {{[]string{``, `>>`}, []*Token{}}},
		`&`:  {{[]string{``, `&`}, []*Token{}}},
		`&^`: {{[]string{``, `&^`}, []*Token{}}},
		`1`:  nil,
	})
}

func TestOperand(t *testing.T) {
	remaining(t, Operand, Tmap{
		`1`:     {{ret}},
		`a.a`:   {{ret, a, dot}, {ret}, {ret}},
		`(a.a)`: {{ret}, {ret}, {ret}},
	})
}

func TestOperandState(t *testing.T) {
	resultState(t, OperandState, map[string][]StateOutput{
		`1`:         {{[]string{``, `1`}, []*Token{ret}}},
		`a`:         {min(`a`)},
		`a.a`:       {{[]string{``, `a`}, []*Token{ret, a, dot}}, min(`a.a`)},
		`((a.a).a)`: {min(`((a.a).a)`)},
		`(1)`:       {min(`(1)`)},
	})
}

func TestOperand_Render(t *testing.T) {
	defect.Equal(t, string(operand{a, false}.Render()), `a`)
	defect.Equal(t, string(operand{a, true}.Render()), `(a)`)
}

func TestOperandName(t *testing.T) {
	remaining(t, OperandName, Tmap{
		`1`:   empty,
		`_`:   {{ret}},
		`a`:   {{ret}},
		`a.a`: {{ret, a, dot}, {ret}},
	})
}

func TestOperandNameState(t *testing.T) {
	resultState(t, OperandNameState, map[string][]StateOutput{
		`1`:   nil,
		`_`:   {min(`_`)},
		`a`:   {min(`a`)},
		`a.a`: {{[]string{``, `a`}, []*Token{ret, a, dot}}},
	})
}

func TestPackageName(t *testing.T) {
	resultState(t, PackageName, map[string][]StateOutput{
		`a`: {{[]string{``, `a`}, []*Token{ret}}},
		`1`: nil,
		`_`: nil,
	})
}

func TestPackageClause(t *testing.T) {
	remaining(t, PackageClause, Tmap{
		`package a`: {{ret}},
		`package _`: empty,
		`package`:   empty,
		`a`:         empty,
	})
}

func TestParameterDeclIDList(t *testing.T) {
	resultState(t, ParameterDeclIDList, map[string][]StateOutput{
		`int`:      nil,
		`...int`:   nil,
		`a, b int`: {{[]string{``, `a,b int`}, []*Token{ret}}},
		`b... int`: {{[]string{``, `b ... int`}, []*Token{ret}}},
		`int, int`: nil,
	})
}

func TestParameterDeclNoList(t *testing.T) {
	resultState(t, ParameterDeclNoList, map[string][]StateOutput{
		`int`:      {{[]string{``, `int`}, []*Token{ret}}},
		`...int`:   {{[]string{``, `... int`}, []*Token{ret}}},
		`a, b int`: {{[]string{``, `a`}, []*Token{ret, _int, b, comma}}},
		`b... int`: {{[]string{``, `b`}, []*Token{ret, _int, ellipsis}}},
		`int, int`: {{[]string{``, `int`}, []*Token{ret, _int, comma}}},
	})
}

func TestParameterDecl_Render(t *testing.T) {
	defect.Equal(t, string(parameterDecl{nil, false, _int}.Render()), `int`)
	defect.Equal(t, string(parameterDecl{nil, true, _int}.Render()), `... int`)
	defect.Equal(t, string(parameterDecl{identifierList{a, b}, false, _int}.Render()), `a,b int`)
	defect.Equal(t, string(parameterDecl{identifierList{a}, true, _int}.Render()), `a ... int`)
}

func TestParamterList(t *testing.T) {
	resultState(t, ParameterList, map[string][]StateOutput{
		`int`:       {{[]string{``, `int`}, []*Token{ret}}},
		`int, int`:  {{[]string{``, `int`}, []*Token{ret, _int, comma}}, {[]string{``, `int,int`}, []*Token{ret}}},
		`a, int`:    {{[]string{``, `a`}, []*Token{ret, _int, comma}}, {[]string{``, `a,int`}, []*Token{ret}}},
		`a int`:     {{[]string{``, `a`}, []*Token{ret, _int}}, {[]string{``, `a int`}, []*Token{ret}}},
		`a ... int`: {{[]string{``, `a`}, []*Token{ret, _int, ellipsis}}, {[]string{``, `a ... int`}, []*Token{ret}}},
		`a, b int`: {
			{[]string{``, `a`}, []*Token{ret, _int, b, comma}},
			{[]string{``, `a,b`}, []*Token{ret, _int}},
			{[]string{``, `a,b int`}, []*Token{ret}}},
		`a, b int, c ... int`: {
			{[]string{``, `a`}, []*Token{ret, _int, ellipsis, c, comma, _int, b, comma}},
			{[]string{``, `a,b`}, []*Token{ret, _int, ellipsis, c, comma, _int}},
			{[]string{``, `a,b int`}, []*Token{ret, _int, ellipsis, c, comma}},
			{[]string{``, `a,b int,c ... int`}, []*Token{ret}}},
	})
}

func TestParameterList_Render(t *testing.T) {
	defect.Equal(t, string(parameterList{parameterDecl{typ: _int}}.Render()), `int`)
	defect.Equal(t, string(parameterList{
		parameterDecl{identifierList{a, b}, false, _int},
		parameterDecl{identifierList{a}, true, _int},
	}.Render()), `a,b int,a ... int`)
}

func TestParameters(t *testing.T) {
	resultState(t, Parameters, map[string][]StateOutput{
		`()`:                    {{[]string{``, `()`}, []*Token{ret}}},
		`(int)`:                 {{[]string{``, `(int)`}, []*Token{ret}}},
		`(a, b)`:                {{[]string{``, `(a,b)`}, []*Token{ret}}},
		`(a, b,)`:               {{[]string{``, `(a,b,)`}, []*Token{ret}}},
		`(a, b int)`:            {{[]string{``, `(a,b int)`}, []*Token{ret}}},
		`(a, b int, c ... int)`: {{[]string{``, `(a,b int,c ... int)`}, []*Token{ret}}},
	})
}

func TestParameters_Render(t *testing.T) {
	defect.Equal(t, string(parameters{}.Render()), `()`)
	defect.Equal(t, string(parameters{comma: true}.Render()), `()`)
	defect.Equal(t, string(parameters{r: parameterList{parameterDecl{typ: _int}}}.Render()), `(int)`)
	defect.Equal(t, string(parameters{r: parameterList{parameterDecl{typ: _int}}, comma: true}.Render()), `(int,)`)
}

func TestPointerType(t *testing.T) {
	resultState(t, PointerType, map[string][]StateOutput{
		`*int`: {{[]string{``, `*int`}, []*Token{ret}}},
		`int`:  nil,
	})
}

func TestPointerType_Renderer(t *testing.T) {
	defect.Equal(t, string(pointerType{&Token{tok: token.IDENT, lit: `a`}}.Render()), `*a`)
}

func TestPrimaryExpr(t *testing.T) {
	remaining(t, PrimaryExpr, Tmap{
		`1`: {{ret}},
		`(a.a)("foo")`: {
			{ret, rparen, {token.STRING, `"foo"`}, lparen},
			{ret, rparen, {token.STRING, `"foo"`}, lparen},
			{ret, rparen, {token.STRING, `"foo"`}, lparen},
			{ret}, {ret}, {ret}, {ret}},
		`a.a`:     {{ret, a, dot}, {ret}, {ret}, {ret}},
		`a[1]`:    {{ret, rbrack, one, lbrack}, {ret}},
		`a[:]`:    {{ret, rbrack, colon, lbrack}, {ret}},
		`a.(int)`: {{ret, rparen, _int, lparen, dot}, {ret}},
		`a(b...)`: {{ret, rparen, ellipsis, b, lparen}, {ret}},
		`a(b...)[:]`: {{ret, rbrack, colon, lbrack, rparen, ellipsis, b, lparen},
			{ret, rbrack, colon, lbrack},
			{ret}},
	})
}

func TestPrimaryExprState(t *testing.T) {
	resultState(t, PrimaryExprState, map[string][]StateOutput{
		`1`:            {{[]string{``, `1`}, []*Token{ret}}},
		`(a)("foo")`:   {{[]string{``, `(a)`}, []*Token{ret, rparen, {tok: token.STRING, lit: `"foo"`}, lparen}}, min(`(a)("foo")`)},
		`(a.a)("foo")`: {{[]string{``, `(a.a)`}, []*Token{ret, rparen, {tok: token.STRING, lit: `"foo"`}, lparen}}, min(`(a.a)("foo")`)},
		`a.a`:          {{[]string{``, `a`}, []*Token{ret, a, dot}}, min(`a.a`)},
		`a.a.a`:        {{[]string{``, `a`}, []*Token{ret, a, dot, a, dot}}, min(`a.a.a`), {[]string{``, `a.a`}, []*Token{ret, a, dot}}},
		`(a).a`:        {min(`(a).a`), {[]string{``, `(a)`}, []*Token{ret, a, dot}}},
		`(a).a.a`:      {{[]string{``, `(a).a`}, []*Token{ret, a, dot}}, {[]string{``, `(a)`}, []*Token{ret, a, dot, a, dot}}, min(`(a).a.a`)},
		`(a.a).a`:      {min(`(a.a).a`), {[]string{``, `(a.a)`}, []*Token{ret, a, dot}}},
		`a[1]`:         {{[]string{``, `a`}, []*Token{ret, rbrack, one, lbrack}}, min(`a[1]`)},
		`a[:]`:         {{[]string{``, `a`}, []*Token{ret, rbrack, colon, lbrack}}, min(`a[:]`)},
		`a.(int)`:      {{[]string{``, `a`}, []*Token{ret, rparen, _int, lparen, dot}}, min(`a.(int)`)},
		`a()`:          {{[]string{``, `a`}, []*Token{ret, rparen, lparen}}, min(`a()`)},
		`a(b...)`:      {{[]string{``, `a`}, []*Token{ret, rparen, ellipsis, b, lparen}}, min(`a(b...)`)},
		`a(b,c)`:       {{[]string{``, `a`}, []*Token{ret, rparen, c, comma, b, lparen}}, min(`a(b,c)`)},
		`a(b)(c)`:      {{[]string{``, `a`}, []*Token{ret, rparen, c, lparen, rparen, b, lparen}}, {[]string{``, `a(b)`}, []*Token{ret, rparen, c, lparen}}, min(`a(b)(c)`)},
		`a(b...)[:]`:   {{[]string{``, `a`}, []*Token{ret, rbrack, colon, lbrack, rparen, ellipsis, b, lparen}}, {[]string{``, `a(b...)`}, []*Token{ret, rbrack, colon, lbrack}}, min(`a(b...)[:]`)},
	})
}

func TestQualifiedIdent(t *testing.T) {
	resultState(t, QualifiedIdent, map[string][]StateOutput{
		`1`:   nil,
		`a`:   nil,
		`a.a`: {{[]string{``, `a.a`}, []*Token{ret}}},
		`_.a`: nil,
		`a._`: {{[]string{``, `a._`}, []*Token{ret}}},
		`_._`: nil,
	})
}

func TestRangeClause(t *testing.T) {
	remaining(t, RangeClause, Tmap{
		`range a`:              {{ret}},
		`a[0], a[1] = range b`: {{ret}},
		`k, v := range a`:      {{ret}},
	})
}

func TestRAppend(t *testing.T) {
	first := []Renderer{one, two, a, b, c, d}
	second := rAppend(first, 2, _int)
	defect.DeepEqual(t, second, []Renderer{one, two, a, b, _int})
	first[0] = zero
	defect.DeepEqual(t, second, []Renderer{one, two, a, b, _int})
}

func TestReceiverType(t *testing.T) {
	remaining(t, ReceiverType, Tmap{
		`1`:      empty,
		`a.a`:    {{ret}, {ret, a, dot}},
		`(a.a)`:  {{ret}},
		`(*a.a)`: {{ret}},
	})
}

func TestReceiverTypeState(t *testing.T) {
	resultState(t, ReceiverTypeState, map[string][]StateOutput{
		`1`:       nil,
		`a.a`:     {min(`a.a`), {[]string{``, `a`}, []*Token{ret, a, dot}}},
		`(a.a)`:   {min(`(a.a)`)},
		`((a.a))`: {min(`((a.a))`)},
		`(*a.a)`:  {min(`(*a.a)`)},
		`*a.a`:    nil,
	})
}

func TestReceiverType_Render(t *testing.T) {
	defect.Equal(t, string(receiverType{r: _int}.Render()), `int`)
	defect.Equal(t, string(receiverType{_int, 2, true}.Render()), `((*int))`)
}

func TestRecvStmt(t *testing.T) {
	remaining(t, RecvStmt, Tmap{
		`<-a`: {{ret}},
		`a[0], b = <-c`: {{ret, c, arrow, assign, b, comma, rbrack, zero, lbrack},
			{ret, c, arrow, assign, b, comma}, {ret}},
		`a, b := <-c`: {{ret, c, arrow, define, b, comma}, {ret}},
		`a[0], b := <-c`: {{ret, c, arrow, define, b, comma, rbrack, zero, lbrack},
			{ret, c, arrow, define, b, comma}},
	})
}

func TestRelOp(t *testing.T) {
	resultState(t, RelOp, map[string][]StateOutput{
		`==`: {{[]string{``, `==`}, []*Token{}}},
		`!=`: {{[]string{``, `!=`}, []*Token{}}},
		`>`:  {{[]string{``, `>`}, []*Token{}}},
		`>=`: {{[]string{``, `>=`}, []*Token{}}},
		`<`:  {{[]string{``, `<`}, []*Token{}}},
		`<=`: {{[]string{``, `<=`}, []*Token{}}},
		`1`:  nil,
	})
}

func TestResult(t *testing.T) {
	resultState(t, Result, map[string][]StateOutput{
		`int`: {{[]string{``, `int`}, []*Token{ret}}},
		`()`:  {{[]string{``, `()`}, []*Token{ret}}},
	})
}

func TestResult_Render(t *testing.T) {
	defect.Equal(t, string(result{typ: typ{_int, 0}}.Render()), `int`)
}

func TestReturnStmt(t *testing.T) {
	remaining(t, ReturnStmt, Tmap{
		`return 1`:    {{ret, one}, {ret}},
		`return 1, 2`: {{ret, two, comma, one}, {ret, two, comma}, {ret}},
	})
}

func TestReturnStmtState(t *testing.T) {
	resultState(t, ReturnStmtState, map[string][]StateOutput{
		`return`:     {min(`return`)},
		`return 1`:   {{[]string{``, `return`}, []*Token{ret, one}}, min(`return 1`)},
		`return 1,2`: {{[]string{``, `return`}, []*Token{ret, two, comma, one}}, {[]string{``, `return 1`}, []*Token{ret, two, comma}}, min(`return 1,2`)},
	})
}

func TestSelector(t *testing.T) {
	resultState(t, Selector, map[string][]StateOutput{
		`1`:  nil,
		`a`:  nil,
		`.a`: {min(`.a`)},
	})
}

func TestSelectStmt(t *testing.T) {
	remaining(t, SelectStmt, Tmap{
		`select {}`:                                            {{ret}},
		`select {default:}`:                                    {{ret}},
		`select {default: a()}`:                                {{ret}},
		`select {default: a();}`:                               {{ret}, {ret}},
		`select {case <-a: ;default:}`:                         {{ret}},
		`select {case <-a: b(); case c<-d: e(); default: f()}`: {{ret}},
	})
}

func TestSendStmt(t *testing.T) {
	remaining(t, SendStmt, Tmap{`a <- 1`: {{ret}}})
}

func TestSendStmtState(t *testing.T) {
	resultState(t, SendStmtState, map[string][]StateOutput{`a<-1`: {min(`a<-1`)}})
}

func TestShortVarDecl(t *testing.T) {
	remaining(t, ShortVarDecl, Tmap{
		`a := 1`:       {{ret}},
		`a, b := 1, 2`: {{ret, two, comma}, {ret}},
	})
}

func TestShortVarDeclState(t *testing.T) {
	resultState(t, ShortVarDeclState, map[string][]StateOutput{
		`a:=1`:     {min(`a:=1`)},
		`a,b:=1,2`: {{[]string{``, `a,b:=1`}, []*Token{ret, two, comma}}, min(`a,b:=1,2`)},
	})
}

func TestSignature(t *testing.T) {
	resultState(t, Signature, map[string][]StateOutput{
		`()`:             {{[]string{``, `()`}, []*Token{ret}}},
		`()()`:           {{[]string{``, `()`}, []*Token{ret, rparen, lparen}}, {[]string{``, `()()`}, []*Token{ret}}},
		`(int, int) int`: {{[]string{``, `(int,int)`}, []*Token{ret, _int}}, {[]string{``, `(int,int)int`}, []*Token{ret}}},
	})
}

func TestSignature_Render(t *testing.T) {
	defect.Equal(t, string(signature{parameters: parameters{r: parameterList{}}}.Render()), `()`)
	defect.Equal(t, string(signature{parameters: parameters{r: parameterList{}},
		result: result{parameters: parameters{r: parameterList{}}}}.Render()), `()()`)
	defect.Equal(t, string(signature{parameters: parameters{r: parameterList{parameterDecl{typ: typ{_int, 0}}, parameterDecl{typ: typ{_int, 0}}}},
		result: result{typ: typ{_int, 0}}}.Render()), `(int,int)int`)
}

func TestSimpleStmt(t *testing.T) {
	remaining(t, SimpleStmt, Tmap{
		`1`:      {{ret, one}, {ret}},
		`a <- 1`: {{ret, one, arrow, a}, {ret, one, arrow}, {ret}},
		`a++`:    {{ret, inc, a}, {ret, inc}, {ret}},
		`a = 1`:  {{ret, one, assign, a}, {ret, one, assign}, {ret}},
		`a := 1`: {{ret, one, define, a}, {ret, one, define}, {ret}},
	})
}

func TestSimpleStmtState(t *testing.T) {
	resultState(t, SimpleStmtState, map[string][]StateOutput{
		`1`:    {{[]string{``}, []*Token{ret, one}}, min(`1`)},
		`a<-1`: {{[]string{``}, []*Token{ret, one, arrow, a}}, {[]string{``, `a`}, []*Token{ret, one, arrow}}, min(`a<-1`)},
		`a++`:  {{[]string{``}, []*Token{ret, inc, a}}, {[]string{``, `a`}, []*Token{ret, inc}}, min(`a++`)},
		`a=1`:  {{[]string{``}, []*Token{ret, one, assign, a}}, {[]string{``, `a`}, []*Token{ret, one, assign}}, min(`a=1`)},
		`a:=1`: {{[]string{``}, []*Token{ret, one, define, a}}, {[]string{``, `a`}, []*Token{ret, one, define}}, min(`a:=1`)},
	})
}

func TestSlice(t *testing.T) {
	resultState(t, Slice, map[string][]StateOutput{
		`[:]`:     {min(`[:]`)},
		`[a:]`:    {min(`[a:]`)},
		`[:b]`:    {min(`[:b]`)},
		`[a:b]`:   {min(`[a:b]`)},
		`[:b:c]`:  {min(`[:b:c]`)},
		`[a:b:c]`: {min(`[a:b:c]`)},
	})
}

func TestSlice_Render(t *testing.T) {
	defect.Equal(t, string(slice{nil, nil, nil}.Render()), `[:]`)
	defect.Equal(t, string(slice{a, nil, nil}.Render()), `[a:]`)
	defect.Equal(t, string(slice{nil, b, nil}.Render()), `[:b]`)
	defect.Equal(t, string(slice{a, b, nil}.Render()), `[a:b]`)
	defect.Equal(t, string(slice{nil, b, c}.Render()), `[:b:c]`)
	defect.Equal(t, string(slice{a, b, c}.Render()), `[a:b:c]`)
}

func TestSliceType(t *testing.T) {
	resultState(t, SliceType, map[string][]StateOutput{
		`[]int`: {{[]string{``, `[]int`}, []*Token{ret}}},
		`int`:   nil,
	})
}

func TestSliceType_Render(t *testing.T) {
	defect.Equal(t, string(sliceType{&Token{tok: token.IDENT, lit: `a`}}.Render()), `[]a`)
}

func TestSourceFile(t *testing.T) {
	remaining(t, SourceFile, Tmap{
		`package p`:             {{}},
		`package p; import "a"`: {{ret, {token.STRING, `"a"`}, {token.IMPORT, `import`}}, {}},
		`package p; var a int`:  {{ret, _int, a, {token.VAR, `var`}}, {}},
		`package p; import "a"; import "b"; var c int; var d int`: {
			{ret, _int, d, {token.VAR, `var`}, semi, _int, c, {token.VAR, `var`}, semi, {token.STRING, `"b"`}, {token.IMPORT, `import`},
				semi, {token.STRING, `"a"`}, {token.IMPORT, `import`}},
			{ret, _int, d, {token.VAR, `var`}, semi, _int, c, {token.VAR, `var`}, semi, {token.STRING, `"b"`}, {token.IMPORT, `import`}},
			{ret, _int, d, {token.VAR, `var`}, semi, _int, c, {token.VAR, `var`}},
			{ret, _int, d, {token.VAR, `var`}}, {}},
	})
}

func TestStatement(t *testing.T) {
	remaining(t, Statement, Tmap{
		`var a int`: {{ret}, {ret, _int, a, {token.VAR, `var`}}},
		`a: var b int`: {{ret},
			{ret, _int, b, {token.VAR, `var`}},
			{ret, _int, b, {token.VAR, `var`}, colon, a},
			{ret, _int, b, {token.VAR, `var`}, colon}},
		`a := 1`:      {{ret, one, define, a}, {ret, one, define}, {ret}},
		`go a()`:      {{ret, rparen, lparen, a, {token.GO, `go`}}, {ret, rparen, lparen}, {ret}},
		`return 1`:    {{ret, one, {token.RETURN, `return`}}, {ret, one}, {ret}},
		`break a`:     {{ret, a, {token.BREAK, `break`}}, {ret, a}, {ret}},
		`continue a`:  {{ret, a, {token.CONTINUE, `continue`}}, {ret, a}, {ret}},
		`goto a`:      {{ret, a, {token.GOTO, `goto`}}, {ret}},
		`fallthrough`: {{ret, _fallthrough}, {ret}},
		`{a()}`:       {{ret, rbrace, rparen, lparen, a, lbrace}, {ret}},
		`if a {}`:     {{ret, rbrace, lbrace, a, _if}, {ret}},
		`switch {}`:   {{ret, rbrace, lbrace, {token.SWITCH, `switch`}}, {ret}},
		`select {}`:   {{ret, rbrace, lbrace, {token.SELECT, `select`}}, {ret}},
		`for {}`:      {{ret, rbrace, lbrace, {token.FOR, `for`}}, {ret}},
		`defer a()`:   {{ret, rparen, lparen, a, {token.DEFER, `defer`}}, {ret, rparen, lparen}, {ret}},
	})
}

func TestStatementState(t *testing.T) {
	resultState(t, StatementState, map[string][]StateOutput{
		`var a int`: {min(`var a int`), {[]string{``}, []*Token{ret, _int, a, _var}}},
		`a: var b int`: {
			min(`a:var b int`),
			{[]string{``, `a:`}, []*Token{ret, _int, b, _var}},
			{[]string{``}, []*Token{ret, _int, b, _var, colon, a}},
			{[]string{``, `a`}, []*Token{ret, _int, b, _var, colon}}},
		`a := 1`: {
			{[]string{``}, []*Token{ret, one, define, a}},
			{[]string{``, `a`}, []*Token{ret, one, define}},
			min(`a:=1`)},
		`go a()`: {
			{[]string{``, `go a`}, []*Token{ret, rparen, lparen}},
			min(`go a()`),
			{[]string{``}, []*Token{ret, rparen, lparen, a, {tok: token.GO, lit: `go`}}}},
		`return 1`: {
			{[]string{``, `return`}, []*Token{ret, one}},
			min(`return 1`),
			{[]string{``}, []*Token{ret, one, {tok: token.RETURN, lit: `return`}}}},
		`break`:    {min(`break`), {[]string{``}, []*Token{ret, {tok: token.BREAK, lit: `break`}}}},
		`continue`: {min(`continue`), {[]string{``}, []*Token{ret, {tok: token.CONTINUE, lit: `continue`}}}},
	})
}

func TestStatementList(t *testing.T) {
	remaining(t, StatementList, Tmap{
		`fallthrough`:              {{ret, _fallthrough}, {ret}, {}},
		`fallthrough;`:             {{semi, _fallthrough}, {semi}, {}},
		`fallthrough; fallthrough`: {{ret, _fallthrough, semi, _fallthrough}, {ret, _fallthrough, semi}, {ret, _fallthrough}, {ret}, {}},
	})
}

func TestStatementListState(t *testing.T) {
	resultState(t, StatementListState, map[string][]StateOutput{
		``:     {{[]string{``, ``}, nil}},
		`a++`:  {{[]string{``, ``}, []*Token{ret, inc, a}}, {[]string{``, `a`}, []*Token{ret, inc}}, min(`a++`), {[]string{``, `a++;`}, []*Token{}}},
		`a++;`: {{[]string{``, ``}, []*Token{semi, inc, a}}, {[]string{``, `a`}, []*Token{semi, inc}}, {[]string{``, `a++`}, []*Token{semi}}, {[]string{``, `a++;`}, []*Token{}}},
		`a++;b++`: {
			{[]string{``, ``}, []*Token{ret, inc, b, semi, inc, a}},
			{[]string{``, `a`}, []*Token{ret, inc, b, semi, inc}},
			{[]string{``, `a++`}, []*Token{ret, inc, b, semi}},
			{[]string{``, `a++;`}, []*Token{ret, inc, b}},
			{[]string{``, `a++;b`}, []*Token{ret, inc}},
			{[]string{``, `a++;b++`}, []*Token{ret}},
			{[]string{``, `a++;b++;`}, []*Token{}}},
	})
}

func TestStatementList_Render(t *testing.T) {
	defect.Equal(t, string(statementList{}.Render()), ``)
	defect.Equal(t, string(statementList{incDecStmt{a, true}}.Render()), `a++`)
	defect.Equal(t, string(statementList{incDecStmt{a, true}, incDecStmt{b, false}}.Render()), `a++;b--`)
}

func TestStructType(t *testing.T) {
	resultState(t, StructType, map[string][]StateOutput{
		`struct{}`:                      {{[]string{``, `struct{}`}, []*Token{ret}}},
		`struct{int}`:                   {{[]string{``, `struct{int;}`}, []*Token{ret}}},
		`struct{int;}`:                  {{[]string{``, `struct{int;}`}, []*Token{ret}}},
		`struct{int;float64;}`:          {{[]string{``, `struct{int;float64;}`}, []*Token{ret}}},
		`struct{a int}`:                 {{[]string{``, `struct{a int;}`}, []*Token{ret}}},
		`struct{a, b int}`:              {{[]string{``, `struct{a,b int;}`}, []*Token{ret}}},
		`struct{a, b int;}`:             {{[]string{``, `struct{a,b int;}`}, []*Token{ret}}},
		`struct{a, b int; c, d string}`: {{[]string{``, `struct{a,b int;c,d string;}`}, []*Token{ret}}},
	})
}

func TestStructType_Render(t *testing.T) {
	defect.Equal(t, string(structType{_fallthrough}.Render()), `struct{fallthrough;}`)
	defect.Equal(t, string(structType{_fallthrough, _fallthrough}.Render()), `struct{fallthrough;fallthrough;}`)
}

func TestSwitchStmt(t *testing.T) {
	remaining(t, SwitchStmt, Tmap{
		`switch {}`:          {{ret}},
		`switch a.(type) {}`: {{ret}},
	})
}

func TestTopLevelDecl(t *testing.T) {
	remaining(t, TopLevelDecl, Tmap{
		`var a int`:        {{ret}},
		`func f(){}`:       {{ret}},
		`func (m M) f(){}`: {{ret}},
	})
}

func TestType(t *testing.T) {
	resultState(t, Type, map[string][]StateOutput{
		`a`:        {{[]string{``, `a`}, []*Token{ret}}},
		`a.a`:      {{[]string{``, `a.a`}, []*Token{ret}}, {[]string{``, `a`}, []*Token{ret, a, dot}}},
		`1`:        nil,
		`_`:        {{[]string{``, `_`}, []*Token{ret}}},
		`(a.a)`:    {{[]string{``, `(a.a)`}, []*Token{ret}}},
		`(((_)))`:  {{[]string{``, `(((_)))`}, []*Token{ret}}},
		`chan int`: {{[]string{``, `chan int`}, []*Token{ret}}},
	})
}

func TestTyp_Render(t *testing.T) {
	defect.Equal(t, string(typ{_int, 0}.Render()), `int`)
	defect.Equal(t, string(typ{_int, 2}.Render()), `((int))`)
}

func TestTypeAssertion(t *testing.T) {
	resultState(t, TypeAssertion, map[string][]StateOutput{
		`.(int)`: {min(`.(int)`)},
		`1`:      nil,
	})
}

func TestTypeCaseClause(t *testing.T) {
	remaining(t, TypeCaseClause, Tmap{
		`case a:`:      {{}},
		`case a: b()`:  {{ret, rparen, lparen, b}, {ret, rparen, lparen}, {ret}, {}},
		`default: a()`: {{ret, rparen, lparen, a}, {ret, rparen, lparen}, {ret}, {}},
		`case a, b: c(); d()`: {
			{ret, rparen, lparen, d, semi, rparen, lparen, c},
			{ret, rparen, lparen, d, semi, rparen, lparen},
			{ret, rparen, lparen, d, semi},
			{ret, rparen, lparen, d},
			{ret, rparen, lparen}, {ret}, {}},
	})
}

func TestTypeDecl(t *testing.T) {
	remaining(t, TypeDecl, Tmap{
		`type a int`:           {{ret}},
		`type (a int)`:         {{ret}},
		`type (a int;)`:        {{ret}},
		`type (a int; b int)`:  {{ret}},
		`type (a int; b int;)`: {{ret}},
	})
}

func TestTypeDeclState(t *testing.T) {
	resultState(t, TypeDeclState, map[string][]StateOutput{
		`type a int`:          {min(`type a int`)},
		`type ()`:             {min(`type ()`)},
		`type (a int)`:        {min(`type (a int)`)},
		`type (a int;)`:       {min(`type (a int;)`)},
		`type (a int; b int)`: {min(`type (a int;b int)`)},
	})
}

func TestTypeDecl_Render(t *testing.T) {
	defect.Equal(t, string(typeDecl{[]Renderer{typeDef{a, b}}, false, false}.Render()), `type a b`)
	defect.Equal(t, string(typeDecl{nil, true, false}.Render()), `type ()`)
	defect.Equal(t, string(typeDecl{[]Renderer{typeDef{a, b}}, true, true}.Render()), `type (a b;)`)
	defect.Equal(t, string(typeDecl{[]Renderer{typeDef{a, b}, typeDef{c, d}}, true, false}.Render()), `type (a b;c d)`)
}

func TestTypeDefState(t *testing.T) {
	resultState(t, TypeDefState, map[string][]StateOutput{`a b`: {min(`a b`)}})
}

func TestTypeList(t *testing.T) {
	remaining(t, TypeList, Tmap{
		`a`:     {{ret}},
		`a, b`:  {{ret, b, comma}, {ret}},
		`a, b,`: {{comma, b, comma}, {comma}},
	})
}

func TestTypeLit(t *testing.T) {
	resultState(t, TypeLit, map[string][]StateOutput{
		`[1]int`:      {{[]string{``, `[1]int`}, []*Token{ret}}},
		`struct{}`:    {{[]string{``, `struct{}`}, []*Token{ret}}},
		`*int`:        {{[]string{``, `*int`}, []*Token{ret}}},
		`func()`:      {{[]string{``, `func()`}, []*Token{ret}}},
		`interface{}`: {{[]string{``, `interface{}`}, []*Token{ret}}},
		`[]int`:       {{[]string{``, `[]int`}, []*Token{ret}}},
		`map[int]int`: {{[]string{``, `map[int]int`}, []*Token{ret}}},
		`chan int`:    {{[]string{``, `chan int`}, []*Token{ret}}},
	})
}

func TestTypeName(t *testing.T) {
	resultState(t, TypeName, map[string][]StateOutput{
		`a`:   {{[]string{``, `a`}, []*Token{ret}}},
		`a.a`: {{[]string{``, `a.a`}, []*Token{ret}}, {[]string{``, `a`}, []*Token{ret, a, dot}}},
		`1`:   nil,
		`_`:   {{[]string{``, `_`}, []*Token{ret}}},
	})
}

func TestTypeSpec(t *testing.T) {
	remaining(t, TypeSpec, Tmap{
		`a int`: {{ret}},
		`a`:     empty,
	})
}

func TestTypeSpecState(t *testing.T) {
	resultState(t, TypeSpecState, map[string][]StateOutput{
		`a b`: {min(`a b`)},
		`a=b`: {min(`a=b`)},
	})
}

func TestTypeSwitchCase(t *testing.T) {
	remaining(t, TypeSwitchCase, Tmap{
		`default`:   {{}},
		`case a`:    {{ret}},
		`case a, b`: {{ret, b, comma}, {ret}},
	})
}

func TestTypeSwitchGuard(t *testing.T) {
	remaining(t, TypeSwitchGuard, Tmap{
		`a.(type)`:      {{ret}},
		`a := a.(type)`: {{ret}},
	})
}

func TestTypeSwitchStmt(t *testing.T) {
	remaining(t, TypeSwitchStmt, Tmap{
		`switch a.(type) {}`:                              {{ret}},
		`switch a := 5; a := a.(type) {}`:                 {{ret}},
		`switch a.(type) { case int: }`:                   {{ret}},
		`switch a.(type) { case int:; }`:                  {{ret}, {ret}},
		`switch a.(type) { case int: b() }`:               {{ret}},
		`switch a.(type) { case int: b(); }`:              {{ret}, {ret}},
		`switch a.(type) { case int: b(); default: c() }`: {{ret}},
	})
}

func TestUnaryExpr(t *testing.T) {
	remaining(t, UnaryExpr, Tmap{
		`1`:  {{ret}},
		`-1`: {{ret}},
		`!a`: {{ret}},
	})
}

func TestUnaryExprState(t *testing.T) {
	resultState(t, UnaryExprState, map[string][]StateOutput{
		`1`:  {{[]string{``, `1`}, []*Token{ret}}},
		`-1`: {min(`-1`)},
		`!a`: {min(`!a`)},
	})
}

func TestUnaryExpr_Render(t *testing.T) {
	defect.Equal(t, string(unaryExpr{a, add, arrow}.Render()), `<-+a`)
}

func TestUnaryOp(t *testing.T) {
	resultState(t, UnaryOp, map[string][]StateOutput{
		`+`:  {{[]string{``, `+`}, []*Token{}}},
		`-`:  {{[]string{``, `-`}, []*Token{}}},
		`!`:  {{[]string{``, `!`}, []*Token{}}},
		`^`:  {{[]string{``, `^`}, []*Token{}}},
		`*`:  {{[]string{``, `*`}, []*Token{}}},
		`&`:  {{[]string{``, `&`}, []*Token{}}},
		`<-`: {{[]string{``, `<-`}, []*Token{}}},
		`1`:  nil,
	})
}

func TestVarDecl(t *testing.T) {
	remaining(t, VarDecl, Tmap{
		`var a int`:                 {{ret}},
		`var (a int)`:               {{ret}},
		`var (a int;)`:              {{ret}},
		`var (a, b = 1, 2)`:         {{ret}},
		`var (a, b = 1, 2; c int;)`: {{ret}},
	})
}

func TestVarDeclState(t *testing.T) {
	resultState(t, VarDeclState, map[string][]StateOutput{
		`var a int`:           {min(`var a int`)},
		`var ()`:              {min(`var ()`)},
		`var (a int)`:         {min(`var (a int)`)},
		`var (a int;)`:        {min(`var (a int;)`)},
		`var (a,b=1,2;c int)`: {min(`var (a,b=1,2;c int)`)},
	})
}

func TestVarDecl_Render(t *testing.T) {
	defect.Equal(t, string(varDecl{[]Renderer{varSpec{a, b, nil}}, false, false}.Render()), `var a b`)
	defect.Equal(t, string(varDecl{nil, true, false}.Render()), `var ()`)
	defect.Equal(t, string(varDecl{[]Renderer{varSpec{a, b, nil}}, true, true}.Render()), `var (a b;)`)
	defect.Equal(t, string(varDecl{[]Renderer{varSpec{a, b, nil}, varSpec{c, d, nil}}, true, false}.Render()), `var (a b;c d)`)
}

func TestVarSpec(t *testing.T) {
	remaining(t, VarSpec, Tmap{
		`a int`:       {{ret}},
		`a int = 1`:   {{ret, one, assign}, {ret}},
		`a = 1`:       {{ret}},
		`a, b = 1, 2`: {{ret, two, comma}, {ret}},
	})
}

func TestVarSpecState(t *testing.T) {
	resultState(t, VarSpecState, map[string][]StateOutput{
		`a int`:   {min(`a int`)},
		`a int=1`: {{[]string{``, `a int`}, []*Token{ret, one, assign}}, min(`a int=1`)},
		`a=1`:     {min(`a=1`)},
		`a,b=1,2`: {{[]string{``, `a,b=1`}, []*Token{ret, two, comma}}, min(`a,b=1,2`)},
	})
}

func TestVarSpec_Render(t *testing.T) {
	defect.Equal(t, string(varSpec{identifierList{a}, _int, nil}.Render()), `a int`)
	defect.Equal(t, string(varSpec{identifierList{a}, nil, expressionList{one}}.Render()), `a=1`)
	defect.Equal(t, string(varSpec{identifierList{a}, _int, expressionList{one}}.Render()), `a int=1`)
}

func TestTokenReader(t *testing.T) {
	toks := [][]*Token{
		{ret, rparen},
		{ret, rparen, a, dot},
	}
	defect.DeepEqual(t, tokenReader(toks, token.RPAREN), [][]*Token{{ret}})
}

func TestTokenParser(t *testing.T) {
	toks := [][]*Token{
		{ret, rparen},
		{ret, rparen, a, dot},
	}
	tree, state := tokenParser(toks, token.RPAREN)
	defect.DeepEqual(t, state, [][]*Token{{ret}})
	defect.DeepEqual(t, tree, []Renderer{rparen})
}

func TestTokenReaderState(t *testing.T) {
	state := []State{
		{[]Renderer{e{}}, []*Token{ret, rparen}},
		{[]Renderer{e{}}, []*Token{ret, rparen, a, dot}},
	}
	newState := tokenReaderState(state, token.RPAREN)
	defect.DeepEqual(t, newState, []State{{[]Renderer{e{}}, []*Token{ret}}})
}

func TestTokenParserState(t *testing.T) {
	state := []State{
		{[]Renderer{e{}}, []*Token{ret, rparen}},
		{[]Renderer{e{}}, []*Token{ret, rparen, a, dot}},
	}
	newState := tokenParserState(state, token.RPAREN)
	defect.DeepEqual(t, newState, []State{{[]Renderer{e{}, rparen}, []*Token{ret}}})
}

func min(s string) StateOutput { return StateOutput{[]string{``, s}, []*Token{ret}} }
