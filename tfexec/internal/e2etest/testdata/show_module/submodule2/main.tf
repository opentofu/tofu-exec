# Copyright (c) The OpenTofu Authors
# SPDX-License-Identifier: MPL-2.0
# Copyright (c) 2023 HashiCorp, Inc.
# SPDX-License-Identifier: MPL-2.0

module "submodule21" {
  source = "./submodule21"
  submodule21_var1 = var.submodule2_var1
  submodule21_var2 = "s21v2"
}

variable "submodule2_var1" {
  type = string
}
