import { action } from '@ember/object';
import Component from '@glimmer/component';

export default class PkiSignIntermediateFormComponent extends Component {
  @action cancel() {
    // TODO
  }
  @action save() {
    // TODO
  }

  get groups() {
    return {
      'Signing options': ['usePss', 'skid', 'signatureBits'],
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
  }
}
