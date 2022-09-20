import { module, test } from 'qunit';
import { setupRenderingTest } from 'ember-qunit';
import { render } from '@ember/test-helpers';
import { hbs } from 'ember-cli-htmlbars';
import { setupMirage } from 'ember-cli-mirage/test-support';
import { statuses } from '../../../mirage/handlers/hcp-link';

const timestamp = '[2022-09-13 14:45:40.666697 -0700 PDT]';

module('Integration | Component | link-status', function (hooks) {
  setupRenderingTest(hooks);
  setupMirage(hooks);

  // this can be removed once feature is released for OSS
  hooks.beforeEach(function () {
    this.owner.lookup('service:version').set('isEnterprise', true);
    this.statuses = statuses;
  });

  test('it does not render banner or error when status is not present', async function (assert) {
    await render(hbs`<LinkStatus @status={{undefined}} />`);

    assert.dom('.navbar-status').doesNotExist('Banner is hidden for disconnected state');
  });

  test('it does not render banner if not enterprise version', async function (assert) {
    this.owner.lookup('service:version').set('isEnterprise', false);

    await render(hbs`<LinkStatus @status={{get this.statuses 0}} />`);

    assert.dom('.navbar-status').doesNotExist('Banner is hidden for disconnected state');
  });

  test('it renders connected status', async function (assert) {
    await render(hbs`<LinkStatus @status={{get this.statuses 0}} />`);

    assert.dom('.navbar-status').hasClass('connected', 'Correct class renders for connected state');
    assert
      .dom('[data-test-link-status]')
      .hasText('This self-managed Vault is linked to HCP.', 'Copy renders for connected state');
    assert
      .dom('[data-test-link-status] a')
      .hasAttribute('href', 'https://portal.cloud.hashicorp.com/sign-in', 'HCP sign in link renders');
  });

  test('it does not render banner for disconnected state with unknown error', async function (assert) {
    await render(hbs`<LinkStatus @status={{get this.statuses 1}} />`);

    assert.dom('.navbar-status').doesNotExist('Banner is hidden for disconnected state');
  });

  test('it should render for disconnected error state', async function (assert) {
    await render(hbs`<LinkStatus @status={{get this.statuses 2}} />`);

    assert.dom('.navbar-status').hasClass('warning', 'Correct class renders for disconnected error state');
    assert
      .dom('[data-test-link-status]')
      .hasText(
        `Vault has been disconnected from HCP since ${timestamp}. Error: some other error other than unknown`,
        'Copy renders for disconnected error state'
      );
  });

  test('it should render for connection refused error state', async function (assert) {
    await render(hbs`<LinkStatus @status={{get this.statuses 3}} />`);

    assert
      .dom('.navbar-status')
      .hasClass('warning', 'Correct class renders for connection refused error state');
    assert
      .dom('[data-test-link-status]')
      .hasText(
        `Vault has been trying to connect to HCP since ${timestamp}, but HCP is not reachable. Vault will try again soon.`,
        'Copy renders for connection refused error state'
      );
  });

  test('it should render for resource id error state', async function (assert) {
    await render(hbs`<LinkStatus @status={{get this.statuses 4}} />`);

    assert.dom('.navbar-status').hasClass('warning', 'Correct class renders for resource id error state');
    assert
      .dom('[data-test-link-status]')
      .hasText(
        `Vault tried connecting to HCP, but the Resource ID is invalid. Check your resource ID. ${timestamp}`,
        'Copy renders for resource id error state'
      );
  });

  test('it should render for unauthorized error state', async function (assert) {
    await render(hbs`<LinkStatus @status={{get this.statuses 5}} />`);

    assert.dom('.navbar-status').hasClass('warning', 'Correct class renders for unauthorized error state');
    assert
      .dom('[data-test-link-status]')
      .hasText(
        `Vault tried connecting to HCP, but the authorization information is wrong. Update it and try again. ${timestamp}`,
        'Copy renders for unauthorized error state'
      );
  });

  test('it should render generic message for unknown error state', async function (assert) {
    await render(hbs`<LinkStatus @status={{get this.statuses 6}} />`);

    assert.dom('.navbar-status').hasClass('warning', 'Correct class renders for unknown error state');
    assert
      .dom('[data-test-link-status]')
      .hasText(
        `Vault has been trying to connect to HCP since ${timestamp}. Vault will try again soon. Error: connection error we are unaware of`,
        'Copy renders for unknown error state'
      );
  });

  // connecting state should always be returned with timestamp and error
  // this case came up in manual testing and should be fixed on the backend but additional checks were added just in case
  test('it renders generic message for connecting state with no timestamp or error', async function (assert) {
    await render(hbs`<LinkStatus @status="connecting" />`);

    assert.dom('.navbar-status').hasClass('warning', 'Correct class renders for unknown error state');
    assert
      .dom('[data-test-link-status]')
      .hasText(
        `Vault has been trying to connect to HCP. Vault will try again soon.`,
        'Copy renders for unknown error state'
      );
  });
});
