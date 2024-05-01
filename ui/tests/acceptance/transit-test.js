/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: BUSL-1.1
 */

import { click, fillIn, find, currentURL, settled, visit, findAll } from '@ember/test-helpers';
import { module, test } from 'qunit';
import { setupApplicationTest } from 'ember-qunit';
import { v4 as uuidv4 } from 'uuid';

import { encodeString } from 'vault/utils/b64';
import authPage from 'vault/tests/pages/auth';
import { deleteEngineCmd, mountEngineCmd, runCmd } from 'vault/tests/helpers/commands';
import codemirror from 'vault/tests/helpers/codemirror';
import { GENERAL } from '../helpers/general-selectors';

const SELECTORS = {
  secretLink: '[data-test-secret-link]',
  popupMenu: '[data-test-popup-menu-trigger]',
  versionsTab: '[data-test-transit-link="versions"]',
  actionsTab: '[data-test-transit-key-actions-link]',
  card: (action) => `[data-test-transit-card="${action}"]`,
  infoRow: (label) => `[data-test-value-div="${label}"]`,
  form: (item) => `[data-test-transit-key="${item}"]`,
  versionRow: (version) => `[data-test-transit-version="${version}"]`,
  rotate: {
    trigger: '[data-test-transit-key-rotate]',
    confirm: '[data-test-confirm-button]',
  },
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
        assert.dom('[data-test-encrypt-modal]').exists(`${key}: Modal opens after encrypt`);
        assert.ok(
          /vault:/.test(find('[data-test-encrypted-value="ciphertext"]').innerText),
          `${key}: ciphertext shows a vault-prefixed ciphertext`
        );
      },
      assertBeforeDecrypt: (key) => {
        assert.dom('[data-test-decrypt-modal]').doesNotExist(`${key}: Modal not open before decrypt`);
        assert
          .dom('[data-test-transit-input="context"]')
          .hasValue(
            'nqR8LiVgNh/lwO2rArJJE9F9DMhh0lKo4JX9DAAkCDw=',
            `${key}: the ui shows the base64-encoded context`
          );
      },
      assertAfterDecrypt: (key) => {
        assert.dom('[data-test-decrypt-modal]').exists(`${key}: Modal opens after decrypt`);
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
        assert.dom('[data-test-encrypt-modal]').exists(`${key}: Modal opens after encrypt`);
        assert.ok(
          /vault:/.test(find('[data-test-encrypted-value="ciphertext"]').innerText),
          `${key}: ciphertext shows a vault-prefixed ciphertext`
        );
      },
      assertBeforeDecrypt: (key) => {
        assert.dom('[data-test-decrypt-modal]').doesNotExist(`${key}: Modal not open before decrypt`);
        assert
          .dom('[data-test-transit-input="context"]')
          .hasValue(encodeString('context'), `${key}: the ui shows the input context`);
      },
      assertAfterDecrypt: (key) => {
        assert.dom('[data-test-decrypt-modal]').exists(`${key}: Modal opens after decrypt`);
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
        assert.dom('[data-test-encrypt-modal]').exists(`${key}: Modal opens after encrypt`);
        assert.ok(
          /vault:/.test(find('[data-test-encrypted-value="ciphertext"]').innerText),
          `${key}: ciphertext shows a vault-prefixed ciphertext`
        );
      },
      assertBeforeDecrypt: (key) => {
        assert.dom('[data-test-decrypt-modal]').doesNotExist(`${key}: Modal not open before decrypt`);
        assert
          .dom('[data-test-transit-input="context"]')
          .hasValue(encodeString('context'), `${key}: the ui shows the input context`);
      },
      assertAfterDecrypt: (key) => {
        assert.dom('[data-test-decrypt-modal]').exists(`${key}: Modal opens after decrypt`);
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
        assert.dom('[data-test-encrypt-modal]').exists(`${key}: Modal opens after encrypt`);
        assert.ok(
          /vault:/.test(find('[data-test-encrypted-value="ciphertext"]').innerText),
          `${key}: ciphertext shows a vault-prefixed ciphertext`
        );
      },
      assertBeforeDecrypt: (key) => {
        assert.dom('[data-test-decrypt-modal]').doesNotExist(`${key}: Modal not open before decrypt`);
        assert
          .dom('[data-test-transit-input="context"]')
          .hasValue(encodeString('secret 2'), `${key}: the ui shows the encoded context`);
      },
      assertAfterDecrypt: (key) => {
        assert.dom('[data-test-decrypt-modal]').exists(`${key}: Modal opens after decrypt`);
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

    codemirror('#plaintext-control').setValue(testCase.plaintext);
    await fillIn('[data-test-transit-input="context"]', testCase.context);

    if (!testCase.encodePlaintext) {
      // If value is already encoded, check the box
      await click('input[data-test-transit-input="encodedBase64"]');
    }
    if (testCase.encodeContext) {
      await click('[data-test-transit-b64-toggle="context"]');
    }
    assert.dom('[data-test-encrypt-modal]').doesNotExist(`${keyName}: is not open before encrypt`);
    await click('[data-test-button-encrypt]');

    if (testCase.assertAfterEncrypt) {
      await settled();
      testCase.assertAfterEncrypt(keyName);
    }
    // store ciphertext for decryption step
    const copiedCiphertext = find('[data-test-encrypted-value="ciphertext"]').innerText;
    await click('dialog button');

    assert.dom('dialog.hds-modal').doesNotExist(`${keyName}: Modal closes after background clicked`);
    await click('[data-test-transit-action-link="decrypt"]');

    if (testCase.assertBeforeDecrypt) {
      await settled();
      testCase.assertBeforeDecrypt(keyName);
    }

    codemirror('#ciphertext-control').setValue(copiedCiphertext);
    await click('[data-test-button-decrypt]');

    if (testCase.assertAfterDecrypt) {
      await settled();
      testCase.assertAfterDecrypt(keyName);
    }

    await click('dialog button');

    assert.dom('dialog.hds-modal').doesNotExist(`${keyName}: Modal closes after background clicked`);
  }
};

