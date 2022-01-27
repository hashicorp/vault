import Model, { attr } from '@ember-data/model';
export default class MonthlyModel extends Model {
  @attr('string') responseTimestamp;
  @attr('array') byNamespace;
  @attr('object') total;
}
