<!--THIS FILE IS AUTO GENERATED. This file is generated from JSDoc comments in app/components/file-to-array-buffer.js. To make changes, first edit that file and run "yarn gen-story-md file-to-array-buffer" to re-generate the content.-->

## FileToArrayBuffer
`FileToArrayBuffer` is a component that will allow you to pick a file from the local file system. Once
loaded, this file will be emitted as a JS ArrayBuffer to the passed `onChange` callback.

**Params**

| Param | Type | Default | Description |
| --- | --- | --- | --- |
| onChange | <code>function</code> | <code></code> | The function to call when the file read is complete. This function recieves the file as a JS ArrayBuffer |
| [label] | <code>String</code> | <code></code> | Text to use as the label for the file input |
| fileHelpText | <code>String</code> | <code></code> | Text to use as help under the file input |

**Example**
  
```js
  <FileToArrayBuffer @onChange={{action (mut file)}} />
```

**See**

- [Uses of FileToArrayBuffer](https://github.com/hashicorp/vault/search?l=Handlebars&q=FileToArrayBuffer+OR+file-to-array-buffer)
- [FileToArrayBuffer Source Code](https://github.com/hashicorp/vault/blob/master/ui/app/components/file-to-array-buffer.js)

---
