import { attr } from '@ember-data/model';
import { copy } from 'ember-copy';
import { computed } from '@ember/object';
import Certificate from './pki/cert';
import { combineFieldGroups } from 'vault/utils/openapi-to-attrs';

export default Certificate.extend({
  signVerbatim: attr('boolean', {
    readOnly: true,
    defaultValue: false,
  }),
  useOpenAPI: true,
  csr: attr('string', {
    label: 'Certificate Signing Request (CSR)',
    editType: 'textarea',
  }),

  fieldGroups: computed('newFields', 'signVerbatim', function () {
    const options = { Options: ['altNames', 'ipSans', 'ttl', 'excludeCnFromSans', 'otherSans'] };
    let groups = [
      {
        default: ['csr', 'commonName', 'format', 'signVerbatim'],
      },
    ];
    if (this.newFields) {
      groups = combineFieldGroups(groups, this.newFields, []);
    }
    if (this.signVerbatim === false) {
      groups.push(options);
    }

    return this.fieldsToAttrs(copy(groups, true));
  }),
});
