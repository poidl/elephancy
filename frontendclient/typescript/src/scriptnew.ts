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

// export function merge<T,S>(ot: Observable<T>, os: Observable<S>) {

// }

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
    // get item(): T {
    //     return this._item
    // }
    // set item(item: T) {
    //     this._item = item
    //     this.notify(item)
    // }
    // next(myitem: any) {
        
    // }
    notify(item: T) {
        for (let e of this.elements) {
            e.next(item)
        }
    }

}

export class Observer implements ObserverInterface {
    next: (property: any) => void
    // next(s: string) {
    //     this.e.value = s
    // }
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
// export class Subject {
//     constructor(
//         public item: Observer = null,
//         public elements: Array<ObserverInterface> = []
//     ) { }
//     subscribe(element: ObserverInterface) {
//         this.elements.push(element)
//     }
//     update = (item: Observer) => {
//         this.item = item
//         this.notify(item)
//     }
//     // get item(): T {
//     //     return this._item
//     // }
//     // set item(item: T) {
//     //     this._item = item
//     //     this.notify(item)
//     // }
//     next(myitem: any) {
//         this.item.next(myitem)
//         this.notify(this.item)
//     }
//     notify(item: Observer) {
//         for (let e of this.elements) {
//             e.next(item)
//         }
//     }

// }



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
        // switch (ev.type) {
        //     case "change": {
        //         let eventdata = (<HTMLInputElement>ev.target).value
        //         this.notify_string(eventdata)
        //     }
        //     case "click": {
        //         this.notify()
        //     }
        // }

        
    }
    // notify_string(s: string) {
    //     for (let e of this.elements) {
    //         e.next(s)
    //     }
    // }
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

// export class ObservableEventDatan extends ObservableEventData {
//     constructor(
//         eventtarget: EventTarget = null,
//         eventtype: string = null,
//         elements: Array<ObserverInterface> = []
//     ) { 
//         super(eventtarget, eventtype, elements)
//      }
//     update = (ev: Event) => {
//         switch (ev.type) {
//             case "click": {
//                 if ev.target
//                 this.notify()
//             }
//         }
        
//     }
// }

// export class AppDrawerGa implements ObserverInterface {
//     constructor(
//         private open: boolean = false,
//         public next: (property: any) => void
//     ) {}
// }


export class AppDrawer {
    constructor(
        // public next: (property: any) => void
        private _open: boolean = false
    ) { }
    // accessors
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
    // next: (property: any) => void
    // settoggle() {
    //     let toggle: (this: AppDrawer)=>void = function(this:AppDrawer) {
    //             if (this.open) {
    //                 this.open = false
    //             } else {
    //                 this.open = true
    //             }
    //         }
    //     this.next = toggle 
    // }
    // next() {
    //     if (this.open) {
    //         this.open = false
    //     } else {
    //         this.open = true
    //     }
    // }
}
export class AppDrawerObserver implements ObserverInterface {
    // constructor(public e: HTMLInputElement) { }
    next(ad: AppDrawer): void {
        console.log(ad.open)
    }
}

// export class AppDrawerElement {
//     constructor(
//         public element: HTMLElement,
//         // public next: (property: any) => void
//         ) { }
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
//     // next: (property: any) => void
//     // settoggle() {
//     //     let toggle: (this: AppDrawer)=>void = function(this:AppDrawer) {
//     //             if (this.open) {
//     //                 this.open = false
//     //             } else {
//     //                 this.open = true
//     //             }
//     //         }
//     //     this.next = toggle 
//     // }
//     // next() {
//     //     if (this.open) {
//     //         this.open = false
//     //     } else {
//     //         this.open = true
//     //     }
//     // }
// }

// // export class AppDrawer implements ObserverInterface {
// //     constructor(
// //         public element: HTMLElement,
// //         ) { }
// //     // accessors
// //     get open() {
// //         return this.element.hasAttribute('open');
// //     }

// //     set open(open: Boolean) {
// //         // Reflect the value of the open property as an HTML attribute.
// //         if (open) {
// //             this.element.setAttribute('open', '');
// //         } else {
// //             this.element.removeAttribute('open');
// //         }
// //     }
// //     next() {
// //         if (this.open) {
// //             this.open = false
// //         } else {
// //             this.open = true
// //         }
// //     }
// // }

// export class AppDrawerCloser extends AppDrawer {
//     constructor(element: HTMLElement) {
//         super(element)
//     }
//     next() {
//         this.open = false
//     }
// }



export class AppDrawerNew implements ObserverInterface {
    constructor(
        public element: HTMLElement, 
        public next: (property: any) => void
        ) {}
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
}




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
