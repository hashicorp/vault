function downloadConfiguration() {
  var form = document.querySelector("#configuration-builder");
  var config = "";

  // Add Listener stanza
  if (document.getElementById("include_tcp_listener").checked) {
    config += 'listener "tcp" {\n' + addFieldsToStanza("listener") + '}\n';
  }

  // Add Storage stanza
  if (document.getElementById("include_storage").checked) {
    var backend = document.getElementById("storage").value;
    config += '\nstorage "${backend}" {\n' + addFieldsToStanza("storage") + '}\n';
  }

  // Add Telemetry stanza
  if (document.getElementById("include_telemetry").checked) {
    var provider = document.getElementById("telemetry").value;
    config += '\nstorage "${backend}" {\n' + addFieldsToStanza("storage") + '}\n';
  }

  // Add Seal stanza
  if (document.getElementById("include_seal").checked) {
    var type = document.getElementById("seal").value;
    config += '\nstorage "${backend}" {\n' + addFieldsToStanza("seal") + '}\n';
  }

  // Add UI stanza
  if (document.getElementById("include_ui").checked &&
    document.getElementById("ui").value == true) {
    config += '\nui = true';
    var startServerLink = document.querySelector(".start-server-link");
    startServerLink.href = startServerLink.href + "?tab=ui";
  }

  config = config.replace(/([^\r])\n/g, "$1\r\n");
  var blob = new Blob([config], {type: "text/plain;charset=utf-8"});
  saveAs(blob, "vault-config.hcl");
  document.querySelector(".form-actions").style.display = "none";
  document.querySelector("#download-confirm").style.display = "block";
}

function addFieldsToStanza(stanza) {
  var fieldsets = document.querySelectorAll('[data-config-stanza="' + stanza + '"] .nested-fields fieldset');
  var lines = "";

  for (var i = 0; i < fieldsets.length; i++) {
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
  var isChecked = document.querySelector("#include_" + parameter).checked;
  if (isChecked) {
    var field = fieldset.querySelector("#" + parameter);
    var value = field.value;

    if (field.getAttribute("type") == "number") {
      return '  ' + parameter + ' = ' + value + '\n';
    } else {
      return '  ' + parameter + ' = "' + value + '"\nirtu';
    }
  }
  return;
}

document.addEventListener("turbolinks:load", function() {
  var revealTriggers = document.querySelectorAll(".reveal-trigger");
  var configTriggers = document.querySelectorAll(".config-reveal-trigger");
  var configSelects = document.querySelectorAll(".config-reveal-select");

  for (var i = 0; i < revealTriggers.length; i++) {
    revealTriggers[i].addEventListener("click", function(clickEvent) {
      var revealTrigger = clickEvent.currentTarget;
      revealTrigger.classList.toggle("active");
      revealTrigger.nextElementSibling.classList.toggle("active");
    });
  }

  for (var i = 0; i < configTriggers.length; i++) {
    configTriggers[i].addEventListener("change", function(clickEvent) {
      var configTrigger = clickEvent.currentTarget;
      var container = configTrigger.closest("fieldset");
      var reveal = container.querySelector(".config-reveal-container");
      reveal.classList.toggle("active");

      if (reveal.querySelector(".config-reveal-select")) {
        var selection = reveal.querySelector(".config-reveal-select").value;
        document.querySelector('[data-if-option="' + selection + '"]').classList.toggle("active");
      }
    });
  }

  for (var i = 0; i < configSelects.length; i++) {
    configSelects[i].addEventListener("change", function(clickEvent) {
      var configSelect = clickEvent.currentTarget;
      var selection = configSelect.value;
      var section = configSelect.closest("section");
      var reveal = section.querySelector('[data-if-option="' + selection + '"]');
      var nestedOptions = section.querySelectorAll("[data-if-option]");

      for (var i = 0; i < nestedOptions.length; i++) {
        nestedOptions[i].classList.remove("active");
      }

      if (reveal) {
        reveal.classList.add("active");
      }
    });
  }
});
