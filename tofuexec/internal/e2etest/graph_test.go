// Copyright (c) The OpenTofu Authors
// SPDX-License-Identifier: MPL-2.0
// Copyright (c) 2023 HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package e2etest

import (
	"context"
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/hashicorp/go-version"

	"github.com/opentofu/tofu-exec/tofuexec"
)

func TestGraph(t *testing.T) {
	runTest(t, "basic", func(t *testing.T, tfv *version.Version, tf *tofuexec.OpenTofu) {
		err := tf.Init(context.Background())
		if err != nil {
			t.Fatalf("error running Init in test directory: %s", err)
		}

		err = tf.Apply(context.Background())
		if err != nil {
			t.Fatalf("error running Apply: %s", err)
		}

		graphOutput, err := tf.Graph(context.Background())
		if err != nil {
			t.Fatalf("error running Graph: %s", err)
		}

		if diff := cmp.Diff(expectedGraphOutput(tfv), graphOutput); diff != "" {
			t.Fatalf("Graph output does not match: %s", diff)
		}
	})
}

func expectedGraphOutput(tfv *version.Version) string {
	// 1.1.0+
	return `digraph {
	compound = "true"
	newrank = "true"
	subgraph "root" {
		"[root] null_resource.foo (expand)" [label = "null_resource.foo", shape = "box"]
		"[root] provider[\"registry.opentofu.org/hashicorp/null\"]" [label = "provider[\"registry.opentofu.org/hashicorp/null\"]", shape = "diamond"]
		"[root] null_resource.foo (expand)" -> "[root] provider[\"registry.opentofu.org/hashicorp/null\"]"
		"[root] provider[\"registry.opentofu.org/hashicorp/null\"] (close)" -> "[root] null_resource.foo (expand)"
		"[root] root" -> "[root] provider[\"registry.opentofu.org/hashicorp/null\"] (close)"
	}
}

`
}
