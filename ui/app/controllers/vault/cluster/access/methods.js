import Ember from 'ember';
import { task } from 'ember-concurrency';

export default Ember.Controller.extend({
  queryParams: {
    page: 'page',
    pageFilter: 'pageFilter',
  },

  page: 1,
  pageFilter: null,
  filter: null,

  disableMethod: task(function*(method) {
    const { type, path } = method.getProperties('type', 'path');
    try {
      yield method.destroyRecord();
      this.get('flashMessages').success(`The ${type} auth method at ${path} has been disabled.`);
    } catch (err) {
      this.get('flashMessages').danger(
        `There was an error disabling auth method at ${path}: ${err.errors.join(' ')}.`
      );
    }
  }).drop(),
});
