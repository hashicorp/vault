import { inject as service } from '@ember/service';
import Base from './cluster-route-base';

export default Base.extend({
  replicationMode: service(),
  beforeModel() {
    this._super(...arguments);
    this.get('replicationMode').setMode('dr');
  },
});
