import Component from '@glimmer/component';
import { tracked } from '@glimmer/tracking';
import { isAfter, startOfMonth } from 'date-fns';
import { action } from '@ember/object';
export default class Current extends Component {
  chartLegend = [
    { key: 'entity_clients', label: 'entity clients' },
    { key: 'non_entity_clients', label: 'non-entity clients' },
  ];
  @tracked firstUpgradeVersion = this.args.model.versionHistory[0].id || null; // return 1.9.0 or earliest upgrade post 1.9.0
  @tracked upgradeDate = this.args.model.versionHistory[0].timestampInstalled || null; // returns RFC3339 timestamp

  @tracked selectedNamespace = null;
  @tracked namespaceArray = this.byNamespaceCurrent.map((namespace) => {
    return { name: namespace['label'], id: namespace['label'] };
  });

  @tracked selectedAuthMethod = null;
  @tracked authMethodOptions = [];

  // Response client count data by namespace for current/partial month
  get byNamespaceCurrent() {
    return this.args.model.monthly?.byNamespace || [];
  }

  get isGatheringData() {
    // return true if tracking IS enabled but no data collected yet
    return this.args.model.config?.enabled === 'On' && this.byNamespaceCurrent.length === 0;
  }

  get hasAttributionData() {
    if (this.selectedAuthMethod) return false;
    if (this.selectedNamespace) {
      return this.authMethodOptions.length > 0;
    }
    return this.totalUsageCounts.clients !== 0 && !!this.totalClientsData;
  }

  get filteredActivity() {
    const namespace = this.selectedNamespace;
    const auth = this.selectedAuthMethod;
    if (!namespace && !auth) {
      return this.getActivityResponse;
    }
    if (!auth) {
      return this.byNamespaceCurrent.find((ns) => ns.label === namespace);
    }
    return this.byNamespaceCurrent
      .find((ns) => ns.label === namespace)
      .mounts?.find((mount) => mount.label === auth);
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
    return this.selectedNamespace ? this.filteredActivity : this.args.model.monthly?.total;
  }

  // total client data for horizontal bar chart in attribution component
  get totalClientsData() {
    if (this.selectedNamespace) {
      return this.filteredActivity?.mounts || null;
    } else {
      return this.byNamespaceCurrent;
    }
  }

  get responseTimestamp() {
    return this.args.model.monthly?.responseTimestamp;
  }

  // ACTIONS
  @action
  selectNamespace([value]) {
    // value comes in as [namespace0]
    this.selectedNamespace = value;
    if (!value) {
      this.authMethodOptions = [];
      // on clear, also make sure auth method is cleared
      this.selectedAuthMethod = null;
    } else {
      // Side effect: set auth namespaces
      const mounts = this.filteredActivity.mounts?.map((mount) => ({
        id: mount.label,
        name: mount.label,
      }));
      this.authMethodOptions = mounts;
    }
  }

  @action
  setAuthMethod([authMount]) {
    this.selectedAuthMethod = authMount;
  }
}
