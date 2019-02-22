package matching

import (
	"sort"

	"github.com/nbutton23/zxcvbn-go/adjacency"
	"github.com/nbutton23/zxcvbn-go/frequency"
	"github.com/nbutton23/zxcvbn-go/match"
)

var (
	DICTIONARY_MATCHERS []match.Matcher
	MATCHERS            []match.Matcher
	ADJACENCY_GRAPHS    []adjacency.AdjacencyGraph
	L33T_TABLE          adjacency.AdjacencyGraph

	SEQUENCES map[string]string
)

const (
	DATE_RX_YEAR_SUFFIX    string = `((\d{1,2})(\s|-|\/|\\|_|\.)(\d{1,2})(\s|-|\/|\\|_|\.)(19\d{2}|200\d|201\d|\d{2}))`
	DATE_RX_YEAR_PREFIX    string = `((19\d{2}|200\d|201\d|\d{2})(\s|-|/|\\|_|\.)(\d{1,2})(\s|-|/|\\|_|\.)(\d{1,2}))`
	DATE_WITHOUT_SEP_MATCH string = `\d{4,8}`
)

func init() {
	loadFrequencyList()
}

func Omnimatch(password string, userInputs []string, filters ...func(match.Matcher) bool) (matches []match.Match) {

	//Can I run into the issue where nil is not equal to nil?
	if DICTIONARY_MATCHERS == nil || ADJACENCY_GRAPHS == nil {
		loadFrequencyList()
	}

	if userInputs != nil {
		userInputMatcher := buildDictMatcher("user_inputs", buildRankedDict(userInputs))
		matches = userInputMatcher(password)
	}

	for _, matcher := range MATCHERS {
		shouldBeFiltered := false
		for i := range filters {
			if filters[i](matcher) {
				shouldBeFiltered = true
				break
			}
		}
		if !shouldBeFiltered {
			matches = append(matches, matcher.MatchingFunc(password)...)
		}
	}
	sort.Sort(match.Matches(matches))
	return matches
}

func loadFrequencyList() {

	for n, list := range frequency.FrequencyLists {
		DICTIONARY_MATCHERS = append(DICTIONARY_MATCHERS, match.Matcher{MatchingFunc: buildDictMatcher(n, buildRankedDict(list.List)), ID: n})
	}

	L33T_TABLE = adjacency.AdjacencyGph["l33t"]

	ADJACENCY_GRAPHS = append(ADJACENCY_GRAPHS, adjacency.AdjacencyGph["qwerty"])
	ADJACENCY_GRAPHS = append(ADJACENCY_GRAPHS, adjacency.AdjacencyGph["dvorak"])
	ADJACENCY_GRAPHS = append(ADJACENCY_GRAPHS, adjacency.AdjacencyGph["keypad"])
	ADJACENCY_GRAPHS = append(ADJACENCY_GRAPHS, adjacency.AdjacencyGph["macKeypad"])

	//l33tFilePath, _ := filepath.Abs("adjacency/L33t.json")
	//L33T_TABLE = adjacency.GetAdjancencyGraphFromFile(l33tFilePath, "l33t")

	SEQUENCES = make(map[string]string)
	SEQUENCES["lower"] = "abcdefghijklmnopqrstuvwxyz"
	SEQUENCES["upper"] = "ABCDEFGHIJKLMNOPQRSTUVWXYZ"
	SEQUENCES["digits"] = "0123456789"

	MATCHERS = append(MATCHERS, DICTIONARY_MATCHERS...)
	MATCHERS = append(MATCHERS, match.Matcher{MatchingFunc: spatialMatch, ID: SPATIAL_MATCHER_NAME})
	MATCHERS = append(MATCHERS, match.Matcher{MatchingFunc: repeatMatch, ID: REPEAT_MATCHER_NAME})
	MATCHERS = append(MATCHERS, match.Matcher{MatchingFunc: sequenceMatch, ID: SEQUENCE_MATCHER_NAME})
	MATCHERS = append(MATCHERS, match.Matcher{MatchingFunc: l33tMatch, ID: L33T_MATCHER_NAME})
	MATCHERS = append(MATCHERS, match.Matcher{MatchingFunc: dateSepMatcher, ID: DATESEP_MATCHER_NAME})
	MATCHERS = append(MATCHERS, match.Matcher{MatchingFunc: dateWithoutSepMatch, ID: DATEWITHOUTSEP_MATCHER_NAME})

}
