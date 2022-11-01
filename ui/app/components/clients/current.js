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
  @tracked namespaceArray = this.byNamespace.map((namespace) => {
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
    // array of upgrade data objects for noteworthy upgrades
    return relevantUpgrades;
  }

  // Response client count data by namespace for current/partial month
  get byNamespace() {
    return this.args.model.monthly?.byNamespace || [];
  }

  get isGatheringData() {
    // return true if tracking IS enabled but no data collected yet
    return this.args.model.config?.enabled === 'On' && this.byNamespace.length === 0;
  }

  get hasAttributionData() {
    if (this.selectedAuthMethod) return false;
    if (this.selectedNamespace) {
      return this.authMethodOptions.length > 0;
    }
    return this.totalUsageCounts.clients !== 0 && !!this.totalClientAttribution;
  }

  get filteredCurrentData() {
    const namespace = this.selectedNamespace;
    const auth = this.selectedAuthMethod;
    if (!namespace && !auth) {
      return this.byNamespace;
    }
    if (!auth) {
      return this.byNamespace.find((ns) => ns.label === namespace);
    }
    return this.byNamespace
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
      return `Vault was upgraded to ${versions} during this month.`;
    } else {
      let version = this.upgradeDuringCurrentMonth[0];
      return `Vault was upgraded to ${version.id} on this month.`;
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
        return ' We added mount level attribution starting in 1.10, so keep that in mind when looking at the data below.';
      }
    }
    // return combined explanation if spans multiple upgrades
    return ' How we count clients changed in 1.9 and we added mount level attribution starting in 1.10. Keep this in mind when looking at the data below.';
  }

  // top level TOTAL client counts for current/partial month
  get totalUsageCounts() {
    return this.selectedNamespace ? this.filteredCurrentData : this.args.model.monthly?.total;
  }

  // total client attribution data for horizontal bar chart in attribution component
  get totalClientAttribution() {
    if (this.selectedNamespace) {
      return this.filteredCurrentData?.mounts || null;
    } else {
      return this.byNamespace;
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
      const mounts = this.filteredCurrentData.mounts?.map((mount) => ({
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
