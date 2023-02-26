package n

import (
	"fmt"
	"gitea.voopsen/n/tree"
	"regexp"
	"strings"
)

func ParseRoute(raw string) (*regexp.Regexp, []tree.Matcher, error) {
	path := splitPath(raw)
	globalPath := make([]string, 0, len(path))

	matchers := make([]tree.Matcher, 0, len(path))
	for _, segment := range path {
		segment = strings.Trim(segment, " \t\n\r")
		if len(segment) == 0 {
			return nil, nil, fmt.Errorf(`bad route: "%s"`, raw)
		}

		if segment[0] == '{' && segment[len(segment)-1] == '}' {
			token := segment[1 : len(segment)-1]

			parts := strings.Split(token, ":")

			var name, regex string
			switch len(parts) {
			case 1:
				name = parts[0]
				regex = `[^/\\]+`
			case 2:
				name = parts[0]
				regex = parts[1]
			default:
				return nil, nil, fmt.Errorf(`bad variable: "%s"`, token)
			}

			matcher, err := createRegexpMatcher(regex)
			if err != nil {
				return nil, nil, fmt.Errorf(`bad regexp in "%s": %w"`, token, err)
			}

			matchers = append(matchers, matcher)
			globalPath = append(globalPath, `(?P<`+name+`>`+regex+`)`)
		} else {
			matchers = append(matchers, tree.NewFixedMatcher(segment))
			globalPath = append(globalPath, segment)
		}
	}

	globalRegex, err := regexp.Compile(strings.Join(globalPath, "/"))
	if err != nil {
		return nil, nil, err
	}

	return globalRegex, matchers, nil
}

func createRegexpMatcher(regex string) (tree.RegexpMatcher, error) {
	matchingRegex := regex
	if matchingRegex[0] != '^' {
		matchingRegex = "^" + matchingRegex
	}

	if matchingRegex[len(matchingRegex)-1] != '$' {
		matchingRegex = matchingRegex + "$"
	}

	return tree.NewRegexpMatcher(matchingRegex)
}

func splitPath(raw string) []string {
	if len(raw) == 0 {
		return nil
	}

	if raw[0] == '/' || raw[0] == '\\' {
		raw = raw[1:]
	}

	var result []string

	var (
		isVar   bool
		current []rune
	)

	for _, char := range []rune(raw) {
		switch char {
		case '/', '\\':
			if isVar {
				break
			}

			result = append(result, string(current))
			current = current[:0:len(current)]
			continue
		case '{':
			isVar = true

		case '}':
			isVar = false
		}

		current = append(current, char)
	}

	if len(current) > 0 {
		result = append(result, string(current))
	}

	return result
}
