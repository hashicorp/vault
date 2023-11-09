# DownloadButton

DownloadButton wraps an `<Hds::Button>` to perform a download action. [HDS button args](https://helios.hashicorp.design/components/button?tab=code)

- NOTE: when using in an engine, remember to add the download service to its dependencies (in /engine.js) and map to it in /app.js
  [ember-docs](https://ember-engines.com/docs/services)

| Param          | Type                  | Default               | Description                                                                                                |
| -------------- | --------------------- | --------------------- | ---------------------------------------------------------------------------------------------------------- |
| [filename]     | <code>string</code>   |                       | name of file that prefixes the ISO timestamp generated at download                                         |
| [data]         | <code>string</code>   |                       | data to download                                                                                           |
| [fetchData]    | <code>function</code> |                       | function that fetches data and returns download content                                                    |
| [extension]    | <code>string</code>   | <code>txt</code>      | file extension, the download service uses this to determine the mimetype                                   |
| [stringify]    | <code>boolean</code>  | <code>false</code>    | argument to stringify the data before passing to the File constructor                                      |
| [onSuccess]    | <code>callback</code> |                       | callback from parent to invoke if download is successful                                                   |
| [hideIcon]     | <code>boolean</code>  | <code>false</code>    | renders the 'download' icon by default, pass true to hide (ex: when download button appears in a dropdown) |
| [text]         | <code>string</code>   | <code>Download</code> | button text, defaults to 'Download'                                                                        |
| [color]        | <code>string</code>   |                       | HDS default is primary, but there are four color options: primary, secondary, tertiary, and critical.      |
| [iconPosition] | <code>string</code>   | <code>leading</code>  | icon position, 'leading' (HDS default) or 'trailing'                                                       |
| [isIconOnly]   | <code>boolean</code>  |                       | button only renders an icon, no text                                                                       |

**Example**

```hbs preview-template
<DownloadButton @text='Download this stuff' @color='secondary' />
```
