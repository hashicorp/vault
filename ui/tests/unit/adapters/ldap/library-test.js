/**
 * Copyright IBM Corp. 2016, 2025
 * SPDX-License-Identifier: BUSL-1.1
 */

import { module, test } from 'qunit';
import { setupTest } from 'ember-qunit';
import { setupMirage } from 'ember-cli-mirage/test-support';

module('Unit | Adapter | ldap/library', function (hooks) {
  setupTest(hooks);
  setupMirage(hooks);

  hooks.beforeEach(function () {
    this.store = this.owner.lookup('service:store');
    this.adapter = this.store.adapterFor('ldap/library');
  });

  module('List LDAP Libraries', function () {
    test('non-hierarchical - should make request to correct endpoint', async function (assert) {
      assert.expect(1);

      this.server.get('/ldap-test/library', (schema, req) => {
        assert.ok(req.queryParams.list, 'GET request made to correct endpoint when listing records');
        return { data: { keys: ['test-library'] } };
      });

      await this.store.query('ldap/library', { backend: 'ldap-test' });
    });

    test('hierarchical - should make request to correct endpoint with path', async function (assert) {
      assert.expect(1);

      this.server.get('/ldap-test/library/service-accounts/', (schema, req) => {
        assert.ok(
          req.queryParams.list,
          'GET request made to correct endpoint when listing hierarchical records'
        );
        return { data: { keys: ['prod-library', 'dev-library'] } };
      });

      await this.store.query('ldap/library', { backend: 'ldap-test', path_to_library: 'service-accounts/' });
    });
  });

  module('Query LDAP Library Record', function () {
    test('non-hierarchical - should make request to correct endpoint', async function (assert) {
      assert.expect(1);

      this.server.get('/ldap-test/library/test-library', () => {
        assert.ok('GET request made to correct endpoint when querying non-hierarchical record');
      });

      await this.store.queryRecord('ldap/library', { backend: 'ldap-test', name: 'test-library' });
    });

    test('hierarchical - should handle URL-encoded paths correctly', async function (assert) {
      assert.expect(3);

      const encodedName = 'service-account1%2Fsa1'; // URL-encoded "service-account1/sa1"
      const expectedData = {
        name: 'Test Library',
        service_account_names: ['test@example.com'],
      };

      this.server.get('/ldap-test/library/service-account1/sa1', () => {
        assert.ok('GET request made with decoded hierarchical path');
        return { data: expectedData };
      });

      const result = await this.store.queryRecord('ldap/library', {
        backend: 'ldap-test',
        name: encodedName,
      });

      assert.strictEqual(result.name, 'sa1', 'Library name extracted correctly from hierarchical path');
      assert.strictEqual(result.path_to_library, 'service-account1/', 'Path to library set correctly');
    });
  });

  module('Create LDAP Library', function () {
    test('non-hierarchical - should make request to correct endpoint', async function (assert) {
      assert.expect(1);

      this.server.post('/ldap-test/library/simple-library', () => {
        assert.ok('POST request made to correct endpoint for non-hierarchical library creation');
      });

      await this.store
        .createRecord('ldap/library', {
          backend: 'ldap-test',
          name: 'simple-library',
          service_account_names: ['test@example.com'],
        })
        .save();
    });

    test('hierarchical - should make request to correct endpoint with full path', async function (assert) {
      assert.expect(1);

      this.server.post('/ldap-test/library/service-account/prod-library', () => {
        assert.ok('POST request made to correct endpoint for hierarchical library creation');
      });

      await this.store
        .createRecord('ldap/library', {
          backend: 'ldap-test',
          name: 'prod-library',
          path_to_library: 'service-account/',
          service_account_names: ['prod@example.com'],
        })
        .save();
    });
  });

  module('Update LDAP Library', function () {
    test('non-hierarchical - should make request to correct endpoint', async function (assert) {
      assert.expect(1);

      this.server.post('/ldap-test/library/simple-library', () => {
        assert.ok('POST request made to correct endpoint for non-hierarchical library update');
      });

      this.store.pushPayload('ldap/library', {
        modelName: 'ldap/library',
        backend: 'ldap-test',
        name: 'simple-library',
        service_account_names: ['test@example.com'],
      });

      const record = this.store.peekRecord('ldap/library', 'simple-library');
      record.service_account_names = ['updated@example.com'];
      await record.save();
    });

    test('hierarchical - should make request using completeLibraryName', async function (assert) {
      assert.expect(1);

      // This is the key test - ensuring the full hierarchical path is used for updates
      this.server.post('/ldap-test/library/service-account/prod-library', () => {
        assert.ok('POST request made to correct hierarchical endpoint for library update');
      });

      this.store.pushPayload('ldap/library', {
        modelName: 'ldap/library',
        backend: 'ldap-test',
        name: 'prod-library',
        path_to_library: 'service-account/',
        service_account_names: ['prod@example.com'],
      });

      const record = this.store.peekRecord('ldap/library', 'prod-library');
      record.service_account_names = ['updated-prod@example.com'];
      await record.save();
    });
  });

  module('Delete LDAP Library', function () {
    test('non-hierarchical - should make request to correct endpoint', async function (assert) {
      assert.expect(1);

      this.server.delete('/ldap-test/library/simple-library', () => {
        assert.ok('DELETE request made to correct endpoint for non-hierarchical library deletion');
      });

      this.store.pushPayload('ldap/library', {
        modelName: 'ldap/library',
        backend: 'ldap-test',
        name: 'simple-library',
        service_account_names: ['test@example.com'],
      });

      await this.store.peekRecord('ldap/library', 'simple-library').destroyRecord();
    });

    test('hierarchical - should make request using completeLibraryName', async function (assert) {
      assert.expect(1);

      this.server.delete('/ldap-test/library/service-account/prod-library', () => {
        assert.ok('DELETE request made to correct hierarchical endpoint for library deletion');
      });

      this.store.pushPayload('ldap/library', {
        modelName: 'ldap/library',
        backend: 'ldap-test',
        name: 'prod-library',
        path_to_library: 'service-account/',
        service_account_names: ['prod@example.com'],
      });

      await this.store.peekRecord('ldap/library', 'prod-library').destroyRecord();
    });
  });

  module('Fetch Library Status', function () {
    test('non-hierarchical - should make request to correct endpoint', async function (assert) {
      assert.expect(1);

      this.server.get('/ldap-test/library/simple-library/status', () => {
        assert.ok('GET request made to correct endpoint when fetching non-hierarchical library status');
      });

      await this.adapter.fetchStatus('ldap-test', 'simple-library');
    });

    test('hierarchical - should make request with full path', async function (assert) {
      assert.expect(1);

      this.server.get('/ldap-test/library/service-account/prod-library/status', () => {
        assert.ok('GET request made to correct hierarchical endpoint for fetchStatus');
        return { data: {} };
      });

      await this.adapter.fetchStatus('ldap-test', 'service-account/prod-library');
    });
  });

  module('Check Out Library Account', function () {
    test('non-hierarchical - should make request to correct endpoint', async function (assert) {
      assert.expect(1);

      this.server.post('/ldap-test/library/simple-library/check-out', (schema, req) => {
        const json = JSON.parse(req.requestBody);
        assert.strictEqual(
          json.ttl,
          '1h',
          'POST request made to correct endpoint when checking out non-hierarchical library'
        );
        return {
          data: { password: 'test', service_account_name: 'foo@bar.com' },
        };
      });

      await this.adapter.checkOutAccount('ldap-test', 'simple-library', '1h');
    });

    test('hierarchical - should make request with full path', async function (assert) {
      assert.expect(2);

      this.server.post('/ldap-test/library/west-region/sa-prod/check-out', (schema, req) => {
        const json = JSON.parse(req.requestBody);
        assert.strictEqual(json.ttl, '2h', 'TTL passed correctly for hierarchical library check-out');
        assert.ok('POST request made to correct hierarchical endpoint for checkOutAccount');
        return {
          data: { password: 'test-password', service_account_name: 'sa-prod@company.com' },
        };
      });

      await this.adapter.checkOutAccount('ldap-test', 'west-region/sa-prod', '2h');
    });
  });

  module('Check In Library Account', function () {
    test('non-hierarchical - should make request to correct endpoint', async function (assert) {
      assert.expect(1);

      this.server.post('/ldap-test/library/simple-library/check-in', (schema, req) => {
        const json = JSON.parse(req.requestBody);
        assert.deepEqual(
          json.service_account_names,
          ['foo@bar.com'],
          'POST request made to correct endpoint when checking in non-hierarchical library service accounts'
        );
        return {
          data: {
            'check-ins': ['foo@bar.com'],
          },
        };
      });

      await this.adapter.checkInAccount('ldap-test', 'simple-library', ['foo@bar.com']);
    });

    test('hierarchical - should make request with full path', async function (assert) {
      assert.expect(2);

      this.server.post('/ldap-test/library/west-region/sa-prod/check-in', (schema, req) => {
        const json = JSON.parse(req.requestBody);
        assert.deepEqual(
          json.service_account_names,
          ['sa-prod@company.com', 'sa-backup@company.com'],
          'Service account names passed correctly for hierarchical library check-in'
        );
        assert.ok('POST request made to correct hierarchical endpoint for checkInAccount');
        return {
          data: {
            'check-ins': ['sa-prod@company.com', 'sa-backup@company.com'],
          },
        };
      });

      await this.adapter.checkInAccount('ldap-test', 'west-region/sa-prod', [
        'sa-prod@company.com',
        'sa-backup@company.com',
      ]);
    });
  });

  module('Edge Cases and Complex Scenarios', function () {
    test('deeply nested hierarchical paths should work correctly', async function (assert) {
      assert.expect(1);

      const deepPath = 'region/country/city/department/team/library-name';

      this.server.get(`/ldap-test/library/${deepPath}/status`, () => {
        assert.ok('GET request made to correct deeply nested hierarchical endpoint');
        return { data: {} };
      });

      await this.adapter.fetchStatus('ldap-test', deepPath);
    });

    test('special characters in hierarchical paths should work correctly', async function (assert) {
      assert.expect(1);

      const specialPath = 'service-account_123/library.name-test';

      this.server.post(`/ldap-test/library/${specialPath}/check-out`, () => {
        assert.ok('POST request made to correct endpoint with special characters in hierarchical path');
        return {
          data: { password: 'test', service_account_name: 'test@domain.com' },
        };
      });

      await this.adapter.checkOutAccount('ldap-test', specialPath, '1h');
    });

    test('complex hierarchical parsing in queryRecord should work correctly', async function (assert) {
      assert.expect(4);

      const encodedName = 'org%2Fteam%2Fservice';
      const responseData = {
        name: 'Test Service Account',
        service_account_names: ['service@org.com'],
        ttl: 3600,
      };

      this.server.get('/ldap-test/library/org/team/service', () => {
        return { data: responseData };
      });

      const result = await this.store.queryRecord('ldap/library', {
        backend: 'ldap-test',
        name: encodedName,
      });

      assert.strictEqual(result.backend, 'ldap-test', 'Backend set correctly');
      assert.strictEqual(result.name, 'service', 'Library name extracted from hierarchical path');
      assert.strictEqual(result.path_to_library, 'org/team/', 'Path to library extracted correctly');
      assert.strictEqual(result.ttl, 3600, 'Library data preserved correctly');
    });

    test('single-level hierarchical path edge case should work correctly', async function (assert) {
      assert.expect(3);

      const encodedName = 'parent%2Fchild';
      const responseData = {
        name: 'Child Library',
        service_account_names: ['child@example.com'],
      };

      this.server.get('/ldap-test/library/parent/child', () => {
        return { data: responseData };
      });

      const result = await this.store.queryRecord('ldap/library', {
        backend: 'ldap-test',
        name: encodedName,
      });

      assert.strictEqual(result.name, 'child', 'Child library name extracted correctly');
      assert.strictEqual(result.path_to_library, 'parent/', 'Parent path extracted correctly');
      assert.strictEqual(result.backend, 'ldap-test', 'Backend preserved correctly');
    });
  });
});
