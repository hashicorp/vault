import Component from '@glimmer/component';
import { tracked } from '@glimmer/tracking';
import { inject as service } from '@ember/service';
import { task } from 'ember-concurrency';
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
 * <InfoTableItemArray
 * @label="Roles"
 * @displayArray={{['test-1','test-2','test-3']}}
 * @isLink={{true}}
 * @rootRoute="vault.cluster.secrets.backend.list-root"
 * @itemRoute="vault.cluster.secrets.backend.show"
 * @modelType="transform/role"
 * @queryParam="role"
 * @backend="transform"
 * ```
 *
 * @param {string} label - used to render lowercased display text for "View all <label>."
 * @param {array} displayArray - The array of data to be displayed. (In InfoTableRow this comes from the @value arg.) If the array length > 10, and @doNotTruncate is false only 5 will show with a count of the number hidden.
 * @param {boolean} [isLink]  - Indicates if the item should contain a link-to component.  Only setup for arrays, but this could be changed if needed.
 * @param {string || array} [rootRoute="vault.cluster.secrets.backend.list-root"] -  Tells what route the link should go to when selecting "view all". If the route requires more than one dynamic param, insert an array.
 * @param {string || array} [itemRoute=vault.cluster.secrets.backend.show] - Tells what route the link should go to when selecting the individual item. If the route requires more than one dynamic param, insert an array.
 * @param {string} [modelType]  - Tells which model you want data for the allOptions to be returned from.  Used in conjunction with the the isLink.
 * @param {string} [wildcardLabel]  - when you want the component to return a count on the model for options returned when using a wildcard you must provide a label of the count e.g. role.  Should be singular.
 * @param {string} [backend] - To specify which backend to point the link to.
 * @param {boolean} [doNotTruncate=false] - Determines whether to show the View all "roles" link.
 */
export default class InfoTableItemArray extends Component {
  @tracked allOptions = null;
  @tracked wildcardInDisplayArray = false;
  @service store;

  get rootRoute() {
    return this.args.rootRoute || 'vault.cluster.secrets.backend.list-root';
  }

  get itemRoute() {
    return this.args.itemRoute || 'vault.cluster.secrets.backend.show';
  }

  get doNotTruncate() {
    return this.args.doNotTruncate || false;
  }

  get displayArrayTruncated() {
    let { displayArray } = this.args;
    if (!displayArray) return null;
    if ((displayArray.length >= 10) & !this.args.doNotTruncate) {
      // if array greater than 10 in length only display the first 5
      return displayArray.slice(0, 5);
    }
    return displayArray;
  }

  async checkWildcardInArray() {
    if (!this.args.displayArray) {
      return;
    }
    let filteredArray = await this.args.displayArray.filter((item) => isWildcardString([item]));
    this.wildcardInDisplayArray = filteredArray.length > 0 ? true : false;
  }

  @task *fetchOptions() {
    if (this.args.isLink && this.args.modelType) {
      let queryOptions = {};

      if (this.args.backend) {
        queryOptions = { backend: this.args.backend };
      }

      let options = yield this.store.query(this.args.modelType, queryOptions);
      this.formatOptions(options);
    }
    this.checkWildcardInArray();
  }

  formatOptions(options) {
    this.allOptions = options.mapBy('id');
  }
}
