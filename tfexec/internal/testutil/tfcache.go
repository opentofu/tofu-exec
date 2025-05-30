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

	"github.com/opentofu/tofudl"
)

const (
	Latest013   = "0.13.7"
	Latest014   = "0.14.11"
	Latest015   = "0.15.5"
	Latest_v1   = "1.9.1"
	Latest_v1_6 = "1.6.1"
	Latest_v1_7 = "1.7.8"
	Latest_v1_8 = "1.8.9"
	Latest_v1_9 = "1.9.1"
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
