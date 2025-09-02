// Copyright (c) The OpenTofu Authors
// SPDX-License-Identifier: MPL-2.0
// Copyright (c) 2023 HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package tfexec

import (
	"context"
	"fmt"
	"io"
	"os/exec"
	"strings"
)

type formatConfig struct {
	recursive bool
	dir       string
}

var defaultFormatConfig = formatConfig{
	recursive: false,
}

type FormatOption interface {
	configureFormat(*formatConfig)
}

func (opt *RecursiveOption) configureFormat(conf *formatConfig) {
	conf.recursive = opt.recursive
}

func (opt *DirOption) configureFormat(conf *formatConfig) {
	conf.dir = opt.path
}

// FormatString formats a passed string.
func (tf *Tofu) FormatString(ctx context.Context, content string) (string, error) {
	in := strings.NewReader(content)
	var outBuf strings.Builder
	err := tf.Format(ctx, in, &outBuf)
	if err != nil {
		return "", err
	}
	return outBuf.String(), nil
}

// Format performs formatting on the unformatted io.Reader (as stdin to the CLI) and returns
// the formatted result on the formatted io.Writer.
func (tf *Tofu) Format(ctx context.Context, unformatted io.Reader, formatted io.Writer) error {
	cmd, err := tf.formatCmd(ctx, nil, Dir("-"))
	if err != nil {
		return err
	}

	cmd.Stdin = unformatted
	cmd.Stdout = mergeWriters(cmd.Stdout, formatted)

	return tf.runTofuCmd(ctx, cmd)
}

// FormatWrite attempts to format and modify all config files in the working or selected (via DirOption) directory.
func (tf *Tofu) FormatWrite(ctx context.Context, opts ...FormatOption) error {
	for _, o := range opts {
		switch o := o.(type) {
		case *DirOption:
			if o.path == "-" {
				return fmt.Errorf("a path of \"-\" is not supported for this method, please use FormatString")
			}
		}
	}

	cmd, err := tf.formatCmd(ctx, []string{"-write=true", "-list=false", "-diff=false"}, opts...)
	if err != nil {
		return err
	}

	return tf.runTofuCmd(ctx, cmd)
}

// FormatCheck returns true if the config files in the working or selected (via DirOption) directory are already formatted.
func (tf *Tofu) FormatCheck(ctx context.Context, opts ...FormatOption) (bool, []string, error) {
	for _, o := range opts {
		switch o := o.(type) {
		case *DirOption:
			if o.path == "-" {
				return false, nil, fmt.Errorf("a path of \"-\" is not supported for this method, please use FormatString")
			}
		}
	}

	cmd, err := tf.formatCmd(ctx, []string{"-write=false", "-list=true", "-diff=false", "-check=true"}, opts...)
	if err != nil {
		return false, nil, err
	}

	var outBuf strings.Builder
	cmd.Stdout = mergeWriters(cmd.Stdout, &outBuf)

	err = tf.runTofuCmd(ctx, cmd)
	if err == nil {
		return true, nil, nil
	}
	if cmd.ProcessState.ExitCode() == 3 {
		// unformatted, parse the file list

		files := []string{}
		lines := strings.Split(strings.Replace(outBuf.String(), "\r\n", "\n", -1), "\n")
		for _, l := range lines {
			l = strings.TrimSpace(l)
			if l == "" {
				continue
			}
			files = append(files, l)
		}

		return false, files, nil
	}
	return false, nil, err
}

func (tf *Tofu) formatCmd(ctx context.Context, args []string, opts ...FormatOption) (*exec.Cmd, error) {
	c := defaultFormatConfig

	for _, o := range opts {
		o.configureFormat(&c)
	}

	args = append([]string{"fmt", "-no-color"}, args...)

	if c.recursive {
		args = append(args, "-recursive")
	}

	if c.dir != "" {
		args = append(args, c.dir)
	}

	return tf.buildTofuCmd(ctx, nil, args...), nil
}
