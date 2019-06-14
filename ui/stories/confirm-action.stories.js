/* eslint-disable import/extensions */
import hbs from 'htmlbars-inline-precompile';
import { storiesOf } from '@storybook/ember';
import { withKnobs, text, boolean } from '@storybook/addon-knobs';
import notes from './confirm-action.md';

storiesOf('ConfirmAction/', module)
  .addParameters({ options: { showPanel: true } })
  .addDecorator(
    withKnobs({
      escapeHTML: false,
    })
  )
  .add(
    `ConfirmAction`, () => ({
      template: hbs`
        <h5 class="title is-5">Confirm Action</h5>
        <ConfirmAction
          @buttonClasses={{buttonClasses}}
          @confirmTitle={{confirmTitle}}
          @confirmMessage={{confirmMessage}}
          @confirmButtonText={{confirmButtonText}}
          @cancelButtonText={{cancelButtonText}}
          @disabled={{disabled}}
          @onConfirmAction={{onComfirmAction}}
        >
          {{buttonText}}
        </ConfirmAction>
      `,
      context: {
        buttonText: text('buttonText', 'Delete'),
        buttonClasses: text('buttonClasses', 'button'),
        confirmTitle: text('confirmTitle', 'Delete this?'),
        confirmMessage: text('confirmMessage', 'You will not be able to recover it later.'),
        confirmButtonText: text('confirmButtonText', 'Delete'),
        cancelButtonText: text('cancelButtonText', 'Cancel'),
        disabled: boolean('disabled', false),
        onComfirmAction: () => {
          console.log('Action!');
        },
      },
    }),
    { notes }
  );
