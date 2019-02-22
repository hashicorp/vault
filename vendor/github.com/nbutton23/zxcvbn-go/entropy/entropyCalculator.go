package entropy

import (
	"github.com/nbutton23/zxcvbn-go/adjacency"
	"github.com/nbutton23/zxcvbn-go/match"
	"github.com/nbutton23/zxcvbn-go/utils/math"
	"math"
	"regexp"
	"unicode"
)

const (
	START_UPPER string = `^[A-Z][^A-Z]+$`
	END_UPPER   string = `^[^A-Z]+[A-Z]$'`
	ALL_UPPER   string = `^[A-Z]+$`
	NUM_YEARS          = float64(119) // years match against 1900 - 2019
	NUM_MONTHS         = float64(12)
	NUM_DAYS           = float64(31)
)

var (
	KEYPAD_STARTING_POSITIONS = len(adjacency.AdjacencyGph["keypad"].Graph)
	KEYPAD_AVG_DEGREE         = adjacency.AdjacencyGph["keypad"].CalculateAvgDegree()
)

func DictionaryEntropy(match match.Match, rank float64) float64 {
	baseEntropy := math.Log2(rank)
	upperCaseEntropy := extraUpperCaseEntropy(match)
	//TODO: L33t
	return baseEntropy + upperCaseEntropy
}

func extraUpperCaseEntropy(match match.Match) float64 {
	word := match.Token

	allLower := true

	for _, char := range word {
		if unicode.IsUpper(char) {
			allLower = false
			break
		}
	}
	if allLower {
		return float64(0)
	}

	//a capitalized word is the most common capitalization scheme,
	//so it only doubles the search space (uncapitalized + capitalized): 1 extra bit of entropy.
	//allcaps and end-capitalized are common enough too, underestimate as 1 extra bit to be safe.

	for _, regex := range []string{START_UPPER, END_UPPER, ALL_UPPER} {
		matcher := regexp.MustCompile(regex)

		if matcher.MatchString(word) {
			return float64(1)
		}
	}
	//Otherwise calculate the number of ways to capitalize U+L uppercase+lowercase letters with U uppercase letters or
	//less. Or, if there's more uppercase than lower (for e.g. PASSwORD), the number of ways to lowercase U+L letters
	//with L lowercase letters or less.

	countUpper, countLower := float64(0), float64(0)
	for _, char := range word {
		if unicode.IsUpper(char) {
			countUpper++
		} else if unicode.IsLower(char) {
			countLower++
		}
	}
	totalLenght := countLower + countUpper
	var possibililities float64

	for i := float64(0); i <= math.Min(countUpper, countLower); i++ {
		possibililities += float64(zxcvbn_math.NChoseK(totalLenght, i))
	}

	if possibililities < 1 {
		return float64(1)
	}

	return float64(math.Log2(possibililities))
}

func SpatialEntropy(match match.Match, turns int, shiftCount int) float64 {
	var s, d float64
	if match.DictionaryName == "qwerty" || match.DictionaryName == "dvorak" {
		//todo: verify qwerty and dvorak have the same length and degree
		s = float64(len(adjacency.BuildQwerty().Graph))
		d = adjacency.BuildQwerty().CalculateAvgDegree()
	} else {
		s = float64(KEYPAD_STARTING_POSITIONS)
		d = KEYPAD_AVG_DEGREE
	}

	possibilities := float64(0)

	length := float64(len(match.Token))

	//TODO: Should this be <= or just < ?
	//Estimate the number of possible patterns w/ length L or less with t turns or less
	for i := float64(2); i <= length+1; i++ {
		possibleTurns := math.Min(float64(turns), i-1)
		for j := float64(1); j <= possibleTurns+1; j++ {
			x := zxcvbn_math.NChoseK(i-1, j-1) * s * math.Pow(d, j)
			possibilities += x
		}
	}

	entropy := math.Log2(possibilities)
	//add extra entropu for shifted keys. ( % instead of 5 A instead of a)
	//Math is similar to extra entropy for uppercase letters in dictionary matches.

	if S := float64(shiftCount); S > float64(0) {
		possibilities = float64(0)
		U := length - S

		for i := float64(0); i < math.Min(S, U)+1; i++ {
			possibilities += zxcvbn_math.NChoseK(S+U, i)
		}

		entropy += math.Log2(possibilities)
	}

	return entropy
}

func RepeatEntropy(match match.Match) float64 {
	cardinality := CalcBruteForceCardinality(match.Token)
	entropy := math.Log2(cardinality * float64(len(match.Token)))

	return entropy
}

//TODO: Validate against python
func CalcBruteForceCardinality(password string) float64 {
	lower, upper, digits, symbols := float64(0), float64(0), float64(0), float64(0)

	for _, char := range password {
		if unicode.IsLower(char) {
			lower = float64(26)
		} else if unicode.IsDigit(char) {
			digits = float64(10)
		} else if unicode.IsUpper(char) {
			upper = float64(26)
		} else {
			symbols = float64(33)
		}
	}

	cardinality := lower + upper + digits + symbols
	return cardinality
}

func SequenceEntropy(match match.Match, dictionaryLength int, ascending bool) float64 {
	firstChar := match.Token[0]
	baseEntropy := float64(0)
	if string(firstChar) == "a" || string(firstChar) == "1" {
		baseEntropy = float64(0)
	} else {
		baseEntropy = math.Log2(float64(dictionaryLength))
		//TODO: should this be just the first or any char?
		if unicode.IsUpper(rune(firstChar)) {
			baseEntropy++
		}
	}

	if !ascending {
		baseEntropy++
	}
	return baseEntropy + math.Log2(float64(len(match.Token)))
}

func ExtraLeetEntropy(match match.Match, password string) float64 {
	var subsitutions float64
	var unsub float64
	subPassword := password[match.I:match.J]
	for index, char := range subPassword {
		if string(char) != string(match.Token[index]) {
			subsitutions++
		} else {
			//TODO: Make this only true for 1337 chars that are not subs?
			unsub++
		}
	}

	var possibilities float64

	for i := float64(0); i <= math.Min(subsitutions, unsub)+1; i++ {
		possibilities += zxcvbn_math.NChoseK(subsitutions+unsub, i)
	}

	if possibilities <= 1 {
		return float64(1)
	}
	return math.Log2(possibilities)
}

func YearEntropy(dateMatch match.DateMatch) float64 {
	return math.Log2(NUM_YEARS)
}

func DateEntropy(dateMatch match.DateMatch) float64 {
	var entropy float64
	if dateMatch.Year < 100 {
		entropy = math.Log2(NUM_DAYS * NUM_MONTHS * 100)
	} else {
		entropy = math.Log2(NUM_DAYS * NUM_MONTHS * NUM_YEARS)
	}

	if dateMatch.Separator != "" {
		entropy += 2 //add two bits for separator selection [/,-,.,etc]
	}
	return entropy
}
