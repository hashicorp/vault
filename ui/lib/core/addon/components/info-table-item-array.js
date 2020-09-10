import { computed } from '@ember/object';
import Component from '@ember/component';
import { inject as service } from '@ember/service';
import { task } from 'ember-concurrency';
import layout from '../templates/components/info-table-item-array';

/**
 * @module InfoTableItemArray
 * `InfoTableItemArray` //
 *
 * @example
 * ```js
 * <InfoTableItemArray @value={{['test-1','test-2','test-3']}} @isLink={{true}} @modelType="transform/role"/>
 * ```
 *
 * @param value=null {array} - This array of data to be displayed.  If there are more than 10 items in the array only five will show and a count of the other number in the array will show.
 * @param [isLink=true] {Boolean} - Indicates if the item should contain a link-to component.  Only setup for arrays, but this could be changed if needed.
 * @param [modelType=null] {string} - Tells what model you want data for the allOptions to be returned from.  Used in conjunction with the the isLink.
 */
export default Component.extend({
  layout,
  allOptions: null,
  value: null,
  store: service(),
  fetchOptions: task(function*() {
    if (this.isLink && this.modelType) {
      try {
        let queryOptions = {};
        let backendModel = yield this.store.peekAll('secret-engine');
        let array = backendModel.toArray().map(option => {
          return option.id;
        });
        if (array) {
          queryOptions = { backend: array.get('firstObject') };
        }
        let options = yield this.store.query(this.modelType, queryOptions);
        this.formatOptions(options);
      } catch (err) {
        throw err;
      }
    }
  }).on('didInsertElement'),

  formatOptions: function(options) {
    let allOptions = options.toArray().map(option => {
      return option.id;
    });
    this.set('allOptions', allOptions);
  },
});
