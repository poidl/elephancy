import { Api } from "./api";
import { Page } from "./api";
import { Link } from "./api";
// import { PagesContainer } from "./apinew";

import * as req from "./myrequest";

let api = new Api()

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

// TODO: why define a class if I already have the interface? So I can use it as 
// a type in the definition of Subject. Does this make sense?
export class Observer implements ObserverInterface {
    next: (property: any) => void
}

// a Subject is an observable for a type which 'next' can be called on (i.e. an 
// observer)
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

