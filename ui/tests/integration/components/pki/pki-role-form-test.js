import { module, test } from 'qunit';
import { setupRenderingTest } from 'ember-qunit';
import { render, click, fillIn } from '@ember/test-helpers';
import { hbs } from 'ember-cli-htmlbars';
import { setupEngine } from 'ember-engines/test-support';
import { SELECTORS, overrideCapabilities, PKI_BASE_URL, clearRecord } from 'vault/tests/helpers/pki-engine';
import { create } from 'ember-cli-page-object';
import { setupMirage } from 'ember-cli-mirage/test-support';
import fm from 'vault/tests/pages/components/flash-message';

const flashMessage = create(fm);

module('Integration | Component | pki/role-form', function (hooks) {
  setupRenderingTest(hooks);
  setupEngine(hooks, 'pki'); // https://github.com/ember-engines/ember-engines/pull/653
  setupMirage(hooks);

  hooks.beforeEach(function () {
    this.store = this.owner.lookup('service:store');
    this.model = this.store.createRecord('pki/pki-role-engine');
    this.server.post('/sys/capabilities-self', () => ({}));
    this.model.backend = 'pki';
  });

  hooks.after(function () {
    //  ARG TODO unsure if need to destroy record.
  });

  test('it should render default fields and toggle groups', async function (assert) {
    assert.expect(12);
    await render(
      hbs`
      <PkiRoleForm
         @model={{this.model}}
         @onCancel={{this.onCancel}}
         @onSave={{this.onSave}}
       />
  `,
      { owner: this.engine }
    );
    assert.dom(SELECTORS.issuerRef).exists('shows form-field issuer ref');
    assert.dom(SELECTORS.backdateValidity).exists('shows form-field backdate validity');
    assert.dom(SELECTORS.maxTtl).exists('shows form-field max ttl');
    assert.dom(SELECTORS.generateLease).exists('shows form-field generateLease');
    assert.dom(SELECTORS.noStore).exists('shows form-field no store');
    assert.dom(SELECTORS.addBasicConstraints).exists('shows form-field add basic constraints');
    assert.dom(SELECTORS.domainHandling).exists('shows form-field group add domain handling');
    assert.dom(SELECTORS.keyParams).exists('shows form-field group key params');
    assert.dom(SELECTORS.keyUsage).exists('shows form-field group key usage');
    assert.dom(SELECTORS.policyIdentifiers).exists('shows form-field group policy identifiers');
    assert.dom(SELECTORS.san).exists('shows form-field group SAN');
    assert.dom(SELECTORS.additionalSubjectFields).exists('shows form-field group additional subject fields');

    //* clean up test state
    await clearRecord(this.store, 'pki/pki-role-engine', 'test-role');
  });

  test('it should save a new pki role meep', async function (assert) {
    assert.expect(12);
    await render(
      hbs`
      <PkiRoleForm
         @model={{this.model}}
         @onCancel={{this.onCancel}}
         @onSave={{this.onSave}}
       />
  `,
      { owner: this.engine }
    );

    await click(SELECTORS.roleCreateButton);
    assert
      .dom(SELECTORS.roleName)
      .hasClass('has-error-border', 'shows border error on role name field when no role name is submitted');
    assert
      .dom('[data-test-inline-error-message]')
      .includesText('Name is required.', 'show correct error message');

    await fillIn(SELECTORS.roleName, 'test-role');
    this.server.post('/sys/capabilities-self', () => overrideCapabilities(PKI_BASE_URL));
    // await click(SELECTORS.roleCreateButton);
    // // must override capabilities permissions // ARG TODO why? ask claire

    // await this.pauseTest();
    assert.strictEqual(
      flashMessage.latestMessage,
      'Successfully created the role test-role.',
      'renders success flash upon role creation'
    );
    assert.strictEqual(
      currentRouteName(),
      'vault.cluster.secrets.backend.pki.pki.roles.role.details',
      'navigates to role detail view after save'
    );
  });
});
