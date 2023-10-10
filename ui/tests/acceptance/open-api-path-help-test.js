import { module, test } from 'qunit';
import Model from '@ember-data/model';
import { setupApplicationTest } from 'vault/tests/helpers';
import authPage from 'vault/tests/pages/auth';
import { deleteAuthCmd, deleteEngineCmd, mountAuthCmd, mountEngineCmd, runCmd } from '../helpers/commands';
import openApiDrivenAttributes from '../helpers/openapi-driven-attributes';

/**
 * This set of tests is for ensuring that backend changes to the OpenAPI spec
 * are known by UI developers and adequately addressed in the UI. When changes
 * are detected from this set of tests, they should be updated to pass and
 * smoke tested to ensure changes to not break the GUI workflow.
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

  module('auth: TLS certificates', function (hooks) {
    hooks.beforeEach(async function () {
      this.backend = 'cert-openapi';
      await runCmd(mountAuthCmd('cert', this.backend), false);
    });
    hooks.afterEach(async function () {
      await runCmd(deleteAuthCmd(this.backend), false);
    });
    test('getPaths returns correct paths', async function (assert) {
      const helpUrl = `auth/${this.backend}`;
      const result = await this.pathHelp.getPaths(helpUrl, this.backend);
      assert.deepEqual(
        result.paths,
        [
          {
            path: '/certs/',
            itemType: 'certificate',
            itemName: 'Certificate',
            operations: ['get', 'list'],
            action: undefined,
            navigation: true,
            param: false,
          },
          {
            path: '/certs/{name}',
            itemType: 'certificate',
            itemName: 'Certificate',
            operations: ['get', 'post', 'delete'],
            action: 'Create',
            navigation: false,
            param: 'name',
          },
        ],
        'correct paths'
      );
      // No paths with navigation=true means they won't show up as tabs on the method page
    });
    test('getProps returns correct model attributes', async function (assert) {
      const config = await this.pathHelp.getProps(`/v1/auth/${this.backend}/config?help=1`, this.backend);
      assert.deepEqual(config, openApiDrivenAttributes.certConfig, 'config attributes');

      const cert = await this.pathHelp.getProps(
        `/v1/auth/${this.backend}/certs/example?help=true`,
        this.backend
      );
      assert.deepEqual(cert, openApiDrivenAttributes.certCert, 'cert attributes');
    });
  });

  module('auth: GCP', function (hooks) {
    hooks.beforeEach(async function () {
      this.backend = 'gcp-openapi';
      await runCmd(mountAuthCmd('gcp', this.backend), false);
    });
    hooks.afterEach(async function () {
      await runCmd(deleteAuthCmd(this.backend), false);
    });
    test('getPaths returns correct paths', async function (assert) {
      const helpUrl = `auth/${this.backend}`;
      const result = await this.pathHelp.getPaths(helpUrl, this.backend);
      assert.deepEqual(result.paths, [], 'correct paths');
    });
    test('getProps returns correct model attributes', async function (assert) {
      const config = await this.pathHelp.getProps(`/v1/auth/${this.backend}/config?help=1`, this.backend);
      assert.deepEqual(config, openApiDrivenAttributes.gcpConfig, 'config');
    });
  });

  module('auth: github', function (hooks) {
    hooks.beforeEach(async function () {
      this.backend = 'github-openapi';
      await runCmd(mountAuthCmd('github', this.backend), false);
    });
    hooks.afterEach(async function () {
      await runCmd(deleteAuthCmd(this.backend), false);
    });
    test('getPaths returns correct paths', async function (assert) {
      const helpUrl = `auth/${this.backend}`;
      const result = await this.pathHelp.getPaths(helpUrl, this.backend);
      assert.deepEqual(result.paths, [], 'correct paths');
    });
    test('getProps returns correct model attributes', async function (assert) {
      const config = await this.pathHelp.getProps(`/v1/auth/${this.backend}/config?help=1`, this.backend);
      assert.deepEqual(config, openApiDrivenAttributes.githubConfig, 'config');
    });
  });

  module('auth: jwt', function (hooks) {
    hooks.beforeEach(async function () {
      this.backend = 'jwt-openapi';
      await runCmd(mountAuthCmd('jwt', this.backend), false);
    });
    hooks.afterEach(async function () {
      await runCmd(deleteAuthCmd(this.backend), false);
    });
    test('getPaths returns correct paths', async function (assert) {
      const helpUrl = `auth/${this.backend}`;
      const result = await this.pathHelp.getPaths(helpUrl, this.backend);
      assert.deepEqual(result.paths, [], 'correct paths');
    });
    test('getProps returns correct model attributes', async function (assert) {
      const config = await this.pathHelp.getProps(`/v1/auth/${this.backend}/config?help=1`, this.backend);
      assert.deepEqual(config, openApiDrivenAttributes.jwtConfig, 'config');
    });
  });

  module('auth: kubernetes', function (hooks) {
    hooks.beforeEach(async function () {
      this.backend = 'k8s-openapi';
      await runCmd(mountAuthCmd('kubernetes', this.backend), false);
    });
    hooks.afterEach(async function () {
      await runCmd(deleteAuthCmd(this.backend), false);
    });
    test('getPaths returns correct paths', async function (assert) {
      const helpUrl = `auth/${this.backend}`;
      const result = await this.pathHelp.getPaths(helpUrl, this.backend);
      assert.deepEqual(
        result.paths,
        [
          {
            action: undefined,
            itemName: 'Role',
            itemType: 'role',
            navigation: true,
            operations: ['get', 'list'],
            param: false,
            path: '/role/',
          },
          {
            action: 'Create',
            itemName: 'Role',
            itemType: 'role',
            navigation: false,
            operations: ['get', 'post', 'delete'],
            param: 'name',
            path: '/role/{name}',
          },
        ],
        'correct paths'
      );
    });
    test('getProps returns correct model attributes', async function (assert) {
      const config = await this.pathHelp.getProps(`/v1/auth/${this.backend}/config?help=1`, this.backend);
      assert.deepEqual(config, openApiDrivenAttributes.k8sConfig, 'config');

      const role = await this.pathHelp.getProps(`/v1/auth/${this.backend}/role/example?help=1`, this.backend);
      assert.deepEqual(role, openApiDrivenAttributes.k8sRole, 'role');
    });
  });

  module('auth: ldap', function (hooks) {
    hooks.beforeEach(async function () {
      this.backend = 'ldap-openapi';
      await runCmd(mountAuthCmd('ldap', this.backend), false);
    });
    hooks.afterEach(async function () {
      await runCmd(deleteAuthCmd(this.backend), false);
    });
    test('getPaths returns correct paths', async function (assert) {
      const helpUrl = `auth/${this.backend}`;
      const result = await this.pathHelp.getPaths(helpUrl, this.backend);
      assert.deepEqual(
        result.paths,
        [
          {
            action: 'Configure',
            itemName: undefined,
            itemType: undefined,
            navigation: false,
            operations: ['get', 'post'],
            param: false,
            path: '/config',
          },
          {
            action: undefined,
            itemName: 'Group',
            itemType: 'group',
            navigation: true,
            operations: ['get', 'list'],
            param: false,
            path: '/groups/',
          },
          {
            action: 'Create',
            itemName: 'Group',
            itemType: 'group',
            navigation: false,
            operations: ['get', 'post', 'delete'],
            param: 'name',
            path: '/groups/{name}',
          },
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
            param: 'name',
            path: '/users/{name}',
          },
        ],
        'correct paths'
      );
    });
    test('getProps returns correct model attributes', async function (assert) {
      const config = await this.pathHelp.getProps(`/v1/auth/${this.backend}/config?help=1`, this.backend);
      assert.deepEqual(config, openApiDrivenAttributes.ldapConfig, 'config');

      const groups = await this.pathHelp.getProps(
        `/v1/auth/${this.backend}/groups/example?help=1`,
        this.backend
      );
      assert.deepEqual(groups, openApiDrivenAttributes.ldapGroup, 'groups');

      const users = await this.pathHelp.getProps(
        `/v1/auth/${this.backend}/users/example?help=1`,
        this.backend
      );
      assert.deepEqual(users, openApiDrivenAttributes.ldapUser, 'users');
    });
  });

  module('auth: okta', function (hooks) {
    hooks.beforeEach(async function () {
      this.backend = 'okta-openapi';
      await runCmd(mountAuthCmd('okta', this.backend), false);
    });
    hooks.afterEach(async function () {
      await runCmd(deleteAuthCmd(this.backend), false);
    });
    test('getPaths returns correct paths', async function (assert) {
      const helpUrl = `auth/${this.backend}`;
      const result = await this.pathHelp.getPaths(helpUrl, this.backend);
      assert.deepEqual(
        result.paths,
        [
          {
            action: 'Configure',
            itemName: undefined,
            itemType: undefined,
            navigation: false,
            operations: ['get', 'post'],
            param: false,
            path: '/config',
          },
          {
            action: undefined,
            itemName: 'Group',
            itemType: 'group',
            navigation: true,
            operations: ['get', 'list'],
            param: false,
            path: '/groups/',
          },
          {
            action: 'Create',
            itemName: 'Group',
            itemType: 'group',
            navigation: false,
            operations: ['get', 'post', 'delete'],
            param: 'name',
            path: '/groups/{name}',
          },
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
            param: 'name',
            path: '/users/{name}',
          },
        ],
        'correct paths'
      );
    });
    test('getProps returns correct model attributes', async function (assert) {
      const config = await this.pathHelp.getProps(`/v1/auth/${this.backend}/config?help=1`, this.backend);
      assert.deepEqual(config, openApiDrivenAttributes.oktaConfig, 'config');

      const groups = await this.pathHelp.getProps(
        `/v1/auth/${this.backend}/groups/example?help=1`,
        this.backend
      );
      assert.deepEqual(groups, openApiDrivenAttributes.oktaGroup, 'groups');

      const users = await this.pathHelp.getProps(
        `/v1/auth/${this.backend}/users/example?help=1`,
        this.backend
      );
      assert.deepEqual(users, openApiDrivenAttributes.oktaUser, 'users');
    });
  });

  module('auth: radius', function (hooks) {
    hooks.beforeEach(async function () {
      this.backend = 'radius-openapi';
      await runCmd(mountAuthCmd('radius', this.backend), false);
    });
    hooks.afterEach(async function () {
      await runCmd(deleteAuthCmd(this.backend), false);
    });
    test('getPaths returns correct paths', async function (assert) {
      const helpUrl = `auth/${this.backend}`;
      const result = await this.pathHelp.getPaths(helpUrl, this.backend);
      assert.deepEqual(
        result.paths,
        [
          {
            action: 'Configure',
            itemName: undefined,
            itemType: undefined,
            navigation: false,
            operations: ['get', 'post'],
            param: false,
            path: '/config',
          },
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
            param: 'name',
            path: '/users/{name}',
          },
        ],
        'correct paths'
      );
    });
    test('getProps returns correct model attributes', async function (assert) {
      const config = await this.pathHelp.getProps(`/v1/auth/${this.backend}/config?help=1`, this.backend);
      assert.deepEqual(config, openApiDrivenAttributes.radiusConfig, 'config');

      const users = await this.pathHelp.getProps(
        `/v1/auth/${this.backend}/users/example?help=1`,
        this.backend
      );
      assert.deepEqual(users, openApiDrivenAttributes.radiusUser, 'users');
    });
  });

  module('auth: aws', function (hooks) {
    hooks.beforeEach(async function () {
      this.backend = 'aws-openapi';
      await runCmd(mountAuthCmd('aws', this.backend), false);
    });
    hooks.afterEach(async function () {
      await runCmd(deleteAuthCmd(this.backend), false);
    });
    test('getPaths returns correct paths', async function (assert) {
      const helpUrl = `auth/${this.backend}`;
      const result = await this.pathHelp.getPaths(helpUrl, this.backend);
      assert.deepEqual(result.paths, [], 'correct paths');
    });
    test('getProps returns correct model attributes', async function (assert) {
      const config = await this.pathHelp.getProps(`/v1/auth/${this.backend}/?help=1`, this.backend);
      assert.deepEqual(config, openApiDrivenAttributes.awsConfig, 'config');
    });
  });
});
