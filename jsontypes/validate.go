// Copyright (c) The OpenTofu Authors
// SPDX-License-Identifier: MPL-2.0
// Copyright (c) 2023 HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package jsontypes

import (
	"encoding/json"
	"errors"
	"fmt"

	"github.com/hashicorp/go-version"
)

// ValidateFormatVersionConstraints defines the versions of the JSON
// validate format that are supported by this package.
var ValidateFormatVersionConstraints = ">= 0.1, < 2.0"

// Pos represents a position in a config file
type Pos struct {
	Line   int `json:"line"`
	Column int `json:"column"`
	Byte   int `json:"byte"`
}

// Range represents a range of bytes between two positions
type Range struct {
	Filename string `json:"filename"`
	Start    Pos    `json:"start"`
	End      Pos    `json:"end"`
}

// ValidateOutput represents JSON output from terraform validate
// (available from 0.12 onwards)
type ValidateOutput struct {
	FormatVersion string `json:"format_version"`

	Valid        bool         `json:"valid"`
	ErrorCount   int          `json:"error_count"`
	WarningCount int          `json:"warning_count"`
	Diagnostics  []Diagnostic `json:"diagnostics"`
}

// Validate checks to ensure that data is present, and the
// version matches the version supported by this library.
func (vo *ValidateOutput) Validate() error {
	if vo == nil {
		return errors.New("validation output is nil")
	}

	if vo.FormatVersion == "" {
		// The format was not versioned in the past
		return nil
	}

	constraint, err := version.NewConstraint(ValidateFormatVersionConstraints)
	if err != nil {
		return fmt.Errorf("invalid version constraint: %w", err)
	}

	v, err := version.NewVersion(vo.FormatVersion)
	if err != nil {
		return fmt.Errorf("invalid format version %q: %w", vo.FormatVersion, err)
	}

	if !constraint.Check(v) {
		return fmt.Errorf("unsupported validation output format version: %q does not satisfy %q",
			v, constraint)
	}

	return nil
}

// UnmarshalJSON implements json.Unmarshaler for ValidateOutput.
//
// As per established convention this method should only ever
// be invoked *indirectly* via [encoding/json] library.
func (vo *ValidateOutput) UnmarshalJSON(b []byte) error {
	type rawOutput ValidateOutput
	var schemas rawOutput

	err := json.Unmarshal(b, &schemas)
	if err != nil {
		return err
	}

	*vo = *(*ValidateOutput)(&schemas)

	return vo.Validate()
}
