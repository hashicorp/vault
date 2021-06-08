package diagnose

import (
	"context"
	"errors"
	"github.com/go-test/deep"
	"os"
	"reflect"
	"strings"
	"testing"
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
	sess.SetSkipList([]string{"dispose-grounds"})
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

	err := Test(ctx, "warm-milk", warmMilk)
	if err != nil {
		return err
	}

	err = brewCoffee(ctx)
	if err != nil {
		return err
	}

	SpotCheck(ctx, "pick-scone", pickScone)

	Test(ctx, "dispose-grounds", Skippable("dispose-grounds", disposeGrounds))
	return nil
}

func warmMilk(ctx context.Context) error {
	// Always succeeds
	return nil
}

func brewCoffee(ctx context.Context) error {
	ctx, span := StartSpan(ctx, "brew-coffee")
	defer span.End()

	//Brewing happens here, successfully
	return nil
}

func pickScone() error {
	return errors.New("no scones")
}

func disposeGrounds(_ context.Context) error {
	//Done!
	return nil
}
