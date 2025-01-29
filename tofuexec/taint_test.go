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

func TestTaintCmd(t *testing.T) {
	td := t.TempDir()

	tf, err := NewOpenTofu(td, tfVersion(t, testutil.Latest_v1))
	if err != nil {
		t.Fatal(err)
	}

	// empty env, to avoid environ mismatch in testing
	tf.SetEnv(map[string]string{})

	t.Run("defaults", func(t *testing.T) {
		taintCmd := tf.taintCmd(context.Background(), "aws_instance.foo")

		assertCmd(t, []string{
			"taint",
			"-no-color",
			"-lock=true",
			"aws_instance.foo",
		}, nil, taintCmd)
	})

	t.Run("override all defaults", func(t *testing.T) {
		taintCmd := tf.taintCmd(context.Background(), "aws_instance.foo",
			State("teststate"),
			AllowMissing(true),
			LockTimeout("200s"),
			Lock(false))

		assertCmd(t, []string{
			"taint",
			"-no-color",
			"-lock-timeout=200s",
			"-state=teststate",
			"-lock=false",
			"-allow-missing",
			"aws_instance.foo",
		}, nil, taintCmd)
	})
}
