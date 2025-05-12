// Copyright (c) The OpenTofu Authors
// SPDX-License-Identifier: MPL-2.0
// Copyright (c) 2023 HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package e2etest

import (
	"context"
	"io"
	"testing"

	"github.com/hashicorp/go-version"

	"github.com/opentofu/tofu-exec/tfexec"
	"github.com/opentofu/tofu-exec/tfexec/internal/testutil"
)

var (
	applyDestroyMinVersion = version.Must(version.NewVersion("0.15.2"))
)

func TestApply(t *testing.T) {
	runTest(t, "basic", func(t *testing.T, tfv *version.Version, tf *tfexec.Tofu) {
		err := tf.Init(context.Background())
		if err != nil {
			t.Fatalf("error running Init in test directory: %s", err)
		}

		err = tf.Apply(context.Background())
		if err != nil {
			t.Fatalf("error running Apply: %s", err)
		}
	})
}

func TestApplyJSON_TF015AndLater(t *testing.T) {
	versions := []string{testutil.Latest015, testutil.Latest_v1, testutil.Latest_v1_1}

	runTestWithVersions(t, versions, "basic", func(t *testing.T, tfv *version.Version, tf *tfexec.Tofu) {
		err := tf.Init(context.Background())
		if err != nil {
			t.Fatalf("error running Init in test directory: %s", err)
		}

		err = tf.ApplyJSON(context.Background(), io.Discard)
		if err != nil {
			t.Fatalf("error running Apply: %s", err)
		}
	})
}

func TestApplyDestroy(t *testing.T) {
	runTest(t, "basic", func(t *testing.T, tfv *version.Version, tf *tfexec.Tofu) {
		if tfv.LessThan(applyDestroyMinVersion) {
			t.Skip("terraform apply -destroy was added in Terraform 0.15.2, so test is not valid")
		}
		err := tf.Init(context.Background())
		if err != nil {
			t.Fatalf("error running Init in test directory: %s", err)
		}

		err = tf.Apply(context.Background())
		if err != nil {
			t.Fatalf("error running Apply: %s", err)
		}

		err = tf.Apply(context.Background(), tfexec.Destroy(true))
		if err != nil {
			t.Fatalf("error running Apply -destroy: %s", err)
		}
	})
}
