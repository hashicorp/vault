import Model, { attr } from '@ember-data/model';
export default class MonthlyModel extends Model {
  @attr('string') responseTimestamp;
  @attr('array') byNamespace;
  @attr('object') total; // total clients during the current/partial month
  @attr('object') new; // total NEW clients during the current/partial
  @attr('array') byNamespaceTotalClients;
  @attr('array') byNamespaceNewClients;
}
