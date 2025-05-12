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

func TestForceUnlockCmd(t *testing.T) {
	td := t.TempDir()

	tf, err := NewTofu(td, tfVersion(t, testutil.Latest_v1_1))
	if err != nil {
		t.Fatal(err)
	}

	// empty env, to avoid environ mismatch in testing
	tf.SetEnv(map[string]string{})

	t.Run("defaults", func(t *testing.T) {
		forceUnlockCmd, err := tf.forceUnlockCmd(context.Background(), "12345")
		if err != nil {
			t.Fatal(err)
		}

		assertCmd(t, []string{
			"force-unlock",
			"-no-color",
			"-force",
			"12345",
		}, nil, forceUnlockCmd)
	})
}
