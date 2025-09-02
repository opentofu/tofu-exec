// Copyright (c) The OpenTofu Authors
// SPDX-License-Identifier: MPL-2.0
// Copyright (c) 2023 HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package tfexec

import (
	"context"
	"os/exec"
	"strconv"
)

type workspaceDeleteConfig struct {
	lock        bool
	lockTimeout string
	force       bool
}

var defaultWorkspaceDeleteOptions = workspaceDeleteConfig{
	lock:        true,
	lockTimeout: "0s",
}

// WorkspaceDeleteCmdOption represents options that are applicable to the WorkspaceDelete method.
type WorkspaceDeleteCmdOption interface {
	configureWorkspaceDelete(*workspaceDeleteConfig)
}

func (opt *LockOption) configureWorkspaceDelete(conf *workspaceDeleteConfig) {
	conf.lock = opt.lock
}

func (opt *LockTimeoutOption) configureWorkspaceDelete(conf *workspaceDeleteConfig) {
	conf.lockTimeout = opt.timeout
}

func (opt *ForceOption) configureWorkspaceDelete(conf *workspaceDeleteConfig) {
	conf.force = opt.force
}

// WorkspaceDelete represents the workspace delete subcommand to the OpenTofu CLI.
func (tf *Tofu) WorkspaceDelete(ctx context.Context, workspace string, opts ...WorkspaceDeleteCmdOption) error {
	cmd, err := tf.workspaceDeleteCmd(ctx, workspace, opts...)
	if err != nil {
		return err
	}
	return tf.runTofuCmd(ctx, cmd)
}

func (tf *Tofu) workspaceDeleteCmd(ctx context.Context, workspace string, opts ...WorkspaceDeleteCmdOption) (*exec.Cmd, error) {
	c := defaultWorkspaceDeleteOptions

	for _, o := range opts {
		o.configureWorkspaceDelete(&c)
	}

	args := []string{"workspace", "delete", "-no-color"}

	if c.force {
		args = append(args, "-force")
	}
	if c.lockTimeout != "" && c.lockTimeout != defaultWorkspaceDeleteOptions.lockTimeout {
		args = append(args, "-lock-timeout="+c.lockTimeout)
	}
	if !c.lock {
		args = append(args, "-lock="+strconv.FormatBool(c.lock))
	}

	args = append(args, workspace)

	cmd := tf.buildTofuCmd(ctx, nil, args...)

	return cmd, nil
}
