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

func TestOutput_noOutputs(t *testing.T) {
	runTest(t, "basic", func(t *testing.T, tfv *version.Version, tf *tfexec.Tofu) {
		if tfv.LessThan(version.Must(version.NewVersion("0.12.14"))) {
			// https://github.com/hashicorp/terraform/blob/v0.12/CHANGELOG.md#01214-november-13-2019
			t.Skip("no outputs being success (instead of error) was changed in 0.12.14")
		}

		err := tf.Init(context.Background())
		if err != nil {
			t.Fatalf("error running Init in test directory: %s", err)
		}

		_, err = tf.Output(context.Background())
		if err != nil {
			t.Fatalf("error running Output: %s", err)
		}
	})
}
