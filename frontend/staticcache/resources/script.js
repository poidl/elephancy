var titlel = document.getElementById('titlelink');
titlel.addEventListener('click', goToRoot, false);

var titlebar = document.getElementById('bartitlelink');
titlebar.addEventListener('click', goToRoot, false);

var container = document.getElementById('leftDrawer');
container.addEventListener('click', navlink_clicked, false);

var pagesmap = {};

var xmlhttp = new XMLHttpRequest();
var url = "/json/pages.json";
xmlhttp.onreadystatechange = function() {
    if (xmlhttp.readyState == 4 && xmlhttp.status == 200) {
        pagesmap = JSON.parse(xmlhttp.responseText);
    }
};
xmlhttp.open("GET", url, true);
xmlhttp.send();

function mapmap(pmap, arr_in, key_out) {
  // arr is [key_in:value_in]. search which of the maps in pmap contains obj, and retrieve value of key2 in that specific map.
  // TODO: return error in case it doesn't find anything
  if (arr_in.length != 2) {
    // TODO: error
  } else {
    key_in =arr_in[0]
    value_in = "/"+arr_in[1].split('/').pop()
  }
  for (var page in pmap) {
    if (pmap[page][key_in] === value_in) {
      return pmap[page][key_out]
    }
  }
  return ""
}

function ajax (e) {
  if ((e.target != e.currentTarget) & (e.target.className === "xhr")) {
    contenturl = mapmap(pagesmap,["Urlpath",e.target.href],"ContentUrl")
    metatitle = mapmap(pagesmap,["Urlpath",e.target.href],"Metatitle")
    e.preventDefault();
    if (e.target.href != location.href) {
      swapMainwindow(contenturl);
      swapTitle(metatitle);
      history.pushState(null, null, e.target.href);
    }
    e.stopPropagation();
  }
}

function goToRoot (e) {
  ajax(e);
}

function navlink_clicked (e) {
  ajax(e);
  navdrawer_toggle(); // close navdrawer after click, in case it is open
}

function swapTitle(metatitle) {
  document.getElementById("metatitle").innerHTML = metatitle;
}

function swapMainwindow(contenturl) {
  var xhttp = new XMLHttpRequest();
  xhttp.onreadystatechange = function() {
    if (xhttp.readyState == 4 && xhttp.status == 200) {
      document.getElementById("mainPanel").innerHTML = xhttp.responseText;
    }
  }

  xhttp.open("GET", "/"+ contenturl, true); // true: asynchronously
  xhttp.setRequestHeader("myheader",'XMLHttpRequest');
  xhttp.send();
}

function navdrawer_close() {
  appbarElement.classList.remove('open');
  navdrawerContainer.classList.remove('open');
}

function navdrawer_toggle() {
  var isOpen = navdrawerContainer.classList.contains('open');
  if(isOpen) {
    navdrawer_close();
  } else {
    appbarElement.classList.add('open');
    navdrawerContainer.classList.add('open');
  }
}

// load document and close navdrawer after click
function loadnav(filename) {
  window.loadDoc(filename);
  window.navdrawer_close();
}

// event listener for app-bar button, to toggle navdrawer on click
var navdrawerContainer = document.querySelector('.navdrawer');
var appbarElement = document.querySelector('.app-bar');

var menuBtn = document.querySelector('.menu');
menuBtn.addEventListener('click', function() {
  window.navdrawer_toggle();
}, true);

window.onload = function() {
  window.addEventListener("popstate", doit, false);

  function doit(e) {
    contenturl = mapmap(pagesmap,["Urlpath",location.href],"ContentUrl")
    metatitle = mapmap(pagesmap,["Urlpath",location.href],"Metatitle")
    swapMainwindow(contenturl);
    swapTitle(metatitle);
  }
}
