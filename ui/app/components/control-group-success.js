import { inject as service } from '@ember/service';
import Component from '@ember/component';
import { task } from 'ember-concurrency';

export default Component.extend({
  router: service(),
  controlGroup: service(),
  store: service(),

  // public attrs
  model: null,
  controlGroupResponse: null,

  //internal state
  error: null,
  unwrapData: null,

  unwrap: task(function* (token) {
    let adapter = this.store.adapterFor('tools');
    this.set('error', null);
    try {
      let response = yield adapter.toolAction('unwrap', null, { clientToken: token });
      this.set('unwrapData', response.auth || response.data);
      this.controlGroup.deleteControlGroupToken(this.model.id);
    } catch (e) {
      this.set('error', `Token unwrap failed: ${e.errors[0]}`);
    }
  }).drop(),

  markAndNavigate: task(function* () {
    this.controlGroup.markTokenForUnwrap(this.model.id);
    let { url } = this.controlGroupResponse.uiParams;
    yield this.router.transitionTo(url);
  }).drop(),
});
