// Copyright (c) The OpenTofu Authors
// SPDX-License-Identifier: MPL-2.0
// Copyright (c) 2023 HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package e2etest

import (
	"context"
	"testing"

	"github.com/hashicorp/go-version"

	"github.com/opentofu/tofu-exec/tfexec"
)

// LockID set in the test fixture
const inmemLockID = "2b6a6738-5dd5-50d6-c0ae-f6352977666b"

func TestForceUnlock(t *testing.T) {
	runTest(t, "inmem_backend_locked", func(t *testing.T, tfv *version.Version, tf *tfexec.Tofu) {
		err := tf.Init(context.Background())
		if err != nil {
			t.Fatalf("error running Init: %v", err)
		}

		err = tf.ForceUnlock(context.Background(), inmemLockID)
		if err != nil {
			t.Fatalf("error running ForceUnlock: %v", err)
		}
	})
	runTest(t, "inmem_backend_locked", func(t *testing.T, tfv *version.Version, tf *tfexec.Tofu) {
		err := tf.Init(context.Background())
		if err != nil {
			t.Fatalf("error running Init: %v", err)
		}

		err = tf.ForceUnlock(context.Background(), "badlockid")
		if err == nil {
			t.Fatalf("expected error when running ForceUnlock with invalid lock id")
		}
	})
}
