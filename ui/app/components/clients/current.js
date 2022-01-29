import Component from '@glimmer/component';
import { tracked } from '@glimmer/tracking';
export default class Current extends Component {
  // TODO CMB delete - just for template view testing
  upgradeDate = new Date('2022-02-29T01:14:38.836Z');
  billingStartDate = new Date('2022-01-29T01:14:38.836Z');
  chartLegend = [
    { key: 'entity_clients', label: 'entity clients' },
    { key: 'non_entity_clients', label: 'non-entity clients' },
  ];
  @tracked selectedNamespace = null;

  // by namespace client count data for partial month
  get byNamespaceCurrent() {
    return this.args.model.monthly?.byNamespace;
  }

  // data for horizontal bar chart in attribution component
  get topTenChartData() {
    return this.selectedNamespace
      ? this.filterByNamespace(this.selectedNamespace).mounts.slice(0, 10)
      : this.byNamespaceCurrent.slice(0, 10);
  }

  // top level TOTAL client counts from response for given month
  get totalUsageCounts() {
    return this.selectedNamespace
      ? this.filterByNamespace(this.selectedNamespace)
      : this.args.model.monthly?.total;
  }

  get responseTimestamp() {
    return this.args.model.monthly?.responseTimestamp;
  }

  // HELPERS
  filterByNamespace(namespace) {
    return this.byNamespaceCurrent.find((ns) => ns.label === namespace);
  }
}
