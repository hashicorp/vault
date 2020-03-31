import { module, test } from 'qunit';
import { setupRenderingTest } from 'ember-qunit';
import { render, click, fillIn } from '@ember/test-helpers';
import sinon from 'sinon';
import hbs from 'htmlbars-inline-precompile';

module('Integration | Component | ttl-picker2', function(hooks) {
  setupRenderingTest(hooks);

  test('it renders time and unit inputs when TTL enabled', async function(assert) {
    let changeSpy = sinon.spy();
    this.set('onChange', changeSpy);
    await render(hbs`
      <TtlPicker2
        @onChange={{onChange}}
        @enableTTL={{true}}
      />
    `);

    assert.dom('[data-test-ttl-value]').exists('TTL Picker time input exists');
    assert.dom('[data-test-ttl-unit]').exists('TTL Picker unit select exists');
  });

  test('it does not show time and unit inputs when TTL disabled', async function(assert) {
    let changeSpy = sinon.spy();
    this.set('onChange', changeSpy);
    await render(hbs`
      <TtlPicker2
        @onChange={{onChange}}
        @enableTTL={{false}}
      />
    `);
    assert.dom('[data-test-ttl-value]').doesNotExist('TTL Picker time input exists');
    assert.dom('[data-test-ttl-unit]').doesNotExist('TTL Picker unit select exists');
  });

  test('it passes the appropriate data to onChange when toggled on', async function(assert) {
    let changeSpy = sinon.spy();
    this.set('onChange', changeSpy);
    await render(hbs`
      <TtlPicker2
        @label="clicktest"
        @unit="m"
        @time="10"
        @onChange={{onChange}}
        @enableTTL={{false}}
      />
    `);
    await click('[data-test-toggle-input="clicktest"]');
    assert.ok(changeSpy.calledOnce, 'it calls the passed onChange');
    assert.ok(
      changeSpy.calledWith({
        enabled: true,
        seconds: 600,
        timeString: '10m',
      }),
      'Passes the default values back to onChange'
    );
  });

  test('it keeps seconds value when unit is changed', async function(assert) {
    let changeSpy = sinon.spy();
    this.set('onChange', changeSpy);
    await render(hbs`
      <TtlPicker2
        @label="clicktest"
        @unit="s"
        @time="360"
        @onChange={{onChange}}
        @enableTTL={{false}}
      />
    `);
    await click('[data-test-toggle-input="clicktest"]');
    assert.ok(changeSpy.calledOnce, 'it calls the passed onChange');
    assert.ok(
      changeSpy.calledWith({
        enabled: true,
        seconds: 360,
        timeString: '360s',
      }),
      'Passes the default values back to onChange'
    );
    await fillIn('[data-test-select="ttl-unit"]', 'm');
    assert.ok(
      changeSpy.calledWith({
        enabled: true,
        seconds: 360,
        timeString: '6m',
      }),
      'Units and time update without changing seconds value'
    );
  });

  test('it recalculates seconds when unit is changed and recalculateSeconds is on', async function(assert) {
    let changeSpy = sinon.spy();
    this.set('onChange', changeSpy);
    await render(hbs`
      <TtlPicker2
        @label="clicktest"
        @unit="s"
        @time="120"
        @onChange={{onChange}}
        @enableTTL={{true}}
        @recalculateSeconds={{true}}
      />
    `);
    await fillIn('[data-test-select="ttl-unit"]', 'm');
    assert.ok(
      changeSpy.calledWith({
        enabled: true,
        seconds: 7200,
        timeString: '120m',
      }),
      'Seconds value is recalculated based on time and unit'
    );
  });
});
