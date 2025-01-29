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

func TestGraphCmd_v1(t *testing.T) {
	td := t.TempDir()

	tf, err := NewOpenTofu(td, tfVersion(t, testutil.Latest_v1))
	if err != nil {
		t.Fatal(err)
	}

	// empty env, to avoid environ mismatch in testing
	tf.SetEnv(map[string]string{})

	t.Run("defaults", func(t *testing.T) {
		graphCmd, _ := tf.graphCmd(context.Background())

		assertCmd(t, []string{
			"graph",
		}, nil, graphCmd)
	})

	t.Run("override all defaults", func(t *testing.T) {
		graphCmd, _ := tf.graphCmd(context.Background(),
			GraphPlan("teststate"),
			DrawCycles(true),
			GraphType("output"))

		assertCmd(t, []string{
			"graph",
			"-plan=teststate",
			"-draw-cycles",
			"-type=output",
		}, nil, graphCmd)
	})
}
