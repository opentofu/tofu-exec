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

func TestMetadataFunctions(t *testing.T) {
	runTest(t, "basic", func(t *testing.T, tfv *version.Version, tf *tfexec.Terraform) {
		if tfv.LessThan(metadataFunctionsMinVersion) {
			t.Skip("metadata functions command is not available in this Terraform version")
		}

		_, err := tf.MetadataFunctions(context.Background())
		if err != nil {
			t.Fatalf("error running MetadataFunctions: %s", err)
		}
	})
}
