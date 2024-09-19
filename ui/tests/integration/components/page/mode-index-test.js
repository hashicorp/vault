/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: BUSL-1.1
 */

import { module, test } from 'qunit';
import { setupRenderingTest } from 'vault/tests/helpers';
import { render } from '@ember/test-helpers';
import { hbs } from 'ember-cli-htmlbars';
import { setupEngine } from 'ember-engines/test-support';

const S = {
  title: 'h1',
  subtitle: 'h2',
  enableForm: '[data-test-replication-enable-form]',
  enableBtn: '[data-test-replication-enable]',
  summary: '[data-test-replication-summary]',
  notAllowed: '[data-test-not-allowed]',
};
module('Integration | Component | replication page/mode-index', function (hooks) {
  setupRenderingTest(hooks);
  setupEngine(hooks, 'replication');

  hooks.beforeEach(function () {
    this.store = this.owner.lookup('service:store');
    this.onEnable = () => {};
    this.clusterModel = { replicationAttrs: {} };
    this.replicationMode = '';
    this.replicationDisabled = true;

    this.renderComponent = () => {
      return render(
        hbs`<Page::ModeIndex @replicationDisabled={{this.replicationDisabled}} @replicationMode={{this.replicationMode}} @cluster={{this.clusterModel}} @onEnableSuccess={{this.onEnable}} />`,
        { owner: this.engine }
      );
    };
  });

  module('DR mode', function (hooks) {
    hooks.beforeEach(function () {
      this.replicationMode = 'dr';
    });
    test('it renders correctly when replication disabled', async function (assert) {
      await this.renderComponent();

      assert.dom(S.title).hasText('Enable Disaster Recovery Replication');
      assert.dom(S.enableForm).exists();
      assert.dom(S.notAllowed).doesNotExist();
      assert.dom(S.enableBtn).exists('Enable button shows by default if no permissions available');
    });
    test('it renders correctly when replication enabled', async function (assert) {
      this.replicationDisabled = false;
      await this.renderComponent();

      assert.dom(S.enableForm).doesNotExist();
      assert.dom(S.summary).exists();
    });

    test('it hides enable button if no permissions', async function (assert) {
      this.clusterModel.canEnablePrimaryDr = false;
      await this.renderComponent();

      assert.dom(S.enableForm).exists();
      assert.dom(S.notAllowed).exists();
      assert.dom(S.enableBtn).doesNotExist();
    });

    test('it shows enable button if has permissions', async function (assert) {
      this.clusterModel.canEnablePrimaryDr = true;
      await this.renderComponent();

      assert.dom(S.enableForm).exists();
      assert.dom(S.notAllowed).doesNotExist();
      assert.dom(S.enableBtn).exists();
    });
  });

  module('Performance mode', function (hooks) {
    hooks.beforeEach(function () {
      this.replicationMode = 'performance';
    });
    test('it renders correctly when replication disabled', async function (assert) {
      await this.renderComponent();

      assert.dom(S.title).hasText('Enable Performance Replication');
      assert.dom(S.enableForm).exists();
      assert.dom(S.notAllowed).doesNotExist();
      assert.dom(S.enableBtn).exists('Enable button shows by default if no permissions available');
    });
    test('it renders correctly when replication enabled', async function (assert) {
      this.replicationDisabled = false;
      await this.renderComponent();

      assert.dom(S.enableForm).doesNotExist();
      assert.dom(S.summary).exists();
    });

    test('it hides enable button if no permissions', async function (assert) {
      this.clusterModel.canEnablePrimaryPerformance = false;
      await this.renderComponent();

      assert.dom(S.enableForm).exists();
      assert.dom(S.notAllowed).exists();
      assert.dom(S.enableBtn).doesNotExist();
    });

    test('it shows enable button if has permissions', async function (assert) {
      this.clusterModel.canEnablePrimaryPerformance = true;
      await this.renderComponent();

      assert.dom(S.enableForm).exists();
      assert.dom(S.notAllowed).doesNotExist();
      assert.dom(S.enableBtn).exists();
    });
  });
});
