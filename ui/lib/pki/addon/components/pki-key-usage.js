import Component from '@glimmer/component';
import { action } from '@ember/object';

/**
 * @module PkiKeyUsage
 * PkiKeyUsage components are used to build out the toggle options for PKI's role create/update key_usage, ext_key_usage and ext_key_usage_oids model params.
 * Instead of having the user search on the following goLang pages for these options we present them in checkbox form and manually add them to the params as an array of strings.
 * key_usage options: https://pkg.go.dev/crypto/x509#KeyUsage
 * ext_key_usage options (not all are include on purpose): https://pkg.go.dev/crypto/x509#ExtKeyUsage
 * @example
 * ```js
 * <PkiKeyUsage @model={@model} @group={group}/>
 * ```
 * @param {class} model - The pki/pki-role-engine model.
 * @param {string} group - The name of the group created in the model. In this case, it's the "Key usage" group.
 */

const KEY_USAGE_FIELDS = {
  DigitalSignature: {
    label: 'Digital Signature',
    value: true,
  },
  ContentCommitment: { label: 'Content Commitment' },
  CrlSign: { label: 'CRL Sign' },
  KeyAgreement: {
    label: 'Key Agreement',
    value: true,
  },
  DataEncipherment: { label: 'Data Encipherment' },
  EncipherOnly: { label: 'Encipher Only' },
  KeyEncipherment: {
    label: 'Key Encipherment',
    value: true,
  },
  CertSign: { label: 'Cert Sign' },
  DecipherOnly: { label: 'Decipher Only' },
};

const EXT_KEY_USAGE_FIELDS = {
  Any: { label: 'Any' },
  EmailProtection: { label: 'Email Protection' },
  TimeStamping: { label: 'Time Stamping' },
  ServerAuth: { label: 'Server Auth' },
  IpsecEndSystem: { label: 'IPSEC End System' },
  OcspSigning: { label: 'OCSP Signing' },
  ClientAuth: { label: 'Client Auth' },
  IpsecTunnel: { label: 'IPSEC Tunnel' },
  IpsecUser: { label: 'IPSEC User' },
  CodeSigning: { label: 'Code Signing' },
};

export default class PkiKeyUsage extends Component {
  constructor() {
    super(...arguments);
    this.keyUsageFields = {};
    this.extKeyUsageFields = {};
    Object.assign(this.keyUsageFields, KEY_USAGE_FIELDS);
    Object.assign(this.extKeyUsageFields, EXT_KEY_USAGE_FIELDS);
    // set default of key_usage to the three params that are true by default.
    this.args.model.set('keyUsage', ['DigitalSignature', 'KeyAgreement', 'KeyEncipherment']);
  }

  @action onStringListChange(value) {
    this.args.model.set('extKeyUsageOids', value);
  }

  _amendList(checkboxName, value, type) {
    let keyUsageList = this.args.model.keyUsage;
    let extKeyUsageList = this.args.model.extKeyUsage;

    /* Process:
    1. We first check if the checkbox change is coming from the checkbox options of key_usage or ext_key_usage.
    // Param key_usage || ext_key_usage accept a comma separated string and an array of strings. E.g. "DigitalSignature,KeyAgreement,KeyEncipherment" || [“DigitalSignature”,“KeyAgreement”,“KeyEncipherment”]
    2. Then we convert the string to an array if it's not already an array (e.g. it's already been converted). This makes it easier to add or remove items.
    3. Then if the value of checkbox is "true" we add it to the arrayList, otherwise remove it.
    */
    if (type === 'keyUsage') {
      let keyUsageListArray = Array.isArray(keyUsageList) ? keyUsageList : keyUsageList.split(',');

      return value ? keyUsageListArray.addObject(checkboxName) : keyUsageListArray.removeObject(checkboxName);
    } else {
      // because there is no default on init for ext_key_usage property (set normally by OpenAPI) we define it as an empty array if it is undefined.
      let extKeyUsageListArray = !extKeyUsageList ? [] : extKeyUsageList;

      return value
        ? extKeyUsageListArray.addObject(checkboxName)
        : extKeyUsageListArray.removeObject(checkboxName);
    }
  }
  @action checkboxChange(type) {
    const checkboxName = event.target.id;
    const value = event.target['checked'];
    this.args.model.set(type, this._amendList(checkboxName, value, type));
  }
}
