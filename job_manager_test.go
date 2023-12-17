package main

import (
	"errors"
	"testing"
)

func TestBadSchemes(t *testing.T) {
	testCases := []struct {
		name                string
		input               string
		output              string
		expect_scheme_error bool
	}{
		{
			name:                "bad input scheme",
			input:               "badscheme://path/to/data/",
			output:              "scheme://path/to/data/",
			expect_scheme_error: true,
		},
		{
			name:                "bad output scheme",
			input:               "file:///path/to/data/",
			output:              "badscheme://path/to/data/",
			expect_scheme_error: true,
		},
		{
			name:                "good schemes",
			input:               "file:///path/to/data/",
			output:              "scheme://path/to/data/",
			expect_scheme_error: false,
		},
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {

			jm := NewJobManager(*NewMockUploader())
			_, err := jm.create(tc.input, tc.output)

			var got *schemeError
			isSchemeError := errors.As(err, &got)

			if tc.expect_scheme_error && (!isSchemeError) {
				t.Fatalf("[%s] expected scheme error %v, got %T", tc.name, tc.expect_scheme_error, err)
			} else if (!tc.expect_scheme_error) && isSchemeError {
				t.Fatalf("[%s] did not expect scheme error %v, got %T", tc.name, tc.expect_scheme_error, err)

			}

		})

	}
}
