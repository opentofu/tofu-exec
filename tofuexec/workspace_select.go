// Copyright (c) The OpenTofu Authors
// SPDX-License-Identifier: MPL-2.0
// Copyright (c) 2023 HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package tofuexec

import "context"

// WorkspaceSelect represents the workspace select subcommand to the OpenTofu CLI.
func (tf *OpenTofu) WorkspaceSelect(ctx context.Context, workspace string) error {
	// TODO: [DIR] param option

	return tf.runTofuCmd(ctx, tf.buildTofuCmd(ctx, nil, "workspace", "select", "-no-color", workspace))
}
