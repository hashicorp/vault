package matching

import (
	"github.com/nbutton23/zxcvbn-go/entropy"
	"github.com/nbutton23/zxcvbn-go/match"
	"strings"
)

func buildDictMatcher(dictName string, rankedDict map[string]int) func(password string) []match.Match {
	return func(password string) []match.Match {
		matches := dictionaryMatch(password, dictName, rankedDict)
		for _, v := range matches {
			v.DictionaryName = dictName
		}
		return matches
	}

}

func dictionaryMatch(password string, dictionaryName string, rankedDict map[string]int) []match.Match {
	length := len(password)
	var results []match.Match
	pwLower := strings.ToLower(password)

	for i := 0; i < length; i++ {
		for j := i; j < length; j++ {
			word := pwLower[i : j+1]
			if val, ok := rankedDict[word]; ok {
				matchDic := match.Match{Pattern: "dictionary",
					DictionaryName: dictionaryName,
					I:              i,
					J:              j,
					Token:          password[i : j+1],
				}
				matchDic.Entropy = entropy.DictionaryEntropy(matchDic, float64(val))

				results = append(results, matchDic)
			}
		}
	}

	return results
}

func buildRankedDict(unrankedList []string) map[string]int {

	result := make(map[string]int)

	for i, v := range unrankedList {
		result[strings.ToLower(v)] = i + 1
	}

	return result
}
