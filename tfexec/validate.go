// Copyright (c) The OpenTofu Authors
// SPDX-License-Identifier: MPL-2.0
// Copyright (c) 2023 HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package tfexec

import (
	"bytes"
	"context"
	"encoding/json"
	tfjson "github.com/hashicorp/terraform-json"
)

// Validate represents the validate subcommand to the OpenTofu CLI.
func (tf *Tofu) Validate(ctx context.Context) (*tfjson.ValidateOutput, error) {
	cmd := tf.buildTofuCmd(ctx, nil, "validate", "-no-color", "-json")

	var outBuf = bytes.Buffer{}
	cmd.Stdout = &outBuf

	err := tf.runTofuCmd(ctx, cmd)
	// TODO: this command should not exit 1 if you pass -json as its hard to differentiate other errors
	if err != nil && cmd.ProcessState.ExitCode() != 1 {
		return nil, err
	}

	var ret tfjson.ValidateOutput
	// TODO: ret.UseJSONNumber(true) validate output should support JSON numbers
	jsonErr := json.Unmarshal(outBuf.Bytes(), &ret)
	if jsonErr != nil {
		// the original call was possibly bad, if it has an error, actually just return that
		if err != nil {
			return nil, err
		}

		return nil, jsonErr
	}

	return &ret, nil
}
