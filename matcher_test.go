package n

import (
	"gitea.voopsen/n/tree"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"regexp"
	"testing"
)

func Test_splitPath(t *testing.T) {
	for _, data := range []struct {
		source string
		result []string
	}{
		{
			source: `users/list`,
			result: []string{"users", "list"},
		},
		{
			source: `users\list`,
			result: []string{"users", "list"},
		},
		{
			source: `users\{var:\d+\w+}`,
			result: []string{"users", `{var:\d+\w+}`},
		},
		{
			source: `users\{var:[\/]+}`,
			result: []string{"users", `{var:[\/]+}`},
		},
	} {
		assert.EqualValues(t, data.result, splitPath(data.source))
		assert.EqualValues(t, data.result, splitPath(`/`+data.source))
		assert.EqualValues(t, data.result, splitPath(data.source+`/`))
		assert.EqualValues(t, data.result, splitPath(`/`+data.source+`/`))
	}
	assert.Nil(t, splitPath(""))
}

type parseRawCase struct {
	raw      string
	global   *regexp.Regexp
	matchers []tree.Matcher
	err      error
	anyErr   bool
}

func mustNewRegexpMatcher(regex string) tree.RegexpMatcher {
	m, _ := tree.NewRegexpMatcher(regex)
	return m
}

var (
	_parseRawTestCases = map[string]parseRawCase{
		"bad/segment/empty": {
			raw:    "users//all",
			anyErr: true,
		},
		"bad/vars/incorrect": {
			raw:    `users/{id:\d+:\w+}/all`,
			anyErr: true,
		},
		"bad/vars/bad-regexp": {
			raw:    `users/{id:[}/all`,
			anyErr: true,
		},
		"bad/fixed/bad-global": {
			raw:    `users/[all`,
			anyErr: true,
		},
		"fixed": {
			raw:    "users/all",
			global: regexp.MustCompile("users/all"),
			matchers: []tree.Matcher{
				tree.NewFixedMatcher("users"),
				tree.NewFixedMatcher("all"),
			},
		},
		"vars/default": {
			raw:    "users/{id}/scripts",
			global: regexp.MustCompile(`users/(?P<id>[^/\\]+)/scripts`),
			matchers: []tree.Matcher{
				tree.NewFixedMatcher("users"),
				mustNewRegexpMatcher("^[^/\\\\]+$"),
				tree.NewFixedMatcher("scripts"),
			},
		},
		"vars/decimals": {
			raw:    `users/{id:\d+}/scripts`,
			global: regexp.MustCompile(`users/(?P<id>\d+)/scripts`),
			matchers: []tree.Matcher{
				tree.NewFixedMatcher("users"),
				mustNewRegexpMatcher(`^\d+$`),
				tree.NewFixedMatcher("scripts"),
			},
		},
		"vars/many": {
			raw:    `users/{id:\d+}/scripts/{scriptId:[a-zA-Z]+}/list`,
			global: regexp.MustCompile(`users/(?P<id>\d+)/scripts/(?P<scriptId>[a-zA-Z]+)/list`),
			matchers: []tree.Matcher{
				tree.NewFixedMatcher("users"),
				mustNewRegexpMatcher(`^\d+$`),
				tree.NewFixedMatcher("scripts"),
				mustNewRegexpMatcher(`^[a-zA-Z]+$`),
				tree.NewFixedMatcher("list"),
			},
		},
	}
)

func Test_ParseRoute(t *testing.T) {
	for name, testCase := range _parseRawTestCases {
		t.Run(name, func(t *testing.T) {
			global, matchers, err := ParseRoute(testCase.raw)

			if testCase.anyErr {
				require.NotNil(t, err)
			} else if testCase.err == nil {
				require.NoError(t, err)
			} else {
				require.Error(t, testCase.err, err)
				return
			}

			for i, m := range testCase.matchers {
				assert.Equal(t, m.String(), matchers[i].String())
			}

			if testCase.global == nil {
				require.Nil(t, global)
				return
			}

			require.NotNil(t, global)
			assert.Equal(t, testCase.global.String(), global.String())
		})
	}
}
