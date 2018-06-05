import Ember from 'ember';

export default Ember.Mixin.create({
  // selector passed to `this.$()` to find the element to focus
  // defaults to `'input'`
  focusOnInsertSelector: null,
  shouldFocus: true,

  // uses Ember.on so that we don't have to worry about calling _super if
  // didInsertElement is overridden
  focusOnInsert: Ember.on('didInsertElement', function() {
    Ember.run.schedule('afterRender', this, 'focusOnInsertFocus');
  }),

  focusOnInsertFocus() {
    if (this.get('shouldFocus') === false) {
      return;
    }
    this.forceFocus();
  },

  forceFocus() {
    var $selector = this.$(this.get('focusOnInsertSelector') || 'input').first();
    if (!$selector.is(':focus')) {
      $selector.focus();
    }
  },
});
