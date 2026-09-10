// Copyright IBM Corp. 2016, 2025
// SPDX-License-Identifier: BUSL-1.1

package random

import (
	"context"
	"crypto/rand"
	"encoding/json"
	"fmt"
	"io"
	"math"
	MRAND "math/rand"
	"reflect"
	"sort"
	"strings"
	"testing"
	"time"
)

func TestStringGenerator_Generate_successful(t *testing.T) {
	type testCase struct {
		timeout   time.Duration
		generator *StringGenerator
	}

	tests := map[string]testCase{
		"common rules": {
			timeout: 1 * time.Second,
			generator: &StringGenerator{
				Length: 20,
				Rules: []Rule{
					CharsetRule{
						Charset:  LowercaseRuneset,
						MinChars: 1,
					},
					CharsetRule{
						Charset:  UppercaseRuneset,
						MinChars: 1,
					},
					CharsetRule{
						Charset:  NumericRuneset,
						MinChars: 1,
					},
					CharsetRule{
						Charset:  ShortSymbolRuneset,
						MinChars: 1,
					},
				},
				charset: AlphaNumericShortSymbolRuneset,
			},
		},
		"charset not explicitly specified": {
			timeout: 1 * time.Second,
			generator: &StringGenerator{
				Length: 20,
				Rules: []Rule{
					CharsetRule{
						Charset:  LowercaseRuneset,
						MinChars: 1,
					},
					CharsetRule{
						Charset:  UppercaseRuneset,
						MinChars: 1,
					},
					CharsetRule{
						Charset:  NumericRuneset,
						MinChars: 1,
					},
				},
			},
		},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			// One context to rule them all, one context to find them, one context to bring them all and in the darkness bind them.
			ctx, cancel := context.WithTimeout(context.Background(), test.timeout)
			defer cancel()

			runeset := map[rune]bool{}
			runesFound := []rune{}

			for i := 0; i < 100; i++ {
				actual, err := test.generator.Generate(ctx, nil)
				if err != nil {
					t.Fatalf("no error expected, but got: %s", err)
				}
				for _, r := range actual {
					if runeset[r] {
						continue
					}
					runeset[r] = true
					runesFound = append(runesFound, r)
				}
			}

			sort.Sort(runes(runesFound))

			expectedCharset := getChars(test.generator.Rules)

			if !reflect.DeepEqual(runesFound, expectedCharset) {
				t.Fatalf("Didn't find all characters from the charset\nActual  : [%s]\nExpected: [%s]", string(runesFound), string(expectedCharset))
			}
		})
	}
}

func TestStringGenerator_Generate_consecutiveChars(t *testing.T) {
	type testCase struct {
		generator *StringGenerator
		// mustContain is a substring that every generated string must include. Empty means no requirement.
		mustContain string
	}

	// Every case runs under a short budget. With rejection sampling the two-character charsets would time out, since
	// almost every random candidate contains an adjacent repeat, so this also proves the string is built constructively.
	tests := map[string]testCase{
		"two character charset": {
			generator: &StringGenerator{
				Length:                  20,
				ConsecutiveCharsAllowed: boolPtr(false),
				Rules: []Rule{
					CharsetRule{
						Charset: []rune("ab"),
					},
				},
			},
		},
		"two character charset long": {
			generator: &StringGenerator{
				Length:                  100,
				ConsecutiveCharsAllowed: boolPtr(false),
				Rules: []Rule{
					CharsetRule{
						Charset: []rune("ab"),
					},
				},
			},
		},
		"same letter in both cases": {
			// The comparison is exact, so "aA" is legal. With only these two characters the string must alternate.
			generator: &StringGenerator{
				Length:                  20,
				ConsecutiveCharsAllowed: boolPtr(false),
				Rules: []Rule{
					CharsetRule{
						Charset: []rune("aA"),
					},
				},
			},
			mustContain: "aA",
		},
		"mixed charset": {
			generator: &StringGenerator{
				Length:                  20,
				ConsecutiveCharsAllowed: boolPtr(false),
				Rules: []Rule{
					CharsetRule{
						Charset: []rune("aAbBcC"),
					},
				},
			},
		},
		"with charset rules": {
			generator: &StringGenerator{
				Length:                  20,
				ConsecutiveCharsAllowed: boolPtr(false),
				Rules: []Rule{
					CharsetRule{
						Charset:  []rune("abc"),
						MinChars: 3,
					},
					CharsetRule{
						Charset:  []rune("ABC"),
						MinChars: 3,
					},
					CharsetRule{
						Charset:  []rune("123"),
						MinChars: 3,
					},
				},
			},
		},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			ctx, cancel := context.WithTimeout(context.Background(), 100*time.Millisecond)
			defer cancel()

			for i := 0; i < 100; i++ {
				actual, err := test.generator.Generate(ctx, nil)
				if err != nil {
					t.Fatalf("no error expected, but got: %s", err)
				}
				if len([]rune(actual)) != test.generator.Length {
					t.Fatalf("expected length %d, got %d: %q", test.generator.Length, len([]rune(actual)), actual)
				}
				if hasConsecutiveChars([]rune(actual)) {
					t.Fatalf("generated string contains consecutive characters: %q", actual)
				}
				if test.mustContain != "" && !strings.Contains(actual, test.mustContain) {
					t.Fatalf("generated string %q does not contain %q", actual, test.mustContain)
				}
				for _, rule := range test.generator.Rules {
					if !rule.Pass([]rune(actual)) {
						t.Fatalf("generated string failed rule %s: %q", rule.Type(), actual)
					}
				}
			}
		})
	}
}

