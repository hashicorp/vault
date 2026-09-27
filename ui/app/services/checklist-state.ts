/**
 * Copyright IBM Corp. 2016, 2026
 * SPDX-License-Identifier: BUSL-1.1
 */

import Service, { service } from '@ember/service';
import { tracked } from '@glimmer/tracking';
import { ResponseError } from '@hashicorp/vault-client-typescript';
import { keepLatestTask, task } from 'ember-concurrency';
import localStorage from 'vault/lib/local-storage';
import { CHECKLIST_LOCAL_STATE_KEY, getChecklistStepConfig } from 'vault/utils/constants/checklist';

import type { ApiResponse } from 'vault/api';
import type ApiService from 'vault/services/api';
import type PermissionsService from 'vault/services/permissions';
import type VersionService from 'vault/services/version';
import type { ChecklistLocalStateData, ChecklistView } from 'vault/utils/constants/checklist';

/**
 * Sparse map of checklist completion state: { checklistId: { stepId: true } }.
 * Step entries absent from the map are treated as incomplete (false) by consumers.
 */
export type ChecklistStateData = Record<string, Record<string, boolean>>;

/**
 * Lifecycle state for a single checklist:
 * - active: visible, not all steps complete
 * - complete: all visible steps complete, checklist not yet hidden
 * - hidden: user dismissed the checklist (regardless of completion)
 * - post-completion: all visible steps complete AND user has hidden it
 */
export type ChecklistLifecycleState = 'active' | 'complete' | 'hidden' | 'post-completion';

const CHECKLIST_STATE_PATH = '/sys/config/ui/checklist-state';

/**
 * Type guard that validates a value has the expected ChecklistStateData shape:
 * { [checklistId]: { [stepId]: boolean } }.
 * An empty object is valid (no checklists have been interacted with yet).
 */
function isChecklistStateData(value: unknown): value is ChecklistStateData {
  if (typeof value !== 'object' || value === null || Array.isArray(value)) return false;

  return Object.values(value as Record<string, unknown>).every(
    (steps) =>
      typeof steps === 'object' &&
      steps !== null &&
      !Array.isArray(steps) &&
      Object.values(steps as Record<string, unknown>).every((v) => typeof v === 'boolean')
  );
}

/**
 * Returns true for HTTP responses that should degrade gracefully rather than
 * surface as errors in the UI:
 * - 403 Forbidden (insufficient permissions)
 * - 404 Not Found (endpoint not yet deployed or CE cluster)
 * - any 5xx server error
 * - TypeError (network failure, connection refused, connection dropped)
 * - SyntaxError (non-JSON body, e.g. panic response without a proper error payload)
 */
function isNonBlockingError(error: unknown): boolean {
  if (error instanceof ResponseError) {
    const { status } = error.response;
    return status === 403 || status === 404 || status >= 500;
  }
  if (error instanceof TypeError || error instanceof SyntaxError) {
    return true;
  }
  return false;
}

/**
 * Type guard that validates a value has the expected ChecklistLocalStateData
 * shape: { [checklistId]: { hidden?: boolean, lastView?: 'checklist' | 'complete-banner' } }.
 * An empty object is valid (no checklists have local preferences yet).
 */
function isChecklistLocalStateData(value: unknown): value is ChecklistLocalStateData {
  if (typeof value !== 'object' || value === null || Array.isArray(value)) return false;

  return Object.values(value as Record<string, unknown>).every((entry) => {
    if (typeof entry !== 'object' || entry === null || Array.isArray(entry)) return false;

    const { hidden, lastView } = entry as Record<string, unknown>;
    if (hidden !== undefined && typeof hidden !== 'boolean') return false;
    if (lastView !== undefined && lastView !== 'checklist' && lastView !== 'complete-banner') return false;
    return true;
  });
}

/**
 * Manages checklist state for the onboarding checklist engine using a two-tier
 * storage model:
 *  - Completion state: read/written via the sys/config/ui/checklist-state API
 *    so that progress is shared across browsers and sessions.
 *  - Local preferences (hidden, last view shown): stored in localStorage,
 *    nested by checklist ID, so they are per-browser and never written to
 *    the backend.
 *
 * The stored completion shape is { checklistId: { stepId: true } }; a missing
 * step entry is treated as false (incomplete) by consumers.
 */
export default class ChecklistStateService extends Service {
  @service declare readonly api: ApiService;
  @service declare readonly version: VersionService;
  @service declare readonly permissions: PermissionsService;

