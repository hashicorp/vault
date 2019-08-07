/* eslint-disable import/extensions */
import hbs from 'htmlbars-inline-precompile';
import { storiesOf } from '@storybook/ember';
import notes from './confirm.md';
import { withKnobs, text } from '@storybook/addon-knobs';

storiesOf('Confirm/', module)
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
      <div class="box">
      <Confirm as |c|>
        <nav class="menu">
          <ul class="menu-list">
            <li class="action">
              <c.Trigger>
                <button
                  type="button"
                  class="link is-destroy"
                  onclick={{action c.onTrigger id}}>
                  Delete
                </button>
              </c.Trigger>
            </li>
          </ul>
        </nav>
        <c.Message
          @id={{item.id}}
          @onCancel={{action c.onCancel}}
          @onConfirm={{onConfirm}}
          @title={{title}}
          @message={{message}}
          @confirmButtonText={{confirmButtonText}}
          @cancelButtonText={{cancelButtonText}}>
        </c.Message>
      </Confirm>
      </div>
    `,
      context: {
        id: 'foo',
        onCancel: () => {
          alert('Cancelled!');
        },
        onConfirm: () => {
          alert('Confirmed!');
        },
        title: text('title', 'Are you sure?'),
        message: text('message', 'You will not be able to undo this action.'),
        confirmButtonText: text('confirmButtonText', 'Confirm'),
        cancelButtonText: text('cancelButtonText', 'Cancel'),
      },
    }),
    { notes }
  );
