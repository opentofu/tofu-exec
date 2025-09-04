// Copyright (c) The OpenTofu Authors
// SPDX-License-Identifier: MPL-2.0
// Copyright (c) 2023 HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package tfexec

import (
	"context"
	"errors"
	"io/ioutil"
	"os"
	"path/filepath"
	"testing"

	"github.com/opentofu/tofu-exec/tfexec/internal/testutil"
)

var tfCache *testutil.TFCache

func TestMain(m *testing.M) {
	os.Exit(func() int {
		var err error
		installDir, err := ioutil.TempDir("", "tfinstall")
		if err != nil {
			panic(err)
		}
		defer os.RemoveAll(installDir)

		tfCache = testutil.NewTFCache(installDir)

		return m.Run()
	}())
}

func TestSetEnv(t *testing.T) {
	td := t.TempDir()

	tf, err := NewTofu(td, tfVersion(t, testutil.Latest_v1))
	if err != nil {
		t.Fatal(err)
	}

	for _, c := range []struct {
		errManual bool
		name      string
	}{
		{false, "OK_ENV_VAR"},

		{true, "TF_LOG"},
		{true, "TF_VAR_foo"},
	} {
		t.Run(c.name, func(t *testing.T) {
			err = tf.SetEnv(map[string]string{c.name: "foo"})

			if c.errManual {
				var evErr *ErrManualEnvVar
				if !errors.As(err, &evErr) {
					t.Fatalf("expected ErrManualEnvVar, got %T %s", err, err)
				}
			} else {
				if !c.errManual && err != nil {
					t.Fatal(err)
				}
			}
		})
	}
}

func TestSetLog(t *testing.T) {
	td := t.TempDir()

	tf, err := NewTofu(td, tfVersion(t, testutil.Latest_v1))

	if err != nil {
		t.Fatalf("unexpected NewTofu error: %s", err)
	}

	// Required so all testing environment variables are not copied.
	err = tf.SetEnv(map[string]string{
		"CLEARENV": "1",
	})

	if err != nil {
		t.Fatalf("unexpected SetEnv error: %s", err)
	}

	t.Run("SetLog TRACE no SetLogPath", func(t *testing.T) {
		err := tf.SetLog("TRACE")

		if err != nil {
			t.Fatalf("unexpected SetLog error: %s", err)
		}

		initCmd, err := tf.initCmd(context.Background())

		if err != nil {
			t.Fatalf("unexpected command error: %s", err)
		}

		assertCmd(t, []string{
			"init",
			"-no-color",
			"-input=false",
			"-backend=true",
			"-get=true",
			"-upgrade=false",
		}, map[string]string{
			"CLEARENV":        "1",
			"TF_LOG":          "",
			"TF_LOG_CORE":     "",
			"TF_LOG_PATH":     "",
			"TF_LOG_PROVIDER": "",
		}, initCmd)
	})

	t.Run("SetLog TRACE and SetLogPath", func(t *testing.T) {
		tfLogPath := filepath.Join(td, "test.log")

		err := tf.SetLog("TRACE")

		if err != nil {
			t.Fatalf("unexpected SetLog error: %s", err)
		}

		err = tf.SetLogPath(tfLogPath)

		if err != nil {
			t.Fatalf("unexpected SetLogPath error: %s", err)
		}

		initCmd, err := tf.initCmd(context.Background())

		if err != nil {
			t.Fatalf("unexpected command error: %s", err)
		}

		assertCmd(t, []string{
			"init",
			"-no-color",
			"-input=false",
			"-backend=true",
			"-get=true",
			"-upgrade=false",
		}, map[string]string{
			"CLEARENV":        "1",
			"TF_LOG":          "TRACE",
			"TF_LOG_CORE":     "",
			"TF_LOG_PATH":     tfLogPath,
			"TF_LOG_PROVIDER": "",
		}, initCmd)
	})

	t.Run("SetLog DEBUG and SetLogPath", func(t *testing.T) {
		tfLogPath := filepath.Join(td, "test.log")

		err := tf.SetLog("DEBUG")

		if err != nil {
			t.Fatalf("unexpected SetLog error: %s", err)
		}

		err = tf.SetLogPath(tfLogPath)

		if err != nil {
			t.Fatalf("unexpected SetLogPath error: %s", err)
		}

		initCmd, err := tf.initCmd(context.Background())

		if err != nil {
			t.Fatalf("unexpected command error: %s", err)
		}

		assertCmd(t, []string{
			"init",
			"-no-color",
			"-input=false",
			"-backend=true",
			"-get=true",
			"-upgrade=false",
		}, map[string]string{
			"CLEARENV":        "1",
			"TF_LOG":          "DEBUG",
			"TF_LOG_CORE":     "",
			"TF_LOG_PATH":     tfLogPath,
			"TF_LOG_PROVIDER": "",
		}, initCmd)
	})
}

