import { Page } from "./api";
import { Pages } from "./api";
import { PagesContainer } from "./api";

import { Observable } from "./observables";
import { ObservableEventData } from "./observables";
import { ObserverInterface } from "./observables";
import { Subject } from "./observables";


import { Api } from "./api";

let api = new Api()

// Some of this is from there:
// https://visualstudiomagazine.com/articles/2013/04/01/introducing-practical-javascript.aspx


export interface PageVMElements {
    linkcontainer: HTMLElement,
    linklist: HTMLElement,
    titledesktop: HTMLElement,
    titlemobile: HTMLElement,
    mainpanel: HTMLElement,
    metatitle: HTMLElement,
    description: HTMLMetaElement,
    topbarmobile: HTMLElement,
    menubutton: HTMLElement,
}

export class PageVM
{
    constructor(
        private elements: PageVMElements,
        public buttonclicked: ObservableEventData = null,
        public linkcontainerclicked: ObservableEventData = null,
        public appdrawerelem1: AppDrawerElement = null,
        public appdrawerelem2: AppDrawerElement = null,
        public appdrawerSubject = new Subject(new AppDrawer),
        public obsPages = new Observable<Pages>(),
        public obsPage = new Observable<Page>(),
        ){

            // set up appdrawer
            let toggle: (this: AppDrawer)=>void = function(this:AppDrawer) {
                if (this.open) {
                    this.open = false
                } else {
                    this.open = true
                }
            }
            let close: (this: AppDrawer)=>void = function(this:AppDrawer) {
                this.open = false
            }

            this.appdrawerelem1 = new AppDrawerElement(this.elements.linkcontainer)
            this.appdrawerelem2 = new AppDrawerElement(this.elements.topbarmobile)

            this.buttonclicked = new ObservableEventData(elements.menubutton,"click").map((x: any) => 'toggle')
            this.appdrawerSubject.subscribe(this.appdrawerelem1)
            this.appdrawerSubject.subscribe(this.appdrawerelem2)
            this.buttonclicked.subscribe(this.appdrawerSubject)

            this.linkcontainerclicked = new ObservableEventData(elements.linkcontainer, "click").map((x: any) => 'close')
            this.linkcontainerclicked.subscribe(this.appdrawerSubject)


            // currently this is useless, since the links are already
            // filled in on the server
            let mylinklist = new Mylinklist(this.elements.linklist)
            this.obsPages.subscribe(mylinklist)

            let mypageview = new Mypageview(
                this.elements.mainpanel,
                this.elements.metatitle,
                this.elements.description
                )
            this.obsPage.subscribe(mypageview)


            this.fetchAllPages()
        }
    async fetchAllPages() { 
        this.pages = await api.listPages()
        // Clicking on the links *before* data has arrived should reload the
        // entire page. *After* data has arrived, attach the AJAX 'click' 
        // event handler
        this.attach_ajax_handlers()
    };
    set pages(pages: Pages) {
        this.obsPages.update(pages)
    };
    get pages(): Pages {
        let pages = this.obsPages.item
        if (pages.length === 0) {
            return null
        } 
        return this.obsPages.item;
    }
    set page(page: Page) {
        this.obsPage.update(page)
    }
    get page(): Page {
        return this.obsPage.item;
    }
    attach_ajax_handlers() {
        if (this.obsPages) {
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
            let p = new PagesContainer(this.obsPages.item)
            let page = p.findPageByKeyValue('prettyurl', a.pathname)

            this.page = page
            history.pushState(null, null, a.href);
            // TODO: what does the next line do?
            e.stopPropagation();
        }
    }
}
HTMLMetaElement

class Mypageview implements ObserverInterface {
    constructor(
        public content: HTMLElement, 
        public metatitle: HTMLElement,
        public description: HTMLMetaElement
        ) { }
    async next(page: Page) {
        let obj = await api.getPageContent(page.id)
        this.content.innerHTML = obj.body
        this.metatitle.innerHTML = page.metatitle
        this.description.content = page.description
    }
}


export class Mylinklist implements ObserverInterface {
    constructor(public e: HTMLElement) { }
    next(pages: Pages) {
        this.e.innerHTML = template(pages)
    }
}

function template(pages: Pages): string {
    return pages.map(
            (page) => 
            `<li><a class="xhr" href="${page.prettyurl}"> ${page.linkname}</a></li>`
        ).join('')
}

class AppDrawer {
    constructor(
        private _open: boolean = false
    ) { }
    get open() {
        return this._open
    }
    set open(bol: boolean) {
        this._open = bol;
    }
    next(s: string) {
        if (s === 'toggle') {
            if (this._open) {
                this._open = false
            } else {
                this._open = true
            }
        } else if (s === 'close') {
            this._open = false
        }  
    }
}

class AppDrawerObserver implements ObserverInterface {
    next(ad: AppDrawer): void {
        console.log(ad.open)
    }
}

class AppDrawerElement {
    constructor(
        public element: HTMLElement,
        ) { }
    get open() {
        return this.element.hasAttribute('open');
    }
    set open(open: Boolean) {
        if (open) {
            this.element.setAttribute('open', '');
        } else {
            this.element.removeAttribute('open');
        }
    }
    next(ad: AppDrawer) {
        if (ad.open) {
            this.open = true
        } else {
            this.open = false
        }
    }
}
