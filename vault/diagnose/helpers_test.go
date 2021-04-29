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

func TestDiagnoseOtelResults(t *testing.T) {
	expected := &Result{
		Name:   "make-coffee",
		Status: WarningStatus,
		Warnings: []string{
			"coffee getting low",
		},
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
		},
	}
	sess := New()
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
	results.Write(os.Stdout)
}

const coffeeLeft = 3

func makeCoffee(ctx context.Context) error {
	if coffeeLeft < 5 {
		Warn(ctx, "coffee getting low")
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
