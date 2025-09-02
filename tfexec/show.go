// Copyright (c) The OpenTofu Authors
// SPDX-License-Identifier: MPL-2.0
// Copyright (c) 2023 HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package tfexec

import (
	"context"
	"fmt"
	"os/exec"
	"strings"

	"github.com/hashicorp/go-version"
	tfjson "github.com/hashicorp/terraform-json"
)

type showConfig struct {
	reattachInfo ReattachInfo
}

var defaultShowOptions = showConfig{}

type ShowOption interface {
	configureShow(*showConfig)
}

func (opt *ReattachOption) configureShow(conf *showConfig) {
	conf.reattachInfo = opt.info
}

// Show reads the default state path and outputs the state.
// To read a state or plan file, ShowState or ShowPlan must be used instead.
func (tf *Tofu) Show(ctx context.Context, opts ...ShowOption) (*tfjson.State, error) {
	c := defaultShowOptions

	for _, o := range opts {
		o.configureShow(&c)
	}

	mergeEnv := map[string]string{}
	if c.reattachInfo != nil {
		reattachStr, err := c.reattachInfo.marshalString()
		if err != nil {
			return nil, err
		}
		mergeEnv[reattachEnvVar] = reattachStr
	}

	showCmd := tf.showCmd(ctx, true, mergeEnv)

	var ret tfjson.State
	ret.UseJSONNumber(true)
	err := tf.runTofuCmdJSON(ctx, showCmd, &ret)
	if err != nil {
		return nil, err
	}

	err = ret.Validate()
	if err != nil {
		return nil, err
	}

	return &ret, nil
}

// ShowStateFile reads a given state file and outputs the state.
func (tf *Tofu) ShowStateFile(ctx context.Context, statePath string, opts ...ShowOption) (*tfjson.State, error) {
	if statePath == "" {
		return nil, fmt.Errorf("statePath cannot be blank: use Show() if not passing statePath")
	}

	c := defaultShowOptions

	for _, o := range opts {
		o.configureShow(&c)
	}

	mergeEnv := map[string]string{}
	if c.reattachInfo != nil {
		reattachStr, err := c.reattachInfo.marshalString()
		if err != nil {
			return nil, err
		}
		mergeEnv[reattachEnvVar] = reattachStr
	}

	showCmd := tf.showCmd(ctx, true, mergeEnv, statePath)

	var ret tfjson.State
	ret.UseJSONNumber(true)
	err := tf.runTofuCmdJSON(ctx, showCmd, &ret)
	if err != nil {
		return nil, err
	}

	err = ret.Validate()
	if err != nil {
		return nil, err
	}

	return &ret, nil
}

// ShowPlanFile reads a given plan file and outputs the plan.
func (tf *Tofu) ShowPlanFile(ctx context.Context, planPath string, opts ...ShowOption) (*tfjson.Plan, error) {
	if planPath == "" {
		return nil, fmt.Errorf("planPath cannot be blank: use Show() if not passing planPath")
	}

	c := defaultShowOptions

	for _, o := range opts {
		o.configureShow(&c)
	}

	mergeEnv := map[string]string{}
	if c.reattachInfo != nil {
		reattachStr, err := c.reattachInfo.marshalString()
		if err != nil {
			return nil, err
		}
		mergeEnv[reattachEnvVar] = reattachStr
	}

	showCmd := tf.showCmd(ctx, true, mergeEnv, planPath)

	var ret tfjson.Plan
	err := tf.runTofuCmdJSON(ctx, showCmd, &ret)
	if err != nil {
		return nil, err
	}

	err = ret.Validate()
	if err != nil {
		return nil, err
	}

	return &ret, nil

}

// ShowPlanFileRaw reads a given plan file and outputs the plan in a
// human-friendly, opaque format.
func (tf *Tofu) ShowPlanFileRaw(ctx context.Context, planPath string, opts ...ShowOption) (string, error) {
	if planPath == "" {
		return "", fmt.Errorf("planPath cannot be blank: use Show() if not passing planPath")
	}

	c := defaultShowOptions

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

	showCmd := tf.showCmd(ctx, false, mergeEnv, planPath)

	var outBuf strings.Builder
	showCmd.Stdout = &outBuf
	err := tf.runTofuCmd(ctx, showCmd)
	if err != nil {
		return "", err
	}

	return outBuf.String(), nil

}

// ShowModule returns module config based on moduleDir located in local filesystem.
// This command was added in tofu version 1.11
func (tf *Tofu) ShowModule(ctx context.Context, moduleDir string) (*Module, error) {
	err := tf.compatible(ctx, version.Must(version.NewVersion("1.11.0-dev")), nil)
	if err != nil {
		return nil, fmt.Errorf("`tofu show -json -module=DIR` was added in tofu 1.11.0: %w", err)
	}
	if moduleDir == "" {
		return nil, fmt.Errorf("moduleDir cannot be blank")
	}

	showCmd := tf.showCmd(ctx, true, nil, "-module="+moduleDir)
	var ret ModuleRoot
	err = tf.runTofuCmdJSON(ctx, showCmd, &ret)
	if err != nil {
		return nil, err
	}

	return &ret.Module, nil
}
func (tf *Tofu) showCmd(ctx context.Context, jsonOutput bool, mergeEnv map[string]string, args ...string) *exec.Cmd {
	allArgs := []string{"show"}
	if mergeEnv == nil {
		mergeEnv = map[string]string{}
	}
	if jsonOutput {
		allArgs = append(allArgs, "-json")
	}
	allArgs = append(allArgs, "-no-color")
	allArgs = append(allArgs, args...)

	return tf.buildTofuCmd(ctx, mergeEnv, allArgs...)
}
