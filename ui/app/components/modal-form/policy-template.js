import Component from '@glimmer/component';
import { action } from '@ember/object';
import { inject as service } from '@ember/service';
import { tracked } from '@glimmer/tracking';

/**
 * @module ModalForm::PolicyTemplate
 * ModalForm::PolicyTemplate components are meant to render within a modal for creating a new policy of unknown type.
 *
 * @example
 *  <ModalForm::PolicyTemplate
 *    @nameInput="new-item-name"
 *    @onSave={{this.closeModal}}
 *    @onCancel={{this.closeModal}}
 *  />
 * ```
 * @callback onCancel - callback triggered when cancel button is clicked
 * @callback onSave - callback triggered when save button is clicked
 * @param {string} nameInput - the name of the newly created policy
 */

export default class PolicyTemplate extends Component {
  @service store;
  @service version;

  @tracked policy = null; // model record passed to policy-form
  @tracked showExamplePolicy = false;

  get policyOptions() {
    return [
      { label: 'ACL Policy', value: 'acl', isDisabled: false },
      { label: 'Role Governing Policy', value: 'rgp', isDisabled: !this.version.hasSentinel },
    ];
  }
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

  @action
  setPolicyType(type) {
    if (this.policy) this.policy.unloadRecord(); // if user selects a different type, clear from store before creating a new record
    // Create form model once type is chosen
    this.policy = this.store.createRecord(`policy/${type}`, { name: this.args.nameInput });
  }

  @action
  onSave(policyModel) {
    this.args.onSave(policyModel);
    // Reset component policy for next use
    this.policy = null;
  }
}
