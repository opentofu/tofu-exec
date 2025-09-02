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

// WorkspaceShow represents the workspace show subcommand to the OpenTofu CLI.
func (tf *Tofu) WorkspaceShow(ctx context.Context) (string, error) {
	workspaceShowCmd, err := tf.workspaceShowCmd(ctx)
	if err != nil {
		return "", err
	}

	var outBuffer strings.Builder
	workspaceShowCmd.Stdout = &outBuffer

	err = tf.runTofuCmd(ctx, workspaceShowCmd)
	if err != nil {
		return "", err
	}

	return strings.TrimSpace(outBuffer.String()), nil
}

func (tf *Tofu) workspaceShowCmd(ctx context.Context) (*exec.Cmd, error) {
	return tf.buildTofuCmd(ctx, nil, "workspace", "show", "-no-color"), nil
}
