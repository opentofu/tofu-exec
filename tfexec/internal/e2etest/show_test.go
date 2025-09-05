// Copyright (c) The OpenTofu Authors
// SPDX-License-Identifier: MPL-2.0
// Copyright (c) 2023 HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package e2etest

import (
	"context"
	"encoding/json"
	"testing"

	"github.com/hashicorp/go-version"
	tfjson "github.com/hashicorp/terraform-json"

	"github.com/opentofu/tofu-exec/tfexec"
)

func TestShow(t *testing.T) {
	runTest(t, "basic_with_state", func(t *testing.T, tfv *version.Version, tf *tfexec.Tofu) {
		providerName := "registry.opentofu.org/hashicorp/null"
		var sensitiveValues json.RawMessage = []byte("{}")
		formatVersion := "1.0"

		expected := &tfjson.State{
			FormatVersion: formatVersion,
			// TerraformVersion is ignored to facilitate latest version testing
			Values: &tfjson.StateValues{
				RootModule: &tfjson.StateModule{
					Resources: []*tfjson.StateResource{{
						Address: "null_resource.foo",
						AttributeValues: map[string]interface{}{
							"id":       "5510719323588825107",
							"triggers": nil,
						},
						SensitiveValues: sensitiveValues,
						Mode:            tfjson.ManagedResourceMode,
						Type:            "null_resource",
						Name:            "foo",
						ProviderName:    providerName,
					}},
				},
			},
		}

		err := tf.Init(context.Background())
		if err != nil {
			t.Fatalf("error running Init in test directory: %s", err)
		}

		actual, err := tf.Show(context.Background())
		if err != nil {
			t.Fatal(err)
		}

		if diff := diffState(expected, actual); diff != "" {
			t.Fatalf("mismatch (-want +got):\n%s", diff)
		}
	})
}

func TestShow_emptyDir(t *testing.T) {
	runTest(t, "empty", func(t *testing.T, tfv *version.Version, tf *tfexec.Tofu) {
		formatVersion := "1.0"

		expected := &tfjson.State{
			FormatVersion: formatVersion,
		}

		actual, err := tf.Show(context.Background())
		if err != nil {
			t.Fatal(err)
		}

		if diff := diffState(expected, actual); diff != "" {
			t.Fatalf("mismatch (-want +got):\n%s", diff)
		}
	})
}

func TestShow_noInitBasic(t *testing.T) {
	t.Parallel()
	// From v1.2.0 onwards, running show before init in the basic case returns
	// an empty state with no error.
	runTest(t, "basic", func(t *testing.T, tfv *version.Version, tf *tfexec.Tofu) {
		expected := &tfjson.State{
			FormatVersion: "1.0",
		}

		actual, err := tf.Show(context.Background())
		if err != nil {
			t.Fatal(err)
		}

		if diff := diffState(expected, actual); diff != "" {
			t.Fatalf("mismatch (-want +got):\n%s", diff)
		}
	})
}

func TestShow_noInitModule(t *testing.T) {
	t.Parallel()

	runTest(t, "registry_module", func(t *testing.T, tfv *version.Version, tf *tfexec.Tofu) {
		expected := &tfjson.State{
			FormatVersion: "1.0",
		}

		actual, err := tf.Show(context.Background())
		if err != nil {
			t.Fatal(err)
		}

		if diff := diffState(expected, actual); diff != "" {
			t.Fatalf("mismatch (-want +got):\n%s", diff)
		}
	})
}

func TestShow_noInitInmemBackend(t *testing.T) {
	runTest(t, "inmem_backend", func(t *testing.T, tfv *version.Version, tf *tfexec.Tofu) {
		_, err := tf.Show(context.Background())
		if err == nil {
			t.Fatalf("expected error, but did not get one")
		}
	})
}

func TestShow_noInitLocalBackendNonDefaultState(t *testing.T) {
	runTest(t, "local_backend_non_default_state", func(t *testing.T, tfv *version.Version, tf *tfexec.Tofu) {
		_, err := tf.Show(context.Background())
		if err == nil {
			t.Fatalf("expected error, but did not get one")
		}
	})
}

func TestShow_noInitCloudBackend(t *testing.T) {
	runTest(t, "cloud_backend", func(t *testing.T, tfv *version.Version, tf *tfexec.Tofu) {
		_, err := tf.Show(context.Background())
		if err == nil {
			t.Fatalf("expected error, but did not get one")
		}
	})
}

