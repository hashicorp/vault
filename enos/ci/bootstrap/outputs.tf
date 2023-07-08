# Copyright (c) HashiCorp, Inc.
# SPDX-License-Identifier: MPL-2.0

output "keys" {
  value = {
    "us-east-1" = {
      name = aws_key_pair.enos_ci_key_us_east_1.key_name
      arn  = aws_key_pair.enos_ci_key_us_east_1.arn
    }
    "us-east-2" = {
      name = aws_key_pair.enos_ci_key_us_east_2.key_name
      arn  = aws_key_pair.enos_ci_key_us_east_2.arn
    }
    "us-west-1" = {
      name = aws_key_pair.enos_ci_key_us_west_1.key_name
      arn  = aws_key_pair.enos_ci_key_us_west_1.arn
    }
    "us-west-2" = {
      name = aws_key_pair.enos_ci_key_us_west_2.key_name
      arn  = aws_key_pair.enos_ci_key_us_west_2.arn
    }
  }
}
