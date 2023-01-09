import { module, test } from 'qunit';
import { setupTest } from 'vault/tests/helpers';
import { setupMirage } from 'ember-cli-mirage/test-support';

module('Unit | Adapter | pki/config', function (hooks) {
  setupTest(hooks);
  setupMirage(hooks);

  hooks.beforeEach(function () {
    this.store = this.owner.lookup('service:store');
    this.secretMountPath = this.owner.lookup('service:secret-mount-path');
    this.backend = 'pki-test';
    this.secretMountPath.currentPath = this.backend;
  });

  test('it exists', function (assert) {
    const adapter = this.owner.lookup('adapter:pki/config');
    assert.ok(adapter);
  });

  test('it calls the correct endpoint on create when formType is import', async function (assert) {
    assert.expect(1);
    const model = {
      formType: 'import',
      pem_bundle: `-----BEGIN CERTIFICATE REQUEST-----
      MIIChDCCAWwCAQAwFjEUMBIGA1UEAxMLZXhhbXBsZS5jb20wggEiMA0GCSqGSIb3
      DQEBAQUAA4IBDwAwggEKAoIBAQCuW9C58M1wO0vdGmtLcJbbCkKyfsHJJae1j4LL
      xdGqs1j9UKD66UALSzZEeMCBdtTNNzThAgYJqCSA5swqpbRf6WZ3K/X7oHbfcrHi
      SAm8v/0QsJDF5Rphiy6wyggaoaHEsbSp83kYy9r+h48vFW5Dr8UvJTsp5kdRn31L
      bTHr56iqOaHQbu6hDj4Ompg/0OElPH1tV2X947o8timR+L89utZzR+d8x/eeTdPl
      H7TEkMEomRvt7NTRHGYRsm3Gzq4AA6PalzIxzwJrNgXfJDutNn/QwcVd5sImwYCO
      GaHsOvGfc02w+Vqqva9EOEQSr6B90kA+vc30I6uSiugzV9TFAgMBAAGgKTAnBgkq
      hkiG9w0BCQ4xGjAYMBYGA1UdEQQPMA2CC2V4YW1wbGUuY29tMA0GCSqGSIb3DQEB
      CwUAA4IBAQAjm6JTU7axU6TzLlXlOp7hZ4+nep2/8vvJ9EOXzL8x/qtTTizctdG9
      Op70gywoUxAS2tatwa4fmW9DbA2eGiLU+Ibj/5b0Veq5DQdp1Qg3MLBP/+AcM/7m
      rrgA9MhkpQahXCj4vXo6NeXYaTh6Jo/s8C9h3WxTD6ptDMiaPFcEuWcx0e3AjjH0
      pe7k9/MfB2wLfQ7+5wee/tCFWZN4tk8YfjQeQA1extXYKM/f8eu3Z/wjbbMOVpwb
      xst+VTY7X9T8cU/hjDEoNG677meI+W5MgiwX0rxTpoz991fqr3vp7PELYj3GMyii
      D1YfvqXieNij4UrduRqCXj1m8SVZlM+X
      -----END CERTIFICATE REQUEST-----`,
    };
    this.server.post(`${this.backend}/issuers/import/bundle`, () => {
      assert.ok(true, 'request made to correct endpoint on create');
      return {};
    });

    await this.store.createRecord('pki/config', model).save();
  });

  test('it calls the correct endpoint on create when formType is generate-root', async function (assert) {
    assert.expect(3);
    this.server.post(`${this.backend}/issuers/generate/root/internal`, () => {
      assert.ok(true, 'request made correctly when type = internal');
      return {};
    });
    this.server.post(`${this.backend}/issuers/generate/root/exported`, () => {
      assert.ok(true, 'request made correctly when type = exported');
      return {};
    });
    this.server.post(`${this.backend}/issuers/generate/root/existing`, () => {
      assert.ok(true, 'request made correctly when type = exising');
      return {};
    });

    await this.store
      .createRecord('pki/config', {
        formType: 'generate-root',
        type: 'internal',
      })
      .save();
    await this.store
      .createRecord('pki/config', {
        formType: 'generate-root',
        type: 'exported',
      })
      .save();
    await this.store
      .createRecord('pki/config', {
        formType: 'generate-root',
        type: 'existing',
      })
      .save();
  });

  test('it calls the correct endpoint on create when formType is generate-csr', async function (assert) {
    assert.expect(3);
    this.server.post(`${this.backend}/issuers/generate/intermediate/internal`, () => {
      assert.ok(true, 'request made correctly when type = internal');
      return {};
    });
    this.server.post(`${this.backend}/issuers/generate/intermediate/exported`, () => {
      assert.ok(true, 'request made correctly when type = exported');
      return {};
    });
    this.server.post(`${this.backend}/issuers/generate/intermediate/existing`, () => {
      assert.ok(true, 'request made correctly when type = exising');
      return {};
    });

    await this.store
      .createRecord('pki/config', {
        formType: 'generate-csr',
        type: 'internal',
      })
      .save();
    await this.store
      .createRecord('pki/config', {
        formType: 'generate-csr',
        type: 'exported',
      })
      .save();
    await this.store
      .createRecord('pki/config', {
        formType: 'generate-csr',
        type: 'existing',
      })
      .save();
  });
});
