import { module, test } from 'qunit';

import { setupMirage } from 'ember-cli-mirage/test-support';
import authPage from 'vault/tests/pages/auth';
import sinon from 'sinon';

import { setupApplicationTest } from 'vault/tests/helpers';
import { PROVIDER_NAME as DummyProviderName } from 'vault/utils/analytics-providers/dummy';

module('Acceptance | analytics', function (hooks) {
  setupApplicationTest(hooks);
  setupMirage(hooks);

  hooks.beforeEach(function () {
    // spy on the analytics service
    this.analyticsService = this.owner.lookup('service:analytics');
    this.analyticsServiceSpy = sinon.spy(this.analyticsService);
  });

  test('initialization flow works for authenticated users', async function (assert) {
    await authPage.login();

    // provider start is called as expected
    assert.true(this.analyticsServiceSpy.start.calledOnce, 'the service is started');
    assert.true(
      this.analyticsServiceSpy.start.calledWith(DummyProviderName, {
        enabled: true,
        API_KEY: 'DUMMY_KEY',
        api_host: 'whatever',
      }),
      'the service is started with the expected config'
    );

    // user identify is called as expected
    assert.true(this.analyticsServiceSpy.identifyUser.calledOnce, 'the user is identified');

    // basic route transitions work
    assert.true(this.analyticsServiceSpy.trackPageView.calledWith('/vault/dashboard'));
  });
});
