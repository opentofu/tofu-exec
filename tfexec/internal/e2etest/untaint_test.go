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

func TestUntaint(t *testing.T) {
	runTest(t, "basic", func(t *testing.T, tfv *version.Version, tf *tfexec.Tofu) {
		err := tf.Init(context.Background())
		if err != nil {
			t.Fatalf("error running Init in test directory: %s", err)
		}

		err = tf.Apply(context.Background())
		if err != nil {
			t.Fatalf("error running Apply: %s", err)
		}

		err = tf.Taint(context.Background(), "null_resource.foo")
		if err != nil {
			t.Fatalf("error running Taint: %s", err)
		}

		err = tf.Untaint(context.Background(), "null_resource.foo")
		if err != nil {
			t.Fatalf("error running Untaint: %s", err)
		}
	})
}
