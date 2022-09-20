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

  @lazyCapabilities(apiPath`${'backend'}/issue/${'id'}`, 'backend', 'id') generatePath;
  get canReadIssue() {
    // ARG TODO was duplicate name, added Issue
    return this.generatePath.get('canUpdate');
  }
  @lazyCapabilities(apiPath`${'backend'}/sign/${'id'}`, 'backend', 'id') signPath;
  get canSign() {
    return this.signPath.get('canUpdate');
  }
  @lazyCapabilities(apiPath`${'backend'}/sign-verbatim/${'id'}`, 'backend', 'id') signVerbatimPath;
  get canSignVerbatim() {
    return this.signVerbatimPath.get('canUpdate');
  }

  // Form Fields not hidden in toggle options
  _attributeMeta = null;
  get formFields() {
    if (!this._attributeMeta) {
      this._attributeMeta = expandAttributeMeta(this, ['name', 'clientType', 'redirectUris']);
    }
    return this._attributeMeta;
  }

  // Form fields hidden behind toggle options
  _fieldToAttrsGroups = null;
  // ARG TODO: I removed 'allowedDomains' but I'm fairly certain it needs to be somewhere. Confirm with design.
  get fieldGroups() {
    if (!this._fieldToAttrsGroups) {
      this._fieldToAttrsGroups = fieldToAttrs(this, [
        { default: ['name'] },
        {
          'Domain handling': [
            'allowedDomains',
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
            'extKeyUsage', // ARG TODO: takes a list, but we have these as checkboxes from the options on the golang site: https://pkg.go.dev/crypto/x509#ExtKeyUsage
          ],
        },
        { 'Policy identifiers': ['policyIdentifiers'] },
        {
          'Subject Alternative Name (SAN) Options': ['allowIpSans', 'allowedUriSans', 'allowedOtherSans'],
        },
        {
          'Additional subject fields': [
            'allowed_serial_numbers',
            'requireCn',
            'useCsrCommonName',
            'useCsrSans',
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
