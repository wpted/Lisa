package lexer

import (
	token "Lisa/lexToken"
)

// Lexer is an instance that is responsible for taking source code as input and output the tokens that represent it.
// It goes through it's preloaded input and output the token it recognize, token by token.
// TODO: To make debugging easier, initialize the Lexer with an io.Reader and the filename of the file the lexer is going through, which helps provide the line number, column number and the filename to a token.
type Lexer struct {
	// input is content that is parsed into the Lexer.
	input string
	// position is the current position in input (which points to current char ch).
	position int
	// readPosition is the current reading position in input (the next character, one after current char ch).
	readPosition int
	// ch is the current character under examination.
	// TODO: While using type byte supports ASCII (which can br compared as integers), UTF-8 is a must (change byte to rune and read it bytes wide).
	ch byte
}

// New creates a new pointer of Lexer.
func New(data string) *Lexer {
	newLexer := &Lexer{
		input:        data,
		position:     0,
		readPosition: 0,
		ch:           0,
	}

	// Initialize the lexer position.
	newLexer.readChar()

	return newLexer
}

// NextToken parses the input and transform it into a defined Token type.
// The function should be called repeatedly.
func (l *Lexer) NextToken() *token.Token {
	var isReserved bool
	var tok *token.Token

	// Skip white space before analyzing.
	l.skipWhiteSpace()

	switch l.ch {
	case '!':
		tok = token.New(token.EXCLAMATION, string(l.ch))
	case '=':
		tok = token.New(token.ASSIGN, string(l.ch))
	case '+':
		tok = token.New(token.PLUS, string(l.ch))
	case '-':
		tok = token.New(token.MINUS, string(l.ch))
	case '*':
		tok = token.New(token.ASTERISK, string(l.ch))
	case '^':
		tok = token.New(token.CARET, string(l.ch))
	case '/':
		tok = token.New(token.SLASH, string(l.ch))
	case '<':
		tok = token.New(token.LESSTHAN, string(l.ch))
	case '>':
		tok = token.New(token.GREATERTHAN, string(l.ch))
	case '(':
		tok = token.New(token.LPAREN, string(l.ch))
	case ')':
		tok = token.New(token.RPAREN, string(l.ch))
	case '{':
		tok = token.New(token.LBRACE, string(l.ch))
	case '}':
		tok = token.New(token.RBRACE, string(l.ch))
	case ';':
		tok = token.New(token.SEMICOLON, string(l.ch))
	case ',':
		tok = token.New(token.COMMA, string(l.ch))
	case 0:
		tok = token.New(token.EOF, "")
	default:
		// Not one of the recognized characters.
		// This is where the lexical reader encounters a letter.
		if isLetter(l.ch) {
			var lt token.LexicalType

			// Move the pointer and get the current identifier.
			identifier := l.readIdentifier()

			// Determine the type of the reading literal.
			lt, isReserved = token.LookUpReservedWord(identifier)
			tok = token.New(lt, identifier)
		} else if isDigit(l.ch) {
			number := l.readNumber()
			tok = token.New(token.INT, number)
		} else {
			tok = token.New(token.ILLEGAL, string(l.ch))
		}
	}

	// After checking token, move the lexical pointer to the next position if it's a reserved word, or when type is token.IDENT, token.INT.
	if isReserved || tok.Type == token.IDENT || tok.Type == token.INT {
		// Early exit here is necessary since the loop for readIdentifier or readNumber jumps one step forward, thus we don't need another jump below.
		return tok
	} else {
		l.readChar()
		return tok
	}
}

// readChar gives us the next character and advances our position in the input string.
// The current character is stored in ch.
// It checks whether it reached the end of the input.
// If we reached the end of the input, we set l.ch to 0, otherwise l.ch is set to the next character.
func (l *Lexer) readChar() {
	// Checking what's coming up next first.
	if l.readPosition >= len(l.input) {
		// EOF.
		l.ch = 0
	} else {
		// Assign the current reading char to ch.
		l.ch = l.input[l.readPosition]
	}

	// Whether read succeeds or not, move the read position to readPosition.
	l.position = l.readPosition
	// Point the read position to the next character in the input.
	l.readPosition++
}

// readIdentifier reads in an identifier and advances our lexers' positions until it encounters a non-letter character.
func (l *Lexer) readIdentifier() string {
	position := l.position

	// Check whether l.ch is a letter until reaches some character that isn't(possibly a white space or a delimiter).
	for isLetter(l.ch) {
		// Move l.position by reading through the characters one by one.
		l.readChar()
	}
	currIdent := l.input[position:l.position]
	return currIdent
}

// readNumber reads in a number and advances our lexers' positions until it encounters a non-number character.
func (l *Lexer) readNumber() string {
	position := l.position
	for isDigit(l.ch) {
		l.readChar()
	}
	currNumStr := l.input[position:l.position]
	return currNumStr
}

// isLetter determines whether an input character is a letter.
func isLetter(ch byte) bool {
	// For ASCII, letter a-z lies within [97, 122] and A-Z within [65, 90].
	// We treat "_" as letter(ASCII: 95), indicating we allow both Camel Case and Snake Case for names of variables or functions.
	if (65 <= ch && ch <= 90) || (97 <= ch && ch <= 122) || ch == 95 {
		return true
	}
	return false
}

// isDigit determines whether an input character is a number.
func isDigit(ch byte) bool {
	// For ASCII, letter 0-9 lies between [48, 57].
	if 48 <= ch && ch <= 57 {
		return true
	}
	return false
}

// skipWhiteSpace skips the encounter white space and advances the lexers pointer forward until it encounters a non-white-space character.
func (l *Lexer) skipWhiteSpace() {
	white := map[byte]struct{}{
		' ':  {},
		'\t': {},
		'\n': {},
		'\r': {},
	}
	_, ok := white[l.ch]
	for ok {
		l.readChar()
		_, ok = white[l.ch]
	}
}
