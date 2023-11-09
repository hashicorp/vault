
# TextFile
TextFile components render a file upload input with the option to toggle a Enter as text button
 that changes the input into a textarea

| Param | Type | Default | Description |
| --- | --- | --- | --- |
| onChange | <code>function</code> |  | Callback function to call when the value of the input changes, returns an object in the shape of { value: fileContents, filename: 'some-file.txt' } |
| [uploadOnly] | <code>boolean</code> | <code>false</code> | When true, renders a static file upload input and removes the option to toggle and input plain text |
| [helpText] | <code>string</code> |  | Text underneath label. |
| [label] | <code>string</code> | <code>&quot;&#x27;File&#x27;&quot;</code> | Text to use as the label for the file input. If none, default of 'File' is rendered |

**Example**  
```hbs preview-template
<TextFile @uploadOnly={{true}} @helpText="help text" @onChange={{this.handleChange}} @label="PEM Bundle" />
```
