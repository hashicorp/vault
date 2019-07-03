import Route from '@ember/routing/route';
import { inject as service } from '@ember/service';

export default Route.extend({
  flashMessages: service(),
  // without an empty model hook here, ember likes to use the parent model, and then things get weird with
  // query params, so here we're no-op'ing the model hook
  model() {},
  afterModel() {
    let warning = `The "Try it out" functionality in this API explorer will make requests to this Vault server on your behalf.

IF YOUR TOKEN HAS THE PROPER CAPABILITIES, THIS WILL CREATE AND DELETE ITEMS ON THE VAULT SERVER.

Your token will also be shown on the screen in the example curl command output.`;
    this.flashMessages.warning(warning, {
      sticky: true,
      preformatted: true,
    });
  },
});
