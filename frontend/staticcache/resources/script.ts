import request = require('request');
import http = require('http');
import { DefaultApi } from "./api";
import { Page } from "./api";
import { Link } from "./api";


// declare module "./api" {
//     export interface Page {
//         getLinkByRel(rel: string): string;
//     }
// }
// Page.prototype.getLinkByRel = function(this: Page, rel: string): string {
//     for (let l of this.links) {
//         if (l.rel === 'self') {
//             return l.href
//         }
//     }
//     // throw new TypeError("Link not found: rel: self ");
//     return ''
// }

// interface linkUtils {
//     getLinkByRel(this: Page, rel: string): string;
// }

// class Page2 extends Page {
//     getLinkByRel(this: this, rel: string): string {
//         for (let l of this.links) {
//             if (l.rel === 'self') {
//                 return l.href
//             }
//         }
//         // throw new TypeError("Link not found: rel: self ");
//         return ''
//     }
// }



let api = new DefaultApi()

async function list() {
    try {
        let r = await api.listPages()
        let pages = r.body
    }
    catch (e) {
        console.log('there was error an calling listPages');
        console.log(e);
    }
}




// list()

// var container = document.getElementById('leftDrawer');
// container.addEventListener('click', navlink_clicked, false);


// async function get_content(link: string) : Promise<{ response: http.ClientResponse; body: string;  }>{

//     let requestOptions: request.Options = {
//             method: 'GET',
//             qs: '',
//             headers: [{key: 'myheader', value: 'XMLHttpRequest'}],
//             uri: link,
//             useQuerystring: false,
//             json: false,
//         };
//         return new Promise<{ response: http.ClientResponse; body: string;  }>((resolve, reject) => {
//         request(requestOptions, (error, response, body) => {
//             if (error) {
//                 reject(error);
//             } else {
//                 if (response.statusCode >= 200 && response.statusCode <= 299) {
//                     resolve({ response: response, body: body });
//                 } else {
//                     reject({ response: response, body: body });
//                 }
//             }
//         });
//     });
// }








// async function ajax(e) {
//     console.log(e.target.pathname)
//     // if ((e.target != e.currentTarget) && (e.target.className === "xhr")) {
//         try {
//             e.preventDefault();
//             e.stopPropagation();
//             let resp = await api.findPageByKeyValue("prettyurl", e.target.pathname)
//             let selflink = resp.body.getLinkByRel('self')

//             // let r = await get_content(link)
//             // let body = r.body
//             // console.log(body)

//         }
//         catch (e) {
//             console.log('there was error an calling ajax()');
//             console.log(e);
//         }
//     //     // metatitle = mapmap(pagesmap,["Urlpath",e.target.href],"Metatitle")
//     //     // if (e.target.href != location.href) {
//     //     //   swapMainwindow(contenturl);
//     //     //   swapTitle(metatitle);
//     //     //   history.pushState(null, null, e.target.href);
//     //     // }
//     // }
// }

// function navlink_clicked(e) {
//     console.log('navlink_clicked')
//     ajax(e);
//     //   na vdrawer_toggle(); // close navdrawer after click, in case it is open
// }



function getLinkByRel(page: Page, rel: string): string {
    for (let l of page.links) {
        if (l.rel === rel) {
            return l.href
        }
    }
    // throw new TypeError("Link not found: rel: self ");
    return ''
}

interface bla extends Page {
    getLinkByRel(rel: string): string;
}


// function applyMixins(p2: Page2, p: Page) {
//     Object.getOwnPropertyNames(p).forEach(key => {
//         p2[key] = p[key]
//     })
// }


// interface Page2 {
//     getLinkByRel(rel: string): string;
// }

// class Page2 {
//     constructor(page: Page) {
//         applyMixins(this, page)
//     }
// }

// Page2.prototype.getLinkByRel = function (this: Page, rel: string): string {
//     for (let l of this.links) {
//         if (l.rel === rel) {
//             return l.href
//         }
//     }
//     // throw new TypeError("Link not found: rel: self ");
//     return ''
// }


async function test() {
    try {
        let resp = await api.findPageByKeyValue("prettyurl", "/")
        let selflink = (<bla>resp.body).getLinkByRel('self')

        // let page2 = new Page2(resp.body)
        // let selflink = page2.getLinkByRel("self")

        // let selflink = page.links.filter(bla("self"))[0].href

        console.log(selflink)

    }
    catch (e) {
        console.log('there was error an calling ajax()');
        console.log(e);
    }
    //     // metatitle = mapmap(pagesmap,["Urlpath",e.target.href],"Metatitle")
    //     // if (e.target.href != location.href) {
    //     //   swapMainwindow(contenturl);
    //     //   swapTitle(metatitle);
    //     //   history.pushState(null, null, e.target.href);
    //     // }
    // }
}

test()


