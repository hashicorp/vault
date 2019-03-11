/* eslint-disable import/extensions */
import hbs from 'htmlbars-inline-precompile';
import { storiesOf } from '@storybook/ember';
import { withKnobs, text, boolean } from '@storybook/addon-knobs';
import notes from './confirm-action.md';

storiesOf('ConfirmAction/', module)
  .addDecorator(
    withKnobs({
      escapeHTML: false,
    })
  )
  .add(
    `ConfirmAction`,
    () => ({
      template: hbs`
        <h5 class="title is-5">Confirm Action</h5>
        <ConfirmAction
          @onConfirmAction={{onComfirmAction}}
          @confirmButtonText={{confirmButtonText}}
          @confirmMessage={{confirmMessage}}
          @cancelButtonText={{cancelButtonText}}
          @disabled={{disabled}}
          >
          Delete
        </ConfirmAction>
    `,
      context: {
        onComfirmAction: () => {
          console.log('Action!');
        },
        confirmButtonText: text('confirmButtonText', 'Yes'),
        confirmMessage: text('confirmMessage', 'Are you sure you want to do this?'),
        cancelButtonText: text('cancelButtonText', 'Cancel'),
        disabled: boolean('disabled', false),
      },
    }),
    { notes }
  );
