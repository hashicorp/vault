import { currentRouteName, currentURL, settled } from '@ember/test-helpers';
import { module, test } from 'qunit';
import { setupApplicationTest } from 'ember-qunit';
import editPage from 'vault/tests/pages/secrets/backend/pki/edit-role';
import listPage from 'vault/tests/pages/secrets/backend/list';
import generatePage from 'vault/tests/pages/secrets/backend/pki/generate-cert';
import configPage from 'vault/tests/pages/settings/configure-secret-backends/pki/section-cert';
import enablePage from 'vault/tests/pages/settings/mount-secret-backend';
import authPage from 'vault/tests/pages/auth';
import { SELECTORS } from 'vault/tests/helpers/pki';
import { csr } from 'vault/tests/helpers/pki/values';
import { v4 as uuidv4 } from 'uuid';

module('Acceptance | secrets/pki/list?tab=cert', function (hooks) {
  setupApplicationTest(hooks);

  hooks.beforeEach(function () {
    return authPage.login();
  });
  // important for this comment to stay here otherwise the formatting mangles the CSR
  // prettier-ignore
  const CSR = csr;

  // mount, generate CA, nav to create role page
  const setup = async (assert, action = 'issue') => {
    const path = `pki-${uuidv4()}`;
    const roleName = 'role';
    await enablePage.enable('pki', path);
    await settled();
    await configPage.visit({ backend: path }).form.generateCA();
    await settled();
    await editPage.visitRoot({ backend: path });
    await settled();
    await editPage.createRole('role', 'example.com');
    await settled();
    await generatePage.visit({ backend: path, id: roleName, action });
    await settled();
    return path;
  };

  test('it issues a cert', async function (assert) {
    assert.expect(10);
    const mount = await setup(assert);
    await settled();
    await generatePage.issueCert('foo');
    await settled();
    assert.strictEqual(currentURL(), `/vault/secrets/${mount}/credentials/role?action=issue`);
    assert.dom(SELECTORS.certificate).exists('displays masked certificate');
    assert.dom(SELECTORS.commonName).exists('displays common name');
    assert.dom(SELECTORS.issueDate).exists('displays issue date');
    assert.dom(SELECTORS.expiryDate).exists('displays expiration date');
    assert.dom(SELECTORS.issuingCa).exists('displays masked issuing CA');
    assert.dom(SELECTORS.privateKey).exists('displays masked private key');
    assert.dom(SELECTORS.serialNumber).exists('displays serial number');
    assert.dom(SELECTORS.caChain).exists('displays the CA chain');

    await generatePage.back();
    await settled();
    assert.notOk(generatePage.commonNameValue, 'the form is cleared');
  });

  test('it signs a csr', async function (assert) {
    assert.expect(4);
    await setup(assert, 'sign');
    await settled();
    await generatePage.sign('common', CSR);
    await settled();
    assert.ok(SELECTORS.certificate, 'displays masked certificate');
    assert.ok(SELECTORS.commonName, 'displays common name');
    assert.ok(SELECTORS.issuingCa, 'displays masked issuing CA');
    assert.ok(SELECTORS.serialNumber, 'displays serial number');
  });

  test('it views a cert', async function (assert) {
    assert.expect(12);
    const path = await setup(assert);
    await generatePage.issueCert('foo');
    await settled();
    await listPage.visitRoot({ backend: path, tab: 'cert' });
    await settled();
    assert.strictEqual(currentURL(), `/vault/secrets/${path}/list?tab=cert`);
    assert.strictEqual(listPage.secrets.length, 2, 'lists certs');
    await listPage.secrets.objectAt(0).click();
    await settled();
    assert.strictEqual(
      currentRouteName(),
      'vault.cluster.secrets.backend.show',
      'navigates to the show page'
    );
    assert.dom(SELECTORS.certificate).exists('displays masked certificate');
    assert.dom(SELECTORS.commonName).exists('displays common name');
    assert.dom(SELECTORS.issueDate).exists('displays issue date');
    assert.dom(SELECTORS.expiryDate).exists('displays expiration date');
    assert.dom(SELECTORS.serialNumber).exists('displays serial number');
    assert.dom(SELECTORS.revocationTime).doesNotExist('does not display revocation time of 0');
    assert.dom(SELECTORS.issuingCa).doesNotExist('does not display empty issuing CA');
    assert.dom(SELECTORS.caChain).doesNotExist('does not display empty CA chain');
    assert.dom(SELECTORS.privateKey).doesNotExist('does not display empty private key');
  });
});
