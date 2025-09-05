// Copyright (c) The OpenTofu Authors
// SPDX-License-Identifier: MPL-2.0
// Copyright (c) 2023 HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package e2etest

import (
	"os"
	"testing"

	"github.com/opentofu/tofu-exec/tfexec/internal/testutil"
)

var tfcache *testutil.TFCache

func TestMain(m *testing.M) {
	os.Exit(func() int {
		installDir, err := os.MkdirTemp("", "tfinstall")
		if err != nil {
			panic(err)
		}
		defer os.RemoveAll(installDir)

		tfcache = testutil.NewTFCache(installDir)
		return m.Run()
	}())
}
