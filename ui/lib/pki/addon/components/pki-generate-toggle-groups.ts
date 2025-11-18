/**
 * Copyright IBM Corp. 2016, 2025
 * SPDX-License-Identifier: BUSL-1.1
 */

import Component from '@glimmer/component';
import { tracked } from '@glimmer/tracking';
import { action } from '@ember/object';
import { keyParamsByType } from 'pki/utils/action-params';

import PkiConfigGenerateForm from 'vault/forms/secrets/pki/config/generate';
import type { ModelValidations } from 'vault/vault/app-types';

interface Args {
  form: PkiConfigGenerateForm;
  actionType: string;
  groups: Map<[key: string], Array<string>> | null;
  modelValidations?: ModelValidations;
}

export default class PkiGenerateToggleGroupsComponent extends Component<Args> {
  @tracked showGroup: string | null = null;

  // shim until sign-intermediate model is migrated to form
  get fieldsKey() {
    return this.args.form instanceof PkiConfigGenerateForm ? 'formFields' : 'allFields';
  }

  get keyParamFields() {
    const { form } = this.args;
    if (form.data.type) {
      const fields = keyParamsByType(form.data.type);
      return fields.map((fieldName) => {
        return form.formFields.find((field) => field.name === fieldName);
      });
    }
    return null;
  }

  get groups() {
    if (this.args.groups) return this.args.groups;
    const groups = {
      'Key parameters': this.keyParamFields,
      'Subject Alternative Name (SAN) Options': ['alt_names', 'ip_sans', 'uri_sans', 'other_sans'],
      'Additional subject fields': [
        'ou',
        'organization',
        'country',
        'locality',
        'province',
        'street_address',
        'postal_code',
      ],
    };
    // excludeCnFromSans and serialNumber are present in default fields for generate-csr -- only include for other types
    if (this.args.actionType !== 'generate-csr') {
      groups['Subject Alternative Name (SAN) Options'].unshift(
        'exclude_cn_from_sans',
        'subject_serial_number'
      );
    }
    return groups;
  }

  @action
  toggleGroup(group: string, isOpen: boolean) {
    this.showGroup = isOpen ? group : null;
  }
}
