// Copyright (c) The OpenTofu Authors
// SPDX-License-Identifier: MPL-2.0
// Copyright (c) 2023 HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package e2etest

import (
	"github.com/google/go-cmp/cmp"
	"github.com/google/go-cmp/cmp/cmpopts"
	"github.com/opentofu/tofu-exec/jsontypes"
	"github.com/zclconf/go-cty-debug/ctydebug"
)

// comparison functions for tfjson structs used in tests

// diffState returns a human-readable report of the differences between two
// state values. It returns an empty string if the two values are equal.
func diffState(expected *jsontypes.State, actual *jsontypes.State) string {
	return cmp.Diff(expected, actual, cmpopts.IgnoreFields(jsontypes.State{}, "TerraformVersion"), cmpopts.IgnoreFields(jsontypes.State{}, "useJSONNumber"))
}

// diffPlan returns a human-readable report of the differences between two
// plan values. It returns an empty string if the two values are equal.
func diffPlan(expected *jsontypes.Plan, actual *jsontypes.Plan) string {
	return cmp.Diff(expected, actual, cmpopts.IgnoreFields(jsontypes.Plan{}, "TerraformVersion"))
}

// diffSchema returns a human-readable report of the differences between two
// schema values. It returns an empty string if the two values are equal.
func diffSchema(expected *jsontypes.ProviderSchemas, actual *jsontypes.ProviderSchemas) string {
	return cmp.Diff(expected, actual, ctydebug.CmpOptions)
}
