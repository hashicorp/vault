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
 * @callback onCancel
 * @callback onSave
 * @param {object} model - The parent's model
 * @param {object} modelData - If @model isn't passed in, @modelData is passed to create the record
 * @param {string} onCancel - callback triggered when cancel button is clicked
 * @param {string} onSave - callback triggered when save button is clicked
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
  @tracked createdModel = null; // set by createRecord() after policyType is selected
  policyOptions = [
    { label: 'ACL Policy', value: 'acl', isDisabled: false },
    { label: 'Role Governing Policy', value: 'rgp', isDisabled: !this.version.hasSentinel },
  ];
  // formatting here is purposeful so that whitespace renders correctly in JsonEditor
  policyTemplates = {
    acl: `
# Grant 'create', 'read' , 'update', and ‘list’ permission to paths prefixed by 'secret/*'
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

# Conditional rule (precond) checks the incoming request endpoint targeted to sys/policies/acl/admin
precond = rule {
    strings.has_prefix(request.path, "sys/policies/admin")
}

# Vault checks to see if the request was made be an entity named James Thomas 
# or the 'Team Lead' role defined as its metadata
main = rule when precond {
    identity.entity.metadata.role is "Team Lead" or
      identity.entity.name is "James Thomas"
}
`,
  };

  get model() {
    // the SS + modal form receives @modelData instead of @model
    return this.args.model ? this.args.model : this.createdModel;
  }

  @task
  *save(event) {
    event.preventDefault();
    try {
      const { isNew, name, policyType } = this.model;
      yield this.model.save();
      this.flashMessages.success(
        `${policyType.toUpperCase()} policy "${name}" was successfully ${isNew ? 'created' : 'updated'}.`
      );
      if (this.wizard.featureState === 'create') {
        this.wizard.transitionFeatureMachine('create', 'CONTINUE', policyType);
      }

      // this form is sometimes used in modal, passing the model notifies
      // the parent if the save was successful
      this.args.onSave(this.model);
    } catch (error) {
      const message = error.errors ? error.errors.join('. ') : error.message;
      this.errorBanner = message;
    }
    this.cleanup();
  }

  @action
  setModelName({ target }) {
    this.model.name = target.value.toLowerCase();
  }

  @action
  async setPolicyType(type) {
    if (this.createdModel) this.cleanup();
    this.createdModel = await this.store.createRecord(`policy/${type}`, {});
    this.createdModel.name = this.args.modelData.name;
  }

  @action
  setPolicyFromFile(index, fileInfo) {
    const { value, fileName } = fileInfo;
    this.model.policy = value;
    if (!this.model.name) {
      const trimmedFileName = trimRight(fileName, ['.json', '.txt', '.hcl', '.policy']);
      this.model.name = trimmedFileName.toLowerCase();
    }
    this.showFileUpload = false;
  }

  @action
  cancel() {
    this.cleanup();
    this.args.onCancel();
  }

  cleanup() {
    const method = this.model.isNew ? 'unloadRecord' : 'rollbackAttributes';
    this.model[method]();
    if (this.createdModel) this.createdModel = null;
  }
}
