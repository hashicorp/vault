/* eslint-disable import/extensions */
import hbs from 'htmlbars-inline-precompile';
import { storiesOf } from '@storybook/ember';
import notes from './auth-form.md';

storiesOf('AuthForm/', module)
  .addParameters({ options: { showPanel: false } })
  .add(
    `AuthForm`,
    () => ({
      template: hbs`
        <h5 class="title is-5">Auth Form</h5>
        <AuthForm />
    `,
      context: {},
    }),
    { notes }
  );
