# Copyright (c) The OpenTofu Authors
# SPDX-License-Identifier: MPL-2.0
# Copyright (c) 2023 HashiCorp, Inc.
# SPDX-License-Identifier: MPL-2.0

variable "submodule1_var1" {
  type = string
}

variable "submodule1_var2" {
  type = string
}

output "submodule1_out" {
  value = concat(var.submodule1_var1, var.submodule1_var2)
}
