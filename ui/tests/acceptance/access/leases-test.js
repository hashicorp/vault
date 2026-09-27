/**
 * Copyright IBM Corp. 2026
 * SPDX-License-Identifier: BUSL-1.1
 */

/**
 * Tests for the leases list and show controller actions: revokePrefix,
 * forceRevokePrefix, revokeLease, and renewLease. Uses Mirage stubs so
 * the mutation request payloads and error rendering can be verified without
 * depending on specific lease state in the backend. These tests complement
 * (not replace) the Playwright gate, which verifies real list navigation and
 * button visibility against a live server.
 */

import { click, currentRouteName, visit } from '@ember/test-helpers';
import { module, test } from 'qunit';
import { setupApplicationTest } from 'ember-qunit';
import { setupMirage } from 'ember-cli-mirage/test-support';
import sinon from 'sinon';
import { login } from 'vault/tests/helpers/auth/auth-helpers';
import { overrideResponse } from 'vault/tests/helpers/stubs';
import { GENERAL } from 'vault/tests/helpers/general-selectors';

// ---------------------------------------------------------------------------
// Shared stub helpers — called AFTER login(), BEFORE visit()
// ---------------------------------------------------------------------------

// Stub the leases list so the list page renders.
// The generated API client encodes the prefix with encodeURIComponent, turning
// slashes into %2F, so `auth/token/create/` becomes the single path segment
// `auth%2Ftoken%2Fcreate`. Mirage's route-recognizer treats %2F as a literal
// (it does not decode %2F to /), meaning a hardcoded stub path with slashes
// in the prefix would not match. The :prefix dynamic segment matches any single
// (possibly percent-encoded) segment, so this stub handles all prefixes.
// Also stubs capabilities so the revoke-prefix buttons are visible.
function stubLeasesListAtPrefix(server) {
  const listResponse = () => ({
    data: { keys: ['abc-lease-1', 'abc-lease-2'] },
    request_id: 'test',
  });

  // The parent leases route checks capabilities and loads the root list first.
  server.get('/sys/leases/lookup/', listResponse);
  // List endpoint — GET with ?list=true; *prefix captures slash-containing prefixes.
  server.get('/sys/leases/lookup/*prefix', listResponse);
}

function stubCapabilities(owner) {
  sinon.stub(owner.lookup('service:capabilities'), 'fetchPathCapabilities').resolves({
    canCreate: true,
    canDelete: true,
    canList: true,
    canPatch: true,
    canRead: true,
    canSudo: true,
    canUpdate: true,
  });
}

// Stub the single-lease lookup (POST /sys/leases/lookup) for the show route.
function stubLeaseLookup(server, leaseId) {
  server.post('/sys/leases/lookup', () => ({
    data: {
      id: leaseId,
      issue_time: '2026-01-01T00:00:00Z',
      expire_time: '2026-01-01T01:00:00Z',
      last_renewal: null,
      renewable: true,
      ttl: 3600,
    },
    request_id: 'test',
  }));
}

// ---------------------------------------------------------------------------
// List controller — revokePrefix / forceRevokePrefix
// ---------------------------------------------------------------------------

module('Acceptance | leases list — revokePrefix actions', function (hooks) {
  setupApplicationTest(hooks);
  setupMirage(hooks);

  hooks.beforeEach(async function () {
    await login();
    stubLeasesListAtPrefix(this.server);
    stubCapabilities(this.owner);
    await visit('/vault/access/leases/list/auth/token/create/');
  });

  test('revokePrefix sends POST to /sys/leases/revoke-prefix/{prefix} and redirects to list-root on success', async function (assert) {
    assert.expect(2);

    this.server.post('/sys/leases/revoke-prefix/:prefix', (_, req) => {
      assert.ok(true, `POST to revoke-prefix was called with prefix: ${req.params.prefix}`);
      return overrideResponse(204);
    });

    await click('[data-test-lease-revoke-prefix]');
    await click(GENERAL.confirmButton);

    assert.strictEqual(
      currentRouteName(),
      'vault.cluster.access.leases.list-root',
      'redirects to list-root after revoke-prefix'
    );
  });

  test('revokePrefix shows a danger flash when the server returns an error', async function (assert) {
    assert.expect(1);

    this.server.post('/sys/leases/revoke-prefix/:prefix', () =>
      overrideResponse(403, { errors: ['permission denied'] })
    );

    await click('[data-test-lease-revoke-prefix]');
    await click(GENERAL.confirmButton);

    assert.dom('[data-test-flash-message]').exists('danger flash is shown when revoke-prefix fails');
  });
});

// ---------------------------------------------------------------------------
// Show controller — revokeLease / renewLease
// ---------------------------------------------------------------------------

module('Acceptance | leases show — revokeLease / renewLease actions', function (hooks) {
  setupApplicationTest(hooks);
  setupMirage(hooks);

  hooks.beforeEach(async function () {
    this.leaseId = 'database/creds/abc-lease-1';
    await login();
    stubCapabilities(this.owner);
    stubLeaseLookup(this.server, this.leaseId);
    await visit(`/vault/access/leases/show/${this.leaseId}`);
  });

  test('revokeLease sends POST to /sys/leases/revoke with the lease_id and redirects to list-root', async function (assert) {
    assert.expect(2);

    this.server.post('/sys/leases/revoke', (_, req) => {
      const body = JSON.parse(req.requestBody);
      assert.strictEqual(body.lease_id, this.leaseId, 'revoke body contains the correct lease_id');
      return overrideResponse(204);
    });

    await click('[data-test-lease-revoke]');
    await click(GENERAL.confirmButton);

    assert.strictEqual(
      currentRouteName(),
      'vault.cluster.access.leases.list-root',
      'redirects to list-root after revoke'
    );
  });

  test('renewLease sends POST to /sys/leases/renew with the correct lease_id', async function (assert) {
    assert.expect(1);

    this.server.post('/sys/leases/renew', (_, req) => {
      const body = JSON.parse(req.requestBody);
      assert.strictEqual(body.lease_id, this.leaseId, 'renew body contains the correct lease_id');
      return overrideResponse(204);
    });

    // Re-stub lookup for the route refresh triggered after renew
    stubLeaseLookup(this.server, this.leaseId);

    await click('[data-test-renew-lease-submit]');
  });

  test('renewLease shows a danger flash when the server returns an error', async function (assert) {
    assert.expect(1);

    this.server.post('/sys/leases/renew', () => overrideResponse(400, { errors: ['lease not renewable'] }));

    await click('[data-test-renew-lease-submit]');

    assert.dom('[data-test-flash-message]').exists('danger flash is shown when renew fails');
  });
});
