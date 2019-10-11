
import hbs from 'htmlbars-inline-precompile';
import { storiesOf } from '@storybook/ember';
import { withKnobs, text } from '@storybook/addon-knobs';

import notes from './file-to-array-buffer.md';

storiesOf('FileToArrayBuffer/', module)
  .addParameters({ options: { showPanel: true } })
  .addDecorator(
    withKnobs()
  )
  .add(`FileToArrayBuffer`, () => ({
    template: hbs`
      <h5 class="title is-5">File To Array Buffer</h5>
      <FileToArrayBuffer @onChange={{this.onChange}} @label={{this.label}}
      @fileHelpText={{this.fileHelpText}} />
      {{#if this.fileName}}
        {{this.fileName}} as bytes: {{this.fileBytes}}
      {{/if}}
    `,
    context: {
      onChange(file, name) {
        console.log(`${name} contents as an ArrayBuffer:`, file);
      },
      label: text('Label'),
      fileHelpText: text('Help text'),
    },
  }),
  {notes}
);
