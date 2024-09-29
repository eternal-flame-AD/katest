package main

import "strings"

var aliases = [][]string{
	{"fu", "hu"},
	{"ti", "chi"},
	{"tu", "tsu"},
	{"si", "shi"},
	{"zi", "ji"},
	{"di", "ji"},
	{"du", "zu"},
	{"nn", "n"},
}

func aliasNormalize(s string) string {
	for _, alias := range aliases {
		matches := false
		for _, a := range alias {
			if strings.ToLower(s) == a {
				matches = true
				break
			}
		}
		if matches {
			return alias[0]
		}
	}
	return s
}

func charMatch(ref Char, input string, requireType bool) bool {
	t := ""
	if colon := strings.Index(input, ":"); colon != -1 {
		t = input[colon+1:]
		input = input[:colon]
		return strings.ToUpper(aliasNormalize(input)) == strings.ToUpper(aliasNormalize(ref.Romaji)) &&
			strings.HasPrefix(ref.Type, strings.ToLower(t))
	}
	return !requireType &&
		strings.ToUpper(aliasNormalize(input)) == strings.ToUpper(aliasNormalize(ref.Romaji))
}

func locateChar(s string, chars []Char) Char {
	for _, c := range chars {
		if c.Romaji == strings.ToUpper(s) {
			return c
		}
	}
	panic("character " + s + " not found")
}

func stringsContains(s string, list []string) bool {
	for _, m := range list {
		if m == s {
			return true
		}
	}
	return false
}

func filter[T comparable](s []T, f func(T) bool) []T {
	var r []T
	for _, v := range s {
		if f(v) {
			r = append(r, v)
		}
	}
	return r
}

func resolveCharsets(desc string) (res []Char) {
	charsetDescs := strings.Split(desc, ",")
	for _, charset := range charsetDescs {
		enableHiragana, enableKatakana := true, true
		if colon := strings.Index(charset, ":"); colon != -1 {
			var sy string
			charset, sy = charset[:colon], charset[colon+1:]
			switch sy {
			case "h":
				enableKatakana = false
			case "k":
				enableHiragana = false
			default:
				panic("unknown writing system")
			}
		}
		switch charset {
		case "50":
			for _, vowel := range []string{"a", "i", "u", "e", "o"} {
				for _, constant := range []string{"", "k", "s", "t", "n", "h", "m", "y", "r", "w"} {
					s := constant + vowel
					if stringsContains(s, []string{"yi", "ye", "wi", "wu", "we"}) {
						continue
					}
					if enableHiragana {
						res = append(res, locateChar(s, hiraganaChars))
					}
					if enableKatakana {
						res = append(res, locateChar(s, katakanaChars))
					}
				}
			}
		case "all":
			res = append(res, hiraganaChars...)
			res = append(res, katakanaChars...)
		default:
			chars := strings.Split(charset, "+")
			for _, s := range chars {
				if enableHiragana {
					res = append(res, locateChar(s, hiraganaChars))
				}
				if enableKatakana {
					res = append(res, locateChar(s, katakanaChars))
				}
			}
		}
	}
	return
}
