package random

import (
	"context"
	"crypto/rand"
	"fmt"
	"io"
	MRAND "math/rand"
	"reflect"
	"sort"
	"testing"
	"time"
)

func TestStringGenerator_Generate_successful(t *testing.T) {
	type testCase struct {
		timeout time.Duration
		charset []rune
		rules   []Rule
	}

	tests := map[string]testCase{
		"common rules": {
			timeout: 1 * time.Second,
			charset: []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789-"),
			rules: []Rule{
				CharsetRestriction{
					Charset:  []rune("abcdefghijklmnopqrstuvwxyz"),
					MinChars: 1,
				},
				CharsetRestriction{
					Charset:  []rune("ABCDEFGHIJKLMNOPQRSTUVWXYZ"),
					MinChars: 1,
				},
				CharsetRestriction{
					Charset:  []rune("0123456789"),
					MinChars: 1,
				},
			},
		},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			sg := StringGenerator{
				Length:  20,
				Charset: test.charset,
				Rules:   test.rules,
			}

			// One context to rule them all, one context to find them, one context to bring them all and in the darkness bind them.
			ctx, cancel := context.WithTimeout(context.Background(), test.timeout)
			defer cancel()

			runeset := map[rune]bool{}
			runesFound := []rune{}

			for i := 0; i < 10000; i++ {
				actual, err := sg.Generate(ctx)
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

func TestStringGenerator_Generate_errors(t *testing.T) {
	type testCase struct {
		timeout time.Duration
		charset string
		rules   []Rule
		rng     io.Reader
	}

	tests := map[string]testCase{
		"already timed out": {
			timeout: 0,
			charset: "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789-",
			rules: []Rule{
				testRule{
					fail: false,
				},
			},
			rng: rand.Reader,
		},
		"impossible rules": {
			timeout: 10 * time.Millisecond, // Keep this short so the test doesn't take too long
			charset: "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789-",
			rules: []Rule{
				testRule{
					fail: true,
				},
			},
			rng: rand.Reader,
		},
		"bad RNG reader": {
			timeout: 10 * time.Millisecond, // Keep this short so the test doesn't take too long
			charset: "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789-",
			rules:   []Rule{},
			rng:     badReader{},
		},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			sg := StringGenerator{
				Length:  20,
				Charset: []rune(test.charset),
				Rules:   test.rules,
				rng:     test.rng,
			}

			// One context to rule them all, one context to find them, one context to bring them all and in the darkness bind them.
			ctx, cancel := context.WithTimeout(context.Background(), test.timeout)
			defer cancel()

			actual, err := sg.Generate(ctx)
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
			charset:  "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789-",
			length:   20,
			expected: "czyhWGUYm3jf-uMFmGp-",
		},
		"max size charset": {
			rngSeed: 1585593298447807002,
			charset: " !\"#$%&'()*+,-./0123456789:;<=>?@ABCDEFGHIJKLMNOPQRSTUVWXYZ[\\]^_" +
				"`abcdefghijklmnopqrstuvwxyz{|}~ĀāĂăĄąĆćĈĉĊċČčĎďĐđĒēĔĕĖėĘęĚěĜĝĞğĠ" +
				"ġĢģĤĥĦħĨĩĪīĬĭĮįİıĲĳĴĵĶķĸĹĺĻļĽľĿŀŁłŃńŅņŇňŉŊŋŌōŎŏŐőŒœŔŕŖŗŘřŚśŜŝŞşŠ" +
				"šŢţŤťŦŧŨũŪūŬŭŮůŰűŲųŴŵŶŷŸŹźŻżŽžſ℀℁ℂ℃℄℅℆ℇ℈℉ℊℋℌℍℎℏℐℑℒℓ℔ℕ№℗℘ℙℚℛℜℝ℞℟",
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
			charset: []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789-"),
			length:  20,
		},
		"max size charset": {
			charset: []rune(
				" !\"#$%&'()*+,-./0123456789:;<=>?@ABCDEFGHIJKLMNOPQRSTUVWXYZ[\\]^_" +
					"`abcdefghijklmnopqrstuvwxyz{|}~ĀāĂăĄąĆćĈĉĊċČčĎďĐđĒēĔĕĖėĘęĚěĜĝĞğĠ" +
					"ġĢģĤĥĦħĨĩĪīĬĭĮįİıĲĳĴĵĶķĸĹĺĻļĽľĿŀŁłŃńŅņŇňŉŊŋŌōŎŏŐőŒœŔŕŖŗŘřŚśŜŝŞşŠ" +
					"šŢţŤťŦŧŨũŪūŬŭŮůŰűŲųŴŵŶŷŸŹźŻżŽžſ℀℁ℂ℃℄℅℆ℇ℈℉ℊℋℌℍℎℏℐℑℒℓ℔ℕ№℗℘ℙℚℛℜℝ℞℟",
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

type badReader struct{}

func (badReader) Read([]byte) (int, error) {
	return 0, fmt.Errorf("test error")
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
				"šŢţŤťŦŧŨũŪūŬŭŮůŰűŲųŴŵŶŷŸŹźŻżŽžſ℀℁ℂ℃℄℅℆ℇ℈℉ℊℋℌℍℎℏℐℑℒℓ℔ℕ№℗℘ℙℚℛℜℝ℞℟" +
				"Σ",
			),
			rng: rand.Reader,
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

	b.Run("default string generator", func(b *testing.B) {
		for _, length := range lengths {
			b.Run(fmt.Sprintf("length=%d", length), func(b *testing.B) {
				for i := 0; i < b.N; i++ {
					ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
					DefaultStringGenerator.Generate(ctx)
					cancel()
				}
			})
		}
	})
	b.Run("large symbol set", func(b *testing.B) {
		sg := StringGenerator{
			Length:  20,
			Charset: []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789 !\"#$%&'()*+,-./0123456789:;<=>?@ABCDEFGHIJKLMNOPQRSTUVWXYZ[\\]^_"),
			Rules: []Rule{
				CharsetRestriction{
					Charset:  []rune("abcdefghijklmnopqrstuvwxyz"),
					MinChars: 1,
				},
				CharsetRestriction{
					Charset:  []rune("ABCDEFGHIJKLMNOPQRSTUVWXYZ"),
					MinChars: 1,
				},
				CharsetRestriction{
					Charset:  []rune("0123456789"),
					MinChars: 1,
				},
				CharsetRestriction{
					Charset:  []rune(" !\"#$%&'()*+,-./0123456789:;<=>?@ABCDEFGHIJKLMNOPQRSTUVWXYZ[\\]^_"),
					MinChars: 1,
				},
			},
		}
		for _, length := range lengths {
			b.Run(fmt.Sprintf("length=%d", length), func(b *testing.B) {
				for i := 0; i < b.N; i++ {
					ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
					str, err := sg.Generate(ctx)
					cancel()
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
