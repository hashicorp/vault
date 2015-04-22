Ember.Application.initializer({
  name: 'load-steps',
  after: 'store',

  initialize: function(container, application) {
    var store = container.lookup('store:main');
    var steps = {
      "steps": [
        { id: 0, name: 'welcome', humanName: "Welcome to the Vault Interactive Demo!"},
        { id: 1, name: 'unseal', humanName: "Step 1: Unsealing your Vault"},
      ]
    };

    application.register('model:step', Demo.Step);

    store.pushPayload('step', steps);
  },
});
