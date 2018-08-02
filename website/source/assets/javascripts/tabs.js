document.addEventListener("turbolinks:load", function() {
  var tabs = document.querySelectorAll(".tabs li");

  function handleTabClick(clickEvent) {
    var clickedLink = clickEvent.currentTarget.querySelector("a");
    var activeContentId = clickedLink.getAttribute("data-tab-for");

    switchTab(activeContentId);

    clickEvent.preventDefault(activeContentId);
    return false;
  }

  function switchTab(id) {
    var tabsContent = document.querySelectorAll(".tabs-content");
    var activeTab = document.querySelector(`[data-tab-for="${id}"]`);
    var activeContent = document.getElementById(id);

    for (var i = 0; i < tabs.length; i++) {
      var tabLink = tabs[i].querySelector("a");
      tabLink.classList.remove("is-active");
    }

    for (i = 0; i < tabsContent.length; i++) {
      tabsContent[i].classList.remove("is-active");
    }

    activeTab.classList.add("is-active");
    activeContent.classList.add("is-active");
  }

  for (i = 0; i < tabs.length; i++) {
    tabs[i].addEventListener("click", handleTabClick)
  }

  var urlParams = new URLSearchParams(window.location.search);
  if (urlParams && urlParams.has("tab")) {
    switchTab(urlParams.get("tab"));
  }
});
