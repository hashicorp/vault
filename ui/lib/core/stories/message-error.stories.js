import hbs from 'htmlbars-inline-precompile';
import { storiesOf } from '@storybook/ember';
import notes from './message-error.md';
import EmberObject from '@ember/object';

let model = EmberObject.create({
  adapterError: {
    message: 'This is an adapterError on the model',
  },
  isError: true,
});

storiesOf('MessageError/', module)
  .addParameters({ options: { showPanel: true } })
  .add(
    `MessageError`,
    () => ({
      template: hbs`
      <h5 class="title is-5">Message Error</h5>
      <MessageError @model={{model}} />
    `,
      context: {
        model,
      },
    }),
    { notes }
  );
