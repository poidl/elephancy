import { Api } from "./api";
import { Page } from "./api";
import { Pages } from "./api";
import { Link } from "./api";
import { PagesContainer } from "./api";

import * as req from "./myrequest";

let host = '127.0.0.1';
let port = 8080;
let basepath = '/api';

let api = new Api()

let pages: Pages

// attach listener to the entire drawer
var container = document.getElementById('leftDrawer');

async function attach_handlers() {
    try {
        pages = await api.listPages()
        container.addEventListener('click', navlink_clicked, false);
    }
    catch (e) {
        console.log('there was error attaching the handlers to left drawer');
        console.log(e);
    }
}

attach_handlers()


function navlink_clicked(e: MouseEvent) {
    console.log('navlink_clicked')
    ajax(e);
    //   na vdrawer_toggle(); // close navdrawer after click, in case it is open
}

async function ajax(e: MouseEvent) {
    let a = (<HTMLAnchorElement>e.target)
    if (a.className === 'xhr') {
        e.preventDefault();
        let p = new PagesContainer(pages)
        let page = p.findPageByKeyValue('prettyurl', a.pathname)

        // get page content
        let link = page.getLinkByRel('self')
        let options = {
            protocol: 'http:',
            hostname: host,
            port: port,
            path: link,
            headers: { myheader: 'XMLHttpRequest' }
        };
        // No need to try/catch here since these throw errors if they fail.
        // TODO: make sense?
        let incoming = await req.make_request(options)
        let obj = await req.get_code_body(incoming)

        document.getElementById("mainPanel").innerHTML = obj.body
        document.getElementById("metatitle").innerHTML = page.metatitle;
        history.pushState(null, null, a.href);
        // TODO: what does the next line do?
        e.stopPropagation();
    }
}

function navdrawer_close() {
  appbarElement.classList.remove('open');
  navdrawerContainer.classList.remove('open');
}

let navdrawerContainer = document.querySelector('.navdrawer');
let appbarElement = document.querySelector('.app-bar');

function navdrawer_toggle() {
    console.log('toggeling')
  let isOpen = navdrawerContainer.classList.contains('open');
  if(isOpen) {
    navdrawer_close();
  } else {
    appbarElement.classList.add('open');
    navdrawerContainer.classList.add('open');
  }
}

let menuBtn = document.querySelector('.menu');
menuBtn.addEventListener('click', function() {
  navdrawer_toggle();
}, true);





