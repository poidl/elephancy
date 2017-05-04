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

var linkcontainer = document.querySelector('.linkcontainer');

async function attach_ajax_handlers() {
    try {
        pages = await api.listPages()
        // Clicking on the links *before* data has arrived should reload the
        // entire page. *After* data has arrived, attach the AJAX 'click' 
        // event handler
        linkcontainer.addEventListener('click', fetch_content, false);
    }
    catch (e) {
        console.log('there was error attaching the handlers to left drawer');
        console.log(e);
    }
}

attach_ajax_handlers()


function fetch_content(e: MouseEvent) {
    ajax(e);
    //   na vdrawer_toggle(); // close linkcontainer after click, in case it is open
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

class AppDrawer  {
    private element: Element;
    constructor(id: string) {
        this.element = document.querySelector(id);

        // listener on element?
        // this.element.addEventListener('click', e => {
        //     do something
        // });
    }
    // accessors
    get open() {
        return this.element.hasAttribute('open');
    }

    set open(open: Boolean) {
        // Reflect the value of the open property as an HTML attribute.
        if (open) {
            this.element.setAttribute('open', '');
        } else {
            this.element.removeAttribute('open');
        }
    }
    toggleDrawer() {
        if (this.open) {
            this.open = false
        } else {
            this.open = true
        }
    }
}

let appdrawer = new AppDrawer('.linkcontainer')
let appbar = new AppDrawer('.app-bar');

let menuBtn = document.querySelector('.menu');
menuBtn.addEventListener('click', function () {
    appdrawer.toggleDrawer();
    appbar.toggleDrawer();
}, true);

linkcontainer.addEventListener('click', close, false)

function close(): void {
    appdrawer.open = false
    appbar.open = false
}