  @tracked private _state: ChecklistStateData = {};

  /**
   * True when the last fetch succeeded. Set to false on 403/5xx so callers
   * can hide checklist UI rather than showing a broken or partially-loaded state.
   */
  @tracked isAvailable = true;

  /**
   * Browser-local checklist preferences (hidden, last view), nested by
   * checklist ID to mirror the backend's checklist_state payload shape.
   * Backed by localStorage so preferences are per-browser and are never
   * written to the backend. Existing localStorage data with an unexpected
   * shape is silently discarded.
   */
  @tracked private _localState: ChecklistLocalStateData = this._loadLocalState();

  private _getLocalStorageItem(key: string, removeMalformed = false): unknown {
    try {
      return localStorage.getItem(key);
    } catch (error) {
      if (error instanceof SyntaxError) {
        if (removeMalformed) localStorage.removeItem(key);
        return null;
      }
      throw error;
    }
  }

  private _loadLocalState(): ChecklistLocalStateData {
    const stored = this._getLocalStorageItem(CHECKLIST_LOCAL_STATE_KEY);
    return isChecklistLocalStateData(stored) ? stored : {};
  }

  /** Merges a partial local-state patch for a single checklist and persists the result. */
  private _patchLocalState(checklistId: string, patch: { hidden?: boolean; lastView?: ChecklistView }): void {
    this._localState = {
      ...this._localState,
      [checklistId]: { ...this._localState[checklistId], ...patch },
    };
    localStorage.setItem(CHECKLIST_LOCAL_STATE_KEY, this._localState);
  }

  /**
   * True when the cluster is running Vault Enterprise.
   * The onboarding checklist is an Enterprise-only feature and should not be
   * rendered in Community Edition builds.
   */
  get isEnterpriseFeature(): boolean {
    return this.version.isEnterprise;
  }

  /**
   * Entry gate for checklist visibility.
   * For phase 1 we rely on existing nav/permission visibility checks: users
   * must be able to configure at least one auth method or secrets engine.
   */
  get hasChecklistEntryAccess(): boolean {
    const canConfigureAuthMethods = this.permissions.hasNavPermission('access', 'methods');
    const canConfigureSecretsEngines = this.permissions.hasPermission('sys/mounts');
    return canConfigureAuthMethods || canConfigureSecretsEngines;
  }

  /**
   * Filters a set of step IDs to only those the current user is permitted to
   * access, based on existing nav permission checks.
   * Pass this filtered list to getChecklistProgress / isChecklistComplete so
   * that progress is computed against the user's visible steps only.
   */
  getVisibleSteps(checklistId: string, stepIds: string[]): string[] {
    return stepIds.filter((id) => {
      const step = getChecklistStepConfig(checklistId, id);
      if (!step) return false;

      const navPermission = step.navPermission;
      if (!navPermission) return true;

      const [navItem, routeParam] = navPermission;
      return this.permissions.hasNavPermission(navItem, routeParam);
    });
  }

  // ---------------------------------------------------------------------------
  // Step completion (backend-backed)
  // ---------------------------------------------------------------------------

  /**
   * Returns whether a step in a checklist has been marked complete.
   * Returns false for any checklist or step that is absent from the stored state.
   */
  isStepCompleted(checklistId: string, stepId: string): boolean {
    return this._state[checklistId]?.[stepId] ?? false;
  }

  /**
   * Returns the completion progress for the given set of visible step IDs.
   * Completed count and total are both derived from the backend completion map
   * so that only persisted completions are counted.
   */
  getChecklistProgress(checklistId: string, stepIds: string[]): { completed: number; total: number } {
    const completed = stepIds.filter((id) => this.isStepCompleted(checklistId, id)).length;
    return { completed, total: stepIds.length };
  }

  /**
   * Returns true when every provided step ID is marked complete in the backend state.
   */
  isChecklistComplete(checklistId: string, stepIds: string[]): boolean {
    return stepIds.length > 0 && stepIds.every((id) => this.isStepCompleted(checklistId, id));
  }

  // ---------------------------------------------------------------------------
  // Local preferences (localStorage-backed, never written to backend)
  // ---------------------------------------------------------------------------

  /**
   * Returns whether the user has hidden this checklist in their current browser.
   */
  isHidden(checklistId: string): boolean {
    return this._localState[checklistId]?.hidden ?? false;
  }

