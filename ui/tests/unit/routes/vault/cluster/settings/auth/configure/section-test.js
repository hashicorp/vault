/**
 * Copyright IBM Corp. 2016, 2025
 * SPDX-License-Identifier: BUSL-1.1
 */

import { module, test } from 'qunit';
import { setupTest } from 'vault/tests/helpers';
import sinon from 'sinon';

module('Unit | Route | vault/cluster/settings/auth/configure/section', function (hooks) {
  setupTest(hooks);

  hooks.beforeEach(function () {
    const { auth } = this.owner.lookup('service:api');
    sinon.stub(auth, 'kubernetesReadAuthConfiguration').resolves({});
    sinon.stub(auth, 'jwtReadConfiguration').resolves({});

    this.route = this.owner.lookup('route:vault/cluster/settings/auth/configure/section');
    this.modelForStub = sinon.stub(this.route, 'modelFor');

    this.testModelForConfiguration = async (methodType, fieldKey) => {
      this.modelForStub.returns({ method: { methodType, path: `${methodType}-test` } });
      const { form } = await this.route.modelForConfiguration('configuration');
      const defaultGroup = form.formFieldGroups[0]['default'];
      const field = defaultGroup.find((field) => field.name === fieldKey);
      return { form, defaultGroup, field };
    };
  });

  test('it should remove jwks_pairs form field for jwt type', async function (assert) {
    const { field } = await this.testModelForConfiguration('jwt', 'jwks_pairs');
    assert.strictEqual(field, undefined, 'jwks_pairs field is removed for jwt type');
  });

  test('it should remove jwks_pairs form field for oidc type', async function (assert) {
    const { field } = await this.testModelForConfiguration('oidc', 'jwks_pairs');
    assert.strictEqual(field, undefined, 'jwks_pairs field is removed for oidc type');
  });

  test('it should update kubernetes_ca_cert form field editType to file', async function (assert) {
    const { field } = await this.testModelForConfiguration('kubernetes', 'kubernetes_ca_cert');
    assert.strictEqual(
      field.options.editType,
      'file',
      'editType is set to file for kubernetes_ca_cert field'
    );
  });
});
