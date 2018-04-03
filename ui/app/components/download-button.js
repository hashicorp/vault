import Ember from 'ember';
import hbs from 'htmlbars-inline-precompile';

const { computed } = Ember;

export default Ember.Component.extend({
  layout: hbs`{{actionText}}`,
  tagName: 'a',
  role: 'button',
  attributeBindings: ['role', 'download', 'href'],
  download: computed('filename', 'extension', function() {
    return `${this.get('filename')}-${new Date().toISOString()}.${this.get('extension')}`;
  }),

  href: computed('data', 'mime', 'stringify', function() {
    let data = this.get('data');
    const mime = this.get('mime');
    if (this.get('stringify')) {
      data = JSON.stringify(data, null, 2);
    }

    const file = new File([data], { type: mime });
    return window.URL.createObjectURL(file);
  }),

  actionText: 'Download',
  data: null,
  filename: null,
  mime: 'text/plain',
  extension: 'txt',
  stringify: false,
});
