/**
 * Copyright IBM Corp. 2016, 2025
 * SPDX-License-Identifier: BUSL-1.1
 */

import { resolve } from 'rsvp';
import { module, test } from 'qunit';
import { setupRenderingTest } from 'vault/tests/helpers';
import { click, fillIn, render } from '@ember/test-helpers';
import { GENERAL } from 'vault/tests/helpers/general-selectors';
import hbs from 'htmlbars-inline-precompile';
import sinon from 'sinon';

module('Integration | Component | replication actions', function (hooks) {
  setupRenderingTest(hooks);

  hooks.beforeEach(function () {
    const { sys } = this.owner.lookup('service:api');
    sinon.stub(sys, 'systemWriteReplicationDrPrimaryDisable').resolves();
    sinon.stub(sys, 'systemWriteReplicationDrSecondaryDisable').resolves();
    sinon.stub(sys, 'systemWriteReplicationPerformancePrimaryDisable').resolves();
    sinon.stub(sys, 'systemWriteReplicationPerformanceSecondaryDisable').resolves();
    sinon.stub(sys, 'systemWriteReplicationDrPrimaryDemote').resolves();
    sinon.stub(sys, 'systemWriteReplicationPerformancePrimaryDemote').resolves();
    sinon.stub(sys, 'systemWriteReplicationDrSecondaryPromote').resolves();
    sinon.stub(sys, 'systemWriteReplicationPerformanceSecondaryPromote').resolves();
    sinon.stub(sys, 'systemWriteReplicationDrSecondaryUpdatePrimary').resolves();
    sinon.stub(sys, 'systemWriteReplicationPerformanceSecondaryUpdatePrimary').resolves();
  });

  const confirmInput = (confirmText) => fillIn('[data-test-confirmation-modal-input]', confirmText);
  const testCases = [
    [
      'dr',
      'primary',
      'disable',
      'Disable Replication',
      () => confirmInput('Disaster Recovery'),
      ['disable', 'primary'],
      false,
    ],
    [
      'performance',
      'primary',
      'disable',
      'Disable Replication',
      () => confirmInput('Performance'),
      ['disable', 'primary'],
      false,
    ],
    [
      'performance',
      'secondary',
      'disable',
      'Disable Replication',
      () => confirmInput('Performance'),
      ['disable', 'secondary'],
      false,
    ],
    ['dr', 'primary', 'recover', 'Recover', null, ['recover'], false],
    ['performance', 'primary', 'recover', 'Recover', null, ['recover'], false],
    ['performance', 'secondary', 'recover', 'Recover', null, ['recover'], false],

    ['dr', 'primary', 'reindex', 'Reindex', () => null, ['reindex'], false],
    ['performance', 'primary', 'reindex', 'Reindex', () => null, ['reindex'], false],
    ['performance', 'secondary', 'reindex', 'Reindex', () => null, ['reindex'], false],

    [
      'dr',
      'primary',
      'demote',
      'Demote cluster',
      () => confirmInput('Disaster Recovery'),
      ['demote', 'primary'],
      false,
    ],
    [
      'performance',
      'primary',
      'demote',
      'Demote cluster',
      () => confirmInput('Performance'),
      ['demote', 'primary'],
      false,
    ],

    // we don't do dr secondary promote in this component so just test perf
    // re-enable this test when the DR secondary disable API endpoint is fixed
    // ['dr', 'secondary', 'disable', 'Disable Replication', null, ['disable', 'secondary'], false],
    // ['dr', 'secondary', 'reindex', 'Reindex', null, ['reindex'], false],
    [
      'performance',
      'secondary',
      'promote',
      'Promote cluster',
      async function () {
        await fillIn('[name="primary_cluster_addr"]', 'cluster addr');
      },
      ['promote', 'secondary', { primary_cluster_addr: 'cluster addr' }],
      false,
    ],
    [
      'performance',
      'secondary',
      'update-primary',
      'Update primary',
      async function () {
        await fillIn('#secondary-token', 'token');
        await fillIn('#primary_api_addr', 'addr');
      },
      ['update-primary', 'secondary', { token: 'token', primary_api_addr: 'addr' }],
      false,
    ],
  ];

  for (const [replicationMode, clusterMode, action, headerText, fillInFn, expectedOnSubmit] of testCases) {
    const testCaseKey = `${replicationMode}-${clusterMode}-${action}`;
    test(`replication and cluster mode action behavior: testCaseKey = ${testCaseKey}`, async function (assert) {
      assert.expect(1);
      this.set('model', {
        replicationAttrs: {
          modeForUrl: clusterMode,
        },
        [replicationMode]: {
          mode: clusterMode,
          modeForUrl: clusterMode,
        },
        replicationModeForDisplay: replicationMode === 'dr' ? 'Disaster Recovery' : 'Performance',
        reload() {
          return resolve();
        },
        rollbackAttributes() {},
      });
      this.set('replicationMode', replicationMode);
      this.set('selectedAction', action);
      this.set('onSubmit', (...actual) => {
        assert.deepEqual(
          JSON.stringify(actual),
          JSON.stringify(expectedOnSubmit),
          `Submitted values match what is expected for the testCaseKey: ${testCaseKey}`
        );
        return resolve();
      });

      await render(
        hbs`<ReplicationActions
          @model={{this.model}}
          @replicationMode={{this.replicationMode}}
          @selectedAction={{this.selectedAction}}
          @onSubmit={{action this.onSubmit}}
        />`
      );
      assert
        .dom(`[data-test-${action}-replication] h3`)
        .hasText(headerText, `Renders the component header for: ${testCaseKey}`);
      await click(GENERAL.button(action));
      if (typeof fillInFn === 'function') {
        await fillInFn.call(this);
      }
      await click(GENERAL.confirmButton);
    });
  }
});
