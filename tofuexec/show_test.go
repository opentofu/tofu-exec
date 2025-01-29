// Copyright (c) The OpenTofu Authors
// SPDX-License-Identifier: MPL-2.0
// Copyright (c) 2023 HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package tofuexec

import (
	"context"
	"testing"

	"github.com/opentofu/tofu-exec/tofuexec/internal/testutil"
)

func TestShowCmd(t *testing.T) {
	td := t.TempDir()

	tf, err := NewOpenTofu(td, tfVersion(t, testutil.Latest_v1))
	if err != nil {
		t.Fatal(err)
	}

	// empty env, to avoid environ mismatch in testing
	tf.SetEnv(map[string]string{})

	// defaults
	showCmd := tf.showCmd(context.Background(), true, nil)

	assertCmd(t, []string{
		"show",
		"-json",
		"-no-color",
	}, nil, showCmd)
}

func TestShowStateFileCmd(t *testing.T) {
	td := t.TempDir()

	tf, err := NewOpenTofu(td, tfVersion(t, testutil.Latest_v1))
	if err != nil {
		t.Fatal(err)
	}

	// empty env, to avoid environ mismatch in testing
	tf.SetEnv(map[string]string{})

	showCmd := tf.showCmd(context.Background(), true, nil, "statefilepath")

	assertCmd(t, []string{
		"show",
		"-json",
		"-no-color",
		"statefilepath",
	}, nil, showCmd)
}

func TestShowPlanFileCmd(t *testing.T) {
	td := t.TempDir()

	tf, err := NewOpenTofu(td, tfVersion(t, testutil.Latest_v1))
	if err != nil {
		t.Fatal(err)
	}

	// empty env, to avoid environ mismatch in testing
	tf.SetEnv(map[string]string{})

	showCmd := tf.showCmd(context.Background(), true, nil, "planfilepath")

	assertCmd(t, []string{
		"show",
		"-json",
		"-no-color",
		"planfilepath",
	}, nil, showCmd)
}

func TestShowPlanFileRawCmd(t *testing.T) {
	td := t.TempDir()

	tf, err := NewOpenTofu(td, tfVersion(t, testutil.Latest_v1))
	if err != nil {
		t.Fatal(err)
	}

	// empty env, to avoid environ mismatch in testing
	tf.SetEnv(map[string]string{})

	showCmd := tf.showCmd(context.Background(), false, nil, "planfilepath")

	assertCmd(t, []string{
		"show",
		"-no-color",
		"planfilepath",
	}, nil, showCmd)
}
