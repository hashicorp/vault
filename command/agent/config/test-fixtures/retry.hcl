vault {
  retry {
    enabled = true
    attempts = 6
    backoff = "1s"
    max_backoff = "1s"
  }
}
