import { reads } from '@ember/object/computed';
import Component from '@ember/component';

export default Component.extend({
  content: null,
  list: reads('content.keys'),
});
