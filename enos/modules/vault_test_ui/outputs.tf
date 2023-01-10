output "ui_test_stderr" {
  value = var.ui_run_tests ? enos_local_exec.test_ui[0].stderr : "No std out tests where not run"
}

output "ui_test_stdout" {
  value = var.ui_run_tests ? enos_local_exec.test_ui[0].stdout : "No std out tests where not run"
}

output "ui_test_environment" {
  value       = join(" \\ \n", local.escaped_ui_test_environment)
  description = "The environment variables that are required in order to run the test:enos yarn target"
}
