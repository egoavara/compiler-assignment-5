package cdtgo

import (
	"errors"
	"fmt"
	"strings"
	"unicode"
)

const (
	NO_KEYWORD = 7
	ID_LENGTH  = 12
)

type Token struct {
	Kind TokenKind
	Data interface{}
}

func (s *Token) String() string {
	result := fmt.Sprintf("Token(%s", s.Kind.String())
	if s.Data != nil{
		result += ", "
		result += fmt.Sprint(s.Data)
	}
	result += ")"
	return result
}

type TokenKind int

const (
	Tnull      TokenKind = -1
	Tnot       TokenKind = 0
	Tnotequ    TokenKind = 1
	Tremainder TokenKind = 2
	TremAssign TokenKind = 3
	Tident     TokenKind = 4
	Tnumber    TokenKind = 5
	Tand       TokenKind = 6
	Tlparen    TokenKind = 7
	Trparen    TokenKind = 8
	Tmul       TokenKind = 9
	TmulAssign TokenKind = 10
	Tplus      TokenKind = 11
	Tinc       TokenKind = 12
	TaddAssign TokenKind = 13
	Tcomma     TokenKind = 14
	Tminus     TokenKind = 15
	Tdec       TokenKind = 16
	TsubAssign TokenKind = 17
	Tdiv       TokenKind = 18
	TdivAssign TokenKind = 19
	Tcolon     TokenKind = 20
	Tsemicolon TokenKind = 21
	Tless      TokenKind = 22
	Tlesse     TokenKind = 23
	Tassign    TokenKind = 24
	Tequal     TokenKind = 25
	Tgreat     TokenKind = 26
	Tgreate    TokenKind = 27
	Tlbracket  TokenKind = 28
	Trbracket  TokenKind = 29
	Teof       TokenKind = 30

	//   ...........    word symbols ................................. //
	Tbreak    TokenKind = 31
	Tcase     TokenKind = 32
	Tconst    TokenKind = 33
	Tcontinue TokenKind = 34
	Tdefault  TokenKind = 35
	Tdo       TokenKind = 36
	Telse     TokenKind = 37
	Tfor      TokenKind = 38
	Tif       TokenKind = 39
	Tint      TokenKind = 40
	Treturn   TokenKind = 41
	Tswitch   TokenKind = 42
	Tvoid     TokenKind = 43
	Twhile    TokenKind = 44
	Tlbrace   TokenKind = 45
	Tor       TokenKind = 46
	Trbrace   TokenKind = 47
)

// 기존 문법
// ':'
// "break", "case", "continue", "default", "do", "for", "switch"
var tokenName = []string{
	"!", "!=", "%", "%=", "%ident", "%number",
	"&&", "(", ")", "*", "*=", "+",
	"++", "+=", ",", "-", "--", "-=",
	"/", "/=", ":", ";", "<", "<=",
	"=", "==", ">", ">=", "[", "]",
	"_|_",
	//   ...........    word symbols ................................. //
	"break", "case", "const", "continue", "default", "do",
	"else", "for", "if", "int", "return", "switch",
	"void", "while", "{", "||", "}",
}
var keywords = map[string]TokenKind{
	"break": Tbreak, "case": Tcase, "continue": Tcontinue, "default": Tdefault, "do": Tdo, "for": Tfor, "switch": Tswitch,
	"const": Tconst, "else": Telse, "if": Tif, "int": Tint, "return": Treturn, "void": Tvoid, "while": Twhile,
}

