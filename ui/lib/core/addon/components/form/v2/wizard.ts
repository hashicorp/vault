/**
 * Copyright IBM Corp. 2016, 2025
 * SPDX-License-Identifier: BUSL-1.1
 */

import { action } from '@ember/object';
import { service } from '@ember/service';
import Component from '@glimmer/component';
import { tracked } from '@glimmer/tracking';
import type { FormConfig, WizardConfig, WizardState, WizardStepState } from 'vault/forms/v2/form-config';
import V2Form from 'vault/forms/v2/v2-form';
import type ApiService from 'vault/services/api';

interface Args {
  config: WizardConfig;
  onCancel: () => void;
  onSuccess: () => void;
}

/**
 * Form::V2::Wizard manages multi-step wizard flows with cross-step data sharing.
 * Delegates form submission to Form::V2 components via callback pattern.
 *
 * Usage:
 * ```handlebars
 * <Form::V2::Wizard @config={{this.wizardConfig}} />
 * ```
 *
 * Features:
 * - Sequential step navigation
 * - Cross-step data flow via wizardState
 * - Dynamic payload resolution (supports functions)
 * - Delegates submission/validation to Form::V2 components
 */
export default class FormV2Wizard extends Component<Args> {
  @service declare readonly api: ApiService;
  @tracked currentStepIndex = 0;
  @tracked wizardState: WizardState = {};

  // Cache the form for the current step (cleared on navigation)
  #currentFormCache?: V2Form<any, any>;

  get config(): WizardConfig {
    return this.args.config;
  }

  get steps() {
    return this.config.steps;
  }

  get currentStep() {
    return this.steps[this.currentStepIndex];
  }

  get currentStepName() {
    return this.currentStep?.name;
  }

  /**
   * Get or create the V2Form instance for the current step.
   * Resolves dynamic payloads using current wizard state.
   * Cached to preserve field changes during a single step.
   */
  get currentForm(): V2Form<any, any> {
    // Return cached form if it exists for current step
    if (this.#currentFormCache) {
      return this.#currentFormCache;
    }

    // Special case: apply step uses the last real step's form
    // When isApplyingChanges is true, currentStepIndex === steps.length,
    // so we need to use the last actual step instead
    const step = this.isApplyingChanges
      ? this.steps[this.steps.length - 1]
      : this.steps[this.currentStepIndex];

    const resolvedPayload = this.#resolvePayload(step?.formConfig.payload);

    const resolvedConfig = {
      ...step?.formConfig,
      payload: resolvedPayload,
    } as FormConfig<any, any>;

    const form = new V2Form<any, any>(resolvedConfig);
    // eslint-disable-next-line ember/no-side-effects
    this.#currentFormCache = form;

    return form;
  }

  get currentStepState(): WizardStepState | undefined {
    return this.currentStepName ? this.wizardState[this.currentStepName] : undefined;
  }

  get stepCount() {
    return this.config.applyChanges ? this.steps.length + 1 : this.steps.length;
  }
  get isFirstStep() {
    return this.currentStepIndex === 0;
  }

  get isLastStep() {
    return this.currentStepIndex === this.stepCount - 1;
  }

  get isApplyingChanges() {
    return this.isLastStep && this.config.applyChanges;
  }

  get canAdvance() {
    // Can advance if current step has a response (successfully completed)
    return !!this.currentStepState?.response;
  }

  /**
   * Resolves a payload that might be a function or a static object.
   *
   * Function payloads enable cross-step data sharing in wizards by reading
   * from wizardState to pre-populate form fields based on previous steps.
   *
   * Example: A mount path entered in step 1 can be reused in steps 2 and 3
   * by defining their payloads as functions that read from wizardState:
   *
   * ```typescript
   * payload: (wizardState) => ({
   *   mount_path: wizardState.step1?.payload?.path || 'default/'
   * })
   * ```
   *
   * This ensures consistency across steps and improves UX by avoiding
   * repetitive data entry.
   */
  #resolvePayload(payload: any): any {
    if (typeof payload === 'function') {
      return payload(this.wizardState);
    }

    return payload;
  }

  /**
   * Updates wizard state after a step submission.
   * Stores only data - execution state is derived from task properties.
   */
  #updateWizardState(stepName: string, payload: any, response: any, error?: string) {
    this.wizardState = {
      ...this.wizardState,
      [stepName]: {
        payload,
        response,
        error,
      },
    };
  }

  /**
   * Handles successful step submission from Form::V2 component.
   * Updates wizard state and advances to next step.
   */
  @action
  onStepSuccess(response: unknown) {
    const stepName = this.currentStepName || '';
    const payload = this.currentForm.payload;

    // Update wizard state with response
    this.#updateWizardState(stepName, payload, response);

    // Auto-advance if not last step
    if (!this.isLastStep) {
      this.nextStep();
    }
  }

  /**
   * Handles failed step submission from Form::V2 component.
   * Updates wizard state with error message.
   */
  @action
  onStepError(errorMessage: string) {
    const stepName = this.currentStepName || '';

    // Store error in wizard state
    this.#updateWizardState(stepName, this.currentForm.payload, null, errorMessage);
  }

  @action
  nextStep() {
    if (this.currentStepIndex < this.stepCount - 1) {
      this.currentStepIndex++;
      // Clear form cache so next step gets fresh form with resolved payload
      this.#currentFormCache = undefined;
    }
  }

  @action
  previousStep() {
    if (this.currentStepIndex > 0) {
      this.currentStepIndex--;
      // Clear form cache when navigating back
      this.#currentFormCache = undefined;
    }
  }

  @action
  onStepChange(_event: Event, stepIndex: number) {
    // Only allow navigating to completed steps or previous steps
    const targetStep = this.steps[stepIndex];
    const targetStepState = this.wizardState[targetStep?.name || ''];

    // Can navigate if:
    // 1. Step is already completed (has response), OR
    // 2. It's a previous step (allow going back)
    if (targetStepState?.response || stepIndex <= this.currentStepIndex) {
      this.currentStepIndex = stepIndex;
      // Clear form cache to re-resolve payload with current wizard state
      this.#currentFormCache = undefined;
    }
  }
}