func TestSetLogCore(t *testing.T) {
	td := t.TempDir()

	tf, err := NewTofu(td, tfVersion(t, testutil.Latest_v1))

	if err != nil {
		t.Fatalf("unexpected NewTofu error: %s", err)
	}

	// Required so all testing environment variables are not copied.
	err = tf.SetEnv(map[string]string{
		"CLEARENV": "1",
	})

	if err != nil {
		t.Fatalf("unexpected SetEnv error: %s", err)
	}

	t.Run("SetLogCore TRACE no SetLogPath", func(t *testing.T) {
		err := tf.SetLogCore("TRACE")

		if err != nil {
			t.Fatalf("unexpected SetLogCore error: %s", err)
		}

		initCmd, err := tf.initCmd(context.Background())

		if err != nil {
			t.Fatalf("unexpected command error: %s", err)
		}

		assertCmd(t, []string{
			"init",
			"-no-color",
			"-input=false",
			"-backend=true",
			"-get=true",
			"-upgrade=false",
		}, map[string]string{
			"CLEARENV":        "1",
			"TF_LOG":          "",
			"TF_LOG_CORE":     "",
			"TF_LOG_PATH":     "",
			"TF_LOG_PROVIDER": "",
		}, initCmd)
	})

	t.Run("SetLogCore TRACE and SetLogPath", func(t *testing.T) {
		tfLogPath := filepath.Join(td, "test.log")

		err := tf.SetLogCore("TRACE")

		if err != nil {
			t.Fatalf("unexpected SetLogCore error: %s", err)
		}

		err = tf.SetLogPath(tfLogPath)

		if err != nil {
			t.Fatalf("unexpected SetLogPath error: %s", err)
		}

		initCmd, err := tf.initCmd(context.Background())

		if err != nil {
			t.Fatalf("unexpected command error: %s", err)
		}

		assertCmd(t, []string{
			"init",
			"-no-color",
			"-input=false",
			"-backend=true",
			"-get=true",
			"-upgrade=false",
		}, map[string]string{
			"CLEARENV":        "1",
			"TF_LOG":          "",
			"TF_LOG_CORE":     "TRACE",
			"TF_LOG_PATH":     tfLogPath,
			"TF_LOG_PROVIDER": "",
		}, initCmd)
	})

	t.Run("SetLogCore DEBUG and SetLogPath", func(t *testing.T) {
		tfLogPath := filepath.Join(td, "test.log")

		err := tf.SetLogCore("DEBUG")

		if err != nil {
			t.Fatalf("unexpected SetLogCore error: %s", err)
		}

		err = tf.SetLogPath(tfLogPath)

		if err != nil {
			t.Fatalf("unexpected SetLogPath error: %s", err)
		}

		initCmd, err := tf.initCmd(context.Background())

		if err != nil {
			t.Fatalf("unexpected command error: %s", err)
		}

		assertCmd(t, []string{
			"init",
			"-no-color",
			"-input=false",
			"-backend=true",
			"-get=true",
			"-upgrade=false",
		}, map[string]string{
			"CLEARENV":        "1",
			"TF_LOG":          "",
			"TF_LOG_CORE":     "DEBUG",
			"TF_LOG_PATH":     tfLogPath,
			"TF_LOG_PROVIDER": "",
		}, initCmd)
	})
}

