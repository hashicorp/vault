import Route from '@ember/routing/route';

export default Route.extend({
  // without an empty model hook here, ember likes to use the parent model, and then things get weird with
  // query params, so here we're no-op'ing the model hook
  model() {},
});
