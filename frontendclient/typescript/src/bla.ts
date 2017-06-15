import { Api } from "./apinew";
import { Page } from "./apinew";
import { Pages } from "./apinew";
import { PagesContainer } from "./apinew";
import { Link } from "./apinew";

import { Observable } from "./scriptnew";
// import { ObservableString } from "./scriptnew";
import { Observablegen } from "./scriptnew";
import { ObservableEventData } from "./scriptnew";
// import { ObservablePages } from "./scriptnew";
// import { ObservablePage } from "./scriptnew";
// import { Myobserver } from "./scriptnew";
import { Mylinklist } from "./scriptnew";
import { Mypageview } from "./scriptnew";
import { Myinput } from "./scriptnew";
import { Myp } from "./scriptnew";

import * as req from "./myrequest";

let host = '127.0.0.1';
let port = 8080;
let basepath = '/api';

let api = new Api()

// Some of this is from there:
// https://visualstudiomagazine.com/articles/2013/04/01/introducing-practical-javascript.aspx

export interface IPageVM 
{
  fetchAllPages(): void;
//   fetchPage(id: number): void;
//   getPage(): Page;
//   setPage(page: Page): void;
//   setPages();
}

export interface Myelements {
    linkcontainer: HTMLElement,
    linklist: HTMLElement,
    titledesktop: HTMLElement,
    titlemobile: HTMLElement,
    mainpanel: HTMLElement,
    metatitle: HTMLElement,
    input: HTMLInputElement,
    paragr: HTMLElement,
}

let me: Myelements = {
    linkcontainer: <HTMLElement>document.querySelector('.linkcontainer'),
    linklist: <HTMLElement>document.getElementById('linklist'),
    titledesktop: <HTMLElement>document.querySelector('.title-desktop'),
    titlemobile: <HTMLElement>document.querySelector('.title-mobile'),
    mainpanel: <HTMLElement>document.getElementById("mainPanel"),
    metatitle: <HTMLElement>document.getElementById("metatitle"),
    input: <HTMLInputElement>document.getElementById("fooinput"),
    paragr: document.getElementById("foop")
}

export class PageVM implements IPageVM 
{
    constructor(
        public string = new Observablegen<string>(),
        public eventdata: ObservableEventData = null,
        public pages = new Observablegen<Pages>(),
        public page = new Observablegen<Page>(),
        private elements = me,
        ){

            this.string.subscribe(new Myinput(elements.input))
            this.string.subscribe(new Myp(elements.paragr))
            this.string.update('55')
            this.eventdata = new ObservableEventData(elements.input,"change")
            this.eventdata.subscribe(new Myp(elements.paragr))

            // currently this is useless, since the links are already
            // filled in on the server
            let mylinklist = new Mylinklist(this.elements.linklist)
            this.pages.subscribe(mylinklist)

            let mypageview = new Mypageview(this.elements.mainpanel, this.elements.metatitle)
            this.page.subscribe(mypageview)

            this.fetchAllPages()
        }
    async fetchAllPages() { 
        this.setPages = await api.listPages()
        // Clicking on the links *before* data has arrived should reload the
        // entire page. *After* data has arrived, attach the AJAX 'click' 
        // event handler
        this.attach_ajax_handlers()
    };
    set setPages(pages: Pages) {
        this.pages.update(pages)
    };
    get getPages(): Pages {
        let pages = this.pages.item
        if (pages.length === 0) {
            return null
        } 
        return this.pages.item;
    }
    attach_ajax_handlers() {
        if (this.getPages) {
            this.elements.linkcontainer.addEventListener('click', this.ajax, false);
            this.elements.titledesktop.addEventListener('click', this.ajax, false);
            this.elements.titlemobile.addEventListener('click', this.ajax, false);
        } else  {
            console.log('there was error attaching the handlers to left drawer: pages not initialized in view model');
        }
    }
    ajax = (e: MouseEvent) => {
        let a = (<HTMLAnchorElement>e.target)
        if (a.className === 'xhr') {
            e.preventDefault();
            let p = new PagesContainer(this.pages.item)
            let page = p.findPageByKeyValue('prettyurl', a.pathname)

            this.setPage = page
            history.pushState(null, null, a.href);
            // TODO: what does the next line do?
            e.stopPropagation();
        }
    }
    set setPage(page: Page) {
        this.page.update(page)
    }

    // fetchPage(id: number) { };
    // getPage() { };
    // getPages() { };
    // setPage(page: Page) { 
}

let vm = new PageVM()

// let objinput = new Binder(<HTMLInputElement>document.getElementById("fooinput"), vm.number);
// vm.string = '999'
// let obj = new Binder(document.getElementById("foo"), vm.getPages);
// vm.fetchAllPages()


window.onload = function () {
    window.addEventListener("popstate", doit, false);

    function doit() {
        let p = new PagesContainer(vm.pages.item)
        let page = p.findPageByKeyValue('prettyurl', '/' + location.href.split('/').pop())
        if (!page) {
            let err = new Error('Error in popstate event handler')
            throw err
        }
        vm.setPage = page
    }
}