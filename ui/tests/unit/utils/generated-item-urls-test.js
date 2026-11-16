/**
 * Copyright IBM Corp. 2026
 * SPDX-License-Identifier: BUSL-1.1
 */

import { module, test } from 'qunit';
import { setupTest } from 'vault/tests/helpers';
import openApiSpec from '@hashicorp/vault-client-typescript/openapi.json';
import { reducePathsByPathName } from 'vault/utils/openapi-helpers';
import { generatedItemHelpUrl, generatedItemUrls } from 'vault/utils/generated-item-urls';

// Mirrors path-help getPaths: a mount's help response lists its paths relative to the mount.
const pathInfoFor = (mountParam, apiPath, itemType) =>
  Object.entries(openApiSpec.paths)
    .filter(([path]) => path.startsWith(`/auth/{${mountParam}}/`))
    .map(([path, details]) => [path.replace(`/auth/{${mountParam}}`, ''), details])
    .reduce(reducePathsByPathName, { apiPath, itemType, itemTypes: [], paths: [] });

// Item URLs for every auth method whose generated item pages are reachable from navigation.
const CASES = [
  { mountParam: 'userpass_mount_path', itemType: 'user', base: 'users', id: 'item-1' },
  { mountParam: 'userpass_mount_path', itemType: 'user', base: 'users', id: 'team a' },
  { mountParam: 'ldap_mount_path', itemType: 'user', base: 'users', id: 'item-1' },
  { mountParam: 'ldap_mount_path', itemType: 'group', base: 'groups', id: 'item-1' },
  { mountParam: 'okta_mount_path', itemType: 'user', base: 'users', id: 'item-1' },
  { mountParam: 'okta_mount_path', itemType: 'group', base: 'groups', id: 'item-1' },
  { mountParam: 'radius_mount_path', itemType: 'user', base: 'users', id: 'item-1' },
  { mountParam: 'cert_mount_path', itemType: 'certificate', base: 'certs', id: 'item-1' },
  { mountParam: 'kubernetes_mount_path', itemType: 'role', base: 'role', id: 'item-1' },
  { mountParam: 'tpm_mount_path', itemType: 'role', base: 'role', id: 'item-1' },
];

module('Unit | Utility | generated item urls', function (hooks) {
  setupTest(hooks);

  hooks.beforeEach(function () {
    // The util returns paths relative to the API service base path, so prefix /v1 to compare full URLs.
    this.urlsFor = ({ pathInfo, itemType, id }) => {
      const urls = generatedItemUrls(pathInfo, itemType);
      return {
        list: `/v1${urls.list}`,
        read: `/v1${urls.read(id)}`,
        create: `/v1${urls.create(id)}`,
        update: `/v1${urls.update(id)}`,
        delete: `/v1${urls.delete(id)}`,
      };
    };
  });

  CASES.forEach(({ mountParam, itemType, base, id }) => {
    test(`${mountParam} ${itemType} "${id}"`, function (assert) {
      const mountPath = 'my-mount';
      const apiPath = `auth/${mountPath}/`;
      const pathInfo = pathInfoFor(mountParam, apiPath, itemType);
      const encoded = id.replace(' ', '%20');
      assert.deepEqual(
        this.urlsFor({ pathInfo, itemType, id }),
        {
          list: `/v1/auth/my-mount/${base}/`,
          read: `/v1/auth/my-mount/${base}/${encoded}`,
          create: `/v1/auth/my-mount/${base}/${id}`,
          update: `/v1/auth/my-mount/${base}/${id}`,
          delete: `/v1/auth/my-mount/${base}/${id}`,
        },
        'list, read, create, update and delete URLs'
      );
    });
  });

  test('help URL points at the create endpoint with example path parameters', function (assert) {
    const pathInfo = pathInfoFor('userpass_mount_path', 'auth/my-mount/', 'user');
    assert.strictEqual(
      generatedItemHelpUrl(pathInfo, 'user'),
      '/v1/auth/my-mount/users/example?help=true',
      'help URL describes the user create endpoint'
    );
  });
});
