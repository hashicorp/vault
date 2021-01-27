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
    // only the combined model for roles (static and dynamic) will return a model name
    // for connections, we follow the normal process of permissions, which automatically has permissions added to model via lazy capabilities
    // therefore no need to peekRecord, just call @item.canEdit in the template
    if (!modelName) {
      return;
    }
    let roleModel = await this.store.peekRecord(modelName, id);
    this.permissions = roleModel;
  }
}
