/* eslint-disable import/extensions */
import hbs from 'htmlbars-inline-precompile';
import { storiesOf } from '@storybook/ember';
import notes from './message.md';

storiesOf('Confirm/Message/', module)
  .addParameters({ options: { showPanel: true } })
  .add(
    `Message`,
    () => ({
      template: hbs`
      <h5 class="title is-5">Message</h5>
      <Message/>
    `,
      context: {},
    }),
    { notes }
  );
