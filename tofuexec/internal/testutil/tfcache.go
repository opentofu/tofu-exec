// Copyright (c) The OpenTofu Authors
// SPDX-License-Identifier: MPL-2.0
// Copyright (c) 2023 HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package testutil

import (
	"context"
	"fmt"
	"github.com/opentofu/tofudl"
	"os"
	"path"
	"runtime"
	"sync"
	"testing"
)

const (
	Latest_v1   = "1.9.0"
	Latest_v1_7 = "1.7.0"
)

type TofuCache struct {
	sync.Mutex

	dir        string
	execs      map[string]string
	downloader tofudl.Downloader
}

func NewTofuCache(dir string) *TofuCache {
	downloader, err := tofudl.New()
	if err != nil {
		panic(err)
	}
	return &TofuCache{
		downloader: downloader,
		dir:        dir,
		execs:      map[string]string{},
	}
}

func (tf *TofuCache) Version(t *testing.T, v string) string {
	t.Helper()

	key := "v:" + v

	return tf.find(t, key, func(ctx context.Context) (string, error) {
		data, err := tf.downloader.Download(ctx, tofudl.DownloadOptVersion(tofudl.Version(v)))
		if err != nil {
			return "", err
		}
		binName := "tofu"
		if runtime.GOOS == "windows" {
			binName += ".exe"
		}
		binPath := path.Join(tf.dir, binName)
		if err := os.WriteFile(binPath, data, 0777); err != nil {
			return "", fmt.Errorf("failed to write %s (%w)", binPath, err)
		}

		return binPath, nil
	})
}
