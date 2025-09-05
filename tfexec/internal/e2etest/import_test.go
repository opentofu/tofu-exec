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

func TestImport(t *testing.T) {
	const (
		expectedID      = "asdlfjksdlfkjsdlfk"
		resourceAddress = "random_string.random_string"
	)

	runTest(t, "import", func(t *testing.T, tfv *version.Version, tf *tfexec.Tofu) {
		ctx := context.Background()

		err := tf.Init(ctx)
		if err != nil {
			t.Fatal(err)
		}

		// Config is unnecessary here since its already the working dir, but just testing an additional flag
		err = tf.Import(ctx, resourceAddress, expectedID, tfexec.DisableBackup(), tfexec.Config(tf.WorkingDir()))
		if err != nil {
			t.Fatal(err)
		}

		state, err := tf.Show(ctx)
		if err != nil {
			t.Fatal(err)
		}

		for _, r := range state.Values.RootModule.Resources {
			if r.Address != resourceAddress {
				continue
			}

			raw, ok := r.AttributeValues["id"]
			if !ok {
				t.Fatal("value not found for \"id\" attribute")
			}
			actual, ok := raw.(string)
			if !ok {
				t.Fatalf("unable to cast %T to string: %#v", raw, raw)
			}

			if actual != expectedID {
				t.Fatalf("expected %q, got %q", expectedID, actual)
			}

			// success
			return
		}

		t.Fatalf("imported resource %q not found", resourceAddress)
	})
}
