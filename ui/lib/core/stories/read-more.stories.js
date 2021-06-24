import hbs from 'htmlbars-inline-precompile';
import { storiesOf } from '@storybook/ember';
import { withKnobs } from '@storybook/addon-knobs';

storiesOf('ReadMore', module)
  .addParameters({ options: { showPanel: true } })
  .addDecorator(withKnobs())
  .add(`ReadMoreLong`, () => ({
    template: hbs`
      <h5 class="title is-5">Read More</h5>
      <ReadMore>
      Lorem ipsum dolor sit amet, consectetur adipiscing elit. Nam volutpat vulputate lacus sit amet lobortis. Nulla fermentum porta consequat. Mauris porttitor libero nibh, ac facilisis ex molestie non. Nulla dolor est, pharetra et maximus vel, varius eu augue. Maecenas eget nisl convallis, vehicula massa quis, pharetra justo.
      </ReadMore>
    `,
    context: {},
  }))
  .add(`ReadMoreShort`, () => ({
    template: hbs`
      <h5 class="title is-5">Read More</h5>
      <ReadMore>
      Short Description here 
      </ReadMore>
    `,
    context: {},
  }));
