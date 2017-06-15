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

// Observables:
// http://www.anasfirdousi.com/understanding-observable-patterns-behind-observables-rxjs-rx.html
// https://netbasal.com/javascript-observables-under-the-hood-2423f760584
export interface ObserverConstructor {
    new (e: any): ObserverInterface
}

export interface ObserverInterface {
    next(property: any): void
}

// function createObserver(ctor: ObserverConstructor, e: any): ObserverInterface {
//     return new ctor(e);
// }

export class Myinput implements ObserverInterface {
    constructor(public e: HTMLInputElement) { }
    next(s: string) {
        this.e.value = s
    }
}

export class Myp implements ObserverInterface {
    constructor(public e: HTMLElement) { }
    next(s: string) {
        this.e.innerHTML = s
    }
}

export class Mylinklist implements ObserverInterface {
    constructor(public e: HTMLElement) { }
    next(pages: Pages) {
        this.e.innerHTML = template(pages)
    }
}

export class Mypageview implements ObserverInterface {
    constructor(
        public content: HTMLElement, 
        public metatitle: HTMLElement
        ) { }
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

export class Observable<T> {
    constructor(
        public item: T = null,
        private elements: Array<ObserverInterface> = []
    ) { }
    subscribe(element: ObserverInterface) {
        this.elements.push(element)
    }
    update = (item: T) => {
        this.item = item
        this.notify(item)
    }
    notify(item: T) {
        for (let e of this.elements) {
            e.next(item)
        }
    }
}

export class ObservableEventData {
    constructor(
        private eventtarget: EventTarget = null,
        private eventtype: string = null,
        private elements: Array<ObserverInterface> = []
    ) { 
        this.eventtarget.addEventListener(this.eventtype, this.update, false)
     }
    subscribe(element: ObserverInterface) {
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
