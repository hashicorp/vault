import Component from '@glimmer/component';

export default class Current extends Component {
  chartLegend = [
    { key: 'entity_clients', label: 'entity clients' },
    { key: 'non_entity_clients', label: 'non-entity clients' },
  ];

  // data for horizontal bar chart in attribution component
  get topTenNamespaces() {
    return this.args.model.monthly?.byNamespace;
  }

  // top level TOTAL client counts from response for given month
  get runningTotals() {
    return this.args.model.monthly?.total;
  }

  get responseTimestamp() {
    return this.args.model.monthly?.responseTimestamp;
  }
}
