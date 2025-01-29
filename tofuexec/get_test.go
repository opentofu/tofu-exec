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

func TestGetCmd(t *testing.T) {
	td := t.TempDir()

	tf, err := NewOpenTofu(td, tfVersion(t, testutil.Latest_v1))
	if err != nil {
		t.Fatal(err)
	}

	// empty env, to avoid environ mismatch in testing
	tf.SetEnv(map[string]string{})

	t.Run("basic", func(t *testing.T) {
		getCmd, err := tf.getCmd(context.Background())
		if err != nil {
			t.Fatal(err)
		}

		assertCmd(t, []string{
			"get",
			"-no-color",
			"-update=false",
		}, nil, getCmd)
	})
}
