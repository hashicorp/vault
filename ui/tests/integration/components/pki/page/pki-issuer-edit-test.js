import { module, test } from 'qunit';
import { setupRenderingTest } from 'ember-qunit';
import { click, fillIn, render } from '@ember/test-helpers';
import { hbs } from 'ember-cli-htmlbars';
import { setupEngine } from 'ember-engines/test-support';
import { setupMirage } from 'ember-cli-mirage/test-support';
import { Response } from 'miragejs';
import sinon from 'sinon';

const selectors = {
  name: '[data-test-input="issuerName"]',
  leaf: '[data-test-field="leafNotAfterBehavior"] select',
  leafOption: '[data-test-field="leafNotAfterBehavior"] option',
  usageCert: '[data-test-usage="Issuing certificates"]',
  usageCrl: '[data-test-usage="Signing CRLs"]',
  usageOcsp: '[data-test-usage="Signing OCSPs"]',
  manualChain: '[data-test-input="manualChain"]',
  certUrls: '[data-test-string-list-input]',
  certUrl1: '[data-test-string-list-input="0"]',
  certUrl2: '[data-test-string-list-input="1"]',
  certUrlAdd: '[data-test-string-list-button="add"]',
  certUrlRemove: '[data-test-string-list-button="delete"]',
  crlDist: '[data-test-input="crlDistributionPoints"] [data-test-string-list-input="0"]',
  ocspServers: '[data-test-input="ocspServers"]  [data-test-string-list-input="0"]',
  save: '[data-test-save]',
  cancel: '[data-test-cancel]',
  error: '[data-test-error] p',
  alert: '[data-test-inline-error-message]',
};

module('Integration | Component | pki | Page::PkiIssuerEditPage::PkiIssuerEdit', function (hooks) {
  setupRenderingTest(hooks);
  setupEngine(hooks, 'pki');
  setupMirage(hooks);

  hooks.beforeEach(function () {
    const router = this.owner.lookup('service:router');
    const transitionSpy = sinon.stub(router, 'transitionTo');
    this.transitionCalled = () =>
      transitionSpy.calledWith('vault.cluster.secrets.backend.pki.issuers.issuer.details');

    const store = this.owner.lookup('service:store');
    store.pushPayload('pki/issuer', {
      modelName: 'pki/issuer',
      data: {
        issuer_id: 'test',
        issuer_name: 'foo-bar',
        leaf_not_after_behavior: 'err',
        usage: 'read-only,issuing-certificates,ocsp-signing',
        manual_chain: 'issuer_ref',
        issuing_certificates: ['http://localhost', 'http://localhost:8200'],
        crl_distribution_points: 'http://localhost',
        ocsp_servers: 'http://localhost',
      },
    });
    this.model = store.peekRecord('pki/issuer', 'test');
    // backend value on model pulled from secretMountPath service
    this.owner.lookup('service:secretMountPath').update('pki');

    this.update = async () => {
      await fillIn(selectors.name, 'bar-baz');
      await click(selectors.usageCrl);
      await click(selectors.certUrlRemove);
    };
  });

  test('it should populate fields with model values', async function (assert) {
    await render(hbs`<Page::PkiIssuerEdit @model={{this.model}} />`, { owner: this.engine });

    assert.dom(selectors.name).hasValue(this.model.issuerName, 'Issuer name field populates');
    assert
      .dom(selectors.leaf)
      .hasValue(this.model.leafNotAfterBehavior, 'Leaf not after behavior option selected');
    assert
      .dom(selectors.leafOption)
      .hasText(
        'Error if the computed NotAfter exceeds that of this issuer',
        'Correct text renders for leaf option'
      );
    assert.dom(selectors.usageCert).isChecked('Usage issuing certificates is checked');
    assert.dom(selectors.usageCrl).isNotChecked('Usage signing crls is not checked');
    assert.dom(selectors.usageOcsp).isChecked('Usage signing ocsps is checked');
    assert.dom(selectors.manualChain).hasValue(this.model.manualChain, 'Manual chain field populates');
    const certUrls = this.model.issuingCertificates.split(',');
    assert.dom(selectors.certUrl1).hasValue(certUrls[0], 'Issuing certificate populates');
    assert.dom(selectors.certUrl2).hasValue(certUrls[1], 'Issuing certificate populates');
    const crlDistributionPoints = this.model.crlDistributionPoints.split(',');
    assert.dom(selectors.crlDist).hasValue(crlDistributionPoints[0], 'Crl distribution points populate');
    const ocspServers = this.model.ocspServers.split(',');
    assert.dom(selectors.ocspServers).hasValue(ocspServers[0], 'Ocsp servers populate');
  });

  test('it should rollback model attributes on cancel', async function (assert) {
    await render(hbs`<Page::PkiIssuerEdit @model={{this.model}} />`, { owner: this.engine });

    await this.update();
    await click(selectors.cancel);

    assert.ok(this.transitionCalled(), 'Transitions to details route on cancel');
    assert.strictEqual(this.model.issuerName, 'foo-bar', 'Issuer name rolled back');
    assert.strictEqual(this.model.usage, 'read-only,issuing-certificates,ocsp-signing', 'Usage rolled back');
    assert.strictEqual(
      this.model.issuingCertificates,
      'http://localhost,http://localhost:8200',
      'Issuing certificates rolled back'
    );
  });

  test('it should update issuer', async function (assert) {
    assert.expect(4);

    this.server.post('/pki/issuer/test', (schema, req) => {
      const data = JSON.parse(req.requestBody);
      assert.strictEqual(data.issuer_name, 'bar-baz', 'Updated issuer name sent in POST request');
      assert.strictEqual(
        data.usage,
        'read-only,issuing-certificates,ocsp-signing,crl-signing',
        'Updated usage sent in POST request'
      );
      assert.strictEqual(
        data.issuing_certificates,
        'http://localhost:8200',
        'Updated issuing certificates sent in POST request'
      );
      return { data };
    });
    await render(hbs`<Page::PkiIssuerEdit @model={{this.model}} />`, { owner: this.engine });

    await this.update();
    await click(selectors.save);
    assert.ok(this.transitionCalled(), 'Transitions to details route on save success');
  });

  test('it should show error messages', async function (assert) {
    this.server.post('/pki/issuer/test', () => new Response(404, {}, { errors: ['Some error occurred'] }));

    await render(hbs`<Page::PkiIssuerEdit @model={{this.model}} />`, { owner: this.engine });
    await click(selectors.save);

    assert
      .dom(selectors.alert)
      .hasText('There was an error submitting this form.', 'Inline error alert renders');
    assert.dom(selectors.error).hasText('Some error occurred', 'Error message renders');
  });
});
