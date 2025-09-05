// Copyright (c) The OpenTofu Authors
// SPDX-License-Identifier: MPL-2.0
// Copyright (c) 2023 HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package testutil

import (
	"context"
	"fmt"
	"os"
	"os/exec"
	"strings"
	"testing"
)

func getInstalledTofuIfVersionsMatch(v string) (string, error) {
	// v will likely have a prefix tofu-
	v = strings.TrimPrefix(v, "tofu-")
	// Check a currently installed tofu version first
	tofuPath, err := exec.LookPath("tofu")
	if err != nil {
		return "", err
	}
	// We run tofu -v here to check if the version matches the desired one
	out, err := exec.CommandContext(context.Background(), tofuPath, "-v").Output()
	if err != nil {
		return "", err
	}

	if strings.Contains(string(out), v) {
		return tofuPath, nil
	}

	return "", fmt.Errorf("installed version: %s", string(out))
}

func (tf *TFCache) find(t *testing.T, key string, execPathFunc func(context.Context) (string, error)) string {
	t.Helper()

	if tf.dir == "" {
		// panic instead of t.fatal as this is going to affect all downstream tests reusing the cache entry
		panic("installDir not yet configured")
	}

	tf.Lock()
	defer tf.Unlock()

	if path, ok := tf.execs[key]; ok {
		return path
	}

	t.Logf("caching exec %q in dir %q", key, tf.dir)

	err := os.MkdirAll(tf.dir, 0777)
	if err != nil {
		// panic instead of t.fatal as this is going to affect all downstream tests reusing the cache entry
		panic(fmt.Sprintf("unable to mkdir %q: %s", tf.dir, err))
	}

	// Try to find it locally
	if path, err := getInstalledTofuIfVersionsMatch(key); err == nil {
		tf.execs[key] = path
		return path
	}

	ctx, cancelFunc := context.WithCancel(context.Background())
	t.Cleanup(cancelFunc)

	execPath, err := execPathFunc(ctx)
	if err != nil {
		// panic instead of t.fatal as this is going to affect all downstream tests reusing the cache entry
		panic(fmt.Sprintf("error installing tofu %q: %s", key, err))
	}

	tf.execs[key] = execPath

	return execPath
}
