import Model, { attr } from '@ember-data/model';
export default class MonthlyModel extends Model {
  @attr('string') responseTimestamp;
  @attr('object') total; // total clients during the current/partial month
  @attr('array') byNamespace;
}
