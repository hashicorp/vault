/* global self */
self.deprecationWorkflow = self.deprecationWorkflow || {};
self.deprecationWorkflow.config = {
  workflow: [
    // ivy-codemirror and ember-radio-button still use send-action
    { handler: 'silence', matchId: 'ember-component.send-action' },
    { handler: 'silence', matchId: 'ember-runtime.deprecate-copy-copyable' },
  ],
};
