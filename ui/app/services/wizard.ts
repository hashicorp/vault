/**
 * Copyright IBM Corp. 2016, 2025
 * SPDX-License-Identifier: BUSL-1.1
 */

import Service from '@ember/service';
import { tracked } from '@glimmer/tracking';
import localStorage from 'vault/lib/local-storage';
import { DISMISSED_WIZARD_KEY } from 'vault/utils/constants/wizard';

import type { WizardId } from 'vault/app-types';

export interface StepConfig {
  title: string;
  component: string;
}

// Unique identifier for each wizard, used to track dismissal state and step state.
export type WizardID = string;

// Dynamic state storage for each wizard defined at the component level.
// This allows for flexible state management across different wizards without
// needing to predefine specific properties for each wizard in the service.
export type WizardState = Record<string, unknown>;

/**
 * WizardService manages the state of wizards across the application,
 * including tracking which wizards have been dismissed by the user.
 * This service provides a centralized way to check and update wizard
 * dismissal state instead of directly accessing localStorage.
 */
export default class WizardService extends Service {
  @tracked dismissedWizards: string[] = this.loadDismissedWizards();
  @tracked introVisibleState: Record<string, boolean> = {};

  /* Tracked properties for step state management */
  @tracked private stepData: Record<WizardID, WizardState> = {};
  @tracked private currentStep: Record<WizardID, number> = {};
  @tracked private steps: Record<WizardID, StepConfig[]> = {};

  /**
   * Load dismissed wizards from localStorage
   */
  private loadDismissedWizards(): string[] {
    return localStorage.getItem(DISMISSED_WIZARD_KEY) ?? [];
  }

  /**
   * Check if a specific wizard has been dismissed by the user
   * @param wizardId - The unique identifier for the wizard
   * @returns true if the wizard has been dismissed, false otherwise
   */
  isDismissed(wizardId: WizardId): boolean {
    return this.dismissedWizards.includes(wizardId);
  }

  /**
   * Mark a wizard as dismissed
   * @param wizardId - The unique identifier for the wizard to dismiss
   */
  dismiss(wizardId: WizardId): void {
    // Only add if not already dismissed
    if (!this.dismissedWizards.includes(wizardId)) {
      this.dismissedWizards = [...this.dismissedWizards, wizardId];
      localStorage.setItem(DISMISSED_WIZARD_KEY, this.dismissedWizards);
    }
  }

  /**
   * Clear the dismissed state for a specific wizard
   * @param wizardId - The unique identifier for the wizard to reset
   */
  reset(wizardId: WizardId): void {
    this.dismissedWizards = this.dismissedWizards.filter((id: string) => id !== wizardId);
    localStorage.setItem(DISMISSED_WIZARD_KEY, this.dismissedWizards);
    // Reset intro visibility when wizard is reset
    this.setIntroVisible(wizardId, true);
  }

  /**
   * Clear all dismissed wizard states
   */
  resetAll(): void {
    this.dismissedWizards = [];
    localStorage.removeItem(DISMISSED_WIZARD_KEY);
    this.introVisibleState = {};
  }

  /**
   * Check if the intro is visible for a specific wizard
   * @param wizardId - The unique identifier for the wizard
   * @returns true if the intro is visible, false otherwise (defaults to true if wizard not dismissed, false if dismissed)
   */
  isIntroVisible(wizardId: WizardId): boolean {
    // If intro visibility has been explicitly set, use that value
    if (this.introVisibleState[wizardId] !== undefined) {
      return this.introVisibleState[wizardId];
    }
    // Otherwise, default to true if wizard is not dismissed (first time showing)
    // and false if wizard is dismissed
    return !this.isDismissed(wizardId);
  }

  /**
   * Set the intro visibility state for a specific wizard
   * @param wizardId - The unique identifier for the wizard
   * @param visible - Whether the intro should be visible
   */
  setIntroVisible(wizardId: WizardId, visible: boolean): void {
    this.introVisibleState = {
      ...this.introVisibleState,
      [wizardId]: visible,
    };
  }

  /* Step state management */

  /**
   * Retrieve the stored state for a wizard, typed by the caller.
   * Returns an empty object if no state has been set yet, so consumers
   * should merge with their own defaults when a fully-typed object is required.
   * @param wizardId - The unique identifier for the wizard
   * @returns The wizard's current state, cast to the caller-supplied type
   */
  getState<T extends object>(wizardId: WizardId): T {
    return (this.stepData[wizardId] ?? {}) as T;
  }

  /**
   * Immutably update a single key in the wizard's stored state.
   * Other keys in the state are left unchanged.
   * @param wizardId - The unique identifier for the wizard
   * @param key - The state key to update
   * @param value - The new value for that key
   */
  updateState(wizardId: WizardId, key: string, value: unknown): void {
    this.stepData = {
      ...this.stepData,
      [wizardId]: { ...this.stepData[wizardId], [key]: value },
    };
  }

  /**
   * Reset the wizard's step data to an empty object and return navigation to step 0.
   * The step configuration (registered via setSteps) is intentionally preserved so
   * the same wizard can be re-entered without re-registering its steps.
   * @param wizardId - The unique identifier for the wizard
   */
  clearWizardState(wizardId: WizardId): void {
    this.stepData = { ...this.stepData, [wizardId]: {} };
    this.currentStep = { ...this.currentStep, [wizardId]: 0 };
  }

  /**
   * Return the index of the currently active step for a wizard.
   * Defaults to 0 if the wizard has not yet navigated to any step.
   * @param wizardId - The unique identifier for the wizard
   * @returns The zero-indexed current step number
   */
  getCurrentStep(wizardId: WizardId): number {
    return this.currentStep[wizardId] ?? 0;
  }

  /**
   * Set the active step index for a wizard.
   * @param wizardId - The unique identifier for the wizard
   * @param step - The zero-indexed step to navigate to
   */
  setCurrentStep(wizardId: WizardId, step: number): void {
    this.currentStep = { ...this.currentStep, [wizardId]: step };
  }

  /**
   * Return the registered step configuration array for a wizard.
   * Returns an empty array if no steps have been registered yet, allowing
   * consumers to supply their own defaults on first render.
   * @param wizardId - The unique identifier for the wizard
   * @returns Array of step configuration objects
   */
  getSteps(wizardId: WizardId): StepConfig[] {
    return this.steps[wizardId] ?? [];
  }

  /**
   * Replace the step configuration array for a wizard.
   * Use this to add, remove, or reorder steps dynamically — for example,
   * skipping a step based on a user's earlier selection.
   * @param wizardId - The unique identifier for the wizard
   * @param steps - The new step configuration to register
   */
  setSteps(wizardId: WizardId, steps: StepConfig[]): void {
    this.steps = { ...this.steps, [wizardId]: steps };
  }

  /**
   * Return true when the current step is the last one in the registered configuration.
   * Returns false if no steps have been registered yet.
   * @param wizardId - The unique identifier for the wizard
   * @returns Whether the wizard is on its final step
   */
  isFinalStep(wizardId: WizardId): boolean {
    const steps = this.getSteps(wizardId);
    return steps.length > 0 && this.getCurrentStep(wizardId) === steps.length - 1;
  }
}
