import Ember from 'ember';

const { get, set } = Ember;

export default Ember.Component.extend({
  tagName: 'button',
  type: 'button',
  toggleTarget: null,
  toggleAttr: null,
  classNameBindings: ['buttonClass'],
  attributeBindings: ['type'],
  buttonClass: 'has-text-info',
  classNames: ['button', 'is-transparent'],
  openLabel: 'Hide options',
  closedLabel: 'More options',
  init() {
    this._super(...arguments);
    const toggleAttr = this.get('toggleAttr');
    Ember.defineProperty(
      this,
      'isOpen',
      Ember.computed(`toggleTarget.${toggleAttr}`, () => {
        const props = this.getProperties('toggleTarget', 'toggleAttr');
        return Ember.get(props.toggleTarget, props.toggleAttr);
      })
    );
  },
  click() {
    const target = this.get('toggleTarget');
    const attr = this.get('toggleAttr');
    const current = get(target, attr);
    set(target, attr, !current);
  },
});