func (t TokenKind) String() string {
	if t == Tnull {
		return "%null"
	}
	return tokenName[t]
}
func Scanning(input string) ([]Token, error) {
	result := make([]Token, 0)
	for left, tk := ScanningStep(input); len(left) > 0 || tk.Kind != Teof; left, tk = ScanningStep(left) {
		if tk.Kind == Tnull {
			if err, ok := tk.Data.(error); ok {
				return nil, fmt.Errorf("%v : %s", err, left)
			}
			return nil, errors.New("unknown error")
		}
		result = append(result, tk)
	}
	result = append(result, Token{
		Kind: Teof,
	})
	return result, nil
}
func ScanningStep(input string) (left string, found Token) {
	found.Kind = Tnull
	// 1단계 : 공백 문자는 생략
	var (
		i = strings.IndexFunc(input, func(r rune) bool { return !unicode.IsSpace(r) })
	)
	if i == -1 {
		found.Kind = Teof
		return "", found
	}
	var (
		c       = []rune(input)[i]
		catched string
	)
	// 2단계 :
	if superLetter(c) {
		// 시작은 [_a-zA-Z]
		// 경우의 수는 2가지, identifier 혹은 reserved word
		catched, left = whileString(input[i+1:], superLetterOrDigit)
		catched = string(c) + catched
		if k, ok := keywords[catched]; ok {
			found.Kind = k
			return left, found
		} else {
			found.Kind = Tident
			found.Data = catched
			return left, found
		}
	} else if isDigit(c) {
		if strings.HasPrefix(input[i:], "0x") || strings.HasPrefix(input[i:], "0X") {
			// 16진법
			catched, left = whileString(input[i+2:], isDigit)
			if len(catched) == 0 {
				found.Data = errors.New("hex, no")
				return input, found
			}
			result := 0
			for _, cat := range catched {
				result = result*16 + hexValue(cat)
			}
			found.Kind = Tnumber
			found.Data = result
			return left, found
		} else if strings.HasPrefix(input[i:], "0") {
			// 8진법
			catched, left = whileString(input[i+1:], isDigit)
			if len(catched) == 0 {
				found.Kind = Tnumber
				found.Data = int(0)
				return input[i+1:], found
			}
			result := 0
			for _, cat := range catched {
				result = result*8 + hexValue(cat)
			}
			found.Kind = Tnumber
			found.Data = result
			return left, found
		} else {
			// 10진법
			catched, left = whileString(input[i:], isDigit)
			if len(catched) == 0 {
				found.Data = errors.New("dec, no")
				return input, found
			}
			result := 0
			for _, cat := range catched {
				result = result*10 + hexValue(cat)
			}
			found.Kind = Tnumber
			found.Data = result
			return left, found
		}
	} else {
		switch c {
		case '/':
			if strings.HasPrefix(input[i:], "/*") {
				idx := strings.Index(input[i:], "*/")
				return ScanningStep(input[idx+2:])
			} else if strings.HasPrefix(input[i:], "//") {
				idx := strings.Index(input[i:], "\n")
				return ScanningStep(input[i+idx+1:])
			} else if strings.HasPrefix(input[i:], "/=") {
				found.Kind = TdivAssign
				return input[i+2:], found
			} else {
				found.Kind = Tdiv
				return input[i+1:], found
			}
		case '!':
			if strings.HasPrefix(input[i:], "!=") {
				found.Kind = Tnotequ
				return input[i+2:], found
			} else {
				found.Kind = Tnot
				return input[i+1:], found
			}
		case '%':
			if strings.HasPrefix(input[i:], "%=") {
				found.Kind = TremAssign
				return input[i+2:], found
			} else {
				found.Kind = Tremainder
				return input[i+1:], found
			}
		case '&':
			if strings.HasPrefix(input[i:], "&&") {
				found.Kind = Tand
				return input[i+2:], found
			} else {
				found.Data = errors.New("next character must be &")
				return input, found
			}
		case '*':
			if strings.HasPrefix(input[i:], "*=") {
				found.Kind = TmulAssign
				return input[i+2:], found
			} else {
				found.Kind = Tmul
				return input[i+1:], found
			}
		case '+':
			if strings.HasPrefix(input[i:], "++") {
				found.Kind = Tinc
				return input[i+2:], found
			} else if strings.HasPrefix(input[i:], "+=") {
				found.Kind = TaddAssign
				return input[i+2:], found
			} else {
				found.Kind = Tplus
				return input[i+1:], found
			}
		case '-':
			if strings.HasPrefix(input[i:], "--") {
				found.Kind = Tdec
				return input[i+2:], found
			} else if strings.HasPrefix(input[i:], "-=") {
				found.Kind = TsubAssign
				return input[i+2:], found
			} else {
				found.Kind = Tminus
				return input[i+1:], found
			}
		case '<':
			if strings.HasPrefix(input[i:], "<=") {
				found.Kind = Tlesse
				return input[i+2:], found
			} else {
				found.Kind = Tless
				return input[i+1:], found
			}
		case '=':
			if strings.HasPrefix(input[i:], "==") {
				found.Kind = Tequal
				return input[i+2:], found
			} else {
				found.Kind = Tassign
				return input[i+1:], found
			}
		case '>':
			if strings.HasPrefix(input[i:], ">=") {
				found.Kind = Tgreate
				return input[i+2:], found
			} else {
				found.Kind = Tgreat
				return input[i+1:], found
			}
		case '|':

			if strings.HasPrefix(input[i:], "||") {
				found.Kind = Tor
				return input[i+2:], found
			} else {
				found.Data = errors.New("next character must be |")
				return input, found
			}
		case '(':
			found.Kind = Tlparen
			return input[i+1:], found
		case ')':
			found.Kind = Trparen
			return input[i+1:], found
		case ',':
			found.Kind = Tcomma
			return input[i+1:], found
		case ':':
			found.Kind = Tcolon
			return input[i+1:], found
		case ';':
			found.Kind = Tsemicolon
			return input[i+1:], found
		case '[':
			found.Kind = Tlbracket
			return input[i+1:], found
		case ']':
			found.Kind = Trbracket
			return input[i+1:], found
		case '{':
			found.Kind = Tlbrace
			return input[i+1:], found
		case '}':
			found.Kind = Trbrace
			return input[i+1:], found
		default:
			found.Data = errors.New("invalid character")
			return input, found
		}
	}
}
func isAlpha(ch rune) bool {
	return 'a' <= ch && ch <= 'z' || 'A' <= ch && ch <= 'Z'
}

