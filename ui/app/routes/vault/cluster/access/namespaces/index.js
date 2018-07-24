import Ember from 'ember';

export default Ember.Route.extend({
  model() {
    return this.store.findAll('namespace').catch(e => {
      if (e.httpStatus === 404) {
        return [];
      }
      throw e;
    });
  },
});
