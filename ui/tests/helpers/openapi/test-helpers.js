import authModelAttributes from './auth-model-attributes';
import secretModelAttributes from './secret-model-attributes';

export const secretEngineHelper = (test, secretEngine) => {
  const engineData = secretModelAttributes[secretEngine];
  if (!engineData)
    throw new Error(`No engine attributes found in secret-model-attributes for ${secretEngine}`);

  const modelNames = Object.keys(engineData);
  // A given secret engine might have multiple models that are openApi driven
  modelNames.forEach((modelName) => {
    test(`${modelName} model getProps returns correct attributes`, async function (assert) {
      const model = this.store.createRecord(modelName, {});
      const helpUrl = model.getHelpUrl(this.backend);
      const result = await this.pathHelp.getProps(helpUrl, this.backend);
      const expected = engineData[modelName];
      assert.deepEqual(result, expected, `getProps returns expected attributes for ${modelName}`);
    });
  });
};

export const authEngineHelper = (test, authBackend) => {
  const authData = authModelAttributes[authBackend];
  if (!authData) throw new Error(`No auth attributes found in auth-model-attributes for ${authBackend}`);

  const itemNames = Object.keys(authData);
  itemNames.forEach((itemName) => {
    if (itemName.startsWith('auth-config/')) {
      // Config test doesn't need to instantiate a new model
      test(`${itemName} model`, async function (assert) {
        const model = this.store.createRecord(itemName, {});
        const helpUrl = model.getHelpUrl(this.mount);
        const result = await this.pathHelp.getProps(helpUrl, this.mount);
        const expected = authData[itemName];
        assert.deepEqual(result, expected, `getProps returns expected attributes for ${itemName}`);
      });
    } else {
      test.skip(`generated-${itemName}-${authBackend} model`, async function (assert) {
        const modelName = `generated-${itemName}-${authBackend}`;
        // Generated items need to instantiate the model first via getNewModel
        await this.pathHelp.getNewModel(modelName, this.mount, `auth/${this.mount}/`, itemName);
        const model = this.store.createRecord(modelName, {});
        // Generated items don't have this method -- helpUrl is calculated in path-help.js line 101
        const helpUrl = model.getHelpUrl(this.mount);
        const result = await this.pathHelp.getProps(helpUrl, this.mount);
        const expected = authData[modelName];
        assert.deepEqual(result, expected, `getProps returns expected attributes for ${modelName}`);
      });
    }
  });
};
