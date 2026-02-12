/**
 * Copyright IBM Corp. 2016, 2025
 * SPDX-License-Identifier: BUSL-1.1
 */

import { service } from '@ember/service';
import Component from '@glimmer/component';
import type WizardService from 'vault/services/wizard';

interface Args {
  /**
   * The unique identifier for the wizard used for handling wizard dismissal and intro visibility state
   */
  wizardId: string;
  /**
   * Whether the intro page is in the default view or in modal view depending on how it is triggered
   */
  isModal: boolean;
  /**
   * Title of the wizard
   */
  title: string;
  /**
   * Whether the current step allows proceeding to the next step
   */
  canProceed: boolean;
  /**
   * The active step. Steps are zero-indexed.
   */
  currentStep: number;
  /**
   * Define step information to be shown in the Stepper Nav
   */
  steps: { title: string; description?: string }[];
  /**
   * Callback to update viewing state when the wizard is exited.
   */
  onDismiss: CallableFunction;
  /**
   * Whether the current step allows proceeding to the next step
   */
  onStepChange: CallableFunction;
  /**
   * State tracked across steps.
   */
  wizardState: unknown;
  /**
   * Callback to update state tracked across steps.
   */
  updateWizardState: CallableFunction;
}

export default class WizardComponent extends Component<Args> {
  @service declare readonly wizard: WizardService;

  get isIntroVisible(): boolean {
    // If wizardId is provided, use the wizard service to check intro visibility
    return this.wizard.isIntroVisible(this.args.wizardId);
  }
}
