package stringcases

import (
	"errors"
	"io"
	"strings"
	"unicode"

	"golang.org/x/text/cases"
	"golang.org/x/text/language"
)

var (
	s        = New(language.English)
	ToKebab  = s.ToKebab
	ToCamel  = s.ToCamel
	ToSnake  = s.ToSnake
	ToPascal = s.ToPascal
)

// https://github.com/golang/lint/blob/6edffad5e6160f5949cdefc81710b2706fbcd4f6/lint.go#LL766-L809
// commonInitialisms is a set of common initialisms.
// Only add entries that are highly unlikely to be non-initialisms.
// For instance, "ID" is fine (Freudian code is rare), but "AND" is not.
var commonInitialisms = map[string]bool{
	"ACL":   true,
	"API":   true,
	"ASCII": true,
	"CPU":   true,
	"CSS":   true,
	"DNS":   true,
	"EOF":   true,
	"GUID":  true,
	"HTML":  true,
	"HTTP":  true,
	"HTTPS": true,
	"ID":    true,
	"IP":    true,
	"JSON":  true,
	"LHS":   true,
	"QPS":   true,
	"RAM":   true,
	"RHS":   true,
	"RPC":   true,
	"SLA":   true,
	"SMTP":  true,
	"SQL":   true,
	"SSH":   true,
	"TCP":   true,
	"TLS":   true,
	"TTL":   true,
	"UDP":   true,
	"UI":    true,
	"UID":   true,
	"UUID":  true,
	"URI":   true,
	"URL":   true,
	"UTF8":  true,
	"VM":    true,
	"XML":   true,
	"XMPP":  true,
	"XSRF":  true,
	"XSS":   true,
}

type String struct {
	uppercase, lowercase, titlecase cases.Caser
}

func New(t language.Tag) *String {
	return &String{
		titlecase: cases.Title(t, cases.NoLower),
		lowercase: cases.Lower(t),
		uppercase: cases.Upper(t),
	}
}

func (str *String) ToSnake(s string) string {
	tokens := tokenize(s)
	runes := make([]string, len(tokens))
	for i, token := range tokens {
		runes[i] = str.lowercase.String(token)
	}

	return strings.Join(runes, "_")
}

func (str *String) ToKebab(s string) string {
	tokens := tokenize(s)
	runes := make([]string, len(tokens))
	for i, token := range tokens {
		runes[i] = str.lowercase.String(token)
	}

	return strings.Join(runes, "-")
}

func (str *String) ToCamel(s string) string {
	tokens := tokenize(s)
	runes := make([]string, len(tokens))
	for i, token := range tokens {
		if i == 0 {
			runes[i] = str.lowercase.String(token)
			continue
		}

		u := str.uppercase.String(token)
		if commonInitialisms[u] {
			runes[i] = u
		} else {
			runes[i] = str.titlecase.String(token)
		}
	}

	return strings.Join(runes, "")
}

func (str *String) ToPascal(s string) string {
	tokens := tokenize(s)
	runes := make([]string, len(tokens))
	for i, token := range tokens {
		u := str.uppercase.String(token)
		if commonInitialisms[u] {
			runes[i] = u
		} else {
			runes[i] = str.titlecase.String(token)
		}
	}

	return strings.Join(runes, "")
}

func tokenize(s string) []string {
	var tokens []string

	reader := strings.NewReader(s)
	for {
		r, _, err := reader.ReadRune()
		if errors.Is(err, io.EOF) {
			break
		}

		switch {
		case unicode.IsNumber(r), unicode.IsLower(r):
			token := extractLower(reader, []rune{r})
			tokens = append(tokens, token)

		case unicode.IsUpper(r):
			token := extractUpper(reader, []rune{r})
			tokens = append(tokens, token)

		default:
			// Skip non-alphanumeric runes.
		}
	}

	return tokens
}

func extractUpper(reader *strings.Reader, runes []rune) string {
	for {
		r, _, err := reader.ReadRune()
		if errors.Is(err, io.EOF) {
			return string(runes)
		}

		switch {
		case unicode.IsUpper(r):
			// Continuous upper unicode indicates the possibility of common
			// initialism word.
			return extractCommonInitialism(reader, append(runes, r))
		case unicode.IsLower(r), unicode.IsNumber(r):
			// Otherwise, it will be camel case word.
			return extractCamel(reader, append(runes, r))
		default:
			// Word breaks when it is non-alphanumeric.
			return string(runes)
		}
	}
}

func extractLower(reader *strings.Reader, runes []rune) string {
	for {
		r, _, err := reader.ReadRune()
		if errors.Is(err, io.EOF) {
			return string(runes)
		}

		switch {
		case unicode.IsUpper(r):
			// Word breaks when the next character is upper.
			if err := reader.UnreadRune(); err != nil {
				panic(err)
			}

			return string(runes)
		case unicode.IsLower(r), unicode.IsNumber(r):
			runes = append(runes, r)
		default:
			// Word breaks when it is non-alphanumeric.
			return string(runes)
		}
	}
}

func extractCommonInitialism(reader *strings.Reader, runes []rune) string {
	for {
		r, _, err := reader.ReadRune()
		if errors.Is(err, io.EOF) {
			return string(runes)
		}

		switch {
		case unicode.IsUpper(r):
			runes = append(runes, r)
			// Common initialism at present has length between 2 and 5.
			if len(runes) >= 2 && len(runes) <= 5 {
				if commonInitialisms[string(runes)] {
					return string(runes)
				}
			}
		// Common initialism pattern breaks at the next lower or number.
		case unicode.IsLower(r), unicode.IsNumber(r):
			if err := reader.UnreadRune(); err != nil {
				panic(err)
			}
			return string(runes)
		default:
			return string(runes)
		}
	}
}

func extractCamel(reader *strings.Reader, runes []rune) string {
	for {
		r, _, err := reader.ReadRune()
		if errors.Is(err, io.EOF) {
			return string(runes)
		}

		switch {
		case unicode.IsUpper(r):
			if err := reader.UnreadRune(); err != nil {
				panic(err)
			}

			return string(runes)
		case unicode.IsLower(r), unicode.IsNumber(r):
			runes = append(runes, r)
		default:
			return string(runes)
		}
	}
}
