// Copyright (c) The OpenTofu Authors
// SPDX-License-Identifier: MPL-2.0
// Copyright (c) 2023 HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package tfexec

import (
	"context"
	"testing"

	"github.com/opentofu/tofu-exec/tfexec/internal/testutil"
)

func TestShowCmd(t *testing.T) {
	td := t.TempDir()

	tf, err := NewTofu(td, tfVersion(t, testutil.Latest_v1))
	if err != nil {
		t.Fatal(err)
	}

	// empty env, to avoid environ mismatch in testing
	_ = tf.SetEnv(map[string]string{})

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

	tf, err := NewTofu(td, tfVersion(t, testutil.Latest_v1))
	if err != nil {
		t.Fatal(err)
	}

	// empty env, to avoid environ mismatch in testing
	_ = tf.SetEnv(map[string]string{})

	showCmd := tf.showCmd(context.Background(), true, nil, "statefilepath")

	assertCmd(t, []string{
		"show",
		"-json",
		"-no-color",
		"statefilepath",
	}, nil, showCmd)
}

func TestShowModule_unsupportedTofuVersion(t *testing.T) {
	tf, err := NewTofu(t.TempDir(), tfVersion(t, testutil.Latest_v1))
	if err != nil {
		t.Fatal(err)
	}

	_ = tf.SetEnv(map[string]string{})

	_, err = tf.ShowModule(context.Background(), "foo/bar")
	if err == nil {
		t.Fatalf("expected error for unsupported version, got nil")
	}
}

func TestShowModuleCmd(t *testing.T) {
	td := t.TempDir()

	tf, err := NewTofu(td, tfVersion(t, testutil.Latest_v1_11))
	if err != nil {
		t.Fatal(err)
	}

	_ = tf.SetEnv(map[string]string{})

	showCmd := tf.showCmd(context.Background(), true, nil, "-module=foo/bar")
	assertCmd(t, []string{
		"show",
		"-json",
		"-no-color",
		"-module=foo/bar",
	}, nil, showCmd)
}

func TestShowModule_blankDir(t *testing.T) {
	td := t.TempDir()
	tf, err := NewTofu(td, tfVersion(t, testutil.Latest_v1_11))
	if err != nil {
		t.Fatal(err)
	}
	_ = tf.SetEnv(map[string]string{})

	_, err = tf.ShowModule(context.Background(), "")
	if err == nil {
		t.Fatalf("expected error for blank moduleDir, got nil")
	}
}

func TestShowPlanFileCmd(t *testing.T) {
	td := t.TempDir()

	tf, err := NewTofu(td, tfVersion(t, testutil.Latest_v1))
	if err != nil {
		t.Fatal(err)
	}

	// empty env, to avoid environ mismatch in testing
	_ = tf.SetEnv(map[string]string{})

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

	tf, err := NewTofu(td, tfVersion(t, testutil.Latest_v1))
	if err != nil {
		t.Fatal(err)
	}

	// empty env, to avoid environ mismatch in testing
	_ = tf.SetEnv(map[string]string{})

	showCmd := tf.showCmd(context.Background(), false, nil, "planfilepath")

	assertCmd(t, []string{
		"show",
		"-no-color",
		"planfilepath",
	}, nil, showCmd)
}
