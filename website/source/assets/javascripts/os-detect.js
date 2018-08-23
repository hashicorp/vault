function getCurrentOS() {
  var userAgent = navigator.userAgent;
  if (userAgent.indexOf("Win") != -1) return "windows";
  if (userAgent.indexOf("Mac") != -1) return "darwin";
  if (userAgent.indexOf("Linux") != -1) return "linux";
  if (userAgent.indexOf("NetBSD") != -1) return "netbsd";
  if (userAgent.indexOf("FreeBSD") != -1) return "freebsd";
  if (userAgent.indexOf("OpenBSD") != -1) return "openbsd";
  if (userAgent.indexOf("SunOS") != -1) return "solaris";
  return "Unkown";
}

function getCurrentOSBit() {
  var userAgent = navigator.userAgent;
  if (userAgent.match( /(Win64|WOW64|Mac OS X 10|amd64|x86)/ )) {
    return "64-bit";
  }
  if (userAgent.match( /arm/ )) {
    return "arm";
  }
  return "32-bit";
}

document.addEventListener("turbolinks:load", function() {
  if (document.querySelector(`[data-os]`)) {
    var osSelector = '[data-os="' + getCurrentOS() + '"]';
    var bitSelector = '[data-os-bit="' + getCurrentOSBit() + '"]';
    var currentOSElement = document.querySelector(osSelector);
    var currentBitLinkElement = document.querySelector(osSelector + ' ' + bitSelector);
    var currentBitLinkHTML = currentBitLinkElement.cloneNode(true);
    var bitList = currentBitLinkElement.parentNode;

    // Move current Bit link to the start of the list
    bitList.removeChild(currentBitLinkElement);
    bitList.prepend(currentBitLinkHTML);

    // Highlight current OS and Bit link
    currentOSElement.classList.add("current");
    document.querySelector(osSelector + ' ' + bitSelector).classList.add("current")
  }
});
