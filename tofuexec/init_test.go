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

func TestInitCmd_v1(t *testing.T) {
	td := t.TempDir()

	tf, err := NewOpenTofu(td, tfVersion(t, testutil.Latest_v1))
	if err != nil {
		t.Fatal(err)
	}

	// empty env, to avoid environ mismatch in testing
	tf.SetEnv(map[string]string{})

	t.Run("defaults", func(t *testing.T) {
		// defaults
		initCmd, err := tf.initCmd(context.Background())
		if err != nil {
			t.Fatal(err)
		}

		assertCmd(t, []string{
			"init",
			"-no-color",
			"-input=false",
			"-backend=true",
			"-get=true",
			"-upgrade=false",
		}, nil, initCmd)
	})

	t.Run("override all defaults", func(t *testing.T) {
		initCmd, err := tf.initCmd(context.Background(), Backend(false), BackendConfig("confpath1"), BackendConfig("confpath2"), FromModule("testsource"), Get(false), PluginDir("testdir1"), PluginDir("testdir2"), Reconfigure(true), Upgrade(true), Dir("initdir"))
		if err != nil {
			t.Fatal(err)
		}

		assertCmd(t, []string{
			"init",
			"-no-color",
			"-input=false",
			"-from-module=testsource",
			"-backend=false",
			"-get=false",
			"-upgrade=true",
			"-reconfigure",
			"-backend-config=confpath1",
			"-backend-config=confpath2",
			"-plugin-dir=testdir1",
			"-plugin-dir=testdir2",
			"initdir",
		}, nil, initCmd)
	})
}
