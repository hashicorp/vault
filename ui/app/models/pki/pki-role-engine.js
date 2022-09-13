import Model, { attr } from '@ember-data/model';
import lazyCapabilities, { apiPath } from 'vault/macros/lazy-capabilities';
import { expandAttributeMeta } from 'vault/utils/field-to-attrs';
import { withModelValidations } from 'vault/decorators/model-validations';

import fieldToAttrs from 'vault/utils/field-to-attrs';
// import { combineFieldGroups } from 'vault/utils/openapi-to-attrs'; // ARG TODO likely I'll need this.

const validations = {
  name: [
    { type: 'presence', message: 'Name is required.' },
    {
      type: 'containsWhiteSpace',
      message: 'Name cannot contain whitespace.',
    },
  ],
};

@withModelValidations(validations)
export default class PkiRolesModel extends Model {
  @attr('string', { readOnly: true }) backend;
  @attr('string', {
    label: 'Role name',
    fieldValue: 'id',
    readOnly: true,
  })
  name;
  // ARG TODO return to
  useOpenAPI = true;
  @lazyCapabilities(apiPath`${'backend'}/roles/${'id'}`, 'backend', 'id') updatePath;
  get canDelete() {
    return this.updatePath.get('canCreate');
  }
  get canEdit() {
    return this.updatePath.get('canEdit');
  }
  get canRead() {
    return this.updatePath.get('canRead');
  }
  // ARG TODO confirm you need this
  @lazyCapabilities(apiPath`${'backend'}/issue/${'id'}`, 'backend', 'id') generatePath;
  get canReadIssue() {
    // ARG TODO was duplicate name, added Issue
    return this.generatePath.get('canUpdate');
  }
  // ARG TODO confirm you need this
  @lazyCapabilities(apiPath`${'backend'}/sign/${'id'}`, 'backend', 'id') signPath;
  get canSign() {
    return this.signPath.get('canUpdate');
  }
  // ARG TODO confirm you need this
  @lazyCapabilities(apiPath`${'backend'}/sign-verbatim/${'id'}`, 'backend', 'id') signVerbatimPath;
  get canSignVerbatim() {
    return this.signVerbatimPath.get('canUpdate');
  }

  // Form Fields not hidden in toggle options
  _attributeMeta = null; // cache initial result of expandAttributeMeta in getter and return
  get formFields() {
    if (!this._attributeMeta) {
      this._attributeMeta = expandAttributeMeta(this, ['name', 'clientType', 'redirectUris']);
    }
    return this._attributeMeta;
  }

  // Form fields hidden behind toggle options
  _fieldToAttrsGroups = null;
  // ARG TODO: I removed 'allowedDomains' but I fairly certain it needs to be somewhere
  get fieldGroups() {
    if (!this._fieldToAttrsGroups) {
      this._fieldToAttrsGroups = fieldToAttrs(this, [
        { default: ['name'] },
        {
          'Domain handling': [
            'allowedDomainTemplate',
            'allowBareDomains',
            'allowSubdomains',
            'allowGlobDomains',
            'allowWildcardCertificates',
            'allowLocalhost',
            'allowAnyName',
            'enforceHostnames',
          ],
        },
        {
          'Key parameters': [
            'keyType',
            'keyBits',
            'signatureBits', // ARG wasn't original a param
          ],
        },
        {
          'Key usage': [
            // ARG TODO Come back to this as there are a lot more here I don't have.
            'DigitalSignature', // ARG it's capitalized in the docs, but confirm
            'KeyAgreement',
            'KeyEncipherment', // ARG wasn't original a param
          ],
        },
        { 'Policy identifiers': ['policy_identifiers'] },
        {
          'Subject Alternative Name (SAN) Options': [
            'allow_ip_sans',
            'allowed_uri_sans',
            'allowed_other_sans',
          ],
        },
        {
          'Additional subject fields': [
            'allowed_serial_numbers',
            'require_cn',
            'use_csr_common_name',
            'ou',
            'organization',
            'country',
            'locality',
            'province',
            'streetAddress',
            'postalCode',
          ],
        },
      ]);
    }
    return this._fieldToAttrsGroups;
  }
}
