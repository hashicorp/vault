import hbs from 'htmlbars-inline-precompile';
import { storiesOf } from '@storybook/ember';
import notes from './confirm.md';
import { withKnobs, text } from '@storybook/addon-knobs';

storiesOf('Confirm/Confirm', module)
  .addParameters({ options: { showPanel: true } })
  .addDecorator(
    withKnobs({
      escapeHTML: false,
    })
  )
  .add(
    `Confirm`,
    () => ({
      template: hbs`
      <h5 class="title is-5">Confirm</h5>
      <div class="popup-menu-content">
        <div class="box">
          <Confirm as |c|>
            <c.Message
              @id={{id}}
              @title={{title}}
              @triggerText={{triggerText}}
              @message={{message}}
              @confirmButtonText={{confirmButtonText}}
              @cancelButtonText={{cancelButtonText}}
              @onConfirm={{onConfirm}}
              />
          </Confirm>
        </div>
      </div>
    `,
      context: {
        id: 'foo',
        onConfirm: () => {
          alert('Confirmed!');
        },
        title: text('title', 'Delete this?'),
        message: text('message', 'You will not be able to recover it later.'),
        confirmButtonText: text('confirmButtonText', 'Delete'),
        cancelButtonText: text('cancelButtonText', 'Cancel'),
        triggerText: text('triggerText', 'Delete'),
      },
    }),
    { notes }
  );
