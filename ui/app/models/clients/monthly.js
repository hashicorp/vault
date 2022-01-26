import Model, { attr } from '@ember-data/model';
export default class MonthlyModel extends Model {
  @attr('array') byNamespace;
  @attr('object') total;
  @attr('string') timestamp;
  // TODO CMB remove 'clients' and use 'total' object?
  @attr('number') clients;
  // new names
  @attr('number') entityClients;
  @attr('number') nonEntityClients;
  // old names
  @attr('number') distinctEntities;
  @attr('number') nonEntityTokens;
}
