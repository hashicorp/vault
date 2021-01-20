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

import { tracked } from '@glimmer/tracking';
import Component from '@glimmer/component';

export default class DatabaseListItem extends Component {
  @tracked keyType;
  get keyTypeValue() {
    const item = this.args.item;
    const internalModel = item._internalModel;
    if (internalModel.modelName === 'database/role') {
      return 'dynamic';
    } else if (internalModel.modelName === 'database/static-role') {
      return 'static';
    } else {
      return '';
    }
  }
}
