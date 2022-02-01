import Model, { attr } from '@ember-data/model';
export default class Activity extends Model {
  @attr('string') responseTimestamp;
  @attr('array') byNamespace;
  @attr('string') endTime;
  @attr('array') formattedEndTime;
  @attr('array') formattedStartTime;
  @attr('string') startTime;
  @attr('object') total;
}
