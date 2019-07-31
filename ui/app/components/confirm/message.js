import Component from '@ember/component';
import { computed } from '@ember/object';

export default Component.extend({
  tagName: '',
  renderedTrigger: null,
  id: null,
  onCancel() {},
  onConfirm() {},
  title: 'Delete this?',
  message: 'You will not be able to recover it later.',
  confirmButtonText: 'Delete',
  cancelButtonText: 'Cancel',
  shouldYield: computed('id', 'renderedTrigger', function() {
    return this.id === this.renderedTrigger;
  }),
});
