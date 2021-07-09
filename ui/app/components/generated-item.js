import AdapterError from '@ember-data/adapter/error';
import { inject as service } from '@ember/service';
import Component from '@ember/component';
import { computed, set } from '@ember/object';
import { task } from 'ember-concurrency';

/**
 * @module GeneratedItem
 * The `GeneratedItem` component is the form to configure generated items related to mounts (e.g. groups, roles, users)
 *
 * @example
 * ```js
 * <GeneratedItem @model={{model}} @mode={{mode}} @itemType={{itemType/>
 * ```
 *
 * @property model=null {DS.Model} - The corresponding item model that is being configured.
 * @property mode {String} - which config mode to use. either `show`, `edit`, or `create`
 * @property itemType {String} - the type of item displayed
 *
 */

export default Component.extend({
  model: null,
  itemType: null,
  flashMessages: service(),
  router: service(),
  validationMessages: null,
  isFormInvalid: true,
  props: computed('model', function() {
    return this.model.serialize();
  }),
  saveModel: task(function*() {
    try {
      yield this.model.save();
    } catch (err) {
      // AdapterErrors are handled by the error-message component
      // in the form
      if (err instanceof AdapterError === false) {
        throw err;
      }
      return;
    }
    this.router.transitionTo('vault.cluster.access.method.item.list').followRedirects();
    this.flashMessages.success(`Successfully saved ${this.itemType} ${this.model.id}.`);
  }).withTestWaiter(),
  init() {
    this._super(...arguments);
    this.set('validationMessages', {});
    if (this.mode === 'edit') {
      // For validation to work in edit mode,
      // reconstruct the model values from field group
      this.model.fieldGroups.forEach(element => {
        if (element.default) {
          element.default.forEach(attr => {
            let fieldValue = attr.options && attr.options.fieldValue;
            if (fieldValue) {
              this.model[attr.name] = this.model[fieldValue];
            }
          });
        }
      });
    }
  },
  actions: {
    onKeyUp(name, value) {
      this.model.set(name, value);
      if (this.model.validations) {
        // Set validation error message for updated attribute
        this.model.validations.attrs[name] && this.model.validations.attrs[name].isValid
          ? set(this.validationMessages, name, '')
          : set(this.validationMessages, name, this.model.validations.attrs[name].message);

        // Set form button state
        this.model.validate().then(({ validations }) => {
          this.set('isFormInvalid', !validations.isValid);
        });
      }
    },
    deleteItem() {
      this.model.destroyRecord().then(() => {
        this.router.transitionTo('vault.cluster.access.method.item.list').followRedirects();
        this.flashMessages.success(`Successfully deleted ${this.itemType} ${this.model.id}.`);
      });
    },
  },
});
