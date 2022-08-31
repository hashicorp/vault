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

  server.get('sys/seal-status', (schema, req) => handleResponse(req, { hcp_link_status: 'connected' }));
  // enterprise only feature initially
  server.get('sys/health', (schema, req) => handleResponse(req, { version: '1.12.0-dev1+ent' }));
}
