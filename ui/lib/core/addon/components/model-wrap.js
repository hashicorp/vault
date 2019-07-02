/**
 * @module ModelWrap
 * ModelWrap components provide a way to call methods on models directly from templates. This is done by yielding callMethod task to the wrapped component.
 *
 * @example
 * ```js
 * <ModelWrap as |m|>
     <button onclick={{action (perform m.callMethod "save" model "Saved!" "Errored!" (transition-to "route")}}>
 * </ModelWrap>
 * ```
 *
 * @yields callMethod {Function}
 *
 */
import { inject as service } from '@ember/service';
import Component from '@ember/component';
import { task } from 'ember-concurrency';
import layout from '../templates/components/model-wrap';

export default Component.extend({
  layout,
  flashMessages: service(),
  tagName: '',

  callMethod: task(function*(method, model, successMessage, failureMessage, successCallback = () => {}) {
    let flash = this.get('flashMessages');
    try {
      yield model[method]();
      flash.success(successMessage);
      successCallback();
    } catch (e) {
      let errString = e.errors.join(' ');
      flash.danger(failureMessage + ' ' + errString);
      model.rollbackAttributes();
    }
  }),
});
