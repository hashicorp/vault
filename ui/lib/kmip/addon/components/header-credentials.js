import Component from '@ember/component';
import { inject as service } from '@ember/service';
import layout from '../templates/components/header-credentials';

export default Component.extend({
  layout,
  tagName: '',
  secretMountPath: service(),
  scope: null,
  role: null,
});