func TestRandomRunesNoConsecutive_successful(t *testing.T) {
	type testCase struct {
		charset []rune // Assumes no duplicate runes
		length  int
	}

	tests := map[string]testCase{
		"two characters": {
			charset: []rune("ab"),
			length:  20,
		},
		"small charset": {
			charset: []rune("abcde"),
			length:  20,
		},
		"common charset": {
			charset: AlphaNumericShortSymbolRuneset,
			length:  20,
		},
		"length 1 with single character": {
			charset: []rune("a"),
			length:  1,
		},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			runeset := map[rune]bool{}
			runesFound := []rune{}

			for i := 0; i < 10000; i++ {
				actual, err := randomRunesNoConsecutive(rand.Reader, test.charset, test.length)
				if err != nil {
					t.Fatalf("no error expected, but got: %s", err)
				}
				if len(actual) != test.length {
					t.Fatalf("expected length %d, got %d: %q", test.length, len(actual), string(actual))
				}
				if hasConsecutiveChars(actual) {
					t.Fatalf("string contains consecutive characters: %q", string(actual))
				}
				for _, r := range actual {
					if runeset[r] {
						continue
					}
					runeset[r] = true
					runesFound = append(runesFound, r)
				}
			}

			sort.Sort(runes(runesFound))

			// Sort the input too just to ensure that they can be compared
			sort.Sort(runes(test.charset))

			// Every rune must be reachable, including the last one which is only selected via the index shift
			if !reflect.DeepEqual(runesFound, test.charset) {
				t.Fatalf("Didn't find all characters from the charset\nActual  : [%s]\nExpected: [%s]", string(runesFound), string(test.charset))
			}
		})
	}
}

func TestRandomRunesNoConsecutive_deterministic(t *testing.T) {
	// Pins the index shift: with charset "ab" every string must strictly alternate, and with a larger charset the
	// seeded output proves the shifted index is applied relative to the previous rune.
	type testCase struct {
		rngSeed  int64
		charset  string
		length   int
		expected string
	}

	tests := map[string]testCase{
		"two characters": {
			rngSeed:  1585593298447807000,
			charset:  "ab",
			length:   10,
			expected: "babababababa"[:10],
		},
		"small charset": {
			rngSeed:  1585593298447807000,
			charset:  "abcde",
			length:   20,
			expected: "dadaebadbdbdaebdcecb",
		},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			rng := MRAND.New(MRAND.NewSource(test.rngSeed))
			actual, err := randomRunesNoConsecutive(rng, []rune(test.charset), test.length)
			if err != nil {
				t.Fatalf("Expected no error, but found: %s", err)
			}
			if string(actual) != test.expected {
				t.Fatalf("Actual: %s  Expected: %s", string(actual), test.expected)
			}
		})
	}
}

