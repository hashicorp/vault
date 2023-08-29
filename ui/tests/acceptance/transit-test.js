/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: BUSL-1.1
 */

import { click, fillIn, find, currentURL, settled, visit, waitUntil, findAll } from '@ember/test-helpers';
import { module, test } from 'qunit';
import { setupApplicationTest } from 'ember-qunit';
import { v4 as uuidv4 } from 'uuid';

import { encodeString } from 'vault/utils/b64';
import authPage from 'vault/tests/pages/auth';
import { deleteEngineCmd, mountEngineCmd, runCmd } from 'vault/tests/helpers/commands';

const SELECTORS = {
  secretLink: '[data-test-secret-link]',
  popupMenu: '[data-test-popup-menu-trigger]',
  versionsTab: '[data-test-transit-link="versions"]',
  actionsTab: '[data-test-transit-key-actions-link]',
  card: (action) => `[data-test-transit-card="${action}"]`,
  rotate: {
    trigger: '[data-test-confirm-action-trigger]',
    confirm: '[data-test-confirm-button]',
  },
};

// convergent
const groupOne = [
  {
    name: (uid = 'name') => `aes-1A-convergent-${uid}`,
    type: 'aes128-gcm96',
    convergent: true,
    supportsEncryption: true,
  },
  {
    name: (uid = 'name') => `aes-1B-convergent-${uid}`,
    type: 'aes256-gcm96',
    convergent: true,
    supportsEncryption: true,
  },
  {
    name: (uid = 'name') => `chacha-1C-convergent-${uid}`,
    type: 'chacha20-poly1305',
    convergent: true,
    supportsEncryption: true,
    autoRotate: true,
  },
];

// exportable, supports encryption
const groupTwo = [
  {
    name: (uid = 'name') => `aes-2A-${uid}`,
    type: 'aes128-gcm96',
    exportable: true,
    supportsEncryption: true,
  },
  {
    name: (uid = 'name') => `aes-2B-${uid}`,
    type: 'aes256-gcm96',
    exportable: true,
    supportsEncryption: true,
  },
  {
    name: (uid = 'name') => `chacha-3B-${uid}`,
    type: 'chacha20-poly1305',
    exportable: true,
    supportsEncryption: true,
  },
];

// exportable, supports signing
const groupThree = [
  {
    name: (uid = 'name') => `ecdsa-3A-${uid}`,
    type: 'ecdsa-p256',
    exportable: true,
    supportsSigning: true,
  },
  {
    name: (uid = 'name') => `ecdsa-3B-${uid}`,
    type: 'ecdsa-p384',
    exportable: true,
    supportsSigning: true,
  },
  {
    name: (uid = 'name') => `ecdsa-3C-${uid}`,
    type: 'ecdsa-p521',
    exportable: true,
    supportsSigning: true,
  },
];
// the rest
const groupFour = [
  {
    name: (uid = 'name') => `ed25519-4A-${uid}`,
    type: 'ed25519',
    derived: true,
    supportsSigning: true,
  },
  {
    name: (uid = 'name') => `rsa-2048-4B-${uid}`,
    type: `rsa-2048`,
    supportsSigning: true,
    supportsEncryption: true,
  },
  {
    name: (uid = 'name') => `rsa-3072-4C-${uid}`,
    type: `rsa-3072`,
    supportsSigning: true,
    supportsEncryption: true,
  },
  {
    name: (uid = 'name') => `rsa-4096-4D-${uid}`,
    type: `rsa-4096`,
    supportsSigning: true,
    supportsEncryption: true,
    autoRotate: true,
  },
];

const generateTransitKey = async function (key, uid) {
  const name = key.name(uid);
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
  if (key.autoRotate) {
    await click('[data-test-toggle-label="Auto-rotation period"]');
  }
  await click('[data-test-transit-key-create]');
  await settled(); // eslint-disable-line
  // link back to the list
  await click('[data-test-secret-root-link]');

  return name;
};

