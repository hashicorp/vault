import Component from '@ember/component';
import { computed } from '@ember/object';
import layout from '../../templates/components/confirm/trigger';

export default Component.extend({
  layout,
  tagName: '',
  renderedTrigger: null,
  id: null,
  onCancel() {},
  onConfirm() {},
  title: 'Delete this?',
  message: 'You will not be able to recover it later.',
  triggerText: null,
  confirmButtonText: 'Delete',
  cancelButtonText: 'Cancel',
  showConfirm: computed('renderedTrigger', function() {
    return !!this.renderedTrigger;
  }),
});
