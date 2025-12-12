/**
 * Copyright IBM Corp. 2016, 2025
 * SPDX-License-Identifier: BUSL-1.1
 */

import Component from '@glimmer/component';
import { action } from '@ember/object';

/**
 * @module QuickStart
 * QuickStart component holds the wizard content pages and navigation controls.
 *
 * @example
 *   <QuickStart @currentStep={{@currentStep}} @steps={{@steps}} @onStepChange={{@onStepChange}} @onDismiss={{@onDismiss}} @hasSubmitBlock={{has-block "submit"}} />
 */

interface Args {
  /**
   * The active step
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
   * Helper arg to conditionally render a custom submit button upon
   * completion of the wizard. Necessary to avoid a nested block error.
   */
  hasSubmitBlock: boolean;
}

export default class QuickStart extends Component<Args> {
  constructor(owner: unknown, args: Args) {
    super(owner, args);
  }

  get isFinalStep() {
    return this.args.currentStep === this.args.steps.length - 1;
  }

  @action
  onStepChange(change: number) {
    const { currentStep, steps, onStepChange } = this.args;
    const target = currentStep + change;

    if (target < 0 || target > steps.length - 1) {
      onStepChange(currentStep);
    } else {
      onStepChange(target);
    }
  }
}
