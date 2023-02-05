import { module, test } from 'qunit';
import { setupRenderingTest } from 'ember-qunit';
import { render, click, fillIn, find } from '@ember/test-helpers';
import { hbs } from 'ember-cli-htmlbars';
import { setupEngine } from 'ember-engines/test-support';
import { SELECTORS } from 'vault/tests/helpers/pki/pki-role-form';
import { setupMirage } from 'ember-cli-mirage/test-support';

module('Integration | Component | pki-role-form', function (hooks) {
  setupRenderingTest(hooks);
  setupMirage(hooks);
  setupEngine(hooks, 'pki'); // https://github.com/ember-engines/ember-engines/pull/653

  hooks.beforeEach(function () {
    this.store = this.owner.lookup('service:store');
    this.model = this.store.createRecord('pki/role');
    this.model.backend = 'pki';
  });

  test('it should render default fields and toggle groups', async function (assert) {
    assert.expect(13);
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
    assert.dom(SELECTORS.customTtl).exists('shows custom yielded form field');
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
  });

  test('it should save a new pki role with various options selected', async function (assert) {
    // Key usage, Key params and Not valid after options are tested in their respective component tests
    assert.expect(9);
    this.server.post(`/${this.model.backend}/roles/test-role`, (schema, req) => {
      assert.ok(true, 'Request made to save role');
      const request = JSON.parse(req.requestBody);
      const allowedDomainsTemplate = request.allowed_domains_template;
      const policyIdentifiers = request.policy_identifiers;
      const allowedUriSansTemplate = request.allow_uri_sans_template;
      const allowedSerialNumbers = request.allowed_serial_numbers;

      assert.true(allowedDomainsTemplate, 'correctly sends allowed_domains_template');
      assert.strictEqual(policyIdentifiers[0], 'some-oid', 'correctly sends policy_identifiers');
      assert.true(allowedUriSansTemplate, 'correctly sends allowed_uri_sans_template');
      assert.strictEqual(
        allowedSerialNumbers[0],
        'some-serial-number',
        'correctly sends allowed_serial_numbers'
      );
      return {};
    });

    this.onSave = () => assert.ok(true, 'onSave callback fires on save success');

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
    await click('[data-test-input="addBasicConstraints"]');
    await click(SELECTORS.domainHandling);
    await click('[data-test-input="allowedDomainsTemplate"]');
    await click(SELECTORS.policyIdentifiers);
    await fillIn('[data-test-input="policyIdentifiers"] [data-test-string-list-input="0"]', 'some-oid');
    await click(SELECTORS.san);
    await click('[data-test-input="allowUriSansTemplate"]');
    await click(SELECTORS.additionalSubjectFields);
    await fillIn(
      '[data-test-input="allowedSerialNumbers"] [data-test-string-list-input="0"]',
      'some-serial-number'
    );
    await click(SELECTORS.keyUsage);
    // check is flexbox by checking the height of the box
    const groupBoxHeight = find('[data-test-toggle-div="Key usage"]').clientHeight;
    assert.strictEqual(
      groupBoxHeight,
      567,
      'renders the correct height of the box element if the component is rending as a flexbox'
    );
    await click(SELECTORS.roleCreateButton);
  });

  test('meep it should rollback attributes or unload record on cancel', async function (assert) {
    assert.expect(1);
    this.onCancel = () => assert.ok(true, 'onCancel callback fires');
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

    await fillIn(SELECTORS.roleName, 'test-role');
    await click(SELECTORS.roleCancelButton);
    // console.log(this.model, 'model name');
    // await this.pauseTest();
  });

  // TODO: ('it should update role', async function (assert) {}

  /* FUTURE TEST TODO:
   * it should update role
   * it should unload the record on cancel
   */
});
