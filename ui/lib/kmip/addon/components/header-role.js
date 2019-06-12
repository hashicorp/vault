import Component from '@ember/component';
import { inject as service } from '@ember/service';
import layout from '../templates/components/header-role';

export default Component.extend({
  layout,
  tagName: '',
  secretMountPath: service(),
  scope: null,
});
