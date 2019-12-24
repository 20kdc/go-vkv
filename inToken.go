package kvkv

import "fmt"

// InTokenize tokenizes a string.
func InTokenize(text string, handleEscapes bool, filename string) ([]Token, error) {
	textRunes := []rune(text)
	// these are the state machine variables. they are annoyingly intricate.
	state := 0
	startOfTokenLocation := Location{
		Filename: filename,
	}
	line := 1
	cursor := 0
	// used to detect backtracking so that line tracking doesn't go out of sync
	highestCursor := -1
	tokens := []Token{}
	hasHandledEOF := false
	// also for building unquoted, too; this used to use rune slices, didn't work out
	buildingQuoted := ""
	for !hasHandledEOF {
		ch := -1
		if cursor < len(textRunes) {
			ch = int(textRunes[cursor])
			if cursor > highestCursor {
				highestCursor = cursor
				if ch == 10 {
					line++
				}
			}
			cursor++
		} else {
			hasHandledEOF = true
		}
		switch state {
			case 0:
				// Whitespace Eater
				startOfTokenLocation.Line = line
				if ch <= 32 {
					// Whitespace or EOF
				} else if ch == '/' {
					// Either it's a comment or it doesn't make sense
					state = 1
				} else if ch == '{' {
					// Open token
					tokens = append(tokens, Token{
						Type: OpenTokenType,
						Text: "{",
						Location: startOfTokenLocation,
					})
				} else if ch == '}' {
					// Close token
					tokens = append(tokens, Token{
						Type: CloseTokenType,
						Text: "}",
						Location: startOfTokenLocation,
					})
				} else if ch == '"' {
					// quoted ID
					state = 4
					buildingQuoted = ""
				} else {
					// unquoted ID ; this action reparses the char
					state = 3
					cursor--
					buildingQuoted = ""
				}
			case 1:
				// Comment - Second character determination
				if ch == '/' {
					state = 2
				} else {
					// Treat as if it was always an unquoted ID
					state = 4
				}
			case 2:
				// Comment - Wait until end of line
				if ch == 10 {
					state = 0
				}
			case 3:
				// Unquoted ID
				if ch <= 32 || ch == '"' || ch == '{' || ch == '}' {
					// end of ID
					tokens = append(tokens, Token{
						Type: TextTokenType,
						Text: buildingQuoted,
						Location: startOfTokenLocation,
					})
					// re-parse this character
					state = 0
					cursor -= 1
				} else {
					buildingQuoted += string(rune(ch))
				}
			case 4:
				// Quoted ID - Main
				if ch == '\\' && handleEscapes {
					state = 5
				} else if ch == '"' {
					state = 0
					tokens = append(tokens, Token{
						Type: TextTokenType,
						Text: buildingQuoted,
						Location: startOfTokenLocation,
					})
				} else if ch == -1 {
					// that wasn't supposed to happen!
					return tokens, fmt.Errorf("hit EOF during quoted string: %s", buildingQuoted)
				} else {
					buildingQuoted += string(rune(ch))
				}
			case 5:
				// always
				state = 4
				// Quoted ID - Escape
				if ch == 'n' {
					buildingQuoted += "\n"
				} else if ch == 't' {
					buildingQuoted += "\t"
				} else if ch == '\\' {
					buildingQuoted += "\\"
				} else if ch == '"' {
					buildingQuoted += "\""
				} else if ch == -1 {
					// that wasn't supposed to happen!
					return tokens, fmt.Errorf("hit EOF during quoted string escape: %s", buildingQuoted)
				} else {
					// just in case anyone gets any funny ideas; not necessarily valid!
					buildingQuoted += string(rune(ch))
				}
		}
	}
	return tokens, nil
}
