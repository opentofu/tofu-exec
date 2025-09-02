// Copyright (c) The OpenTofu Authors
// SPDX-License-Identifier: MPL-2.0
// Copyright (c) 2023 HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package tfexec

import (
	"context"
	"os/exec"

	tfjson "github.com/hashicorp/terraform-json"
)

// MetadataFunctions represents the tofu metadata functions -json subcommand.
func (tf *Tofu) MetadataFunctions(ctx context.Context) (*tfjson.MetadataFunctions, error) {
	functionsCmd := tf.metadataFunctionsCmd(ctx)

	var ret tfjson.MetadataFunctions
	err := tf.runTofuCmdJSON(ctx, functionsCmd, &ret)
	if err != nil {
		return nil, err
	}

	return &ret, nil
}

func (tf *Tofu) metadataFunctionsCmd(ctx context.Context, args ...string) *exec.Cmd {
	allArgs := []string{"metadata", "functions", "-json"}
	allArgs = append(allArgs, args...)

	return tf.buildTofuCmd(ctx, nil, allArgs...)
}
