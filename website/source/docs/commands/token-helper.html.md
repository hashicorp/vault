---
layout: "docs"
page_title: "Token Helpers"
sidebar_current: "docs-commands-token-helper"
description: |-
  The Vault CLI supports external token helpers that make retrieving, setting and erasing tokens simpler to use.
---

# Token Helpers

The Vault CLI provides a built in tool for authenticating to any of the enabled auth backends. By default the Vault CLI will take the generated token after a successful authentication and store it on disk in the `~/.vault_token` file. This functionality can change in Vault via the use of a token helper. A token helper is an external program that Vault calls to save, retrieve or erase a saved token. The token helper could be a very simple script or a more complex program depending on your needs. The interface to the external token helper is extremely simple.

## Configuration

To configure a token helper, edit (or create) the file `~/.vault` and add a line similar to:

```
token_helper = "/path/to/token/helper.sh"
```

You will need to use the fully qualified path to the token helper script. The script should be executable.

## Developing a Token Helper

The interface to a token helper is extremely simple: the script is passed with one argument that could be `get`, `store` or `erase`. If the argument is `get`, the script should do whatever work it needs to do to retrieve the stored token and then print the token to `STDOUT`. If the argument is `store`, Vault is asking you to store the token. Finally, if the argument is `erase`, your program should erase the stored token.

If your program succeeds, it should exit with status code 0. If it encounters an issue that prevents it from working, it should exit with some other status code. You should write a user-friendly error message to `STDERR`. You should never write anything other than the token to `STDOUT`, as Vault assumes whatever it gets on `STDOUT` is the token.

### Example Token Helper

This is an example token helper written in Ruby that stores and retrieves tokens in a json file called `~/.vault_tokens`. The key is the environment variable $VAULT_ADDR, this allows the Vault user to easily store and retrieve tokens from a number of different Vault servers.

```
#!/usr/bin/env ruby

require 'json'

unless ENV['VAULT_ADDR']
  STDERR.puts "No VAULT_ADDR environment variable set. Set it and run me again!"
  exit 100
end

begin
  tokens = JSON.parse(File.read("#{ENV['HOME']}/.vault_tokens"))
rescue Errno::ENOENT => e
  # file doesn't exist so create a blank hash for it
  tokens = {}
end

case ARGV.first
when 'get'
  print tokens[ENV['VAULT_ADDR']] if tokens[ENV['VAULT_ADDR']]
  exit 0
when 'store'
  tokens[ENV['VAULT_ADDR']] = STDIN.read
when 'erase'
  tokens.delete!(ENV['VAULT_ADDR'])
end

File.open("#{ENV['HOME']}/.vault_tokens", 'w') { |file| file.write(tokens.to_json) }
```


