// Copyright (c) The OpenTofu Authors
// SPDX-License-Identifier: MPL-2.0
// Copyright (c) 2023 HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package tofuexec

import (
	"bytes"
	"context"
	"os/exec"
)

type statePullConfig struct {
	reattachInfo ReattachInfo
}

var defaultStatePullConfig = statePullConfig{}

type StatePullOption interface {
	configureShow(*statePullConfig)
}

func (opt *ReattachOption) configureStatePull(conf *statePullConfig) {
	conf.reattachInfo = opt.info
}

func (tf *OpenTofu) StatePull(ctx context.Context, opts ...StatePullOption) (string, error) {
	c := defaultStatePullConfig

	for _, o := range opts {
		o.configureShow(&c)
	}

	mergeEnv := map[string]string{}
	if c.reattachInfo != nil {
		reattachStr, err := c.reattachInfo.marshalString()
		if err != nil {
			return "", err
		}
		mergeEnv[reattachEnvVar] = reattachStr
	}

	cmd := tf.statePullCmd(ctx, mergeEnv)

	var ret bytes.Buffer
	cmd.Stdout = &ret
	err := tf.runTofuCmd(ctx, cmd)
	if err != nil {
		return "", err
	}

	return ret.String(), nil
}

func (tf *OpenTofu) statePullCmd(ctx context.Context, mergeEnv map[string]string) *exec.Cmd {
	args := []string{"state", "pull"}

	return tf.buildTofuCmd(ctx, mergeEnv, args...)
}