func TestSetLogPath(t *testing.T) {
	td := t.TempDir()

	tf, err := NewTofu(td, tfVersion(t, testutil.Latest_v1))

	if err != nil {
		t.Fatalf("unexpected NewTofu error: %s", err)
	}

	// Required so all testing environment variables are not copied.
	err = tf.SetEnv(map[string]string{
		"CLEARENV": "1",
	})

	if err != nil {
		t.Fatalf("unexpected SetEnv error: %s", err)
	}

	t.Run("case 1: No SetLogPath", func(t *testing.T) {
		initCmd, err := tf.initCmd(context.Background())

		if err != nil {
			t.Fatalf("unexpected command error: %s", err)
		}

		assertCmd(t, []string{
			"init",
			"-no-color",
			"-input=false",
			"-backend=true",
			"-get=true",
			"-upgrade=false",
		}, map[string]string{
			"CLEARENV":        "1",
			"TF_LOG":          "",
			"TF_LOG_CORE":     "",
			"TF_LOG_PATH":     "",
			"TF_LOG_PROVIDER": "",
		}, initCmd)
	})

	t.Run("case 2: SetLogPath sets TF_LOG (if no TF_LOG_CORE or TF_LOG_PROVIDER) and TF_LOG_PATH", func(t *testing.T) {
		tfLogPath := filepath.Join(td, "test.log")

		err = tf.SetLogPath(tfLogPath)

		if err != nil {
			t.Fatalf("unexpected SetLogPath error: %s", err)
		}

		initCmd, err := tf.initCmd(context.Background())

		if err != nil {
			t.Fatalf("unexpected command error: %s", err)
		}

		assertCmd(t, []string{
			"init",
			"-no-color",
			"-input=false",
			"-backend=true",
			"-get=true",
			"-upgrade=false",
		}, map[string]string{
			"CLEARENV":        "1",
			"TF_LOG":          "TRACE",
			"TF_LOG_CORE":     "",
			"TF_LOG_PATH":     tfLogPath,
			"TF_LOG_PROVIDER": "",
		}, initCmd)
	})

	t.Run("case 3: SetLogPath does not set TF_LOG if TF_LOG_CORE", func(t *testing.T) {
		tfLogPath := filepath.Join(td, "test.log")

		err := tf.SetLog("")

		if err != nil {
			t.Fatalf("unexpected SetLog error: %s", err)
		}

		err = tf.SetLogCore("TRACE")

		if err != nil {
			t.Fatalf("unexpected SetLogCore error: %s", err)
		}

		err = tf.SetLogProvider("")

		if err != nil {
			t.Fatalf("unexpected SetLogProvider error: %s", err)
		}

		err = tf.SetLogPath(tfLogPath)

		if err != nil {
			t.Fatalf("unexpected SetLogPath error: %s", err)
		}

		initCmd, err := tf.initCmd(context.Background())

		if err != nil {
			t.Fatalf("unexpected command error: %s", err)
		}

		assertCmd(t, []string{
			"init",
			"-no-color",
			"-input=false",
			"-backend=true",
			"-get=true",
			"-upgrade=false",
		}, map[string]string{
			"CLEARENV":        "1",
			"TF_LOG":          "",
			"TF_LOG_CORE":     "TRACE",
			"TF_LOG_PATH":     tfLogPath,
			"TF_LOG_PROVIDER": "",
		}, initCmd)
	})

	t.Run("case 4: SetLogPath does not set TF_LOG if TF_LOG_PROVIDER", func(t *testing.T) {
		tfLogPath := filepath.Join(td, "test.log")

		err := tf.SetLog("")

		if err != nil {
			t.Fatalf("unexpected SetLog error: %s", err)
		}

		err = tf.SetLogCore("")

		if err != nil {
			t.Fatalf("unexpected SetLogCore error: %s", err)
		}

		err = tf.SetLogProvider("TRACE")

		if err != nil {
			t.Fatalf("unexpected SetLogProvider error: %s", err)
		}

		err = tf.SetLogPath(tfLogPath)

		if err != nil {
			t.Fatalf("unexpected SetLogPath error: %s", err)
		}

		initCmd, err := tf.initCmd(context.Background())

		if err != nil {
			t.Fatalf("unexpected command error: %s", err)
		}

		assertCmd(t, []string{
			"init",
			"-no-color",
			"-input=false",
			"-backend=true",
			"-get=true",
			"-upgrade=false",
		}, map[string]string{
			"CLEARENV":        "1",
			"TF_LOG":          "",
			"TF_LOG_CORE":     "",
			"TF_LOG_PATH":     tfLogPath,
			"TF_LOG_PROVIDER": "TRACE",
		}, initCmd)
	})
}