func TestRandomRunesNoConsecutive_errors(t *testing.T) {
	type testCase struct {
		charset []rune
		length  int
		rng     io.Reader
	}

	tests := map[string]testCase{
		"nil charset": {
			charset: nil,
			length:  20,
			rng:     rand.Reader,
		},
		"single character with length > 1": {
			charset: []rune("a"),
			length:  2,
			rng:     rand.Reader,
		},
		"zero length": {
			charset: []rune("ab"),
			length:  0,
			rng:     rand.Reader,
		},
		"bad RNG reader": {
			charset: []rune("ab"),
			length:  20,
			rng:     badReader{},
		},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			actual, err := randomRunesNoConsecutive(test.rng, test.charset, test.length)
			if err == nil {
				t.Fatalf("Expected error but none found")
			}
			if actual != nil {
				t.Fatalf("Expected nil result but found: %q", string(actual))
			}
		})
	}
}

// hasConsecutiveChars returns true if any two adjacent runes are identical. Test helper only: production code builds
// strings so that this can never be true rather than checking after the fact.
func hasConsecutiveChars(value []rune) bool {
	for i := 1; i < len(value); i++ {
		if value[i] == value[i-1] {
			return true
		}
	}
	return false
}

func TestStringGenerator_Generate_consecutiveCharsAllowedByDefault(t *testing.T) {
	// With a single-character charset every generated string necessarily repeats that character, which proves the
	// restriction is off unless explicitly disabled.
	gen := &StringGenerator{
		Length: 5,
		Rules: []Rule{
			CharsetRule{
				Charset: []rune("a"),
			},
		},
	}

	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
	defer cancel()

	actual, err := gen.Generate(ctx, nil)
	if err != nil {
		t.Fatalf("no error expected, but got: %s", err)
	}
	if actual != "aaaaa" {
		t.Fatalf("expected %q, got %q", "aaaaa", actual)
	}
}

func TestHasConsecutiveChars(t *testing.T) {
	tests := map[string]struct {
		value    string
		expected bool
	}{
		"empty":                     {value: "", expected: false},
		"single char":               {value: "a", expected: false},
		"no repeats":                {value: "abcdef", expected: false},
		"repeat separated":          {value: "aba", expected: false},
		"lowercase repeat":          {value: "abccd", expected: true},
		"uppercase repeat":          {value: "ABBC", expected: true},
		"different case is allowed": {value: "abaAd", expected: false},
		"repeat at start":           {value: "aab", expected: true},
		"repeat at end":             {value: "baa", expected: true},
		"digit repeat":              {value: "a11b", expected: true},
		"symbol repeat":             {value: "a--b", expected: true},
		"accented char is distinct": {value: "aá", expected: false},
		"non-ascii repeat":          {value: "éé", expected: true},
		"non-ascii different case":  {value: "éÉ", expected: false},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			actual := hasConsecutiveChars([]rune(test.value))
			if actual != test.expected {
				t.Fatalf("hasConsecutiveChars(%q) = %t, expected %t", test.value, actual, test.expected)
			}
		})
	}
}

func TestStringGenerator_Generate_errors(t *testing.T) {
	type testCase struct {
		timeout   time.Duration
		generator *StringGenerator
		rng       io.Reader
	}

	tests := map[string]testCase{
		"already timed out": {
			timeout: 0,
			generator: &StringGenerator{
				Length: 20,
				Rules: []Rule{
					testCharsetRule{
						fail: false,
					},
				},
				charset: AlphaNumericShortSymbolRuneset,
			},
			rng: rand.Reader,
		},
		"impossible rules": {
			timeout: 10 * time.Millisecond, // Keep this short so the test doesn't take too long
			generator: &StringGenerator{
				Length: 20,
				Rules: []Rule{
					testCharsetRule{
						fail: true,
					},
				},
				charset: AlphaNumericShortSymbolRuneset,
			},
			rng: rand.Reader,
		},
		"bad RNG reader": {
			timeout: 10 * time.Millisecond, // Keep this short so the test doesn't take too long
			generator: &StringGenerator{
				Length:  20,
				Rules:   []Rule{},
				charset: AlphaNumericShortSymbolRuneset,
			},
			rng: badReader{},
		},
		"0 length": {
			timeout: 10 * time.Millisecond,
			generator: &StringGenerator{
				Length: 0,
				Rules: []Rule{
					CharsetRule{
						Charset:  []rune("abcde"),
						MinChars: 0,
					},
				},
				charset: []rune("abcde"),
			},
			rng: rand.Reader,
		},
		"-1 length": {
			timeout: 10 * time.Millisecond,
			generator: &StringGenerator{
				Length: -1,
				Rules: []Rule{
					CharsetRule{
						Charset:  []rune("abcde"),
						MinChars: 0,
					},
				},
				charset: []rune("abcde"),
			},
			rng: rand.Reader,
		},
		"no charset": {
			timeout: 10 * time.Millisecond,
			generator: &StringGenerator{
				Length: 20,
				Rules:  []Rule{},
			},
			rng: rand.Reader,
		},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			// One context to rule them all, one context to find them, one context to bring them all and in the darkness bind them.
			ctx, cancel := context.WithTimeout(context.Background(), test.timeout)
			defer cancel()

			actual, err := test.generator.Generate(ctx, test.rng)
			if err == nil {
				t.Fatalf("Expected error but none found")
			}
			if actual != "" {
				t.Fatalf("Random string returned: %s", actual)
			}
		})
	}
}