module('Acceptance | transit (flaky)', function (hooks) {
  setupApplicationTest(hooks);

  hooks.beforeEach(async function () {
    const uid = uuidv4();
    await authPage.login();
    this.uid = uid;
    this.path = `transit-${uid}`;

    this.generateTransitKey = async function (key) {
      const name = key.name(uid);
      const config = [];

      if (key.exportable) config.push('exportable=true');
      if (key.derived) config.push('derived=true');
      if (key.convergent) config.push('convergent_encryption=true');
      if (key.autoRotate) config.push('auto_rotate_period=720h');

      await runCmd([`vault write ${this.path}/keys/${name} type=${key.type} ${config.join(' ')} -f`]);
      return name;
    };
    await runCmd(mountEngineCmd('transit', this.path));
    // Start test on backend main page
    return visit(`/vault/secrets/${this.path}/list`);
  });

  hooks.afterEach(async function () {
    await authPage.login();
    await runCmd(deleteEngineCmd(this.mountPath));
  });

  test('it generates a key', async function (assert) {
    assert.expect(8);
    const type = 'chacha20-poly1305';
    const name = `test-generate-${this.uid}`;
    await click('[data-test-secret-create]');

    await fillIn(SELECTORS.form('name'), name);
    await fillIn(SELECTORS.form('type'), type);
    await click(SELECTORS.form('exportable'));
    await click(SELECTORS.form('derived'));
    await click(SELECTORS.form('convergent-encryption'));
    await click('[data-test-toggle-label="Auto-rotation period"]');
    await click(SELECTORS.form('create'));

    assert.strictEqual(
      currentURL(),
      `/vault/secrets/${this.path}/show/${name}?tab=details`,
      'it navigates to show page'
    );
    assert.dom(SELECTORS.infoRow('Auto-rotation period')).hasText('30 days');
    assert.dom(SELECTORS.infoRow('Deletion allowed')).hasText('false');
    assert.dom(SELECTORS.infoRow('Derived')).hasText('Yes');
    assert.dom(SELECTORS.infoRow('Convergent encryption')).hasText('Yes');
    await click(GENERAL.breadcrumbLink(this.path));
    await click(SELECTORS.popupMenu);
    const actions = findAll('.hds-dropdown__list li');
    assert.strictEqual(actions.length, 2, 'shows 2 items in popup menu');

    await click(SELECTORS.secretLink);
    assert.strictEqual(
      currentURL(),
      `/vault/secrets/${this.path}/show/${name}?tab=actions`,
      'navigates to key actions tab'
    );
    await click(SELECTORS.actionsTab);
    assert.strictEqual(
      currentURL(),
      `/vault/secrets/${this.path}/show/${name}?tab=actions`,
      'navigates back to transit actions'
    );
  });

  test('create form renders supported options for each key type', async function (assert) {
    assert.expect(30);
    await visit(`/vault/secrets/${this.path}/create`);
    const KEY_OPTIONS = [
      {
        type: 'ed25519',
        derived: true,
        exportable: true,
      },
      {
        type: 'rsa-2048',
        exportable: true,
      },
      {
        type: 'rsa-3072',
        exportable: true,
      },
      {
        type: 'rsa-4096',
        exportable: true,
        supportsEncryption: true,
      },
      {
        type: 'ecdsa-p256',
        exportable: true,
      },
      {
        type: 'ecdsa-p384',
        exportable: true,
      },
      {
        type: 'ecdsa-p521',
        exportable: true,
      },
      {
        type: 'aes128-gcm96',
        'convergent-encryption': true,
        derived: true,
        exportable: true,
      },
      {
        type: 'aes256-gcm96',
        'convergent-encryption': true,
        derived: true,
        exportable: true,
      },
      {
        type: 'chacha20-poly1305',
        'convergent-encryption': true,
        derived: true,
        exportable: true,
      },
    ];
    for (const key of KEY_OPTIONS) {
      const { type } = key;
      await fillIn(SELECTORS.form('type'), type);

      for (const checkbox of ['exportable', 'derived', 'convergent-encryption']) {
        const assertion = key[checkbox] ? 'exists' : 'doesNotExist';
        assert.dom(SELECTORS.form(checkbox))[assertion](`${type} ${checkbox} ${assertion}`);
      }
    }
  });

  test('it rotates, encrypts and decrypts key type chacha20-poly1305', async function (assert) {
    assert.expect(42);
    const keyData = {
      name: (uid) => `chacha-convergent-${uid}`,
      type: 'chacha20-poly1305',
      convergent: true,
      derived: true,
      supportsEncryption: true,
      autoRotate: true,
    };

    const name = await this.generateTransitKey(keyData);
    await visit(`vault/secrets/${this.path}/show/${name}`);
    assert
      .dom(SELECTORS.infoRow('Auto-rotation period'))
      .hasText('30 days', 'Has expected auto rotate value');

    await click(SELECTORS.versionsTab);
    assert.dom(SELECTORS.versionRow(1)).hasTextContaining('Version 1', `${name}: only one key version`);

    await click(SELECTORS.rotate.trigger);
    await click(SELECTORS.rotate.confirm);

    assert.dom(SELECTORS.versionRow(2)).exists('two key versions after rotate');

    // navigate back to actions tab
    await click(SELECTORS.actionsTab);

    assert.dom(SELECTORS.card('encrypt')).exists(`renders encrypt action card for ${name}`);
    await click(SELECTORS.card('encrypt'));
    assert
      .dom('[data-test-transit-key-version-select]')
      .exists(`${name}: the rotated key allows you to select versions`);
    assert
      .dom('[data-test-transit-action-link="export"]')
      .doesNotExist(`${name}: non-exportable key does not link to export action`);
    await testConvergentEncryption(assert, name);
  });

  const KEY_TYPE_COMBINATIONS = [
    {
      name: (uid) => `aes-${uid}`,
      type: 'aes128-gcm96',
      exportable: true,
      supportsEncryption: true,
    },
    {
      name: (uid) => `aes-convergent-${uid}`,
      type: 'aes128-gcm96',
      convergent: true,
      derived: true,
      supportsEncryption: true,
    },
    {
      name: (uid) => `aes-${uid}`,
      type: 'aes256-gcm96',
      exportable: true,
      supportsEncryption: true,
    },
    {
      name: (uid) => `aes-convergent-${uid}`,
      type: 'aes256-gcm96',
      convergent: true,
      derived: true,
      supportsEncryption: true,
    },
    {
      name: (uid) => `chacha-${uid}`,
      type: 'chacha20-poly1305',
      exportable: true,
      supportsEncryption: true,
    },
    {
      name: (uid) => `chacha-convergent-${uid}`,
      type: 'chacha20-poly1305',
      convergent: true,
      derived: true,
      supportsEncryption: true,
      autoRotate: true,
    },
    {
      name: (uid) => `ecdsa-${uid}`,
      type: 'ecdsa-p256',
      exportable: true,
      supportsSigning: true,
    },
    {
      name: (uid) => `ecdsa-${uid}`,
      type: 'ecdsa-p384',
      exportable: true,
      supportsSigning: true,
    },
    {
      name: (uid) => `ecdsa-${uid}`,
      type: 'ecdsa-p521',
      exportable: true,
      supportsSigning: true,
    },
    {
      name: (uid) => `ed25519-${uid}`,
      type: 'ed25519',
      derived: true,
      supportsSigning: true,
    },
    {
      name: (uid) => `rsa-2048-${uid}`,
      type: `rsa-2048`,
      supportsSigning: true,
      supportsEncryption: true,
    },
    {
      name: (uid) => `rsa-3072-${uid}`,
      type: `rsa-3072`,
      supportsSigning: true,
      supportsEncryption: true,
    },
    {
      name: (uid) => `rsa-4096-${uid}`,
      type: `rsa-4096`,
      supportsSigning: true,
      supportsEncryption: true,
      autoRotate: true,
    },
  ];

  for (const key of KEY_TYPE_COMBINATIONS) {
    test(`transit backend: ${key.type}`, async function (assert) {
      assert.expect(key.convergent ? 43 : 7);
      const name = await this.generateTransitKey(key);
      await visit(`vault/secrets/${this.path}/show/${name}`);

      const expectedRotateValue = key.autoRotate ? '30 days' : 'Key will not be automatically rotated';
      assert
        .dom('[data-test-row-value="Auto-rotation period"]')
        .hasText(expectedRotateValue, 'Has expected auto rotate value');

      await click(SELECTORS.versionsTab);
      // wait for capabilities

      assert.dom('[data-test-transit-version]').exists({ count: 1 }, `${name}: only one key version`);
      await click(SELECTORS.rotate.trigger);

      await click(SELECTORS.rotate.confirm);
      assert
        .dom('[data-test-transit-version]')
        .exists({ count: 2 }, `${name}: two key versions after rotate`);
      await click('[data-test-transit-key-actions-link]');

      assert.strictEqual(
        currentURL(),
        `/vault/secrets/${this.path}/show/${name}?tab=actions`,
        `${name}: navigates to transit actions`
      );

      const keyAction = key.supportsEncryption ? 'encrypt' : 'sign';

      assert
        .dom(`[data-test-transit-action-title=${keyAction}]`)
        .exists(`shows a card with title that links to the ${name} transit action`);

      await click(SELECTORS.card(keyAction));

      assert
        .dom('[data-test-transit-key-version-select]')
        .exists(`${name}: the rotated key allows you to select versions`);
      if (key.exportable) {
        assert
          .dom('[data-test-transit-action-link="export"]')
          .exists(`${name}: exportable key has a link to export action`);
      } else {
        assert
          .dom('[data-test-transit-action-link="export"]')
          .doesNotExist(`${name}: non-exportable key does not link to export action`);
      }
      if (key.convergent && key.supportsEncryption) {
        await testConvergentEncryption(assert, name);
        await settled();
      }
      await settled();
    });
  }
});
