package tree

import "regexp"

type Matcher interface {
	String() string
	MatchString(s string) bool
}

type FixedMatcher struct {
	s string
}

func NewFixedMatcher(s string) FixedMatcher {
	return FixedMatcher{
		s: s,
	}
}

func (f FixedMatcher) String() string {
	return f.s
}

func (f FixedMatcher) MatchString(s string) bool {
	return s == f.s
}

type RegexpMatcher struct {
	regexp *regexp.Regexp
}

func NewRegexpMatcher(regex string) (RegexpMatcher, error) {
	r, err := regexp.Compile(regex)
	if err != nil {
		return RegexpMatcher{}, err
	}

	return RegexpMatcher{
		regexp: r,
	}, nil
}

func (r RegexpMatcher) String() string {
	return r.regexp.String()
}

func (r RegexpMatcher) MatchString(s string) bool {
	return r.regexp.MatchString(s)
}
