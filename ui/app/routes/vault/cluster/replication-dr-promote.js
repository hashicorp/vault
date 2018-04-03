import Ember from 'ember';
import Base from './cluster-route-base';

export default Base.extend({
  replicationMode: Ember.inject.service(),
  beforeModel() {
    this._super(...arguments);
    this.get('replicationMode').setMode('dr');
  },
});
