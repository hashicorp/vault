import Model, { attr } from '@ember-data/model';
export default class Activity extends Model {
  @attr('string') responseTimestamp;
  @attr('array') byNamespace;
  @attr('string') endTime;
  @attr('string') startTime;
  @attr('object') total;
}
