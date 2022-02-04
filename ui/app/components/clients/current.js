import Component from '@glimmer/component';
import { tracked } from '@glimmer/tracking';
import { compareAsc, startOfMonth } from 'date-fns';
import { action } from '@ember/object';
export default class Current extends Component {
  chartLegend = [
    { key: 'entity_clients', label: 'entity clients' },
    { key: 'non_entity_clients', label: 'non-entity clients' },
  ];
  @tracked namespaceArray = this.args.model.monthly?.byNamespace.map((namespace) => {
    return { name: namespace['label'], id: namespace['label'] };
  });
  @tracked selectedNamespace = null;

  get upgradeDate() {
    let keyInfoObject = this.args.model.versionHistory.key_info;
    if (!keyInfoObject) {
      return false;
    }
    let firstKey = Object.keys(keyInfoObject);
    // compare against this month and this year to show message or not.
    let versionDate = new Date(keyInfoObject[firstKey].timestamp_installed);
    let compare = compareAsc(versionDate, startOfMonth(new Date())); // if the versionDate happened after the first of today's month
    // Compare the two dates and return 1 if the first date is after the second, -1 if the first date is before the second or 0 if dates are equal.
    if (compare === 1) {
      return keyInfoObject[firstKey].timestamp_installed || false;
    }
    return false;
  }

  get licenseStartDate() {
    return this.args.licenseStartDate || null;
  }

  // API client count data by namespace for current/partial month
  get byNamespaceCurrent() {
    return this.args.model.monthly?.byNamespace || null;
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
