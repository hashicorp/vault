import Component from '@glimmer/component';

export default class Current extends Component {
  chartLegend = [
    { key: 'entity_clients', label: 'entity clients' },
    { key: 'non_entity_clients', label: 'non-entity clients' },
  ];

  // data for horizontal bar chart in attribution component
  get topTenNamespaces() {
    if (!this.args.model.monthly || !this.args.model.monthly.byNamespace) {
      return null;
    }
    return this.args.model.monthly.byNamespace;
  }

  // top level TOTAL client counts from response for given month
  get runningTotals() {
    if (!this.args.model.monthly || !this.args.model.monthly.total) {
      return null;
    }
    return this.args.model.monthly.total;
  }

  get responseTimestamp() {
    if (!this.args.model.monthly || !this.args.model.monthly.responseTimestamp) {
      return null;
    }
    return this.args.model.monthly.responseTimestamp;
  }
}
