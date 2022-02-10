import { isPresent, isVisible, text } from 'ember-cli-page-object';

export default {
  title: text('[data-test-component=json-editor-title]'),
  hasToolbar: isPresent('[data-test-component=json-editor-toolbar]'),
  hasJSONEditor: isPresent('[data-test-component=json-editor]'),
  canEdit: isVisible('div.CodeMirror-gutters'),
};
