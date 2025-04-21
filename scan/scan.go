package scan

import (
	"fmt"
)

type Type int

type Token struct {
	Type Type
	Line int
	Text string
}

func (token Token) String() string {
	return fmt.Sprintf("%d %q", token.Type, token.Text)
}

const (
	EOF Type = iota
	Newline
	// Single
	LeftParen    // (
	RightParen   // )
	LeftBracket  // [
	RightBracket // ]
	LeftCurly    // {
	RightCurly   // }
	LeftAngle    // <
	RightAngle   // >
	Assign       // =
	Comma        // ,
	Dot          // .
	Colon        // :
	SemiColon    // ;
	Bang         // !
	Slash        // /
	Star         // *
	Plus         // +
	Minus        // -
	Pipe         // |

	// Multiple
	Equals        // ==
	NotEquals     // !=
	GreaterEquals // >=
	LesserEquals  // <=

	// Literals
	Identifier // foo
	String     // "foo"
	Number     // 1337
	True       // true
	False      // false

	// Keywords
	Struct // struct
	Return // return
	Int    // int
	Double // double
	Float  // float
	Bool   // bool
	For    // for
	In     // in
	Let    // let
	If     // if
	Else   // else
)

func keywordOrIdentifier(text string) Type {
	switch text {
	case "struct":
		return Struct
	case "return":
		return Return
	case "int":
		return Int
	case "double":
		return Double
	case "float":
		return Float
	case "bool":
		return Bool
	case "for":
		return For
	case "in":
		return In
	case "let":
		return Let
	case "if":
		return If
	case "else":
		return Else
	case "true":
		return True
	case "false":
		return False
	default:
		return Identifier
	}
}

type Scanner struct {
	tokens  []Token
	source  string
	start   int
	current int
	line    int
	errors  []string
}

func NewScanner(source string) Scanner {
	return Scanner{
		tokens:  make([]Token, 0),
		source:  source,
		start:   0,
		current: 0,
		line:    1,
		errors:  make([]string, 0),
	}
}

func (scanner *Scanner) Scan() ([]Token, []string) {
	for !scanner.end() {
		scanner.start = scanner.current
		scanner.scanToken()
	}

	scanner.addToken(Token{Type: EOF, Line: scanner.line})
	return scanner.tokens, scanner.errors
}

func (scanner *Scanner) scanToken() {
	c := scanner.advance()

	switch c {
	case '(':
		scanner.addToken(scanner.newToken(LeftParen, string(c)))
	case ')':
		scanner.addToken(scanner.newToken(RightParen, string(c)))
	case '[':
		scanner.addToken(scanner.newToken(LeftBracket, string(c)))
	case ']':
		scanner.addToken(scanner.newToken(RightBracket, string(c)))
	case '{':
		scanner.addToken(scanner.newToken(LeftCurly, string(c)))
	case '}':
		scanner.addToken(scanner.newToken(RightCurly, string(c)))
	case '<':
		if scanner.match('=') {
			scanner.addToken(scanner.newToken(LesserEquals, scanner.lexeme()))
		} else {
			scanner.addToken(scanner.newToken(LeftAngle, string(c)))
		}
	case '>':
		if scanner.match('=') {
			scanner.addToken(scanner.newToken(GreaterEquals, scanner.lexeme()))
		} else {
			scanner.addToken(scanner.newToken(RightAngle, string(c)))
		}
	case '=':
		if scanner.match('=') {
			scanner.addToken(scanner.newToken(Equals, scanner.lexeme()))
		} else {
			scanner.addToken(scanner.newToken(Assign, string(c)))
		}
	case '!':
		if scanner.match('=') {
			scanner.addToken(scanner.newToken(NotEquals, scanner.lexeme()))
		} else {
			scanner.addToken(scanner.newToken(Bang, string(c)))
		}
	case ',':
		scanner.addToken(scanner.newToken(Comma, string(c)))
	case '.':
		scanner.addToken(scanner.newToken(Dot, string(c)))
	case ':':
		scanner.addToken(scanner.newToken(Colon, string(c)))
	case ';':
		scanner.addToken(scanner.newToken(SemiColon, string(c)))
	case '/':
		if scanner.match('/') {
			for scanner.peek() != '\n' && !scanner.end() {
				scanner.advance()
			}
		} else {
			scanner.addToken(scanner.newToken(Slash, string(c)))
		}
	case '*':
		scanner.addToken(scanner.newToken(Star, string(c)))
	case '+':
		scanner.addToken(scanner.newToken(Plus, string(c)))
	case '-':
		scanner.addToken(scanner.newToken(Minus, string(c)))
	case '|':
		scanner.addToken(scanner.newToken(Pipe, string(c)))
	case '"':
		scanner.stringLiteral()
	case ' ':
	case '\r':
	case '\t':
	case '\n':
		scanner.line++
		scanner.addToken(scanner.newToken(Newline, string(c)))
	default:
		if isDigit(c) {
			scanner.numberLiteral()
		} else if isAlpha(c) {
			scanner.identifier()
		} else {
			if c != 0 {
				scanner.err(fmt.Sprintf("Unexpected character '%c'", c))
			}
		}
	}
}

