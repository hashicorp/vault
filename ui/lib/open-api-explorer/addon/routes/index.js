import Route from '@ember/routing/route';
import { inject as service } from '@ember/service';

export default Route.extend({
  flashMessages: service(),
  // without an empty model hook here, ember likes to use the parent model, and then things get weird with
  // query params, so here we're no-op'ing the model hook
  model() {},
  afterModel() {
    this.flashMessages.warning(
      'The "Try it out" functionality in this API explorer will communicate to this Vault server\'s endpoints with your current token. Your token will also be shown on the screen in the example curl command output.',
      {
        sticky: true,
      }
    );
  },
});
