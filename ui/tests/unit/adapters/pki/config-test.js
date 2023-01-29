import { module, test } from 'qunit';
import { setupTest } from 'vault/tests/helpers';
import { setupMirage } from 'ember-cli-mirage/test-support';
import { allowAllCapabilitiesStub } from 'vault/tests/helpers/stubs';

module('Unit | Adapter | pki/config', function (hooks) {
  setupTest(hooks);
  setupMirage(hooks);

  hooks.beforeEach(function () {
    this.store = this.owner.lookup('service:store');
    this.secretMountPath = this.owner.lookup('service:secret-mount-path');
    this.backend = 'pki-test';
    this.secretMountPath.currentPath = this.backend;
    this.server.post('/sys/capabilities-self', allowAllCapabilitiesStub());
  });

  test('it exists', function (assert) {
    const adapter = this.owner.lookup('adapter:pki/config');
    assert.ok(adapter);
  });

  module('formType import', function (hooks) {
    hooks.beforeEach(function () {
      this.payload = {
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
    });

    test('it calls the correct endpoint when useIssuer = false', async function (assert) {
      assert.expect(1);

      this.server.post(`${this.backend}/config/ca`, () => {
        assert.ok(true, 'request made to correct endpoint on create');
        return {};
      });

      await this.store
        .createRecord('pki/config', this.payload)
        .save({ adapterOptions: { formType: 'import', useIssuer: false } });
    });

    test('it calls the correct endpoint when useIssuer = true', async function (assert) {
      assert.expect(1);
      this.server.post(`${this.backend}/issuers/import/bundle`, () => {
        assert.ok(true, 'request made to correct endpoint on create');
        return {};
      });

      await this.store
        .createRecord('pki/config', this.payload)
        .save({ adapterOptions: { formType: 'import', useIssuer: true } });
    });
  });

  module('formType generate-root', function () {
    test('it calls the correct endpoint when useIssuer = false', async function (assert) {
      assert.expect(4);
      const adapterOptions = { adapterOptions: { formType: 'generate-root', useIssuer: false } };
      this.server.post(`${this.backend}/root/generate/internal`, () => {
        assert.ok(true, 'request made correctly when type = internal');
        return {};
      });
      this.server.post(`${this.backend}/root/generate/exported`, () => {
        assert.ok(true, 'request made correctly when type = exported');
        return {};
      });
      this.server.post(`${this.backend}/root/generate/existing`, () => {
        assert.ok(true, 'request made correctly when type = exising');
        return {};
      });
      this.server.post(`${this.backend}/root/generate/kms`, () => {
        assert.ok(true, 'request made correctly when type = kms');
        return {};
      });

      await this.store
        .createRecord('pki/config', {
          type: 'internal',
        })
        .save(adapterOptions);
      await this.store
        .createRecord('pki/config', {
          type: 'exported',
        })
        .save(adapterOptions);
      await this.store
        .createRecord('pki/config', {
          type: 'existing',
        })
        .save(adapterOptions);
      await this.store
        .createRecord('pki/config', {
          type: 'kms',
        })
        .save(adapterOptions);
    });

    test('it calls the correct endpoint when useIssuer = true', async function (assert) {
      assert.expect(4);
      const adapterOptions = { adapterOptions: { formType: 'generate-root', useIssuer: true } };
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
      this.server.post(`${this.backend}/issuers/generate/root/kms`, () => {
        assert.ok(true, 'request made correctly when type = kms');
        return {};
      });

      await this.store
        .createRecord('pki/config', {
          type: 'internal',
        })
        .save(adapterOptions);
      await this.store
        .createRecord('pki/config', {
          type: 'exported',
        })
        .save(adapterOptions);
      await this.store
        .createRecord('pki/config', {
          type: 'existing',
        })
        .save(adapterOptions);
      await this.store
        .createRecord('pki/config', {
          type: 'kms',
        })
        .save(adapterOptions);
    });
  });

  module('formType generate-csr', function () {
    test('it calls the correct endpoint when useIssuer = false', async function (assert) {
      assert.expect(4);
      const adapterOptions = { adapterOptions: { formType: 'generate-csr', useIssuer: false } };
      this.server.post(`${this.backend}/intermediate/generate/internal`, () => {
        assert.ok(true, 'request made correctly when type = internal');
        return {};
      });
      this.server.post(`${this.backend}/intermediate/generate/exported`, () => {
        assert.ok(true, 'request made correctly when type = exported');
        return {};
      });
      this.server.post(`${this.backend}/intermediate/generate/existing`, () => {
        assert.ok(true, 'request made correctly when type = exising');
        return {};
      });
      this.server.post(`${this.backend}/intermediate/generate/kms`, () => {
        assert.ok(true, 'request made correctly when type = kms');
        return {};
      });

      await this.store
        .createRecord('pki/config', {
          type: 'internal',
        })
        .save(adapterOptions);
      await this.store
        .createRecord('pki/config', {
          type: 'exported',
        })
        .save(adapterOptions);
      await this.store
        .createRecord('pki/config', {
          type: 'existing',
        })
        .save(adapterOptions);
      await this.store
        .createRecord('pki/config', {
          type: 'kms',
        })
        .save(adapterOptions);
    });

    test('it calls the correct endpoint when useIssuer = true', async function (assert) {
      assert.expect(4);
      const adapterOptions = { adapterOptions: { formType: 'generate-csr', useIssuer: true } };
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
      this.server.post(`${this.backend}/issuers/generate/intermediate/kms`, () => {
        assert.ok(true, 'request made correctly when type = kms');
        return {};
      });

      await this.store
        .createRecord('pki/config', {
          type: 'internal',
        })
        .save(adapterOptions);
      await this.store
        .createRecord('pki/config', {
          type: 'exported',
        })
        .save(adapterOptions);
      await this.store
        .createRecord('pki/config', {
          type: 'existing',
        })
        .save(adapterOptions);
      await this.store
        .createRecord('pki/config', {
          type: 'kms',
        })
        .save(adapterOptions);
    });
  });
});
