import { Api } from "./apinew";
import { Page } from "./apinew";
import { Pages } from "./apinew";
import { Link } from "./apinew";
// import { PagesContainer } from "./apinew";

import * as req from "./myrequest";

let host = '127.0.0.1';
let port = 8080;
let basepath = '/api';

let api = new Api()


let pages: Pages
// let pagescontainer: PagesContainer




export class Myobserver {
    constructor(public e: any){}
    next(property: any) {};
}

export class Myinput extends Myobserver {
    constructor(public e: HTMLInputElement) {
        super(e)
    }
    next(s: string) {
        this.e.value = s
    }
}

export class Myp extends Myobserver {
    constructor(public e: HTMLElement) {
        super(e)
    }
    next(s: string) {
        this.e.innerHTML = s
    }
}

export class Mylinklist extends Myobserver {
    constructor(public e: HTMLElement) {
        super(e)
    }
    next(pages: Pages) {
        this.e.innerHTML = template(pages)
    }
}

export class Mypageview extends Myobserver {
    constructor(
        public content: HTMLElement, 
        public metatitle: HTMLElement
        ) {
        super(content)
    }
    async next(page: Page) {
        let obj = await api.getPageContent(page.id)
        this.content.innerHTML = obj.body
        this.metatitle.innerHTML = page.metatitle
    }
}

function template(pages: Pages): string {
    return pages.map(
            (page) => 
            `<li><a class="xhr" href="${page.prettyurl}"> ${page.linkname}</a></li>`
        ).join('')
}

export class Observable {
    subscribe(subscriber: any) {}
    update = (element: any) => {}
    // notify(element: any) { }
}

export class ObservableString extends Observable {
    constructor(
        private s: string = null,
        private elements: Array<Myobserver> = []
    ) { super() }
    subscribe(element: Myobserver) {
        this.elements.push(element)
    }
    update = (s: string) => {
        this.s = s
        this.notify(s)
    }
    notify(s: string) {
        for (let e of this.elements) {
            e.next(s)
        }
    }
}

export class ObservableEventData extends Observable {
    constructor(
        private eventtarget: EventTarget = null,
        private eventtype: string = null,
        private elements: Array<Myobserver> = []
    ) { 
        super()
        this.eventtarget.addEventListener(this.eventtype, this.update, false)
     }
    subscribe(element: Myobserver) {
        this.elements.push(element)
    }
    update = (ev: Event) => {
        switch (ev.type) {
            case "change": {
                let eventdata = (<HTMLInputElement>ev.target).value
                this.notify(eventdata)
            }
        }
        
    }
    notify(s: string) {
        for (let e of this.elements) {
            e.next(s)
        }
    }
}

export class ObservablePages extends Observable {
    constructor(
        public pages: Pages,
        private observers: Array<Myobserver> = []
    ) { 
        super()
     }
    subscribe(observer: Myobserver) {
        this.observers.push(observer)
    }
    update = (pages: Pages) => {
        this.pages = pages
        this.notify(pages)
    }
    notify(pages: Pages) {
        for (let o of this.observers) {
            o.next(pages)
        }
    }
}

export class ObservablePage extends Observable {
    constructor(
        public page: Page,
        private observers: Array<Myobserver> = []
    ) { 
        super()
     }
    subscribe(observer: Myobserver) {
        this.observers.push(observer)
    }
    update = (page: Page) => {
        this.page = page
        this.notify(page)
    }
    notify(page: Page) {
        for (let o of this.observers) {
            o.next(page)
        }
    }
}

// class AppDrawer {
//     private element: Element;
//     constructor(id: string) {
//         this.element = document.querySelector(id);

//         // listener on element?
//         // this.element.addEventListener('click', e => {
//         //     do something
//         // });
//     }
//     // accessors
//     get open() {
//         return this.element.hasAttribute('open');
//     }

//     set open(open: Boolean) {
//         // Reflect the value of the open property as an HTML attribute.
//         if (open) {
//             this.element.setAttribute('open', '');
//         } else {
//             this.element.removeAttribute('open');
//         }
//     }
//     toggleDrawer() {
//         if (this.open) {
//             this.open = false
//         } else {
//             this.open = true
//         }
//     }
// }

// let appdrawer = new AppDrawer('.linkcontainer')
// let appbar = new AppDrawer('.top-bar-mobile');

// let menuBtn = document.querySelector('.menubutton');
// menuBtn.addEventListener('click', function () {
//     appdrawer.toggleDrawer();
//     appbar.toggleDrawer();
// }, true);

// linkcontainer.addEventListener('click', close, false)

// function close(): void {
//     appdrawer.open = false
//     appbar.open = false
// }

// window.onload = function () {
//     window.addEventListener("popstate", doit, false);

//     function doit() {
//         let p = pagescontainer
//         let page = p.findPageByKeyValue('prettyurl', '/' + location.href.split('/').pop())
//         if (!page) {
//             let err = new Error('Error in popstate event handler')
//             throw err
//         }
//         update(page)
//     }
// }
