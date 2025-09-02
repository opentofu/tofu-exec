// Copyright (c) The OpenTofu Authors
// SPDX-License-Identifier: MPL-2.0
// Copyright (c) 2023 HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package tfexec

import (
	"context"
	"os/exec"
	"strings"
)

type graphConfig struct {
	plan       string
	drawCycles bool
	graphType  string
}

var defaultGraphOptions = graphConfig{}

type GraphOption interface {
	configureGraph(*graphConfig)
}

func (opt *GraphPlanOption) configureGraph(conf *graphConfig) {
	conf.plan = opt.file
}

func (opt *DrawCyclesOption) configureGraph(conf *graphConfig) {
	conf.drawCycles = opt.drawCycles
}

func (opt *GraphTypeOption) configureGraph(conf *graphConfig) {
	conf.graphType = opt.graphType
}

func (tf *Tofu) Graph(ctx context.Context, opts ...GraphOption) (string, error) {
	graphCmd, err := tf.graphCmd(ctx, opts...)
	if err != nil {
		return "", err
	}
	var outBuf strings.Builder
	graphCmd.Stdout = &outBuf
	err = tf.runTofuCmd(ctx, graphCmd)
	if err != nil {
		return "", err
	}

	return outBuf.String(), nil

}

func (tf *Tofu) graphCmd(ctx context.Context, opts ...GraphOption) (*exec.Cmd, error) {
	c := defaultGraphOptions

	for _, o := range opts {
		o.configureGraph(&c)
	}

	args := []string{"graph"}

	if c.plan != "" {
		args = append(args, "-plan="+c.plan)
	}

	if c.drawCycles {
		args = append(args, "-draw-cycles")
	}

	if c.graphType != "" {
		args = append(args, "-type="+c.graphType)
	}

	return tf.buildTofuCmd(ctx, nil, args...), nil
}
