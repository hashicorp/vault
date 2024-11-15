// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package ignorefiles

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"text/scanner"
)

func readRules(input io.Reader) ([]rule, error) {
	rules := defaultExclusions
	scanner := bufio.NewScanner(input)
	scanner.Split(bufio.ScanLines)
	currentRuleIndex := len(defaultExclusions) - 1

	for scanner.Scan() {
		pattern := scanner.Text()
		// Ignore blank lines
		if len(pattern) == 0 {
			continue
		}
		// Trim spaces
		pattern = strings.TrimSpace(pattern)
		// Ignore comments
		if pattern[0] == '#' {
			continue
		}
		// New rule structure
		rule := rule{}
		// Exclusions
		if pattern[0] == '!' {
			rule.negated = true
			pattern = pattern[1:]
			// Mark all previous rules as having negations after it
			for i := currentRuleIndex; i >= 0; i-- {
				if rules[i].negationsAfter {
					break
				}
				rules[i].negationsAfter = true
			}
		}
		// If it is a directory, add ** so we catch descendants
		if pattern[len(pattern)-1] == os.PathSeparator {
			pattern = pattern + "**"
		}
		// If it starts with /, it is absolute
		if pattern[0] == os.PathSeparator {
			pattern = pattern[1:]
		} else {
			// Otherwise prepend **/
			pattern = "**" + string(os.PathSeparator) + pattern
		}
		rule.val = pattern
		rules = append(rules, rule)
		currentRuleIndex += 1
	}

	if err := scanner.Err(); err != nil {
		return nil, fmt.Errorf("syntax error in .terraformignore: %w", err)
	}
	return rules, nil
}

type rule struct {
	val            string         // the value of the rule itself
	negated        bool           // prefixed by !, a negated rule
	negationsAfter bool           // negatied rules appear after this rule
	regex          *regexp.Regexp // regular expression to match for the rule
}

func (r *rule) match(path string) (bool, error) {
	if r.regex == nil {
		if err := r.compile(); err != nil {
			return false, filepath.ErrBadPattern
		}
	}

	b := r.regex.MatchString(path)
	return b, nil
}

func (r *rule) compile() error {
	regStr := "^"
	pattern := r.val
	// Go through the pattern and convert it to a regexp.
	// Use a scanner to support utf-8 chars.
	var scan scanner.Scanner
	scan.Init(strings.NewReader(pattern))

	sl := string(os.PathSeparator)
	escSL := sl
	if sl == `\` {
		escSL += `\`
	}

	for scan.Peek() != scanner.EOF {
		ch := scan.Next()
		if ch == '*' {
			if scan.Peek() == '*' {
				// is some flavor of "**"
				scan.Next()

				// Treat **/ as ** so eat the "/"
				if string(scan.Peek()) == sl {
					scan.Next()
				}

				if scan.Peek() == scanner.EOF {
					// is "**EOF" - to align with .gitignore just accept all
					regStr += ".*"
				} else {
					// is "**"
					// Note that this allows for any # of /'s (even 0) because
					// the .* will eat everything, even /'s
					regStr += "(.*" + escSL + ")?"
				}
			} else {
				// is "*" so map it to anything but "/"
				regStr += "[^" + escSL + "]*"
			}
		} else if ch == '?' {
			// "?" is any char except "/"
			regStr += "[^" + escSL + "]"
		} else if ch == '.' || ch == '$' {
			// Escape some regexp special chars that have no meaning
			// in golang's filepath.Match
			regStr += `\` + string(ch)
		} else if ch == '\\' {
			// escape next char. Note that a trailing \ in the pattern
			// will be left alone (but need to escape it)
			if sl == `\` {
				// On windows map "\" to "\\", meaning an escaped backslash,
				// and then just continue because filepath.Match on
				// Windows doesn't allow escaping at all
				regStr += escSL
				continue
			}
			if scan.Peek() != scanner.EOF {
				regStr += `\` + string(scan.Next())
			} else {
				regStr += `\`
			}
		} else {
			regStr += string(ch)
		}
	}

	regStr += "$"
	re, err := regexp.Compile(regStr)
	if err != nil {
		return err
	}

	r.regex = re
	return nil
}

/*
	Default rules as they would appear in .terraformignore:
	.git/
	.terraform/
	!.terraform/modules/
*/

var defaultExclusions = []rule{
	{
		val:            strings.Join([]string{"**", ".terraform", "**"}, string(os.PathSeparator)),
		negated:        false,
		negationsAfter: true,
	},
	// Place negation rules as high as possible in the list
	{
		val:            strings.Join([]string{"**", ".terraform", "modules", "**"}, string(os.PathSeparator)),
		negated:        true,
		negationsAfter: false,
	},
	{
		val:            strings.Join([]string{"**", ".git", "**"}, string(os.PathSeparator)),
		negated:        false,
		negationsAfter: false,
	},
}

func init() {
	// We'll precompile all of the default rules at initialization, so we
	// don't need to recompile them every time we encounter a package that
	// doesn't have any rules (the common case).
	for _, r := range defaultExclusions {
		err := r.compile()
		if err != nil {
			panic(fmt.Sprintf("invalid default rule %q: %s", r.val, err))
		}
	}
}
