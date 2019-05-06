/* eslint-disable import/extensions */
import hbs from 'htmlbars-inline-precompile';
import { storiesOf } from '@storybook/ember';
import notes from './message-error.md';

let model = Ember.Object.create({
  adapterError: {
    message: 'This is an adapterError on the model'
  },
  isError: true
});

storiesOf('MessageError/', module)
  .addParameters({ options: { showPanel: true } })
  .add(`MessageError`, () => ({
    template: hbs`
      <h5 class="title is-5">Message Error</h5>
      <MessageError @model={{model}} />
    `,
    context: {
      model
    }
  }),
  {notes}
);
