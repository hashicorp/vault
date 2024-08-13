/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: BUSL-1.1
 */

import { module, test } from 'qunit';
import sinon from 'sinon';
import { setupRenderingTest } from 'vault/tests/helpers';
import { GENERAL } from 'vault/tests/helpers/general-selectors';
import { SECRET_ENGINE_SELECTORS as SES } from 'vault/tests/helpers/secret-engine/secret-engine-selectors';
import { render, click, fillIn } from '@ember/test-helpers';
import { setupMirage } from 'ember-cli-mirage/test-support';
import { hbs } from 'ember-cli-htmlbars';
import { v4 as uuidv4 } from 'uuid';
import { createConfig } from 'vault/tests/helpers/secret-engine/secret-engine-helpers';

module('Integration | Component | SecretEngine/configure-aws', function (hooks) {
  setupRenderingTest(hooks);
  setupMirage(hooks);

  hooks.beforeEach(function () {
    this.store = this.owner.lookup('service:store');

    this.uid = uuidv4();
    this.id = `aws-${this.uid}`;
    this.model = createConfig(this.store, this.id, 'aws-lease'); // currently when you queryRecord for secret-engine type aws it returns the lease/config. This is going to change in the refactor.
    this.saveAWSLease = sinon.stub();
    this.saveAWSRoot = sinon.stub();

    this.renderComponent = () => {
      return render(hbs`
        <SecretEngine::ConfigureAws @model={{this.model}} @saveAWSLease={{this.saveAWSLease}} @saveAWSRoot={{this.saveAWSRoot}} @tab="root" @region="" />
        `);
    };
  });

  test('it renders fields', async function (assert) {
    await this.renderComponent();
    assert.dom(SES.aws.rootForm).exists('it lands on the aws root configuration form.');
    assert.dom(GENERAL.inputByAttr('accessKey')).exists(`accessKey shows for Access section.`);
    assert.dom(GENERAL.inputByAttr('secretKey')).exists(`secretKey shows for Access section.`);

    await click(GENERAL.hdsTab('lease'));
    assert.dom('[data-test-ttl-form-label="Lease"]').exists('Lease TTL is rendered');
    assert.dom('[data-test-ttl-form-label="Maximum Lease"]').exists('Maximum Lease TTL is rendered');
  });

  test('it calls saveAWSRoot on save root config', async function (assert) {
    await this.renderComponent();
    await fillIn(GENERAL.inputByAttr('accessKey'), 'foo');
    await fillIn(GENERAL.inputByAttr('secretKey'), 'bar');
    await click(SES.aws.saveRootConfig);
    assert.ok(this.saveAWSRoot.calledOnce, 'saveAWSRoot was called once');
    assert.ok(this.saveAWSLease.notCalled, 'saveAWSLease was not called');
  });

  test('it calls saveAWSLease on save lease config', async function (assert) {
    await this.renderComponent();
    // createLease config already has ttls set so just save the values
    await click(SES.aws.saveLeaseConfig);
    assert.ok(this.saveAWSLease.calledOnce, 'saveAWSLease was called once');
    assert.ok(this.saveAWSRoot.notCalled, 'saveAWSRoot was not called');
  });
});
