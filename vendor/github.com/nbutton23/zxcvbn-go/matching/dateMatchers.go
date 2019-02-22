package matching

import (
	"regexp"
	"strconv"
	"strings"

	"github.com/nbutton23/zxcvbn-go/entropy"
	"github.com/nbutton23/zxcvbn-go/match"
)

const (
	DATESEP_MATCHER_NAME        = "DATESEP"
	DATEWITHOUTSEP_MATCHER_NAME = "DATEWITHOUT"
)

func FilterDateSepMatcher(m match.Matcher) bool {
	return m.ID == DATESEP_MATCHER_NAME
}

func FilterDateWithoutSepMatcher(m match.Matcher) bool {
	return m.ID == DATEWITHOUTSEP_MATCHER_NAME
}

func checkDate(day, month, year int64) (bool, int64, int64, int64) {
	if (12 <= month && month <= 31) && day <= 12 {
		day, month = month, day
	}

	if day > 31 || month > 12 {
		return false, 0, 0, 0
	}

	if !((1900 <= year && year <= 2019) || (0 <= year && year <= 99)) {
		return false, 0, 0, 0
	}

	return true, day, month, year
}

func dateSepMatcher(password string) []match.Match {
	dateMatches := dateSepMatchHelper(password)

	var matches []match.Match
	for _, dateMatch := range dateMatches {
		match := match.Match{
			I:              dateMatch.I,
			J:              dateMatch.J,
			Entropy:        entropy.DateEntropy(dateMatch),
			DictionaryName: "date_match",
			Token:          dateMatch.Token,
		}

		matches = append(matches, match)
	}

	return matches
}
func dateSepMatchHelper(password string) []match.DateMatch {

	var matches []match.DateMatch

	matcher := regexp.MustCompile(DATE_RX_YEAR_SUFFIX)
	for _, v := range matcher.FindAllString(password, len(password)) {
		splitV := matcher.FindAllStringSubmatch(v, len(v))
		i := strings.Index(password, v)
		j := i + len(v)
		day, _ := strconv.ParseInt(splitV[0][4], 10, 16)
		month, _ := strconv.ParseInt(splitV[0][2], 10, 16)
		year, _ := strconv.ParseInt(splitV[0][6], 10, 16)
		match := match.DateMatch{Day: day, Month: month, Year: year, Separator: splitV[0][5], I: i, J: j, Token: password[i:j]}
		matches = append(matches, match)
	}

	matcher = regexp.MustCompile(DATE_RX_YEAR_PREFIX)
	for _, v := range matcher.FindAllString(password, len(password)) {
		splitV := matcher.FindAllStringSubmatch(v, len(v))
		i := strings.Index(password, v)
		j := i + len(v)
		day, _ := strconv.ParseInt(splitV[0][4], 10, 16)
		month, _ := strconv.ParseInt(splitV[0][6], 10, 16)
		year, _ := strconv.ParseInt(splitV[0][2], 10, 16)
		match := match.DateMatch{Day: day, Month: month, Year: year, Separator: splitV[0][5], I: i, J: j, Token: password[i:j]}
		matches = append(matches, match)
	}

	var out []match.DateMatch
	for _, match := range matches {
		if valid, day, month, year := checkDate(match.Day, match.Month, match.Year); valid {
			match.Pattern = "date"
			match.Day = day
			match.Month = month
			match.Year = year
			out = append(out, match)
		}
	}
	return out

}

type DateMatchCandidate struct {
	DayMonth string
	Year     string
	I, J     int
}

type DateMatchCandidateTwo struct {
	Day   string
	Month string
	Year  string
	I, J  int
}

func dateWithoutSepMatch(password string) []match.Match {
	dateMatches := dateWithoutSepMatchHelper(password)

	var matches []match.Match
	for _, dateMatch := range dateMatches {
		match := match.Match{
			I:              dateMatch.I,
			J:              dateMatch.J,
			Entropy:        entropy.DateEntropy(dateMatch),
			DictionaryName: "date_match",
			Token:          dateMatch.Token,
		}

		matches = append(matches, match)
	}

	return matches
}

//TODO Has issues with 6 digit dates
func dateWithoutSepMatchHelper(password string) (matches []match.DateMatch) {
	matcher := regexp.MustCompile(DATE_WITHOUT_SEP_MATCH)
	for _, v := range matcher.FindAllString(password, len(password)) {
		i := strings.Index(password, v)
		j := i + len(v)
		length := len(v)
		lastIndex := length - 1
		var candidatesRoundOne []DateMatchCandidate

		if length <= 6 {
			//2-digit year prefix
			candidatesRoundOne = append(candidatesRoundOne, buildDateMatchCandidate(v[2:], v[0:2], i, j))

			//2-digityear suffix
			candidatesRoundOne = append(candidatesRoundOne, buildDateMatchCandidate(v[0:lastIndex-2], v[lastIndex-2:], i, j))
		}
		if length >= 6 {
			//4-digit year prefix
			candidatesRoundOne = append(candidatesRoundOne, buildDateMatchCandidate(v[4:], v[0:4], i, j))

			//4-digit year sufix
			candidatesRoundOne = append(candidatesRoundOne, buildDateMatchCandidate(v[0:lastIndex-3], v[lastIndex-3:], i, j))
		}

		var candidatesRoundTwo []DateMatchCandidateTwo
		for _, c := range candidatesRoundOne {
			if len(c.DayMonth) == 2 {
				candidatesRoundTwo = append(candidatesRoundTwo, buildDateMatchCandidateTwo(c.DayMonth[0:0], c.DayMonth[1:1], c.Year, c.I, c.J))
			} else if len(c.DayMonth) == 3 {
				candidatesRoundTwo = append(candidatesRoundTwo, buildDateMatchCandidateTwo(c.DayMonth[0:2], c.DayMonth[2:2], c.Year, c.I, c.J))
				candidatesRoundTwo = append(candidatesRoundTwo, buildDateMatchCandidateTwo(c.DayMonth[0:0], c.DayMonth[1:3], c.Year, c.I, c.J))
			} else if len(c.DayMonth) == 4 {
				candidatesRoundTwo = append(candidatesRoundTwo, buildDateMatchCandidateTwo(c.DayMonth[0:2], c.DayMonth[2:4], c.Year, c.I, c.J))
			}
		}

		for _, candidate := range candidatesRoundTwo {
			intDay, err := strconv.ParseInt(candidate.Day, 10, 16)
			if err != nil {
				continue
			}

			intMonth, err := strconv.ParseInt(candidate.Month, 10, 16)

			if err != nil {
				continue
			}

			intYear, err := strconv.ParseInt(candidate.Year, 10, 16)
			if err != nil {
				continue
			}

			if ok, _, _, _ := checkDate(intDay, intMonth, intYear); ok {
				matches = append(matches, match.DateMatch{Token: password, Pattern: "date", Day: intDay, Month: intMonth, Year: intYear, I: i, J: j})
			}

		}
	}

	return matches
}

func buildDateMatchCandidate(dayMonth, year string, i, j int) DateMatchCandidate {
	return DateMatchCandidate{DayMonth: dayMonth, Year: year, I: i, J: j}
}

func buildDateMatchCandidateTwo(day, month string, year string, i, j int) DateMatchCandidateTwo {

	return DateMatchCandidateTwo{Day: day, Month: month, Year: year, I: i, J: j}
}
