// Copyright (c) The OpenTofu Authors
// SPDX-License-Identifier: MPL-2.0
// Copyright (c) 2023 HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package tofuexec

import (
	"context"
	"fmt"
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
// The given io.Writer, if specified, will receive
// [machine-readable](https://opentofu.org/docs/internals/machine-readable-ui/)
// JSON from OpenTofu including test results.
func (tf *OpenTofu) Test(ctx context.Context, w io.Writer, opts ...TestOption) error {
	err := tf.compatible(ctx, tofu1_6_0, nil)

	if err != nil {
		return fmt.Errorf("terraform test was added in 1.6.0: %w", err)
	}

	tf.SetStdout(w)

	testCmd := tf.testCmd(ctx)

	err = tf.runTofuCmd(ctx, testCmd)

	if err != nil {
		return err
	}

	return nil
}

func (tf *OpenTofu) testCmd(ctx context.Context, opts ...TestOption) *exec.Cmd {
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
