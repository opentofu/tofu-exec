// Copyright (c) The OpenTofu Authors
// SPDX-License-Identifier: MPL-2.0
// Copyright (c) 2023 HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package tfexec

import (
	"context"
	"testing"

	"github.com/opentofu/tofu-exec/tfexec/internal/testutil"
)

func TestFormatCmd(t *testing.T) {
	td := t.TempDir()

	tf, err := NewTofu(td, tfVersion(t, testutil.Latest_v1))
	if err != nil {
		t.Fatal(err)
	}

	// empty env, to avoid environ mismatch in testing
	_ = tf.SetEnv(map[string]string{})

	t.Run("defaults", func(t *testing.T) {
		fmtCmd, err := tf.formatCmd(context.Background(), []string{})
		if err != nil {
			t.Fatal(err)
		}

		assertCmd(t, []string{
			"fmt",
			"-no-color",
		}, nil, fmtCmd)
	})

	t.Run("override all defaults", func(t *testing.T) {
		fmtCmd, err := tf.formatCmd(context.Background(),
			[]string{"-diff"},
			Recursive(true),
			Dir("mydir"))
		if err != nil {
			t.Fatal(err)
		}

		assertCmd(t, []string{
			"fmt",
			"-no-color",
			"-diff",
			"-recursive",
			"--",
			"mydir",
		}, nil, fmtCmd)
	})
}
