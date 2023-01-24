/* eslint qunit/no-conditional-assertions: "warn" */
import { module, test } from 'qunit';
import { setupTest } from 'ember-qunit';
import { DOMAIN_STRINGS, PROVIDER_WITH_LOGO } from 'vault/models/role-jwt';

module('Unit | Model | role-jwt', function (hooks) {
  setupTest(hooks);

  test('it exists', function (assert) {
    const model = this.owner.lookup('service:store').createRecord('role-jwt');
    assert.ok(!!model);
    assert.strictEqual(model.providerName, null, 'no providerName');
    assert.strictEqual(model.providerButtonComponent, null, 'no providerButtonComponent');
  });

  test('it computes providerName when known provider url match fails', function (assert) {
    const model = this.owner.lookup('service:store').createRecord('role-jwt', {
      authUrl: 'http://example.com',
    });

    assert.strictEqual(model.providerName, null, 'no providerName');
    assert.strictEqual(model.providerButtonComponent, null, 'no providerButtonComponent');
  });

  test('it provides a providerName for listed known providers', function (assert) {
    assert.expect(12);
    Object.keys(DOMAIN_STRINGS).forEach((domainPart) => {
      const model = this.owner.lookup('service:store').createRecord('role-jwt', {
        authUrl: `http://provider-${domainPart}.com`,
      });

      const expectedName = DOMAIN_STRINGS[domainPart];
      assert.strictEqual(model.providerName, expectedName, `computes providerName: ${expectedName}`);
      if (PROVIDER_WITH_LOGO.includes(expectedName)) {
        assert.strictEqual(
          model.providerButtonComponent,
          `auth-button-${domainPart}`,
          `computes providerButtonComponent: ${domainPart}`
        );
      } else {
        assert.strictEqual(
          model.providerButtonComponent,
          null,
          `computes providerButtonComponent: ${domainPart}`
        );
      }
    });
  });
});
