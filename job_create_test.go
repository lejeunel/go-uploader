package main

import (
	"errors"
	"testing"
)

func TestJobCreateScheme(t *testing.T) {
	testCases := []struct {
		name              string
		input             string
		output            string
		expectSchemeError bool
	}{
		{
			name:              "bad input scheme",
			input:             "badscheme://path/to/data/",
			output:            "scheme://path/to/data/",
			expectSchemeError: true,
		},
		{
			name:              "bad output scheme",
			input:             "file:///path/to/data/",
			output:            "badscheme://path/to/data/",
			expectSchemeError: true,
		},
		{
			name:              "good schemes",
			input:             "file:///path/to/data/",
			output:            "scheme://path/to/data/",
			expectSchemeError: false,
		},
	}
	jm := NewJobManager(*NewMockUploader())
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {

			_, err := jm.create(tc.input, tc.output)

			var got *schemeError
			isSchemeError := errors.As(err, &got)

			if tc.expectSchemeError && (!isSchemeError) {
				t.Fatalf("[%s] expected scheme error %v, got %T", tc.name, tc.expectSchemeError, err)
			} else if (!tc.expectSchemeError) && isSchemeError {
				t.Fatalf("[%s] did not expect scheme error %v, got %T", tc.name, tc.expectSchemeError, err)

			}

		})

	}
}

func TestJobCreateSourceError(t *testing.T) {
	testCases := []struct {
		name              string
		input             string
		output            string
		expectSourceError bool
	}{
		{
			name:              "good source",
			input:             "file:///path/to/data/",
			output:            "scheme://path/to/data/",
			expectSourceError: false,
		},
		{
			name:              "bad source",
			input:             "file:///non-existing-path/to/data/",
			output:            "scheme://path/to/data/",
			expectSourceError: true,
		},
	}
	jm := NewJobManager(*NewMockUploader())

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {

			_, err := jm.create(tc.input, tc.output)

			var got *sourceError
			isSourceError := errors.As(err, &got)

			if tc.expectSourceError && (!isSourceError) {
				t.Fatalf("[%s] expected source error %v, got %T", tc.name, tc.expectSourceError, err)
			} else if (!tc.expectSourceError) && isSourceError {
				t.Fatalf("[%s] did not expect source error %v, got %T", tc.name, tc.expectSourceError, err)

			}

		})

	}

}
