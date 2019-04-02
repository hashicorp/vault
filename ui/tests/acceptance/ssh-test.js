import { click, fillIn, findAll, currentURL, find, settled } from '@ember/test-helpers';
import { module, test } from 'qunit';
import { setupApplicationTest } from 'ember-qunit';
import authPage from 'vault/tests/pages/auth';
import enablePage from 'vault/tests/pages/settings/mount-secret-backend';

module('Acceptance | ssh secret backend', function(hooks) {
  setupApplicationTest(hooks);

  hooks.beforeEach(function() {
    return authPage.login();
  });

  const PUB_KEY = `ssh-rsa AAAAB3NzaC1yc2EAAAADAQABAAABAQCn9p5dHNr4aU4R2W7ln+efzO5N2Cdv/SXk6zbCcvhWcblWMjkXf802B0PbKvf6cJIzM/Xalb3qz1cK+UUjCSEAQWefk6YmfzbOikfc5EHaSKUqDdE+HlsGPvl42rjCr28qYfuYh031YfwEQGEAIEypo7OyAj+38NLbHAQxDxuaReee1YCOV5rqWGtEgl2VtP5kG+QEBza4ZfeglS85f/GGTvZC4Jq1GX+wgmFxIPnd6/mUXa4ecoR0QMfOAzzvPm4ajcNCQORfHLQKAcmiBYMiyQJoU+fYpi9CJGT1jWTmR99yBkrSg6yitI2qqXyrpwAbhNGrM0Fw0WpWxh66N9Xp meirish@Macintosh-3.local`;

  const ROLES = [
    {
      type: 'ca',
      name: 'carole',
      async fillInCreate() {
        await click('[data-test-input="allowUserCertificates"]');
      },
      async fillInGenerate() {
        await fillIn('[data-test-input="publicKey"]', PUB_KEY);
      },
      assertAfterGenerate(assert, sshPath) {
        assert.equal(currentURL(), `/vault/secrets/${sshPath}/sign/${this.name}`, 'ca sign url is correct');
        assert.dom('[data-test-row-label="Signed key"]').exists({ count: 1 }, 'renders the signed key');
        assert
          .dom('[data-test-row-value="Signed key"]')
          .exists({ count: 1 }, "renders the signed key's value");
        assert.dom('[data-test-row-label="Serial number"]').exists({ count: 1 }, 'renders the serial');
        assert.dom('[data-test-row-value="Serial number"]').exists({ count: 1 }, 'renders the serial value');
      },
    },
    {
      type: 'otp',
      name: 'otprole',
      async fillInCreate() {
        await fillIn('[data-test-input="defaultUser"]', 'admin');
        await click('[data-test-toggle-more]');
        await fillIn('[data-test-input="cidrList"]', '1.2.3.4/32');
      },
      async fillInGenerate() {
        await fillIn('[data-test-input="username"]', 'admin');
        await fillIn('[data-test-input="ip"]', '1.2.3.4');
      },
      assertAfterGenerate(assert, sshPath) {
        assert.equal(
          currentURL(),
          `/vault/secrets/${sshPath}/credentials/${this.name}`,
          'otp credential url is correct'
        );
        assert.dom('[data-test-row-label="Key"]').exists({ count: 1 }, 'renders the key');
        assert.dom('[data-test-masked-input]').exists({ count: 1 }, 'renders mask for key value');
        assert.dom('[data-test-row-label="Port"]').exists({ count: 1 }, 'renders the port');
        assert.dom('[data-test-row-value="Port"]').exists({ count: 1 }, "renders the port's value");
      },
    },
  ];
  test('ssh backend', async function(assert) {
    const now = new Date().getTime();
    const sshPath = `ssh-${now}`;

    await enablePage.enable('ssh', sshPath);
    await click('[data-test-secret-backend-configure]');
    assert.equal(currentURL(), `/vault/settings/secrets/configure/${sshPath}`);
    assert.ok(findAll('[data-test-ssh-configure-form]').length, 'renders the empty configuration form');

    // default has generate CA checked so we just submit the form
    await click('[data-test-ssh-input="configure-submit"]');

    assert.ok(findAll('[data-test-ssh-input="public-key"]').length, 'a public key is fetched');
    await click('[data-test-backend-view-link]');

    assert.equal(currentURL(), `/vault/secrets/${sshPath}/list`, `redirects to ssh index`);

    for (let role of ROLES) {
      // create a role
      await click('[ data-test-secret-create]');
      assert.ok(
        find('[data-test-secret-header]').textContent.includes('SSH role'),
        `${role.type}: renders the create page`
      );

      await fillIn('[data-test-input="name"]', role.name);
      await fillIn('[data-test-input="keyType"]', role.type);
      await role.fillInCreate();

      // save the role
      await click('[data-test-role-ssh-create]');
      assert.equal(
        currentURL(),
        `/vault/secrets/${sshPath}/show/${role.name}`,
        `${role.type}: navigates to the show page on creation`
      );

      // sign a key with this role
      await click('[data-test-backend-credentials]');
      await role.fillInGenerate();

      // generate creds
      await click('[data-test-secret-generate]');
      role.assertAfterGenerate(assert, sshPath);

      // click the "Back" button
      await click('[data-test-secret-generate-back]');
      assert.ok(
        findAll('[data-test-secret-generate-form]').length,
        `${role.type}: back takes you back to the form`
      );

      await click('[data-test-secret-generate-cancel]');
      assert.equal(
        currentURL(),
        `/vault/secrets/${sshPath}/list`,
        `${role.type}: cancel takes you to ssh index`
      );
      assert.ok(
        findAll(`[data-test-secret-link="${role.name}"]`).length,
        `${role.type}: role shows in the list`
      );

      //and delete
      await click(`[data-test-secret-link="${role.name}"] [data-test-popup-menu-trigger]`);
      await click(`[data-test-ssh-role-delete="${role.name}"] button`);
      await click(`[data-test-confirm-button]`);

      await settled();
      assert
        .dom(`[data-test-secret-link="${role.name}"]`)
        .doesNotExist(`${role.type}: role is no longer in the list`);
    }
  });
});
