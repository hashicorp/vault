/**
 * Copyright IBM Corp. 2016, 2025
 * SPDX-License-Identifier: BUSL-1.1
 */

import { module, test } from 'qunit';
import { setupRenderingTest } from 'vault/tests/helpers';
import { render } from '@ember/test-helpers';
import { hbs } from 'ember-cli-htmlbars';
import { getErrorResponse } from 'vault/tests/helpers/api/error-response';
import { GENERAL } from 'vault/tests/helpers/general-selectors';
import sinon from 'sinon';

module('Integration | Component | page/error', function (hooks) {
  setupRenderingTest(hooks);

  hooks.beforeEach(function () {
    this.error = undefined;
    this.isFullPage = undefined;
    this.titleTag = undefined;
    this.renderComponent = () =>
      render(hbs`<Page::Error @error={{this.error}} @isFullPage={{this.isFullPage}} @titleTag={{this.titleTag}} />
    `);
  });

  test('it handles undefined args', async function (assert) {
    await this.renderComponent();

    assert.dom(GENERAL.pageError.title()).hasText('Error');
    assert
      .dom(GENERAL.pageError.message)
      .hasText('A problem has occurred. Check the Vault logs or console for more details.');
    assert.dom(GENERAL.icon()).doesNotExist('it does not render an icon when there is no error code');
  });

  test('it should render 404 error', async function (assert) {
    this.error = {
      httpStatus: 404,
      path: '/v1/kubernetes/config',
    };
    await this.renderComponent();

    assert.dom(GENERAL.icon('alert-circle')).exists('it renders alert-circle icon for 404 error');
    assert.dom(GENERAL.pageError.title(404)).hasText('ERROR 404 Not found', 'Error title renders');
    assert
      .dom(GENERAL.pageError.message)
      .hasText(`Sorry, we were unable to find any content at ${this.error.path}.`, 'Error message renders');
  });

  test('it should render 403 error', async function (assert) {
    this.error = {
      httpStatus: 403,
      path: '/v1/kubernetes/config',
    };
    await this.renderComponent();

    assert.dom(GENERAL.icon('skip')).exists('it renders skip icon for 403 error');
    assert.dom(GENERAL.pageError.title(403)).hasText('ERROR 403 Not authorized', 'Error title renders');
    assert
      .dom(GENERAL.pageError.message)
      .hasText(`You are not authorized to access content at ${this.error.path}.`, 'Error message renders');
  });

  test('it should render error codes that do not have default messages', async function (assert) {
    this.error = {
      httpStatus: 400,
      errors: ['something has gone wrong'],
    };
    await this.renderComponent();

    assert.dom(GENERAL.icon('alert-circle')).exists('it renders alert-circle icon for other error codes');
    assert.dom(GENERAL.pageError.title(400)).hasText('ERROR 400 Error', 'Error title renders');
    assert
      .dom(GENERAL.pageError.message)
      .hasText('A problem has occurred. Check the Vault logs or console for more details.');
    assert.dom(GENERAL.pageError.details).hasText('something has gone wrong');
  });

  test('it should render default message when message contains "permission denied"', async function (assert) {
    this.error = {
      message: '1 error occurred:\n\t* permission denied\n\n',
      status: 403,
      path: '/v1/sys/config/ui/login/default-auth/?list=true',
      errorURL: '/vault/config-ui/login-settings',
    };
    await this.renderComponent();

    assert.dom(GENERAL.pageError.title(403)).hasText('ERROR 403 Not authorized');
    assert
      .dom(GENERAL.pageError.message)
      .hasText(
        'You are not authorized to access content at /v1/sys/config/ui/login/default-auth/?list=true.'
      );
  });

  test('it should use default message when error contains "Ember Data"', async function (assert) {
    this.error = {
      message:
        'Ember Data Request GET /v1/sys/policies/acl returned a 403\nPayload (application/json)\n{\n  "errors": [\n    "1 error occurred:\\n\\t* permission denied\\n\\n"\n  ]\n}',
      errors: ['1 error occurred:\n\t* permission denied\n\n'],
      httpStatus: 403,
      path: '/v1/sys/policies/acl',
      errorURL: '/vault/policies/acl',
    };
    await this.renderComponent();

    assert.dom(GENERAL.pageError.title(403)).hasText('ERROR 403 Not authorized');
    assert
      .dom(GENERAL.pageError.message)
      .hasText('You are not authorized to access content at /v1/sys/policies/acl.');
  });

  test('it updates when error arg changes', async function (assert) {
    this.error = { httpStatus: 403, path: '/v1/kubernetes/config' };
    await this.renderComponent();

    assert
      .dom(GENERAL.pageError.message)
      .hasText(
        'You are not authorized to access content at /v1/kubernetes/config.',
        'Initial message renders'
      );
    // Change error arg using "set" to trigger reactivity
    this.set('error', { httpStatus: 403, path: '/v1/kubernetes/roles' });
    assert
      .dom(GENERAL.pageError.message)
      .hasText(
        'You are not authorized to access content at /v1/kubernetes/roles.',
        'Updated message renders'
      );
  });

  test('it should render general error without http status', async function (assert) {
    this.error = {
      message: 'An unexpected error occurred',
      errors: ['This is one thing that went wrong', 'Unfortunately something else went wrong too'],
    };
    await this.renderComponent();

    assert.dom(GENERAL.pageError.title()).hasText('Error');
    assert.dom(GENERAL.pageError.message).hasText(this.error.message);
    this.error.errors.forEach((error, index) => {
      assert.dom(`[data-test-page-error-details="${index}"]`).hasText(this.error.errors[index]);
    });
  });

  test('it should handle 404 api client errors', async function (assert) {
    this.error = getErrorResponse();
    await this.renderComponent();

    assert.dom(GENERAL.pageError.title(404)).hasText('ERROR 404 Not found');
    assert
      .dom(GENERAL.pageError.message)
      .hasText('Sorry, we were unable to find any content at /v1/test/error/parsing.');
  });

  test('it should handle 403 api client errors', async function (assert) {
    this.error = getErrorResponse({ errors: ['permission denied'] }, 403);
    await this.renderComponent();

    assert.dom(GENERAL.pageError.title(403)).hasText('ERROR 403 Not authorized');
    assert
      .dom(GENERAL.pageError.message)
      .hasText('You are not authorized to access content at /v1/test/error/parsing.');
  });

  test('it should handle api client errors that are not 403 or 404', async function (assert) {
    const error = { errors: ['bad things occurred'] };
    this.error = getErrorResponse(error, 500);
    await this.renderComponent();

    assert.dom(GENERAL.pageError.title(500)).hasText('ERROR 500 Error');
    assert.dom(GENERAL.pageError.message).hasText(error.errors[0]);
  });

  test('it should handle api client errors that are already parsed', async function (assert) {
    const error = { errors: ['oh dear!'] };
    const api = this.owner.lookup('service:api');
    this.error = await api.parseError(getErrorResponse(error, 500));
    await this.renderComponent();

    assert.dom(GENERAL.pageError.title(500)).hasText('ERROR 500 Error');
    assert.dom(GENERAL.pageError.message).hasText(error.errors[0]);
  });

  // COMPONENT STYLING
  test('it defaults to inline error when @isFullPage is not provided', async function (assert) {
    this.error = {
      httpStatus: 404,
      path: '/v1/kubernetes/config',
    };
    await this.renderComponent();

    assert.dom('h1').doesNotExist('it does not render title in an h1 tag');
    assert
      .dom(`${GENERAL.pageError.error} div`)
      .doesNotHaveClass('align-self-center', 'defaults to inline styling');
    assert.dom(`${GENERAL.pageError.error} div`).hasClass('top-padding-32', 'has top padding by default');
  });

  test('it renders full page with h1 title, center alignment and footer', async function (assert) {
    this.isFullPage = true;
    this.error = {
      message: 'An unexpected error occurred',
      errors: ['Something went wrong'],
    };
    await this.renderComponent();

    assert.dom('h1').exists().hasText('Error', 'it renders title as an h1 tag');
    assert.dom(`${GENERAL.pageError.error} div`).hasClass('align-self-center');
    assert.dom(`${GENERAL.pageError.error} div`).doesNotHaveClass('top-padding-32');
    assert.dom(GENERAL.pageError.message).hasText(this.error.message);
    assert
      .dom(GENERAL.pageError.error)
      .hasTextContaining(
        'Double check the URL or return to the dashboard.',
        'additional message renders for full page'
      );
    assert.dom('a').hasText('Go to dashboard', 'Dashboard link renders');
  });

  test('it renders passed @titleTag even when @isFullPage is true', async function (assert) {
    this.titleTag = 'h3';
    this.isFullPage = true;
    await this.renderComponent();

    assert.dom('h3').hasText('Error');
    assert.dom('h1').doesNotExist();
  });

  test('it renders media block when provided', async function (assert) {
    this.error = {
      httpStatus: 404,
      path: '/v1/kubernetes/config',
    };
    await render(hbs`
      <Page::Error @error={{this.error}}>
        <:media>
          <span data-test-page-error-media>Custom media</span>
        </:media>
      </Page::Error>
    `);

    assert.dom('[data-test-page-error-media]').exists().hasText('Custom media');
  });

  test('it renders customFooter block instead of default footer when isFullPage is true', async function (assert) {
    this.error = {
      httpStatus: 404,
      path: '/v1/kubernetes/config',
    };
    await render(hbs`
      <Page::Error @error={{this.error}} @isFullPage={{true}}>
        <:customFooter as |A|>
          <A.Body @text="My custom footer situation" />
          <A.Footer as |F|>
            <F.LinkStandalone @icon="wand" @text="Custom action" @href="/" data-test-custom-footer />
          </A.Footer>
        </:customFooter>
      </Page::Error>
    `);

    assert
      .dom('[data-test-custom-footer]')
      .exists('custom footer renders instead of "Go to dashboard" message');
    assert.dom('a').exists({ count: 1 }, 'only 1 link renders');
    assert.dom('a').hasText('Custom action');
  });

  test('it renders path instead of errorURL if both exist', async function (assert) {
    this.error = {
      httpStatus: 404,
      path: '/v1/kubernetes/config',
      errorURL: 'vault/secrets-engines/pki_int',
    };
    await this.renderComponent();

    assert
      .dom(GENERAL.pageError.message)
      .hasText(`Sorry, we were unable to find any content at ${this.error.path}.`);
  });

  test('it renders errorURL if path does not exist', async function (assert) {
    const router = this.owner.lookup('service:router');
    const currentURLStub = sinon.stub(router, 'currentURL');
    currentURLStub.value('/vault/secret-engines/kv');
    this.error = {
      httpStatus: 404,
      errorURL: 'vault/secrets-engines/pki_int',
    };
    await this.renderComponent();

    assert
      .dom(GENERAL.pageError.message)
      .hasText(`Sorry, we were unable to find any content at ${this.error.errorURL}.`);
    currentURLStub.restore();
  });

  test('it renders current URL from router if no path or errorURL exists', async function (assert) {
    const router = this.owner.lookup('service:router');
    const currentURLStub = sinon.stub(router, 'currentURL');
    currentURLStub.value('/vault/secret-engines/kv');
    this.error = { httpStatus: 404 };
    await this.renderComponent();

    assert
      .dom(GENERAL.pageError.message)
      .hasText(`Sorry, we were unable to find any content at /vault/secret-engines/kv.`);
    currentURLStub.restore();
  });
});