func isDigit(ch rune) bool {
	return '0' <= ch && ch <= '9'
}

// [_a-zA-Z]
func superLetter(ch rune) bool {
	return isAlpha(ch) || ch == '_'
}

// [_a-zA-Z0-9]
func superLetterOrDigit(ch rune) bool {
	return isAlpha(ch) || isDigit(ch) || ch == '_'
}

func hexValue(ch rune) int {
	switch ch {
	case '0':
		fallthrough
	case '1':
		fallthrough
	case '2':
		fallthrough
	case '3':
		fallthrough
	case '4':
		fallthrough
	case '5':
		fallthrough
	case '6':
		fallthrough
	case '7':
		fallthrough
	case '8':
		fallthrough
	case '9':
		return int(ch - '0')
	case 'A':
		fallthrough
	case 'B':
		fallthrough
	case 'C':
		fallthrough
	case 'D':
		fallthrough
	case 'E':
		fallthrough
	case 'F':
		return int((ch - 'A') + 10)
	case 'a':
		fallthrough
	case 'b':
		fallthrough
	case 'c':
		fallthrough
	case 'd':
		fallthrough
	case 'e':
		fallthrough
	case 'f':
		return int((ch - 'a') + 10)
	default:
		return -1
	}
}

func char(ch rune) func(rune) bool {
	return func(r rune) bool {
		return r == ch
	}
}
func not(param0 func(rune) bool) func(rune) bool {
	return func(r rune) bool {
		return !param0(r)
	}
}
func whileString(str string, param0 func(rune) bool) (catched string, left string) {
	for i, c := range str {
		if !param0(c) {
			return str[:i], str[i:]
		}
	}
	return str, ""
}
