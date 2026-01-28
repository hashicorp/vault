/**
 * Copyright IBM Corp. 2016, 2025
 * SPDX-License-Identifier: BUSL-1.1
 */

import { service } from '@ember/service';
import { action } from '@ember/object';
import { tracked } from '@glimmer/tracking';
import Component from '@glimmer/component';
import localStorage from 'vault/lib/local-storage';
import { SecurityPolicy } from 'vault/components/wizard/namespaces/step-1';
import { CreationMethod } from 'vault/components/wizard/namespaces/step-3';
import { DISMISSED_WIZARD_KEY } from 'vault/components/wizard';

import type ApiService from 'vault/services/api';
import type Block from 'vault/components/wizard/namespaces/step-2';
import type FlashMessageService from 'vault/services/flash-messages';
import type NamespaceService from 'vault/services/namespace';
import type RouterService from '@ember/routing/router-service';

const DEFAULT_STEPS = [
  { title: 'Select setup', component: 'wizard/namespaces/step-1' },
  { title: 'Map out namespaces', component: 'wizard/namespaces/step-2' },
  { title: 'Apply changes', component: 'wizard/namespaces/step-3' },
];

interface Args {
  onDismiss: CallableFunction;
  onRefresh: CallableFunction;
}

interface WizardState {
  securityPolicyChoice: SecurityPolicy | null;
  namespacePaths: string[] | null;
  namespaceBlocks: Block[] | null;
  creationMethod: CreationMethod | null;
  codeSnippet: string | null;
}

export default class WizardNamespacesWizardComponent extends Component<Args> {
  @service declare readonly api: ApiService;
  @service declare readonly router: RouterService;
  @service declare readonly flashMessages: FlashMessageService;
  @service declare namespace: NamespaceService;

  @tracked steps = DEFAULT_STEPS;
  @tracked wizardState: WizardState = {
    securityPolicyChoice: null,
    namespacePaths: null,
    namespaceBlocks: null,
    creationMethod: null,
    codeSnippet: null,
  };
  @tracked currentStep = 0;

  methods = CreationMethod;
  policy = SecurityPolicy;

  wizardId = 'namespace';

  // Whether the current step requirements have been met to proceed to the next step
  get canProceed() {
    switch (this.currentStep) {
      case 0: // Step 1 - requires security policy choice
        return Boolean(this.wizardState.securityPolicyChoice);
      case 1: // Step 2 - requires valid namespace inputs
        return Boolean(this.wizardState.namespacePaths);
      case 2: // Step 3 - no validation is needed
        return true;
      default:
        return true;
    }
  }

  get exitText() {
    return this.currentStep === this.steps.length - 1 &&
      this.wizardState.securityPolicyChoice === SecurityPolicy.STRICT
      ? 'Done & Exit'
      : 'Exit';
  }

  updateSteps() {
    if (this.wizardState.securityPolicyChoice === SecurityPolicy.FLEXIBLE) {
      this.steps = [
        { title: 'Select setup', component: 'wizard/namespaces/step-1' },
        { title: 'Apply changes', component: 'wizard/namespaces/step-3' },
      ];
    } else {
      this.steps = DEFAULT_STEPS;
    }
  }

  @action
  onStepChange(step: number) {
    this.currentStep = step;
    // if user policy selection changes which steps we show, update upon page navigation
    // instead of flashing the changes when toggling
    this.updateSteps();
  }

  @action
  updateWizardState(key: string, value: unknown) {
    this.wizardState = {
      ...this.wizardState,
      [key]: value,
    };
  }

  @action
  async onSubmit() {
    switch (this.wizardState.creationMethod) {
      case CreationMethod.UI:
        await this.createNamespacesFromWizard();
        break;
      default:
        // The other creation methods require the user to execute the commands on their own
        // In these cases, there is no submit button
        break;
    }
  }

  @action
  async onDismiss() {
    const item = localStorage.getItem(DISMISSED_WIZARD_KEY) ?? [];
    localStorage.setItem(DISMISSED_WIZARD_KEY, [...item, this.wizardId]);
    await this.args.onRefresh();
    this.args.onDismiss();
  }

  @action
  async createNamespacesFromWizard() {
    try {
      const { namespacePaths } = this.wizardState;
      if (!namespacePaths) return;

      for (const nsPath of namespacePaths) {
        const parts = nsPath.split('/');
        const namespaceName = parts[parts.length - 1] as string;
        const parentPath = parts.length > 1 ? parts.slice(0, -1).join('/') : undefined;
        // this provides the full nested path for the header
        const fullPath = parentPath ? this.namespace.path + '/' + parentPath : undefined;
        await this.createNamespace(namespaceName, fullPath);
      }

      this.flashMessages.success(`The namespaces have been successfully created.`);
    } catch (error) {
      const { message } = await this.api.parseError(error);
      this.flashMessages.danger(`Error creating namespaces: ${message}`);
    } finally {
      await this.args.onRefresh();
      this.onDismiss();
    }
  }

  @action
  switchNamespace(targetNamespace: string) {
    this.router.transitionTo('vault.cluster.dashboard', {
      queryParams: { namespace: targetNamespace },
    });
  }

  async createNamespace(path: string, header?: string) {
    const headers = header ? this.api.buildHeaders({ namespace: header }) : undefined;
    await this.api.sys.systemWriteNamespacesPath(path, {}, headers);
  }
}
