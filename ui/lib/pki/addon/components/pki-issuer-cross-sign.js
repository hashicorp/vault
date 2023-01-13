import { action } from '@ember/object';
import Component from '@glimmer/component';

/**
 * @module PkiIssuerCrossSign
 * PkiIssuerCrossSign components are used to...
 *
 * @example
 * ```js
 * <PkiIssuerCrossSign @requiredParam={requiredParam} @optionalParam={optionalParam} @param1={{param1}}/>
 * ```
 * @param {object} requiredParam - requiredParam is...
 * @param {string} [optionalParam] - optionalParam is...
 * @param {string} [param1=defaultValue] - param1 is...
 */

export default class PkiIssuerCrossSign extends Component {
  // TODO key names will map to model attrs?
  objectKeys = [
    { label: 'Mount path', key: 'intermediateMount', placeholder: 'Mount path' },
    { label: "Issuer's current name", key: 'intermediateName', placeholder: 'Current issuer name' },
    { label: 'New issuer name', key: 'newCertName', placeholder: 'Enter a new issuer name' },
  ];

  @action
  handleChange(stuff) {
    // do fancy cross-sign stuff
    // eslint-disable-next-line no-console
    console.log(stuff);
  }
}
