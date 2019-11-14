/* global self */
self.deprecationWorkflow = self.deprecationWorkflow || {};
//self.deprecationWorkflow.config = {
//throwOnUnhandled: true
//}
self.deprecationWorkflow.config = {
  workflow: [
    // after ED 3.9 this shouldn't be necessary
    { handler: 'silence', matchId: 'deprecate-fetch-ember-data-support' },
  ],
};
