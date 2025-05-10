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

func TestStateMvCmd(t *testing.T) {
	td := t.TempDir()

	tf, err := NewTofu(td, tfVersion(t, testutil.Latest_v1))
	if err != nil {
		t.Fatal(err)
	}

	// empty env, to avoid environ mismatch in testing
	tf.SetEnv(map[string]string{})

	t.Run("defaults", func(t *testing.T) {
		stateMvCmd, err := tf.stateMvCmd(context.Background(), "testsource", "testdestination")
		if err != nil {
			t.Fatal(err)
		}

		assertCmd(t, []string{
			"state",
			"mv",
			"-no-color",
			"-lock-timeout=0s",
			"-lock=true",
			"testsource",
			"testdestination",
		}, nil, stateMvCmd)
	})

	t.Run("override all defaults", func(t *testing.T) {
		stateMvCmd, err := tf.stateMvCmd(context.Background(), "testsrc", "testdest", Backup("testbackup"), BackupOut("testbackupout"), LockTimeout("200s"), State("teststate"), StateOut("teststateout"), Lock(false))
		if err != nil {
			t.Fatal(err)
		}

		assertCmd(t, []string{
			"state",
			"mv",
			"-no-color",
			"-backup=testbackup",
			"-backup-out=testbackupout",
			"-lock-timeout=200s",
			"-state=teststate",
			"-state-out=teststateout",
			"-lock=false",
			"testsrc",
			"testdest",
		}, nil, stateMvCmd)
	})
}
