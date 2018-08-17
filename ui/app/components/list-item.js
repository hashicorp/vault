import Ember from 'ember';
import { task } from 'ember-concurrency';

const { inject } = Ember;
export default Ember.Component.extend({
  flashMessages: inject.service(),
  tagName: '',
  linkParams: null,
  componentName: null,
  hasMenu: false,

  callMethod: task(function*(method, model, successMessage, failureMessage, successCallback = () => {}) {
    let flash = this.get('flashMessages');
    try {
      yield model[method]();
      flash.success(successMessage);
      successCallback();
    } catch (e) {
      let errString = e.errors.join(' ');
      flash.danger(failureMessage + errString);
      model.rollbackAttributes();
    }
  }),
});