func TestRandomRunes_deterministic(t *testing.T) {
	// These tests are to ensure that the charset selection doesn't do anything weird like selecting the same character
	// over and over again. The number of test cases here should be kept to a minimum since they are sensitive to changes
	type testCase struct {
		rngSeed  int64
		charset  string
		length   int
		expected string
	}

	tests := map[string]testCase{
		"small charset": {
			rngSeed:  1585593298447807000,
			charset:  "abcde",
			length:   20,
			expected: "ddddddcdebbeebdbdbcd",
		},
		"common charset": {
			rngSeed:  1585593298447807001,
			charset:  AlphaNumericShortSymbolCharset,
			length:   20,
			expected: "ON6lVjnBs84zJbUBVEzb",
		},
		"max size charset": {
			rngSeed: 1585593298447807002,
			charset: " !\"#$%&'()*+,-./0123456789:;<=>?@ABCDEFGHIJKLMNOPQRSTUVWXYZ[\\]^_" +
				"`abcdefghijklmnopqrstuvwxyz{|}~ĀāĂăĄąĆćĈĉĊċČčĎďĐđĒēĔĕĖėĘęĚěĜĝĞğĠ" +
				"ġĢģĤĥĦħĨĩĪīĬĭĮįİıĲĳĴĵĶķĸĹĺĻļĽľĿŀŁłŃńŅņŇňŉŊŋŌōŎŏŐőŒœŔŕŖŗŘřŚśŜŝŞşŠ" +
				"šŢţŤťŦŧŨũŪūŬŭŮůŰűŲųŴŵŶŷŸŹźŻżŽžſ℀℁ℂ℃℄℅℆ℇ℈℉ℊℋℌℍℎℏℐℑℒℓ℔ℕ№℗℘ℙℚℛℜℝ℞℟℠",
			length:   20,
			expected: "tųŎ℄ņ℃Œ.@řHš-ℍ}ħGĲLℏ",
		},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			rng := MRAND.New(MRAND.NewSource(test.rngSeed))
			runes, err := randomRunes(rng, []rune(test.charset), test.length)
			if err != nil {
				t.Fatalf("Expected no error, but found: %s", err)
			}

			str := string(runes)

			if str != test.expected {
				t.Fatalf("Actual: %s  Expected: %s", str, test.expected)
			}
		})
	}
}

