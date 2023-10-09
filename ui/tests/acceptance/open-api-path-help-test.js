import { module, test } from 'qunit';
import Model from '@ember-data/model';
import { setupApplicationTest } from 'vault/tests/helpers';
import authPage from 'vault/tests/pages/auth';
import { deleteAuthCmd, deleteEngineCmd, mountAuthCmd, mountEngineCmd, runCmd } from '../helpers/commands';
import openApiDrivenAttributes from '../helpers/openapi-driven-attributes';

/**
 * This set of tests is for ensuring that backend changes to the OpenAPI spec
 * are known by UI developers and adequately addressed in the UI. In addition
 * to updating the response
 */
module('Acceptance | OpenAPI path help test', function (hooks) {
  setupApplicationTest(hooks);
  hooks.beforeEach(function () {
    this.pathHelp = this.owner.lookup('service:pathHelp');
    this.newModel = Model.extend({});
    return authPage.login();
  });

  module('engine: ssh role', function (hooks) {
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
      assert.deepEqual(result, openApiDrivenAttributes.sshRole);
    });
  });

  module('engine: kmip enterprise', function (hooks) {
    hooks.beforeEach(async function () {
      this.backend = 'kmip-openapi';
      await runCmd(mountEngineCmd('kmip', this.backend), false);
    });
    hooks.afterEach(async function () {
      await runCmd(deleteEngineCmd(this.backend), false);
    });

    test('getProps returns correct model attributes', async function (assert) {
      const config = await this.pathHelp.getProps(`/v1/${this.backend}/config?help=1`, this.backend);
      assert.deepEqual(config, openApiDrivenAttributes.kmipConfig, 'kmip config');

      const role = await this.pathHelp.getProps(
        `/v1/${this.backend}/scope/example/role/example?help=1`,
        this.backend
      );
      assert.deepEqual(role, openApiDrivenAttributes.kmipRole, 'kmip role');
    });
  });

  module('engine: pki', function (hooks) {
    hooks.beforeEach(async function () {
      this.backend = 'pki-openapi';
      await runCmd(mountEngineCmd('pki', this.backend), false);
    });
    hooks.afterEach(async function () {
      await runCmd(deleteEngineCmd(this.backend), false);
    });

    test('getProps returns correct model attributes', async function (assert) {
      const role = await this.pathHelp.getProps(`/v1/${this.backend}/roles/example?help=1`, this.backend);
      assert.deepEqual(role, openApiDrivenAttributes.pkiRole, 'pki role');

      const signCsr = await this.pathHelp.getProps(
        `/v1/${this.backend}/issuer/example/sign-intermediate?help=1`,
        this.backend
      );
      assert.deepEqual(signCsr, openApiDrivenAttributes.pkiSignCsr, 'pki sign intermediate');

      const tidy = await this.pathHelp.getProps(`/v1/${this.backend}/config/auto-tidy?help=1`, this.backend);
      assert.deepEqual(tidy, openApiDrivenAttributes.pkiTidy, 'pki tidy');

      const certGenerate = await this.pathHelp.getProps(
        `/v1/${this.backend}/issue/example?help=1`,
        this.backend
      );
      assert.deepEqual(certGenerate, openApiDrivenAttributes.pkiCertGenerate, 'pki cert generate');

      const certSign = await this.pathHelp.getProps(`/v1/${this.backend}/sign/example?help=1`, this.backend);
      assert.deepEqual(certSign, openApiDrivenAttributes.pkiCertSign, 'pki cert generate');

      const acme = await this.pathHelp.getProps(`/v1/${this.backend}/config/acme?help=1`, this.backend);
      assert.deepEqual(acme, openApiDrivenAttributes.pkiAcme, 'pki acme');

      const cluster = await this.pathHelp.getProps(`/v1/${this.backend}/config/cluster?help=1`, this.backend);
      assert.deepEqual(cluster, openApiDrivenAttributes.pkiCluster, 'pki cluster');

      const urls = await this.pathHelp.getProps(`/v1/${this.backend}/config/urls?help=1`, this.backend);
      assert.deepEqual(urls, openApiDrivenAttributes.pkiUrls, 'pki urls');
    });
  });

  module('auth: userpass', function (hooks) {
    hooks.beforeEach(async function () {
      this.backend = 'userpass7';
      await runCmd(mountAuthCmd('userpass', this.backend), false);
    });
    hooks.afterEach(async function () {
      await runCmd(deleteAuthCmd(this.backend), false);
    });
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
      assert.deepEqual(result, openApiDrivenAttributes.userpassUser);
    });
  });

  module('auth: approle', function (hooks) {
    hooks.beforeEach(async function () {
      this.backend = 'approle-openapi';
      await runCmd(mountAuthCmd('approle', this.backend), false);
    });
    hooks.afterEach(async function () {
      await runCmd(deleteAuthCmd(this.backend), false);
    });
    test('getPaths returns correct paths', async function (assert) {
      const helpUrl = `auth/${this.backend}`;
      const result = await this.pathHelp.getPaths(helpUrl, this.backend);
      assert.deepEqual(result.paths, [], 'correct paths');
    });
  });

  module('auth: azure', function (hooks) {
    hooks.beforeEach(async function () {
      this.backend = 'azure-openapi';
      await runCmd(mountAuthCmd('azure', this.backend), false);
    });
    hooks.afterEach(async function () {
      await runCmd(deleteAuthCmd(this.backend), false);
    });
    test('getPaths returns correct paths', async function (assert) {
      const helpUrl = `auth/${this.backend}`;
      const result = await this.pathHelp.getPaths(helpUrl, this.backend);
      assert.deepEqual(result.paths, [], 'correct paths');
      // No paths with navigation=true means they won't show up as tabs on the method page
    });
    test('getProps returns correct model attributes', async function (assert) {
      const helpUrl = `/v1/auth/${this.backend}/config?help=1`;
      const result = await this.pathHelp.getProps(helpUrl, this.backend);
      assert.deepEqual(result, openApiDrivenAttributes.azureConfig);
    });
  });

  /* TODO: fill in these other OpenAPI tests -- supported auth backends (extends AuthConfig model)

    cert
    gcp
    github
    jwt
    kubernetes
    ldap
    okta
    radius
    aws/client
    aws/tidy
  */
});
