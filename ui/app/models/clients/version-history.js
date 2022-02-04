import Model, { attr } from '@ember-data/model';
export default class VersionHistoryModel extends Model {
  @attr('array') keys;
  @attr('object') key_info;
}
