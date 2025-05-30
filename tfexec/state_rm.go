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

type stateRmConfig struct {
	backup      string
	backupOut   string
	dryRun      bool
	lock        bool
	lockTimeout string
	state       string
	stateOut    string
}

var defaultStateRmOptions = stateRmConfig{
	lock:        true,
	lockTimeout: "0s",
}

// StateRmCmdOption represents options used in the Refresh method.
type StateRmCmdOption interface {
	configureStateRm(*stateRmConfig)
}

func (opt *BackupOption) configureStateRm(conf *stateRmConfig) {
	conf.backup = opt.path
}

func (opt *BackupOutOption) configureStateRm(conf *stateRmConfig) {
	conf.backupOut = opt.path
}

func (opt *DryRunOption) configureStateRm(conf *stateRmConfig) {
	conf.dryRun = opt.dryRun
}

func (opt *LockOption) configureStateRm(conf *stateRmConfig) {
	conf.lock = opt.lock
}

func (opt *LockTimeoutOption) configureStateRm(conf *stateRmConfig) {
	conf.lockTimeout = opt.timeout
}

func (opt *StateOption) configureStateRm(conf *stateRmConfig) {
	conf.state = opt.path
}

func (opt *StateOutOption) configureStateRm(conf *stateRmConfig) {
	conf.stateOut = opt.path
}

// StateRm represents the terraform state rm subcommand.
func (tf *Tofu) StateRm(ctx context.Context, address string, opts ...StateRmCmdOption) error {
	cmd, err := tf.stateRmCmd(ctx, address, opts...)
	if err != nil {
		return err
	}
	return tf.runTofuCmd(ctx, cmd)
}

func (tf *Tofu) stateRmCmd(ctx context.Context, address string, opts ...StateRmCmdOption) (*exec.Cmd, error) {
	c := defaultStateRmOptions

	for _, o := range opts {
		o.configureStateRm(&c)
	}

	args := []string{"state", "rm", "-no-color"}

	// string opts: only pass if set
	if c.backup != "" {
		args = append(args, "-backup="+c.backup)
	}
	if c.backupOut != "" {
		args = append(args, "-backup-out="+c.backupOut)
	}
	if c.lockTimeout != "" {
		args = append(args, "-lock-timeout="+c.lockTimeout)
	}
	if c.state != "" {
		args = append(args, "-state="+c.state)
	}
	if c.stateOut != "" {
		args = append(args, "-state-out="+c.stateOut)
	}

	// boolean and numerical opts: always pass
	args = append(args, "-lock="+strconv.FormatBool(c.lock))

	// unary flags: pass if true
	if c.dryRun {
		args = append(args, "-dry-run")
	}

	// positional arguments
	args = append(args, address)

	return tf.buildTofuCmd(ctx, nil, args...), nil
}
