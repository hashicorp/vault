import Component from '@glimmer/component';
import { tracked } from '@glimmer/tracking';
import { isAfter, startOfMonth } from 'date-fns';
import { action } from '@ember/object';
export default class Current extends Component {
  chartLegend = [
    { key: 'entity_clients', label: 'entity clients' },
    { key: 'non_entity_clients', label: 'non-entity clients' },
  ];
  @tracked selectedNamespace = null;
  @tracked namespaceArray = this.byNamespaceCurrent.map((namespace) => {
    return { name: namespace['label'], id: namespace['label'] };
  });
  @tracked firstUpgradeVersion = this.args.model.versionHistory[0].id || null; // return 1.9.0 or earliest upgrade post 1.9.0
  @tracked upgradeDate = this.args.model.versionHistory[0].timestampInstalled || null; // returns RFC3339 timestamp

  // API client count data by namespace for current/partial month
  get byNamespaceCurrent() {
    return this.args.model.monthly?.byNamespace || [];
  }

  get countsIncludeOlderData() {
    let firstUpgrade = this.args.model.versionHistory[0];
    if (!firstUpgrade) {
      return false;
    }
    let versionDate = new Date(firstUpgrade.timestampInstalled);
    // compare against this month and this year to show message or not.
    return isAfter(versionDate, startOfMonth(new Date())) ? versionDate : false;
  }

  // top level TOTAL client counts for current/partial month
  get totalUsageCounts() {
    return this.selectedNamespace
      ? this.filterByNamespace(this.selectedNamespace)
      : this.args.model.monthly?.total;
  }

  // total client data for horizontal bar chart in attribution component
  get totalClientsData() {
    if (this.selectedNamespace) {
      let filteredNamespace = this.filterByNamespace(this.selectedNamespace);
      return filteredNamespace.mounts ? this.filterByNamespace(this.selectedNamespace).mounts : null;
    } else {
      return this.byNamespaceCurrent;
    }
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
  selectNamespace([value]) {
    // value comes in as [namespace0]
    this.selectedNamespace = value;
  }
}
