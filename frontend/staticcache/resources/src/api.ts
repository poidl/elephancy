import * as req from "./myrequest";

let host = '127.0.0.1';
let port = 8080;
let basepath = '/api';
// let basepath = 'http://127.0.0.1:8080/api';

export class Link {
    'rel': string;
    'href': string;
}

export class Links extends Array<Link> {
}

export class Pages extends Array<Page> {
}

interface PageRaw {
    'id': number;
    'links': Links;
    'prettyurl': string;
    'linkname': string;
    'linkweight': string;
    'metatitle': string;
    [key: string]: any;
}


function applyMixins(p1: PageRaw, p2: Page) {
    Object.getOwnPropertyNames(p1).forEach(key => {
        p2[key] = p1[key]
    })
}

export class Page {
    'id': number;
    'links': Links;
    'prettyurl': string;
    'linkname': string;
    'linkweight': string;
    'metatitle': string;
    constructor(p: PageRaw) {
        applyMixins(p, this)
    };
    // currently unused
    // getLinkByRel(this: Page, rel: string): string {
    //     for (let l of this.links) {
    //         if (l.rel === 'self') {
    //             return l.href
    //         }
    //     }
    //     // throw new TypeError("Link not found: rel: self ");
    //     return ''
    // }
    [key: string]: any;
}

export class PagesContainer {
    pages: Pages
    // Array cannot be a super class
    // https://github.com/Microsoft/TypeScript/wiki/FAQ#why-doesnt-extending-built-ins-like-error-array-and-map-work
    constructor(pages: Pages) {
        this.pages = pages
    }
    findPageByKeyValue = function (this: PagesContainer, key: string, value: any): Page {
        let p = this.pages.filter(function (p: Page) {
            return p[key] === value
        })
        return p[0]
    }
}

function toPageArray(val: PageRaw): Page {
    return new Page(val)
}

function parse_pages(obj: { code: number, body: string }): Promise<Pages> {
    return new Promise((resolve, reject) => {
        if (!(obj.code >= 200 && obj.code <= 299)) {
            // If this happens, something is wrong with the app (independent
            // of the user's request). Shutdown.
            shutdown({ code: obj.code, body: obj.body })
        } else {
            let pages = JSON.parse(obj.body)
            resolve(pages.map(toPageArray))
        }
    })
}

export function shutdown(obj: { code: number, body: string }) {
    let err = new Error('Server responded: ' + JSON.stringify(obj))
    throw err
}

export class Api {

    // public listPages(): Promise<Pages> {
    //     return new Promise<Pages>((resolve, reject) => {
    //         request(basepath + '/pages', (error, response, body) => {
    //             if (error) {
    //                 reject(error);
    //             } else {
    //                 if (response.statusCode >= 200 && response.statusCode <= 299) {
    //                     let obj = JSON.parse(body)
    //                     resolve(obj.map(toPageArray))
    //                     // resolve(body);
    //                 } else {
    //                     reject({ response: response, body: body });
    //                 }
    //             }
    //         });
    //     });

    // }
    public listPages(): Promise<Pages> {
        let options = {
            protocol: 'http:',
            hostname: host,
            port: port,
            path: basepath + '/pages',
        };


        return req.make_request(options)
            .then(req.get_code_body)
            .then(parse_pages)
        // .catch((err: Error) => {
        //     console.log(err)
        // })
    }
    public getPageContent(id: number): Promise<{ code: number, body: string }> {
        let options = {
            protocol: 'http:',
            hostname: host,
            port: port,
            path: basepath + '/content/' + id,
            headers: { myheader: 'XMLHttpRequest' }
        };


        return req.make_request(options)
            .then(req.get_code_body)
        // .catch((err: Error) => {
        //     console.log(err)
        // })
    }
}

