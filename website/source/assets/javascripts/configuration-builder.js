function downloadConfiguration() {
  var form = document.querySelector("#configuration-builder");
  var config = "";

  // Add Listener stanza
  if (document.getElementById("include_tcp_listener").checked) {
    config += `listener "tcp" {
${addFieldsToStanza("listener")}}
`;
  }

  // Add Storage stanza
  if (document.getElementById("include_storage").checked) {
    var backend = document.getElementById("storage").value;
    config += `
storage "${backend}" {
${addFieldsToStanza("storage")}}
`;
  }

  // Add Telemetry stanza
  if (document.getElementById("include_telemetry").checked) {
    var provider = document.getElementById("telemetry").value;
    config += `
telemetry {
${addFieldsToStanza("telemetry")}}
`;
  }

  // Add Seal stanza
  if (document.getElementById("include_seal").checked) {
    var type = document.getElementById("seal").value;
    config += `
seal "${type}" {
${addFieldsToStanza("seal")}}
`;
  }

  // Add UI stanza
  if (document.getElementById("include_ui").checked) {
    config += `
ui = true`;
    var startServerLink = document.querySelector(".start-server-link")
    startServerLink.href = `${startServerLink.href}?tab=ui`;
  }

  config = config.replace(/([^\r])\n/g, "$1\r\n");
  var blob = new Blob([config], {type: "text/plain;charset=utf-8"});
  saveAs(blob, "vault-config.hcl");
  document.querySelector(".form-actions").style.display = "none";
  document.querySelector("#download-confirm").style.display = "block";
  return false;
}

function addFieldsToStanza(stanza) {
  var fieldsets = document.querySelectorAll(`[data-config-stanza="${stanza}"] .nested-fields fieldset`);
  var lines = "";

  for (i = 0; i < fieldsets.length; i++) {
    var fieldset = fieldsets[i];
    if (fieldset.offsetWidth > 0 && fieldset.offsetHeight > 0) {
      var line = fieldsetToLine(fieldset);
      if (line) {
        lines += line;
      }
    }
  }
  return lines;
}

function fieldsetToLine(fieldset) {
  var parameter = fieldset.getAttribute("name");
  var isChecked = document.querySelector(`#include_${parameter}`).checked;
  if (isChecked) {
    var field = fieldset.querySelector(`#${parameter}`);
    var value = field.value;

    if (field.getAttribute("type") == "number") {
      return `  ${parameter} = ${value}
`;
    } else {
      return `  ${parameter} = "${value}"
`;
    }
  }
  return;
}

document.addEventListener("turbolinks:load", function() {
  var revealTriggers = document.querySelectorAll(".reveal-trigger");
  var configTriggers = document.querySelectorAll(".config-reveal-trigger");
  var configSelects = document.querySelectorAll(".config-reveal-select");

  for (i = 0; i < revealTriggers.length; i++) {
    revealTriggers[i].addEventListener("click", function(clickEvent) {
      var revealTrigger = clickEvent.currentTarget;
      revealTrigger.classList.toggle("active");
      revealTrigger.nextElementSibling.classList.toggle("active");
    });
  }

  for (i = 0; i < configTriggers.length; i++) {
    configTriggers[i].addEventListener("change", function(clickEvent) {
      var configTrigger = clickEvent.currentTarget;
      var container = configTrigger.closest("fieldset");
      var reveal = container.querySelector(".config-reveal-container");
      reveal.classList.toggle("active");

      if (reveal.querySelector(".config-reveal-select")) {
        var selection = reveal.querySelector(".config-reveal-select").value;
        document.querySelector(`[data-if-option="${selection}"]`).classList.toggle("active");
      }
    });
  }

  for (i = 0; i < configSelects.length; i++) {
    configSelects[i].addEventListener("change", function(clickEvent) {
      var configSelect = clickEvent.currentTarget;
      var selection = configSelect.value;
      var section = configSelect.closest("section");
      var reveal = section.querySelector(`[data-if-option='${selection}']`);
      var nestedOptions = section.querySelectorAll("[data-if-option]");

      for (i = 0; i < nestedOptions.length; i++) {
        nestedOptions[i].classList.remove("active");
      }

      if (reveal) {
        reveal.classList.add("active");
      }
    });
  }
});
