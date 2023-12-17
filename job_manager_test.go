package main

import (
	"errors"
	"testing"
)

func TestBadSchemes(t *testing.T) {
	testCases := []struct {
		name   string
		input  string
		output string
	}{
		{
			name:   "bad input scheme",
			input:  "badscheme://path/to/data/",
			output: "scheme://path/to/data/",
		},
		{
			name:   "bad output scheme",
			input:  "file:///path/to/data/",
			output: "badscheme://path/to/data/",
		},
	}
	for _, tc := range testCases {
		jm := NewJobManager(*NewMockUploader())
		_, err := jm.create(tc.input, tc.output)

		var got *schemeError
		isSchemeError := errors.As(err, &got)

		if !isSchemeError {
			t.Fatalf("[%s] expected a schemeError, got %T", tc.name, err)
		}

	}
}

func TestCreate(t *testing.T) {
	jm := NewJobManager(*NewMockUploader())
	_, err := jm.create("file:///path/to/dir/", "scheme://path/to/data/")
	if err != nil {
		t.Fatalf("expected null error creating job: %v", err)
	}

}
