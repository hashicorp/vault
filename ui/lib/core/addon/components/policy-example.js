/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: BUSL-1.1
 */

import Component from '@glimmer/component';

/**
 * @module PolicyExample
 * The PolicyExample component receives a policy type ('acl', 'rgp', or 'egp') and renders a copyable policy example of
 * that type using the <JsonEditor> component. Inside a modal, the PolicyExample component must be wrapped in a conditional
 * (example below), otherwise the <JsonEditor> value won't render until it's focused.
 *
 * @example
 *  <PolicyExample @policyType="acl" @container="#search-select-modal" />
 *
 * @param {string} policyType - policy type to decide which template to render; can either be "acl" or "rgp"
 * @param {string} container - selector for the container the example renders inside, passed to the copy button in JsonEditor
 */

export default class PolicyExampleComponent extends Component {
  // formatting here is purposeful so that whitespace renders correctly in JsonEditor
  policyTemplates = {
    acl: `
# Grant 'create', 'read' , 'update', and ‘list’ permission
# to paths prefixed by 'secret/*'
path "secret/*" {
  capabilities = [ "create", "read", "update", "list" ]
}

# Even though we allowed secret/*, this line explicitly denies
# secret/super-secret. This takes precedence.
path "secret/super-secret" {
  capabilities = ["deny"]
}
`,
    rgp: `
# Import strings library that exposes common string operations
import "strings"

# Conditional rule (precond) checks the incoming request endpoint
# targeted to sys/policies/acl/admin
precond = rule {
    strings.has_prefix(request.path, "sys/policies/admin")
}

# Vault checks to see if the request was made by an entity
# named James Thomas or Team Lead role defined as its metadata
main = rule when precond {
    identity.entity.metadata.role is "Team Lead" or
      identity.entity.name is "James Thomas"
}
`,
    egp: `
import "time"

# Expect requests to only happen during work days (Monday
# through Friday) 0 for Sunday and 6 for Saturday
workdays = rule {
    time.now.weekday > 0 and time.now.weekday < 6
}

# Expect requests to only happen during work hours (7:00 am -
# 6:00 pm)
workhours = rule {
    time.now.hour > 7 and time.now.hour < 18
}
main = rule {
    workdays and workhours
}
`,
  };
  moreInformationLinks = {
    acl: '/vault/docs/concepts/policies#capabilities',
    rgp: '/vault/tutorials/policies/sentinel#role-governing-policies-rgps',
    egp: '/vault/docs/enterprise/sentinel#endpoint-governing-policies-egps',
  };
}
