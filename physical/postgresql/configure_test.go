package postgresql

import (
	"os"
	"testing"
)

func TestConnectionURL(t *testing.T) {
	type input struct {
		envar string
		conf  map[string]string
	}

	var cases = []struct {
		name  string
		want  string
		input input
	}{
		{
			name: "environment_variable_not_set_use_config_value",
			want: "abc",
			input: input{
				envar: "",
				conf:  map[string]string{"connection_url": "abc"},
			},
		},
		{
			name: "no_value_connection_url_set_key_exists",
			want: "",
			input: input{
				envar: "",
				conf:  map[string]string{"connection_url": ""},
			},
		},
		{
			name: "no_value_connection_url_set_key_doesnt_exist",
			want: "",
			input: input{
				envar: "",
				conf:  map[string]string{},
			},
		},
		{
			name: "environment_variable_set",
			want: "abc",
			input: input{
				envar: "abc",
				conf:  map[string]string{"connection_url": "def"},
			},
		},
	}

	for _, tt := range cases {
		t.Run(tt.name, func(t *testing.T) {
			os.Setenv("PG_CONNECTION_URL", tt.input.envar)
			defer os.Setenv("PG_CONNECTION_URL", "")

			got := connectionURL(tt.input.conf)

			if got != tt.want {
				t.Errorf("connectionURL(%s): want '%s', got '%s'", tt.input, tt.want, got)
			}
		})
	}
}
