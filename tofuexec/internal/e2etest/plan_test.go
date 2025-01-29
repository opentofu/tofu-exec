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

func TestPlan(t *testing.T) {
	runTest(t, "basic", func(t *testing.T, tfv *version.Version, tf *tofuexec.OpenTofu) {
		err := tf.Init(context.Background())
		if err != nil {
			t.Fatalf("error running Init in test directory: %s", err)
		}

		hasChanges, err := tf.Plan(context.Background())
		if err != nil {
			t.Fatalf("error running Plan: %s", err)
		}
		if !hasChanges {
			t.Fatalf("expected: true, got: %t", hasChanges)
		}
	})

}

func TestPlanWithState(t *testing.T) {
	runTest(t, "basic_with_state", func(t *testing.T, tfv *version.Version, tf *tofuexec.OpenTofu) {
		if tfv.LessThan(providerAddressMinVersion) {
			t.Skip("state file provider FQNs not compatible with this OpenTofu version")
		}
		err := tf.Init(context.Background())
		if err != nil {
			t.Fatalf("error running Init in test directory: %s", err)
		}

		hasChanges, err := tf.Plan(context.Background())
		if err != nil {
			t.Fatalf("error running Plan: %s", err)
		}
		if hasChanges {
			t.Fatalf("expected: false, got: %t", hasChanges)
		}
	})
}

func TestPlanJSON_Tofu16AndLater(t *testing.T) {
	versions := []string{testutil.Latest_v1, testutil.Latest_v1_7}

	runTestWithVersions(t, versions, "basic", func(t *testing.T, tfv *version.Version, tf *tofuexec.OpenTofu) {
		err := tf.Init(context.Background())
		if err != nil {
			t.Fatalf("error running Init in test directory: %s", err)
		}

		hasChanges, err := tf.PlanJSON(context.Background(), io.Discard)
		if err != nil {
			t.Fatalf("error running Apply: %s", err)
		}
		if !hasChanges {
			t.Fatalf("expected: true, got: %t", hasChanges)
		}
	})
}
