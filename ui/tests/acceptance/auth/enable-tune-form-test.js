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
      'seal_wrap',
      'config.listing_visibility',
      'config.default_lease_ttl',
      'config.max_lease_ttl',
      'config.token_type',
      'config.audit_non_hmac_request_keys',
      'config.audit_non_hmac_response_keys',
      'config.passthrough_request_headers',
      'config.allowed_response_headers',
      'config.plugin_version',
    ];
    this.tokensGroup = {
      Tokens: [
        'token_bound_cidrs',
        'token_explicit_max_ttl',
        'token_max_ttl',
        'token_no_default_policy',
        'token_num_uses',
        'token_period',
        'token_policies',
        'token_ttl',
        'token_type',
      ],
    };
    this.oidcJwtGroup = {
      'OIDC/JWT Options': [
        'oidc_client_id',
        'oidc_client_secret',
        'oidc_discovery_ca_pem',
        'jwt_validation_pubkeys',
        'jwt_supported_algs',
        'bound_issuer',
      ],
    };
  });

  module('azure', function (hooks) {
    hooks.beforeEach(async function () {
      this.type = 'azure';
      this.path = `${this.type}-${uuidv4()}`;
      this.tuneFields = [
        'environment',
        'identity_token_audience',
        'identity_token_ttl',
        'max_retries',
        'max_retry_delay',
        'resource',
        'retry_delay',
        'root_password_ttl',
        'tenant_id',
      ];
      this.tuneToggles = { 'Azure Options': ['client_id', 'client_secret'] };
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
        provider_config: `${GENERAL.fieldByAttr('provider_config')} .cm-editor`,
      };
      this.tuneFields = [
        'default_role',
        'jwks_ca_pem',
        'jwks_url',
        'namespace_in_state',
        'oidc_discovery_url',
        'oidc_response_mode',
        'oidc_response_types',
        'provider_config',
        'unsupported_critical_cert_extensions',
      ];
      this.tuneToggles = this.oidcJwtGroup;
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
        'case_sensitive_names',
        'connection_timeout',
        'dereference_aliases',
        'max_page_size',
        'password_policy',
        'request_timeout',
        'use_pre111_group_cn_behavior',
        'username_as_alias',
      ];
      this.tuneToggles = {
        'LDAP Options': [
          'starttls',
          'insecure_tls',
          'discoverdn',
          'deny_null_bind',
          'tls_min_version',
          'tls_max_version',
          'certificate',
          'client_tls_cert',
          'client_tls_key',
          'userattr',
          'upndomain',
          'anonymous_group_search',
        ],
        'Customize User Search': ['binddn', 'userdn', 'bindpass', 'userfilter'],
        'Customize Group Membership Search': ['groupfilter', 'groupattr', 'groupdn', 'use_token_groups'],
        ...this.tokensGroup,
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
        provider_config: `${GENERAL.fieldByAttr('provider_config')} .cm-editor`,
      };
      this.tuneFields = [
        'oidc_discovery_url',
        'default_role',
        'jwks_ca_pem',
        'jwks_url',
        'oidc_response_mode',
        'oidc_response_types',
        'namespace_in_state',
        'provider_config',
        'unsupported_critical_cert_extensions',
      ];
      this.tuneToggles = this.oidcJwtGroup;
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
      this.tuneFields = ['org_name', 'api_token', 'base_url', 'bypass_okta_mfa'];
      this.tuneToggles = this.tokensGroup;
      await login();
      return visit('/vault/settings/auth/enable');
    });
    hooks.afterEach(async function () {
      await runCmd(deleteAuthCmd(this.path), false);
    });
    testHelper(test);
  });
});
