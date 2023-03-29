/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: MPL-2.0
 */

import Component from '@glimmer/component';
import { tracked } from '@glimmer/tracking';
import { action } from '@ember/object';
import { keyParamsByType } from 'pki/utils/action-params';
import PkiActionModel from 'vault/models/pki/action';

interface Args {
  model: PkiActionModel;
  groups: Map<[key: string], Array<string>> | null;
}

export default class PkiGenerateToggleGroupsComponent extends Component<Args> {
  @tracked showGroup: string | null = null;

  get keyParamFields() {
    const { type } = this.args.model;
    if (!type) return null;
    const fields = keyParamsByType(type);
    return fields.map((fieldName) => {
      return this.args.model.allFields.find((attr) => attr.name === fieldName);
    });
  }

  get groups() {
    if (this.args.groups) return this.args.groups;
    const groups = {
      'Key parameters': this.keyParamFields,
      'Subject Alternative Name (SAN) Options': ['altNames', 'ipSans', 'uriSans', 'otherSans'],
      'Additional subject fields': [
        'ou',
        'organization',
        'country',
        'locality',
        'province',
        'streetAddress',
        'postalCode',
      ],
    };
    // excludeCnFromSans and serialNumber are present in default fields for generate-csr -- only include for other types
    if (this.args.model.actionType !== 'generate-csr') {
      groups['Subject Alternative Name (SAN) Options'].unshift('excludeCnFromSans', 'serialNumber');
    }
    return groups;
  }

  @action
  toggleGroup(group: string, isOpen: boolean) {
    this.showGroup = isOpen ? group : null;
  }
}
