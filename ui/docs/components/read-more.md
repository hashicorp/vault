# ReadMore

ReadMore components are used to wrap long text that wed like to show as one line initially with the option to expand and read.
Text which is shorter than the surrounding div will not truncate or show the See More button.

**Example**

```hbs preview-template
<div class='box linked-block-item'>
  <ReadMore>
    My super long text goes in here. Sed ut perspiciatis unde omnis iste natus error sit voluptatem
    accusantium doloremque laudantium, totam rem aperiam, eaque ipsa quae ab illo inventore veritatis et quasi
    architecto beatae vitae dicta sunt explicabo. Nemo enim ipsam voluptatem quia voluptas sit aspernatur aut
    odit aut fugit, sed quia consequuntur magni dolores eos qui ratione voluptatem sequi nesciunt. Neque porro
    quisquam est, qui dolorem ipsum quia dolor sit amet, consectetur, adipisci velit, sed quia non numquam
    eius modi tempora incidunt ut labore et dolore magnam aliquam quaerat voluptatem. Ut enim ad minima
    veniam, quis nostrum exercitationem ullam corporis suscipit laboriosam, nisi ut aliquid ex ea commodi
    consequatur? Quis autem vel eum iure reprehenderit qui in ea voluptate velit esse quam nihil molestiae
    consequatur, vel illum qui dolorem eum fugiat quo voluptas nulla pariatur?
  </ReadMore>
</div>

<div class='linked-block-item'>
  <ReadMore>
    Text that fits does not truncate
  </ReadMore>
</div>
```
