import { inject as service } from '@ember/service';
import Component from '@ember/component';
import { task } from 'ember-concurrency';

export default Component.extend({
  flashMessages: service(),
  tagName: '',
  linkParams: null,
  componentName: null,
  hasMenu: true,

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
