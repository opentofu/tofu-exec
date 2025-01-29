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

func TestOutputCmd(t *testing.T) {
	td := t.TempDir()

	tf, err := NewOpenTofu(td, tfVersion(t, testutil.Latest_v1))
	if err != nil {
		t.Fatal(err)
	}

	// empty env, to avoid environ mismatch in testing
	tf.SetEnv(map[string]string{})

	t.Run("defaults", func(t *testing.T) {
		outputCmd := tf.outputCmd(context.Background())

		assertCmd(t, []string{
			"output",
			"-no-color",
			"-json",
		}, nil, outputCmd)
	})

	t.Run("override all defaults", func(t *testing.T) {
		outputCmd := tf.outputCmd(context.Background(),
			State("teststate"))

		assertCmd(t, []string{
			"output",
			"-no-color",
			"-json",
			"-state=teststate",
		}, nil, outputCmd)
	})
}
