// Copyright (c) The OpenTofu Authors
// SPDX-License-Identifier: MPL-2.0
// Copyright (c) 2023 HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package testutil

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"runtime"
	"sync"
	"testing"

	"github.com/hashicorp/hc-install/build"
	"github.com/hashicorp/hc-install/product"
	"github.com/opentofu/tofudl"
)

const (
	Latest013   = "0.13.7"
	Latest014   = "0.14.11"
	Latest015   = "0.15.5"
	Latest_v1   = "1.9.1"
	Latest_v1_1 = "1.1.9"
	Latest_v1_5 = "1.5.3"
	Latest_v1_6 = "1.6.0-alpha20230719"
)

const appendUserAgent = "tfexec-testutil"

type TFCache struct {
	sync.Mutex

	dir   string
	execs map[string]string
}

func NewTFCache(dir string) *TFCache {
	return &TFCache{
		dir:   dir,
		execs: map[string]string{},
	}
}

func (tf *TFCache) GitRef(t *testing.T, ref string) string {
	t.Helper()

	key := "gitref:" + ref

	return tf.find(t, key, func(ctx context.Context) (string, error) {
		gr := &build.GitRevision{
			Product: product.Terraform,
			Ref:     ref,
		}
		gr.SetLogger(TestLogger())

		return gr.Build(ctx)
	})
}

func (tf *TFCache) Version(t *testing.T, v string) string {
	t.Helper()

	key := "tofu-v" + v // example: tofu-v1.9.1

	return tf.find(t, key, func(ctx context.Context) (string, error) {
		dl, err := tofudl.New()
		if err != nil {
			return "", fmt.Errorf("error when instantiating tofudl %w", err)
		}

		ver := tofudl.Version(v)
		opts := tofudl.DownloadOptVersion(ver)
		binary, err := dl.Download(ctx, opts)
		if err != nil {
			return "", fmt.Errorf("error when downloading %w", err)
		}

		// Write out the tofu binary to the disk:
		file := filepath.Join(tf.dir, key)
		if runtime.GOOS == "windows" {
			file += ".exe"
		}

		if err := os.WriteFile(file, binary, 0755); err != nil {
			return "", fmt.Errorf("error when writing the file %s: %w", file, err)
		}

		return file, nil
	})
}
