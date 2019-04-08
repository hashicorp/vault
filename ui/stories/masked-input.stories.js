/* eslint-disable import/extensions */
import hbs from 'htmlbars-inline-precompile';
import { storiesOf } from '@storybook/ember';
import { withKnobs, text, boolean } from '@storybook/addon-knobs';
import notes from './masked-input.md';

storiesOf('MaskedInput/', module)
  .addParameters({ options: { showPanel: true } })
  .addDecorator(withKnobs())
  .add(
    `MaskedInput`,
    () => ({
      template: hbs`
        <h5 class="title is-5">Masked Input</h5>
        <MaskedInput
          @value={{value}}
          @placeholder={{placeholder}}
          @allowCopy={{allowCopy}}
          @displayOnly={{displayOnly}}
        />
    `,
      context: {
        value: text('value', ''),
        placeholder: text('placeholder', 'super-secret'),
        allowCopy: boolean('allowCopy', false),
        displayOnly: boolean('displayOnly', false),
      },
    }),
    { notes }
  );
