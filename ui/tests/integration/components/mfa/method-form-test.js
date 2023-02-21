import { module, test } from 'qunit';
import { setupRenderingTest } from 'ember-qunit';
import { render, fillIn, click } from '@ember/test-helpers';
import { hbs } from 'ember-cli-htmlbars';
import { setupMirage } from 'ember-cli-mirage/test-support';

module('Integration | Component | mfa-method-form', function (hooks) {
  setupRenderingTest(hooks);
  setupMirage(hooks);

  hooks.beforeEach(function () {
    this.store = this.owner.lookup('service:store');
    this.model = this.store.createRecord('mfa-method');
    this.model.type = 'totp';
  });

  test('it should render correct fields', async function (assert) {
    assert.expect(6);

    await render(hbs`
      <Mfa::MethodForm
        @model={{this.model}}
        @hasActions="true"
      />
      <div id="modal-wormhole"></div>
    `);
    assert.dom('[data-test-input="issuer"]').exists(`Issuer field input renders`);
    assert.dom('[data-test-input="period"]').exists('Period field ttl renders');
    assert.dom('[data-test-input="key_size"]').exists('Key size field input renders');
    assert.dom('[data-test-input="qr_size"]').exists('QR size field input renders');
    assert.dom('[data-test-input="algorithm"]').exists(`Algorithm field radio input renders`);
    assert
      .dom('[data-test-input="max_validation_attempts"]')
      .exists(`Max validation attempts field input renders`);
  });

  test('it should create new mfa method', async function (assert) {
    assert.expect(3);

    this.server.post('/identity/mfa/method/totp', () => {
      assert.ok(true, 'create request sent to server');
      return {};
    });

    await render(hbs`
      <Mfa::MethodForm
        @hasActions="true"
        @model={{this.model}}
        @onSave={{fn (mut this.didSave) true}}
      />
      <div id="modal-wormhole"></div>
    `);

    await fillIn('[data-test-input="issuer"]', 'Vault');
    await click('[data-test-mfa-save]');
    await fillIn('[data-test-confirmation-modal-input="Edit totp configuration?"]', 'totp');
    await click('[data-test-confirm-button="Edit totp configuration?"]');
    assert.true(this.didSave, 'onSave callback triggered');
    assert.strictEqual(this.model.issuer, 'Vault', 'Issuer property set on model');
  });

  test('it should populate form fields with model data', async function (assert) {
    assert.expect(3);

    this.model.issuer = 'Vault';
    this.model.period = '30s';
    this.model.algorithm = 'SHA512';

    await render(hbs`
      <Mfa::MethodForm
        @hasActions="true"
        @model={{this.model}}
      />
      <div id="modal-wormhole"></div>
    `);
    assert.dom('[data-test-input="issuer"]').hasValue('Vault', 'Issuer input is populated');
    assert.dom('[data-test-ttl-value="Period"]').hasValue('30', 'Period input ttl is populated');
    const checkedAlgorithm = this.element.querySelector('input[name=algorithm]:checked');
    assert.dom(checkedAlgorithm).hasValue('SHA512', 'SHA512 radio input is selected');
  });
});
