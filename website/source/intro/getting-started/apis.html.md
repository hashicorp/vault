---
layout: "intro"
page_title: "Using the REST APIs"
sidebar_current: "gettingstarted-apis"
description: |-
  Using the REST APIs for authentication and secret access.
---

# Using the REST APIs
All Vault capabilities can be accessed via HTTP, rather than the CLI. In fact some calls, for example [app-id](/docs/auth/app-id.html) authentication cannot be called by the CLI at all. Once you have started your server, you can use curl, or any other http client to make API calls. For example, if you have started Vault in dev mode, you could validated initialization status like this:

```
$ curl http://localhost:8200/v1/sys/init

{"initialized":true}
```

# Access Secrets via the REST APIs
Machines will most likely access Vault via the REST APIs. Assuming a machine is using the [app-id](/docs/auth/app-id.html) backend for authentication, the flow would look like this:
![REST Sequence](/assets/images/app-id-api-sequence.png)
Executing this flow on the command line via curl would look something like this:
```
$ curl

```
