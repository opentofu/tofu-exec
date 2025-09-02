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

func TestParsePlaintextVersionOutput(t *testing.T) {
	for i, c := range []struct {
		expectedV         *version.Version
		expectedProviders map[string]*version.Version

		stdout string
	}{
		// 0.13 tests
		{
			mustVersion(t, "0.13.0-dev"), nil, `
OpenTofu v0.13.0-dev`,
		},
		{
			mustVersion(t, "0.13.0-dev"), map[string]*version.Version{
				"registry.opentofu.org/hashicorp/null": mustVersion(t, "2.1.2"),
				"registry.opentofu.org/paultyng/null":  mustVersion(t, "0.1.0"),
			}, `
OpenTofu v0.13.0-dev
+ provider registry.opentofu.org/hashicorp/null v2.1.2
+ provider registry.opentofu.org/paultyng/null v0.1.0`,
		},
		{
			mustVersion(t, "0.13.0-dev"), nil, `
OpenTofu v0.13.0-dev

Your version of OpenTofu is out of date! The latest version
is 0.13.1. You can update by downloading from https://www.terraform.io/downloads.html`,
		},
		{
			mustVersion(t, "0.13.0-dev"), map[string]*version.Version{
				"registry.opentofu.org/hashicorp/null": mustVersion(t, "2.1.2"),
				"registry.opentofu.org/paultyng/null":  mustVersion(t, "0.1.0"),
			}, `
OpenTofu v0.13.0-dev
+ provider registry.opentofu.org/hashicorp/null v2.1.2
+ provider registry.opentofu.org/paultyng/null v0.1.0

Your version of OpenTofu is out of date! The latest version
is 0.13.1. You can update by downloading from https://www.terraform.io/downloads.html`,
		},

		// 0.12 tests
		{
			mustVersion(t, "0.12.26"), nil, `
OpenTofu v0.12.26
`,
		},
		{
			mustVersion(t, "0.12.26"), map[string]*version.Version{
				"null": mustVersion(t, "2.1.2"),
			}, `
OpenTofu v0.12.26
+ provider.null v2.1.2
`,
		},
		{
			mustVersion(t, "0.12.18"), nil, `
OpenTofu v0.12.18

Your version of OpenTofu is out of date! The latest version
is 0.12.26. You can update by downloading from https://www.terraform.io/downloads.html
`,
		},
		{
			mustVersion(t, "0.12.18"), map[string]*version.Version{
				"null": mustVersion(t, "2.1.2"),
			}, `
OpenTofu v0.12.18
+ provider.null v2.1.2

Your version of OpenTofu is out of date! The latest version
is 0.12.26. You can update by downloading from https://www.terraform.io/downloads.html
`,
		},
	} {
		t.Run(fmt.Sprintf("%d %s", i, c.expectedV), func(t *testing.T) {
			actualV, actualProv, err := parsePlaintextVersionOutput(c.stdout)
			if err != nil {
				t.Fatal(err)
			}

			if !c.expectedV.Equal(actualV) {
				t.Fatalf("expected %s, got %s", c.expectedV, actualV)
			}

			for k, v := range c.expectedProviders {
				if actual := actualProv[k]; actual == nil || !v.Equal(actual) {
					t.Fatalf("expected %s for %s, got %s", v, k, actual)
				}
			}

			if len(c.expectedProviders) != len(actualProv) {
				t.Fatalf("expected %d providers, got %d", len(c.expectedProviders), len(actualProv))
			}
		})
	}
}

func TestParseJsonVersionOutput(t *testing.T) {
	testStdout := []byte(`{
  "terraform_version": "0.15.0-beta1",
  "platform": "darwin_amd64",
  "provider_selections": {
    "registry.opentofu.org/hashicorp/aws": "3.31.0",
    "registry.opentofu.org/hashicorp/google": "3.58.0"
  },
  "terraform_outdated": false
}
`)
	tfVersion, pvs, err := parseJsonVersionOutput(testStdout)
	if err != nil {
		t.Fatal(err)
	}
	expectedTfVer := mustVersion(t, "0.15.0-beta1")

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
		{true, "", "0.12.26", ""},
		{true, "", "0.13.0-beta3", ""},

		{false, "", "0.12.26", "0.12.25"},
		{false, "", "0.12.26", "0.12.26"},
		{false, "0.12.27", "0.12.26", ""},
		{true, "", "0.12.26", "0.13.0"},
		{true, "0.12.25", "0.12.26", ""},
		{true, "0.12.26", "0.12.26", ""},
		{true, "0.12.26", "0.12.26", "0.12.27"},
		{true, "0.12.26", "0.12.26", "0.13.0"},

		{false, "0.12.26", "0.13.0-beta3", "0.13.0"},
		{true, "0.12.26", "0.13.0-beta3", ""},
		{true, "0.13.0", "0.13.0-beta3", ""},
		{true, "0.13.0", "0.13.0-beta3", "0.14.0"},
		{true, "", "0.13.0-beta3", "0.14.0"},
		{expected: true, min: "1.11.0-dev", tfv: "1.11.0", max: ""},
		{expected: false, min: "1.11.0-dev", tfv: "1.10.0", max: ""},
	} {
		t.Run(fmt.Sprintf("%d", i), func(t *testing.T) {
			tfv, err := version.NewVersion(c.tfv)
			if err != nil {
				t.Fatal(err)
			}

			var minV *version.Version
			if c.min != "" {
				minV, err = version.NewVersion(c.min)
				if err != nil {
					t.Fatal(err)
				}
			}

			var maxV *version.Version
			if c.max != "" {
				maxV, err = version.NewVersion(c.max)
				if err != nil {
					t.Fatal(err)
				}
			}

			actual := versionInRange(tfv, minV, maxV)
			if actual != c.expected {
				t.Fatalf("expected %v, got %v: %s <= %s < %s", c.expected, actual, minV, tfv, maxV)
			}
		})
	}
}
