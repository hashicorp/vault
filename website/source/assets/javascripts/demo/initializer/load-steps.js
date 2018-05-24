Ember.Application.initializer({
  name: 'load-steps',
  after: 'store',

  initialize: function(container, application) {
    var store = container.lookup('store:main');
    var steps = {
      "steps": [
        { id: 0, name: 'welcome', humanName: "Welcome to the Vault Interactive Tutorial!"},
        { id: 1, name: 'steps', humanName: "Step 1: Overview"},
        { id: 2, name: 'init', humanName: "Step 2: Initialize your Vault"},
        { id: 3, name: 'unseal', humanName: "Step 3: Unsealing your Vault"},
        { id: 4, name: 'auth', humanName: "Step 4: Authorize your requests"},
        { id: 5, name: 'list', humanName: "Step 5: List available secret engines"},
        { id: 6, name: 'secrets', humanName: "Step 6: Read and write secrets"},
        { id: 7, name: 'update', humanName: "Step 7: Update the secret data"},
        { id: 8, name: 'patch', humanName: "Step 8: Update the data without overwriting"},
        { id: 9, name: 'versions', humanName: "Step 9: Work with different data versions"},
        { id: 10, name: 'delete', humanName: "Step 10: Delete the data"},
        { id: 11, name: 'recover', humanName: "Step 11: Recover the deleted data"},
        { id: 12, name: 'destroy', humanName: "Step 12: Permanently delete data"},
        { id: 13, name: 'help', humanName: "Step 13: Get Help"},
        { id: 14, name: 'seal', humanName: "Step 14: Seal your Vault"},
        { id: 15, name: 'finish', humanName: "You're finished!"},
      ]
    };

    application.register('model:step', Demo.Step);

    store.pushPayload('step', steps);
  },
});
