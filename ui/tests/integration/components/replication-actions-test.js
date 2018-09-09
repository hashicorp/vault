import { resolve } from 'rsvp';
import Service from '@ember/service';
import EmberObject, { computed } from '@ember/object';
import { module, test } from 'qunit';
import { setupRenderingTest } from 'ember-qunit';
import { click, fillIn, blur } from '@ember/test-helpers';
import hbs from 'htmlbars-inline-precompile';

const CapabilitiesStub = EmberObject.extend({
  canUpdate: computed('capabilities', function() {
    return (this.get('capabilities') || []).includes('root');
  }),
});

const storeStub = Service.extend({
  callArgs: null,
  capabilitiesReturnVal: null,
  findRecord(_, path) {
    const self = this;
    self.set('callArgs', { path });
    const caps = CapabilitiesStub.create({
      path,
      capabilities: self.get('capabilitiesReturnVal') || [],
    });
    return resolve(caps);
  },
});

module('Integration | Component | replication actions', function(hooks) {
  setupRenderingTest(hooks);

  hooks.beforeEach(function() {
    this.owner.register('service:store', storeStub);
    this.storeService = this.owner.lookup('service:store');
  });

  async function testAction(
    assert,
    replicationMode,
    clusterMode,
    action,
    headerText,
    capabilitiesPath,
    fillInFn,
    expectedOnSubmit
  ) {
    const testKey = `${replicationMode}-${clusterMode}-${action}`;
    if (replicationMode) {
      this.set('model', {
        replicationAttrs: {
          modeForUrl: clusterMode,
        },
        [replicationMode]: {
          mode: clusterMode,
          modeForUrl: clusterMode,
        },
      });
      this.set('replicationMode', replicationMode);
    } else {
      this.set('model', { mode: clusterMode });
    }
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
    this.render(
      hbs`{{replication-actions model=model replicationMode=replicationMode selectedAction=selectedAction onSubmit=(action onSubmit)}}`
    );

    assert.equal(
      this.$(`h4:contains(${headerText})`).length,
      1,
      `${testKey}: renders the correct partial as default`
    );

    if (typeof fillInFn === 'function') {
      fillInFn.call(this);
    }
    await click('button');
    await click('button.red');
  }

  function callTest(context, assert) {
    return function() {
      testAction.call(context, assert, ...arguments);
    };
  }

  test('actions', function(assert) {
    const t = callTest(this, assert);
    //TODO move to table test so we don't share the same store
    //t('dr', 'primary', 'disable', 'Disable dr replication', 'sys/replication/dr/primary/disable', null, ['disable', 'primary']);
    //t('performance', 'primary', 'disable', 'Disable performance replication', 'sys/replication/performance/primary/disable', null, ['disable', 'primary']);
    t('dr', 'secondary', 'disable', 'Disable replication', 'sys/replication/dr/secondary/disable', null, [
      'disable',
      'secondary',
    ]);
    t(
      'performance',
      'secondary',
      'disable',
      'Disable replication',
      'sys/replication/performance/secondary/disable',
      null,
      ['disable', 'secondary']
    );

    t('dr', 'primary', 'recover', 'Recover', 'sys/replication/recover', null, ['recover']);
    t('performance', 'primary', 'recover', 'Recover', 'sys/replication/recover', null, ['recover']);
    t('performance', 'secondary', 'recover', 'Recover', 'sys/replication/recover', null, ['recover']);

    t('dr', 'primary', 'reindex', 'Reindex', 'sys/replication/reindex', null, ['reindex']);
    t('performance', 'primary', 'reindex', 'Reindex', 'sys/replication/reindex', null, ['reindex']);
    t('dr', 'secondary', 'reindex', 'Reindex', 'sys/replication/reindex', null, ['reindex']);
    t('performance', 'secondary', 'reindex', 'Reindex', 'sys/replication/reindex', null, ['reindex']);

    t('dr', 'primary', 'demote', 'Demote cluster', 'sys/replication/dr/primary/demote', null, [
      'demote',
      'primary',
    ]);
    t(
      'performance',
      'primary',
      'demote',
      'Demote cluster',
      'sys/replication/performance/primary/demote',
      null,
      ['demote', 'primary']
    );
    // we don't do dr secondary promote in this component so just test perf
    t(
      'performance',
      'secondary',
      'promote',
      'Promote cluster',
      'sys/replication/performance/secondary/promote',
      async function() {
        await fillIn('[name="primary_cluster_addr"]', 'cluster addr');
        await blur('[name="primary_cluster_addr"]');
      },
      ['promote', 'secondary', { primary_cluster_addr: 'cluster addr' }]
    );

    // don't yet update-primary for dr
    t(
      'performance',
      'secondary',
      'update-primary',
      'Update primary',
      'sys/replication/performance/secondary/update-primary',
      async function() {
        await fillIn('#secondary-token', 'token');
        await blur('#secondary-token');
        await fillIn('#primary_api_addr', 'addr');
        await blur('#primary_api_addr');
      },
      ['update-primary', 'secondary', { token: 'token', primary_api_addr: 'addr' }]
    );
  });
});
