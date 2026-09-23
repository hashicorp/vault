/**
 * Copyright IBM Corp. 2016, 2026
 * SPDX-License-Identifier: BUSL-1.1
 */

import { setupMirage } from 'ember-cli-mirage/test-support';
import { setupTest } from 'ember-qunit';
import { module, test } from 'qunit';
import sinon from 'sinon';
import localStorage from 'vault/lib/local-storage';
import { CHECKLIST_LOCAL_STATE_KEY } from 'vault/utils/constants/checklist';

import type { Server } from 'miragejs';
import { Response as MirageResponse } from 'miragejs';
import type ChecklistStateService from 'vault/services/checklist-state';
import type PermissionsService from 'vault/services/permissions';

declare module '@ember/test-helpers' {
  interface TestContext {
    service: ChecklistStateService;
    server: Server;
  }
}

const CHECKLIST_ENDPOINT = 'sys/config/ui/checklist-state';

module('Unit | Service | checklist-state', function (hooks) {
  setupTest(hooks);
  setupMirage(hooks);

  hooks.beforeEach(function () {
    this.service = this.owner.lookup('service:checklist-state') as ChecklistStateService;
  });

  test('isStepCompleted returns false before any state has been fetched', function (assert) {
    assert.false(
      this.service.isStepCompleted('cluster-startup', 'create-admin'),
      'Returns false for unknown checklist'
    );
    assert.false(
      this.service.isStepCompleted('cluster-startup', 'unknown-step'),
      'Returns false for unknown step within a known checklist'
    );
  });

  module('#fetchState', function () {
    test('populates in-memory state from the API response', async function (assert) {
      const state = { 'cluster-startup': { 'create-admin': true, 'enable-ui': false } };
      this.server.get(CHECKLIST_ENDPOINT, () => ({ data: state }));

      await this.service.fetchState.perform();

      assert.true(
        this.service.isStepCompleted('cluster-startup', 'create-admin'),
        'Completed step returns true'
      );
      assert.false(
        this.service.isStepCompleted('cluster-startup', 'enable-ui'),
        'Step explicitly set to false returns false'
      );
      assert.false(
        this.service.isStepCompleted('cluster-startup', 'missing-step'),
        'Step absent from state returns false'
      );
    });

    test('accepts an empty state document from the API', async function (assert) {
      this.server.get(CHECKLIST_ENDPOINT, () => ({ data: {} }));

      await this.service.fetchState.perform();

      assert.false(
        this.service.isStepCompleted('cluster-startup', 'create-admin'),
        'Returns false when state is empty'
      );
    });

    test('falls back to empty state when the API response shape is unexpected', async function (assert) {
      this.server.get(CHECKLIST_ENDPOINT, () => ({
        data: { 'cluster-startup': 'not-an-object' },
      }));

      await this.service.fetchState.perform();

      assert.true(this.service.isAvailable, 'isAvailable remains true after unexpected shape');
      assert.false(
        this.service.isStepCompleted('cluster-startup', 'create-admin'),
        'All steps are false when state falls back to empty'
      );
    });

    test('falls back to empty state when the API returns null data', async function (assert) {
      this.server.get(CHECKLIST_ENDPOINT, () => ({ data: null }));

      await this.service.fetchState.perform();

      assert.true(this.service.isAvailable, 'isAvailable remains true after null data');
      assert.false(
        this.service.isStepCompleted('cluster-startup', 'create-admin'),
        'All steps are false when state falls back to empty'
      );
    });

    test('sets isAvailable to false and does not throw on 403', async function (assert) {
      this.server.get(CHECKLIST_ENDPOINT, () => new MirageResponse(403, {}, ''));

      await this.service.fetchState.perform();

      assert.false(this.service.isAvailable, 'isAvailable is false after a 403 response');
      assert.false(
        this.service.isStepCompleted('cluster-startup', 'create-admin'),
        'All steps remain false when fetch is forbidden'
      );
    });

    test('sets isAvailable to false and does not throw on 5xx', async function (assert) {
      this.server.get(CHECKLIST_ENDPOINT, () => new MirageResponse(503, {}, ''));

      await this.service.fetchState.perform();

      assert.false(this.service.isAvailable, 'isAvailable is false after a 5xx response');
    });

    test('resets isAvailable to true after a subsequent successful fetch', async function (assert) {
      this.server.get(CHECKLIST_ENDPOINT, () => new MirageResponse(403, {}, ''));
      await this.service.fetchState.perform();
      assert.false(this.service.isAvailable, 'isAvailable is false after failed fetch');

      this.server.get(CHECKLIST_ENDPOINT, () => ({ data: {} }));
      await this.service.fetchState.perform();
      assert.true(this.service.isAvailable, 'isAvailable resets to true after successful fetch');
    });
  });

  module('#updateStep', function () {
    test('sends the correct payload and updates in-memory state', async function (assert) {
      assert.expect(3);

      const mergedState = { 'cluster-startup': { 'create-admin': true } };

      this.server.post(CHECKLIST_ENDPOINT, (_schema: unknown, request: { requestBody: string }) => {
        const body = JSON.parse(request.requestBody);
        assert.deepEqual(
          body.checklist_state,
          { 'cluster-startup': { 'create-admin': true } },
          'POST body contains the correct checklist_state payload'
        );
        return { data: mergedState };
      });

      await this.service.updateStep.perform('cluster-startup', 'create-admin', true);

      assert.true(
        this.service.isStepCompleted('cluster-startup', 'create-admin'),
        'Step is marked complete after update'
      );
      assert.false(
        this.service.isStepCompleted('cluster-startup', 'enable-ui'),
        'Steps absent from merged state remain false'
      );
    });

    test('persists false (incomplete) state for a step', async function (assert) {
      this.server.post(CHECKLIST_ENDPOINT, () => ({
        data: { 'cluster-startup': { 'create-admin': false } },
      }));

      await this.service.updateStep.perform('cluster-startup', 'create-admin', false);

      assert.false(
        this.service.isStepCompleted('cluster-startup', 'create-admin'),
        'Step explicitly marked incomplete remains false'
      );
    });

    test('reflects the full merged state returned by the server', async function (assert) {
      const mergedState = {
        'cluster-startup': { 'create-admin': true, 'enable-ui': false },
        resilience: { 'enable-replication': false },
      };

      this.server.post(CHECKLIST_ENDPOINT, () => ({ data: mergedState }));

      await this.service.updateStep.perform('cluster-startup', 'create-admin', true);

      assert.true(
        this.service.isStepCompleted('cluster-startup', 'create-admin'),
        'Updated step is complete'
      );
      assert.false(
        this.service.isStepCompleted('resilience', 'enable-replication'),
        'Merged step from another checklist is reflected'
      );
    });

    test('retains in-memory state when the update response shape is unexpected', async function (assert) {
      // Pre-populate state so we can assert it is preserved.
      this.server.get(CHECKLIST_ENDPOINT, () => ({
        data: { 'cluster-startup': { 'create-admin': true } },
      }));
      await this.service.fetchState.perform();

      // Server returns null — should silently retain the prior in-memory state.
      this.server.post(CHECKLIST_ENDPOINT, () => ({ data: null }));
      await this.service.updateStep.perform('cluster-startup', 'enable-ui', true);

      assert.true(
        this.service.isStepCompleted('cluster-startup', 'create-admin'),
        'Prior state is preserved when update response shape is unexpected'
      );
      assert.false(
        this.service.isStepCompleted('cluster-startup', 'enable-ui'),
        'Unconfirmed step remains false when update response is null'
      );
    });

    test('absorbs 403 silently and leaves in-memory state unchanged', async function (assert) {
      this.server.get(CHECKLIST_ENDPOINT, () => ({
        data: { 'cluster-startup': { 'create-admin': true } },
      }));
      await this.service.fetchState.perform();

      this.server.post(CHECKLIST_ENDPOINT, () => new MirageResponse(403, {}, ''));
      await this.service.updateStep.perform('cluster-startup', 'enable-ui', true);

      assert.true(
        this.service.isStepCompleted('cluster-startup', 'create-admin'),
        'Prior state is preserved after a forbidden update'
      );
      assert.false(
        this.service.isStepCompleted('cluster-startup', 'enable-ui'),
        'Unconfirmed step remains false after a forbidden update'
      );
    });

    test('absorbs 5xx silently and leaves in-memory state unchanged', async function (assert) {
      this.server.get(CHECKLIST_ENDPOINT, () => ({
        data: { 'cluster-startup': { 'create-admin': true } },
      }));
      await this.service.fetchState.perform();

      this.server.post(CHECKLIST_ENDPOINT, () => new MirageResponse(500, {}, ''));
      await this.service.updateStep.perform('cluster-startup', 'enable-ui', true);

      assert.true(
        this.service.isStepCompleted('cluster-startup', 'create-admin'),
        'Prior state is preserved after a server error on update'
      );
      assert.false(
        this.service.isStepCompleted('cluster-startup', 'enable-ui'),
        'Unconfirmed step remains false after a server error on update'
      );
    });
  });

  module('#markComplete', function () {
    test('persists the step to the backend and updates in-memory state', async function (assert) {
      this.server.post(CHECKLIST_ENDPOINT, (_schema: unknown, request: { requestBody: string }) => {
        const body = JSON.parse(request.requestBody);
        assert.deepEqual(
          body.checklist_state,
          { 'cluster-startup': { 'tvp-cli': true } },
          'POST body sends the correct delta for the completed step'
        );
        return { data: { 'cluster-startup': { 'tvp-cli': true } } };
      });

      await this.service.markComplete.perform('cluster-startup', 'tvp-cli');

      assert.true(
        this.service.isStepCompleted('cluster-startup', 'tvp-cli'),
        'Step is marked complete in in-memory state after confirmed backend write'
      );
    });

    test('propagates 403 errors so the UI can surface them', async function (assert) {
      this.server.post(CHECKLIST_ENDPOINT, () => new MirageResponse(403, {}, ''));

      await assert.rejects(
        this.service.markComplete.perform('cluster-startup', 'tvp-cli'),
        'Rejects with a 403 error'
      );

      assert.false(
        this.service.isStepCompleted('cluster-startup', 'tvp-cli'),
        'Step remains incomplete when the backend rejects the write'
      );
    });

    test('propagates 5xx errors so the UI can surface them', async function (assert) {
      this.server.post(CHECKLIST_ENDPOINT, () => new MirageResponse(500, {}, ''));

      await assert.rejects(
        this.service.markComplete.perform('cluster-startup', 'tvp-cli'),
        'Rejects with a 5xx error'
      );

      assert.false(
        this.service.isStepCompleted('cluster-startup', 'tvp-cli'),
        'Step remains incomplete after a server error'
      );
    });

    test('state is consistent on retry — succeeds after a prior failure', async function (assert) {
      this.server.post(CHECKLIST_ENDPOINT, () => new MirageResponse(500, {}, ''));
      await assert.rejects(this.service.markComplete.perform('cluster-startup', 'tvp-cli'));

      this.server.post(CHECKLIST_ENDPOINT, () => ({
        data: { 'cluster-startup': { 'tvp-cli': true } },
      }));
      await this.service.markComplete.perform('cluster-startup', 'tvp-cli');

      assert.true(
        this.service.isStepCompleted('cluster-startup', 'tvp-cli'),
        'Step is complete after a successful retry'
      );
    });
  });
});

