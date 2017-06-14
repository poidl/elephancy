import { Api } from "./apinew";
import { Page } from "./apinew";
import { Pages } from "./apinew";
import { PagesContainer } from "./apinew";
import { Link } from "./apinew";

import { Observable } from "./scriptnew";
import { ObservableString } from "./scriptnew";
import { ObservableEventData } from "./scriptnew";
import { ObservablePages } from "./scriptnew";
import { ObservablePage } from "./scriptnew";
import { Myobserver } from "./scriptnew";
import { Mylinklist } from "./scriptnew";
import { Mypageview } from "./scriptnew";
import { Myinput } from "./scriptnew";
import { Myp } from "./scriptnew";

import * as req from "./myrequest";

let host = '127.0.0.1';
let port = 8080;
let basepath = '/api';

let api = new Api()

export interface IPageVM 
{
  fetchAllPages(): void;
//   fetchPage(id: number): void;
//   getPage(): Page;
//   setPage(page: Page): void;
//   setPages();
}

export class PageVM implements IPageVM 
{
    constructor(
        public string: string = '66',
        public observablestring = new ObservableString(),
        public observableeventdata: ObservableEventData = null,
        public observablepages: ObservablePages = new ObservablePages([]),
        public observablepage: ObservablePage = new ObservablePage(null),
        private linkcontainer = <HTMLElement>document.querySelector('.linkcontainer'),
        private linklist = <HTMLElement>document.getElementById('linklist'),
        private titledesktop = <HTMLElement>document.querySelector('.title-desktop'),
        private titlemobile = <HTMLElement>document.querySelector('.title-mobile'),
        private mainpanel = <HTMLElement>document.getElementById("mainPanel"),
        private metatitle = <HTMLElement>document.getElementById("metatitle")
        ){
            let input = <HTMLInputElement>document.getElementById("fooinput")
            let paragr = document.getElementById("foop")

            this.observablestring.subscribe(new Myinput(input))
            this.observablestring.subscribe(new Myp(paragr))
            this.observablestring.update(this.string)

            this.observableeventdata = new ObservableEventData(input,"change")
            this.observableeventdata.subscribe(new Myp(paragr))

            // this.observablepages = new ObservablePages()
            let mylinklist = new Mylinklist(this.linklist)
            this.observablepages.subscribe(mylinklist)
            let mypageview = new Mypageview(this.mainpanel, this.metatitle)
            this.observablepage.subscribe(mypageview)

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
        this.observablepages.update(pages)
    };
    get getPages(): Pages {
        let pages = this.observablepages.pages
        if (pages.length === 0) {
            return null
        } 
        return this.observablepages.pages;
    }
    attach_ajax_handlers() {
        console.log(this.getPages)
        if (this.getPages) {
            this.linkcontainer.addEventListener('click', this.ajax, false);
            this.titledesktop.addEventListener('click', this.ajax, false);
            this.titlemobile.addEventListener('click', this.ajax, false);
        } else  {
            console.log('there was error attaching the handlers to left drawer: pages not initialized in view model');
        }
    }
    ajax = (e: MouseEvent) => {
        let a = (<HTMLAnchorElement>e.target)
        if (a.className === 'xhr') {
            e.preventDefault();
            let p = new PagesContainer(this.observablepages.pages)
            let page = p.findPageByKeyValue('prettyurl', a.pathname)

            this.setPage = page
            history.pushState(null, null, a.href);
            // TODO: what does the next line do?
            e.stopPropagation();
        }
    }
    set setPage(page: Page) {
        this.observablepage.update(page)
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

