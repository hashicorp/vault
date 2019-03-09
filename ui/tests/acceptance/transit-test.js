import { click, fillIn, find, currentURL, settled, visit } from '@ember/test-helpers';
import { module, test } from 'qunit';
import { setupApplicationTest } from 'ember-qunit';
import { encodeString } from 'vault/utils/b64';
import authPage from 'vault/tests/pages/auth';
import enablePage from 'vault/tests/pages/settings/mount-secret-backend';

const keyTypes = [
  {
    name: ts => `aes-${ts}`,
    type: 'aes256-gcm96',
    exportable: true,
    supportsEncryption: true,
  },
  {
    name: ts => `aes-convergent-${ts}`,
    type: 'aes256-gcm96',
    convergent: true,
    supportsEncryption: true,
  },
  {
    name: ts => `chacha-${ts}`,
    type: 'chacha20-poly1305',
    exportable: true,
    supportsEncryption: true,
  },
  {
    name: ts => `chacha-convergent-${ts}`,
    type: 'chacha20-poly1305',
    convergent: true,
    supportsEncryption: true,
  },
  {
    name: ts => `ecdsa-${ts}`,
    type: 'ecdsa-p256',
    exportable: true,
    supportsSigning: true,
  },
  {
    name: ts => `ed25519-${ts}`,
    type: 'ed25519',
    derived: true,
    supportsSigning: true,
  },
  {
    name: ts => `rsa-2048-${ts}`,
    type: `rsa-2048`,
    supportsSigning: true,
    supportsEncryption: true,
  },
  {
    name: ts => `rsa-4096-${ts}`,
    type: `rsa-4096`,
    supportsSigning: true,
    supportsEncryption: true,
  },
];

let generateTransitKey = async function(key, now) {
  let name = key.name(now);
  await click('[data-test-secret-create]');
  await fillIn('[data-test-transit-key-name]', name);
  await fillIn('[data-test-transit-key-type]', key.type);
  if (key.exportable) {
    await click('[data-test-transit-key-exportable]');
  }
  if (key.derived) {
    await click('[data-test-transit-key-derived]');
  }
  if (key.convergent) {
    await click('[data-test-transit-key-convergent-encryption]');
  }
  await click('[data-test-transit-key-create]');
  await settled();

  // link back to the list
  await click('[data-test-secret-root-link]');
  await settled();
  return name;
};

