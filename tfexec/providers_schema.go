// Copyright (c) The OpenTofu Authors
// SPDX-License-Identifier: MPL-2.0
// Copyright (c) 2023 HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package tfexec

import (
	"context"
	"os/exec"

	"github.com/opentofu/tofu-exec/jsontypes"
)

// ProvidersSchema represents the tofu providers schema -json subcommand.
func (tf *Tofu) ProvidersSchema(ctx context.Context) (*jsontypes.ProviderSchemas, error) {
	schemaCmd := tf.providersSchemaCmd(ctx)

	var ret jsontypes.ProviderSchemas
	err := tf.runTofuCmdJSON(ctx, schemaCmd, &ret)
	if err != nil {
		return nil, err
	}

	err = ret.Validate()
	if err != nil {
		return nil, err
	}

	return &ret, nil
}

func (tf *Tofu) providersSchemaCmd(ctx context.Context, args ...string) *exec.Cmd {
	allArgs := []string{"providers", "schema", "-json", "-no-color"}
	allArgs = append(allArgs, args...)

	return tf.buildTofuCmd(ctx, nil, allArgs...)
}
