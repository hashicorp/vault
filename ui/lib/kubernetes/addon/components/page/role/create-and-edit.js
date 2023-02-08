import Component from '@glimmer/component';
import { inject as service } from '@ember/service';
import { tracked } from '@glimmer/tracking';
import { action } from '@ember/object';
import { task } from 'ember-concurrency';
import { waitFor } from '@ember/test-waiters';
import { getRules } from '../../../utils/generated-role-rules';
import { htmlSafe } from '@ember/template';

/**
 * @module CreateAndEditRolePage
 * CreateAndEditRolePage component is a child component for create and edit role pages.
 *
 * @param {object} model - role model that contains role record and backend
 */

export default class CreateAndEditRolePageComponent extends Component {
  @service router;
  @service flashMessages;

  @tracked roleRulesTemplates;
  @tracked selectedTemplateId;
  @tracked modelValidations;
  @tracked invalidFormAlert;
  @tracked errorBanner;

  constructor() {
    super(...arguments);
    this.initRoleRules();
    // if editing and annotations or labels exist expand the section
    const { extraAnnotations, extraLabels } = this.args.model;
    if (extraAnnotations || extraLabels) {
      this.showAnnotations = true;
    }
  }

  get generationPreferences() {
    return [
      {
        title: 'Generate token only using existing service account',
        description:
          'Enter a service account that already exists in Kubernetes and Vault will dynamically generate a token.',
        value: 'basic',
      },
      {
        title: 'Generate token, service account, and role binding objects',
        description:
          'Enter a pre-existing role (or ClusterRole) to use. Vault will generate a token, a service account and role binding objects.',
        value: 'expanded',
      },
      {
        title: 'Generate entire Kubernetes object chain',
        description:
          'Vault will generate the entire chain— a role, a token, a service account, and role binding objects— based on rules you supply.',
        value: 'full',
      },
    ];
  }

  get extraFields() {
    return [
      {
        type: 'annotations',
        key: 'extraAnnotations',
        description: 'Attach arbitrary non-identifying metadata to objects.',
      },
      {
        type: 'labels',
        key: 'extraLabels',
        description:
          'Labels specify identifying attributes of objects that are meaningful and relevant to users.',
      },
    ];
  }

  get roleRulesHelpText() {
    const message =
      'This specifies the Role or ClusterRole rules to use when generating a role. Kubernetes documentation is';
    const link =
      '<a href="https://kubernetes.io/docs/reference/access-authn-authz/rbac/" target="_blank" rel="noopener noreferrer">available here</>';
    return htmlSafe(`${message} ${link}.`);
  }

  @action
  initRoleRules() {
    // first check if generatedRoleRules matches one of the templates, the user may have chosen a template and not made changes
    // in this case we need to select the corresponding template in the dropdown
    // if there is no match then replace the example rules with the user defined value for no template option
    const { generatedRoleRules } = this.args.model;
    const rulesTemplates = getRules();
    this.selectedTemplateId = '1';

    if (generatedRoleRules) {
      const template = rulesTemplates.findBy('rules', generatedRoleRules);
      if (template) {
        this.selectedTemplateId = template.id;
      } else {
        rulesTemplates.findBy('id', '1').rules = generatedRoleRules;
      }
    }
    this.roleRulesTemplates = rulesTemplates;
  }

  @action
  resetRoleRules() {
    this.roleRulesTemplates = getRules();
  }

  @action
  selectTemplate(event) {
    this.selectedTemplateId = event.target.value;
  }

  @action
  changePreference(pref) {
    if (pref === 'full') {
      this.initRoleRules();
    } else {
      this.selectedTemplateId = null;
    }
    this.args.model.generationPreference = pref;
  }

  @task
  @waitFor
  *save() {
    try {
      // set generatedRoleRoles to value of selected template
      const selectedTemplate = this.roleRulesTemplates.findBy('id', this.selectedTemplateId);
      if (selectedTemplate) {
        this.args.model.generatedRoleRules = selectedTemplate.rules;
      }
      yield this.args.model.save();
      this.router.transitionTo(
        'vault.cluster.secrets.backend.kubernetes.roles.role.details',
        this.args.model.name
      );
    } catch (error) {
      const message = error.errors ? error.errors.join('. ') : error.message;
      this.errorBanner = message;
      this.invalidFormAlert = 'There was an error submitting this form.';
    }
  }

  @action
  async onSave(event) {
    event.preventDefault();
    const { isValid, state, invalidFormMessage } = await this.args.model.validate();
    if (isValid) {
      this.modelValidations = null;
      this.save.perform();
    } else {
      this.invalidFormAlert = invalidFormMessage;
      this.modelValidations = state;
    }
  }

  @action
  cancel() {
    const { model } = this.args;
    const method = model.isNew ? 'unloadRecord' : 'rollbackAttributes';
    model[method]();
    this.router.transitionTo('vault.cluster.secrets.backend.kubernetes.roles');
  }
}
