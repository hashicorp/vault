export default function (server) {
  const handleResponse = (req, props) => {
    const xhr = req.passthrough();
    xhr.onreadystatechange = () => {
      if (xhr.readyState === 4 && xhr.status < 300) {
        // XMLHttpRequest response prop only has a getter -- redefine as writable and set value
        Object.defineProperty(xhr, 'response', {
          writable: true,
          value: JSON.stringify({
            ...JSON.parse(xhr.responseText),
            ...props,
          }),
        });
      }
    };
  };

  server.get('sys/seal-status', (schema, req) => {
    // randomly return one of the various states to test polling
    // 401 and 500 are stubs -- update with actual API values once determined
    const hcp_link_status = ['connected', 'disconnected', '401', '500'][Math.floor(Math.random() * 2)];
    return handleResponse(req, { hcp_link_status });
  });
  // enterprise only feature initially
  server.get('sys/health', (schema, req) => handleResponse(req, { version: '1.12.0-dev1+ent' }));
}
