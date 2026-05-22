// Copyright (c) The OpenTofu Authors
// SPDX-License-Identifier: MPL-2.0
// Copyright (c) 2023 HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package jsontypes

import (
	"encoding/json"
	"testing"
)

func TestConfig_Validate(t *testing.T) {
	var c *Config
	if err := c.Validate(); err == nil || err.Error() != "config is nil" {
		t.Errorf("invalid error returned: %s", err)
	}
	c = &Config{}
	if err := c.Validate(); err != nil {
		t.Errorf("expected no error but got: %s", err)
	}
}

func TestConfig_UnmarshalJSON(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		c := &Config{}
		if err := c.UnmarshalJSON([]byte(`{"provider_config":{"test1": {}, "test2": {}}}`)); err != nil {
			t.Fatalf("unexpected error during unmarshalling the configuration: %s", err)
		}
		if c.ProviderConfigs == nil {
			t.Fatalf("unmarshalling failed: provider configs is missing")
		}
		if _, ok := c.ProviderConfigs["test1"]; !ok {
			t.Fatalf("unmarshalling failed: provider test1 config is missing")
		}
		if _, ok := c.ProviderConfigs["test2"]; !ok {
			t.Fatalf("unmarshalling failed: provider test2 config is missing")
		}
	})
	t.Run("error", func(t *testing.T) {
		var c *Config
		if err := json.Unmarshal([]byte(`invalid json content`), &c); err == nil || err.Error() != "invalid character 'i' looking for beginning of value" {
			t.Fatalf("unexpected error: %s", err)
		}
	})
}
