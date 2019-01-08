import { addMinutes } from 'date-fns';
import { module, test } from 'qunit';
import { setupRenderingTest } from 'ember-qunit';
import { render } from '@ember/test-helpers';
import hbs from 'htmlbars-inline-precompile';
import sinon from 'sinon';
import { create } from 'ember-cli-page-object';
import license from '../../pages/components/license-info';
import { allFeatures } from 'vault/helpers/all-features';

const FEATURES = allFeatures();

const component = create(license);

module('Integration | Component | license info', function(hooks) {
  setupRenderingTest(hooks);

  const LICENSE_WARNING_TEXT = `Warning Your temporary license expires in 30 minutes and your vault will seal. Please enter a valid license below.`;

  test('it renders properly for temporary license', async function(assert) {
    const now = Date.now();
    this.set('licenseId', 'temporary');
    this.set('expirationTime', addMinutes(now, 30));
    this.set('startTime', now);
    this.set('features', ['HSM', 'Namespaces']);
    await render(
      hbs`<LicenseInfo @licenseId={{this.licenseId}} @expirationTime={{this.expirationTime}} @startTime={{this.startTime}} @features={{this.features}}/>`
    );
    assert.equal(component.warning, LICENSE_WARNING_TEXT, 'it renders warning text including time left');
    assert.equal(component.hasSaveButton, true, 'it renders the save button');
    assert.equal(component.hasTextInput, true, 'it renders text input for new license');
    assert.equal(component.featureRows.length, FEATURES.length, 'it renders all of the features');
    assert.equal(component.featureRows[0].featureName, 'HSM', 'it renders HSM feature');
    assert.equal(component.featureRows[0].featureStatus, 'Active', 'it renders Active for HSM feature');
    assert.equal(
      component.featureRows[1].featureName,
      'Performance Replication',
      'it renders Performance Replication feature name'
    );
    assert.equal(
      component.featureRows[1].featureStatus,
      'Not Active',
      'it renders Not Active for Performance Replication'
    );
  });

  test('it renders feature status properly for features associated with license', async function(assert) {
    const now = Date.now();
    this.set('licenseId', 'temporary');
    this.set('expirationTime', addMinutes(now, 30));
    this.set('startTime', now);
    this.set('features', ['HSM', 'Namespaces']);
    await render(
      hbs`<LicenseInfo @licenseId={{this.licenseId}} @expirationTime={{this.expirationTime}} @startTime={{this.startTime}} @features={{this.features}}/>`
    );
    assert.equal(component.featureRows.length, FEATURES.length, 'it renders all of the features');
    let activeFeatures = component.featureRows.filter(f => f.featureStatus === 'Active');
    assert.equal(activeFeatures.length, 2);
  });

  test('it renders properly for non-temporary license', async function(assert) {
    const now = Date.now();
    this.set('licenseId', 'test');
    this.set('expirationTime', addMinutes(now, 30));
    this.set('startTime', now);
    this.set('features', ['HSM', 'Namespaces']);
    await render(
      hbs`<LicenseInfo @licenseId={{this.licenseId}} @expirationTime={{this.expirationTime}} @startTime={{this.startTime}} @features={{this.features}}/>`
    );
    assert.equal(component.hasWarning, false, 'it does not have a warning');
    assert.equal(component.hasSaveButton, false, 'it does not render the save button');
    assert.equal(component.hasTextInput, false, 'it does not render the text input for new license');
    assert.equal(component.hasEnterButton, true, 'it renders the button to toggle license form');
  });

  test('it shows and hides license form when enter and cancel buttons are clicked', async function(assert) {
    const now = Date.now();
    this.set('licenseId', 'test');
    this.set('expirationTime', addMinutes(now, 30));
    this.set('startTime', now);
    this.set('features', ['HSM', 'Namespaces']);
    await render(
      hbs`<LicenseInfo @licenseId={{this.licenseId}} @expirationTime={{this.expirationTime}} @startTime={{this.startTime}} @features={{this.features}}/>`
    );
    await component.enterButton();
    assert.equal(component.hasSaveButton, true, 'it does not render the save button');
    assert.equal(component.hasTextInput, true, 'it does not render the text input for new license');
    assert.equal(component.hasEnterButton, false, 'it renders the button to toggle license form');
    await component.cancelButton();
    assert.equal(component.hasSaveButton, false, 'it does not render the save button');
    assert.equal(component.hasTextInput, false, 'it does not render the text input for new license');
    assert.equal(component.hasEnterButton, true, 'it renders the button to toggle license form');
  });

  test('it calls saveModel when save button is clicked', async function(assert) {
    const now = Date.now();
    this.set('licenseId', 'temporary');
    this.set('expirationTime', addMinutes(now, 30));
    this.set('startTime', now);
    this.set('features', ['HSM', 'Namespaces']);
    this.set('saveModel', sinon.spy());
    await render(
      hbs`<LicenseInfo @licenseId={{this.licenseId}} @expirationTime={{this.expirationTime}} @startTime={{this.startTime}} @features={{this.features}} @saveModel={{this.saveModel}}/>`
    );
    await component.text('ABCDE12345');
    await component.saveButton();
    assert.ok(this.get('saveModel').calledOnce);
  });

  test('it renders Performance Standby as inactive if count is 0', async function(assert) {
    const now = Date.now();
    this.set('licenseId', 'temporary');
    this.set('expirationTime', addMinutes(now, 30));
    this.set('startTime', now);
    this.set('model', { performanceStandbyCount: 0 });
    this.set('features', ['Performance Standby', 'Namespaces']);

    await render(
      hbs`<LicenseInfo @licenseId={{this.licenseId}} @expirationTime={{this.expirationTime}} @startTime={{this.startTime}} @features={{this.features}} @model={{this.model}}/>`
    );

    let row = component.featureRows.filterBy('featureName', 'Performance Standby')[0];
    assert.equal(row.featureStatus, 'Not Active', 'renders feature as inactive because when count is 0');
  });

  test('it renders Performance Standby as active and shows count', async function(assert) {
    const now = Date.now();
    this.set('licenseId', 'temporary');
    this.set('expirationTime', addMinutes(now, 30));
    this.set('startTime', now);
    this.set('model', { performanceStandbyCount: 4 });
    this.set('features', ['Performance Standby', 'Namespaces']);

    await render(
      hbs`<LicenseInfo @licenseId={{this.licenseId}} @expirationTime={{this.expirationTime}} @startTime={{this.startTime}} @features={{this.features}} @model={{this.model}}/>`
    );

    let row = component.featureRows.filterBy('featureName', 'Performance Standby')[0];
    assert.equal(row.featureStatus, 'Active — 4 standby nodes allotted', 'renders active and displays count');
  });
});
