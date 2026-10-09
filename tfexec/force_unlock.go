// Copyright (c) The OpenTofu Authors
// SPDX-License-Identifier: MPL-2.0
// Copyright (c) 2023 HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package tfexec

import (
	"context"
	"os/exec"
)

type forceUnlockConfig struct {
	dir string
}

var defaultForceUnlockOptions = forceUnlockConfig{}

type ForceUnlockOption interface {
	configureForceUnlock(*forceUnlockConfig)
}

func (opt *DirOption) configureForceUnlock(conf *forceUnlockConfig) {
	conf.dir = opt.path
}

// ForceUnlock represents the `tofu force-unlock` command
func (tf *Tofu) ForceUnlock(ctx context.Context, lockID string, opts ...ForceUnlockOption) error {
	unlockCmd, err := tf.forceUnlockCmd(ctx, lockID, opts...)
	if err != nil {
		return err
	}

	if err := tf.runTofuCmd(ctx, unlockCmd); err != nil {
		return err
	}

	return nil
}

func (tf *Tofu) forceUnlockCmd(ctx context.Context, lockID string, opts ...ForceUnlockOption) (*exec.Cmd, error) {
	c := defaultForceUnlockOptions

	for _, o := range opts {
		o.configureForceUnlock(&c)
	}
	args := []string{"force-unlock", "-no-color", "-force"}

	// optional positional arguments
	var positional []string
	if c.dir != "" {
		positional = append(positional, c.dir)
	}

	// positional arguments
	positional = append(positional, lockID)

	return tf.buildTofuCmd(ctx, nil, positional, args...), nil
}
