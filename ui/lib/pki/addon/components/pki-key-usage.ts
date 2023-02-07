/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: MPL-2.0
 */

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

interface Field {
  key: string;
  label: string;
}

const KEY_USAGE_FIELDS: Field[] = [
  { key: 'DigitalSignature', label: 'Digital Signature' },
  { key: 'ContentCommitment', label: 'Content Commitment' },
  { key: 'CrlSign', label: 'CRL Sign' },
  { key: 'KeyAgreement', label: 'Key Agreement' },
  { key: 'DataEncipherment', label: 'Data Encipherment' },
  { key: 'EncipherOnly', label: 'Encipher Only' },
  { key: 'KeyEncipherment', label: 'Key Encipherment' },
  { key: 'CertSign', label: 'Cert Sign' },
  { key: 'DecipherOnly', label: 'Decipher Only' },
];

const EXT_KEY_USAGE_FIELDS: Field[] = [
  { key: 'Any', label: 'Any' },
  { key: 'EmailProtection', label: 'Email Protection' },
  { key: 'TimeStamping', label: 'Time Stamping' },
  { key: 'ServerAuth', label: 'Server Auth' },
  { key: 'IpsecEndSystem', label: 'IPSEC End System' },
  { key: 'OcspSigning', label: 'OCSP Signing' },
  { key: 'ClientAuth', label: 'Client Auth' },
  { key: 'IpsecTunnel', label: 'IPSEC Tunnel' },
  { key: 'IpsecUser', label: 'IPSEC User' },
  { key: 'CodeSigning', label: 'Code Signing' },
];

interface PkiKeyUsageArgs {
  group: string;
  model: {
    keyUsage: string[];
    extKeyUsageOids: string[];
    extKeyUsage: string[];
  };
}

export default class PkiKeyUsage extends Component<PkiKeyUsageArgs> {
  keyUsageFields = KEY_USAGE_FIELDS;
  extKeyUsageFields = EXT_KEY_USAGE_FIELDS;

  @action onStringListChange(value: string[]) {
    this.args.model.extKeyUsageOids = value;
  }

  _amendList(checkboxName: string, value: boolean, type: string): string[] {
    const list = type === 'keyUsage' ? this.args.model.keyUsage : this.args.model.extKeyUsage;
    const idx = list.indexOf(checkboxName);
    if (value === true && idx < 0) {
      list.push(checkboxName);
    } else if (value === false && idx >= 0) {
      list.splice(idx, 1);
    }
    return list;
  }

  @action checkboxChange(name: string, value: string[]) {
    // Make sure we can set this value type to this model key
    if (name === 'keyUsage' || name === 'extKeyUsage') {
      this.args.model[name] = value;
    }
  }
}
