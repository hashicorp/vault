import Model, { attr } from '@ember-data/model';
// ARG TODO copied from before, modify for what you need
export default class Monthly extends Model {
  @attr('string') responseTimestamp;
  @attr('array') byNamespace;
  @attr('object') total;
  // TODO CMB remove 'clients' and use 'total' object?
  @attr('number') clients;
  // new names
  @attr('number') entityClients;
  @attr('number') nonEntityClients;
  // old names
  @attr('number') distinctEntities;
  @attr('number') nonEntityTokens;
}
