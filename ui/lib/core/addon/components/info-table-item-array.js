import { computed } from '@ember/object';
import Component from '@ember/component';
import { inject as service } from '@ember/service';
import { task } from 'ember-concurrency';
import layout from '../templates/components/info-table-item-array';
import { isWildcardString } from 'vault/helpers/is-wildcard-string';

/**
 * @module InfoTableItemArray
 * The `InfoTableItemArray` component handles arrays in the info-table-row component.
 * If an array has more than 10 items, then only 5 are displayed and a count of total items is displayed next to the five.
 * If a isLink is true than a link can be set for the use to click on the specific array item
 * If a wildcard is a potential variable in the string of an item, then you can use the modelType and wildcardLabel parameters to
 * return a wildcard count similar to what is done in the searchSelect component.
 *
 * @example
 * ```js
 * <InfoTableItemArray @displayArray={{['test-1','test-2','test-3']}} @isLink={{true}} @modelType="transform/role"/ @queryParam="role" @backend="transform" viewAll="roles">
 * ```
 *
 * @param displayArray=null {array} - This array of data to be displayed.  If there are more than 10 items in the array only five will show and a count of the other number in the array will show.
 * @param [isLink=true] {Boolean} - Indicates if the item should contain a link-to component.  Only setup for arrays, but this could be changed if needed.
 * @param [modelType=null] {string} - Tells what model you want data for the allOptions to be returned from.  Used in conjunction with the the isLink.
 * @param [wildcardLabel] {String} - when you want the component to return a count on the model for options returned when using a wildcard you must provide a label of the count e.g. role.  Should be singular.
 * @param [queryParam] {String} - If you want to specific a tab for the View All XX to display to.  Ex: role
 * @param [backend] {String} - To specify which backend to point the link to.
 * @param [viewAll] {String} - Specify the word at the end of the link View all xx.
 */
export default Component.extend({
  layout,
  'data-test-info-table-item-array': true,
  allOptions: null,
  displayArray: null,
  wildcardInDisplayArray: false,
  store: service(),
  displayArrayAmended: computed('displayArray', function() {
    let { displayArray } = this;
    if (displayArray.length >= 10) {
      // if array greater than 10 in length only display the first 5
      displayArray = displayArray.slice(0, 5);
    }

    return displayArray;
  }),

  checkWildcardInArray: task(function*() {
    let filteredArray = yield this.displayArray.filter(item => isWildcardString(item));
    this.set('wildcardInDisplayArray', filteredArray.length > 0 ? true : false);
  }).on('didInsertElement'),

  fetchOptions: task(function*() {
    if (this.isLink && this.modelType) {
      let queryOptions = {};

      if (this.backend) {
        queryOptions = { backend: this.backend };
      }

      let options = yield this.store.query(this.modelType, queryOptions);
      this.formatOptions(options);
    }
  }).on('didInsertElement'),

  formatOptions: function(options) {
    let allOptions = options.toArray().map(option => {
      return option.id;
    });
    this.set('allOptions', allOptions);
  },
});
