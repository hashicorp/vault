/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: BUSL-1.1
 */

import { module, test } from 'qunit';
import { setupTest } from 'ember-qunit';

module('Unit | Model | ldap/library', function (hooks) {
  setupTest(hooks);

  test('completeLibraryName should return name for non-hierarchical libraries', function (assert) {
    const store = this.owner.lookup('service:store');
    const model = store.createRecord('ldap/library', {
      name: 'simple-library',
      backend: 'ldap-test',
    });

    assert.strictEqual(
      model.completeLibraryName,
      'simple-library',
      'completeLibraryName returns name for non-hierarchical libraries'
    );
  });

  test('completeLibraryName should combine path_to_library and name for hierarchical libraries', function (assert) {
    const store = this.owner.lookup('service:store');
    const model = store.createRecord('ldap/library', {
      name: 'sa-prod',
      path_to_library: 'service-account/',
      backend: 'ldap-test',
    });

    assert.strictEqual(
      model.completeLibraryName,
      'service-account/sa-prod',
      'completeLibraryName combines path_to_library and name correctly'
    );
  });

  test('completeLibraryName should handle deeply nested hierarchical libraries', function (assert) {
    const store = this.owner.lookup('service:store');
    const model = store.createRecord('ldap/library', {
      name: 'library',
      path_to_library: 'region/country/city/department/',
      backend: 'ldap-test',
    });

    assert.strictEqual(
      model.completeLibraryName,
      'region/country/city/department/library',
      'completeLibraryName handles deeply nested paths correctly'
    );
  });

  test('completeLibraryName should remove trailing slash for directory-only models', function (assert) {
    const store = this.owner.lookup('service:store');
    const model = store.createRecord('ldap/library', {
      name: 'service-account/',
      backend: 'ldap-test',
    });

    assert.strictEqual(
      model.completeLibraryName,
      'service-account/',
      'completeLibraryName removes trailing slash for directory-only models'
    );
  });

  test('completeLibraryName should handle empty path_to_library', function (assert) {
    const store = this.owner.lookup('service:store');
    const model = store.createRecord('ldap/library', {
      name: 'root-library',
      path_to_library: '',
      backend: 'ldap-test',
    });

    assert.strictEqual(
      model.completeLibraryName,
      'root-library',
      'completeLibraryName handles empty path_to_library correctly'
    );
  });

  test('completeLibraryName should handle null path_to_library', function (assert) {
    const store = this.owner.lookup('service:store');
    const model = store.createRecord('ldap/library', {
      name: 'root-library',
      path_to_library: null,
      backend: 'ldap-test',
    });

    assert.strictEqual(
      model.completeLibraryName,
      'root-library',
      'completeLibraryName handles null path_to_library correctly'
    );
  });

  test('fetchStatus should use completeLibraryName', function (assert) {
    assert.expect(2);

    const store = this.owner.lookup('service:store');
    const adapter = store.adapterFor('ldap/library');

    // Mock the adapter's fetchStatus method
    adapter.fetchStatus = function (backend, libraryName) {
      assert.strictEqual(backend, 'ldap-test', 'Backend passed correctly');
      assert.strictEqual(libraryName, 'service-account/sa', 'Complete library name passed to adapter');
      return Promise.resolve([]);
    };

    const model = store.createRecord('ldap/library', {
      name: 'sa',
      path_to_library: 'service-account/',
      backend: 'ldap-test',
    });

    model.fetchStatus();
  });

  test('checkOutAccount should use completeLibraryName', function (assert) {
    assert.expect(3);

    const store = this.owner.lookup('service:store');
    const adapter = store.adapterFor('ldap/library');

    // Mock the adapter's checkOutAccount method
    adapter.checkOutAccount = function (backend, libraryName, ttl) {
      assert.strictEqual(backend, 'ldap-test', 'Backend passed correctly');
      assert.strictEqual(libraryName, 'service-account/sa', 'Complete library name passed to adapter');
      assert.strictEqual(ttl, '2h', 'TTL passed correctly');
      return Promise.resolve({});
    };

    const model = store.createRecord('ldap/library', {
      name: 'sa',
      path_to_library: 'service-account/',
      backend: 'ldap-test',
    });

    model.checkOutAccount('2h');
  });

  test('checkInAccount should use completeLibraryName', function (assert) {
    assert.expect(3);

    const store = this.owner.lookup('service:store');
    const adapter = store.adapterFor('ldap/library');

    // Mock the adapter's checkInAccount method
    adapter.checkInAccount = function (backend, libraryName, accounts) {
      assert.strictEqual(backend, 'ldap-test', 'Backend passed correctly');
      assert.strictEqual(libraryName, 'service-account/sa', 'Complete library name passed to adapter');
      assert.deepEqual(accounts, ['test@example.com'], 'Account array passed correctly');
      return Promise.resolve({});
    };

    const model = store.createRecord('ldap/library', {
      name: 'sa',
      path_to_library: 'service-account/',
      backend: 'ldap-test',
    });

    model.checkInAccount('test@example.com');
  });

  test('completeLibraryName should handle path_to_library without trailing slash', function (assert) {
    const store = this.owner.lookup('service:store');
    const model = store.createRecord('ldap/library', {
      name: 'sa',
      path_to_library: 'service-account', // No trailing slash
      backend: 'ldap-test',
    });

    assert.strictEqual(
      model.completeLibraryName,
      'service-accountsa',
      'completeLibraryName concatenates without adding slash when not present'
    );
  });

  test('completeLibraryName should handle complex nested directory structure', function (assert) {
    const store = this.owner.lookup('service:store');
    const model = store.createRecord('ldap/library', {
      name: 'production-service',
      path_to_library: 'org/division/team/environment/',
      backend: 'ldap-test',
    });

    assert.strictEqual(
      model.completeLibraryName,
      'org/division/team/environment/production-service',
      'completeLibraryName handles complex nested directory structure correctly'
    );
  });

  test('completeLibraryName should handle directory names with special characters', function (assert) {
    const store = this.owner.lookup('service:store');
    const model = store.createRecord('ldap/library', {
      name: 'service.account-123',
      path_to_library: 'org_division/team-2024/',
      backend: 'ldap-test',
    });

    assert.strictEqual(
      model.completeLibraryName,
      'org_division/team-2024/service.account-123',
      'completeLibraryName handles special characters in names and paths correctly'
    );
  });
});
