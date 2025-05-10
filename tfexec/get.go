// Copyright (c) The OpenTofu Authors
// SPDX-License-Identifier: MPL-2.0
// Copyright (c) 2023 HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package tfexec

import (
	"context"
	"fmt"
	"os/exec"
)

type getCmdConfig struct {
	dir    string
	update bool
}

// GetCmdOption represents options used in the Get method.
type GetCmdOption interface {
	configureGet(*getCmdConfig)
}

func (opt *DirOption) configureGet(conf *getCmdConfig) {
	conf.dir = opt.path
}

func (opt *UpdateOption) configureGet(conf *getCmdConfig) {
	conf.update = opt.update
}

// Get represents the terraform get subcommand.
func (tf *Tofu) Get(ctx context.Context, opts ...GetCmdOption) error {
	cmd, err := tf.getCmd(ctx, opts...)
	if err != nil {
		return err
	}
	return tf.runTofuCmd(ctx, cmd)
}

func (tf *Tofu) getCmd(ctx context.Context, opts ...GetCmdOption) (*exec.Cmd, error) {
	c := getCmdConfig{}

	for _, o := range opts {
		o.configureGet(&c)
	}

	args := []string{"get", "-no-color"}

	args = append(args, "-update="+fmt.Sprint(c.update))

	if c.dir != "" {
		args = append(args, c.dir)
	}

	return tf.buildTofuCmd(ctx, nil, args...), nil
}
