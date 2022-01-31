import Model, { attr } from '@ember-data/model';
export default class Activity extends Model {
  @attr('string') responseTimestamp;
  @attr('array') byNamespace;
  @attr('string') endTime;
  @attr('string') formattedEndTime;
  @attr('string') formattedStartTime;
  @attr('string') startTime;
  @attr('object') total;
}
