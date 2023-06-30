import Model, { attr } from '@ember-data/model';

export default class SecretListModel extends Model {
  @attr() secrets;
  @attr() keyInfo;
}
