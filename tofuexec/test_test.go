// Copyright (c) The OpenTofu Authors
// SPDX-License-Identifier: MPL-2.0
// Copyright (c) 2023 HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package tofuexec

import (
	"context"
	"testing"

	"github.com/opentofu/tofu-exec/tofuexec/internal/testutil"
)

func TestTestCmd(t *testing.T) {
	td := t.TempDir()

	tf, err := NewOpenTofu(td, tfVersion(t, testutil.Latest_v1_7))

	if err != nil {
		t.Fatal(err)
	}

	// empty env, to avoid environ mismatch in testing
	tf.SetEnv(map[string]string{})

	t.Run("defaults", func(t *testing.T) {
		testCmd := tf.testCmd(context.Background())

		assertCmd(t, []string{
			"test",
			"-json",
		}, nil, testCmd)
	})

	t.Run("override all defaults", func(t *testing.T) {
		testCmd := tf.testCmd(context.Background(), TestsDirectory("test"))

		assertCmd(t, []string{
			"test",
			"-json",
			"-tests-directory=test",
		}, nil, testCmd)
	})
}
