---
layout: "docs"
page_title: "Vault Enterprise Control Groups"
sidebar_current: "docs-vault-enterprise-control-groups"
description: |-
  Vault Enterprise has support for Control Group Authorization.

---

# Vault Enterprise Control Group Support

Vault Enterprise has support for Control Group Authorization. Control Groups
add additional authorization factors to be required before satisfying a request.  

When a Control Group is required for a request, a limited duration response
wrapping token is returned to the user instead of the requested data. The
accessor of the response wrapping token can be passed to the authorizers 
required by the control group policy. Once all authorizations are satisified,
the wrapping token can be used to unwrap and process the original request.

## Control Group Factors

Control Groups can verify the following factors:
 
- `Identity Groups` - Require an authorizer to be in a specific set of identity
groups.

## Control Groups In ACL Policies

Control Group requirements on paths are specified as `control_group` along 
with other ACL parameters.

### Sample ACL Policies

```
path "secret/foo" {
    capabilities = ["read"]
    control_group = {
        factor "ops_manager" {
            identity {
                group_names = ["managers"]
                approvals = 1
            }
        }
    }
}
```

The above policy grants `read` access to `secret/foo` only after one member of
the "managers" group authorizes the request.

```
path "secret/foo" {
    capabilities = ["create", "update"]
    control_group = {
        ttl = "4h"
        factor "tech leads" {
            identity {
                group_names = ["managers", "leads"]
                approvals = 2
            }
        }
        factor "super users" {
            identity {
                group_names = ["superusers"]
                approvals = 1
            }
        }
    }
}
```

The above policy grants `create` and `update` access to `secret/foo` only after 
two member of the "managers" or "leads" group and one member of the "superusers"
group authorizes the request.  If an authorizer is a member of both the 
"managers" and "superusers" group, one authorization for both factors will be 
satisfied.

## Control Groups in Sentinel

Control Groups are also supported in Sentinel policies using the `controlgroup`
import.  See [Sentinel Documentation](/docs/enterprise/sentinel/index.html) for more
details on available properties.

### Sample Sentinel Policy

```
import "time"
import "controlgroup"

control_group = func() {
    numAuthzs = 0
    for controlgroup.authorizations as authz {
		if "managers" in authz.groups.by_name {
			if time.load(authz.time).unix > time.now.unix - 3600 {
				numAuthzs = numAuthzs + 1
			}
		}
    }
    if numAuthzs >= 2 {
        return true
    }
    return false
}

main = rule {
    control_group()
}
```

The above policy will reject the request unless two members of the `managers`
group have authorized the request. Additionally it verifies the authorizations
happened in the last hour.

### API

Control Groups can be managed over the HTTP API. Please see 
[Control Groups API](/api/system/control-group.html) for more details.
