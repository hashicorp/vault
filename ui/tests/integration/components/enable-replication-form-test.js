/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: BUSL-1.1
 */

import sinon from 'sinon';
import { module, test } from 'qunit';
import { setupRenderingTest } from 'vault/tests/helpers';
import { setupEngine } from 'ember-engines/test-support';
import { render, fillIn, click, settled } from '@ember/test-helpers';
import { hbs } from 'ember-cli-htmlbars';
import { setupMirage } from 'ember-cli-mirage/test-support';
import { overrideResponse } from 'vault/tests/helpers/stubs';

const ENABLE_FORM = {
  clusterMode: '[data-test-replication-cluster-mode-select]',
  clusterAddr: '[data-test-input="primary_cluster_addr"]',
  secondaryToken: '[data-test-textarea="secondary-token"]',
  primaryAddr: '[data-test-input="primary_api_addr"]',
  caFile: '[data-test-input="ca_file"]',
  caPath: '[data-test-input="ca_path"]',
  submitButton: '[data-test-replication-enable]',
  notAllowed: '[data-test-not-allowed]',
  inlineMessage: '[data-test-inline-error-message]',
  cannotEnable: '[data-test-disable-to-continue]',
  cannotEnableExplanation: '[data-test-disable-explanation]',
  error: '[data-test-message-error-description]',
};
module('Integration | Component | enable-replication-form', function (hooks) {
  setupRenderingTest(hooks);
  setupMirage(hooks);
  setupEngine(hooks, 'replication');

  hooks.beforeEach(function () {
    this.context = { owner: this.engine };
    this.version = this.owner.lookup('service:version');
  });

  ['performance', 'dr'].forEach((replicationMode) => {
    test(`it renders correct form inputs when ${replicationMode} replication mode`, async function (assert) {
      assert.expect(10);
      this.version.features = ['Performance Replication', 'DR Replication'];
      this.set('replicationMode', replicationMode);
      await render(
        hbs`<EnableReplicationForm
      @replicationMode={{this.replicationMode}}
      @canEnablePrimary={{true}}
      @canEnableSecondary={{true}}
      @performanceMode="disabled"
    />`,
        this.context
      );

      assert.dom(ENABLE_FORM.clusterMode).hasValue('primary');
      ['clusterAddr'].forEach((field) => {
        assert.dom(ENABLE_FORM[field]).hasNoValue();
      });
      assert.dom(ENABLE_FORM.submitButton).isNotDisabled();

      await fillIn(ENABLE_FORM.clusterMode, 'secondary');
      assert.dom(ENABLE_FORM.inlineMessage).hasText('This will immediately clear all data in this cluster!');
      ['secondaryToken', 'primaryAddr', 'caFile', 'caPath'].forEach((field) => {
        assert.dom(ENABLE_FORM[field]).hasNoValue();
      });
      assert.dom(ENABLE_FORM.submitButton).isDisabled();
      await fillIn(ENABLE_FORM.secondaryToken, 'some-token');
      await fillIn(ENABLE_FORM.primaryAddr, 'some-addr');
      assert.dom(ENABLE_FORM.submitButton).isNotDisabled();
    });
    test(`it shows warning when capabilities restricted for ${replicationMode} replication mode`, async function (assert) {
      assert.expect(10);
      this.version.features = ['Performance Replication', 'DR Replication'];
      this.set('replicationMode', replicationMode);
      await render(
        hbs`<EnableReplicationForm
          @replicationMode={{this.replicationMode}}
          @canEnablePrimary={{false}}
          @canEnableSecondary={{false}}
          @performanceMode="disabled"
        />`,
        this.context
      );
      assert.dom(ENABLE_FORM.clusterMode).hasValue('primary');
      assert
        .dom(ENABLE_FORM.notAllowed)
        .hasText('The token you are using is not authorized to enable primary replication.');
      ['clusterAddr', 'submitButton'].forEach((field) => {
        assert.dom(ENABLE_FORM[field]).doesNotExist();
      });

      await fillIn(ENABLE_FORM.clusterMode, 'secondary');
      assert
        .dom(ENABLE_FORM.notAllowed)
        .hasText('The token you are using is not authorized to enable secondary replication.');
      ['secondaryToken', 'primaryAddr', 'caFile', 'caPath', 'submitButton'].forEach((field) => {
        assert.dom(ENABLE_FORM[field]).doesNotExist();
      });
    });
  });

  test('enable DR when cluster is perf primary', async function (assert) {
    this.version.features = ['Performance Replication', 'DR Replication'];
    this.set('replicationMode', 'dr');
    this.set('performanceMode', 'primary');
    await render(
      hbs`<EnableReplicationForm
        @replicationMode={{this.replicationMode}}
        @canEnablePrimary={{true}}
        @canEnableSecondary={{true}}
        @performanceMode={{this.performanceMode}}
      />`,
      this.context
    );
    assert.dom(ENABLE_FORM.clusterMode).hasValue('primary');
    ['clusterAddr'].forEach((field) => {
      assert.dom(ENABLE_FORM[field]).hasNoValue();
    });
    assert.dom(ENABLE_FORM.submitButton).isNotDisabled();

    await fillIn(ENABLE_FORM.clusterMode, 'secondary');
    assert
      .dom(ENABLE_FORM.cannotEnable)
      .hasText('Disable Performance Replication in order to enable this cluster as a DR secondary.');
    await click(ENABLE_FORM.cannotEnable);
    assert
      .dom(ENABLE_FORM.cannotEnableExplanation)
      .hasText(
        "When running as a DR Secondary Vault is read only. For this reason, we don't allow other Replication modes to operate at the same time. This cluster is also currently operating as a Performance Primary."
      );
    assert.dom(ENABLE_FORM.submitButton).isDisabled();

    this.set('performanceMode', 'secondary');
    await settled();
    assert
      .dom(ENABLE_FORM.cannotEnableExplanation)
      .hasText(
        "When running as a DR Secondary Vault is read only. For this reason, we don't allow other Replication modes to operate at the same time. This cluster is also currently operating as a Performance Secondary."
      );
  });

  module('only DR replication in features', function (hooks) {
    hooks.beforeEach(function () {
      this.version.features = ['DR Replication'];
    });
    test('attempting to enable performance replication', async function (assert) {
      await render(
        hbs`<EnableReplicationForm
          @replicationMode="performance"
          @canEnablePrimary={{true}}
          @canEnableSecondary={{true}}
          @performanceMode="disabled"
        />`,
        this.context
      );
      assert.dom(ENABLE_FORM.submitButton).isDisabled();
    });
  });

  module('successful enable', function (hooks) {
    hooks.beforeEach(function () {
      this.version.features = ['Performance Replication', 'DR Replication'];
      this.successSpy = sinon.spy();
      this.set('onSuccess', this.successSpy);
    });
    ['dr', 'performance'].forEach((replicationMode) => {
      test(`${replicationMode} primary`, async function (assert) {
        assert.expect(4);
        this.set('replicationMode', replicationMode);
        this.server.post(`/sys/replication/${replicationMode}/primary/enable`, (_, req) => {
          const body = JSON.parse(req.requestBody);
          assert.deepEqual(body, {
            primary_cluster_addr: 'some-addr',
          });
          return {
            returned: 'value',
          };
        });
        await render(
          hbs`<EnableReplicationForm
            @replicationMode={{this.replicationMode}}
            @canEnablePrimary={{true}}
            @canEnableSecondary={{true}}
            @performanceMode="disabled"
            @onSuccess={{this.onSuccess}}
            @doTransition={{false}}
          />`,
          this.context
        );
        await fillIn(ENABLE_FORM.clusterAddr, 'some-addr');
        await click(ENABLE_FORM.submitButton);
        // after success
        assert.dom(ENABLE_FORM.clusterAddr).hasNoValue();
        assert.true(this.successSpy.calledOnce, 'called once');
        assert.deepEqual(
          this.successSpy.getCall(0).args,
          [{ returned: 'value' }, replicationMode, 'primary', false],
          'called with correct args'
        );
      });
      test(`${replicationMode} secondary`, async function (assert) {
        assert.expect(5);
        this.set('replicationMode', replicationMode);
        this.server.post(`/sys/replication/${replicationMode}/secondary/enable`, (_, req) => {
          const body = JSON.parse(req.requestBody);
          assert.deepEqual(
            body,
            {
              primary_api_addr: 'http://127.0.0.1:8200',
              token: 'some-token-value',
            },
            'does not include empty values'
          );
          return {
            returned: 'value',
          };
        });
        await render(
          hbs`<EnableReplicationForm
          @replicationMode={{this.replicationMode}}
            @canEnablePrimary={{true}}
            @canEnableSecondary={{true}}
            @performanceMode="disabled"
            @onSuccess={{this.onSuccess}}
            @doTransition={{true}}
          />`,
          this.context
        );
        await fillIn(ENABLE_FORM.clusterMode, 'secondary');
        await fillIn(ENABLE_FORM.secondaryToken, 'some-token-value');
        await fillIn(ENABLE_FORM.primaryAddr, 'http://127.0.0.1:8200');
        // Fill in then clear ca path
        await fillIn(ENABLE_FORM.caPath, 'some-path');
        await fillIn(ENABLE_FORM.caPath, '');
        await click(ENABLE_FORM.submitButton);
        // after success
        assert.dom(ENABLE_FORM.secondaryToken).hasValue('');
        assert.dom(ENABLE_FORM.primaryAddr).hasNoValue();
        assert.true(this.successSpy.calledOnce, 'called once');
        assert.deepEqual(
          this.successSpy.getCall(0).args,
          [{ returned: 'value' }, replicationMode, 'secondary', true],
          'called with correct args'
        );
      });
    });
  });

  module('shows API errors', function (hooks) {
    hooks.beforeEach(function () {
      this.version.features = ['Performance Replication', 'DR Replication'];
      this.successSpy = sinon.spy();
      this.set('onSuccess', this.successSpy);
    });
    ['dr', 'performance'].forEach((replicationMode) => {
      test(`${replicationMode} primary`, async function (assert) {
        this.set('replicationMode', replicationMode);
        this.server.post(`/sys/replication/${replicationMode}/primary/enable`, overrideResponse(403));
        await render(
          hbs`<EnableReplicationForm
            @replicationMode={{this.replicationMode}}
            @canEnablePrimary={{true}}
            @canEnableSecondary={{true}}
            @performanceMode="disabled"
            @onSuccess={{this.onSuccess}}
          />`,
          this.context
        );
        await fillIn(ENABLE_FORM.clusterAddr, 'some-addr');
        await click(ENABLE_FORM.submitButton);
        assert.dom(ENABLE_FORM.error).hasText('permission denied', 'shows error returned from API');
        assert.dom(ENABLE_FORM.clusterAddr).hasValue('some-addr', 'does not clear form');
        assert.false(this.successSpy.calledOnce, 'success spy not called');
      });
      test(`${replicationMode} secondary`, async function (assert) {
        this.set('replicationMode', replicationMode);
        this.server.post(`/sys/replication/${replicationMode}/secondary/enable`, overrideResponse(403));
        await render(
          hbs`<EnableReplicationForm
          @replicationMode={{this.replicationMode}}
            @canEnablePrimary={{true}}
            @canEnableSecondary={{true}}
            @performanceMode="disabled"
            @onSuccess={{this.onSuccess}}
          />`,
          this.context
        );
        await fillIn(ENABLE_FORM.clusterMode, 'secondary');
        await fillIn(ENABLE_FORM.secondaryToken, 'some-token-value');
        await fillIn(ENABLE_FORM.primaryAddr, 'http://127.0.0.1:8200');
        await click(ENABLE_FORM.submitButton);
        // after error
        assert.dom(ENABLE_FORM.error).hasText('permission denied', 'shows error returned from API');
        assert.dom(ENABLE_FORM.secondaryToken).hasValue('some-token-value', 'does not clear form');
        assert.false(this.successSpy.calledOnce, 'success spy not called');
      });
    });
  });
});
