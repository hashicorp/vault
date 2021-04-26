package diagnose

import (
	"context"
	"errors"
	"os"
	"reflect"
	"testing"
)

func TestDiagnoseOtelResults(t *testing.T) {
	expected := &Result{
		Name:   "make-coffee",
		Status: warnStatus,
		Warnings: []string{
			"coffee getting low",
			"no scones",
		},
		Children: []*Result{
			{
				Name:   "warm-milk",
				Status: okStatus,
			},
			{
				Name:   "brew-coffee",
				Status: okStatus,
			},
		},
	}
	Init()

	func() {
		ctx, span := StartSpan(context.Background(), "make-coffee")
		defer span.End()

		makeCoffee(ctx)
	}()

	results := Shutdown()
	if !reflect.DeepEqual(results, expected) {
		t.Fail()
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

	err = pickScone(ctx)
	if err != nil {
		Warn(ctx, err.Error())
	}

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

func pickScone(ctx context.Context) error {
	return errors.New("no scones")
}
