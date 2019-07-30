/* global self */
self.deprecationWorkflow = self.deprecationWorkflow || {};
self.deprecationWorkflow.config = {
  workflow: [
    { handler: 'throw', matchId: 'ember-component.send-action' },
    { handler: 'throw', matchId: 'transition-state' },
    { handler: 'throw', matchId: 'ember-polyfills.deprecate-merge' },
    { handler: 'throw', matchId: 'events.inherited-function-listeners' },
    { handler: 'throw', matchId: 'ember-runtime.deprecate-copy-copyable' },
  ],
};
