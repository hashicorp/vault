import Ember from 'ember';
import hbs from 'htmlbars-inline-precompile';

const { computed, get, set } = Ember;

export default Ember.Component.extend({
  // passed from outside
  onChange: null,
  wrapResponse: true,

  ttl: null,

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
    Ember.assert(
      '`onChange` handler is a required attr in `' + this.toString() + '`.',
      get(this, 'onChange')
    );
  },

  layout: hbs`
    <div class="field">
      <div class="b-checkbox">
        <input
          id="wrap-response"
          class="styled"
          name="wrap-response"
          type="checkbox"
          checked={{wrapResponse}}
          onchange={{action 'changedValue' 'wrapResponse'}}
          />
        <label for="wrap-response" class="is-label">
          Wrap response
        </label>
      </div>
      {{#if wrapResponse}}
        {{ttl-picker data-test-wrap-ttl-picker=true labelText='Wrap TTL' onChange=(action (mut ttl))}}
      {{/if}}
    </div>
  `,

  actions: {
    changedValue(key, event) {
      const { type, value, checked } = event.target;
      const val = type === 'checkbox' ? checked : value;
      set(this, key, val);
      get(this, 'onChange')(get(this, 'wrapTTL'));
    },
  },
});
