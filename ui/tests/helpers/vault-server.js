export default class VaultServerTestHelper {
  isOpen = false;

  constructor() {
    this.socket = new WebSocket('ws://127.0.0.1:9201');
    this.socket.addEventListener('error', console.error); // eslint-disable-line
    this.socket.addEventListener('open', () => {
      this.isOpen = true;
    });
  }

  restart() {
    // do some polling to ensure connection is open
    if (this.pollRestart) {
      clearTimeout(this.pollRestart);
    }
    return new Promise((resolve) => {
      const sendMessage = () => {
        if (!this.isOpen) {
          this.pollRestart = setTimeout(sendMessage, 500);
        } else {
          this.socket.send('restart vault');
          setTimeout(resolve, 500);
        }
      };
      sendMessage();
    });
  }
}
