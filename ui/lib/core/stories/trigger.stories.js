/* eslint-disable import/extensions */
import hbs from 'htmlbars-inline-precompile';
import { storiesOf } from '@storybook/ember';
import notes from './trigger.md';

storiesOf('Confirm/Trigger/', module)
  .addParameters({ options: { showPanel: true } })
  .add(
    `Trigger`,
    () => ({
      template: hbs`
      <h5 class="title is-5">Trigger</h5>
      <Trigger/>
    `,
      context: {},
    }),
    { notes }
  );
