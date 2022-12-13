import Controller from '@ember/controller';
import { task } from 'ember-concurrency';
import { inject as service } from '@ember/service';

export default Controller.extend({
  flashMessages: service(),

  queryParams: {
    page: 'page',
    pageFilter: 'pageFilter',
  },

  page: 1,
  pageFilter: null,
  filter: null,

  disableMethod: task(function* (method) {
    const { type, path } = method;
    try {
      yield method.destroyRecord();
      this.flashMessages.success(`The ${type} Auth Method at ${path} has been disabled.`);
    } catch (err) {
      this.flashMessages.danger(
        `There was an error disabling Auth Method at ${path}: ${err.errors.join(' ')}.`
      );
    }
  }).drop(),
});
