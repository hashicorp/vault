/**
 * Copyright IBM Corp. 2016, 2025
 * SPDX-License-Identifier: BUSL-1.1
 */

import Component from '@glimmer/component';
import { tracked } from '@glimmer/tracking';

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

// each wizard implementation can track whether the user has already dismissed the wizard via local storage
export const DISMISSED_WIZARD_KEY = 'dismissed-wizards';

export default class Wizard extends Component<Args> {
  @tracked showWelcome = true;
}
