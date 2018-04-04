import { moduleForComponent, test } from 'ember-qunit';
import hbs from 'htmlbars-inline-precompile';
import Ember from 'ember';

const CapabilitiesStub = Ember.Object.extend({
  canUpdate: Ember.computed('capabilities', function() {
    return (this.get('capabilities') || []).includes('root');
  }),
});

const storeStub = Ember.Service.extend({
  callArgs: null,
  capabilitiesReturnVal: null,
  findRecord(_, path) {
    const self = this;
    self.set('callArgs', { path });
    const caps = CapabilitiesStub.create({
      path,
      capabilities: self.get('capabilitiesReturnVal') || [],
    });
    return Ember.RSVP.resolve(caps);
  },
});

moduleForComponent('replication-actions', 'Integration | Component | replication actions', {
  integration: true,
  beforeEach: function() {
    this.register('service:store', storeStub);
    this.inject.service('store', { as: 'storeService' });
  },
});

function testAction(
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
    return Ember.RSVP.resolve();
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
  this.$('button').click();
  this.$('button.red').click();
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
    function() {
      this.$('[name="primary_cluster_addr"]').val('cluster addr').change();
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
    function() {
      this.$('#secondary-token').val('token').change();
      this.$('#primary_api_addr').val('addr').change();
    },
    ['update-primary', 'secondary', { token: 'token', primary_api_addr: 'addr' }]
  );
});
