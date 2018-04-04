import Ember from 'ember';
import hbs from 'htmlbars-inline-precompile';

const { computed } = Ember;

export default Ember.Component.extend({
  layout: hbs`<a href="{{href-to 'vault.cluster' 'vault'}}" class={{class}}>
      {{#if hasBlock}}
        {{yield}}
      {{else}}
        {{text}}
      {{/if}}
      </a>
  `,

  tagName: '',

  text: computed(function() {
    return 'home';
  }),

  computedClasses: computed('classNames', function() {
    return this.get('classNames').join(' ');
  }),
});