func (scanner *Scanner) identifier() {
	for isAlphaNumeric(scanner.peek()) {
		scanner.advance()
	}

	text := scanner.lexeme()
	typ := keywordOrIdentifier(text)
	scanner.addToken(scanner.newToken(typ, text))
}

func (scanner *Scanner) numberLiteral() {
	for isDigit(scanner.peek()) {
		scanner.advance()
	}

	if scanner.peek() == '.' && isDigit(scanner.peekNext()) {
		scanner.advance()

		for isDigit(scanner.peek()) {
			scanner.advance()
		}
	}

	scanner.addToken(scanner.newToken(Number, scanner.lexeme()))
}

func (scanner *Scanner) stringLiteral() {
	for scanner.peek() != '"' && !scanner.end() {
		if scanner.peek() == '\n' {
			scanner.err("unterminated string")
		}
		scanner.advance()
	}

	if scanner.end() {
		scanner.err("unterminated string")
		return
	}

	scanner.advance()

	literal := scanner.source[scanner.start+1 : scanner.current-1]
	scanner.addToken(scanner.newToken(String, literal))
}

func (scanner *Scanner) err(msg string) {
	scanner.errors = append(scanner.errors, fmt.Sprintf("%s on line %d", msg, scanner.line))
}

func isDigit(c byte) bool {
	return c >= '0' && c <= '9'
}

func isAlpha(c byte) bool {
	return (c >= 'a' && c <= 'z') ||
		(c >= 'A' && c <= 'Z') ||
		c == '_'
}

func isAlphaNumeric(c byte) bool {
	return isAlpha(c) || isDigit(c)
}

func (scanner *Scanner) match(c byte) bool {
	if scanner.end() {
		return false
	}

	if scanner.source[scanner.current] != c {
		return false
	}

	scanner.current++
	return true
}

func (scanner *Scanner) peek() byte {
	if scanner.end() {
		return 0
	}
	return scanner.source[scanner.current]
}

func (scanner *Scanner) peekNext() byte {
	if scanner.current+1 >= len(scanner.source) {
		return 0
	}
	return scanner.source[scanner.current+1]
}

func (scanner *Scanner) advance() byte {
	scanner.current++
	if !scanner.end() {
		return scanner.source[scanner.current-1]
	}
	return 0
}

func (scanner *Scanner) addToken(token Token) {
	scanner.tokens = append(scanner.tokens, token)
}

func (scanner *Scanner) newToken(tokenType Type, text string) Token {
	return Token{
		Type: tokenType,
		Text: text,
		Line: scanner.line,
	}
}

func (scanner *Scanner) lexeme() string {
	return scanner.source[scanner.start:scanner.current]
}

func (scanner *Scanner) end() bool {
	return scanner.current >= len(scanner.source)
}
