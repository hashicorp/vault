document.addEventListener("turbolinks:load", function() {
  var downloadLinks = document.querySelectorAll(".download-arches .download-link");

  for (i = 0; i < downloadLinks.length; i++) {
    downloadLinks[i].addEventListener("click", handleDownloadLinkClick);
  }

  function handleDownloadLinkClick(clickEvent) {
    var clickedLink = clickEvent.currentTarget;
    var bit = clickedLink.innerHTML;
    var container = clickedLink.closest(".download");
    var name = container.querySelector(".os-name").innerHTML;
    var icon = container.querySelector(".icon svg").outerHTML;
    var confirm = document.querySelector("#download-confirm");

    document.querySelector(".download-arches").style.display = "none";
    confirm.style.display = "flex";
    confirm.querySelector(".chosen-os-name").innerHTML = name;
    confirm.querySelector(".chosen-os-bit").innerHTML = bit;
    confirm.querySelector(".icon").innerHTML = icon;
  }
});
