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

	// positional arguments
	args = append(args, lockID)

	// optional positional arguments
	if c.dir != "" {
		args = append(args, c.dir)
	}

	return tf.buildTofuCmd(ctx, nil, args...), nil
}
