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
    public listPages(): Promise<Array<Page>> {
        let options = {
            protocol: 'http:',
            hostname: host,
            port: port,
            path: basepath + '/pag',
        };


        return make_request(options)
            .then(get_code_body)
            .then(parse)
        // .catch((err: Error) => {
        //     console.log(err)
        // })
    }
}

function make_request(options: any): Promise<http.IncomingMessage> {
    let req = http.request(options)
    req.end()
    return new Promise<http.IncomingMessage>((resolve, reject) => {
        req.on('response', function (message: http.IncomingMessage) {
            resolve(message)
        })
        // an error is thrown here if e.g. connection is refused (e.g. wrong port number)
        req.on('error', function (err: Error) {
            console.log(err)
            throw err
        })
    })
}

function get_code_body(message: http.IncomingMessage): Promise<{ code: number, body: string }> {
    let body: string = ''
    let err: Error
    message.on('error', function (err: Error) {
        err = err
    });
    message.on('data', function (chunk) {
        body += chunk;
    });
    let code = message.statusCode
    return new Promise<{ code: number, body: string }>((resolve, reject) => {
        if (err) {
            console.log(err)
            throw (err)
        }
        message.on('end', function () {
            if (code >= 200 && code <= 299) {
                resolve({ code: code, body: body })
            } else {
                // shutdown
                handle_nosuccess({ code: code, body: body })
            }
        });
    })
}

function parse(obj: { code: number, body: string }): Promise<Array<Page>> {
    return new Promise((resolve, reject) => {
        let pages = JSON.parse(obj.body)
        resolve(pages.map(toPageArray))
    })
}

function handle_nosuccess(obj: { code: number, body: string }) {
    let err = new Error('Server responded: ' + JSON.stringify(obj))
    throw err
}
