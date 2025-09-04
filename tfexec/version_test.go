// Copyright (c) The OpenTofu Authors
// SPDX-License-Identifier: MPL-2.0
// Copyright (c) 2023 HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package tfexec

import (
	"fmt"
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/hashicorp/go-version"
)

func mustVersion(t *testing.T, s string) *version.Version {
	v, err := version.NewVersion(s)
	if err != nil {
		t.Fatal(err)
	}
	return v
}

func TestParseJsonVersionOutput(t *testing.T) {
	testStdout := []byte(`{
  "terraform_version": "1.9.2",
  "platform": "darwin_amd64",
  "provider_selections": {
    "registry.opentofu.org/hashicorp/aws": "3.31.0",
    "registry.opentofu.org/hashicorp/google": "3.58.0"
  }
}
`)
	tfVersion, pvs, err := parseJsonVersionOutput(testStdout)
	if err != nil {
		t.Fatal(err)
	}
	expectedTfVer := mustVersion(t, "1.9.2")

	if !expectedTfVer.Equal(tfVersion) {
		t.Fatalf("version doesn't match (%q != %q)",
			expectedTfVer.String(), tfVersion.String())
	}

	expectedPvs := map[string]*version.Version{
		"registry.opentofu.org/hashicorp/aws":    mustVersion(t, "3.31.0"),
		"registry.opentofu.org/hashicorp/google": mustVersion(t, "3.58.0"),
	}
	if diff := cmp.Diff(expectedPvs, pvs); diff != "" {
		t.Fatalf("provider versions don't match: %s", diff)
	}
}

func TestVersionInRange(t *testing.T) {
	for i, c := range []struct {
		expected bool
		min      string
		tfv      string
		max      string
	}{
		{true, "", "1.6.2", ""},
		{true, "", "1.7.0-beta3", ""},

		{false, "", "1.6.1", "1.5.3"},
		{false, "", "1.6.1", "1.6.1"},
		{false, "1.7.0", "1.6.0", ""},
		{true, "", "1.6.0", "1.7.0"},
		{true, "1.5.0", "1.6.0", ""},
		{true, "1.6.1", "1.6.1", ""},
		{true, "1.6.0", "1.6.0", "1.7.0"},
		{true, "1.6.0", "1.6.0", "1.8.0"},

		{false, "1.6.0", "1.7.0-beta3", "1.7.0"},
		{true, "1.6.0", "1.7.0-beta3", ""},
		{true, "1.7.0", "1.7.0-beta3", ""},
		{true, "1.7.0", "1.7.0-beta3", "1.8.0"},
		{true, "", "1.7.0-beta3", "1.8.0"},
		{expected: true, min: "1.11.0-dev", tfv: "1.11.0", max: ""},
		{expected: false, min: "1.11.0-dev", tfv: "1.10.0", max: ""},
	} {
		t.Run(fmt.Sprintf("%d", i), func(t *testing.T) {
			tfv, err := version.NewVersion(c.tfv)
			if err != nil {
				t.Fatal(err)
			}

			var min *version.Version
			if c.min != "" {
				min, err = version.NewVersion(c.min)
				if err != nil {
					t.Fatal(err)
				}
			}

			var max *version.Version
			if c.max != "" {
				max, err = version.NewVersion(c.max)
				if err != nil {
					t.Fatal(err)
				}
			}

			actual := versionInRange(tfv, min, max)
			if actual != c.expected {
				t.Fatalf("expected %v, got %v: %s <= %s < %s", c.expected, actual, min, tfv, max)
			}
		})
	}
}
