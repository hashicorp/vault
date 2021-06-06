HashiCorp-internal libs
=================

Do not use these unless you know what you're doing.

These libraries are used by other HashiCorp software to reduce code duplication
and increase consistency. They are not libraries needed by Vault plugins --
those are in the sdk/ module.

There are no compatibility guarantees. Things in here may change or move or
disappear at any time.

If you are a Vault plugin author and think you need a library in here in your
plugin, please open an issue for discussion.