func TestRandomRunes_successful(t *testing.T) {
	type testCase struct {
		charset []rune // Assumes no duplicate runes
		length  int
	}

	tests := map[string]testCase{
		"small charset": {
			charset: []rune("abcde"),
			length:  20,
		},
		"common charset": {
			charset: AlphaNumericShortSymbolRuneset,
			length:  20,
		},
		"max size charset": {
			charset: []rune(
				" !\"#$%&'()*+,-./0123456789:;<=>?@ABCDEFGHIJKLMNOPQRSTUVWXYZ[\\]^_" +
					"`abcdefghijklmnopqrstuvwxyz{|}~ĀāĂăĄąĆćĈĉĊċČčĎďĐđĒēĔĕĖėĘęĚěĜĝĞğĠ" +
					"ġĢģĤĥĦħĨĩĪīĬĭĮįİıĲĳĴĵĶķĸĹĺĻļĽľĿŀŁłŃńŅņŇňŉŊŋŌōŎŏŐőŒœŔŕŖŗŘřŚśŜŝŞşŠ" +
					"šŢţŤťŦŧŨũŪūŬŭŮůŰűŲųŴŵŶŷŸŹźŻżŽžſ℀℁ℂ℃℄℅℆ℇ℈℉ℊℋℌℍℎℏℐℑℒℓ℔ℕ№℗℘ℙℚℛℜℝ℞℟℠",
			),
			length: 20,
		},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			runeset := map[rune]bool{}
			runesFound := []rune{}

			for i := 0; i < 10000; i++ {
				actual, err := randomRunes(rand.Reader, test.charset, test.length)
				if err != nil {
					t.Fatalf("no error expected, but got: %s", err)
				}
				for _, r := range actual {
					if runeset[r] {
						continue
					}
					runeset[r] = true
					runesFound = append(runesFound, r)
				}
			}

			sort.Sort(runes(runesFound))

			// Sort the input too just to ensure that they can be compared
			sort.Sort(runes(test.charset))

			if !reflect.DeepEqual(runesFound, test.charset) {
				t.Fatalf("Didn't find all characters from the charset\nActual  : [%s]\nExpected: [%s]", string(runesFound), string(test.charset))
			}
		})
	}
}

func TestRandomRunes_errors(t *testing.T) {
	type testCase struct {
		charset []rune
		length  int
		rng     io.Reader
	}

	tests := map[string]testCase{
		"nil charset": {
			charset: nil,
			length:  20,
			rng:     rand.Reader,
		},
		"empty charset": {
			charset: []rune{},
			length:  20,
			rng:     rand.Reader,
		},
		"charset is too long": {
			charset: []rune(" !\"#$%&'()*+,-./0123456789:;<=>?@ABCDEFGHIJKLMNOPQRSTUVWXYZ[\\]^_" +
				"`abcdefghijklmnopqrstuvwxyz{|}~ĀāĂăĄąĆćĈĉĊċČčĎďĐđĒēĔĕĖėĘęĚěĜĝĞğĠ" +
				"ġĢģĤĥĦħĨĩĪīĬĭĮįİıĲĳĴĵĶķĸĹĺĻļĽľĿŀŁłŃńŅņŇňŉŊŋŌōŎŏŐőŒœŔŕŖŗŘřŚśŜŝŞşŠ" +
				"šŢţŤťŦŧŨũŪūŬŭŮůŰűŲųŴŵŶŷŸŹźŻżŽžſ℀℁ℂ℃℄℅℆ℇ℈℉ℊℋℌℍℎℏℐℑℒℓ℔ℕ№℗℘ℙℚℛℜℝ℞℟℠" +
				"Σ",
			),
			length: 20,
			rng:    rand.Reader,
		},
		"length is zero": {
			charset: []rune("abcde"),
			length:  0,
			rng:     rand.Reader,
		},
		"length is negative": {
			charset: []rune("abcde"),
			length:  -3,
			rng:     rand.Reader,
		},
		"reader failed": {
			charset: []rune("abcde"),
			length:  20,
			rng:     badReader{},
		},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			actual, err := randomRunes(test.rng, test.charset, test.length)
			if err == nil {
				t.Fatalf("Expected error but none found")
			}
			if actual != nil {
				t.Fatalf("Expected no value, but found [%s]", string(actual))
			}
		})
	}
}

