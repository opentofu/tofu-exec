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

func TestStatePushCmd(t *testing.T) {
	tf, err := NewTofu(t.TempDir(), tfVersion(t, testutil.Latest_v1))
	if err != nil {
		t.Fatal(err)
	}

	tf.SetEnv(map[string]string{})

	t.Run("defaults", func(t *testing.T) {
		statePushCmd, err := tf.statePushCmd(context.Background(), "testpath")
		if err != nil {
			t.Fatal(err)
		}

		assertCmd(t, []string{
			"state",
			"push",
			"-lock=false",
			"-lock-timeout=0s",
			"testpath",
		}, nil, statePushCmd)
	})

	t.Run("override all defaults", func(t *testing.T) {
		statePushCmd, err := tf.statePushCmd(context.Background(), "testpath", Force(true), Lock(true), LockTimeout("10s"))
		if err != nil {
			t.Fatal(err)
		}

		assertCmd(t, []string{
			"state",
			"push",
			"-force",
			"-lock=true",
			"-lock-timeout=10s",
			"testpath",
		}, nil, statePushCmd)
	})
}
