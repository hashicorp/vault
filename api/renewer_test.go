package api

import (
	"reflect"
	"testing"
)

func TestRenewer_NewRenewer(t *testing.T) {
	t.Parallel()

	client, err := NewClient(DefaultConfig())
	if err != nil {
		t.Fatal(err)
	}

	cases := []struct {
		name string
		i    *RenewerInput
		e    *Renewer
		err  bool
	}{
		{
			"nil",
			nil,
			nil,
			true,
		},
		{
			"missing_secret",
			&RenewerInput{
				Secret: nil,
			},
			nil,
			true,
		},
		{
			"default_grace",
			&RenewerInput{
				Secret: &Secret{},
			},
			&Renewer{
				secret: &Secret{},
			},
			false,
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			v, err := client.NewRenewer(tc.i)
			if (err != nil) != tc.err {
				t.Fatal(err)
			}

			if v == nil {
				return
			}

			// Zero-out channels because reflect
			v.client = nil
			v.random = nil
			v.doneCh = nil
			v.renewCh = nil
			v.stopCh = nil

			if !reflect.DeepEqual(tc.e, v) {
				t.Errorf("not equal\nexp: %#v\nact: %#v", tc.e, v)
			}
		})
	}
}
