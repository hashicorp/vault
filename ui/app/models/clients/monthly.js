import Model, { attr } from '@ember-data/model';
// ARG TODO copied from before, modify for what you need
export default class MonthlyModel extends Model {
  @attr('object') byNamespace;
  @attr('number') clients;
  @attr('number') distinct_entities;
  @attr('number') non_entity_tokens;
}
