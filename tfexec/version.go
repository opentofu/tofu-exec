// Copyright (c) The OpenTofu Authors
// SPDX-License-Identifier: MPL-2.0
// Copyright (c) 2023 HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package tfexec

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"strconv"
	"strings"

	"github.com/hashicorp/go-version"
	tfjson "github.com/hashicorp/terraform-json"
)

// Version returns structured output from the tofu version command including both the OpenTofu CLI version
// and any initialized provider versions. This will read cached values when present unless the skipCache parameter
// is set to true.
func (tf *Tofu) Version(ctx context.Context, skipCache bool) (tfVersion *version.Version, providerVersions map[string]*version.Version, err error) {
	tf.versionLock.Lock()
	defer tf.versionLock.Unlock()

	if tf.execVersion == nil || skipCache {
		tf.execVersion, tf.provVersions, err = tf.version(ctx)
		if err != nil {
			return nil, nil, err
		}
	}

	return tf.execVersion, tf.provVersions, nil
}

// version does not use the locking on the Tofu instance and should probably not be used directly, prefer Version.
func (tf *Tofu) version(ctx context.Context) (*version.Version, map[string]*version.Version, error) {
	versionCmd := tf.buildTofuCmd(ctx, nil, "version", "-json")

	var outBuf bytes.Buffer
	versionCmd.Stdout = &outBuf

	err := tf.runTofuCmd(ctx, versionCmd)
	if err != nil {
		return nil, nil, err
	}

	tfVersion, providerVersions, err := parseJsonVersionOutput(outBuf.Bytes())
	if err != nil {
		return nil, nil, err
	}

	return tfVersion, providerVersions, err
}

func parseJsonVersionOutput(stdout []byte) (*version.Version, map[string]*version.Version, error) {
	var out tfjson.VersionOutput
	err := json.Unmarshal(stdout, &out)
	if err != nil {
		return nil, nil, err
	}

	tfVersion, err := version.NewVersion(out.Version)
	if err != nil {
		return nil, nil, fmt.Errorf("unable to parse version %q: %w", out.Version, err)
	}

	providerVersions := make(map[string]*version.Version)
	for provider, versionStr := range out.ProviderSelections {
		v, err := version.NewVersion(versionStr)
		if err != nil {
			return nil, nil, fmt.Errorf("unable to parse %q version %q: %w",
				provider, versionStr, err)
		}
		providerVersions[provider] = v
	}

	return tfVersion, providerVersions, nil
}

func errorVersionString(v *version.Version) string {
	if v == nil {
		return "-"
	}
	return v.String()
}

// compatible asserts compatibility of the cached tofu version with the executable, and returns a well known error if not.
func (tf *Tofu) compatible(ctx context.Context, minInclusive *version.Version, maxExclusive *version.Version) error {
	tfv, _, err := tf.Version(ctx, false)
	if err != nil {
		return err
	}
	if ok := versionInRange(tfv, minInclusive, maxExclusive); !ok {
		return &ErrVersionMismatch{
			MinInclusive: errorVersionString(minInclusive),
			MaxExclusive: errorVersionString(maxExclusive),
			Actual:       errorVersionString(tfv),
		}
	}

	return nil
}

func stripPrereleaseAndMeta(v *version.Version) *version.Version {
	if v == nil {
		return nil
	}
	segs := []string{}
	for _, s := range v.Segments() {
		segs = append(segs, strconv.Itoa(s))
	}
	vs := strings.Join(segs, ".")
	clean, _ := version.NewVersion(vs)
	return clean
}

// versionInRange checks compatibility of the tofu version. The minimum is inclusive and the max
// is exclusive, equivalent to min <= expected version < max.
//
// Pre-release information is ignored for comparison.
func versionInRange(tfv *version.Version, minInclusive *version.Version, maxExclusive *version.Version) bool {
	if minInclusive == nil && maxExclusive == nil {
		return true
	}
	tfv = stripPrereleaseAndMeta(tfv)
	minInclusive = stripPrereleaseAndMeta(minInclusive)
	maxExclusive = stripPrereleaseAndMeta(maxExclusive)
	if minInclusive != nil && !tfv.GreaterThanOrEqual(minInclusive) {
		return false
	}
	if maxExclusive != nil && !tfv.LessThan(maxExclusive) {
		return false
	}

	return true
}
