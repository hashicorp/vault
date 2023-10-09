import { module, test } from 'qunit';
import Model from '@ember-data/model';
import { setupApplicationTest } from 'vault/tests/helpers';
import authPage from 'vault/tests/pages/auth';
import { deleteAuthCmd, deleteEngineCmd, mountAuthCmd, mountEngineCmd, runCmd } from '../helpers/commands';
import openapiDrivenAttributes from '../helpers/openapi-driven-attributes';

module('Acceptance | OpenAPI path help test', function (hooks) {
  setupApplicationTest(hooks);
  hooks.beforeEach(function () {
    this.pathHelp = this.owner.lookup('service:pathHelp');
    this.newModel = Model.extend({});
    return authPage.login();
  });

  module('auth: userpass', function (hooks) {
    hooks.beforeEach(async function () {
      this.backend = 'userpass7';
      await runCmd(mountAuthCmd('userpass', this.backend), false);
    });
    hooks.afterEach(async function () {
      await runCmd(deleteAuthCmd(this.backend), false);
    });
    // First we fetch the base openAPI to get relevant paths
    test('getPaths returns correct paths', async function (assert) {
      const helpUrl = `auth/${this.backend}/`;
      const baseResult = await this.pathHelp.getPaths(helpUrl, this.backend);
      assert.deepEqual(
        baseResult.paths,
        [
          {
            action: undefined,
            itemName: 'User',
            itemType: 'user',
            navigation: true,
            operations: ['get', 'list'],
            param: false,
            path: '/users/',
          },
          {
            action: 'Create',
            itemName: 'User',
            itemType: 'user',
            navigation: false,
            operations: ['get', 'post', 'delete'],
            param: 'username',
            path: '/users/{username}',
          },
        ],
        'userpass has correct paths response when no itemType'
      );
      const resultWithItemType = await this.pathHelp.getPaths(helpUrl, this.backend, 'user');
      assert.deepEqual(
        resultWithItemType.paths,
        [
          {
            path: '/users/',
            itemType: 'user',
            itemName: 'User',
            operations: ['get', 'list'],
            action: undefined,
            navigation: true,
            param: false,
          },
          {
            path: '/users/{username}',
            itemType: 'user',
            itemName: 'User',
            operations: ['get', 'post', 'delete'],
            action: 'Create',
            navigation: false,
            param: 'username',
          },
        ],
        'userpass has correct paths response when itemType = user'
      );
    });
    test('getProps for user returns correct model attributes', async function (assert) {
      const helpUrl = `/v1/auth/${this.backend}/users/example?help=true`;
      const result = await this.pathHelp.getProps(helpUrl, this.backend);
      assert.deepEqual(result, openapiDrivenAttributes.userpassUser);
    });
  });

  module('engine: role-ssh', function (hooks) {
    hooks.beforeEach(async function () {
      this.backend = 'ssh-openapi';
      await runCmd(mountEngineCmd('ssh', this.backend), false);
    });
    hooks.afterEach(async function () {
      await runCmd(deleteEngineCmd(this.backend), false);
    });

    test('getProps returns correct model attributes', async function (assert) {
      const helpUrl = `/v1/${this.backend}/roles/example?help=1`;
      const result = await this.pathHelp.getProps(helpUrl, this.backend);
      assert.deepEqual(result, openapiDrivenAttributes.sshRole);
    });
  });
  /* TODO: fill in these other OpenAPI tests
  module('auth: azure', function (hooks) {});
  module('auth: gcp', function (hooks) {});
  module('auth: github', function (hooks) {});
  module('auth: jwt', function (hooks) {});
  module('auth: kubernetes', function (hooks) {});
  module('auth: ldap', function (hooks) {});
  module('auth: okta', function (hooks) {});
  module('auth: radius', function (hooks) {});
  module('auth: kmip/config', function (hooks) {});
  module('auth: kmip/role', function (hooks) {});
  module('engine: pki/role', function (hooks) {});
  module('engine: pki/tidy', function (hooks) {});
  module('engine: pki/certificate/base', function (hooks) {});
  module('engine: pki/acme', function (hooks) {});
  module('engine: pki/config', function (hooks) {});
  */
});
