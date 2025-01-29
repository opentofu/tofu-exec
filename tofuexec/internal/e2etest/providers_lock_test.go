// Copyright (c) The OpenTofu Authors
// SPDX-License-Identifier: MPL-2.0
// Copyright (c) 2023 HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package e2etest

import (
	"context"
	"testing"

	"github.com/hashicorp/go-version"

	"github.com/opentofu/tofu-exec/tofuexec"
)

var (
	providersLockMinVersion = version.Must(version.NewVersion("0.14.0"))
)

func TestProvidersLock(t *testing.T) {
	runTest(t, "basic", func(t *testing.T, tfv *version.Version, tf *tofuexec.OpenTofu) {
		if tfv.LessThan(providersLockMinVersion) {
			t.Skip("tofu providers lock was added in Terraform 0.14, so test is not valid")
		}
		err := tf.Init(context.Background())
		if err != nil {
			t.Fatalf("error running Init in test directory: %s", err)
		}

		err = tf.ProvidersLock(context.Background())
		if err != nil {
			t.Fatalf("error running provider lock: %s", err)
		}
	})

}
