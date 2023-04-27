// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package diagnose

import (
	"context"
	"errors"
	"os"
	"reflect"
	"strings"
	"testing"

	"github.com/go-test/deep"
)

const getMoreCoffee = "You'll find more coffee in the freezer door, or consider buying more for the office."

func TestDiagnoseOtelResults(t *testing.T) {
	expected := &Result{
		Name:   "make-coffee",
		Status: ErrorStatus,
		Warnings: []string{
			"coffee getting low",
		},
		Advice: getMoreCoffee,
		Children: []*Result{
			{
				Name:   "prepare-kitchen",
				Status: ErrorStatus,
				Children: []*Result{
					{
						Name:   "build-microwave",
						Status: ErrorStatus,
						Children: []*Result{
							{
								Name:    "buy-parts",
								Status:  ErrorStatus,
								Message: "no stores sell microwave parts, please buy a microwave instead.",
								Warnings: []string{
									"warning: you are about to try to build a microwave from scratch.",
								},
							},
						},
					},
				},
			},
			{
				Name:   "warm-milk",
				Status: OkStatus,
			},
			{
				Name:   "brew-coffee",
				Status: OkStatus,
			},
			{
				Name:    "pick-scone",
				Status:  ErrorStatus,
				Message: "no scones",
			},
			{
				Name:    "dispose-grounds",
				Status:  SkippedStatus,
				Message: "skipped as requested",
			},
		},
	}
	sess := New(os.Stdout)
	sess.SkipFilters = []string{"dispose-grounds"}
	ctx := Context(context.Background(), sess)

	func() {
		ctx, span := StartSpan(ctx, "make-coffee")
		defer span.End()

		makeCoffee(ctx)
	}()

	results := sess.Finalize(ctx)
	results.ZeroTimes()

	if !reflect.DeepEqual(results, expected) {
		t.Fatalf("results mismatch: %s", strings.Join(deep.Equal(results, expected), "\n"))
	}
	results.Write(os.Stdout, 0)
}

const coffeeLeft = 3

func makeCoffee(ctx context.Context) error {
	if coffeeLeft < 5 {
		Warn(ctx, "coffee getting low")
		Advise(ctx, getMoreCoffee)
	}

	// To mimic listener TLS checks, we'll see if we can nest a Test and add errors in the function
	Test(ctx, "prepare-kitchen", func(ctx context.Context) error {
		return Test(ctx, "build-microwave", func(ctx context.Context) error {
			buildMicrowave(ctx)
			return nil
		})
	})

	Test(ctx, "warm-milk", func(ctx context.Context) error {
		return warmMilk(ctx)
	})

	brewCoffee(ctx)

	SpotCheck(ctx, "pick-scone", pickScone)

	Test(ctx, "dispose-grounds", disposeGrounds)
	return nil
}

// buildMicrowave will throw an error in the function itself to fail the span,
// but will return nil so the caller test doesn't necessarily throw an error.
// The intended behavior is that the superspan will detect the failed subspan
// and fail regardless. This happens when Fail is used to fail the span, but not
// when Error is used. See the comment in the function itself.
func buildMicrowave(ctx context.Context) error {
	ctx, span := StartSpan(ctx, "buy-parts")

	Fail(ctx, "no stores sell microwave parts, please buy a microwave instead.")

	// The error line here does not actually yield an error in the output.
	// TODO: Debug this. In the meantime, always use Fail over Error.
	// Error(ctx, errors.New("no stores sell microwave parts, please buy a microwave instead."))

	Warn(ctx, "warning: you are about to try to build a microwave from scratch.")
	span.End()
	return nil
}

func warmMilk(ctx context.Context) error {
	// Always succeeds
	return nil
}

func brewCoffee(ctx context.Context) error {
	ctx, span := StartSpan(ctx, "brew-coffee")
	defer span.End()

	// Brewing happens here, successfully
	return nil
}

func pickScone() error {
	return errors.New("no scones")
}

func disposeGrounds(_ context.Context) error {
	// Done!
	return nil
}

func TestCapitalizeFirstLetter(t *testing.T) {
	s := "this is a test."
	if CapitalizeFirstLetter(s) != "This is a test." {
		t.Fatalf("first word of string was not capitalized: got %s", CapitalizeFirstLetter(s))
	}
	s = "this"
	if CapitalizeFirstLetter(s) != "This" {
		t.Fatalf("first word of string was not capitalized: got %s", CapitalizeFirstLetter(s))
	}
	s = "."
	if CapitalizeFirstLetter(s) != "." {
		t.Fatalf("String without letters was not unchanged: got %s", CapitalizeFirstLetter(s))
	}
}
