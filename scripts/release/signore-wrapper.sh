#!/usr/bin/env bash
# Copyright (c) The OpenTofu Authors
# SPDX-License-Identifier: MPL-2.0
# Copyright (c) 2023 HashiCorp, Inc.
# SPDX-License-Identifier: MPL-2.0

printf '\n[GNUPG:] SIG_CREATED ' >&${1#--status-fd=}
signore sign --file /dev/stdin --signer $3 2>/dev/null