func TestShow_noInitRemoteBackend(t *testing.T) {
	runTest(t, "remote_backend", func(t *testing.T, tfv *version.Version, tf *tfexec.Tofu) {
		_, err := tf.Show(context.Background())
		if err == nil {
			t.Fatalf("expected error, but did not get one")
		}
	})
}

func TestShow_statefileDoesNotExist(t *testing.T) {
	runTest(t, "basic", func(t *testing.T, tfv *version.Version, tf *tfexec.Tofu) {
		err := tf.Init(context.Background())
		if err != nil {
			t.Fatalf("error running Init in test directory: %s", err)
		}

		_, err = tf.ShowStateFile(context.Background(), "statefilefoo")
		if err == nil {
			t.Fatalf("expected error, but did not get one")
		}
	})
}

func TestShowBigInt(t *testing.T) {
	runTest(t, "bigint", func(t *testing.T, tfv *version.Version, tf *tfexec.Tofu) {
		providerName := "registry.opentofu.org/hashicorp/random"
		var sensitiveValues json.RawMessage = []byte("{}")
		formatVersion := "1.0"

		expected := &tfjson.State{
			FormatVersion: formatVersion,
			// TerraformVersion is ignored to facilitate latest version testing
			Values: &tfjson.StateValues{
				RootModule: &tfjson.StateModule{
					Resources: []*tfjson.StateResource{{
						Address: "random_integer.bigint",
						AttributeValues: map[string]interface{}{
							"id":      "7227701560655103598",
							"max":     json.Number("7227701560655103598"),
							"min":     json.Number("7227701560655103597"),
							"result":  json.Number("7227701560655103598"),
							"seed":    "12345",
							"keepers": nil,
						},
						SensitiveValues: sensitiveValues,
						Mode:            tfjson.ManagedResourceMode,
						Type:            "random_integer",
						Name:            "bigint",
						ProviderName:    providerName,
					}},
				},
			},
		}

		err := tf.Init(context.Background())
		if err != nil {
			t.Fatalf("error running Init in test directory: %s", err)
		}

		err = tf.Apply(context.Background())
		if err != nil {
			t.Fatalf("error running Apply in test directory: %s", err)
		}

		actual, err := tf.Show(context.Background())
		if err != nil {
			t.Fatal(err)
		}

		if diff := diffState(expected, actual); diff != "" {
			t.Fatalf("mismatch (-want +got):\n%s", diff)
		}
	})
}

// E2E test for ShowModule using deep_module testdata
func TestShowModule_ShowModule(t *testing.T) {
	runTest(t, "show_module", func(t *testing.T, tfv *version.Version, tf *tfexec.Tofu) {
		// Only run for OpenTofu version >= 1.11
		minVer := version.Must(version.NewVersion("1.11.0-dev"))
		if tfv.LessThan(minVer) {
			t.Skip("ShowModule requires OpenTofu version >= 1.11.0")
		}

		// ShowModule for submodule1
		mod, err := tf.ShowModule(context.Background(), "submodule1")
		if err != nil {
			t.Fatalf("ShowModule failed: %v", err)
		}
		if mod == nil {
			t.Fatalf("ShowModule returned nil module")
		}
		if len(mod.Variables) != 2 {
			t.Fatalf("expected 2 input variables, got %d", len(mod.Variables))
		}
		if len(mod.Outputs) != 1 {
			t.Fatalf("expected 1 output, got %d", len(mod.Outputs))
		}

		// Submodule2
		mod, err = tf.ShowModule(context.Background(), "submodule2")
		if err != nil {
			t.Fatalf("ShowModule failed: %v", err)
		}
		if mod == nil {
			t.Fatalf("ShowModule returned nil module")
		}
		if len(mod.Variables) != 1 {
			t.Fatalf("expected 1 input variables, got %d", len(mod.Variables))
		}

		if len(mod.ModuleCalls) != 1 {
			t.Fatalf("expected 1 nested submodule, got %d", len(mod.ModuleCalls))
		}

		// Submodule21 with 2 outputs and 2 inputs
		mod, err = tf.ShowModule(context.Background(), "submodule2/submodule21")
		if err != nil {
			t.Fatalf("ShowModule failed: %v", err)
		}
		if mod == nil {
			t.Fatalf("ShowModule returned nil module")
		}
		if len(mod.Variables) != 2 {
			t.Fatalf("expected 2 input variables, got %d", len(mod.Variables))
		}
		if len(mod.Outputs) != 2 {
			t.Fatalf("expected 1 nested submodule, got %d", len(mod.ModuleCalls))
		}
		if len(mod.ModuleCalls) != 0 {
			t.Fatalf("expected 0 nested submodule, got %d", len(mod.ModuleCalls))
		}

	})
}
