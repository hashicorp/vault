import Model, { attr } from '@ember-data/model';

export default class Activity extends Model {
  @attr('object') total;
  @attr('string') endTime;
  @attr('string') startTime;
  @attr('array') byNamespace;
}
