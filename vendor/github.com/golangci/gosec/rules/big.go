// (c) Copyright 2016 Hewlett Packard Enterprise Development LP
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package rules

import (
	"go/ast"

	"github.com/golangci/gosec"
)

type usingBigExp struct {
	gosec.MetaData
	pkg   string
	calls []string
}

func (r *usingBigExp) ID() string {
	return r.MetaData.ID
}

func (r *usingBigExp) Match(n ast.Node, c *gosec.Context) (gi *gosec.Issue, err error) {
	if _, matched := gosec.MatchCallByType(n, c, r.pkg, r.calls...); matched {
		return gosec.NewIssue(c, n, r.ID(), r.What, r.Severity, r.Confidence), nil
	}
	return nil, nil
}

// NewUsingBigExp detects issues with modulus == 0 for Bignum
func NewUsingBigExp(id string, conf gosec.Config) (gosec.Rule, []ast.Node) {
	return &usingBigExp{
		pkg:   "*math/big.Int",
		calls: []string{"Exp"},
		MetaData: gosec.MetaData{
			ID:         id,
			What:       "Use of math/big.Int.Exp function should be audited for modulus == 0",
			Severity:   gosec.Low,
			Confidence: gosec.High,
		},
	}, []ast.Node{(*ast.CallExpr)(nil)}
}