func BenchmarkStringGenerator_Generate(b *testing.B) {
	lengths := []int{
		8, 12, 16, 20, 24, 28,
	}

	type testCase struct {
		generator *StringGenerator
	}

	benches := map[string]testCase{
		"no restrictions": {
			generator: &StringGenerator{
				Rules: []Rule{
					CharsetRule{
						Charset: AlphaNumericFullSymbolRuneset,
					},
				},
			},
		},
		"default generator": {
			generator: DefaultStringGenerator,
		},
		"large symbol set": {
			generator: &StringGenerator{
				Rules: []Rule{
					CharsetRule{
						Charset:  LowercaseRuneset,
						MinChars: 1,
					},
					CharsetRule{
						Charset:  UppercaseRuneset,
						MinChars: 1,
					},
					CharsetRule{
						Charset:  NumericRuneset,
						MinChars: 1,
					},
					CharsetRule{
						Charset:  FullSymbolRuneset,
						MinChars: 1,
					},
				},
			},
		},
		"max symbol set": {
			generator: &StringGenerator{
				Rules: []Rule{
					CharsetRule{
						Charset: []rune(" !\"#$%&'()*+,-./0123456789:;<=>?@ABCDEFGHIJKLMNOPQRSTUVWXYZ[\\]^_" +
							"`abcdefghijklmnopqrstuvwxyz{|}~ĀāĂăĄąĆćĈĉĊċČčĎďĐđĒēĔĕĖėĘęĚěĜĝĞğĠ" +
							"ġĢģĤĥĦħĨĩĪīĬĭĮįİıĲĳĴĵĶķĸĹĺĻļĽľĿŀŁłŃńŅņŇňŉŊŋŌōŎŏŐőŒœŔŕŖŗŘřŚśŜŝŞşŠ" +
							"šŢţŤťŦŧŨũŪūŬŭŮůŰűŲųŴŵŶŷŸŹźŻżŽžſ℀℁ℂ℃℄℅℆ℇ℈℉ℊℋℌℍℎℏℐℑℒℓ℔ℕ№℗℘ℙℚℛℜℝ℞℟℠"),
					},
					CharsetRule{
						Charset:  LowercaseRuneset,
						MinChars: 1,
					},
					CharsetRule{
						Charset:  UppercaseRuneset,
						MinChars: 1,
					},
					CharsetRule{
						Charset:  []rune("ĩĪīĬĭĮįİıĲĳĴĵĶķĸĹĺĻļĽľĿŀŁłŃńŅņŇňŉŊŋŌōŎŏŐőŒ"),
						MinChars: 1,
					},
				},
			},
		},
		"restrictive charset rules": {
			generator: &StringGenerator{
				Rules: []Rule{
					CharsetRule{
						Charset: AlphaNumericShortSymbolRuneset,
					},
					CharsetRule{
						Charset:  []rune("A"),
						MinChars: 1,
					},
					CharsetRule{
						Charset:  []rune("1"),
						MinChars: 1,
					},
					CharsetRule{
						Charset:  []rune("a"),
						MinChars: 1,
					},
					CharsetRule{
						Charset:  []rune("-"),
						MinChars: 1,
					},
				},
			},
		},
	}

	for name, bench := range benches {
		b.Run(name, func(b *testing.B) {
			for _, length := range lengths {
				bench.generator.Length = length
				b.Run(fmt.Sprintf("length=%d", length), func(b *testing.B) {
					ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
					defer cancel()

					b.ResetTimer()
					for i := 0; i < b.N; i++ {
						str, err := bench.generator.Generate(ctx, nil)
						if err != nil {
							b.Fatalf("Failed to generate string: %s", err)
						}
						if str == "" {
							b.Fatalf("Didn't error but didn't generate a string")
						}
					}
				})
			}
		})
	}

	// Mimic what the SQLCredentialsProducer is doing
	b.Run("SQLCredentialsProducer", func(b *testing.B) {
		sg := StringGenerator{
			Length:  16, // 16 because the SQLCredentialsProducer prepends 4 characters to a 20 character password
			charset: AlphaNumericRuneset,
			Rules:   nil,
		}

		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()

		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			str, err := sg.Generate(ctx, nil)
			if err != nil {
				b.Fatalf("Failed to generate string: %s", err)
			}
			if str == "" {
				b.Fatalf("Didn't error but didn't generate a string")
			}
		}
	})
}

// Ensure the StringGenerator can be properly JSON-ified
func TestStringGenerator_JSON(t *testing.T) {
	expected := StringGenerator{
		Length:  20,
		charset: deduplicateRunes([]rune("teststring" + ShortSymbolCharset)),
		Rules: []Rule{
			testCharsetRule{
				String:  "teststring",
				Integer: 123,
			},
			CharsetRule{
				Charset:  ShortSymbolRuneset,
				MinChars: 1,
			},
		},
	}

	b, err := json.Marshal(expected)
	if err != nil {
		t.Fatalf("Failed to marshal to JSON: %s", err)
	}

	parser := PolicyParser{
		RuleRegistry: Registry{
			Rules: map[string]ruleConstructor{
				"testrule": newTestRule,
				"charset":  ParseCharset,
			},
		},
	}
	actual, err := parser.ParsePolicy(string(b))
	if err != nil {
		t.Fatalf("Failed to parse JSON: %s", err)
	}

	if !reflect.DeepEqual(actual, expected) {
		t.Fatalf("Actual: %#v\nExpected: %#v", actual, expected)
	}
}

