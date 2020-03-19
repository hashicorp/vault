import hbs from 'htmlbars-inline-precompile';
import { storiesOf } from '@storybook/ember';
import notes from './switch.md';
import { withKnobs, object, text, boolean, select } from '@storybook/addon-knobs';

storiesOf('Switch/', module)
  .addParameters({ options: { showPanel: true } })
  .addDecorator(withKnobs())
  .add(
    `Switch`,
    () => ({
      template: hbs`
      <h5 class="title is-5">Switch</h5>
      <Switch
        @inputId={{inputId}}
        @isChecked={{isChecked}}
        @onChange={{onChange}}
        @disabled={{disabled}}
        @size={{size}}
        @status={{status}}
        @round={{round}}
      >
        {{yielded}}
      </Switch>
    `,
      context: {
        inputId: text('id', 'my-switch'),
        name: text('name', 'my-checkbox'),
        yielded: text('yield', 'Label content here ✔️'),
        isChecked: boolean('isChecked', true),
        disabled: boolean('disabled', false),
        size: select('size', ['small', 'medium', 'large'], 'small'),
        status: select('status', ['normal', 'success'], 'normal'),
        round: boolean('round', true),
        onChange() {
          this.set('isChecked', !this.isChecked);
        },
      },
    }),
    { notes }
  );
