/* global self */
self.deprecationWorkflow = self.deprecationWorkflow || {};
//self.deprecationWorkflow.config = {
//throwOnUnhandled: true
//}
self.deprecationWorkflow.config = {
  // current output from deprecationWorkflow.flushDeprecations();
  // deprecations that will not be removed until 4.0.0 are filtered by deprecation-filter initializer rather than silencing below
  workflow: [
    { handler: 'log', matchId: 'routing.transition-methods' },
    { handler: 'log', matchId: 'implicit-injections' },
    { handler: 'log', matchId: 'ember-metal.get-with-default' },
    { handler: 'log', matchId: 'manager-capabilities.modifiers-3-13' },
    { handler: 'log', matchId: 'computed-property.override' },
    { handler: 'log', matchId: 'ember-glimmer.link-to.positional-arguments' },
    { handler: 'log', matchId: 'this-property-fallback' },
    { handler: 'log', matchId: 'ember-source.deprecation-without-for' },
    { handler: 'log', matchId: 'ember-source.deprecation-without-since' },
    { handler: 'log', matchId: 'ember-composable-helpers.contains-helper' },
    { handler: 'log', matchId: 'ember-component.is-visible' },
    { handler: 'log', matchId: 'ember-engines.deprecation-camelized-engine-names' },
    { handler: 'log', matchId: 'ember-engines.deprecation-router-service-from-host' },
  ],
};
