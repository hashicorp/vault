/**
 * @module DatabaseListItem
 * DatabaseListItem components are used for the list items for the Database Secret Engines for Roles.
 * This component automatically handles read-only list items if capabilities are not granted or the item is internal only.
 *
 * @example
 * ```js
 * <DatabaseListItem @item={item} />
 * ```
 * @param {object} item - item refers to the model item used on the list item partial
 */

import Component from '@glimmer/component';
import { tracked } from '@glimmer/tracking';
import { inject as service } from '@ember/service';

export default class DatabaseListItem extends Component {
  @tracked permissions = null;
  @tracked roleType = '';
  @service store;
  constructor() {
    super(...arguments);
    this.fetchPermissions();
  }

  get keyTypeValue() {
    const item = this.args.item;
    if (item.modelName === 'database/role') {
      return 'dynamic';
    } else if (item.modelName === 'database/static-role') {
      return 'static';
    } else {
      return '';
    }
  }

  async fetchPermissions() {
    let { id, modelName } = this.args.item;
    let roleModel = await this.store.peekRecord(modelName, id);
    this.permissions = roleModel;
  }
}