type badReader struct{}

func (badReader) Read([]byte) (int, error) {
	return 0, fmt.Errorf("test error")
}

func TestValidate(t *testing.T) {
	type testCase struct {
		generator *StringGenerator
		expectErr bool
	}

	tests := map[string]testCase{
		"default generator": {
			generator: DefaultStringGenerator,
			expectErr: false,
		},
		"length is 0": {
			generator: &StringGenerator{
				Length: 0,
			},
			expectErr: true,
		},
		"length is negative": {
			generator: &StringGenerator{
				Length: -2,
			},
			expectErr: true,
		},
		"nil charset, no rules": {
			generator: &StringGenerator{
				Length:  5,
				charset: nil,
			},
			expectErr: true,
		},
		"zero length charset, no rules": {
			generator: &StringGenerator{
				Length:  5,
				charset: []rune{},
			},
			expectErr: true,
		},
		"rules require password longer than length": {
			generator: &StringGenerator{
				Length:  5,
				charset: []rune("abcde"),
				Rules: []Rule{
					CharsetRule{
						Charset:  []rune("abcde"),
						MinChars: 6,
					},
				},
			},
			expectErr: true,
		},
		"charset has non-printable characters": {
			generator: &StringGenerator{
				Length: 0,
				charset: []rune{
					'a',
					'b',
					0, // Null character
					'd',
					'e',
				},
			},
			expectErr: true,
		},
		"consecutive chars disallowed with single character": {
			generator: &StringGenerator{
				Length:                  5,
				ConsecutiveCharsAllowed: boolPtr(false),
				charset:                 []rune("a"),
			},
			expectErr: true,
		},
		"consecutive chars disallowed with same letter in both cases": {
			generator: &StringGenerator{
				Length:                  5,
				ConsecutiveCharsAllowed: boolPtr(false),
				charset:                 []rune("aA"),
			},
			expectErr: false,
		},
		"consecutive chars disallowed with single character and length 1": {
			generator: &StringGenerator{
				Length:                  1,
				ConsecutiveCharsAllowed: boolPtr(false),
				charset:                 []rune("a"),
			},
			expectErr: false,
		},
		"consecutive chars disallowed with multiple characters": {
			generator: &StringGenerator{
				Length:                  5,
				ConsecutiveCharsAllowed: boolPtr(false),
				charset:                 []rune("ab"),
			},
			expectErr: false,
		},
		"consecutive chars allowed with single character": {
			generator: &StringGenerator{
				Length:                  5,
				ConsecutiveCharsAllowed: boolPtr(true),
				charset:                 []rune("a"),
			},
			expectErr: false,
		},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			err := test.generator.validateConfig()
			if test.expectErr && err == nil {
				t.Fatalf("err expected, got nil")
			}
			if !test.expectErr && err != nil {
				t.Fatalf("no error expected, got: %s", err)
			}
		})
	}
}

type testNonCharsetRule struct {
	String string `mapstructure:"string" json:"string"`
}

func (tr testNonCharsetRule) Pass([]rune) bool { return true }
func (tr testNonCharsetRule) Type() string     { return "testNonCharsetRule" }

func TestGetChars(t *testing.T) {
	type testCase struct {
		rules    []Rule
		expected []rune
	}

	tests := map[string]testCase{
		"nil rules": {
			rules:    nil,
			expected: []rune(nil),
		},
		"empty rules": {
			rules:    []Rule{},
			expected: []rune(nil),
		},
		"rule without chars": {
			rules: []Rule{
				testNonCharsetRule{
					String: "teststring",
				},
			},
			expected: []rune(nil),
		},
		"rule with chars": {
			rules: []Rule{
				CharsetRule{
					Charset:  []rune("abcdefghij"),
					MinChars: 1,
				},
			},
			expected: []rune("abcdefghij"),
		},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			actual := getChars(test.rules)
			if !reflect.DeepEqual(actual, test.expected) {
				t.Fatalf("Actual: %v\nExpected: %v", actual, test.expected)
			}
		})
	}
}

