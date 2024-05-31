/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: BUSL-1.1
 */

import { run } from '@ember/runloop';
import { resolve } from 'rsvp';
import Service from '@ember/service';
import { module, test } from 'qunit';
import { setupRenderingTest } from 'ember-qunit';
import { click, fillIn, render } from '@ember/test-helpers';
import hbs from 'htmlbars-inline-precompile';
import sinon from 'sinon';

const storeStub = Service.extend({
  callArgs: null,
  adapterFor() {
    return {
      replicationAction() {
        return {
          then(cb) {
            cb();
          },
        };
      },
    };
  },
});

const routerService = Service.extend({
  transitionTo: sinon.stub().returns(resolve()),
});

module('Integration | Component | replication actions', function (hooks) {
  setupRenderingTest(hooks);

  hooks.beforeEach(function () {
    run(() => {
      this.owner.register('service:router', routerService);
      this.owner.unregister('service:store');
      this.owner.register('service:store', storeStub);
      this.storeService = this.owner.lookup('service:store');
    });
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
    test(`replication mode ${replicationMode}, cluster mode: ${clusterMode}, action: ${action}`, async function (assert) {
      assert.expect(1);
      const testKey = `${replicationMode}-${clusterMode}-${action}`;
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
          `${testKey}: submitted values match expected`
        );
        return resolve();
      });
      this.set('storeService.capabilitiesReturnVal', ['root']);
      await render(
        hbs`
                <ReplicationActions
          @model={{this.model}}
          @replicationMode={{this.replicationMode}}
          @selectedAction={{this.selectedAction}}
          @onSubmit={{action this.onSubmit}}
        />
        `
      );
      assert
        .dom(`[data-test-${action}-replication] h3`)
        .hasText(headerText, `${testKey}: renders the ${action} component header`);

      await click(`[data-test-replication-action-trigger="${action}"]`);
      if (typeof fillInFn === 'function') {
        await fillInFn.call(this);
      }
      await click('[data-test-confirm-button]');
    });
  }
});
