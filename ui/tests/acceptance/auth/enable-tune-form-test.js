/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: BUSL-1.1
 */

import { module, test } from 'qunit';
import { setupApplicationTest } from 'ember-qunit';
import { v4 as uuidv4 } from 'uuid';

import { login } from 'vault/tests/helpers/auth/auth-helpers';
import { visit } from '@ember/test-helpers';
import { deleteAuthCmd, runCmd } from 'vault/tests/helpers/commands';
import testHelper from './test-helper';
import { GENERAL } from 'vault/tests/helpers/general-selectors';

// These models use openAPI so we assert the form inputs using an acceptance test
// The default selector is to use GENERAL.inputByAttr()
// custom fields should be added to the this.customSelectorss object
module('Acceptance | auth enable tune form test', function (hooks) {
  setupApplicationTest(hooks);
  hooks.beforeEach(async function () {
    // these tend to be the same across models because they share the same mount-config model
    // if necessary, they can be overridden in the individual module
    this.mountFields = [
      'path',
      'description',
      'local',
      'sealWrap',
      'config.listingVisibility',
      'config.defaultLeaseTtl',
      'config.maxLeaseTtl',
      'config.tokenType',
      'config.auditNonHmacRequestKeys',
      'config.auditNonHmacResponseKeys',
      'config.passthroughRequestHeaders',
      'config.allowedResponseHeaders',
      'config.pluginVersion',
    ];
  });

  module('azure', function (hooks) {
    hooks.beforeEach(async function () {
      this.type = 'azure';
      this.path = `${this.type}-${uuidv4()}`;
      this.tuneFields = [
        'environment',
        'identityTokenAudience',
        'identityTokenTtl',
        'maxRetries',
        'maxRetryDelay',
        'resource',
        'retryDelay',
        'rootPasswordTtl',
        'tenantId',
      ];
      this.tuneToggles = { 'Azure Options': ['clientId', 'clientSecret'] };
      await login();
      return visit('/vault/settings/auth/enable');
    });
    hooks.afterEach(async function () {
      await runCmd(deleteAuthCmd(this.path), false);
    });
    testHelper(test);
  });

  module('jwt', function (hooks) {
    hooks.beforeEach(async function () {
      this.type = 'jwt';
      this.path = `${this.type}-${uuidv4()}`;
      this.customSelectors = {
        providerConfig: `${GENERAL.fieldByAttr('providerConfig')} textarea`,
      };
      this.tuneFields = [
        'defaultRole',
        'jwksCaPem',
        'jwksUrl',
        'namespaceInState',
        'oidcDiscoveryUrl',
        'oidcResponseMode',
        'oidcResponseTypes',
        'providerConfig',
        'unsupportedCriticalCertExtensions',
      ];
      this.tuneToggles = {
        'JWT Options': [
          'oidcClientId',
          'oidcClientSecret',
          'oidcDiscoveryCaPem',
          'jwtValidationPubkeys',
          'jwtSupportedAlgs',
          'boundIssuer',
        ],
      };
      await login();
      return visit('/vault/settings/auth/enable');
    });
    hooks.afterEach(async function () {
      await runCmd(deleteAuthCmd(this.path), false);
    });
    testHelper(test);
  });

  module('ldap', function (hooks) {
    hooks.beforeEach(async function () {
      this.type = 'ldap';
      this.path = `${this.type}-${uuidv4()}`;
      this.tuneFields = [
        'url',
        'caseSensitiveNames',
        'connectionTimeout',
        'dereferenceAliases',
        'maxPageSize',
        'passwordPolicy',
        'requestTimeout',
        'tokenBoundCidrs',
        'tokenExplicitMaxTtl',
        'tokenMaxTtl',
        'tokenNoDefaultPolicy',
        'tokenNumUses',
        'tokenPeriod',
        'tokenPolicies',
        'tokenTtl',
        'tokenType',
        'usePre111GroupCnBehavior',
        'usernameAsAlias',
      ];
      this.tuneToggles = {
        'LDAP Options': [
          'starttls',
          'insecureTls',
          'discoverdn',
          'denyNullBind',
          'tlsMinVersion',
          'tlsMaxVersion',
          'certificate',
          'clientTlsCert',
          'clientTlsKey',
          'userattr',
          'upndomain',
          'anonymousGroupSearch',
        ],
        'Customize User Search': ['binddn', 'userdn', 'bindpass', 'userfilter'],
        'Customize Group Membership Search': ['groupfilter', 'groupattr', 'groupdn', 'useTokenGroups'],
      };
      await login();
      return visit('/vault/settings/auth/enable');
    });
    hooks.afterEach(async function () {
      await runCmd(deleteAuthCmd(this.path), false);
    });
    testHelper(test);
  });

  module('oidc', function (hooks) {
    hooks.beforeEach(async function () {
      this.type = 'oidc';
      this.path = `${this.type}-${uuidv4()}`;
      this.customSelectors = {
        providerConfig: `${GENERAL.fieldByAttr('providerConfig')} textarea`,
      };
      this.tuneFields = [
        'oidcDiscoveryUrl',
        'defaultRole',
        'jwksCaPem',
        'jwksUrl',
        'oidcResponseMode',
        'oidcResponseTypes',
        'namespaceInState',
        'providerConfig',
        'unsupportedCriticalCertExtensions',
      ];
      this.tuneToggles = {
        'OIDC Options': [
          'oidcClientId',
          'oidcClientSecret',
          'oidcDiscoveryCaPem',
          'jwtValidationPubkeys',
          'jwtSupportedAlgs',
          'boundIssuer',
        ],
      };
      await login();
      return visit('/vault/settings/auth/enable');
    });
    hooks.afterEach(async function () {
      await runCmd(deleteAuthCmd(this.path), false);
    });
    testHelper(test);
  });

  module('okta', function (hooks) {
    hooks.beforeEach(async function () {
      this.type = 'okta';
      this.path = `${this.type}-${uuidv4()}`;
      this.tuneFields = [
        'orgName',
        'tokenBoundCidrs',
        'tokenExplicitMaxTtl',
        'tokenMaxTtl',
        'tokenNoDefaultPolicy',
        'tokenNumUses',
        'tokenPeriod',
        'tokenPolicies',
        'tokenTtl',
        'tokenType',
      ];
      this.tuneToggles = { Options: ['apiToken', 'baseUrl', 'bypassOktaMfa'] };
      await login();
      return visit('/vault/settings/auth/enable');
    });
    hooks.afterEach(async function () {
      await runCmd(deleteAuthCmd(this.path), false);
    });
    testHelper(test);
  });
});
