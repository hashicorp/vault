import Component from '@glimmer/component';
import { service } from '@ember/service';
import { action } from '@ember/object';
import { tracked } from '@glimmer/tracking';
import { task } from 'ember-concurrency';
import { waitFor } from '@ember/test-waiters';

/**
 * @module LicenseProofOfValueComponent component TODO.
 * 
 * @example
 * <LicenseProofOfValue
    @model={{this.model}}
    @replication={{this.replication}}
    />
 *
 * @param {object} model - this is the license/status model which has features as a param.
 * @param {object} replication - this comes from the cluster adapter which allows you to fetch the replication status response.
 */

// success = enabled and on license
// warning = not enabled but on license
// nuetral = not enabled and not on license

const NAMESPACE_FEATURES = [
  { key: 'sentinel', label: 'Sentinel' },
  { key: 'sealWrapping', label: 'Seal Wrapping' },
  { key: 'controlGroups', label: 'Control Groups' },
  { key: 'kmip', label: 'KMIP Secret Engine' },
  { key: 'transform', label: 'Transform Secret Engine' },
  { key: 'keymgmt', label: 'Key Management Secret Engine' },
  { key: 'keyEncipherment', label: 'Key Encipherment' },
];

export default class LicenseProofOfValueComponent extends Component {
  namespaceFeatures = NAMESPACE_FEATURES;

  @service flags;
  @service store;
  @tracked namespaceFeaturesSelected = [];
  @tracked secretEngineData;
  @tracked openKmipModal = false;
  @tracked totalKmipMounts = 0;

  licenseHasFeature(featureName) {
    const { features } = this.args.model;
    return features ? features.includes(featureName) : false;
  }

  @action checkboxChange(name, value) {
    this.namespaceFeaturesSelected = value;
  }

  get networkRequestCounter() {
    return this.namespaceFeaturesSelected.length * this.namespaceCount;
  }

  get drStatus() {
    const { dr } = this.args.replication;
    const color = this.licenseHasFeature('DR Replication')
      ? dr.mode === 'enabled'
        ? 'success'
        : 'warning'
      : 'neutral';

    const text = this.licenseHasFeature('DR Replication')
      ? dr.mode === 'enabled'
        ? 'Enabled!'
        : 'Not enabled'
      : 'NA';

    return { color, text };
  }

  get performanceStatus() {
    const { performance } = this.args.replication;
    const color = this.licenseHasFeature('DR Replication')
      ? performance.mode === 'enabled'
        ? 'success'
        : 'warning'
      : 'neutral';

    const text = this.licenseHasFeature('DR Replication')
      ? performance.mode === 'enabled'
        ? 'Enabled!'
        : 'Not enabled'
      : 'NA';

    return { color, text };
  }

  get secretsSyncStatus() {
    const color = this.licenseHasFeature('Secrets Sync')
      ? this.flags.secretSyncIsActivated
        ? 'success'
        : 'warning'
      : 'neutral';
    const text = this.licenseHasFeature('Secrets Sync')
      ? this.flags.secretSyncIsActivated
        ? 'Activated!'
        : 'Not activated'
      : 'NA';

    return { color, text };
  }

  get kmipStatus() {
    // TODO return, maybe disable if not on license?
    const color = this.licenseHasFeature('KMIP') ? 'critical' : 'critical';
    const text = this.licenseHasFeature('KMIP') ? 'revist' : 'revisit';

    return { color, text };
  }

  get transformStatus() {
    const color = this.licenseHasFeature('Transform Secrets Engine') ? 'critical' : 'critical';
    const text = this.licenseHasFeature('Transform Secrets Engine') ? 'revist' : 'revisit';

    return { color, text };
  }

  get namespaceCount() {
    const { data } = this.args.namespaces;
    // if no data.keys does not exists return 0
    // TODO handle error case
    return !data.keys ? 0 : data.keys.length;
    // example data
    // data": {
    //     "keys": [
    //         "ns1/",
    //         "ns1/ns-child/"
    //     ]
    // },
  }

  countNestedItems(arr) {
    let count = 0;
    for (let i = 0; i < arr.length; i++) {
      if (typeof arr[i] === 'string' && arr[i].includes('/')) {
        count++;
      }
    }
    return count;
  }

