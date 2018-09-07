import { hash } from 'rsvp';
import Route from '@ember/routing/route';

export default Route.extend({
  model() {
    return hash({
      cluster: this.modelFor('vault.cluster'),
      seal: this.store.findRecord('capabilities', 'sys/seal'),
    });
  },
});
