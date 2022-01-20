import Model, { attr } from '@ember-data/model';
// ARG TODO copied from before, modify for what you need
export default class NewInitActivityModel extends Model {
  @attr('object') total;
  @attr('string') endTime;
  @attr('string') startTime;
  @attr('array') byNamespace;
  @attr('array') byMonth;
  @attr('array') byMonthNewClients;
}
