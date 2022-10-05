package mapq

import (
	"errors"
	_ "errors"
	"fmt"
	_ "fmt"
	"strings"
)

var loc int

const (
	TYPE_PLUS      = iota // "+"0
	TYPE_SUB              // "-"0
	TYPE_MUL              // "*"0
	TYPE_DIV              // "/"0
	TYPE_LP               // "("0
	TYPE_RP               // ")"0
	TYPE_VAR              // "([a-z]|[A-Z])([a-z]|[A-Z]|[0-9])*"0
	TYPE_RES_TRUE         // "true"
	TYPE_RES_FALSE        // "false"
	TYPE_AND              // "&&"0
	TYPE_OR               // "||"
	TYPE_EQ               // "=="
	TYPE_LG               // ">"0
	TYPE_SM               // "<"0
	TYPE_LEQ              // ">="0
	TYPE_SEQ              // "<="0
	TYPE_NEQ              // "!="0
	TYPE_STR              // a quoted string(单引号)0
	TYPE_INT              // an integer
	TYPE_FLOAT            // 小数，x.y这种
	TYPE_UNKNOWN          // 未知的类型
	TYPE_NOT              // "!"0
	TYPE_DOT              // "."0
	TYPE_RES_NULL         // "null"
)

// Lexer 词法分析器
type Lexer struct {
	input string
	pos   int //开始是0
	runes []rune
}

// Peek 看下一个字符
func (l *Lexer) Peek() (ch rune, end bool) {
	l.pos += 1
	return
}

// some finction maybe useful for your implementation
func isLetter(ch rune) bool {

	return ('a' <= ch && ch <= 'z') || ('A' <= ch && ch <= 'Z')
}
func isLetterOrUnderscore(ch rune) bool {
	return isLetter(ch) || ch == '_' || ch == '-'
}
func isNum(ch rune) bool {
	return '0' <= ch && ch <= '9'
}

// Checkpoint 检查点
type Checkpoint struct {
	pos int
}

// SetCheckpoint 设置检查点
func (l *Lexer) SetCheckpoint() Checkpoint {

	panic("not implemented")
}

// GobackTo 回到一个检查点
func (l *Lexer) GobackTo(c Checkpoint) {
	panic("not implemented")
}
func (l *Lexer) fin() bool {
	return l.pos >= len(l.runes)-1
}

// SetInput 设置输入
func (l *Lexer) SetInput(s string) {
	for i := 0; i <= len(s); i++ {
		l.runes = append([]rune(s))
	}
	//fmt.Println(l.runes)

}

// Scan scan a token
func (l *Lexer) Scan() (code int, token string, eos bool) {
	for l.pos <= len(l.runes) {
		var char = l.runes[loc]
		switch char {
		case '(': //左括号
			code = TYPE_LP
			token = string('(')
			eos = true
			l.Peek()
			return
		case ')': //右括号
			code = TYPE_RP
			token = string(')')
			eos = true
			l.Peek()
			return
		case '+':
			code = TYPE_PLUS
			token = string('+')
			eos = true
			l.Peek()
			return
		case '-':
			code = TYPE_SUB
			token = string('-')
			eos = true
			l.Peek()
			return
		case '*':
			code = TYPE_MUL
			token = string('*')
			eos = true
			l.Peek()
			return
		case '/':
			code = TYPE_DIV
			token = string('/')
			eos = true
			l.Peek()
			return
		case '<':
			if l.fin() {
				l.Peek()
				char = l.runes[loc]
				switch char {
				case '=':
					code = TYPE_SEQ
					token = string("<=")
					eos = true
				default:
					code = TYPE_SM
					token = string('<')
					eos = true
				}
			} else {
				code = TYPE_SM
				token = string('<')
				eos = true
			}
			return

		case '>':
			if l.fin() {
				l.Peek()
				char = l.runes[loc]
				switch char {
				case '=':
					code = TYPE_LEQ
					token = string(">=")
					eos = true
				default:
					code = TYPE_LG
					token = string('>')
					eos = true
				}
			} else {
				code = TYPE_LG
				token = string('<')
				eos = true
			}
			return

		case '!':
			if l.fin() {
				l.Peek()
				char = l.runes[loc]
				switch char {
				case '=':
					code = TYPE_NEQ
					token = "!="
					eos = true
				default:
					code = TYPE_NOT
					token = "!"
					eos = true
				}
			} else {
				code = TYPE_NOT
				token = string('>')
				eos = true
			}
			return
		case '=':
			if l.fin() {
				l.Peek()
				char = l.runes[loc]
				if char == '=' {
					code = TYPE_EQ
					token = "=="
					eos = true
				}
			} else {
				code = TYPE_EQ
				token = string('>')
				eos = true
			}

			return
		case '&':
			l.Peek()
			char = l.runes[loc]
			if char == '&' {
				code = TYPE_AND
				token = "&&"
				eos = true
			}
			return
		case '|':
			l.Peek()
			char = l.runes[loc]
			if char == '&' {
				code = TYPE_OR
				token = "||"
				eos = true
			}
			return
		case '.':
			l.Peek()
			code = TYPE_DOT
			token = "."
			eos = true
			return
		case 39: // string

			l.Peek()
			char = l.runes[loc]
			for char != 39 {

				token += string(char)
				code = TYPE_STR
				eos = true

				l.Peek()
				char = l.runes[loc]
			}
			if char == 39 {

				return
			}
		default: // integers or identifiers
			if isLetterOrUnderscore(char) {
				for isLetterOrUnderscore(char) || isNum(char) {
					if l.fin() {
						token += string(char)
						if token == "true" {
							code = TYPE_RES_TRUE
							eos = true
							return
						} else if token == "false" {
							code = TYPE_RES_FALSE
							eos = true
							return
						} else {
							code = TYPE_VAR
							eos = true
							return
						}
					}

					token += string(char)
					code = TYPE_VAR
					eos = true
					l.Peek()
					char = l.runes[loc]
					if !isLetterOrUnderscore(char) && !isNum(char) {
						fmt.Println(loc)
						return
					}
				}

				return
			}

			if isNum(char) {
				for isNum(char) {
					if l.fin() {
						if !strings.Contains(token, ".") {
							code = TYPE_INT
							eos = true
						}
						token += string(char)
						return
					}
					token += string(char)
					l.Peek()
					char = l.runes[loc]
				}

				if char == '.' {
					code = TYPE_FLOAT
					eos = true
					token += string(char)
					l.Peek()
					char = l.runes[loc]
					for isNum(char) {
						if l.fin() {
							token += string(char)
							return
						}
						token += string(char)
						l.Peek()
						char = l.runes[loc]
					}
				}
			}
		}
	}
	return
}

// ScanType 扫描一个特定Token，下一个token不是这个类型则自动回退，返回err
var cod int
var tok string
var eos bool

func (l *Lexer) ScanType(code int) (token string, err error) {

	cod, tok, eos = l.Scan()
	if code == cod {
		token += string(tok)
		return
	} else {
		err = errors.New("true")
	}
	return
}
