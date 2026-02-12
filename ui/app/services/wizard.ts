/**
 * Copyright IBM Corp. 2016, 2025
 * SPDX-License-Identifier: BUSL-1.1
 */

import Service from '@ember/service';
import { tracked } from '@glimmer/tracking';
import localStorage from 'vault/lib/local-storage';

const DISMISSED_WIZARD_KEY = 'dismissed-wizards';

/**
 * WizardService manages the state of wizards across the application,
 * particularly tracking which wizards have been dismissed by the user.
 * This service provides a centralized way to check and update wizard
 * dismissal state instead of directly accessing localStorage.
 */
export default class WizardService extends Service {
  @tracked dismissedWizards: string[] = this.loadDismissedWizards();
  @tracked introVisibleState: Record<string, boolean> = {};

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
  isDismissed(wizardId: string): boolean {
    return this.dismissedWizards.includes(wizardId);
  }

  /**
   * Mark a wizard as dismissed
   * @param wizardId - The unique identifier for the wizard to dismiss
   */
  dismiss(wizardId: string): void {
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
  reset(wizardId: string): void {
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
  isIntroVisible(wizardId: string): boolean {
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
  setIntroVisible(wizardId: string, visible: boolean): void {
    this.introVisibleState = {
      ...this.introVisibleState,
      [wizardId]: visible,
    };
  }
}
