// Copyright (c) The OpenTofu Authors
// SPDX-License-Identifier: MPL-2.0
// Copyright (c) 2023 HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package e2etest

import (
	"context"
	"fmt"
	"hash/crc32"
	"io"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/hashicorp/go-version"

	"github.com/opentofu/tofu-exec/tfexec"
	"github.com/opentofu/tofu-exec/tfexec/internal/testutil"
)

const testFixtureDir = "testdata"

func runTest(t *testing.T, fixtureName string, cb func(t *testing.T, tfVersion *version.Version, tf *tfexec.Tofu)) {
	t.Helper()

	versions := []string{
		testutil.Latest_v1_11,
		testutil.Latest_v1,
		testutil.Latest_v1_9,
		testutil.Latest_v1_8,
		testutil.Latest_v1_7,
		testutil.Latest_v1_6,
	}
	if override := os.Getenv("TFEXEC_E2ETEST_VERSIONS"); override != "" {
		versions = strings.Split(override, ",")
	}

	// If the env var TFEXEC_E2ETEST_TOFU_PATH is set to the path of a
	// valid OpenTofu executable, only tests appropriate to that
	// executable's version will be run.
	if localBinPath := os.Getenv("TFEXEC_E2ETEST_TOFU_PATH"); localBinPath != "" {
		// By convention, every new Tofu struct is given a clean
		// temp dir, even if we are only invoking tf.Version(). This
		// prevents any possible confusion that could result from
		// reusing an t.TempDir() (for example) that already contained
		// OpenTofu files.
		td := t.TempDir()
		ltf, err := tfexec.NewTofu(td, localBinPath)
		if err != nil {
			t.Fatal(err)
		}

		_ = ltf.SetAppendUserAgent("tfexec-e2etest")

		lVersion, _, err := ltf.Version(context.Background(), false)
		if err != nil {
			t.Fatalf("unable to determine version of Tofu binary at %s: %s", localBinPath, err)
		}

		versions = []string{lVersion.String()}
	}

	runTestWithVersions(t, versions, fixtureName, cb)
}

func runTestWithVersions(t *testing.T, versions []string, fixtureName string, cb func(t *testing.T, tfVersion *version.Version, tf *tfexec.Tofu)) {
	t.Helper()

	alreadyRunVersions := map[string]bool{}
	for _, tfv := range versions {
		t.Run(fmt.Sprintf("%s-%s", fixtureName, tfv), func(t *testing.T) {
			if alreadyRunVersions[tfv] {
				t.Skipf("already run version %q", tfv)
			}
			alreadyRunVersions[tfv] = true

			td, err := os.MkdirTemp("", "tf")
			if err != nil {
				t.Fatalf("error creating temporary test directory: %s", err)
			}
			t.Cleanup(func() {
				_ = os.RemoveAll(td)
			})

			var execPath string
			if localBinPath := os.Getenv("TFEXEC_E2ETEST_TOFU_PATH"); localBinPath != "" {
				execPath = localBinPath
			} else {
				execPath = tfcache.Version(t, tfv)
			}

			tf, err := tfexec.NewTofu(td, execPath)
			if err != nil {
				t.Fatal(err)
			}

			_ = tf.SetAppendUserAgent("tfexec-e2etest")

			runningVersion, _, err := tf.Version(context.Background(), false)
			if err != nil {
				t.Fatalf("unable to determine running version (expected %q): %s", tfv, err)
			}

			// Check that the runningVersion matches the expected
			// test version. This ensures non-matching tests are
			// skipped when using a local tofu executable.
			if !strings.HasPrefix(tfv, "refs/") {
				testVersion, err := version.NewVersion(tfv)
				if err != nil {
					t.Fatalf("unable to parse version %s: %s", testVersion, err)
				}
				if !testVersion.Equal(runningVersion) {
					t.Skipf("test applies to version %s, but local executable is version %s", tfv, runningVersion)
				}
			}

			if fixtureName != "" {
				err = copyFiles(filepath.Join(testFixtureDir, fixtureName), td)
				if err != nil {
					t.Fatalf("error copying config file into test dir: %s", err)
				}
			}

			// Separate strings.Builder because it's not concurrent safe
			var stdout strings.Builder
			tf.SetStdout(&stdout)
			var stderr strings.Builder
			tf.SetStderr(&stderr)

			tf.SetLogger(&testingPrintfer{t})

			// TODO: capture panics here?
			cb(t, runningVersion, tf)

			t.Logf("CLI Output:\n%s", stdout.String())
			if len(stderr.String()) > 0 {
				t.Logf("CLI Error:\n%s", stderr.String())
			}
		})
	}
}

type testingPrintfer struct {
	t *testing.T
}

func (t *testingPrintfer) Printf(format string, v ...interface{}) {
	t.t.Logf(format, v...)
}

func copyFiles(path string, dstPath string) error {
	infos, err := ioutil.ReadDir(path)
	if err != nil {
		return err
	}

	for _, info := range infos {
		srcPath := filepath.Join(path, info.Name())
		if info.IsDir() {
			newDir := filepath.Join(dstPath, info.Name())
			err = os.MkdirAll(newDir, info.Mode())
			if err != nil {
				return err
			}
			err = copyFiles(srcPath, newDir)
			if err != nil {
				return err
			}
		} else {
			err = copyFile(srcPath, dstPath)
			if err != nil {
				return err
			}
		}

	}
	return nil
}

func copyFile(path string, dstPath string) error {
	srcF, err := os.Open(path)
	if err != nil {
		return err
	}
	defer srcF.Close()

	di, err := os.Stat(dstPath)
	if err != nil {
		return err
	}
	if di.IsDir() {
		_, file := filepath.Split(path)
		dstPath = filepath.Join(dstPath, file)
	}

	dstF, err := os.Create(dstPath)
	if err != nil {
		return err
	}
	defer dstF.Close()

	if _, err := io.Copy(dstF, srcF); err != nil {
		return err
	}

	return nil
}

// filesEqual asserts that two files have the same contents.
func textFilesEqual(t *testing.T, expected, actual string) {
	eb, err := os.ReadFile(expected)
	if err != nil {
		t.Fatal(err)
	}

	ab, err := os.ReadFile(actual)
	if err != nil {
		t.Fatal(err)
	}

	es := string(eb)
	es = strings.ReplaceAll(es, "\r\n", "\n")

	as := string(ab)
	as = strings.ReplaceAll(as, "\r\n", "\n")

	if as != es {
		t.Fatalf("expected:\n%s\n\ngot:\n%s\n", es, as)
	}
}

func checkSum(t *testing.T, filename string) uint32 {
	b, err := os.ReadFile(filename)
	if err != nil {
		t.Fatal(err)
	}
	return crc32.ChecksumIEEE(b)
}
