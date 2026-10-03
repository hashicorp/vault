/**
 * Copyright IBM Corp. 2016, 2025
 * SPDX-License-Identifier: BUSL-1.1
 */

import Component from '@glimmer/component';
import { action } from '@ember/object';
import { service } from '@ember/service';
import { task } from 'ember-concurrency';
import trimRight from 'vault/utils/trim-right';
import { tracked } from '@glimmer/tracking';
import {
  formatStanzas,
  policySnippetArgs,
  PolicyStanza,
  PolicyTypes,
} from 'core/utils/code-generators/policy';
import errorMessage from 'vault/utils/error-message';

import type FlashMessageService from 'ember-cli-flash/services/flash-messages';
import type { HTMLElementEvent } from 'vault/forms';
import type { PolicyData } from 'core/components/code-generator/policy/builder';
import type { FormField } from 'vault/vault/app-types';

/**
 * @module PolicyForm
 * PolicyForm components are the forms to create and edit all types of policies. This is only the form, not the outlying layout, and expects that the form model is passed from the parent.
 *
 * @example
 *  <PolicyForm
 *    @model={{this.model}}
 *    @onSave={{transition-to "vault.cluster.policy.show" this.model.policyType this.model.name}}
 *    @onCancel={{transition-to "vault.cluster.policies.index"}}
 *    @isCompact={{false}}
 *  />
 * ```
 * @callback onCancel - callback triggered when cancel button is clicked
 * @callback onSave - callback triggered when save button is clicked. Passes saved model
 * @param {object} model - ember data model from createRecord
 * @param {boolean} isCompact - renders a compact version of the form component, such as when rendering in a modal (see policy-template.hbs)
 */

enum EditorTypes {
  CODE = 'code',
  VISUAL = 'visual',
}

interface PolicyModel {
  name: string;
  policy: string;
  policyType: PolicyTypes;
  isNew: boolean;
  additionalAttrs?: FormField[]; // Only exist for "rgp" and "egp" policy types
  save: () => Promise<void>;
  unloadRecord: () => void;
  rollbackAttributes: () => void;
}

interface Args {
  onCancel: () => void;
  onSave: (model: PolicyModel) => void;
  model: PolicyModel;
  isCompact?: boolean;
}

export default class PolicyFormComponent extends Component<Args> {
  @service declare readonly flashMessages: FlashMessageService;

  editTypes = { [EditorTypes.VISUAL]: 'Visual editor', [EditorTypes.CODE]: 'Code editor' } as const;

  @tracked editType: EditorTypes = EditorTypes.VISUAL;
  @tracked errorBanner = '';
  @tracked showFileUpload = false;
  @tracked showSwitchEditorsModal = false;
  @tracked showTemplateModal = false;
  @tracked stanzas: PolicyStanza[] = [new PolicyStanza()];

  constructor(owner: unknown, args: Args) {
    super(owner, args);
    // Only ACL policies support the visual editor
    this.editType = this.args.model.policyType === PolicyTypes.ACL ? EditorTypes.VISUAL : EditorTypes.CODE;
  }

  isActiveEditor = (type: string): boolean => type === this.editType;

  get hasPolicyDiff() {
    const { policy } = this.args.model;
    // Make sure policy has a value (if it's undefined, neither editor has been used)
    // Return true if there is a difference between stanzas and policy arg
    // which means the user has made changes using the code editor
    return policy && formatStanzas(this.stanzas) !== policy;
  }

  get snippetArgs() {
    const policyName = this.args.model.name || '<policy name>';
    const policy = formatStanzas(this.stanzas);
    return policySnippetArgs(policyName, policy);
  }

  get visualEditorSupported() {
    const { model, isCompact } = this.args;
    return model.isNew && model.policyType === PolicyTypes.ACL && !isCompact;
  }

  @action
  confirmEditorSwitch() {
    // User has confirmed discarding changes so switch to "visual" editor
    this.editType = EditorTypes.VISUAL;
    this.showSwitchEditorsModal = false;
    // Reset this.args.model.policy to match visual editor stanzas
    this.setPolicy(formatStanzas(this.stanzas));
  }

  @action
  handleNameInput(event: HTMLElementEvent<HTMLInputElement>) {
    const { value } = event.target;
    this.setName(value);
  }

  @action
  handlePolicyChange({ policy, stanzas }: PolicyData) {
    this.setPolicy(policy);
    this.stanzas = stanzas;
  }

  @action
  handleRadioChange(event: HTMLElementEvent<HTMLInputElement>) {
    const { value } = event.target;

    if (!Object.values(EditorTypes).includes(value as EditorTypes)) {
      console.debug(`Invalid editor type: ${value}`); // eslint-disable-line
      return;
    }

    const editorType = value as EditorTypes;

    // Users cannot make changes using the code editor and have those parsed BACK to the visual editor
    if (editorType === EditorTypes.VISUAL && this.hasPolicyDiff) {
      // Open modal to confirm user wants to switch back to "visual" editor and lose changes
      this.showSwitchEditorsModal = true;
    } else {
      this.editType = editorType;
    }
  }

  @task
  *save(event: HTMLElementEvent<HTMLFormElement>) {
    event.preventDefault();
    try {
      const { name, policyType, isNew } = this.args.model;
      yield this.args.model.save();
      this.flashMessages.success(
        `${policyType.toUpperCase()} policy "${name}" was successfully ${isNew ? 'created' : 'updated'}.`
      );
      this.args.onSave(this.args.model);
    } catch (error) {
      this.errorBanner = errorMessage(error);
    }
  }

  @action
  setName(name: string) {
    this.args.model.name = name.toLowerCase();
  }

  @action
  setPolicyFromFile(fileInfo: { value: string; filename: string }) {
    const { value, filename } = fileInfo;
    this.setPolicy(value);
    if (!this.args.model.name) {
      const trimmedFileName = trimRight(filename, ['.json', '.txt', '.hcl', '.policy']);
      this.setName(trimmedFileName);
    }
    this.showFileUpload = false;
    // Switch to the code editor if they've uploaded a policy
    this.editType = EditorTypes.CODE;
  }

  @action
  setPolicy(policy: string) {
    this.args.model.policy = policy;
  }

  @action
  cancel() {
    const method = this.args.model.isNew ? 'unloadRecord' : 'rollbackAttributes';
    this.args.model[method]();
    this.args.onCancel();
  }
}
