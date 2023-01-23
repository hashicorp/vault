import Component from '@glimmer/component';
import { tracked } from '@glimmer/tracking';
import { action } from '@ember/object';
import { keyParamsByType } from 'pki/utils/action-params';
import PkiActionModel from 'vault/models/pki/action';

interface Args {
  model: PkiActionModel;
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
    return {
      'Key parameters': this.keyParamFields,
      'Subject Alternative Name (SAN) Options': [
        'excludeCnFromSans',
        'serialNumber',
        'altNames',
        'ipSans',
        'uriSans',
        'otherSans',
      ],
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
  }

  @action
  toggleGroup(group: string, isOpen: boolean) {
    this.showGroup = isOpen ? group : null;
  }
}
