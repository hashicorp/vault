import { module, test } from 'qunit';
import { setupTest } from 'ember-qunit';
import { DOMAIN_STRINGS, PROVIDER_WITH_LOGO } from 'vault/models/role-jwt';

module('Unit | Model | role-jwt', function(hooks) {
  setupTest(hooks);

  test('it exists', function(assert) {
    let model = this.owner.lookup('service:store').createRecord('role-jwt');
    assert.ok(!!model);
    assert.equal(model.providerName, null, 'no providerName');
    assert.equal(model.providerButtonComponent, null, 'no providerButtonComponent');
  });

  test('it computes providerName when known provider url match fails', function(assert) {
    let model = this.owner.lookup('service:store').createRecord('role-jwt', {
      authUrl: 'http://example.com',
    });

    assert.equal(model.providerName, null, 'no providerName');
    assert.equal(model.providerButtonComponent, null, 'no providerButtonComponent');
  });

  test('it provides a providerName for listed known providers', function(assert) {
    Object.keys(DOMAIN_STRINGS).forEach(domainPart => {
      let model = this.owner.lookup('service:store').createRecord('role-jwt', {
        authUrl: `http://provider-${domainPart}.com`,
      });

      let expectedName = DOMAIN_STRINGS[domainPart];
      assert.equal(model.providerName, expectedName, `computes providerName: ${expectedName}`);
      if (PROVIDER_WITH_LOGO.includes(expectedName)) {
        assert.equal(
          model.providerButtonComponent,
          `auth-button-${domainPart}`,
          `computes providerButtonComponent: ${domainPart}`
        );
      } else {
        assert.equal(model.providerButtonComponent, null, `computes providerButtonComponent: ${domainPart}`);
      }
    });
  });
});
