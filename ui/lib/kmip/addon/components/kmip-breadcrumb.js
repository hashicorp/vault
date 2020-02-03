import Component from '@ember/component';
import { inject as service } from '@ember/service';
import layout from '../templates/components/kmip-breadcrumb';
import { or } from '@ember/object/computed';

export default Component.extend({
  layout,
  tagName: '',
  secretMountPath: service(),
  shouldShowPath: or('showPath', 'scope', 'role'),
  showPath: false,
  path: null,
  scope: null,
  role: null,
});
