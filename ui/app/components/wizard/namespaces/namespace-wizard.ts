/**
 * Copyright IBM Corp. 2016, 2025
 * SPDX-License-Identifier: BUSL-1.1
 */

import { service } from '@ember/service';
import { action } from '@ember/object';
import Component from '@glimmer/component';
import { SecurityPolicy } from 'vault/components/wizard/namespaces/step-1';
import { CreationMethod } from 'vault/utils/constants/snippet';
import { WIZARD_ID_MAP } from 'vault/utils/constants/wizard';
import { INTRO_NAMESPACES_CTA_CLICKED, NAMESPACE_CREATED } from 'vault/utils/analytic-events';

import type ApiService from 'vault/services/api';
import type Block from 'vault/components/wizard/namespaces/step-2';
import type FlashMessageService from 'vault/services/flash-messages';
import type NamespaceService from 'vault/services/namespace';
import type RouterService from '@ember/routing/router-service';
import type WizardService from 'vault/services/wizard';
import type AnalyticsService from 'vault/services/analytics';
import type { StepConfig } from 'vault/services/wizard';

const DEFAULT_STEPS: StepConfig[] = [
  { title: 'Select setup', component: 'wizard/namespaces/step-1' },
  { title: 'Map out namespaces', component: 'wizard/namespaces/step-2' },
  { title: 'Apply changes', component: 'wizard/namespaces/step-3' },
];

interface Args {
  isIntroModal: boolean;
  onRefresh: CallableFunction;
  onFlexiblePolicyComplete: CallableFunction;
}

interface WizardState {
  securityPolicyChoice: SecurityPolicy | null;
  namespacePaths: string[] | null;
  namespaceBlocks: Block[] | null;
  creationMethod: CreationMethod | null;
  codeSnippet: string | null;
}

const DEFAULT_WIZARD_STATE: WizardState = {
  securityPolicyChoice: null,
  namespacePaths: null,
  namespaceBlocks: null,
  creationMethod: null,
  codeSnippet: null,
};

export default class WizardNamespacesWizardComponent extends Component<Args> {
  @service declare readonly api: ApiService;
  @service declare readonly router: RouterService;
  @service declare readonly flashMessages: FlashMessageService;
  @service declare readonly wizard: WizardService;
  @service declare readonly analytics: AnalyticsService;
  @service declare namespace: NamespaceService;

  methods = CreationMethod;
  policy = SecurityPolicy;

  wizardId = WIZARD_ID_MAP.namespace;

  get currentStep() {
    return this.wizard.getCurrentStep(this.wizardId);
  }

  get steps() {
    const steps = this.wizard.getSteps(this.wizardId);
    return steps.length > 0 ? steps : DEFAULT_STEPS;
  }

  get wizardState(): WizardState {
    return { ...DEFAULT_WIZARD_STATE, ...this.wizard.getState<Partial<WizardState>>(this.wizardId) };
  }

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

  get isFinalStep() {
    return this.wizard.isFinalStep(this.wizardId);
  }

  get shouldShowExitButton() {
    // Show exit button unless we're on the final step with UI creation method
    return !(this.wizardState.creationMethod === CreationMethod.UI && this.isFinalStep);
  }

  get exitText() {
    return this.isFinalStep && this.wizardState.securityPolicyChoice === SecurityPolicy.STRICT
      ? 'Done & Exit'
      : 'Exit';
  }

  updateSteps() {
    if (this.wizardState.securityPolicyChoice === SecurityPolicy.FLEXIBLE) {
      this.wizard.setSteps(this.wizardId, [
        { title: 'Select setup', component: 'wizard/namespaces/step-1' },
        { title: 'Apply changes', component: 'wizard/namespaces/step-3' },
      ]);
    } else {
      this.wizard.setSteps(this.wizardId, DEFAULT_STEPS);
    }
  }

  @action
  onStepChange(step: number) {
    this.wizard.setCurrentStep(this.wizardId, step);
    // if user policy selection changes which steps we show, update upon page navigation
    // instead of flashing the changes when toggling
    this.updateSteps();
  }

  @action
  updateWizardState(key: string, value: unknown) {
    this.wizard.updateState(this.wizardId, key, value);
  }

  @action
  async onDone() {
    await this.onDismiss({ trackExit: false });
    this.args.onFlexiblePolicyComplete();
    this.flashMessages.success(`Your current setup is 1 namespace.`, { title: 'Guided start complete' });
  }

  @action
  async onDismiss({ trackExit = true }: { trackExit?: boolean } = {}) {
    if (trackExit) {
      const isOnIntro = this.wizard.isIntroVisible(this.wizardId);
      const CTA = isOnIntro ? (this.args.isIntroModal ? 'Close' : 'Skip') : 'Exit';
      const location = isOnIntro ? 'intro-page' : 'wizard';
      this.trackCtaEvent(CTA, location, 'dismissed', 'intro-dismiss-button');
    }
    this.wizard.dismiss(this.wizardId);
    this.wizard.clearWizardState(this.wizardId);
    await this.args.onRefresh();
  }

  @action
  trackClickEvent(cta: string) {
    this.trackCtaEvent(cta, 'intro', 'clicked', 'intro-cta-button');
  }

  // `variation` distinguishes the modal from the full-page intro, which is
  // otherwise only implied by the CTA label.
  private trackCtaEvent(CTA: string, location: string, action: string, uiElement: string) {
    this.analytics.trackEvent(INTRO_NAMESPACES_CTA_CLICKED, {
      CTA,
      channel: 'webpage',
      location,
      objectType: 'namespace',
      variation: this.args.isIntroModal ? 'modal' : 'page',
      uiElement,
      type: 'Button',
      action,
    });
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
  onIntroChange(visible: boolean) {
    // Hiding the intro here means the user clicked "Guided start"
    if (!visible) {
      this.trackClickEvent('Guided start');
    }
    this.wizard.setIntroVisible(this.wizardId, visible);
  }

  // Namespaces have no subtype, so `object` carries the security policy choice
  // (strict/flexible) the user selected in the wizard.
  private trackNamespaceCreationEvent(quantity: number, successFlag: boolean) {
    this.analytics.trackEvent(NAMESPACE_CREATED, {
      objectType: 'namespace',
      object: this.wizardState.securityPolicyChoice ?? undefined,
      process: 'UI',
      quantity,
      successFlag,
    });
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

      this.trackNamespaceCreationEvent(namespacePaths.length, true);

      this.flashMessages.success('Your new configuration has been applied.', { title: 'Namespaces created' });
    } catch (error) {
      this.trackNamespaceCreationEvent(this.wizardState.namespacePaths?.length ?? 0, false);

      const { message } = await this.api.parseError(error);
      this.flashMessages.danger(`Error creating namespaces: ${message}`);
    } finally {
      this.onDismiss({ trackExit: false });
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
