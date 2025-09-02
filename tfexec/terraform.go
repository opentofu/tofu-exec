// Copyright (c) The OpenTofu Authors
// SPDX-License-Identifier: MPL-2.0
// Copyright (c) 2023 HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package tfexec

import (
	"fmt"
	"io"
	"log"
	"os"
	"sync"

	"github.com/hashicorp/go-version"
)

type printfer interface {
	Printf(format string, v ...interface{})
}

// Tofu represents the Tofu CLI executable and working directory.
//
// Typically this is constructed against the root module of a Tofu configuration
// but you can override paths used in some commands depending on the available
// options.
//
// All functions that execute CLI commands take a context.Context. It should be noted that
// exec.Cmd.Run will not return context.DeadlineExceeded or context.Canceled by default, we
// have augmented our wrapped errors to respond true to errors.Is for context.DeadlineExceeded
// and context.Canceled if those are present on the context when the error is parsed. See
// https://github.com/golang/go/issues/21880 for more about the Go limitations.
//
// By default, the instance inherits the environment from the calling code (using os.Environ)
// but it ignores certain environment variables that are managed within the code and prohibits
// setting them through SetEnv:
//
//   - TF_APPEND_USER_AGENT
//   - TF_IN_AUTOMATION
//   - TF_INPUT
//   - TF_LOG
//   - TF_LOG_PATH
//   - TF_REATTACH_PROVIDERS
//   - TF_DISABLE_PLUGIN_TLS
//   - TF_SKIP_PROVIDER_VERIFY
type Tofu struct {
	execPath           string
	workingDir         string
	appendUserAgent    string
	disablePluginTLS   bool
	skipProviderVerify bool
	env                map[string]string

	stdout io.Writer
	stderr io.Writer
	logger printfer

	// TF_LOG environment variable, defaults to TRACE if logPath is set.
	log string

	// TF_LOG_CORE environment variable
	logCore string

	// TF_LOG_PATH environment variable
	logPath string

	// TF_LOG_PROVIDER environment variable
	logProvider string

	versionLock  sync.Mutex
	execVersion  *version.Version
	provVersions map[string]*version.Version
}

// NewTofu returns a Tofu struct with default values for all fields.
// If a blank execPath is supplied, NewTofu will error.
// Use tofudl or output from os.LookPath to get a desirable execPath.
func NewTofu(workingDir string, execPath string) (*Tofu, error) {
	if workingDir == "" {
		return nil, fmt.Errorf("failed to initialize tofu-exec (NewTofu): cannot be initialised with empty working directory")
	}

	if _, err := os.Stat(workingDir); err != nil {
		return nil, fmt.Errorf("failed to initialize tofu-exec (NewTofu): error with working directory %s: %s", workingDir, err)
	}

	if execPath == "" {
		err := fmt.Errorf("failed to initialize tofu-exec (NewTofu): please supply the path to a Tofu executable using execPath, e.g. using the github.com/opentofu/tofudl library")
		return nil, &ErrNoSuitableBinary{
			err: err,
		}
	}
	tf := Tofu{
		execPath:   execPath,
		workingDir: workingDir,
		env:        nil, // explicit nil means copy os.Environ
		logger:     log.New(io.Discard, "", 0),
	}

	return &tf, nil
}

// SetEnv allows you to override environment variables, this should not be used for any well known
// OpenTofu environment variables that are already covered in options. Pass nil to copy the values
// from os.Environ. Attempting to set environment variables that should be managed manually will
// result in ErrManualEnvVar being returned.
func (tf *Tofu) SetEnv(env map[string]string) error {
	prohibited := ProhibitedEnv(env)
	if len(prohibited) > 0 {
		// just error on the first instance
		return &ErrManualEnvVar{prohibited[0]}
	}

	tf.env = env
	return nil
}

// SetLogger specifies a logger for tfexec to use.
func (tf *Tofu) SetLogger(logger printfer) {
	tf.logger = logger
}

// SetStdout specifies a writer to stream stdout to for every command.
//
// This should be used for information or logging purposes only, not control
// flow. Any parsing necessary should be added as functionality to this package.
func (tf *Tofu) SetStdout(w io.Writer) {
	tf.stdout = w
}

// SetStderr specifies a writer to stream stderr to for every command.
//
// This should be used for information or logging purposes only, not control
// flow. Any parsing necessary should be added as functionality to this package.
func (tf *Tofu) SetStderr(w io.Writer) {
	tf.stderr = w
}

// SetLog sets the TF_LOG environment variable for OpenTofu CLI execution.
// This must be combined with a call to SetLogPath to take effect.
func (tf *Tofu) SetLog(log string) error {
	tf.log = log
	return nil
}

// SetLogCore sets the TF_LOG_CORE environment variable for OpenTofu CLI
// execution. This must be combined with a call to SetLogPath to take effect.
func (tf *Tofu) SetLogCore(logCore string) error {
	tf.logCore = logCore
	return nil
}

// SetLogPath sets the TF_LOG_PATH environment variable for OpenTofu CLI
// execution.
func (tf *Tofu) SetLogPath(path string) error {
	tf.logPath = path
	// Prevent setting the log path without enabling logging
	if tf.log == "" && tf.logCore == "" && tf.logProvider == "" {
		tf.log = "TRACE"
	}
	return nil
}

// SetLogProvider sets the TF_LOG_PROVIDER environment variable for OpenTofu
// CLI execution. This must be combined with a call to SetLogPath to take
// effect.
func (tf *Tofu) SetLogProvider(logProvider string) error {
	tf.logProvider = logProvider
	return nil
}

// SetAppendUserAgent sets the TF_APPEND_USER_AGENT environment variable for
// OpenTofu CLI execution.
func (tf *Tofu) SetAppendUserAgent(ua string) error {
	tf.appendUserAgent = ua
	return nil
}

// SetDisablePluginTLS sets the TF_DISABLE_PLUGIN_TLS environment variable for
// OpenTofu CLI execution.
func (tf *Tofu) SetDisablePluginTLS(disabled bool) error {
	tf.disablePluginTLS = disabled
	return nil
}

// WorkingDir returns the working directory for OpenTofu.
func (tf *Tofu) WorkingDir() string {
	return tf.workingDir
}

// ExecPath returns the path to the tofu executable.
func (tf *Tofu) ExecPath() string {
	return tf.execPath
}
