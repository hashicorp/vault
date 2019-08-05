import Component from '@ember/component';
import layout from '../../templates/components/confirm/message';

export default Component.extend({
  layout,
  tagName: '',
  renderedTrigger: null,
  id: null,
  onCancel() {},
  onConfirm() {},
  title: 'Delete this?',
  message: 'You will not be able to recover it later.',
  confirmButtonText: 'Delete',
  cancelButtonText: 'Cancel',
});
