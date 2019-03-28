/* eslint-disable import/extensions */
import hbs from 'htmlbars-inline-precompile';
import { storiesOf } from '@storybook/ember';
import notes from './config.md';

storiesOf('AuthConfigForm/Config/', module)
  .addParameters({ options: { showPanel: false } })
  .add(
    `Config`,
    () => ({
      template: hbs`
        <h5 class="title is-5">Config</h5>
        {{auth-config-form/config}}
    `,
      context: {},
    }),
    { notes }
  );
