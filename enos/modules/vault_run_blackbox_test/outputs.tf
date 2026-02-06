# Copyright IBM Corp. 2016, 2025
# SPDX-License-Identifier: BUSL-1.1

output "test_result" {
  description = "The result of the blackbox test execution (human-readable)"
  value       = enos_local_exec.run_blackbox_test.stdout
}


output "test_results_summary" {
  description = "Summary of test results for dashboards"
  value = {
    status       = local.test_status
    passed       = local.test_status == "PASSED"
    exit_code    = local.test_exit_code
    timestamp    = timestamp()
    json_file    = local.json_file_path
    test_filter  = length(var.test_names) > 0 ? join(", ", var.test_names) : "all tests"
    test_package = var.test_package
  }
}
