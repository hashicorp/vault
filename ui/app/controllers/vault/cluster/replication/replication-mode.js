import Ember from 'ember';

export default Ember.Controller.extend({
  rm: Ember.inject.service('replication-mode'),
  replicationMode: Ember.computed.alias('rm.mode'),
});
