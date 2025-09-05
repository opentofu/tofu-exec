// Copyright (c) The OpenTofu Authors
// SPDX-License-Identifier: MPL-2.0
// Copyright (c) 2023 HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package e2etest

import (
	"context"
	"io"
	"regexp"
	"testing"

	"github.com/hashicorp/go-version"

	"github.com/opentofu/tofu-exec/tfexec"
)

func TestTest(t *testing.T) {
	runTest(t, "test_command_passing", func(t *testing.T, tfv *version.Version, tf *tfexec.Tofu) {
		err := tf.Test(context.Background(), nil)

		if err != nil {
			t.Fatalf("error running test command: %s", err)
		}
	})
}

func TestTestError(t *testing.T) {
	runTest(t, "test_command_failing", func(t *testing.T, tfv *version.Version, tf *tfexec.Tofu) {
		err := tf.Test(context.Background(), io.Discard)

		if err == nil {
			t.Fatal("expected error, got none")
		}

		got := err.Error()
		expected := regexp.MustCompile("exit status 1")

		if !expected.MatchString(got) {
			t.Fatalf("expected error matching '%s', got: %s", expected, got)
		}
	})
}
