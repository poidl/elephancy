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

export interface ObserverInterfaceNew {
    next: (property: any) => void
}

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
        public elements: Array<ObserverInterface> = []
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

export class Observer implements ObserverInterface {
    next: (property: any) => void
}

export class Subject extends Observable<Observer> {
    constructor(
        item: Observer = null,
        elements: Array<ObserverInterface> = []
    ) { super(item, elements) }
    next(myitem: any) {
        this.item.next(myitem)
        this.notify(this.item)
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
        this.notify(ev)
    }
    notify(ev: Event) {
        for (let e of this.elements) {
            e.next(this.mapf(ev))
        }
    }
    map = (f: (input: any) => any) =>  {
        this.mapf = f
        return this
    }
    mapf = (input: any) => input
}


export class AppDrawer {
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

export class AppDrawerObserver implements ObserverInterface {
    next(ad: AppDrawer): void {
        console.log(ad.open)
    }
}

export class AppDrawerElement {
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

