package stringcases

import (
	"fmt"
	"regexp"
	"strings"
	"unicode"

	"golang.org/x/text/cases"
	"golang.org/x/text/language"
)

var commonRe = regexp.MustCompile(`(ACL|API|ASCII|CPU|CSS|DNS|EOF|GUID|HTML|HTTP|HTTPS|ID|IP|JSON|LHS|QPS|RAM|RHS|RPC|SLA|SMTP|SQL|SSH|TCP|TLS|TTL|UDP|UI|UID|URI|URL|UTF8|UUID|VM|XML|XMPP|XSRF|XSS)`)
var repeatRe = regexp.MustCompile(`([^a-zA-Z]+)`)
var camelRe = regexp.MustCompile(`([a-z][A-Z])`)
var titleCaser = cases.Title(language.English, cases.NoLower)

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
	// Tokenize the string, but keeping the common initialism.
	s = commonRe.ReplaceAllStringFunc(s, func(s string) string {
		return fmt.Sprintf("_%s_", strings.ToLower(s))
	})

	// Splits aA to a_a.
	s = camelRe.ReplaceAllStringFunc(s, func(s string) string {
		return fmt.Sprintf("%s_%s", s[:1], strings.ToLower(s[1:]))
	})
	s = strings.ToLower(s)

	var sb strings.Builder
	defer sb.Reset()

	for i, r := range s {
		if unicode.IsNumber(r) || unicode.IsLetter(r) {
			sb.WriteRune(r)
		} else {
			if i != 0 && i != len(s)-1 {
				sb.WriteRune('_')
			}
		}
	}

	s = repeatRe.ReplaceAllStringFunc(sb.String(), func(s string) string {
		return s[:1]
	})

	return s
}

func ToKebab(s string) string {
	s = ToSnake(s)

	return strings.ReplaceAll(s, "_", "-")
}

func ToPascal(s string) string {
	s = ToSnake(s)

	words := strings.Split(s, "_")
	res := make([]string, len(words))
	for i, word := range words {
		if v := strings.ToUpper(word); commonInitialisms[v] {
			res[i] = v
		} else {
			res[i] = titleCaser.String(word)
		}
	}

	return strings.Join(res, "")
}

func ToCamel(s string) string {
	s = ToSnake(s)

	words := strings.Split(s, "_")
	res := make([]string, len(words))
	for i, word := range words {
		if i == 0 {
			res[i] = word
			continue
		}

		if v := strings.ToUpper(word); commonInitialisms[v] {
			res[i] = v
		} else {
			res[i] = titleCaser.String(word)
		}
	}

	return strings.Join(res, "")
}
