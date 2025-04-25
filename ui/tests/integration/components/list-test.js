import { module, test } from 'qunit';
import { setupRenderingTest } from 'vault/tests/helpers';
import { render, click } from '@ember/test-helpers';
import { hbs } from 'ember-cli-htmlbars';
import { v4 as uuidv4 } from 'uuid';
import sinon from 'sinon';
import { setupMirage } from 'ember-cli-mirage/test-support';
import { overrideResponse } from 'vault/tests/helpers/stubs';

import { createSecretsEngine } from 'vault/tests/helpers/secret-engine/secret-engine-helpers';
import { SECRET_ENGINE_SELECTORS as SES } from 'vault/tests/helpers/secret-engine/secret-engine-selectors';
import { GENERAL } from 'vault/tests/helpers/general-selectors';

module('Integration | Component | secret-engine/list', function (hooks) {
  setupRenderingTest(hooks);
  setupMirage(hooks);

  // Incoming tests:
  // allows you to disable an engine
  // adds disabled css styling to unsupported secret engines
  // filters by name and engine type
  // applies overflow styling

  // Add these tests:
  // shows only displayable engines (whatever that means)
  // sorts the engines first by supported and then name
  //

  hooks.beforeEach(function () {
    this.server.post('/sys/capabilities-self', () => ({
      data: {
        capabilities: ['root'],
      },
    }));
    this.store = this.owner.lookup('service:store');
    this.version = this.owner.lookup('service:version');
    this.flashMessages = this.owner.lookup('service:flash-messages');
    this.flashMessages.registerTypes(['success', 'danger']);
    this.flashSuccessSpy = sinon.spy(this.flashMessages, 'success');
    this.flashDangerSpy = sinon.spy(this.flashMessages, 'danger');
    this.uid = uuidv4();
    // generate a model of cubbyhole, kv2, and nomad
    this.secretEngineModels = [
      createSecretsEngine(this.store, 'cubbyhole', 'cubbyhole-test'),
      createSecretsEngine(this.store, 'kv', 'kv2-test'),
      createSecretsEngine(this.store, 'nomad', 'nomad-test'),
    ];
  });

  test('it allows you to disable an engine', async function (assert) {
    const enginePath = 'kv2-test';
    this.server.delete(`sys/mounts/${enginePath}`, () => {
      assert.ok(true, 'Destroy record is called and deletes the engine');
      return overrideResponse(204);
    });
    await render(hbs`<SecretEngine::List @secretEngineModels={{this.secretEngineModels}} />`);

    assert.dom(SES.secretsBackendLink(enginePath)).exists('shows the link for the kvv2 secrets engine');
    const row = SES.secretsBackendLink(enginePath);
    await click(`${row} ${GENERAL.menuTrigger}`);
    await click(`${row} ${GENERAL.confirmTrigger}`);
    await click(GENERAL.confirmButton);

    assert.true(
      this.flashSuccessSpy.calledWith(`The kv Secrets Engine at ${enginePath}/ has been disabled.`),
      'Flash message shows that engine was disabled.'
    );
  });
});
