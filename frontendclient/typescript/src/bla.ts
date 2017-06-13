import { Api } from "./apinew";
import { Page } from "./apinew";
import { Pages } from "./apinew";
import { Link } from "./apinew";

import { Observable } from "./scriptnew";
import { ObservableString } from "./scriptnew";
import { ObservableEventData } from "./scriptnew";
import { ObservablePages } from "./scriptnew";
import { Myobserver } from "./scriptnew";
import { Mylinklist } from "./scriptnew";
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
        private page: Page = null,
        private pages: Pages = [],
        public string: string = '66',
        public observablestring = new ObservableString(),
        public observableeventdata: ObservableEventData = null,
        public observablepages: ObservablePages = null,
        ){
            let input = <HTMLInputElement>document.getElementById("fooinput")
            let paragr = document.getElementById("foop")

            this.observablestring.subscribe(new Myinput(input))
            this.observablestring.subscribe(new Myp(paragr))
            this.observablestring.update(this.string)

            this.observableeventdata = new ObservableEventData(input,"change")
            this.observableeventdata.subscribe(new Myp(paragr))

            this.observablepages = new ObservablePages(this.pages)
            let linklist = new Mylinklist(document.getElementById("linklist"))
            this.observablepages.subscribe(linklist)
            this.fetchAllPages()
        }
    async fetchAllPages() { 
        this.setPages = await api.listPages()
    };
    set setPages(pages: Pages) {
        this.observablepages.update(pages)
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

