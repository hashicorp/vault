/**
 * Copyright IBM Corp. 2016, 2025
 * SPDX-License-Identifier: BUSL-1.1
 */

import Component from '@glimmer/component';
import { action } from '@ember/object';

interface Args {
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
   * Callback to update the current step when navigating backwards or
   * forwards through the wizard
   */
  onStepChange: CallableFunction;
  /**
   * Whether the current step allows proceeding to the next step
   */
  canProceed?: boolean;
  /**
   * State tracked across steps.
   */
  wizardState?: unknown;
  /**
   * Callback to update state tracked across steps.
   */
  updateWizardState?: CallableFunction;
}

export default class GuidedSetup extends Component<Args> {
  get isFinalStep() {
    return this.args.currentStep === this.args.steps.length - 1;
  }

  @action
  onStepChange(change: number) {
    const { currentStep, onStepChange } = this.args;
    const target = currentStep + change;
    onStepChange(target);
  }

  @action
  onNavStepChange(_event: Event, stepIndex: number) {
    const { onStepChange } = this.args;
    onStepChange(stepIndex);
  }
}