  get namespaceDetails() {
    // return nested and level (ex: 66% of your namespaces are nested. you have a total of 4 levels of nesting)
    const { data } = this.args.namespaces;
    if (!data.keys) return null;
    // nested levels —probably turn into a helper?
    let nestedCount = 0;
    let maxSlashes = -1;
    // let itemWithMaxSlashes;
    for (const item of data.keys) {
      const slashCount = item.split('/').length - 1;
      if (slashCount > 1) {
        nestedCount++;
      }

      if (slashCount > maxSlashes) {
        maxSlashes = slashCount;
        // itemWithMaxSlashes = item; // todo not using but maybe should
      }
    }
    const percentNested = Number(nestedCount / data.keys.length).toFixed(2) * 100;
    return `${percentNested}% of your namespaces are nested.
    You have a total of ${maxSlashes} levels of nesting.`;
  }

  get namespaceLicenseFeaturesOnLicense() {
    // count the features returned on the model
    // TODO account for namespace only feature
    return this.args.model.features.length;
  }

  get allNamespaceOnlyFeatures() {
    return 17; // need to think on this
  }

  get kmipMountData() {
    // secretEngine data includes all secret-engine mounts for all namespaces
    // we want to pull out the kmip specific data making it easier to pass into the modal table
    // ex: secretEngineData = { "ns1": { "transform": 2, "kmip": 1, "keymgmt": 0}, " ": { "transform": 0, "kmip": 2, "keymgmt": 1}}
    // return { "ns1": 1, " ": 2}
    const kmipData = {};
    for (const ns in this.secretEngineData) {
      kmipData[ns] = this.secretEngineData[ns].kmip || 0;
    }
    // we also want to count the total number of kmip mounts and set the tracked property totalKmipMounts
    this.updateTotalKmipMounts(kmipData);
    // sort so the empty key is first (root namespace)
    return this.sortObjectsWithEmptyKeyFirst(kmipData);
  }

  updateTotalKmipMounts(kmipData) {
    this.totalKmipMounts = Object.values(kmipData).reduce((acc, val) => acc + val, 0);
  }

  sortObjectsWithEmptyKeyFirst(obj) {
    const keys = Object.keys(obj);

    keys.sort((a, b) => {
      if (a === '') return -1;
      if (b === '') return 1;
      return a.localeCompare(b);
    });

    const sortedObject = {};
    keys.forEach((key) => {
      sortedObject[key] = obj[key];
    });

    return sortedObject;
  }

  @task
  @waitFor
  *fetchNamespaceFeaturesData() {
    // we are fetching the internal mounts endpoint when the user selects any of the secret engine license features
    if (
      this.namespaceFeaturesSelected.some(
        (feature) => feature === 'kmip' || feature === 'transform' || feature === 'keymgmt'
      )
    ) {
      const response = yield this.fetchMountsByAllNamespaces();
      const featureCounts = {};
      for (const ns in response) {
        const secretEngineMounts = response[ns];

        if (Object.keys(secretEngineMounts).length < 1) {
          // if there are no secretEngineMounts in a namespace, add the namespace to the featureCounts object with a value of 0 for each feature
          featureCounts[ns] = { kmip: 0, transform: 0, keymgmt: 0 };
          continue;
        }

        Object.values(secretEngineMounts).forEach((mount) => {
          featureCounts[ns]
            ? featureCounts[ns][mount.type]
              ? featureCounts[ns][mount.type]++
              : (featureCounts[ns][mount.type] = 1)
            : (featureCounts[ns] = { [mount.type]: 1 });
        });
      }

      this.secretEngineData = featureCounts;
      // determine which modal to open
      if (this.namespaceFeaturesSelected.includes('kmip')) {
        this.openKmipModal = true;
      }
    }
  }

  // TODO naming here could use some love because we filter out only Secret Engine types
  async fetchMountsByAllNamespaces() {
    // ideally this would be on a route or service. but hackweek for now.
    const { data } = this.args.namespaces; // data.keys has the array of namespaces
    if (!data?.keys) return; // todo something better because some folks might want root only.
    const adapter = this.store.adapterFor('application');
    const mountResponseByNamespace = {};
    // add an empty key for the root namespace
    data.keys.push(' ');
    for (const ns of data.keys) {
      try {
        const response = await adapter.ajax('/v1/sys/internal/ui/mounts', 'GET', { namespace: ns });
        const mountSecretData = response.data.secret;
        // filterMountSecretData to only relevant secret-engines and their counts
        for (const key in mountSecretData) {
          // todo make these three a const
          if (!['transform', 'kmip', 'keymgmt'].includes(mountSecretData[key].type)) {
            delete mountSecretData[key];
          }
        }
        mountResponseByNamespace[ns] = mountSecretData;
      } catch (e) {
        // TODO handle error better
        throw new Error('Error fetching mounts');
      }
    }
    // mountResponseByNamespace returns an object with the namespace as the key and the secret engine mounts as the value. It shows all namespaces regardless of if they have secret engine mounts.
    return mountResponseByNamespace;
  }
}
