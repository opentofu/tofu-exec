// Copyright (c) The OpenTofu Authors
// SPDX-License-Identifier: MPL-2.0
// Copyright (c) 2023 HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package tfexec

import "context"

// WorkspaceSelect represents the workspace select subcommand to the Terraform CLI.
func (tf *Tofu) WorkspaceSelect(ctx context.Context, workspace string) error {
	// TODO: [DIR] param option

	return tf.runTerraformCmd(ctx, tf.buildTerraformCmd(ctx, nil, "workspace", "select", "-no-color", workspace))
}
