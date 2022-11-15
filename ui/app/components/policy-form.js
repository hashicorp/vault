import Component from '@glimmer/component';
import { action } from '@ember/object';
import { inject as service } from '@ember/service';
import { task } from 'ember-concurrency';
import { tracked } from '@glimmer/tracking';
import trimRight from 'vault/utils/trim-right';

/**
 * @module PolicyForm
 * PolicyForm components are used to display the create and edit forms for all types of policies
 *
 * @example
 *  <PolicyForm
 *    @model={{this.model}}
 *    @onSave={{transition-to "vault.cluster.policy.show" this.model.policyType this.model.name}}
 *    @onCancel={{transition-to "vault.cluster.policies.index"}}
 *  />
 * ```
 * @callback onCancel - callback triggered when cancel button is clicked
 * @callback onSave - callback triggered when save button is clicked
 * @param {object} model - ember data model from createRecord
 * @param {boolean} [isInline=false] - true when form is rendered within a modal
 * * params when form renders within search-select-with-modal.hbs:
 * @param {object} [nameInput] - search input from SS passed as name attr when firing createSearchSelectModel callback
 * @callback createSearchSelectModel - callback to fire when new item is selected to create in SS+Modal
 */

export default class PolicyFormComponent extends Component {
  @service store;
  @service wizard;
  @service version;
  @service flashMessages;
  @tracked errorBanner;
  @tracked file = null;
  @tracked showFileUpload = false;
  @tracked showExamplePolicy = false;
  policyOptions = [
    { label: 'ACL Policy', value: 'acl', isDisabled: false },
    { label: 'Role Governing Policy', value: 'rgp', isDisabled: !this.version.hasSentinel },
  ];
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

  @task
  *save(event) {
    event.preventDefault();
    try {
      const { isNew, name, policyType } = this.args.model;
      yield this.args.model.save();
      this.flashMessages.success(
        `${policyType.toUpperCase()} policy "${name}" was successfully ${isNew ? 'created' : 'updated'}.`
      );
      if (this.wizard.featureState === 'create') {
        this.wizard.transitionFeatureMachine('create', 'CONTINUE', policyType);
      }

      // this form is sometimes used in modal, passing the model notifies
      // the parent if the save was successful
      this.args.onSave(this.args.model);
    } catch (error) {
      const message = error.errors ? error.errors.join('. ') : error.message;
      this.errorBanner = message;
    }
    this.cleanup();
  }

  @action
  setModelName({ target }) {
    this.args.model.name = target.value.toLowerCase();
  }

  @action
  setPolicyType(type) {
    // selecting a type only happens in the modal form
    // cleanup any model argument before firing parent action to create a new record in identity/edit-form.js
    if (this.args.model) this.cleanup();
    this.args.createSearchSelectModel({ type, name: this.args.nameInput });
  }

  @action
  setPolicyFromFile(index, fileInfo) {
    const { value, fileName } = fileInfo;
    this.args.model.policy = value;
    if (!this.args.model.name) {
      const trimmedFileName = trimRight(fileName, ['.json', '.txt', '.hcl', '.policy']);
      this.args.model.name = trimmedFileName.toLowerCase();
    }
    this.showFileUpload = false;
  }

  @action
  cancel() {
    this.cleanup();
    this.args.onCancel();
  }

  cleanup() {
    const method = this.args.model.isNew ? 'unloadRecord' : 'rollbackAttributes';
    this.args.model[method]();
  }
}
