// Copyright (c) The OpenTofu Authors
// SPDX-License-Identifier: MPL-2.0
// Copyright (c) 2023 HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package tofuexec

import (
	"context"
	"runtime"
	"testing"

	"github.com/opentofu/tofu-exec/tofuexec/internal/testutil"
)

func TestFormatCmd(t *testing.T) {
	if runtime.GOOS == "darwin" && runtime.GOARCH == "arm64" {
		t.Skip("OpenTofu for darwin/arm64 is not available until v1")
	}

	td := t.TempDir()

	tf, err := NewOpenTofu(td, tfVersion(t, testutil.Latest_v1))
	if err != nil {
		t.Fatal(err)
	}

	// empty env, to avoid environ mismatch in testing
	tf.SetEnv(map[string]string{})

	t.Run("defaults", func(t *testing.T) {
		fmtCmd, err := tf.formatCmd(context.Background(), []string{})
		if err != nil {
			t.Fatal(err)
		}

		assertCmd(t, []string{
			"fmt",
			"-no-color",
		}, nil, fmtCmd)
	})

	t.Run("override all defaults", func(t *testing.T) {
		fmtCmd, err := tf.formatCmd(context.Background(),
			[]string{"string1", "string2"},
			Recursive(true),
			Dir("mydir"))
		if err != nil {
			t.Fatal(err)
		}

		assertCmd(t, []string{
			"fmt",
			"-no-color",
			"string1",
			"string2",
			"-recursive",
			"mydir",
		}, nil, fmtCmd)
	})
}
