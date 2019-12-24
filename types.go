package kvkv

import "fmt"

// Location describes a location in the file, for debugging
type Location struct {
	Filename string
	Line int
}

func (l Location) String() string {
	if l.Filename != "" {
		return fmt.Sprintf("%v:%v", l.Filename, l.Line)
	}
	return fmt.Sprintf("line %v", l.Line)
}

// Value is either a ([]Object) or a string.
type Value interface {}

// Object is an object (in particular, this is the root node of parsable documents)
type Object []Entry

// Find finds a specific key
func (o Object) Find(name string) (Entry, error) {
	for _, v := range o {
		if v.Key == name {
			return v, nil
		}
	}
	return Entry{}, fmt.Errorf("unable to find %s", name)
}

// FindString is equivalent to Find followed by ValueString
func (o Object) FindString(name string) (string, error) {
	val, err := o.Find(name)
	if err != nil {
		return "", err
	}
	return val.ValueString()
}

// FindObject is equivalent to Find followed by ValueObject
func (o Object) FindObject(name string) (Object, error) {
	val, err := o.Find(name)
	if err != nil {
		return nil, err
	}
	return val.ValueObject()
}

// Entry is a key in a Map.
type Entry struct {
	Key string
	Value Value
	Location
}

// ValueString attempts conversion to string w/ error return
func (e Entry) ValueString() (string, error) {
	switch res := e.Value.(type) {
		case string:
			return res, nil
	}
	return "", fmt.Errorf("asserted %s as string, value not string", e.Key)
}

// ValueObject attempts conversion to object w/ error return
func (e Entry) ValueObject() (Object, error) {
	switch res := e.Value.(type) {
		case Object:
			return res, nil
	}
	return nil, fmt.Errorf("asserted %s as object, value not object", e.Key)
}

// TokenType represents a type of token. Linters represent annoyance.
type TokenType int

// Please be aware:
// 1. comments are not tokens
// 2. text is *after* handling of escapes and quotes
const (
	// TextTokenType represents a piece of text (quoted or unquoted); the quotes are not in the resulting token text.
	TextTokenType TokenType = iota
	// OpenTokenType represents the { character outside of quotes
	OpenTokenType
	// CloseTokenType represents the } character outside of quotes
	CloseTokenType
)

type Token struct {
	Text string
	Type TokenType
	Location
}

