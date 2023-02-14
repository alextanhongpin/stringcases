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
	titleCaser = cases.Title(language.English, cases.NoLower)
	lowerCaser = cases.Lower(language.English)
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

func ToSnake(s string) string {
	tokens := tokenize(s)
	res := make([]string, len(tokens))
	for i, token := range tokens {
		res[i] = strings.ToLower(token)
	}

	return strings.Join(res, "_")
}

func ToKebab(s string) string {
	tokens := tokenize(s)
	res := make([]string, len(tokens))
	for i, token := range tokens {
		res[i] = strings.ToLower(token)
	}

	return strings.Join(res, "-")
}

func ToCamel(s string) string {
	tokens := tokenize(s)
	res := make([]string, len(tokens))
	for i, token := range tokens {
		if i == 0 {
			res[i] = lowerCaser.String(token)
			continue
		}

		u := strings.ToUpper(token)
		if commonInitialisms[u] {
			res[i] = u
		} else {
			res[i] = titleCaser.String(token)
		}
	}

	return strings.Join(res, "")
}

func ToPascal(s string) string {
	tokens := tokenize(s)
	res := make([]string, len(tokens))
	for i, token := range tokens {
		u := strings.ToUpper(token)
		if commonInitialisms[u] {
			res[i] = u
		} else {
			res[i] = titleCaser.String(token)
		}
	}

	return strings.Join(res, "")
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
		case
			unicode.IsNumber(r),
			unicode.IsLower(r):
			if err := reader.UnreadRune(); err != nil {
				panic(err)
			}

			token := extractLower(reader)
			tokens = append(tokens, string(token))
		case unicode.IsUpper(r):
			if err := reader.UnreadRune(); err != nil {
				panic(err)
			}

			token := extractUpper(reader)
			tokens = append(tokens, string(token))
		default:
		}
	}

	return tokens
}

func extractUpper(reader *strings.Reader) []rune {
	var res []rune
	var isProbablyCommonInitialism bool

	for {
		r, _, err := reader.ReadRune()
		if errors.Is(err, io.EOF) {
			break
		}
		switch {
		case unicode.IsUpper(r):
			// The first and second character is uppercase.
			if len(res) == 1 {
				isProbablyCommonInitialism = true
			}

			// Non-common initialism pattern breaks on the next uppercase rune.
			if len(res) > 1 && !isProbablyCommonInitialism {
				if err := reader.UnreadRune(); err != nil {
					panic(err)
				}

				return res
			}
		case unicode.IsLower(r), unicode.IsNumber(r):
			// Common initialism pattern breaks on the next lowercase rune.
			if isProbablyCommonInitialism {
				if err := reader.UnreadRune(); err != nil {
					panic(err)
				}

				return res
			}
		default:
			return res
		}

		res = append(res, r)
		if len(res) >= 2 && len(res) <= 5 {
			if commonInitialisms[string(res)] {
				return res
			}
		}
	}

	return res
}

func extractLower(reader *strings.Reader) []rune {
	var res []rune

	for {
		r, _, err := reader.ReadRune()
		if errors.Is(err, io.EOF) {
			break
		}

		switch {
		case unicode.IsUpper(r):
			if err := reader.UnreadRune(); err != nil {
				panic(err)
			}

			return res
		case unicode.IsLower(r), unicode.IsNumber(r):
			res = append(res, r)
		default:
			return res
		}
	}

	return res
}
