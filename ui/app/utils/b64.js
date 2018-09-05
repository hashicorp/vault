import base64js from 'base64-js';

export function encodeString(string) {
  var encoded = new TextEncoderLite('utf-8').encode(string);
  return base64js.fromByteArray(encoded);
}

export function decodeString(b64String) {
  var uint8array = base64js.toByteArray(b64String);
  return new TextDecoderLite('utf-8').decode(uint8array);
}
