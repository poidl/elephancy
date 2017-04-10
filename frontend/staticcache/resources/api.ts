import request = require('request');
import http = require('http');

let basepath = 'http://127.0.0.1:8080/api';

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

    // public async findPageByKeyValue (key: string, value: string) : {

    // }
    /**
     * 
     * Returns all pages
     */
    public listPages(): Promise<Array<Page>> {
        return new Promise<Array<Page>>((resolve, reject) => {
            request(basepath + '/pages', (error, response, body) => {
                if (error) {
                    reject(error);
                } else {
                    if (response.statusCode >= 200 && response.statusCode <= 299) {
                        let obj = JSON.parse(body)
                        resolve(obj.map(toPageArray))
                        // resolve(body);
                    } else {
                        reject({ response: response, body: body });
                    }
                }
            });
        });
        // return new Promise(function (resolve, reject) {
        //     // https://developers.google.com/web/fundamentals/getting-started/primers/promises#promisifying_xmlhttprequest
            

        //     request(basepath + '/pages', function(error, response, body) {
        //         console.log(body);
        //         });
        //     let req = new XMLHttpRequest();
        //     req.open('GET', basepath + '/pages');

        //     req.onload = function () {
        //         // This is called even on 404 etc
        //         // so check the status
        //         if (req.status == 200) {
        //             // Resolve the promise with the response text
        //             let obj = JSON.parse(req.response)
        //             resolve(obj.map(toPageArray))
        //         }
        //         else {
        //             // Otherwise reject with the status text
        //             // which will hopefully be a meaningful error
        //             reject(Error(req.statusText));
        //         }
        //     };

        //     // Handle network errors
        //     req.onerror = function () {
        //         reject(Error("Network Error"));
        //     };

        //     // Make the request
        //     req.send();
        // });

    }
}
