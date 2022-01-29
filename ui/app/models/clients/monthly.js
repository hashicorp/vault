import Model, { attr } from '@ember-data/model';
// ARG TODO copied from before, modify for what you need
export default class Monthly extends Model {
  @attr('string') responseTimestamp;
  @attr('array') byNamespace;
  @attr('object') total;
}