const testConvergentEncryption = async function(assert, keyName) {
  const tests = [
    // raw bytes for plaintext and context
    {
      plaintext: 'NaXud2QW7KjyK6Me9ggh+zmnCeBGdG93LQED49PtoOI=',
      context: 'nqR8LiVgNh/lwO2rArJJE9F9DMhh0lKo4JX9DAAkCDw=',
      encodePlaintext: false,
      encodeContext: false,
      decodeAfterDecrypt: false,
      assertAfterEncrypt: key => {
        assert.ok(
          /vault:/.test(find('[data-test-transit-input="ciphertext"]').value),
          `${key}: ciphertext shows a vault-prefixed ciphertext`
        );
      },
      assertBeforeDecrypt: key => {
        assert
          .dom('[data-test-transit-input="context"]')
          .hasValue(
            'nqR8LiVgNh/lwO2rArJJE9F9DMhh0lKo4JX9DAAkCDw=',
            `${key}: the ui shows the base64-encoded context`
          );
      },

      assertAfterDecrypt: key => {
        assert
          .dom('[data-test-transit-input="plaintext"]')
          .hasValue(
            'NaXud2QW7KjyK6Me9ggh+zmnCeBGdG93LQED49PtoOI=',
            `${key}: the ui shows the base64-encoded plaintext`
          );
      },
    },
    // raw bytes for plaintext, string for context
    {
      plaintext: 'NaXud2QW7KjyK6Me9ggh+zmnCeBGdG93LQED49PtoOI=',
      context: encodeString('context'),
      encodePlaintext: false,
      encodeContext: false,
      decodeAfterDecrypt: false,
      assertAfterEncrypt: key => {
        assert.ok(
          /vault:/.test(find('[data-test-transit-input="ciphertext"]').value),
          `${key}: ciphertext shows a vault-prefixed ciphertext`
        );
      },
      assertBeforeDecrypt: key => {
        assert
          .dom('[data-test-transit-input="context"]')
          .hasValue(encodeString('context'), `${key}: the ui shows the input context`);
      },
      assertAfterDecrypt: key => {
        assert
          .dom('[data-test-transit-input="plaintext"]')
          .hasValue(
            'NaXud2QW7KjyK6Me9ggh+zmnCeBGdG93LQED49PtoOI=',
            `${key}: the ui shows the base64-encoded plaintext`
          );
      },
    },
    // base64 input
    {
      plaintext: encodeString('This is the secret'),
      context: encodeString('context'),
      encodePlaintext: false,
      encodeContext: false,
      decodeAfterDecrypt: true,
      assertAfterEncrypt: key => {
        assert.ok(
          /vault:/.test(find('[data-test-transit-input="ciphertext"]').value),
          `${key}: ciphertext shows a vault-prefixed ciphertext`
        );
      },
      assertBeforeDecrypt: key => {
        assert
          .dom('[data-test-transit-input="context"]')
          .hasValue(encodeString('context'), `${key}: the ui shows the input context`);
      },
      assertAfterDecrypt: key => {
        assert
          .dom('[data-test-transit-input="plaintext"]')
          .hasValue('This is the secret', `${key}: the ui decodes plaintext`);
      },
    },

    // string input
    {
      plaintext: 'There are many secrets 🤐',
      context: 'secret 2',
      encodePlaintext: true,
      encodeContext: true,
      decodeAfterDecrypt: true,
      assertAfterEncrypt: key => {
        assert.ok(find('[data-test-transit-input="ciphertext"]'), `${key}: ciphertext box shows`);
        assert.ok(
          /vault:/.test(find('[data-test-transit-input="ciphertext"]').value),
          `${key}: ciphertext shows a vault-prefixed ciphertext`
        );
      },
      assertBeforeDecrypt: key => {
        assert
          .dom('[data-test-transit-input="context"]')
          .hasValue(encodeString('secret 2'), `${key}: the ui shows the encoded context`);
      },
      assertAfterDecrypt: key => {
        assert.ok(find('[data-test-transit-input="plaintext"]'), `${key}: plaintext box shows`);
        assert
          .dom('[data-test-transit-input="plaintext"]')
          .hasValue('There are many secrets 🤐', `${key}: the ui decodes plaintext`);
      },
    },
  ];

  for (let testCase of tests) {
    await click('[data-test-transit-action-link="encrypt"]');
    await fillIn('[data-test-transit-input="plaintext"]', testCase.plaintext);
    await fillIn('[data-test-transit-input="context"]', testCase.context);
    if (testCase.encodePlaintext) {
      await click('[data-test-transit-b64-toggle="plaintext"]');
    }
    if (testCase.encodeContext) {
      await click('[data-test-transit-b64-toggle="context"]');
    }
    await click('[data-test-button-encrypt]');
    await settled();
    if (testCase.assertAfterEncrypt) {
      testCase.assertAfterEncrypt(keyName);
    }
    await click('[data-test-transit-action-link="decrypt"]');
    await settled();
    if (testCase.assertBeforeDecrypt) {
      testCase.assertBeforeDecrypt(keyName);
    }
    await click('[data-test-button-decrypt]');
    await settled();

    if (testCase.assertAfterDecrypt) {
      if (testCase.decodeAfterDecrypt) {
        await click('[data-test-transit-b64-toggle="plaintext"]');
        testCase.assertAfterDecrypt(keyName);
      } else {
        testCase.assertAfterDecrypt(keyName);
      }
    }
  }
};
module('Acceptance | transit', function(hooks) {
  setupApplicationTest(hooks);
  let path;
  let now;

  hooks.beforeEach(async function() {
    await authPage.login();
    now = new Date().getTime();
    path = `transit-${now}`;

    await enablePage.enable('transit', path);
    await settled();
  });

  for (let key of keyTypes) {
    test(`transit backend: ${key.type}`, async function(assert) {
      let name = await generateTransitKey(key, now);
      await visit(`vault/secrets/${path}/show/${name}`);
      await settled();
      await click('[data-test-transit-link="versions"]');
      // wait for capabilities
      await settled();
      assert.dom('[data-test-transit-key-version-row]').exists({ count: 1 }, `${name}: only one key version`);
      await click('[data-test-confirm-action-trigger');
      await click('[data-test-confirm-button]');
      // wait for rotate call
      await settled();
      assert
        .dom('[data-test-transit-key-version-row]')
        .exists({ count: 2 }, `${name}: two key versions after rotate`);
      await click('[data-test-transit-key-actions-link]');
      await settled();
      assert.ok(
        currentURL().startsWith(`/vault/secrets/${path}/actions/${name}`),
        `${name}: navigates to tranist actions`
      );
      assert.ok(
        find('[data-test-transit-key-version-select]'),
        `${name}: the rotated key allows you to select versions`
      );
      if (key.exportable) {
        assert.ok(
          find('[data-test-transit-action-link="export"]'),
          `${name}: exportable key has a link to export action`
        );
      } else {
        assert
          .dom('[data-test-transit-action-link="export"]')
          .doesNotExist(`${name}: non-exportable key does not link to export action`);
      }
      if (key.convergent && key.supportsEncryption) {
        await testConvergentEncryption(assert, name);
      }
      await settled();
    });
  }
});
