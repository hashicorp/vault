import Model, { attr } from '@ember-data/model';
import lazyCapabilities, { apiPath } from 'vault/macros/lazy-capabilities';
import { expandAttributeMeta } from 'vault/utils/field-to-attrs';
import { withModelValidations } from 'vault/decorators/model-validations';

import fieldToAttrs from 'vault/utils/field-to-attrs';

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
export default class PkiRolesEngineModel extends Model {
  @attr('string', { readOnly: true }) backend;
  @attr('string', {
    label: 'Role name',
    fieldValue: 'id',
    readOnly: true,
  })
  name;

  // must be a getter so it can be added to the prototype needed in the pathHelp service on the line here: if (newModel.merged || modelProto.useOpenAPI !== true) {
  get useOpenAPI() {
    return true;
  }
  getHelpUrl(backend) {
    return `/v1/${backend}/roles/example?help=1`;
  }
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
          'Key parameters': ['keyType', 'keyBits', 'signatureBits'],
        },
        {
          'Key usage': [
            'DigitalSignature', // ARG TODO: capitalized in the docs, but should confirm
            'KeyAgreement',
            'KeyEncipherment',
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
