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
    var tabsContents = document.querySelectorAll(".tabs-content");
    var activeTab = document.querySelector('[data-tab-for="' + id + '"]');
    var activeContent = document.getElementById(id);

    tabs.forEach(function(tab) {
      var tabLink = tab.querySelector("a");
      tabLink.classList.remove("is-active");
    });

    tabsContents.forEach(function(tabsContent) {
      tabsContent.classList.remove("is-active");
    });

    activeTab.classList.add("is-active");
    activeContent.classList.add("is-active");
  }


  tabs.forEach(function(tab) {
    tab.addEventListener("click", handleTabClick);
  });

  var urlParams = new URLSearchParams(window.location.search);
  if (urlParams && urlParams.has("tab")) {
    switchTab(urlParams.get("tab"));
  }
});
