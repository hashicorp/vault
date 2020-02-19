disable_cache = true
disable_mlock = true

backend "consul" {
    foo = "faz"
}


# this does not override previous disable_cache entries
disable_cache = false

# this does
disable_clustering = false