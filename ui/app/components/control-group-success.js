import Ember from 'ember';
import { task } from 'ember-concurrency';

const { computed, inject } = Ember;

export default Ember.Component.extend({
  router: inject.service(),
  controlGroup: inject.service(),
  store: inject.service(),

  // public attrs
  model: null,
  controlGroupResponse: null,

  //internal state
  error: null,
  unwrapData: null,

  unwrap: task(function*(token) {
    let adapter = this.get('store').adapterFor('tools');
    try {
      let response = yield adapter.toolAction('unwrap', null, { clientToken: token });
      this.set('unwrapData', response.auth || response.data);
      this.get('controlGroup').deleteControlGroupToken(this.get('model.id'));
    } catch (e) {
      this.set('error', `Token unwrap failed: ${e.errors[0]}`);
    }
  }).drop(),

  markAndNavigate: task(function*() {
    this.get('controlGroup').markTokenForUnwrap(this.get('model.id'));
    let { url, name, contexts, queryParams } = this.get('controlGroupResponse.uiParams');
    let router = this.get('router');

    return url ? router.transitionTo(url) : router.transitionTo(name, ...(contexts || []), { queryParams });
  }).drop(),
});
