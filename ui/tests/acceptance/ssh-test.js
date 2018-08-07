import { test } from 'qunit';
import moduleForAcceptance from 'vault/tests/helpers/module-for-acceptance';

moduleForAcceptance('Acceptance | ssh secret backend', {
  beforeEach() {
    return authLogin();
  },
  afterEach() {
    return authLogout();
  },
});

const PUB_KEY = `ssh-rsa AAAAB3NzaC1yc2EAAAADAQABAAABAQCn9p5dHNr4aU4R2W7ln+efzO5N2Cdv/SXk6zbCcvhWcblWMjkXf802B0PbKvf6cJIzM/Xalb3qz1cK+UUjCSEAQWefk6YmfzbOikfc5EHaSKUqDdE+HlsGPvl42rjCr28qYfuYh031YfwEQGEAIEypo7OyAj+38NLbHAQxDxuaReee1YCOV5rqWGtEgl2VtP5kG+QEBza4ZfeglS85f/GGTvZC4Jq1GX+wgmFxIPnd6/mUXa4ecoR0QMfOAzzvPm4ajcNCQORfHLQKAcmiBYMiyQJoU+fYpi9CJGT1jWTmR99yBkrSg6yitI2qqXyrpwAbhNGrM0Fw0WpWxh66N9Xp meirish@Macintosh-3.local`;

const ROLES = [
  {
    type: 'ca',
    name: 'carole',
    fillInCreate() {
      click('[data-test-input="allowUserCertificates"]');
    },
    fillInGenerate() {
      fillIn('[data-test-input="publicKey"]', PUB_KEY);
    },
    assertAfterGenerate(assert, sshPath) {
      assert.equal(currentURL(), `/vault/secrets/${sshPath}/sign/${this.name}`, 'ca sign url is correct');
      assert.dom('[data-test-row-label="Signed key"]').exists({ count: 1 }, 'renders the signed key');
      assert.dom('[data-test-row-value="Signed key"]').exists({ count: 1 }, "renders the signed key's value");
      assert.dom('[data-test-row-label="Serial number"]').exists({ count: 1 }, 'renders the serial');
      assert.dom('[data-test-row-value="Serial number"]').exists({ count: 1 }, 'renders the serial value');
    },
  },
  {
    type: 'otp',
    name: 'otprole',
    fillInCreate() {
      fillIn('[data-test-input="defaultUser"]', 'admin');
      click('[data-test-toggle-more]');
      fillIn('[data-test-input="cidrList"]', '1.2.3.4/32');
    },
    fillInGenerate() {
      fillIn('[data-test-input="username"]', 'admin');
      fillIn('[data-test-input="ip"]', '1.2.3.4');
    },
    assertAfterGenerate(assert, sshPath) {
      assert.equal(
        currentURL(),
        `/vault/secrets/${sshPath}/credentials/${this.name}`,
        'otp credential url is correct'
      );
      assert.dom('[data-test-row-label="Key"]').exists({ count: 1 }, 'renders the key');
      assert.dom('[data-test-row-value="Key"]').exists({ count: 1 }, "renders the key's value");
      assert.dom('[data-test-row-label="Port"]').exists({ count: 1 }, 'renders the port');
      assert.dom('[data-test-row-value="Port"]').exists({ count: 1 }, "renders the port's value");
    },
  },
];
test('ssh backend', function(assert) {
  const now = new Date().getTime();
  const sshPath = `ssh-${now}`;

  mountSupportedSecretBackend(assert, 'ssh', sshPath);
  click('[data-test-secret-backend-configure]');
  andThen(() => {
    assert.equal(currentURL(), `/vault/settings/secrets/configure/${sshPath}`);
    assert.ok(find('[data-test-ssh-configure-form]').length, 'renders the empty configuration form');
  });

  // default has generate CA checked so we just submit the form
  click('[data-test-ssh-input="configure-submit"]');
  andThen(() => {
    assert.ok(find('[data-test-ssh-input="public-key"]').length, 'a public key is fetched');
  });
  click('[data-test-backend-view-link]');

  //back at the roles list
  andThen(() => {
    assert.equal(currentURL(), `/vault/secrets/${sshPath}/list`, `redirects to ssh index`);
  });

  ROLES.forEach(role => {
    // create a role
    click('[ data-test-secret-create]');
    andThen(() => {
      assert.ok(
        find('[data-test-secret-header]').text().includes('SSH Role'),
        `${role.type}: renders the create page`
      );
    });

    fillIn('[data-test-input="name"]', role.name);
    fillIn('[data-test-input="keyType"]', role.type);
    role.fillInCreate();

    // save the role
    click('[data-test-role-ssh-create]');
    andThen(() => {
      assert.equal(
        currentURL(),
        `/vault/secrets/${sshPath}/show/${role.name}`,
        `${role.type}: navigates to the show page on creation`
      );
    });

    // sign a key with this role
    click('[data-test-backend-credentials]');
    role.fillInGenerate();

    // generate creds
    click('[data-test-secret-generate]');
    andThen(() => {
      role.assertAfterGenerate(assert, sshPath);
    });

    // click the "Back" button
    click('[data-test-secret-generate-back]');
    andThen(() => {
      assert.ok(
        find('[data-test-secret-generate-form]').length,
        `${role.type}: back takes you back to the form`
      );
    });

    click('[data-test-secret-generate-cancel]');
    //back at the roles list
    andThen(() => {
      assert.equal(
        currentURL(),
        `/vault/secrets/${sshPath}/list`,
        `${role.type}: cancel takes you to ssh index`
      );
      assert.ok(
        find(`[data-test-secret-link="${role.name}"]`).length,
        `${role.type}: role shows in the list`
      );
    });

    //and delete
    click(`[data-test-secret-link="${role.name}"] [data-test-popup-menu-trigger]`);
    andThen(() => {
      click(`[data-test-ssh-role-delete="${role.name}"] button`);
    });
    click(`[data-test-confirm-button]`);

    andThen(() => {
      assert
        .dom(`[data-test-secret-link="${role.name}"]`)
        .doesNotExist(`${role.type}: role is no longer in the list`);
    });
  });
});
