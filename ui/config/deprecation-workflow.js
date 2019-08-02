/* global self */
self.deprecationWorkflow = self.deprecationWorkflow || {};
//self.deprecationWorkflow.config = {
//throwOnUnhandled: true
//}
self.deprecationWorkflow.config = {
  workflow: [
    // this seems to crop up in relation to pretender things
    { handler: 'silence', matchId: 'ember-polyfills.deprecate-merge' },
  ],
};