func TestDeduplicateRunes(t *testing.T) {
	type testCase struct {
		input    []rune
		expected []rune
	}

	tests := map[string]testCase{
		"empty string": {
			input:    []rune(""),
			expected: []rune(nil),
		},
		"no duplicates": {
			input:    []rune("abcde"),
			expected: []rune("abcde"),
		},
		"in order duplicates": {
			input:    []rune("aaaabbbbcccccccddddeeeee"),
			expected: []rune("abcde"),
		},
		"out of order duplicates": {
			input:    []rune("abcdeabcdeabcdeabcde"),
			expected: []rune("abcde"),
		},
		"unicode no duplicates": {
			input:    []rune("日本語"),
			expected: []rune("日本語"),
		},
		"unicode in order duplicates": {
			input:    []rune("日日日日本本本語語語語語"),
			expected: []rune("日本語"),
		},
		"unicode out of order duplicates": {
			input:    []rune("日本語日本語日本語日本語"),
			expected: []rune("日本語"),
		},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			actual := deduplicateRunes(test.input)
			if !reflect.DeepEqual(actual, test.expected) {
				t.Fatalf("Actual: %#v\nExpected:%#v", actual, test.expected)
			}
		})
	}
}

func TestRandomRunes_Bias(t *testing.T) {
	type testCase struct {
		charset   []rune
		maxStdDev float64
	}

	tests := map[string]testCase{
		"small charset": {
			charset:   []rune("abcde"),
			maxStdDev: 2700,
		},
		"lowercase characters": {
			charset:   LowercaseRuneset,
			maxStdDev: 1000,
		},
		"alphabetical characters": {
			charset:   AlphabeticRuneset,
			maxStdDev: 800,
		},
		"alphanumeric": {
			charset:   AlphaNumericRuneset,
			maxStdDev: 800,
		},
		"alphanumeric with symbol": {
			charset:   AlphaNumericShortSymbolRuneset,
			maxStdDev: 800,
		},
		"charset evenly divisible into 256": {
			charset:   append(AlphaNumericRuneset, '!', '@'),
			maxStdDev: 800,
		},
		"large charset": {
			charset:   FullSymbolRuneset,
			maxStdDev: 800,
		},
		"just under half size charset": {
			charset: []rune(" !\"#$%&'()*+,-./0123456789:;<=>?@ABCDEFGHIJKLMNOPQRSTUVWXYZ[\\]^_" +
				"`abcdefghijklmnopqrstuvwxyz{|}~ĀāĂăĄąĆćĈĉĊċČčĎďĐđĒēĔĕĖėĘęĚěĜĝĞğ"),
			maxStdDev: 800,
		},
		"half size charset": {
			charset: []rune(" !\"#$%&'()*+,-./0123456789:;<=>?@ABCDEFGHIJKLMNOPQRSTUVWXYZ[\\]^_" +
				"`abcdefghijklmnopqrstuvwxyz{|}~ĀāĂăĄąĆćĈĉĊċČčĎďĐđĒēĔĕĖėĘęĚěĜĝĞğĠ"),
			maxStdDev: 800,
		},
	}

	for name, test := range tests {
		t.Run(fmt.Sprintf("%s (%d chars)", name, len(test.charset)), func(t *testing.T) {
			runeCounts := map[rune]int{}

			generations := 50000
			length := 100
			for i := 0; i < generations; i++ {
				str, err := randomRunes(nil, test.charset, length)
				if err != nil {
					t.Fatal(err)
				}
				for _, r := range str {
					runeCounts[r]++
				}
			}

			chars := charCounts{}

			var sum float64
			for r, count := range runeCounts {
				chars = append(chars, charCount{r, count})
				sum += float64(count)
			}

			mean := sum / float64(len(runeCounts))
			var stdDev float64
			for _, count := range runeCounts {
				stdDev += math.Pow(float64(count)-mean, 2)
			}

			stdDev = math.Sqrt(stdDev / float64(len(runeCounts)))
			t.Logf("Mean  : %10.4f", mean)

			if stdDev > test.maxStdDev {
				t.Fatalf("Standard deviation is too large: %.2f > %.2f", stdDev, test.maxStdDev)
			}
		})
	}
}

type charCount struct {
	r     rune
	count int
}

type charCounts []charCount

func (s charCounts) Len() int           { return len(s) }
func (s charCounts) Less(i, j int) bool { return s[i].r < s[j].r }
func (s charCounts) Swap(i, j int)      { s[i], s[j] = s[j], s[i] }
