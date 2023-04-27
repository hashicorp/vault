/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: MPL-2.0
 */

import { run } from '@ember/runloop';
import { resolve } from 'rsvp';
import Service from '@ember/service';
import { module, test } from 'qunit';
import { setupRenderingTest } from 'ember-qunit';
import { click, fillIn, blur, render } from '@ember/test-helpers';
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
    ['dr', 'primary', 'recover', 'Recover', () => confirmInput('Disaster Recovery'), ['recover'], false],
    ['performance', 'primary', 'recover', 'Recover', () => confirmInput('Performance'), ['recover'], false],
    ['performance', 'secondary', 'recover', 'Recover', () => confirmInput('Performance'), ['recover'], false],

    ['dr', 'primary', 'reindex', 'Reindex', () => confirmInput('Disaster Recovery'), ['reindex'], false],
    ['performance', 'primary', 'reindex', 'Reindex', () => confirmInput('Performance'), ['reindex'], false],
    ['performance', 'secondary', 'reindex', 'Reindex', () => confirmInput('Performance'), ['reindex'], false],

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
        await blur('[name="primary_cluster_addr"]');
        await confirmInput('Performance');
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
        await blur('#secondary-token');
        await fillIn('#primary_api_addr', 'addr');
        await blur('#primary_api_addr');
        await confirmInput('Performance');
      },
      ['update-primary', 'secondary', { token: 'token', primary_api_addr: 'addr' }],
      false,
    ],
  ];

  for (const [
    replicationMode,
    clusterMode,
    action,
    headerText,
    fillInFn,
    expectedOnSubmit,
    oldVersion,
  ] of testCases) {
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
        <div id="modal-wormhole"></div>
        <ReplicationActions
          @model={{this.model}}
          @replicationMode={{this.replicationMode}}
          @selectedAction={{this.selectedAction}}
          @onSubmit={{action this.onSubmit}}
        />
        `
      );

      const selector = oldVersion ? 'h4' : `[data-test-${action}-replication] h4`;
      assert
        .dom(selector)
        .hasText(headerText, `${testKey}: renders the correct component header (${oldVersion})`);

      if (oldVersion) {
        await click('[data-test-confirm-action-trigger]');
        await click('[data-test-confirm-button]');
      } else {
        await click('[data-test-replication-action-trigger]');
        if (typeof fillInFn === 'function') {
          await fillInFn.call(this);
        }
        await blur('[data-test-confirmation-modal-input]');
        await click('[data-test-confirm-button]');
      }
    });
  }
});
