// Copyright (c) The OpenTofu Authors
// SPDX-License-Identifier: MPL-2.0
// Copyright (c) 2023 HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package e2etest

import (
	"context"
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/hashicorp/go-version"

	"github.com/opentofu/tofu-exec/jsontypes"
	"github.com/opentofu/tofu-exec/tfexec"
)

func TestValidate(t *testing.T) {
	runTest(t, "basic", func(t *testing.T, tfv *version.Version, tf *tfexec.Tofu) {
		err := tf.Init(context.Background())
		if err != nil {
			t.Fatal(err)
		}

		validation, err := tf.Validate(context.Background())
		if err != nil {
			t.Fatal(err)
		}

		if !validation.Valid {
			t.Fatalf("expected valid, got %#v", validation)
		}
	})

	runTest(t, "invalid", func(t *testing.T, tfv *version.Version, tf *tfexec.Tofu) {
		err := tf.Init(context.Background())
		if err != nil {
			t.Logf("error initializing: %s", err)
		}

		var expectedDiags []jsontypes.Diagnostic

		expectedDiags = []jsontypes.Diagnostic{
			{
				Severity: "error",
				Summary:  "Unsupported block type",
				Detail:   "Blocks of type \"bad_block\" are not expected here.",
				Range: &jsontypes.Range{
					Filename: "main.tf",
					Start: jsontypes.Pos{
						Line:   1,
						Column: 1,
					},
					End: jsontypes.Pos{
						Line:   1,
						Column: 10,
					},
				},
				Snippet: &jsontypes.DiagnosticSnippet{
					Code:                 "bad_block {",
					StartLine:            1,
					HighlightStartOffset: 0,
					HighlightEndOffset:   9,
					Values:               []jsontypes.DiagnosticExpressionValue{},
				},
			},
			{
				Severity: "error",
				Summary:  "Unsupported argument",
				Detail:   "An argument named \"bad_attribute\" is not expected here.",
				Range: &jsontypes.Range{
					Filename: "main.tf",
					Start: jsontypes.Pos{
						Line:   5,
						Column: 5,
					},
					End: jsontypes.Pos{
						Line:   5,
						Column: 18,
					},
				},
				Snippet: &jsontypes.DiagnosticSnippet{
					Context:              ptrToString("terraform"),
					Code:                 "    bad_attribute = \"string\"",
					StartLine:            5,
					HighlightStartOffset: 4,
					HighlightEndOffset:   17,
					Values:               []jsontypes.DiagnosticExpressionValue{},
				},
			},
		}

		actual, err := tf.Validate(context.Background())
		if err != nil {
			t.Fatal(err)
		}

		// reset byte locations in actual as CRLF issues render them off between operating systems
		var cleanActual []jsontypes.Diagnostic
		for _, diag := range actual.Diagnostics {
			diag.Range.Start.Byte = 0
			diag.Range.End.Byte = 0
			cleanActual = append(cleanActual, diag)
		}

		if diff := cmp.Diff(expectedDiags, cleanActual); diff != "" {
			t.Fatalf("diags do not match: %s", diff)
		}
	})
}

func ptrToString(value string) *string {
	return &value
}
