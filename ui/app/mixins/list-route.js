import Ember from 'ember';

export default Ember.Mixin.create({
  queryParams: {
    page: {
      refreshModel: true,
    },
    pageFilter: {
      refreshModel: true,
    },
  },
});
