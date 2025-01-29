// Copyright (c) The OpenTofu Authors
// SPDX-License-Identifier: MPL-2.0
// Copyright (c) 2023 HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package e2etest

import (
	"context"
	"encoding/json"
	"errors"
	"strings"
	"testing"

	"github.com/hashicorp/go-version"
	tfjson "github.com/hashicorp/terraform-json"

	"github.com/opentofu/tofu-exec/tofuexec"
)

var (
	v1_0_1 = version.Must(version.NewVersion("1.0.1"))
)

func TestShow(t *testing.T) {
	runTest(t, "basic_with_state", func(t *testing.T, tfv *version.Version, tf *tofuexec.OpenTofu) {

		if tfv.LessThan(providerAddressMinVersion) {
			t.Skip("state file provider FQNs not compatible with this OpenTofu version")
		}

		providerName := "registry.opentofu.org/hashicorp/null"
		if tfv.LessThan(providerAddressMinVersion) {
			providerName = "null"
		}

		formatVersion := "1.0"
		sensitiveValues := json.RawMessage("{}")

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
	runTest(t, "empty", func(t *testing.T, tfv *version.Version, tf *tofuexec.OpenTofu) {
		if tfv.LessThan(showMinVersion) {
			t.Skip("tofu show was added in OpenTofu 0.12, so test is not valid")
		}

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
	// From v1.2.0 onwards, running show before init in the basic case returns
	// an empty state with no error.
	runTest(t, "basic", func(t *testing.T, tfv *version.Version, tf *tofuexec.OpenTofu) {
		// HACK KEM: Really I mean tfv.LessThan(version.Must(version.NewVersion("1.2.0"))),
		// but I want this test to run for refs/heads/main prior to the release of v1.2.0.
		if tfv.LessThan(version.Must(version.NewVersion("1.2.0"))) {

			t.Skip("test applies only to v1.2.0 and greater")
		}
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
	// From v1.2.0 onwards, running show before init in the basic case returns
	// an empty state with no error.
	runTest(t, "registry_module", func(t *testing.T, tfv *version.Version, tf *tofuexec.OpenTofu) {
		if tfv.LessThanOrEqual(version.Must(version.NewVersion("1.2.0"))) {
			t.Skip("test applies only to v1.2.0 and greater")
		}
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
	runTest(t, "inmem_backend", func(t *testing.T, tfv *version.Version, tf *tofuexec.OpenTofu) {
		if tfv.LessThan(showMinVersion) {
			t.Skip("tofu show was added in OpenTofu 0.12, so test is not valid")
		}

		_, err := tf.Show(context.Background())
		if err == nil {
			t.Fatalf("expected error, but did not get one")
		}
	})
}

func TestShow_noInitLocalBackendNonDefaultState(t *testing.T) {
	runTest(t, "local_backend_non_default_state", func(t *testing.T, tfv *version.Version, tf *tofuexec.OpenTofu) {
		if tfv.LessThan(showMinVersion) {
			t.Skip("tofu show was added in OpenTofu 0.12, so test is not valid")
		}

		_, err := tf.Show(context.Background())
		if err == nil {
			t.Fatalf("expected error, but did not get one")
		}
	})
}

func TestShow_noInitCloudBackend(t *testing.T) {
	runTest(t, "cloud_backend", func(t *testing.T, tfv *version.Version, tf *tofuexec.OpenTofu) {
		if tfv.LessThan(version.Must(version.NewVersion("1.1.0"))) {
			t.Skip("cloud backend was added in OpenTofu 1.1, so test is not valid")
		}

		_, err := tf.Show(context.Background())
		if err == nil {
			t.Fatalf("expected error, but did not get one")
		}
	})
}

func TestShow_noInitEtcdBackend(t *testing.T) {
	runTest(t, "etcd_backend", func(t *testing.T, tfv *version.Version, tf *tofuexec.OpenTofu) {
		if tfv.LessThan(showMinVersion) {
			t.Skip("tofu show was added in OpenTofu 0.12, so test is not valid")
		}

		if tfv.GreaterThanOrEqual(version.Must(version.NewVersion("1.3.0"))) || tfv.Prerelease() != "" {
			t.Skip("etcd backend was removed in OpenTofu 1.3, so test is not valid")
		}

		_, err := tf.Show(context.Background())
		if err == nil {
			t.Fatalf("expected error, but did not get one")
		}
	})
}

func TestShow_noInitRemoteBackend(t *testing.T) {
	runTest(t, "remote_backend", func(t *testing.T, tfv *version.Version, tf *tofuexec.OpenTofu) {
		if tfv.LessThan(showMinVersion) {
			t.Skip("tofu show was added in OpenTofu 0.12, so test is not valid")
		}

		_, err := tf.Show(context.Background())
		if err == nil {
			t.Fatalf("expected error, but did not get one")
		}
	})
}

func TestShow_statefileDoesNotExist(t *testing.T) {
	runTest(t, "basic", func(t *testing.T, tfv *version.Version, tf *tofuexec.OpenTofu) {
		if tfv.LessThan(showMinVersion) {
			t.Skip("tofu show was added in OpenTofu 0.12, so test is not valid")
		}

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

func TestShow_versionMismatch(t *testing.T) {
	runTest(t, "basic", func(t *testing.T, tfv *version.Version, tf *tofuexec.OpenTofu) {
		// only testing versions without show
		if tfv.GreaterThanOrEqual(showMinVersion) {
			t.Skip("tofu show was added in OpenTofu 0.12, so test is not valid")
		}

		var mismatch *tofuexec.ErrVersionMismatch
		_, err := tf.Show(context.Background())
		if !errors.As(err, &mismatch) {
			t.Fatalf("expected version mismatch error, got %T %s", err, err)
		}
		if mismatch.Actual != "0.11.15" {
			t.Fatalf("expected version 0.11.15, got %q", mismatch.Actual)
		}
		if mismatch.MinInclusive != "0.12.0" {
			t.Fatalf("expected min 0.12.0, got %q", mismatch.MinInclusive)
		}
		if mismatch.MaxExclusive != "-" {
			t.Fatalf("expected max -, got %q", mismatch.MaxExclusive)
		}
	})
}

func TestShowBigInt(t *testing.T) {
	runTest(t, "bigint", func(t *testing.T, tfv *version.Version, tf *tofuexec.OpenTofu) {
		if tfv.LessThan(showMinVersion) {
			t.Skip("tofu show was added in OpenTofu 0.12, so test is not valid")
		}

		providerName := "registry.opentofu.org/hashicorp/random"
		if tfv.LessThan(providerAddressMinVersion) {
			providerName = "random"
		}

		formatVersion := "1.0"
		sensitiveValues := json.RawMessage("{}")

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

// Since our plan strings are not large, prefer simple cross-platform
// normalization handling over pulling in a dependency.
func normalizePlanOutput(str string) string {
	// Ignore any extra newlines at the beginning or end of output
	str = strings.TrimSpace(str)
	// Normalize CRLF to LF for cross-platform testing
	str = strings.Replace(str, "\r\n", "\n", -1)

	return str
}
