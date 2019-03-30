import { run } from '@ember/runloop';
import { resolve } from 'rsvp';
import Service from '@ember/service';
import { module, test } from 'qunit';
import { setupRenderingTest } from 'ember-qunit';
import { click, fillIn, blur, render, find } from '@ember/test-helpers';
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

module('Integration | Component | replication actions', function(hooks) {
  setupRenderingTest(hooks);

  hooks.beforeEach(function() {
    run(() => {
      this.owner.register('service:router', routerService);
      this.owner.unregister('service:store');
      this.owner.register('service:store', storeStub);
      this.storeService = this.owner.lookup('service:store');
    });
  });

  let testCases = [
    ['dr', 'primary', 'disable', 'Disable Replication', null, ['disable', 'primary']],
    ['performance', 'primary', 'disable', 'Disable Replication', null, ['disable', 'primary']],
    ['dr', 'secondary', 'disable', 'Disable Replication', null, ['disable', 'secondary']],
    ['performance', 'secondary', 'disable', 'Disable Replication', null, ['disable', 'secondary']],
    ['dr', 'primary', 'recover', 'Recover', null, ['recover']],
    ['performance', 'primary', 'recover', 'Recover', null, ['recover']],
    ['performance', 'secondary', 'recover', 'Recover', null, ['recover']],

    ['dr', 'primary', 'reindex', 'Reindex', null, ['reindex']],
    ['performance', 'primary', 'reindex', 'Reindex', null, ['reindex']],
    ['dr', 'secondary', 'reindex', 'Reindex', null, ['reindex']],
    ['performance', 'secondary', 'reindex', 'Reindex', null, ['reindex']],

    ['dr', 'primary', 'demote', 'Demote cluster', null, ['demote', 'primary']],
    ['performance', 'primary', 'demote', 'Demote cluster', null, ['demote', 'primary']],
    // we don't do dr secondary promote in this component so just test perf
    [
      'performance',
      'secondary',
      'promote',
      'Promote cluster',
      async function() {
        await fillIn('[name="primary_cluster_addr"]', 'cluster addr');
        await blur('[name="primary_cluster_addr"]');
      },
      ['promote', 'secondary', { primary_cluster_addr: 'cluster addr' }],
    ],

    // don't yet update-primary for dr
    [
      'performance',
      'secondary',
      'update-primary',
      'Update primary',
      async function() {
        await fillIn('#secondary-token', 'token');
        await blur('#secondary-token');
        await fillIn('#primary_api_addr', 'addr');
        await blur('#primary_api_addr');
      },
      ['update-primary', 'secondary', { token: 'token', primary_api_addr: 'addr' }],
    ],
  ];

  for (let [replicationMode, clusterMode, action, headerText, fillInFn, expectedOnSubmit] of testCases) {
    test(`replication mode ${replicationMode}, cluster mode: ${clusterMode}, action: ${action}`, async function(assert) {
      const testKey = `${replicationMode}-${clusterMode}-${action}`;
      this.set('model', {
        replicationAttrs: {
          modeForUrl: clusterMode,
        },
        [replicationMode]: {
          mode: clusterMode,
          modeForUrl: clusterMode,
        },
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
        hbs`{{replication-actions model=model replicationMode=replicationMode selectedAction=selectedAction onSubmit=(action onSubmit)}}`
      );

      assert.equal(
        find('h4').textContent.trim(),
        headerText,
        `${testKey}: renders the correct partial as default`
      );

      if (typeof fillInFn === 'function') {
        await fillInFn.call(this);
      }
      await click('[data-test-confirm-action-trigger]');
      await click('[data-test-confirm-button]');
    });
  }
});
