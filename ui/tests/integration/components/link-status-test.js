import { module, test } from 'qunit';
import { setupRenderingTest } from 'ember-qunit';
import { render } from '@ember/test-helpers';
import { hbs } from 'ember-cli-htmlbars';
import { setupMirage } from 'ember-cli-mirage/test-support';

module('Integration | Component | link-status', function (hooks) {
  setupRenderingTest(hooks);
  setupMirage(hooks);

  // this can be removed once feature is released for OSS
  hooks.beforeEach(function () {
    this.owner.lookup('service:version').set('isEnterprise', true);
  });

  test('it does not render disconnected status', async function (assert) {
    this.server.get('sys/seal-status', () => ({ hcp_link_status: 'disconnected' }));

    await render(hbs`<LinkStatus />`);

    assert.dom('.navbar-status').doesNotExist('Banner is hidden for disconnected state');
  });

  test('it renders connected status', async function (assert) {
    this.server.get('sys/seal-status', () => ({ hcp_link_status: 'connected' }));

    await render(hbs`<LinkStatus />`);

    assert.dom('.navbar-status').hasClass('connected', 'Correct class renders for connected state');
    assert
      .dom('[data-test-link-status]')
      .hasText(
        'This self-managed Vault is linked to the HashiCorp Cloud Platform',
        'Copy renders for connected state'
      );
    assert
      .dom('[data-test-link-status] a')
      .hasAttribute('href', 'https://portal.cloud.hashicorp.com/sign-in', 'HCP sign in link renders');
  });
});
