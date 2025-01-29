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

	"github.com/opentofu/tofu-exec/tofuexec"
	"github.com/opentofu/tofu-exec/tofuexec/internal/testutil"
)

func TestRefresh(t *testing.T) {
	runTest(t, "basic", func(t *testing.T, tfv *version.Version, tf *tofuexec.OpenTofu) {
		err := tf.Init(context.Background())
		if err != nil {
			t.Fatalf("error running Init in test directory: %s", err)
		}

		err = tf.Apply(context.Background())
		if err != nil {
			t.Fatalf("error running Apply: %s", err)
		}

		err = tf.Refresh(context.Background())
		if err != nil {
			t.Fatalf("error running Refresh: %s", err)
		}
	})
}

func TestRefreshJSON_Tofu16AndLater(t *testing.T) {
	versions := []string{testutil.Latest_v1, testutil.Latest_v1_7}

	runTestWithVersions(t, versions, "basic", func(t *testing.T, tfv *version.Version, tf *tofuexec.OpenTofu) {
		err := tf.Init(context.Background())
		if err != nil {
			t.Fatalf("error running Init in test directory: %s", err)
		}

		err = tf.RefreshJSON(context.Background(), io.Discard)
		if err != nil {
			t.Fatalf("error running Apply: %s", err)
		}
	})
}
