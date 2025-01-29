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

func TestWorkspaceDeleteCmd(t *testing.T) {
	tf, err := NewOpenTofu(t.TempDir(), tfVersion(t, testutil.Latest_v1))
	if err != nil {
		t.Fatal(err)
	}

	// empty env, to avoid environ mismatch in testing
	tf.SetEnv(map[string]string{})

	t.Run("defaults", func(t *testing.T) {
		workspaceDeleteCmd, err := tf.workspaceDeleteCmd(context.Background(), "workspace-name")
		if err != nil {
			t.Fatal(err)
		}

		assertCmd(t, []string{
			"workspace", "delete",
			"-no-color",
			"workspace-name",
		}, nil, workspaceDeleteCmd)
	})

	t.Run("override all defaults", func(t *testing.T) {
		workspaceDeleteCmd, err := tf.workspaceDeleteCmd(context.Background(), "workspace-name",
			LockTimeout("200s"),
			Force(true),
			Lock(false))
		if err != nil {
			t.Fatal(err)
		}

		assertCmd(t, []string{
			"workspace", "delete",
			"-no-color",
			"-force",
			"-lock-timeout=200s",
			"-lock=false",
			"workspace-name",
		}, nil, workspaceDeleteCmd)
	})
}
