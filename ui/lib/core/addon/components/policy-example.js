/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: MPL-2.0
 */

import Component from '@glimmer/component';

/**
 * @module PolicyExample
 * PolicyExample component is meant to render within a PolicyForm component to show an example of a policy.
 * edit this *** 
 * PolicyExample does not render in a modal, as it is wrapped in a conditional within PolicyForm.
 * 
 *
 * @example
 *  <PolicyExample 
 *    @policyType={{@model.policyType}} 
 *  />
 * 
 * @example (in modal)
 *  <Modal
 *    @onClose={{fn (mut this.showTemplateModal) false}}
 *    @isActive={{this.showTemplateModal}}
 *  >
 *    <section class="modal-card-body">
 *      {{! code-mirror modifier does not render value initially until focus event fires }}
 *      {{! wait until the Modal is rendered and then show the PolicyExample (contains JsonEditor) }}
        {{#if this.showTemplateModal}}
          <PolicyExample @policyType={{@model.policyType}}/>
        {{/if}}
      </section>
      <div class="modal-card-head has-border-top-light">
        <button type="button" class="button" {{on "click" (fn (mut this.showTemplateModal) false)}} data-test-close-modal>
          Close
        </button>
      </div>
    </Modal>
 * ```
 * talk about getter, policyTemplate, policyType ?
 * @param {string} policyType - policy type to decide which template to render; can either be "acl" or "rgp"
 * 
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
  };
}
