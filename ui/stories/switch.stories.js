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
        @id={{id}}
        @name={{name}}
        @isChecked={{isChecked}}
        @disabled={{disabled}}
        @size={{size}}
        @round={{round}}
      >
        {{yielded}}
      </Switch>
    `,
      context: {
        id: text('id', 'my-switch'),
        name: text('name', 'my-checkbox'),
        yielded: text('yield', 'Inner content here'),
        isChecked: boolean('isChecked', false),
        disabled: boolean('disabled', false),
        size: select('size', ['small', 'medium', 'large'], 'small'),
        round: boolean('round', true),
        onChange(key, value) {
          console.log(`${key} =  ${value}`);
        },
      },
    }),
    { notes }
  );
