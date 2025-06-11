/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: BUSL-1.1
 */

import { module, test } from 'qunit';
import { setupRenderingTest } from 'ember-qunit';
import { setupMirage } from 'ember-cli-mirage/test-support';
import { setupTotpMfaResponse } from 'vault/tests/helpers/mfa/mfa-helpers';
import { mfaTests, setupTestContext } from './test-helper';
import { ERROR_JWT_LOGIN } from 'vault/components/auth/form/oidc-jwt';
import { overrideResponse } from 'vault/tests/helpers/stubs';
import sinon from 'sinon';
import { windowStub } from 'vault/tests/helpers/oidc-window-stub';

module('Integration | Component | auth | page | mfa', function (hooks) {
  setupRenderingTest(hooks);
  setupMirage(hooks);

  hooks.beforeEach(async function () {
    setupTestContext(this);
  });

  module('github', function (hooks) {
    hooks.beforeEach(async function () {
      this.authType = 'github';
      this.loginData = { token: 'mysupersecuretoken' };
      this.path = this.authType;
      this.stubRequests = () => {
        this.server.post(`/auth/${this.path}/login`, () => setupTotpMfaResponse(this.path));
      };
    });

    mfaTests(test);
  });

  module('jwt', function (hooks) {
    hooks.beforeEach(async function () {
      this.authType = 'jwt';
      this.loginData = { role: 'some-dev', jwt: 'jwttoken' };
      this.path = this.authType;
      this.routerStub = sinon.stub(this.owner.lookup('service:router'), 'urlFor').returns('123-example.com');

      this.stubRequests = () => {
        this.server.post('/auth/:path/oidc/auth_url', () =>
          overrideResponse(400, { errors: [ERROR_JWT_LOGIN] })
        );
        this.server.post(`/auth/${this.path}/login`, () => setupTotpMfaResponse(this.path));
      };
    });

    mfaTests(test);
  });

  module('oidc', function (hooks) {
    hooks.beforeEach(async function () {
      this.authType = 'oidc';
      this.loginData = { role: 'some-dev' };
      this.path = this.authType;
      // Requests are stubbed in the order they are hit
      this.stubRequests = () => {
        this.server.post(`/auth/${this.path}/oidc/auth_url`, () => {
          return { data: { auth_url: 'http://dev-foo-bar.com' } };
        });
        this.server.get(`/auth/${this.path}/oidc/callback`, () => setupTotpMfaResponse(this.path));
      };

      // additional OIDC setup
      this.routerStub = sinon.stub(this.owner.lookup('service:router'), 'urlFor').returns('123-example.com');
      this.windowStub = windowStub();
    });

    hooks.afterEach(function () {
      this.routerStub.restore();
      this.windowStub.restore();
    });

    mfaTests(test);
  });

  module('username and password methods', function (hooks) {
    hooks.beforeEach(async function () {
      this.loginData = { username: 'matilda', password: 'password' };
      this.stubRequests = () => {
        this.server.post(`/auth/${this.path}/login/matilda`, () => setupTotpMfaResponse(this.path));
      };
    });

    module('ldap', function (hooks) {
      hooks.beforeEach(async function () {
        this.authType = 'ldap';
        this.path = this.authType;
      });

      mfaTests(test);
    });

    module('okta', function (hooks) {
      hooks.beforeEach(async function () {
        this.authType = 'okta';
        this.path = this.authType;
      });

      mfaTests(test);
    });

    module('radius', function (hooks) {
      hooks.beforeEach(async function () {
        this.authType = 'radius';
        this.path = this.authType;
      });

      mfaTests(test);
    });

    module('userpass', function (hooks) {
      hooks.beforeEach(async function () {
        this.authType = 'userpass';
        this.path = this.authType;
      });

      mfaTests(test);
    });
  });

  // ENTERPRISE METHODS
  module('saml', function (hooks) {
    hooks.beforeEach(async function () {
      this.version.type = 'enterprise';
      this.authType = 'saml';
      this.path = this.authType;
      this.loginData = { role: 'some-dev' };
      // Requests are stubbed in the order they are hit
      this.stubRequests = () => {
        this.server.put(`/auth/${this.path}/sso_service_url`, () => ({
          data: {
            sso_service_url: 'test/fake/sso/route',
            token_poll_id: '1234',
          },
        }));
        this.server.put(`/auth/${this.path}/token`, () => setupTotpMfaResponse(this.authType));
      };
      this.windowStub = windowStub();
    });

    hooks.afterEach(function () {
      this.windowStub.restore();
    });

    mfaTests(test);
  });
});
