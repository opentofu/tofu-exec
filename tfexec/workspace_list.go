// Copyright (c) The OpenTofu Authors
// SPDX-License-Identifier: MPL-2.0
// Copyright (c) 2023 HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package tfexec

import (
	"context"
	"strings"
)

// WorkspaceList represents the workspace list subcommand to the Terraform CLI.
func (tf *Tofu) WorkspaceList(ctx context.Context) ([]string, string, error) {
	// TODO: [DIR] param option
	wlCmd := tf.buildTofuCmd(ctx, nil, "workspace", "list", "-no-color")

	var outBuf strings.Builder
	wlCmd.Stdout = &outBuf

	err := tf.runTofuCmd(ctx, wlCmd)
	if err != nil {
		return nil, "", err
	}

	ws, current := parseWorkspaceList(outBuf.String())

	return ws, current, nil
}

const currentWorkspacePrefix = "* "

func parseWorkspaceList(stdout string) ([]string, string) {
	lines := strings.Split(stdout, "\n")

	current := ""
	workspaces := []string{}
	for _, line := range lines {
		line = strings.TrimSpace(line)
		if line == "" {
			continue
		}
		if strings.HasPrefix(line, currentWorkspacePrefix) {
			line = strings.TrimPrefix(line, currentWorkspacePrefix)
			current = line
		}
		workspaces = append(workspaces, line)
	}

	return workspaces, current
}
