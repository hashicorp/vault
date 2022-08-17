import Model, { attr } from '@ember-data/model';
import lazyCapabilities, { apiPath } from 'vault/macros/lazy-capabilities';

export default class OidcKeyModel extends Model {
  @attr('string') name;

  @attr('string', {
    defaultValue: 'RS256',
    possibleValues: ['RS256', 'RS384', 'RS512', 'ES256', 'ES384', 'ES512', 'EdDSA'],
  })
  algorithm;

  @attr({ editType: 'ttl', helperTextDisabled: 'Vault will use the default lease duration' }) rotationPeriod;

  @attr({
    label: 'Verification TTL',
    editType: 'ttl',
    helperTextDisabled: 'Vault will use the default lease duration',
    hideToggle: true,
  })
  verificationTtl;

  @attr('array', { label: 'Allowed applications' }) allowedClientIds; // no editType because does not use form-field component

  @lazyCapabilities(apiPath`identity/oidc/key/${'name'}`, 'name') keyPath;
  @lazyCapabilities(apiPath`identity/oidc/key/${'name'}/rotate`, 'name') rotatePath;
  @lazyCapabilities(apiPath`identity/oidc/key`) keysPath;
  get canCreate() {
    return this.keyPath.get('canCreate');
  }
  get canRead() {
    return this.keyPath.get('canRead');
  }
  get canEdit() {
    return this.keyPath.get('canUpdate');
  }
  get canRotate() {
    return this.rotatePath.get('canUpdate');
  }
  get canDelete() {
    return this.keyPath.get('canDelete');
  }
  get canList() {
    return this.keysPath.get('canList');
  }
}
