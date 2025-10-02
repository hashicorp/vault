/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: BUSL-1.1
 */

import { module, test } from 'qunit';
import { setupRenderingTest } from 'ember-qunit';
import { setupMirage } from 'ember-cli-mirage/test-support';
import { click, render } from '@ember/test-helpers';
import hbs from 'htmlbars-inline-precompile';
import { GENERAL } from 'vault/tests/helpers/general-selectors';
import sinon from 'sinon';
import Form from 'vault/forms/form';
import { overrideResponse } from 'vault/tests/helpers/stubs';
import AuthMethodResource from 'vault/resources/auth/method';

const TEST_CASES = [
  {
    methodType: 'aws',
    endpoint: '/auth/aws-test/config/client',
    section: 'client',
  },
  {
    methodType: 'aws',
    endpoint: '/auth/aws-test/config/tidy/identity-accesslist',
    section: 'identity-accesslist',
  },
  {
    methodType: 'aws',
    endpoint: '/auth/aws-test/config/tidy/roletag-denylist',
    section: 'roletag-denylist',
  },
  { methodType: 'azure', endpoint: '/auth/azure-test/config' },
  { methodType: 'github', endpoint: '/auth/github-test/config' },
  { methodType: 'gcp', endpoint: '/auth/gcp-test/config' },
  { methodType: 'jwt', endpoint: '/auth/jwt-test/config' },
  { methodType: 'oidc', endpoint: '/auth/oidc-test/config' },
  { methodType: 'kubernetes', endpoint: '/auth/kubernetes-test/config' },
  { methodType: 'ldap', endpoint: '/auth/ldap-test/config' },
  { methodType: 'okta', endpoint: '/auth/okta-test/config' },
  { methodType: 'radius', endpoint: '/auth/radius-test/config' },
];

module('Integration | Component | auth-config-form config', function (hooks) {
  setupRenderingTest(hooks);
  setupMirage(hooks);

  hooks.beforeEach(function () {
    this.flashSuccessSpy = sinon.spy(this.owner.lookup('service:flash-messages'), 'success');
    this.router = this.owner.lookup('service:router');
    this.transitionStub = sinon
      .stub(this.router, 'transitionTo')
      .returns({ followRedirects: () => Promise.resolve() });

    this.renderComponent = (path, methodType, section = 'configuration') => {
      // This component actually receives and OpenApiForm class, but this component test
      // is not testing the openapi parameters and just the form submit so we just use the Form class.
      this.form = new Form();
      this.section = section;
      this.method = new AuthMethodResource({ path, type: methodType }, this);
      return render(hbs`
        <AuthConfigForm::Config
         @form={{this.form}}
         @section={{this.section}}
         @method={{this.method}}
         />`);
    };
  });

  test.each('it makes save request to the expected endpoint', TEST_CASES, async function (assert, data) {
    assert.expect(3);
    const { methodType, endpoint, section } = data;
    const testId = methodType + `${section ? ` "${section}" section` : ''}`;
    this.server.post(endpoint, () => {
      assert.true(true, `${testId}: it calls expected endpoint "${endpoint}" on save.`);
    });
    await this.renderComponent(`${methodType}-test`, methodType, section);
    await click(GENERAL.submitButton);
    assert.true(this.flashSuccessSpy.calledOnce, `${testId}: flash success is called`);
    assert.true(this.transitionStub.calledOnce, `${testId}: transitionTo is called`);
  });

  test.each('it renders error banner', TEST_CASES, async function (assert, data) {
    assert.expect(3);
    const { methodType, endpoint, section } = data;
    const testId = methodType + `${section ? ` "${section}" section` : ''}`;
    this.server.post(endpoint, () => {
      return overrideResponse(400, { errors: ['uh oh'] });
    });
    await this.renderComponent(`${methodType}-test`, methodType, section);
    await click(GENERAL.submitButton);
    assert.false(this.flashSuccessSpy.calledOnce, `${testId}: flash success is NOT called`);
    assert.false(this.transitionStub.calledOnce, `${testId}: transitionTo is NOT called`);
    assert.dom(GENERAL.messageError).hasText('Error uh oh', `${testId}: it renders error message`);
  });
});
