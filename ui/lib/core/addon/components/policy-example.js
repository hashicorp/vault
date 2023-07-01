import Component from '@glimmer/component';

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
  };
}
