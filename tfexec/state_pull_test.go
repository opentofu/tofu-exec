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

func TestStatePull(t *testing.T) {
	tf, err := NewTofu(t.TempDir(), tfVersion(t, testutil.Latest_v1))
	if err != nil {
		t.Fatal(err)
	}

	tf.SetEnv(map[string]string{})

	t.Run("tfstate", func(t *testing.T) {
		statePullCmd := tf.statePullCmd(context.Background(), nil)

		assertCmd(t, []string{
			"state",
			"pull",
		}, nil, statePullCmd)
	})
}
