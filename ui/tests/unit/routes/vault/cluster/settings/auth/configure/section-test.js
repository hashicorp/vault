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

  module('options', function (hooks) {
    hooks.beforeEach(function () {
      this.optionsFor = (config) => {
        this.modelForStub.returns({
          method: { id: 'userpass/' },
          methodOptions: { type: 'userpass', config: { default_lease_ttl: 0, ...config } },
        });
        return this.route.modelForOptions().form;
      };
    });

    test('it shows the lockout settings saved on the mount', function (assert) {
      const form = this.optionsFor({
        user_lockout_config: {
          user_lockout_threshold: 7,
          user_lockout_duration: 600,
          user_lockout_counter_reset_duration: 1800,
          user_lockout_disable: true,
        },
      });

      assert.propEqual(
        form.data.user_lockout_config,
        {
          lockout_threshold: '7',
          lockout_duration: '600s',
          lockout_counter_reset_duration: '1800s',
          lockout_disable: true,
        },
        'saved settings use the keys the form sends back'
      );
    });

    test('it leaves unset lockout settings empty', function (assert) {
      const form = this.optionsFor({
        user_lockout_config: {
          user_lockout_threshold: 3,
          user_lockout_duration: 0,
          user_lockout_counter_reset_duration: 0,
          user_lockout_disable: false,
        },
      });

      assert.propEqual(
        form.data.user_lockout_config,
        { lockout_threshold: '3', lockout_disable: false },
        'zero durations are not turned into a configured value'
      );
    });

    test('it does not send back a threshold too large to read exactly', function (assert) {
      // 18446744073709551615 is a valid threshold, but JSON parsing rounds it to this number
      const form = this.optionsFor({
        user_lockout_config: { user_lockout_threshold: 18446744073709552000, user_lockout_duration: 600 },
      });

      assert.propEqual(
        form.data.user_lockout_config,
        { lockout_duration: '600s' },
        'the rounded threshold is left out so saving does not change or reject it'
      );
    });

    test('it starts empty when the mount has no lockout settings', function (assert) {
      const form = this.optionsFor({});
      assert.propEqual(form.data.user_lockout_config, {}, 'nothing is pre-filled or sent');
    });

    test('it does not pass the saved lockout object on as mount config', function (assert) {
      const form = this.optionsFor({ user_lockout_config: { user_lockout_threshold: 7 } });
      assert.false('user_lockout_config' in form.data.config, 'the read-only shape is not resubmitted');
    });
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

  test('it should include token_reviewer_jwt field for kubernetes auth method', async function (assert) {
    const { form } = await this.testModelForConfiguration('kubernetes', 'token_reviewer_jwt');
    // token_reviewer_jwt lives in a non-default field group, so search across all groups
    const allFields = form.formFieldGroups.flatMap((group) => Object.values(group)[0] || []);
    const field = allFields.find((f) => f.name === 'token_reviewer_jwt');
    assert.ok(field, 'token_reviewer_jwt field is present in the kubernetes auth method form');
  });
});
