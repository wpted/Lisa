package repl

import (
	token "Lisa/lexToken"
	"Lisa/lexer"
	"bufio"
	"fmt"
	"io"
	"sync"
)

const PROMPT = "LISA >>> "

var lexerPool = sync.Pool{
	New: func() any { return new(lexer.Lexer) },
}

func Start(in io.Reader, out io.Reader) {
	scanner := bufio.NewScanner(in)

	for {
		fmt.Printf(PROMPT)
		scanned := scanner.Scan()
		if !scanned {
			return
		}

		line := scanner.Text()
		if line == "QUIT" {
			break
		}
		l := lexerPool.Get().(*lexer.Lexer)
		l = lexer.New(line)
		for tok := l.ReadNextToken(); tok.Type != token.EOF; tok = l.ReadNextToken() {
			fmt.Printf("%+v\n", tok)
		}

		l.Free()
		lexerPool.Put(l)
	}
}