func TestSetLogProvider(t *testing.T) {
	td := t.TempDir()

	tf, err := NewTofu(td, tfVersion(t, testutil.Latest_v1))

	if err != nil {
		t.Fatalf("unexpected NewTofu error: %s", err)
	}

	// Required so all testing environment variables are not copied.
	err = tf.SetEnv(map[string]string{
		"CLEARENV": "1",
	})

	if err != nil {
		t.Fatalf("unexpected SetEnv error: %s", err)
	}

	t.Run("SetLogProvider TRACE no SetLogPath", func(t *testing.T) {
		err := tf.SetLogProvider("TRACE")

		if err != nil {
			t.Fatalf("unexpected SetLogProvider error: %s", err)
		}

		initCmd, err := tf.initCmd(context.Background())

		if err != nil {
			t.Fatalf("unexpected command error: %s", err)
		}

		assertCmd(t, []string{
			"init",
			"-no-color",
			"-input=false",
			"-backend=true",
			"-get=true",
			"-upgrade=false",
		}, map[string]string{
			"CLEARENV":        "1",
			"TF_LOG":          "",
			"TF_LOG_CORE":     "",
			"TF_LOG_PATH":     "",
			"TF_LOG_PROVIDER": "",
		}, initCmd)
	})

	t.Run("SetLogProvider TRACE and SetLogPath", func(t *testing.T) {
		tfLogPath := filepath.Join(td, "test.log")

		err := tf.SetLogProvider("TRACE")

		if err != nil {
			t.Fatalf("unexpected SetLogProvider error: %s", err)
		}

		err = tf.SetLogPath(tfLogPath)

		if err != nil {
			t.Fatalf("unexpected SetLogPath error: %s", err)
		}

		initCmd, err := tf.initCmd(context.Background())

		if err != nil {
			t.Fatalf("unexpected command error: %s", err)
		}

		assertCmd(t, []string{
			"init",
			"-no-color",
			"-input=false",
			"-backend=true",
			"-get=true",
			"-upgrade=false",
		}, map[string]string{
			"CLEARENV":        "1",
			"TF_LOG":          "",
			"TF_LOG_CORE":     "",
			"TF_LOG_PATH":     tfLogPath,
			"TF_LOG_PROVIDER": "TRACE",
		}, initCmd)
	})

	t.Run("SetLogProvider DEBUG and SetLogPath", func(t *testing.T) {
		tfLogPath := filepath.Join(td, "test.log")

		err := tf.SetLogProvider("DEBUG")

		if err != nil {
			t.Fatalf("unexpected SetLogProvider error: %s", err)
		}

		err = tf.SetLogPath(tfLogPath)

		if err != nil {
			t.Fatalf("unexpected SetLogPath error: %s", err)
		}

		initCmd, err := tf.initCmd(context.Background())

		if err != nil {
			t.Fatalf("unexpected command error: %s", err)
		}

		assertCmd(t, []string{
			"init",
			"-no-color",
			"-input=false",
			"-backend=true",
			"-get=true",
			"-upgrade=false",
		}, map[string]string{
			"CLEARENV":        "1",
			"TF_LOG":          "",
			"TF_LOG_CORE":     "",
			"TF_LOG_PATH":     tfLogPath,
			"TF_LOG_PROVIDER": "DEBUG",
		}, initCmd)
	})
}

func TestCheckpointDisablePropagation_v1(t *testing.T) {
	td := t.TempDir()

	tf, err := NewTofu(td, tfVersion(t, testutil.Latest_v1))
	if err != nil {
		t.Fatal(err)
	}

	err = os.Setenv("CHECKPOINT_DISABLE", "1")
	if err != nil {
		t.Fatal(err)
	}
	defer os.Unsetenv("CHECKPOINT_DISABLE")

	t.Run("case 1: env var is set in environment and not overridden", func(t *testing.T) {

		err = tf.SetEnv(map[string]string{
			"FOOBAR": "1",
		})
		if err != nil {
			t.Fatal(err)
		}

		initCmd, err := tf.initCmd(context.Background())
		if err != nil {
			t.Fatal(err)
		}

		assertCmd(t, []string{
			"init",
			"-no-color",
			"-input=false",
			"-backend=true",
			"-get=true",
			"-upgrade=false",
		}, map[string]string{
			"CHECKPOINT_DISABLE": "1",
			"FOOBAR":             "1",
		}, initCmd)
	})

	t.Run("case 2: env var is set in environment and overridden with SetEnv", func(t *testing.T) {
		err = tf.SetEnv(map[string]string{
			"CHECKPOINT_DISABLE": "",
			"FOOBAR":             "2",
		})
		if err != nil {
			t.Fatal(err)
		}

		initCmd, err := tf.initCmd(context.Background())
		if err != nil {
			t.Fatal(err)
		}

		assertCmd(t, []string{
			"init",
			"-no-color",
			"-input=false",
			"-backend=true",
			"-get=true",
			"-upgrade=false",
		}, map[string]string{
			"CHECKPOINT_DISABLE": "",
			"FOOBAR":             "2",
		}, initCmd)
	})
}

// test that a suitable error is returned if NewTofu is called without a valid
// executable path
func TestNoTofuBinary(t *testing.T) {
	td := t.TempDir()

	_, err := NewTofu(td, "")
	if err == nil {
		t.Fatal("expected NewTofu to error, but it did not")
	}

	var e *ErrNoSuitableBinary
	if !errors.As(err, &e) {
		t.Fatal("expected error to be ErrNoSuitableBinary")
	}
}

func tfVersion(t *testing.T, v string) string {
	if tfCache == nil {
		t.Fatalf("tfCache not yet configured, TestMain must run first")
	}

	return tfCache.Version(t, v)
}
