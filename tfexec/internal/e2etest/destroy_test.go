// Copyright (c) The OpenTofu Authors
// SPDX-License-Identifier: MPL-2.0
// Copyright (c) 2023 HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package e2etest

import (
	"context"
	"io"
	"regexp"
	"testing"

	"github.com/hashicorp/go-version"

	"github.com/opentofu/tofu-exec/tfexec"
	"github.com/opentofu/tofu-exec/tfexec/internal/testutil"
)

func TestDestroy(t *testing.T) {
	runTest(t, "basic", func(t *testing.T, tfv *version.Version, tf *tfexec.Tofu) {
		err := tf.Init(context.Background())
		if err != nil {
			t.Fatalf("error running Init in test directory: %s", err)
		}

		err = tf.Apply(context.Background())
		if err != nil {
			t.Fatalf("error running Apply: %s", err)
		}

		err = tf.Destroy(context.Background())
		if err != nil {
			t.Fatalf("error running Destroy: %s", err)
		}
	})
}

func TestDestroyJSON_TF014AndEarlier(t *testing.T) {
	versions := []string{testutil.Latest013, testutil.Latest014}

	runTestWithVersions(t, versions, "basic", func(t *testing.T, tfv *version.Version, tf *tfexec.Tofu) {
		err := tf.Init(context.Background())
		if err != nil {
			t.Fatalf("error running Init in test directory: %s", err)
		}

		re := regexp.MustCompile("terraform destroy -json was added in 0.15.3")

		err = tf.DestroyJSON(context.Background(), io.Discard)
		if err != nil && !re.MatchString(err.Error()) {
			t.Fatalf("error running Apply: %s", err)
		}
	})
}

func TestDestroyJSON_TF015AndLater(t *testing.T) {
	versions := []string{testutil.Latest015, testutil.Latest_v1}

	runTestWithVersions(t, versions, "basic", func(t *testing.T, tfv *version.Version, tf *tfexec.Tofu) {
		err := tf.Init(context.Background())
		if err != nil {
			t.Fatalf("error running Init in test directory: %s", err)
		}

		err = tf.DestroyJSON(context.Background(), io.Discard)
		if err != nil {
			t.Fatalf("error running Apply: %s", err)
		}
	})
}
