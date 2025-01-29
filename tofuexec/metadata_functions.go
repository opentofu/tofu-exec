// Copyright (c) The OpenTofu Authors
// SPDX-License-Identifier: MPL-2.0
// Copyright (c) 2023 HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package tofuexec

import (
	"context"
	"fmt"
	"os/exec"

	tfjson "github.com/hashicorp/terraform-json"
)

// MetadataFunctions represents the tofu metadata functions -json subcommand.
func (tf *OpenTofu) MetadataFunctions(ctx context.Context) (*tfjson.MetadataFunctions, error) {
	err := tf.compatible(ctx, tf1_4_0, nil)
	if err != nil {
		return nil, fmt.Errorf("tofu metadata functions was added in 1.4.0: %w", err)
	}

	functionsCmd := tf.metadataFunctionsCmd(ctx)

	var ret tfjson.MetadataFunctions
	err = tf.runTofuCmdJSON(ctx, functionsCmd, &ret)
	if err != nil {
		return nil, err
	}

	return &ret, nil
}

func (tf *OpenTofu) metadataFunctionsCmd(ctx context.Context, args ...string) *exec.Cmd {
	allArgs := []string{"metadata", "functions", "-json"}
	allArgs = append(allArgs, args...)

	return tf.buildTofuCmd(ctx, nil, allArgs...)
}
