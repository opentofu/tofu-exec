// Copyright (c) The OpenTofu Authors
// SPDX-License-Identifier: MPL-2.0
// Copyright (c) 2023 HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package e2etest

import (
	"context"
	"path/filepath"
	"testing"

	"github.com/hashicorp/go-version"

	"github.com/opentofu/tofu-exec/tfexec"
	"github.com/opentofu/tofu-exec/tfexec/internal/testutil"
)

func TestUpgrade012(t *testing.T) {
	runTestWithVersions(t, []string{testutil.Latest012}, "pre_011_syntax", func(t *testing.T, tfv *version.Version, tf *tfexec.Terraform) {
		err := tf.Init(context.Background())
		if err != nil {
			t.Fatalf("error running Init in test directory: %s", err)
		}

		err = tf.Upgrade012(context.Background())
		if err != nil {
			t.Fatalf("error from FormatWrite: %T %q", err, err)
		}

		for file, golden := range map[string]string{
			"file1.tf": "file1.golden.txt",
			"file2.tf": "file2.golden.txt",
		} {
			textFilesEqual(t, filepath.Join(tf.WorkingDir(), golden), filepath.Join(tf.WorkingDir(), file))
		}
	})
}
