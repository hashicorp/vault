/**
 * Copyright IBM Corp. 2016, 2025
 * SPDX-License-Identifier: BUSL-1.1
 */

import Component from '@glimmer/component';
import { service } from '@ember/service';
import { tracked } from '@glimmer/tracking';
import { action } from '@ember/object';
import { task } from 'ember-concurrency';
import { waitFor } from '@ember/test-waiters';
import { getRules } from 'kubernetes/utils/generated-role-rules';

import type RouterService from '@ember/routing/router-service';
import type ApiService from 'vault/services/api';
import type SecretMountPath from 'vault/services/secret-mount-path';
import type FlashMessageService from 'vault/services/flash-messages';
import type { EngineOwner, ValidationMap } from 'vault/app-types';
import type KubernetesRoleForm from 'vault/forms/secrets/kubernetes/role';
import type { HTMLElementEvent } from 'vault/forms';
import type { EditorView } from '@codemirror/view';

/**
 * @module CreateAndEditRolePage
 * CreateAndEditRolePage component is a child component for create and edit role pages.
 *
 * @param {KubernetesRoleForm} form - form class
 */

interface Args {
  form: KubernetesRoleForm;
}
interface RulesTemplate {
  id: string;
  label: string;
  rules: string;
}

export default class CreateAndEditRolePageComponent extends Component<Args> {
  @service('app-router') declare readonly router: RouterService;
  @service declare readonly api: ApiService;
  @service declare readonly secretMountPath: SecretMountPath;
  @service declare readonly flashMessages: FlashMessageService;

  @tracked showAnnotations = false;
  @tracked declare roleRulesTemplates: RulesTemplate[];
  @tracked selectedTemplateId = '';
  @tracked modelValidations: ValidationMap | null = null;
  @tracked invalidFormAlert = '';
  @tracked errorBanner = '';
  @tracked declare codemirrorEditor: EditorView;

  constructor(owner: EngineOwner, args: Args) {
    super(owner, args);
    // generated role rules are only rendered for the full object chain option
    if (this.args.form.generationPreference === 'full') {
      this.initRoleRules();
    }
    // if editing and annotations or labels exist expand the section
    const { extra_annotations, extra_labels } = this.args.form.data;
    if (extra_annotations || extra_labels) {
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
        key: 'extra_annotations',
        description: 'Attach arbitrary non-identifying metadata to objects.',
      },
      {
        type: 'labels',
        key: 'extra_labels',
        description:
          'Labels specify identifying attributes of objects that are meaningful and relevant to users.',
      },
    ];
  }

  get roleRulesHelpText() {
    const message =
      'This specifies the Role or ClusterRole rules to use when generating a role. Kubernetes documentation is';
    const link =
      '<a href="https://kubernetes.io/docs/reference/access-authn-authz/rbac/" target="_blank" rel="noopener noreferrer" class="has-text-white">available here</>';
    return `${message} ${link}.`;
  }

  @action
  initRoleRules() {
    // first check if generatedRoleRules matches one of the templates, the user may have chosen a template and not made changes
    // in this case we need to select the corresponding template in the dropdown
    // if there is no match then replace the example rules with the user defined value for no template option
    const { generated_role_rules } = this.args.form.data;
    const rulesTemplates = getRules();
    this.selectedTemplateId = '1';

    if (generated_role_rules) {
      const template = rulesTemplates.find((t) => t.rules === generated_role_rules);
      if (template) {
        this.selectedTemplateId = template.id;
      } else {
        (rulesTemplates.find((t) => t.id === '1') as RulesTemplate).rules = generated_role_rules;
      }
    }
    this.roleRulesTemplates = rulesTemplates;
  }

  @action
  resetRoleRules() {
    // Reset tracked rule templates to initial values
    this.roleRulesTemplates = getRules();
    // Make sure editor renders the reset template
    this.updateCodeMirror();
  }

  @action
  updateCodeMirror() {
    const template = this.roleRulesTemplates.find((t) => t.id === this.selectedTemplateId);
    this.codemirrorEditor.dispatch({
      changes: [
        {
          from: 0,
          to: this.codemirrorEditor.state.doc.length,
          insert: template?.rules,
        },
      ],
    });
  }

  @action
  selectTemplate(event: HTMLElementEvent<HTMLSelectElement>) {
    this.selectedTemplateId = event.target.value;
    // Dispatch the event to codemirror so the code editor updates when a template is selected
    this.updateCodeMirror();
  }

  @action
  changePreference(pref: 'basic' | 'expanded' | 'full') {
    if (pref === 'full') {
      this.initRoleRules();
    } else {
      this.selectedTemplateId = '';
    }
    this.args.form.generationPreference = pref;
  }

  save = task(
    waitFor(async (event: SubmitEvent) => {
      event.preventDefault();

      try {
        const { currentPath } = this.secretMountPath;
        const { form } = this.args;
        const { isValid, state, invalidFormMessage, data } = form.toJSON();

        if (isValid) {
          // set generated_role_rules to value of selected template
          const selectedTemplate = this.roleRulesTemplates?.find((t) => t.id === this.selectedTemplateId);
          if (selectedTemplate) {
            data.generated_role_rules = selectedTemplate.rules;
          }
          const { name, ...payload } = data;
          await this.api.secrets.kubernetesWriteRole(name as string, currentPath, payload);
          this.router.transitionTo('vault.cluster.secrets.backend.kubernetes.roles.role.details', name);
        } else {
          this.invalidFormAlert = invalidFormMessage;
          this.modelValidations = state;
        }
      } catch (error) {
        const { message } = await this.api.parseError(
          error,
          'Error saving role. Please try again or contact support'
        );
        this.errorBanner = message;
        this.invalidFormAlert = 'There was an error submitting this form.';
      }
    })
  );

  @action
  cancel() {
    this.router.transitionTo('vault.cluster.secrets.backend.kubernetes.roles');
  }
}
