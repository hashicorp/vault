import { inject as service } from '@ember/service';
import { not } from '@ember/object/computed';
import Component from '@ember/component';
import { computed } from '@ember/object';

export default Component.extend({
  classNames: 'config-pki-ca',
  store: service('store'),
  flashMessages: service(),
  errors: null,

  /*
   * @param boolean
   * @private
   * bool that gets flipped if you have a CA cert and click the Replace Cert button
   */
  replaceCA: false,

  /*
   * @param boolean
   * @private
   * bool that gets flipped if you push the click the "Set signed intermediate" button
   */
  setSignedIntermediate: false,

  /*
   * @param boolean
   * @private
   * bool that gets flipped if you push the click the "Set signed intermediate" button
   */
  signIntermediate: false,

  /*
   * @param boolean
   * @private
   *
   * true when there's no CA cert currently configured
   */
  needsConfig: not('config.pem'),

  /*
   * @param DS.Model
   * @private
   *
   * a `pki-ca-certificate` model used to back the form when uploading or creating a CA cert
   * created and set on `init`, and unloaded on willDestroy
   *
   */
  model: null,

  /*
   * @param DS.Model
   * @public
   *
   * a `pki-config` model - passed in in the component usage
   *
   */
  config: null,

  /*
   * @param Function
   * @public
   *
   * function that gets called to refresh the config model
   *
   */
  onRefresh() {},

  loading: false,

  willDestroy() {
    const ca = this.model;
    if (ca) {
      ca.unloadRecord();
    }
    this._super(...arguments);
  },

  createOrReplaceModel(modelType) {
    const ca = this.model;
    const config = this.config;
    const store = this.store;
    const backend = config.get('backend');
    if (ca) {
      ca.unloadRecord();
    }
    const caCert = store.createRecord(modelType || 'pki-ca-certificate', {
      id: `${backend}-ca-cert`,
      backend,
    });
    this.set('model', caCert);
  },

  /*
   * @private
   * @returns array
   *
   * When a CA is configured, we let them download
   * the CA in der, pem, and the CA Chain in pem (if one exists)
   *
   * This array provides the text and download hrefs for those links.
   *
   */
  downloadHrefs: computed('config', 'config.{backend,pem,caChain,der}', function () {
    const config = this.config;
    const { backend, pem, caChain, der } = config;

    if (!pem) {
      return [];
    }

    const pemFile = new Blob([pem], { type: 'text/plain' });
    const links = [
      {
        display: 'Download CA Certificate in PEM format',
        name: `${backend}_ca.pem`,
        url: URL.createObjectURL(pemFile),
      },
      {
        display: 'Download CA Certificate in DER format',
        name: `${backend}_ca.der`,
        url: URL.createObjectURL(der),
      },
    ];
    if (caChain) {
      const caChainFile = new Blob([caChain], { type: 'text/plain' });
      links.push({
        display: 'Download CA Certificate Chain',
        name: `${backend}_ca_chain.pem`,
        url: URL.createObjectURL(caChainFile),
      });
    }
    return links;
  }),

  actions: {
    saveCA(method) {
      this.set('loading', true);
      const model = this.model;
      const isUpload = this.model.uploadPemBundle;
      model
        .save({ adapterOptions: { method } })
        .then(() => {
          if (method === 'setSignedIntermediate' || isUpload) {
            this.send('refresh');
            this.flashMessages.success('The certificate for this backend has been updated.');
          }
        })
        .catch((e) => {
          this.set('errors', e.errors);
        })
        .finally(() => {
          this.set('loading', false);
        });
    },
    refresh() {
      this.setProperties({
        setSignedIntermediate: false,
        signIntermediate: false,
        replaceCA: false,
      });
      this.onRefresh();
    },
    toggleReplaceCA() {
      if (!this.replaceCA) {
        this.createOrReplaceModel();
      }
      this.toggleProperty('replaceCA');
    },
    toggleVal(name, val) {
      if (!name) {
        return;
      }
      const model = name === 'signIntermediate' ? 'pki-ca-certificate-sign' : null;
      if (!this.get(name)) {
        this.createOrReplaceModel(model);
      }
      if (val !== undefined) {
        this.set(name, val);
      } else {
        this.toggleProperty(name);
      }
    },
  },
});