  /**
   * Hides the checklist for the current browser session and future sessions.
   * Does not affect backend completion state.
   */
  hideChecklist(checklistId: string): void {
    this._patchLocalState(checklistId, { hidden: true });
  }

  /**
   * Restores the checklist to visible in the current browser.
   */
  showChecklist(checklistId: string): void {
    this._patchLocalState(checklistId, { hidden: false });
  }

  /**
   * Returns which panel (checklist vs. completion) was last shown for a
   * checklist. Defaults to 'complete-banner' — the original completion panel —
   * when no preference has been recorded yet.
   */
  getLastView(checklistId: string): ChecklistView {
    return this._localState[checklistId]?.lastView ?? 'complete-banner';
  }

  /**
   * Records which panel (checklist vs. completion) is currently shown for a
   * checklist so it can be restored to the same view later.
   */
  setLastView(checklistId: string, view: ChecklistView): void {
    this._patchLocalState(checklistId, { lastView: view });
  }

  // ---------------------------------------------------------------------------
  // Lifecycle state
  // ---------------------------------------------------------------------------

  /**
   * Returns the combined lifecycle state of a checklist based on backend
   * completion and the browser-local hidden preference.
   */
  getLifecycleState(checklistId: string, stepIds: string[]): ChecklistLifecycleState {
    const hidden = this.isHidden(checklistId);
    const complete = this.isChecklistComplete(checklistId, stepIds);

    if (hidden && complete) return 'post-completion';
    if (hidden) return 'hidden';
    if (complete) return 'complete';
    return 'active';
  }

  /**
   * Fetches the full checklist state from the API and updates the in-memory state.
   * Cancels any in-flight fetch when a newer call is made.
   * Sets isAvailable to false on 403/5xx so the UI can degrade gracefully.
   * Falls back to an empty state when the response body is null or has an
   * unexpected shape — this is the correct starting state for a fresh cluster.
   */
  fetchState = keepLatestTask(async () => {
    try {
      const response = await this.api.request.get(CHECKLIST_STATE_PATH);
      const body = (await response.json()) as ApiResponse;
      const data = body?.data;

      // Treat a null/undefined/unexpected shape as an empty state. This
      // happens on the very first request to a fresh cluster where no
      // checklist progress has been written yet. An empty object is the
      // correct starting state and can be updated via the UI.
      this._state = isChecklistStateData(data) ? data : {};
      this.isAvailable = true;
      return this._state;
    } catch (error) {
      if (isNonBlockingError(error)) {
        this.isAvailable = false;
        return;
      }
      throw error;
    }
  });

  /**
   * Marks a step as complete or incomplete and persists the change via the API.
   * The server performs a deep merge, so other checklist entries are preserved.
   * Updates in-memory state with the full merged result returned by the server.
   * Silently absorbs 403/5xx errors so a transient failure does not disrupt the UI.
   * Retains the current in-memory state when the response body is null or has
   * an unexpected shape so a transient body issue does not wipe local state.
   */
  updateStep = task(async (checklistId: string, stepId: string, completed: boolean) => {
    try {
      const response = await this.api.request.post(CHECKLIST_STATE_PATH, {
        checklist_state: { [checklistId]: { [stepId]: completed } },
      });
      const body = (await response.json()) as ApiResponse;
      const data = body?.data;

      // Fall back to current in-memory state when the server returns a
      // null or unexpected shape so a transient body issue does not wipe
      // the locally-known state.
      if (isChecklistStateData(data)) {
        this._state = data;
      }
      return this._state;
    } catch (error) {
      if (isNonBlockingError(error)) {
        return;
      }
      throw error;
    }
  });

  /**
   * Marks a step complete as the result of an explicit user action.
   * Unlike updateStep (used by background detection), this task propagates all
   * failures — including 403 and 5xx — so the UI can surface actionable errors.
   * In-memory state is only updated after the backend confirms the write; on
   * failure the step remains in its prior state so the user can retry.
   */
  markComplete = task(async (checklistId: string, stepId: string) => {
    const response = await this.api.request.post(CHECKLIST_STATE_PATH, {
      checklist_state: { [checklistId]: { [stepId]: true } },
    });
    const body = (await response.json()) as ApiResponse;
    const data = body?.data;

    if (!isChecklistStateData(data)) {
      throw new Error('Received checklist state in an unexpected format after update');
    }

    this._state = data;
    return this._state;
  });
}