module('Unit | Service | checklist-state (local concerns)', function (hooks) {
  setupTest(hooks);

  let getItemStub = sinon.stub();
  let setItemStub = sinon.stub();

  hooks.beforeEach(function () {
    getItemStub = sinon.stub(localStorage, 'getItem');
    setItemStub = sinon.stub(localStorage, 'setItem');
    sinon.stub(localStorage, 'removeItem');
    getItemStub.withArgs(CHECKLIST_LOCAL_STATE_KEY).returns(null);
    let service: ChecklistStateService | undefined;
    Object.defineProperty(this, 'service', {
      configurable: true,
      get: () => (service ??= this.owner.lookup('service:checklist-state') as ChecklistStateService),
    });
  });

  hooks.afterEach(function () {
    sinon.restore();
  });

  module('#getLastView / #setLastView', function () {
    test('getLastView defaults to complete-banner when no preference has been recorded', function (assert) {
      assert.strictEqual(
        this.service.getLastView('cluster-startup'),
        'complete-banner',
        'Defaults to the completion panel'
      );
    });

    test('setLastView persists the view and getLastView reflects it', function (assert) {
      this.service.setLastView('cluster-startup', 'checklist');

      assert.strictEqual(this.service.getLastView('cluster-startup'), 'checklist');
      assert.true(
        setItemStub.calledWith(CHECKLIST_LOCAL_STATE_KEY, {
          'cluster-startup': { lastView: 'checklist' },
        }),
        'Saves the last view nested under the checklist ID'
      );
    });

    test('loads pre-existing last-view data from localStorage on service creation', function (assert) {
      getItemStub.withArgs(CHECKLIST_LOCAL_STATE_KEY).returns({
        'cluster-startup': { lastView: 'checklist' },
      });
      const service = this.owner.lookup('service:checklist-state') as ChecklistStateService;

      assert.strictEqual(service.getLastView('cluster-startup'), 'checklist');
    });

    test('ignores malformed localStorage data without breaking page load', function (assert) {
      getItemStub.withArgs(CHECKLIST_LOCAL_STATE_KEY).throws(new SyntaxError('Invalid JSON'));
      const service = this.owner.lookup('service:checklist-state') as ChecklistStateService;

      assert.strictEqual(
        service.getLastView('cluster-startup'),
        'complete-banner',
        'Falls back gracefully to the default view'
      );
    });
  });

  module('#isHidden / #hideChecklist / #showChecklist', function () {
    test('isHidden returns false for a checklist that has not been hidden', function (assert) {
      assert.false(this.service.isHidden('cluster-startup'), 'Returns false when no hidden state exists');
    });

    test('hideChecklist persists the checklist as hidden, nested under its ID', function (assert) {
      this.service.hideChecklist('cluster-startup');

      assert.true(this.service.isHidden('cluster-startup'), 'isHidden returns true after hiding');
      assert.true(
        setItemStub.calledWith(CHECKLIST_LOCAL_STATE_KEY, {
          'cluster-startup': { hidden: true },
        }),
        'Saves updated local state to localStorage'
      );
    });

    test('hideChecklist persists repeated hidden-state patches', function (assert) {
      this.service.hideChecklist('cluster-startup');
      this.service.hideChecklist('cluster-startup');

      assert.true(this.service.isHidden('cluster-startup'));
      assert.strictEqual(setItemStub.callCount, 2, 'Each call still persists (idempotent patch)');
    });

    test('showChecklist marks the checklist as not hidden', function (assert) {
      this.service.hideChecklist('cluster-startup');
      this.service.showChecklist('cluster-startup');

      assert.false(this.service.isHidden('cluster-startup'), 'isHidden returns false after showing');
      assert.true(
        setItemStub.calledWith(CHECKLIST_LOCAL_STATE_KEY, {
          'cluster-startup': { hidden: false },
        }),
        'Saves updated local state to localStorage'
      );
    });

    test('loads pre-existing hidden checklists from localStorage on service creation', function (assert) {
      getItemStub.withArgs(CHECKLIST_LOCAL_STATE_KEY).returns({ 'cluster-startup': { hidden: true } });
      const service = this.owner.lookup('service:checklist-state') as ChecklistStateService;

      assert.true(service.isHidden('cluster-startup'), 'Pre-existing hidden state is loaded');
    });

    test('ignores unparseable localStorage data without breaking page load', function (assert) {
      getItemStub.withArgs(CHECKLIST_LOCAL_STATE_KEY).throws(new SyntaxError('Invalid JSON'));
      const service = this.owner.lookup('service:checklist-state') as ChecklistStateService;

      assert.false(service.isHidden('cluster-startup'), 'Falls back gracefully to no hidden checklists');
    });

    test('ignores a localStorage value that is a plain string instead of an object', function (assert) {
      getItemStub.withArgs(CHECKLIST_LOCAL_STATE_KEY).returns('not-an-object');
      const service = this.owner.lookup('service:checklist-state') as ChecklistStateService;

      assert.false(service.isHidden('cluster-startup'), 'String-typed storage value is discarded');
    });

    test('ignores a localStorage value that is an array instead of an object', function (assert) {
      getItemStub.withArgs(CHECKLIST_LOCAL_STATE_KEY).returns(['cluster-startup']);
      const service = this.owner.lookup('service:checklist-state') as ChecklistStateService;

      assert.false(service.isHidden('cluster-startup'), 'Array-typed storage value is discarded');
    });

    test('ignores a localStorage entry whose hidden field is a non-boolean string', function (assert) {
      getItemStub.withArgs(CHECKLIST_LOCAL_STATE_KEY).returns({ 'cluster-startup': { hidden: 'yes' } });
      const service = this.owner.lookup('service:checklist-state') as ChecklistStateService;

      assert.false(service.isHidden('cluster-startup'), 'Entry with string hidden value is discarded');
    });

    test('ignores a localStorage entry whose lastView field is an unrecognised value', function (assert) {
      getItemStub
        .withArgs(CHECKLIST_LOCAL_STATE_KEY)
        .returns({ 'cluster-startup': { lastView: 'unknown-view' } });
      const service = this.owner.lookup('service:checklist-state') as ChecklistStateService;

      // The whole stored object is rejected by isChecklistLocalStateData, so
      // _localState falls back to {}, and getLastView returns its default.
      assert.strictEqual(
        service.getLastView('cluster-startup'),
        'complete-banner',
        'Entry with unrecognised lastView is discarded and falls back to default'
      );
    });

    test('preserves other checklists and preferences when patching a single checklist', function (assert) {
      getItemStub.withArgs(CHECKLIST_LOCAL_STATE_KEY).returns({
        'other-checklist': { hidden: true },
        'cluster-startup': { lastView: 'checklist' },
      });
      const service = this.owner.lookup('service:checklist-state') as ChecklistStateService;

      service.hideChecklist('cluster-startup');

      assert.true(service.isHidden('other-checklist'), 'Unrelated checklist is untouched');
      assert.strictEqual(
        service.getLastView('cluster-startup'),
        'checklist',
        'Existing preference for the patched checklist is preserved'
      );
      assert.true(service.isHidden('cluster-startup'), 'New preference for the patched checklist is applied');
    });
  });

  module('#getChecklistProgress / #isChecklistComplete', function () {
    test('reports zero progress when no steps are complete', function (assert) {
      const progress = this.service.getChecklistProgress('cluster-startup', ['create-admin', 'enable-ui']);
      assert.deepEqual(progress, { completed: 0, total: 2 }, 'Returns zero completed out of total');
    });

    test('counts only backend-confirmed completions', function (assert) {
      // Directly seed the private state to simulate a successful fetchState
      this.service['_state'] = { 'cluster-startup': { 'create-admin': true } };

      const progress = this.service.getChecklistProgress('cluster-startup', ['create-admin', 'enable-ui']);
      assert.deepEqual(progress, { completed: 1, total: 2 }, 'Counts only confirmed completions');
    });

    test('isChecklistComplete returns true only when all visible steps are done', function (assert) {
      this.service['_state'] = {
        'cluster-startup': { 'create-admin': true, 'enable-ui': true },
      };
      const stepIds = ['create-admin', 'enable-ui'];

      assert.true(this.service.isChecklistComplete('cluster-startup', stepIds), 'All steps complete');
      assert.false(
        this.service.isChecklistComplete('cluster-startup', [...stepIds, 'missing-step']),
        'Returns false when a visible step is not yet complete'
      );
    });

    test('isChecklistComplete returns false for an empty step list', function (assert) {
      assert.false(
        this.service.isChecklistComplete('cluster-startup', []),
        'Empty step list is never considered complete'
      );
    });
  });

  module('#getLifecycleState', function () {
    const STEPS = ['create-admin', 'enable-ui'];

    test('returns active when not hidden and not complete', function (assert) {
      assert.strictEqual(this.service.getLifecycleState('cluster-startup', STEPS), 'active');
    });

    test('returns complete when all steps done and not hidden', function (assert) {
      this.service['_state'] = { 'cluster-startup': { 'create-admin': true, 'enable-ui': true } };

      assert.strictEqual(this.service.getLifecycleState('cluster-startup', STEPS), 'complete');
    });

    test('returns hidden when user dismissed before completing', function (assert) {
      this.service.hideChecklist('cluster-startup');

      assert.strictEqual(this.service.getLifecycleState('cluster-startup', STEPS), 'hidden');
    });

    test('returns post-completion when all steps done and checklist is hidden', function (assert) {
      this.service['_state'] = { 'cluster-startup': { 'create-admin': true, 'enable-ui': true } };
      this.service.hideChecklist('cluster-startup');

      assert.strictEqual(this.service.getLifecycleState('cluster-startup', STEPS), 'post-completion');
    });
  });

  module('#isEnterpriseFeature', function () {
    test('returns true on Enterprise', function (assert) {
      (this.owner.lookup('service:version') as unknown as { type: unknown }).type = 'enterprise';
      const service = this.owner.lookup('service:checklist-state') as ChecklistStateService;
      assert.true(service.isEnterpriseFeature, 'isEnterpriseFeature is true for Enterprise clusters');
    });

    test('returns false on Community Edition', function (assert) {
      (this.owner.lookup('service:version') as unknown as { type: unknown }).type = 'community';
      const service = this.owner.lookup('service:checklist-state') as ChecklistStateService;
      assert.false(service.isEnterpriseFeature, 'isEnterpriseFeature is false for CE clusters');
    });
  });

  module('#getVisibleSteps', function () {
    test('returns all steps on Community Edition when permissions are granted', function (assert) {
      (this.owner.lookup('service:version') as unknown as { type: unknown }).type = 'community';
      const service = this.owner.lookup('service:checklist-state') as ChecklistStateService;
      const permissions = this.owner.lookup('service:permissions') as unknown as PermissionsService;
      sinon.stub(permissions, 'hasNavPermission').returns(true);

      const visible = service.getVisibleSteps('cluster-startup', [
        'tvp-cli',
        'namespaces',
        'policy',
        'auth',
        'kv',
      ]);

      assert.deepEqual(
        visible,
        ['tvp-cli', 'namespaces', 'policy', 'auth', 'kv'],
        'step filtering is permission-based; CE gating happens at checklist entry'
      );
    });

    test('returns all steps when Enterprise and all permissions granted', function (assert) {
      (this.owner.lookup('service:version') as unknown as { type: unknown }).type = 'enterprise';
      const service = this.owner.lookup('service:checklist-state') as ChecklistStateService;
      const permissions = this.owner.lookup('service:permissions') as unknown as PermissionsService;
      sinon.stub(permissions, 'hasNavPermission').returns(true);

      const visible = service.getVisibleSteps('cluster-startup', [
        'tvp-cli',
        'namespaces',
        'policy',
        'auth',
        'kv',
      ]);

      assert.deepEqual(visible, ['tvp-cli', 'namespaces', 'policy', 'auth', 'kv'], 'All steps visible');
    });

    test('filters out steps for which the token lacks nav permission', function (assert) {
      (this.owner.lookup('service:version') as unknown as { type: unknown }).type = 'enterprise';
      const service = this.owner.lookup('service:checklist-state') as ChecklistStateService;
      const permissions = this.owner.lookup('service:permissions') as unknown as PermissionsService;
      sinon.stub(permissions, 'hasNavPermission').callsFake((navItem: string) => navItem !== 'policies');

      const visible = service.getVisibleSteps('cluster-startup', [
        'tvp-cli',
        'namespaces',
        'policy',
        'auth',
        'kv',
      ]);

      assert.false(visible.includes('policy'), 'policy step excluded when token lacks policies permission');
      assert.true(visible.includes('auth'), 'auth step included when token has access permission');
    });

    test('getChecklistProgress computed over visible steps only', function (assert) {
      (this.owner.lookup('service:version') as unknown as { type: unknown }).type = 'enterprise';
      const service = this.owner.lookup('service:checklist-state') as ChecklistStateService;
      service['_state'] = { 'cluster-startup': { policy: true, auth: true } };
      const permissions = this.owner.lookup('service:permissions') as unknown as PermissionsService;
      sinon.stub(permissions, 'hasNavPermission').returns(true);

      const allSteps = ['tvp-cli', 'namespaces', 'policy', 'auth', 'kv'];
      const visible = service.getVisibleSteps('cluster-startup', allSteps);
      const { completed, total } = service.getChecklistProgress('cluster-startup', visible);

      assert.strictEqual(total, 5, 'Total equals number of visible steps (all visible here)');
      assert.strictEqual(completed, 2, 'Only backend-confirmed completions count toward progress');
    });

    test('getChecklistProgress denominator reflects permission-filtered step count', function (assert) {
      (this.owner.lookup('service:version') as unknown as { type: unknown }).type = 'enterprise';
      const service = this.owner.lookup('service:checklist-state') as ChecklistStateService;
      service['_state'] = { 'cluster-startup': { policy: true, auth: true } };
      const permissions = this.owner.lookup('service:permissions') as unknown as PermissionsService;
      // User can only see policy and auth (not namespaces)
      sinon.stub(permissions, 'hasNavPermission').callsFake((navItem: string) => navItem !== 'access');

      const allSteps = ['tvp-cli', 'namespaces', 'policy', 'auth', 'kv'];
      const visible = service.getVisibleSteps('cluster-startup', allSteps);
      const { completed, total } = service.getChecklistProgress('cluster-startup', visible);

      assert.strictEqual(
        total,
        3,
        'Denominator is 3: tvp-cli (no permission check), policy (permitted), kv (no check) — not 5'
      );
      assert.strictEqual(completed, 1, 'Only policy is complete among the visible steps');
    });
  });
});
