import { currentRouteName } from '@ember/test-helpers';
import { module, test } from 'qunit';
import { setupApplicationTest } from 'ember-qunit';
import editPage from 'vault/tests/pages/secrets/backend/pki/edit-role';
import listPage from 'vault/tests/pages/secrets/backend/list';
import generatePage from 'vault/tests/pages/secrets/backend/pki/generate-cert';
import showPage from 'vault/tests/pages/secrets/backend/pki/show';
import configPage from 'vault/tests/pages/settings/configure-secret-backends/pki/section-cert';
import enablePage from 'vault/tests/pages/settings/mount-secret-backend';
import authPage from 'vault/tests/pages/auth';

module('Acceptance | secrets/pki/list?tab=certs', function(hooks) {
  setupApplicationTest(hooks);

  hooks.beforeEach(function() {
    return authPage.login();
  });
  // important for this comment to stay here otherwise the formatting mangles the CSR
  // prettier-ignore
  const CSR = `-----BEGIN CERTIFICATE REQUEST-----
MIICdDCCAVwCAQAwDjEMMAoGA1UEAxMDbG9sMIIBIjANBgkqhkiG9w0BAQEFAAOC
AQ8AMIIBCgKCAQEA4Dz2b/nAP/M6bqyk5mctqqYAAcoME//xPBy0wREHuZ776Pu4
l45kDL3dPXiY8U2P9pn8WIr2KpLK6oWUfSsiG2P082bpWDL20UymkWqDhhrA4unf
ZRq68UIDbcetlLw15YKnlNdvNZ7Qr8Se8KV0YGR/wFqI7QfS6VE3lhxZWEBUayI0
egqOuDbXAcZTON1AZ92/F+WFSbc43iYdDk16XfAPFKhtvLr6zQQuzebAb7HG04Hc
GhRskixxyJ8XY6XUplfsa1HcpUXE4f1GeUvq3g6ltVCSJ0p7qI9FFjV4t+DCLVVV
LnwHUi9Vzz6i2wjMt7P6+gHR+RrOWBgRMn38fwIDAQABoCEwHwYJKoZIhvcNAQkO
MRIwEDAOBgNVHREEBzAFggNsb2wwDQYJKoZIhvcNAQELBQADggEBAAm3AHQ1ctdV
8HCrMOXGVLgI2cB1sFd6VYVxPBxIk812Y4wyO8Q6POE5VZNTIgMcSeIaFu5lgHNL
Peeb54F+zEa+OJYkcWgCAX5mY/0HoML4p2bxFTSjllSpcX7ktjq4IEIY/LRpqSgc
jgZHHRwanFfkeIOhN4Q5qJWgBPNhDAcNPE7T0M/4mxqYDqMSJvMYmC67hq1UOOug
/QVDUDJRC1C0aDw9if+DbG/bt1V6HpMQhDIEUjzfu4zG8pcag3cJpOA8JhW1hnG0
XA2ZOCA7s34/szr2FczXtIoKiYmv3UzPyO9/4mc0Q2+/nR4CG8NU9WW/XJCne9ID
elRplAzrMF4=
-----END CERTIFICATE REQUEST-----`;

  // mount, generate CA, nav to create role page
  const setup = async (assert, action = 'issue') => {
    const path = `pki-${new Date().getTime()}`;
    const roleName = 'role';
    await enablePage.enable('pki', path);
    await configPage.visit({ backend: path }).form.generateCA();
    await editPage.visitRoot({ backend: path });
    await editPage.createRole('role', 'example.com');
    await generatePage.visit({ backend: path, id: roleName, action });
    return path;
  };

  test('it issues a cert', async function(assert) {
    await setup(assert);

    await generatePage.issueCert('foo');
    assert.ok(generatePage.hasCert, 'displays the cert');

    await generatePage.back();
    assert.notOk(generatePage.commonNameValue, 'the form is cleared');
  });

  test('it signs a csr', async function(assert) {
    await setup(assert, 'sign');
    await generatePage.sign('common', CSR);
    assert.ok(generatePage.hasCert, 'displays the cert');
  });

  test('it views a cert', async function(assert) {
    const path = await setup(assert);
    await generatePage.issueCert('foo');
    await listPage.visitRoot({ backend: path, tab: 'certs' });
    assert.ok(listPage.secrets.length > 0, 'lists certs');

    await listPage.secrets.objectAt(0).click();
    assert.equal(currentRouteName(), 'vault.cluster.secrets.backend.show', 'navigates to the show page');
    assert.ok(showPage.hasCert, 'shows the cert');
  });
});
