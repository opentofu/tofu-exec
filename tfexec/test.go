// Copyright (c) The OpenTofu Authors
// SPDX-License-Identifier: MPL-2.0
// Copyright (c) 2023 HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package tfexec

import (
	"context"
	"io"
	"os/exec"
)

type testConfig struct {
	testsDirectory string
}

var defaultTestOptions = testConfig{}

type TestOption interface {
	configureTest(*testConfig)
}

func (opt *TestsDirectoryOption) configureTest(conf *testConfig) {
	conf.testsDirectory = opt.testsDirectory
}

// Test represents the tofu test -json subcommand.
//
// The given io.Writer, if specified, will receive machine-readable
// JSON from Terraform including test results.
func (tf *Tofu) Test(ctx context.Context, w io.Writer, opts ...TestOption) error {
	tf.SetStdout(w)

	testCmd := tf.testCmd(ctx)

	err := tf.runTofuCmd(ctx, testCmd)

	if err != nil {
		return err
	}

	return nil
}

func (tf *Tofu) testCmd(ctx context.Context, opts ...TestOption) *exec.Cmd {
	c := defaultTestOptions

	for _, o := range opts {
		o.configureTest(&c)
	}

	args := []string{"test", "-json"}

	if c.testsDirectory != "" {
		args = append(args, "-tests-directory="+c.testsDirectory)
	}

	return tf.buildTofuCmd(ctx, nil, args...)
}
