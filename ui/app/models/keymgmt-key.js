import Model, { attr } from '@ember-data/model';

export default class KeymgmtKeyModel extends Model {
  @attr('string') title;
  @attr('string') name;
}
