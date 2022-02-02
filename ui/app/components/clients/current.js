import Component from '@glimmer/component';
import { tracked } from '@glimmer/tracking';
import { action } from '@ember/object';
export default class Current extends Component {
  chartLegend = [
    { key: 'entity_clients', label: 'entity clients' },
    { key: 'non_entity_clients', label: 'non-entity clients' },
  ];
  searchSelectModel = ['namespace'];
  @tracked selectedNamespace = null;

  // TODO CMB get from model
  get upgradeDate() {
    return this.args.upgradeDate || null;
  }

  get licenseStartDate() {
    return this.args.licenseStartDate || null;
  }

  // by namespace client count data for partial month
  get byNamespaceCurrent() {
    return this.args.model.monthly?.byNamespace || null;
  }

  // data for horizontal bar chart in attribution component
  get topTenChartData() {
    if (this.selectedNamespace) {
      let filteredNamespace = this.filterByNamespace(this.selectedNamespace);
      return filteredNamespace.mounts
        ? this.filterByNamespace(this.selectedNamespace).mounts.slice(0, 10)
        : null;
    } else {
      return this.byNamespaceCurrent;
    }
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
  // ACTIONS
  @action
  selectNamespace(value) {
    // value comes in as [namespace0]
    this.selectedNamespace = value[0];
  }
}
