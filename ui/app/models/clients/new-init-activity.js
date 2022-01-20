import Model, { attr } from '@ember-data/model';
// ARG TODO copied from before, modify for what you need
export default class NewInitActivityModel extends Model {
  @attr('object') total;
  @attr('object') byNamespace;
  @attr('string') endTime;
  @attr('string') startTime;
  @attr('number') clients;
  @attr('number') distinct_entities;
  @attr('number') non_entity_tokens;
}
