disable_cache = true
disable_mlock = true
enable_cors     = true
allowed_origins = "http://localhost:8[0-9]{3}"

backend "consul" {
    foo = "bar"
    disable_clustering = "true"
}
