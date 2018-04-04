import Ember from 'ember';
import { task } from 'ember-concurrency';

export default Ember.Controller.extend({
  showRoute: 'vault.cluster.access.identity.show',
  showTab: 'details',
  navToShow: task(function*(model) {
    yield this.transitionToRoute(this.get('showRoute'), model.id, this.get('showTab'));
  }),
});
