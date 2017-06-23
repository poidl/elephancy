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
}


export class Page {
    'id': number;
    'links': Links;
    'prettyurl': string;
    'linkname': string;
    'linkweight': string;
    'metatitle': string;
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


function parse_pages(obj: { code: number, body: string }): Promise<Pages> {
    return new Promise((resolve, reject) => {
        if (!(obj.code >= 200 && obj.code <= 299)) {
            // If this happens, something is wrong with the app (independent
            // of the user's request). Shutdown.
            shutdown({ code: obj.code, body: obj.body })
        } else {
            let pages = <Pages>JSON.parse(obj.body)
            resolve(pages)
        }
    })
}

export function shutdown(obj: { code: number, body: string }) {
    let err = new Error('Server responded: ' + JSON.stringify(obj))
    throw err
}

export class Api {

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