const testConvergentEncryption = async function (assert, keyName) {
  const tests = [
    // raw bytes for plaintext and context
    {
      plaintext: 'NaXud2QW7KjyK6Me9ggh+zmnCeBGdG93LQED49PtoOI=',
      context: 'nqR8LiVgNh/lwO2rArJJE9F9DMhh0lKo4JX9DAAkCDw=',
      encodePlaintext: false,
      encodeContext: false,
      assertAfterEncrypt: (key) => {
        assert.dom('.modal.is-active').exists(`${key}: Modal opens after encrypt`);
        assert.ok(
          /vault:/.test(find('[data-test-encrypted-value="ciphertext"]').innerText),
          `${key}: ciphertext shows a vault-prefixed ciphertext`
        );
      },
      assertBeforeDecrypt: (key) => {
        assert.dom('.modal.is-active').doesNotExist(`${key}: Modal not open before decrypt`);
        assert
          .dom('[data-test-transit-input="context"]')
          .hasValue(
            'nqR8LiVgNh/lwO2rArJJE9F9DMhh0lKo4JX9DAAkCDw=',
            `${key}: the ui shows the base64-encoded context`
          );
      },

      assertAfterDecrypt: (key) => {
        assert.dom('.modal.is-active').exists(`${key}: Modal opens after decrypt`);
        assert.strictEqual(
          find('[data-test-encrypted-value="plaintext"]').innerText,
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
      assertAfterEncrypt: (key) => {
        assert.dom('.modal.is-active').exists(`${key}: Modal opens after encrypt`);
        assert.ok(
          /vault:/.test(find('[data-test-encrypted-value="ciphertext"]').innerText),
          `${key}: ciphertext shows a vault-prefixed ciphertext`
        );
      },
      assertBeforeDecrypt: (key) => {
        assert.dom('.modal.is-active').doesNotExist(`${key}: Modal not open before decrypt`);
        assert
          .dom('[data-test-transit-input="context"]')
          .hasValue(encodeString('context'), `${key}: the ui shows the input context`);
      },
      assertAfterDecrypt: (key) => {
        assert.dom('.modal.is-active').exists(`${key}: Modal opens after decrypt`);
        assert.strictEqual(
          find('[data-test-encrypted-value="plaintext"]').innerText,
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
      assertAfterEncrypt: (key) => {
        assert.dom('.modal.is-active').exists(`${key}: Modal opens after encrypt`);
        assert.ok(
          /vault:/.test(find('[data-test-encrypted-value="ciphertext"]').innerText),
          `${key}: ciphertext shows a vault-prefixed ciphertext`
        );
      },
      assertBeforeDecrypt: (key) => {
        assert.dom('.modal.is-active').doesNotExist(`${key}: Modal not open before decrypt`);
        assert
          .dom('[data-test-transit-input="context"]')
          .hasValue(encodeString('context'), `${key}: the ui shows the input context`);
      },
      assertAfterDecrypt: (key) => {
        assert.dom('.modal.is-active').exists(`${key}: Modal opens after decrypt`);
        assert.strictEqual(
          find('[data-test-encrypted-value="plaintext"]').innerText,
          encodeString('This is the secret'),
          `${key}: the ui decodes plaintext`
        );
      },
    },

    // string input
    {
      plaintext: 'There are many secrets 🤐',
      context: 'secret 2',
      encodePlaintext: true,
      encodeContext: true,
      assertAfterEncrypt: (key) => {
        assert.dom('.modal.is-active').exists(`${key}: Modal opens after encrypt`);
        assert.ok(
          /vault:/.test(find('[data-test-encrypted-value="ciphertext"]').innerText),
          `${key}: ciphertext shows a vault-prefixed ciphertext`
        );
      },
      assertBeforeDecrypt: (key) => {
        assert.dom('.modal.is-active').doesNotExist(`${key}: Modal not open before decrypt`);
        assert
          .dom('[data-test-transit-input="context"]')
          .hasValue(encodeString('secret 2'), `${key}: the ui shows the encoded context`);
      },
      assertAfterDecrypt: (key) => {
        assert.dom('.modal.is-active').exists(`${key}: Modal opens after decrypt`);
        assert.strictEqual(
          find('[data-test-encrypted-value="plaintext"]').innerText,
          encodeString('There are many secrets 🤐'),
          `${key}: the ui decodes plaintext`
        );
      },
    },
  ];

  for (const testCase of tests) {
    await click('[data-test-transit-action-link="encrypt"]');

    find('#plaintext-control .CodeMirror').CodeMirror.setValue(testCase.plaintext);
    await fillIn('[data-test-transit-input="context"]', testCase.context);

    if (!testCase.encodePlaintext) {
      // If value is already encoded, check the box
      await click('input[data-test-transit-input="encodedBase64"]');
    }
    if (testCase.encodeContext) {
      await click('[data-test-transit-b64-toggle="context"]');
    }
    assert.dom('.modal.is-active').doesNotExist(`${name}: is not open before encrypt`);
    await click('[data-test-button-encrypt]');

    if (testCase.assertAfterEncrypt) {
      await settled();
      testCase.assertAfterEncrypt(keyName);
    }
    // store ciphertext for decryption step
    const copiedCiphertext = find('[data-test-encrypted-value="ciphertext"]').innerText;
    await click('.modal.is-active [data-test-modal-background]');

    assert.dom('.modal.is-active').doesNotExist(`${name}: Modal closes after background clicked`);
    await click('[data-test-transit-action-link="decrypt"]');

    if (testCase.assertBeforeDecrypt) {
      await settled();
      testCase.assertBeforeDecrypt(keyName);
    }
    find('#ciphertext-control .CodeMirror').CodeMirror.setValue(copiedCiphertext);
    await click('[data-test-button-decrypt]');

    if (testCase.assertAfterDecrypt) {
      await settled();
      testCase.assertAfterDecrypt(keyName);
    }

    await click('.modal.is-active [data-test-modal-background]');

    assert.dom('.modal.is-active').doesNotExist(`${name}: Modal closes after background clicked`);
  }
};
module('Acceptance | transit', function (hooks) {
  setupApplicationTest(hooks);

  hooks.beforeEach(async function () {
    const uid = uuidv4();
    await authPage.login();
    await settled();
    this.uid = uid;
    this.path = `transit-${uid}`;

    await runCmd(mountEngineCmd('transit', this.path));
    // Start test on backend main page
    return visit(`/vault/secrets/${this.path}/list`);
  });

  hooks.afterEach(async function () {
    await authPage.login();
    await runCmd(deleteEngineCmd(this.mountPath));
  });

  test(`transit backend: list menu`, async function (assert) {
    assert.expect(3);
    const name = await generateTransitKey(groupTwo[0], this.uid);
    await click(SELECTORS.popupMenu);
    const actions = findAll('.ember-basic-dropdown-content li');
    assert.strictEqual(actions.length, 2, 'shows 2 items in popup menu');

    await click(SELECTORS.secretLink);
    assert.strictEqual(
      currentURL(),
      `/vault/secrets/${this.path}/show/${name}?tab=actions`,
      'navigates to key actions tab'
    );
    await click(SELECTORS.actionsTab);
    assert.strictEqual(currentURL(), `/vault/secrets/${this.path}/show/${name}?tab=actions`),
      'navigates back to transit actions';
  });

  // convergent keys
  for (const key of groupOne) {
    test(`transit backend: group 1 ${key.name()}`, async function (assert) {
      assert.expect(42);
      const name = await generateTransitKey(key, this.uid);
      await visit(`vault/secrets/${this.path}/show/${name}`);

      const expectedRotateValue = key.autoRotate ? '30 days' : 'Key will not be automatically rotated';
      assert
        .dom('[data-test-row-value="Auto-rotation period"]')
        .hasText(expectedRotateValue, 'Has expected auto rotate value');

      await click(SELECTORS.versionsTab);
      assert.dom('[data-test-transit-key-version-row]').exists({ count: 1 }, `${name}: only one key version`);
      await waitUntil(() => find(SELECTORS.rotate.trigger));
      await click(SELECTORS.rotate.trigger);
      await click(SELECTORS.rotate.confirm);
      // wait for rotate call
      await waitUntil(() => findAll('[data-test-transit-key-version-row]').length >= 2);
      assert
        .dom('[data-test-transit-key-version-row]')
        .exists({ count: 2 }, `${name}: two key versions after rotate`);

      // navigate back to actions tab
      await visit(`/vault/secrets/${this.path}/show/${name}?tab=actions`);

      const keyAction = 'encrypt';
      await waitUntil(() => find(SELECTORS.card(keyAction)));
      assert.dom(SELECTORS.card(keyAction)).exists(`renders ${keyAction} card for ${key.name()}`);
      await click(SELECTORS.card(keyAction));
      assert
        .dom('[data-test-transit-key-version-select]')
        .exists(`${name}: the rotated key allows you to select versions`);
      assert
        .dom('[data-test-transit-action-link="export"]')
        .doesNotExist(`${name}: non-exportable key does not link to export action`);
      await testConvergentEncryption(assert, name);
    });
  }
  // exportable, supports encryption
  for (const key of groupTwo) {
    test(`transit backend: group 2 ${key.name()}`, async function (assert) {
      assert.expect(6);
      const name = await generateTransitKey(key, this.uid);
      await visit(`vault/secrets/${this.path}/show/${name}`);

      assert
        .dom('[data-test-row-value="Auto-rotation period"]')
        .hasText('Key will not be automatically rotated', 'key will not auto rotate');

      await click(SELECTORS.versionsTab);
      assert.dom('[data-test-transit-key-version-row]').exists({ count: 1 }, `${name}: only one key version`);
      await waitUntil(() => find(SELECTORS.rotate.trigger));
      await click(SELECTORS.rotate.trigger);
      await click(SELECTORS.rotate.confirm);

      // wait for rotate call
      await waitUntil(() => findAll('[data-test-transit-key-version-row]').length >= 2);
      assert
        .dom('[data-test-transit-key-version-row]')
        .exists({ count: 2 }, `${name}: two key versions after rotate`);

      // navigate back to actions tab
      await visit(`/vault/secrets/${this.path}/show/${name}?tab=actions`);
      const keyAction = 'encrypt';
      await waitUntil(() => find(SELECTORS.card(keyAction)));
      assert.dom(SELECTORS.card(keyAction)).exists(`renders ${keyAction} card for ${key.name()}`);
      await click(SELECTORS.card(keyAction));

      assert
        .dom('[data-test-transit-key-version-select]')
        .exists(`${name}: the rotated key allows you to select versions`);
      assert
        .dom('[data-test-transit-action-link="export"]')
        .exists(`${name}: exportable key has a link to export action`);
    });
  }
  // exportable, supports signing
  for (const key of groupThree) {
    test(`transit backend: group 2 ${key.name()}`, async function (assert) {
      assert.expect(6);
      const name = await generateTransitKey(key, this.uid);
      await visit(`vault/secrets/${this.path}/show/${name}`);

      assert
        .dom('[data-test-row-value="Auto-rotation period"]')
        .hasText('Key will not be automatically rotated', 'key will not auto rotate');

      await click(SELECTORS.versionsTab);
      assert.dom('[data-test-transit-key-version-row]').exists({ count: 1 }, `${name}: only one key version`);
      await waitUntil(() => find(SELECTORS.rotate.trigger));
      await click(SELECTORS.rotate.trigger);
      await click(SELECTORS.rotate.confirm);

      // wait for rotate call
      await waitUntil(() => findAll('[data-test-transit-key-version-row]').length >= 2);
      assert
        .dom('[data-test-transit-key-version-row]')
        .exists({ count: 2 }, `${name}: two key versions after rotate`);

      // navigate back to actions tab
      await visit(`/vault/secrets/${this.path}/show/${name}?tab=actions`);
      const keyAction = 'sign';
      await waitUntil(() => find(SELECTORS.card(keyAction)));
      assert.dom(SELECTORS.card(keyAction)).exists(`renders ${keyAction} card for ${key.name()}`);
      await click(SELECTORS.card(keyAction));

      assert
        .dom('[data-test-transit-key-version-select]')
        .exists(`${name}: the rotated key allows you to select versions`);
      assert
        .dom('[data-test-transit-action-link="export"]')
        .exists(`${name}: exportable key has a link to export action`);
    });
  }

  for (const key of groupFour) {
    test(`transit backend: group 3 ${key.name()}`, async function (assert) {
      assert.expect(6);
      const name = await generateTransitKey(key, this.uid);
      await visit(`vault/secrets/${this.path}/show/${name}`);

      const expectedRotateValue = key.autoRotate ? '30 days' : 'Key will not be automatically rotated';
      assert
        .dom('[data-test-row-value="Auto-rotation period"]')
        .hasText(expectedRotateValue, 'Has expected auto rotate value');

      await click(SELECTORS.versionsTab);

      assert.dom('[data-test-transit-key-version-row]').exists({ count: 1 }, `${name}: only one key version`);
      await waitUntil(() => find(SELECTORS.rotate.trigger));
      await click(SELECTORS.rotate.trigger);
      await click(SELECTORS.rotate.confirm);

      // wait for rotate call
      await waitUntil(() => findAll('[data-test-transit-key-version-row]').length >= 2);
      assert
        .dom('[data-test-transit-key-version-row]')
        .exists({ count: 2 }, `${name}: two key versions after rotate`);

      // navigate back to actions tab
      await visit(`/vault/secrets/${this.path}/show/${name}?tab=actions`);

      const keyAction = key.supportsEncryption ? 'encrypt' : 'sign';
      await waitUntil(() => find(SELECTORS.card(keyAction)));
      assert.dom(SELECTORS.card(keyAction)).exists(`renders ${keyAction} card for ${key.name()}`);
      await click(SELECTORS.card(keyAction));

      assert
        .dom('[data-test-transit-key-version-select]')
        .exists(`${name}: the rotated key allows you to select versions`);

      assert
        .dom('[data-test-transit-action-link="export"]')
        .doesNotExist(`${name}: non-exportable key does not link to export action`);
    });
  }
});
