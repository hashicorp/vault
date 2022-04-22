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
  @tracked namespaceArray = this.byNamespaceTotalClients.map((namespace) => {
    return { name: namespace['label'], id: namespace['label'] };
  });

  @tracked selectedAuthMethod = null;
  @tracked authMethodOptions = [];

  get upgradeVersionHistory() {
    const versionHistory = this.args.model.versionHistory;
    if (!versionHistory || versionHistory.length === 0) {
      return null;
    }

    // get upgrade data for initial upgrade to 1.9 and/or 1.10
    let relevantUpgrades = [];
    const importantUpgrades = ['1.9', '1.10'];
    importantUpgrades.forEach((version) => {
      let findUpgrade = versionHistory.find((versionData) => versionData.id.match(version));
      if (findUpgrade) relevantUpgrades.push(findUpgrade);
    });

    // if no history for 1.9 or 1.10, customer skipped these releases so get first stored upgrade
    if (relevantUpgrades.length === 0) {
      relevantUpgrades.push({
        id: versionHistory[0].id,
        previousVersion: versionHistory[0].previousVersion,
        timestampInstalled: versionHistory[0].timestampInstalled,
      });
    }
    // array of upgrade data objects for noteworthy upgrades
    return relevantUpgrades;
  }

  // Response total client count data by namespace for current/partial month
  get byNamespaceTotalClients() {
    return this.args.model.monthly?.byNamespaceTotalClients || [];
  }

  // Response new client count data by namespace for current/partial month
  get byNamespaceNewClients() {
    return this.args.model.monthly?.byNamespaceNewClients || [];
  }

  get isGatheringData() {
    // return true if tracking IS enabled but no data collected yet
    return (
      this.args.model.config?.enabled === 'On' &&
      this.byNamespaceTotalClients.length === 0 &&
      this.byNamespaceNewClients.length === 0
    );
  }

  get hasAttributionData() {
    if (this.selectedAuthMethod) return false;
    if (this.selectedNamespace) {
      return this.authMethodOptions.length > 0;
    }
    return this.totalUsageCounts.clients !== 0 && !!this.totalClientsData;
  }

  get filteredTotalData() {
    const namespace = this.selectedNamespace;
    const auth = this.selectedAuthMethod;
    if (!namespace && !auth) {
      return this.byNamespaceTotalClients;
    }
    if (!auth) {
      return this.byNamespaceTotalClients.find((ns) => ns.label === namespace);
    }
    return this.byNamespaceTotalClients
      .find((ns) => ns.label === namespace)
      .mounts?.find((mount) => mount.label === auth);
  }

  get filteredNewData() {
    const namespace = this.selectedNamespace;
    const auth = this.selectedAuthMethod;
    if (!namespace && !auth) {
      return this.byNamespaceNewClients;
    }
    if (!auth) {
      return this.byNamespaceNewClients.find((ns) => ns.label === namespace);
    }
    return this.byNamespaceNewClients
      .find((ns) => ns.label === namespace)
      .mounts?.find((mount) => mount.label === auth);
  }

  get upgradeDuringCurrentMonth() {
    if (!this.upgradeVersionHistory) {
      return null;
    }
    const upgradesWithinData = this.upgradeVersionHistory.filter((upgrade) => {
      // TODO how do timezones affect this?
      let upgradeDate = new Date(upgrade.timestampInstalled);
      return isAfter(upgradeDate, startOfMonth(new Date()));
    });
    // return all upgrades that happened within date range of queried activity
    return upgradesWithinData.length === 0 ? null : upgradesWithinData;
  }

  get upgradeVersionAndDate() {
    if (!this.upgradeDuringCurrentMonth) {
      return null;
    }
    if (this.upgradeDuringCurrentMonth.length === 2) {
      let versions = this.upgradeDuringCurrentMonth.map((upgrade) => upgrade.id).join(' and ');
      return `Vault was upgraded to ${versions} during this month`;
    } else {
      let version = this.upgradeDuringCurrentMonth[0];
      return `Vault was upgraded to ${version.id} on this month`;
    }
  }

  get versionSpecificText() {
    if (!this.upgradeDuringCurrentMonth) {
      return null;
    }
    if (this.upgradeDuringCurrentMonth.length === 1) {
      let version = this.upgradeDuringCurrentMonth[0].id;
      if (version.match('1.9')) {
        return ' How we count clients changed in 1.9, so keep that in mind when looking at the data below.';
      }
      if (version.match('1.10')) {
        return ' We added new client breakdowns starting in 1.10, so keep that in mind when looking at the data below.';
      }
    }
    // return combined explanation if spans multiple upgrades, or customer skipped 1.9 and 1.10
    return ' How we count clients changed in 1.9 and we added new client breakdowns starting in 1.10. Keep this in mind when looking at the data below.';
  }

  // top level TOTAL client counts for current/partial month
  get totalUsageCounts() {
    return this.selectedNamespace ? this.filteredTotalData : this.args.model.monthly?.total;
  }

  get newUsageCounts() {
    return this.selectedNamespace ? this.filteredNewData : this.args.model.monthly?.new;
  }

  // total client data for horizontal bar chart in attribution component
  get totalClientsData() {
    if (this.selectedNamespace) {
      return this.filteredTotalData?.mounts || null;
    } else {
      return this.byNamespaceTotalClients;
    }
  }

  // new client data for horizontal bar chart in attribution component
  get newClientsData() {
    if (this.selectedNamespace) {
      return this.filteredNewData?.mounts || null;
    } else {
      return this.byNamespaceNewClients;
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
      const mounts = this.filteredTotalData.mounts?.map((mount) => ({
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
