import Ember from 'ember';

const { computed, inject } = Ember;

export default Ember.Component.extend({
  model: null,
  controlGroupResponse: null,
  router: inject.service(),

  linkURL: computed('controlGroupResponse.uiParams', function() {
    let { name, contexts, queryParams } = this.get('controlGroupResponse.uiParams');
    let router = this.get('router');

    return router.urlFor(name, ...(contexts || []), { queryParams });
  }),
});
