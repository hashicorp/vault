import { module, test } from 'qunit';
import { setupRenderingTest } from 'vault/tests/helpers';
import { render } from '@ember/test-helpers';
import { hbs } from 'ember-cli-htmlbars';
import { setupEngine } from 'ember-engines/test-support';
import { parseCertificate } from 'vault/utils/parse-pki-cert';
import { unsupportedOids } from 'vault/tests/helpers/pki/values';

module('Integration | Component | parsed-certificate-info-rows', function (hooks) {
  setupRenderingTest(hooks);
  setupEngine(hooks, 'pki');

  test('it renders nothing if no valid attributes passed', async function (assert) {
    this.set('parsedCertificate', {
      foo: '',
      common_name: 'not-shown',
    });
    await render(hbs`<ParsedCertificateInfoRows @model={{this.parsedCertificate}} />`, {
      owner: this.engine,
    });

    assert.dom(this.element).hasText('');
  });

  test('it renders only valid attributes with values', async function (assert) {
    this.set('parsedCertificate', {
      common_name: 'not-shown',
      use_pss: false,
      alt_names: ['something', 'here'],
      ttl: undefined,
    });
    await render(hbs`<ParsedCertificateInfoRows @model={{this.parsedCertificate}} />`, {
      owner: this.engine,
    });
    assert.dom('[data-test-component="info-table-row"]').exists({ count: 2 }, 'renders 2 rows');

    assert.dom('[data-test-value-div="Common name"]').doesNotExist('common name is never rendered');
    assert.dom('[data-test-row-value="Subject Alternative Names (SANs)"]').hasText('something,here');
    assert.dom('[data-test-value-div="Use PSS"]').hasText('No', 'Booleans are rendered');
    assert.dom('[data-test-value-div="ttl"]').doesNotExist('ttl is not rendered because value undefined');
    assert.dom('[data-test-alert-banner="alert"]').doesNotExist('does not render parsing error info banner');
  });

  test('it renders info banner when parsing fails and no parsing errors', async function (assert) {
    this.set('parsedCertificate', {
      can_parse: false,
    });
    await render(hbs`<ParsedCertificateInfoRows @model={{this.parsedCertificate}} />`, {
      owner: this.engine,
    });

    assert
      .dom('[data-test-alert-banner="alert"]')
      .hasText(
        `There was an error parsing certificate metadata Vault cannot display unparsed values, but this will not interfere with the certificate's functionality.`
      );
  });

  test('it renders info banner when parsing fails and parsing errors exist', async function (assert) {
    this.set('parsedCertificate', {
      can_parse: false,
      parsing_errors: [new Error('some parsing error')],
    });
    await render(hbs`<ParsedCertificateInfoRows @model={{this.parsedCertificate}} />`, {
      owner: this.engine,
    });

    assert
      .dom('[data-test-alert-banner="alert"]')
      .hasText(
        `There was an error parsing certificate metadata Vault cannot display unparsed values, but this will not interfere with the certificate's functionality. Parsing error(s): some parsing error`
      );
  });

  test('it renders info banner when parsing is successful but unsupported OIDs return parsing errors', async function (assert) {
    const { parsing_errors } = parseCertificate(unsupportedOids);
    this.set('parsedCertificate', {
      can_parse: true,
      parsing_errors,
    });
    await render(hbs`<ParsedCertificateInfoRows @model={{this.parsedCertificate}} />`, {
      owner: this.engine,
    });

    assert
      .dom('[data-test-alert-banner="alert"]')
      .hasText(
        `There was an error parsing certificate metadata Vault cannot display unparsed values, but this will not interfere with the certificate's functionality. Parsing error(s): certificate contains unsupported subject OIDs: 1.2.840.113549.1.9.1, certificate contains unsupported extension OIDs: 2.5.29.37`
      );
  });
});
