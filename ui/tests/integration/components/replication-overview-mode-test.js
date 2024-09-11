/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: BUSL-1.1
 */

import { module, test } from 'qunit';
import { setupRenderingTest } from 'vault/tests/helpers';
import { render, settled } from '@ember/test-helpers';
import { hbs } from 'ember-cli-htmlbars';
import { setupEngine } from 'ember-engines/test-support';

const OVERVIEW_MODE = {
  title: '[data-test-overview-mode-title]',
  body: '[data-test-overview-mode-body]',
  detailsLink: '[data-test-replication-details-link]',
};
module('Integration | Component | replication-overview-mode', function (hooks) {
  setupRenderingTest(hooks);
  setupEngine(hooks, 'replication');

  hooks.beforeEach(function () {
    this.versionService = this.owner.lookup('service:version');
    this.versionService.features = [];
    this.mode = 'dr';
    this.clusterName = 'foobar';
    this.modeDetails = { mode: 'disabled' };

    this.renderComponent = async () => {
      return render(
        hbs`
        <ReplicationOverviewMode
          @clusterName={{this.clusterName}}
          @mode={{this.mode}}
          @model={{this.modeDetails}}
        />`,
        { owner: this.engine }
      );
    };
  });

  test('without features', async function (assert) {
    await this.renderComponent();
    assert.dom(OVERVIEW_MODE.title).hasText('Disaster Recovery (DR)');
    assert
      .dom(OVERVIEW_MODE.body)
      .includesText('Disaster Recovery is a feature of Vault Enterprise Premium. Upgrade');
    assert.dom(OVERVIEW_MODE.detailsLink).doesNotExist('does not show link to replication (dr)');

    this.set('mode', 'performance');
    await settled();
    assert.dom(OVERVIEW_MODE.title).hasText('Performance');
    assert
      .dom(OVERVIEW_MODE.body)
      .includesText('Performance Replication is a feature of Vault Enterprise Premium. Upgrade');
    assert.dom(OVERVIEW_MODE.detailsLink).doesNotExist('does not show link to replication (perf)');
  });

  module('with features', function (hooks) {
    hooks.beforeEach(function () {
      this.versionService.features = ['DR Replication', 'Performance Replication'];
    });

    test('it renders when replication disabled', async function (assert) {
      await this.renderComponent();
      assert.dom(OVERVIEW_MODE.title).hasText('Disaster Recovery (DR)');
      assert
        .dom(OVERVIEW_MODE.body)
        .hasText(
          'Disaster Recovery Replication is designed to protect against catastrophic failure of entire clusters. Secondaries do not forward service requests until they are elected and become a new primary.'
        );
      assert.dom(OVERVIEW_MODE.detailsLink).hasText('Enable');

      this.set('mode', 'performance');
      await settled();
      assert.dom(OVERVIEW_MODE.title).hasText('Performance');
      assert
        .dom(OVERVIEW_MODE.body)
        .hasText(
          'Performance Replication scales workloads horizontally across clusters to make requests faster. Local secondaries handle read requests but forward writes to the primary to be handled.'
        );
      assert.dom(OVERVIEW_MODE.detailsLink).hasText('Enable');
    });

    test('it renders when replication enabled', async function (assert) {
      this.mode = 'performance';
      this.modeDetails = {
        replicationEnabled: true,
        mode: 'primary',
        modeForUrl: 'primary',
        clusterIdDisplay: 'foobar12',
      };
      await this.renderComponent();
      assert.dom(OVERVIEW_MODE.title).hasText('Performance');
      assert
        .dom(OVERVIEW_MODE.body)
        .includesText('ENABLED Primary foobar12', 'renders mode type and cluster ID if passed');
      assert.dom(OVERVIEW_MODE.detailsLink).hasText('Details');

      this.set('modeDetails', {
        replicationEnabled: true,
        mode: 'secondary',
        modeForUrl: 'secondary',
        clusterIdDisplay: 'foobar12',
        secondaryId: 'some-secondary',
      });
      await settled();
      assert.dom(OVERVIEW_MODE.title).hasText('Performance');
      assert.dom(OVERVIEW_MODE.body).includesText('ENABLED Secondary some-secondary foobar12');
      assert.dom(OVERVIEW_MODE.detailsLink).hasText('Details');
    });

    test('it renders when replication bootstrapping', async function (assert) {
      this.modeDetails = {
        replicationEnabled: true,
        mode: 'bootstrapping',
        modeForUrl: 'bootstrapping',
      };
      await this.renderComponent();
      assert.dom(OVERVIEW_MODE.title).hasText('Disaster Recovery (DR)');
      assert.dom(OVERVIEW_MODE.body).includesText('ENABLED Bootstrapping');
      assert.dom(OVERVIEW_MODE.detailsLink).hasText('Details');
    });
  });
});
