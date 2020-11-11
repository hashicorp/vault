import { assert } from '@ember/debug';
import Component from '@ember/component';
import { set, get, computed } from '@ember/object';
import hbs from 'htmlbars-inline-precompile';

export default Component.extend({
  // passed from outside
  onChange: null,
  wrapResponse: true,

  ttl: '30m',

  wrapTTL: computed('wrapResponse', 'ttl', function() {
    const { wrapResponse, ttl } = this.getProperties('wrapResponse', 'ttl');
    return wrapResponse ? ttl : null;
  }),

  didRender() {
    this._super(...arguments);
    get(this, 'onChange')(get(this, 'wrapTTL'));
  },

  init() {
    this._super(...arguments);
    assert('`onChange` handler is a required attr in `' + this.toString() + '`.', get(this, 'onChange'));
  },

  layout: hbs`
    <div class="field">
      {{ttl-picker2
        data-test-wrap-ttl-picker=true
        label='Wrap response'
        helperTextDisabled='Will not wrap response'
        helperTextEnabled='Will wrap response with a lease of'
        enableTTL=wrapResponse
        initialValue=ttl
        onChange=(action 'changedValue')
      }}
    </div>
  `,

  actions: {
    changedValue(ttlObj) {
      set(this, 'wrapResponse', ttlObj.enabled);
      set(this, 'ttl', `${ttlObj.seconds}s`);
      get(this, 'onChange')(get(this, 'wrapTTL'));
    },
  },
});
