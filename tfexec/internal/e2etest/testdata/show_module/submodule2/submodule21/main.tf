# Copyright (c) The OpenTofu Authors
# SPDX-License-Identifier: MPL-2.0
# Copyright (c) 2023 HashiCorp, Inc.
# SPDX-License-Identifier: MPL-2.0

variable "submodule21_var1" {
  type = string
}

variable "submodule21_var2" {
  type = string
}



output "submodule21_out1" {
  value = var.submodule21_var1
}

output "submodule21_out2" {
  value = var.submodule21_var2
}
