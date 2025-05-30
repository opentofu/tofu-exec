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

func TestRefreshCmd(t *testing.T) {
	td := t.TempDir()

	tf, err := NewTofu(td, tfVersion(t, testutil.Latest_v1))
	if err != nil {
		t.Fatal(err)
	}

	// empty env, to avoid environ mismatch in testing
	tf.SetEnv(map[string]string{})

	t.Run("defaults", func(t *testing.T) {
		refreshCmd, err := tf.refreshCmd(context.Background())
		if err != nil {
			t.Fatal(err)
		}

		assertCmd(t, []string{
			"refresh",
			"-no-color",
			"-input=false",
			"-lock-timeout=0s",
			"-lock=true",
		}, nil, refreshCmd)
	})

	t.Run("override all defaults", func(t *testing.T) {
		refreshCmd, err := tf.refreshCmd(context.Background(), Backup("testbackup"), LockTimeout("200s"), State("teststate"), StateOut("teststateout"), VarFile("testvarfile"), Lock(false), Target("target1"), Target("target2"), Var("var1=foo"), Var("var2=bar"), Dir("refreshdir"))
		if err != nil {
			t.Fatal(err)
		}

		assertCmd(t, []string{
			"refresh",
			"-no-color",
			"-input=false",
			"-backup=testbackup",
			"-lock-timeout=200s",
			"-state=teststate",
			"-state-out=teststateout",
			"-var-file=testvarfile",
			"-lock=false",
			"-target=target1",
			"-target=target2",
			"-var", "var1=foo",
			"-var", "var2=bar",
			"refreshdir",
		}, nil, refreshCmd)
	})
}

func TestRefreshJSONCmd(t *testing.T) {
	td := t.TempDir()

	tf, err := NewTofu(td, tfVersion(t, testutil.Latest_v1))
	if err != nil {
		t.Fatal(err)
	}

	// empty env, to avoid environ mismatch in testing
	tf.SetEnv(map[string]string{})

	t.Run("defaults", func(t *testing.T) {
		refreshCmd, err := tf.refreshJSONCmd(context.Background())
		if err != nil {
			t.Fatal(err)
		}

		assertCmd(t, []string{
			"refresh",
			"-no-color",
			"-input=false",
			"-lock-timeout=0s",
			"-lock=true",
			"-json",
		}, nil, refreshCmd)
	})

	t.Run("override all defaults", func(t *testing.T) {
		refreshCmd, err := tf.refreshJSONCmd(context.Background(), Backup("testbackup"), LockTimeout("200s"), State("teststate"), StateOut("teststateout"), VarFile("testvarfile"), Lock(false), Target("target1"), Target("target2"), Var("var1=foo"), Var("var2=bar"), Dir("refreshdir"))
		if err != nil {
			t.Fatal(err)
		}

		assertCmd(t, []string{
			"refresh",
			"-no-color",
			"-input=false",
			"-backup=testbackup",
			"-lock-timeout=200s",
			"-state=teststate",
			"-state-out=teststateout",
			"-var-file=testvarfile",
			"-lock=false",
			"-target=target1",
			"-target=target2",
			"-var", "var1=foo",
			"-var", "var2=bar",
			"-json",
			"refreshdir",
		}, nil, refreshCmd)
	})
}
