// import request = require('request');
import http = require('http')
// import url = require('url')

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
    getLinkByRel(this: Page, rel: string): string {
        for (let l of this.links) {
            if (l.rel === 'self') {
                return l.href
            }
        }
        // throw new TypeError("Link not found: rel: self ");
        return ''
    }
    [key: string]: any;
}

function toPageArray(val: PageRaw): Page {
    return new Page(val)
}



export class Api {

    // public listPages(): Promise<Array<Page>> {
    //     return new Promise<Array<Page>>((resolve, reject) => {
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
    public test(): void {
        let myoptions = {
            protocol: 'http:',
            hostname: host,
            port: port,
            path: basepath + '/pages',
        };
        function cb(message: http.IncomingMessage) {
            let body: string = ''
            message.on('data', function (chunk) {
                body += chunk;
            });
            message.on('end', function () {
                let obj = JSON.parse(body)
                let pa = obj.map(toPageArray)
                console.log(pa)
            });
        }

        http.request(myoptions, cb).end()
    }
}
