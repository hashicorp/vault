import { moduleForComponent, test } from 'ember-qunit';
import hbs from 'htmlbars-inline-precompile';

moduleForComponent('toggle-button', 'Integration | Component | toggle button', {
  integration: true,
});

test('toggle functionality', function(assert) {
  this.set('toggleTarget', {});

  this.render(hbs`{{toggle-button toggleTarget=toggleTarget toggleAttr="toggled"}}`);

  assert.equal(this.$('button').text().trim(), 'More options', 'renders default closedLabel');

  this.$('button').click();
  assert.equal(this.get('toggleTarget.toggled'), true, 'it toggles the attr on the target');
  assert.equal(this.$('button').text().trim(), 'Hide options', 'renders default openLabel');
  this.$('button').click();
  assert.equal(this.get('toggleTarget.toggled'), false, 'it toggles the attr on the target');

  this.set('closedLabel', 'Open the options!');
  this.set('openLabel', 'Close the options!');
  this.render(
    hbs`{{toggle-button toggleTarget=toggleTarget toggleAttr="toggled" closedLabel=closedLabel openLabel=openLabel}}`
  );

  assert.equal(this.$('button').text().trim(), 'Open the options!', 'renders passed closedLabel');
  this.$('button').click();
  assert.equal(this.$('button').text().trim(), 'Close the options!', 'renders passed openLabel');
});
