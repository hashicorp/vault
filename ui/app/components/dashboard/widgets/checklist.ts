/**
 * Copyright IBM Corp. 2016, 2026
 * SPDX-License-Identifier: BUSL-1.1
 */

import { action } from '@ember/object';
import { service } from '@ember/service';
import Component from '@glimmer/component';
import { tracked } from '@glimmer/tracking';
import type ChecklistStateService from 'vault/services/checklist-state';
import type { ChecklistConfig, ChecklistStepConfig, StepCta } from 'vault/utils/constants/checklist';
import { getChecklistStepConfig } from 'vault/utils/constants/checklist';

interface StepDisplay {
  id: string;
  label: string;
  description: string;
  ctas: StepCta[];
  isComplete: boolean;
  /** True when the step requires user action rather than automatic detection. */
  isUserExplicit: boolean;
}

interface Args {
  /** Full checklist configuration used to render labels, CTAs, and copy. */
  checklist: ChecklistConfig;
  /** The checklist ID (e.g. 'cluster-startup'). */
  checklistId: string;
  /** Ordered list of step IDs visible to the current user. */
  visibleSteps: string[];
  /** Called when the user clicks "Hide guide". */
  onHide: () => void;
  /** Called after a step completion state is changed. */
  onStepCompletionChange?: () => void;
}

/**
 * @module DashboardWidgetsChecklist
 * Renders the onboarding checklist step list with progress for the active lifecycle state.
 */
export default class DashboardWidgetsChecklist extends Component<Args> {
  @service('checklist-state') declare readonly checklistState: ChecklistStateService;
  @tracked expandedStepId: string | null | undefined; // internal accordion state

  private getStepConfig(stepId: string): ChecklistStepConfig | null {
    return getChecklistStepConfig(this.args.checklistId, stepId);
  }

  /** Steps enriched with display label, description, completion status, and interaction type. */
  get stepsWithStatus(): StepDisplay[] {
    return this.args.visibleSteps.map((id) => {
      const step = this.getStepConfig(id);
      return {
        id,
        label: step?.label ?? id,
        description: step?.description ?? '',
        ctas: step?.ctas ?? ([] as StepCta[]),
        isComplete: this.checklistState.isStepCompleted(this.args.checklistId, id),
        isUserExplicit: step?.completion === 'explicit',
      };
    });
  }

  get completedCount(): number {
    return this.stepsWithStatus.filter((s) => s.isComplete).length;
  }

  get visibleStepCount(): number {
    return this.args.visibleSteps.length;
  }

  get progressPercent(): number {
    if (!this.visibleStepCount) return 0;
    return Math.round((this.completedCount / this.visibleStepCount) * 100);
  }

  /** The ID of the first step that is not yet complete — auto-expanded on render. */
  get nextIncompleteActionableStepId(): string | null {
    return this.stepsWithStatus.find((s) => !s.isComplete)?.id ?? null;
  }

  get firstOpenableStepId(): string | null {
    return this.stepsWithStatus[0]?.id ?? null;
  }

  /** Single controlled open state for accordion items. */
  get openStepId(): string | null {
    if (this.expandedStepId === undefined) {
      return this.nextIncompleteActionableStepId ?? this.firstOpenableStepId;
    }

    if (this.expandedStepId) {
      const step = this.stepsWithStatus.find((s) => s.id === this.expandedStepId);
      if (step) {
        return this.expandedStepId;
      }
    }

    return this.nextIncompleteActionableStepId ?? this.firstOpenableStepId;
  }

  private getNextActionableStepId(currentStepId: string): string | null {
    const steps = this.stepsWithStatus;
    const currentIndex = steps.findIndex((step) => step.id === currentStepId);

    if (currentIndex === -1) {
      return this.nextIncompleteActionableStepId;
    }

    const nextIncomplete = steps.slice(currentIndex + 1).find((step) => !step.isComplete);

    return nextIncomplete?.id ?? this.nextIncompleteActionableStepId;
  }

  @action
  async toggleStepCompletion(stepId: string, isComplete: boolean) {
    if (isComplete) {
      await this.checklistState.updateStep.perform(this.args.checklistId, stepId, false);
      this.args.onStepCompletionChange?.();
      this.expandedStepId = stepId;
      return;
    }

    await this.checklistState.markComplete.perform(this.args.checklistId, stepId);
    this.args.onStepCompletionChange?.();
    this.expandedStepId = this.getNextActionableStepId(stepId);
  }

  @action
  handleItemToggle(stepId: string) {
    this.expandedStepId = stepId;
  }
}
