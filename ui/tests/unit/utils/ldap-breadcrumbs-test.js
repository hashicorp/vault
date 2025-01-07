/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: BUSL-1.1
 */

import { ldapBreadcrumbs, roleRoutes } from 'ldap/utils/ldap-breadcrumbs';
import { module, test } from 'qunit';

module('Unit | Utility | ldap breadcrumbs', function (hooks) {
  hooks.beforeEach(async function () {
    this.mountPath = 'my-engine';
    this.roleType = 'static';
    const routeParams = (childResource) => {
      return [this.mountPath, this.roleType, childResource];
    };
    this.testCrumbs = (path, { lastItemCurrent }) => {
      return ldapBreadcrumbs(path, routeParams, roleRoutes, lastItemCurrent);
    };
  });

  test('it generates crumbs when the path is a directory', function (assert) {
    const path = 'prod/org/';
    let actual = this.testCrumbs(path, { lastItemCurrent: true });
    let expected = [
      { label: 'prod', route: 'roles.subdirectory', models: [this.mountPath, this.roleType, 'prod/'] },
      { label: 'org' },
    ];
    assert.propEqual(actual, expected, 'crumbs are correct when lastItemCurrent = true');

    actual = this.testCrumbs(path, { lastItemCurrent: false });
    expected = [
      { label: 'prod', route: 'roles.subdirectory', models: [this.mountPath, this.roleType, 'prod/'] },
      { label: 'org', route: 'roles.subdirectory', models: [this.mountPath, this.roleType, 'prod/org/'] },
    ];
    assert.propEqual(actual, expected, 'crumbs are correct when lastItemCurrent = false');
  });

  test('it generates crumbs when the path is not a directory', function (assert) {
    const path = 'prod/org/admin';
    let actual = this.testCrumbs(path, { lastItemCurrent: true });
    let expected = [
      { label: 'prod', route: 'roles.subdirectory', models: [this.mountPath, this.roleType, 'prod/'] },
      { label: 'org', route: 'roles.subdirectory', models: [this.mountPath, this.roleType, 'prod/org/'] },
      { label: 'admin' },
    ];
    assert.propEqual(actual, expected, 'crumbs are correct when lastItemCurrent = true');

    actual = this.testCrumbs(path, { lastItemCurrent: false });
    expected = [
      { label: 'prod', route: 'roles.subdirectory', models: [this.mountPath, this.roleType, 'prod/'] },
      { label: 'org', route: 'roles.subdirectory', models: [this.mountPath, this.roleType, 'prod/org/'] },
      {
        label: 'admin',
        route: 'roles.role.details',
        models: [this.mountPath, this.roleType, 'prod/org/admin'],
      },
    ];
    assert.propEqual(actual, expected, 'crumbs are correct when lastItemCurrent = false');
  });

  test('it generates crumbs when the path is the top-level', function (assert) {
    const path = 'prod/';
    let actual = this.testCrumbs(path, { lastItemCurrent: true });
    let expected = [{ label: 'prod' }];
    assert.propEqual(actual, expected, 'crumbs are correct when lastItemCurrent = true');

    actual = this.testCrumbs(path, { lastItemCurrent: false });
    expected = [
      { label: 'prod', route: 'roles.subdirectory', models: [this.mountPath, this.roleType, 'prod/'] },
    ];
    assert.propEqual(actual, expected, 'crumbs are correct when lastItemCurrent = false');
  });

  test('it fails gracefully when no path', function (assert) {
    const path = undefined;
    const actual = this.testCrumbs(path, { lastItemCurrent: false });
    assert.propEqual(actual, [], 'returns empty array when path is null');
  });
});
