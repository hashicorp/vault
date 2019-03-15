/* eslint-disable import/extensions */
import hbs from 'htmlbars-inline-precompile';
import { storiesOf } from '@storybook/ember';
import notes from './options.md';

storiesOf('AuthConfigForm/Options/', module)
  .addParameters({ options: { showPanel: true } })
  .add(
    `Options`,
    () => ({
      template: hbs`
        <h5 class="title is-5">Options</h5>
        {{auth-config-form/options}}
    `,
      context: {},
    }),
    { notes }
  );
