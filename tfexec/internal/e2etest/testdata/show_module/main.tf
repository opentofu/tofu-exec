# Copyright (c) The OpenTofu Authors
# SPDX-License-Identifier: MPL-2.0
# Copyright (c) 2023 HashiCorp, Inc.
# SPDX-License-Identifier: MPL-2.0

module "mod1" {
  source = "./submodule1"
  submodule1_var1 = "s1v1"
  submodule1_var2 = "s1v2"
}

module "mod2" {
  source = "./submodule2"
  submodule2_var1 = "s2v1"
}
