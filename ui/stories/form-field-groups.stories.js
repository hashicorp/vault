/* eslint-disable import/extensions */
import hbs from 'htmlbars-inline-precompile';
import { storiesOf } from '@storybook/ember';
import { withKnobs, select } from '@storybook/addon-knobs';
import notes from './form-field-groups.md';

// This will need to be replaced with a fake model, since the form fields associated with
// each model come from OpenApi and Storybook doesn't have a Vault server to call OpenApi from.
// Without OpenApi, not all of the models' form fields will show up in the Storybook UI.
const MODELS = {
  Approle: 'approle',
  AWS: 'aws/client',
  Azure: 'azure',
  Cert: 'cert',
  GCP: 'gcp',
  Github: 'github',
  JWT: 'jwt',
  Kubernetes: 'kubernetes',
  LDAP: 'ldap',
  OKTA: 'okta',
  Radius: 'radius',
  Userpass: 'userpass',
};

const DEFAULT_VALUE = 'aws/client';

storiesOf('Form/FormFieldGroups/', module)
  .addParameters({ options: { showPanel: true } })
  .addDecorator(withKnobs())
  .add(
    `FormFieldGroups`,
    () => ({
      template: hbs`
        <h5 class="title is-5">Form Field Groups</h5>
        <FormFieldGroups @model={{compute (action 'getModel' model)}} />
    `,
      context: {
        actions: {
          getModel(modelType) {
            return Ember.getOwner(this)
              .lookup('service:store')
              .createRecord(`auth-config/${modelType}`);
          },
        },
        model: select('model', MODELS, DEFAULT_VALUE),
      },
    }),
    { notes }
  );
