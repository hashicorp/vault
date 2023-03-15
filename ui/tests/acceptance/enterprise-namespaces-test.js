import { click, settled, visit, fillIn, currentURL } from '@ember/test-helpers';
import { module, test } from 'qunit';
import { setupApplicationTest } from 'ember-qunit';
import { create } from 'ember-cli-page-object';
import { setupMirage } from 'ember-cli-mirage/test-support';

import consoleClass from 'vault/tests/pages/components/console/ui-panel';
import authPage from 'vault/tests/pages/auth';
import logout from 'vault/tests/pages/logout';

const shell = create(consoleClass);

const createNS = async (name) => {
  await shell.runCommands(`write sys/namespaces/${name} -force`);
};

module('Acceptance | Enterprise | namespaces', function (hooks) {
  setupApplicationTest(hooks);

  hooks.beforeEach(function () {
    return authPage.login();
  });

  test('it clears namespaces when you log out', async function (assert) {
    const ns = 'foo';
    await createNS(ns);
    await shell.runCommands(`write -field=client_token auth/token/create policies=default`);
    const token = shell.lastLogOutput;
    await logout.visit();
    await authPage.login(token);
    assert.dom('[data-test-namespace-toggle]').doesNotExist('does not show the namespace picker');
    await logout.visit();
  });

  test('it shows nested namespaces if you log in with a namspace starting with a /', async function (assert) {
    const nses = ['beep', 'boop', 'bop'];
    for (const [i, ns] of nses.entries()) {
      await createNS(ns);
      await settled();
      // this is usually triggered when creating a ns in the form, here we'll trigger a reload of the
      // namespaces manually
      await this.owner.lookup('service:namespace').findNamespacesForUser.perform();
      if (i === nses.length - 1) {
        break;
      }
      // the namespace path will include all of the namespaces up to this point
      const targetNamespace = nses.slice(0, i + 1).join('/');
      const url = `/vault/secrets?namespace=${targetNamespace}`;
      // check if namespace is in the toggle
      await click('[data-test-namespace-toggle]');

      // check that the single namespace "beep" or "boop" not "beep/boop" shows in the toggle display
      assert
        .dom(`[data-test-namespace-link="${targetNamespace}"]`)
        .hasText(ns, 'shows the namespace in the toggle component');
      // close toggle
      await click('[data-test-namespace-toggle]');
      // because quint does not like page reloads, visiting url directing instead of clicking on namespace in toggle
      await visit(url);
    }
    await logout.visit();
    await settled();
    await authPage.visit({ namespace: '/beep/boop' });
    await settled();
    await authPage.tokenInput('root').submit();
    await settled();
    await click('[data-test-namespace-toggle]');

    assert.dom('[data-test-current-namespace]').hasText('/beep/boop/', 'current namespace begins with a /');
    assert
      .dom('[data-test-namespace-link="beep/boop/bop"]')
      .exists('renders the link to the nested namespace');
  });

  test('it shows the regular namespace toolbar when not managed', async function (assert) {
    // This test is the opposite of the test in managed-namespace-test
    await logout.visit();
    assert.strictEqual(currentURL(), '/vault/auth?with=token', 'Does not redirect');
    assert.dom('[data-test-namespace-toolbar]').exists('Normal namespace toolbar exists');
    assert
      .dom('[data-test-managed-namespace-toolbar]')
      .doesNotExist('Managed namespace toolbar does not exist');
    assert.dom('input#namespace').hasAttribute('placeholder', '/ (Root)');
    await fillIn('input#namespace', '/foo');
    const encodedNamespace = encodeURIComponent('/foo');
    assert.strictEqual(
      currentURL(),
      `/vault/auth?namespace=${encodedNamespace}&with=token`,
      'Does not prepend root to namespace'
    );
  });

  module('auth form', function (hooks) {
    setupMirage(hooks);
    const SELECTORS = {
      authTab: (path) => `[data-test-auth-method="${path}"] a`,
      authSubmit: '[data-test-auth-submit]',
    };
    hooks.beforeEach(async function () {
      this.namespace = 'testns';
      this.rootOidc = 'root-oidc';
      this.nsOidc = 'ns-oidc';

      const enableOidc = async (path, role = null) => {
        this.server.post(`/auth/${path}/config`, () => {});
        const hasRole = role || '';
        await shell.runCommands([
          `write sys/auth/${path} type=oidc`,
          `write auth/${path}/config default_role="${hasRole}" oidc_discovery_url="https://example.com"`,
          // show method as tab
          `write sys/auth/${path}/tune listing_visibility="unauth"`,
        ]);
      };
      await authPage.login();
      // enable oidc in root namespace, without default role
      await enableOidc(this.rootOidc);
      // create child namespace to enable oidc
      await createNS(this.namespace);
      await logout.visit();

      // enable oidc in child namespace with default role
      await authPage.loginNs(this.namespace);
      await enableOidc(this.nsOidc, `${this.nsOidc}-role`);
      return await authPage.logout();
    });

    hooks.afterEach(async function () {
      const disableOidc = async (path) => {
        await shell.runCommands([`delete /sys/auth/${path}`]);
      };

      await authPage.loginNs(this.namespace);
      await visit(`/vault/access?namespace=${this.namespace}`);
      // disable methods to cleanup test state for re-running
      await disableOidc(this.rootOidc);
      await disableOidc(this.nsOidc);
      await authPage.logout();

      await authPage.login();
      await shell.runCommands([`delete /sys/auth/${this.namespace}`]);
    });

    test('auth form updates when a namespace is entered', async function (assert) {
      assert.expect(5);
      this.server.post(`/auth/${this.rootOidc}/oidc/auth_url`, (schema, req) => {
        const request = JSON.parse(req.requestBody);
        assert.deepEqual(
          request.redirect_uri,
          `http://localhost:7357/ui/vault/auth/${this.rootOidc}/oidc/callback`,
          'request made to auth_url when the login page is visited'
        );
      });
      this.server.post(`/auth/${this.nsOidc}/oidc/auth_url`, (schema, req) => {
        const request = JSON.parse(req.requestBody);
        assert.deepEqual(
          request.redirect_uri,
          `http://localhost:7357/ui/vault/auth/${this.nsOidc}/oidc/callback?namespace=testns`,
          'request made to correct auth_url when namespace is filled in'
        );
      });
      await visit('/vault/auth');
      assert.dom(SELECTORS.authTab(this.rootOidc)).exists('renders oidc method tab for root');
      await authPage.namespaceInput(this.namespace);
      assert.strictEqual(
        currentURL(),
        `/vault/auth?namespace=${this.namespace}&with=ns-oidc%2F`,
        'url updates with namespace value'
      );
      assert.dom(SELECTORS.authTab(this.nsOidc)).exists('renders oidc method tab for namespace');
    });
  });
});
