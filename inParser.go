package kvkv

import "fmt"

// InParse parses a slice as an Object.
func InParse(text []Token) (Object, error) {
	obj := Object{}
	for len(text) > 0 {
		if text[0].Type != TextTokenType {
			return obj, fmt.Errorf("object key non-text (%v)", text[0].Location)
		}
		entry := Entry{
			Key: text[0].Text,
			Location: text[0].Location,
		}
		if len(text) < 2 {
			return obj, fmt.Errorf("object key %v, no value (%v)", entry.Key, text[0].Location)
		}
		text = text[1:]
		if text[0].Type == TextTokenType {
			entry.Value = text[0].Text
			text = text[1:]
		} else if text[0].Type == OpenTokenType {
			// sub-object! oh, what "fun"
			// now we have to use a second loop to scan for the end
			counter := 0
			for k, v := range text {
				if v.Type == CloseTokenType {
					counter--
					if counter == 0 {
						val, err := InParse(text[1:k])
						if err != nil {
							return obj, fmt.Errorf("(%s @ %v): %s", entry.Key, entry.Location, err)
						}
						entry.Value = val
						text = text[k + 1:]
						break
					}
				} else if v.Type == OpenTokenType {
					counter++
				}
			}
			if entry.Value == nil {
				return obj, fmt.Errorf("sub-object without end (%v)", text[0].Location)
			}
		} else {
			return obj, fmt.Errorf("object value start token type invalid (%v)", text[0].Location)
		}
		obj = append(obj, entry)
	}
	return obj, nil
}

