// import { Api } from "./api";
// // import { Page } from "./api";
// import { Link } from "./api";


// let api = new Api()

// // async function list() {
// //     try {
// //         let pages = await api.listPages()
// //         console.log(pages[0].getLinkByRel('self'))
// //     }
// //     catch (e) {
// //         console.log('there was error an calling listPages');
// //         console.log(e);
// //     }
// // }

// // list()

// api.test()





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




// async function test() {
//     try {
//         let resp = await api.findPageByKeyValue("prettyurl", "/")
//         let selflink = (<bla>resp.body).getLinkByRel('self')

//         // let page2 = new Page2(resp.body)
//         // let selflink = page2.getLinkByRel("self")

//         // let selflink = page.links.filter(bla("self"))[0].href

//         console.log(selflink)

//     }
//     catch (e) {
//         console.log('there was error an calling ajax()');
//         console.log(e);
//     }
//     //     // metatitle = mapmap(pagesmap,["Urlpath",e.target.href],"Metatitle")
//     //     // if (e.target.href != location.href) {
//     //     //   swapMainwindow(contenturl);
//     //     //   swapTitle(metatitle);
//     //     //   history.pushState(null, null, e.target.href);
//     //     // }
//     // }
// }

// test()


import { Api } from "./api";
// import { Page } from "./api";
import { Link } from "./api";


let api = new Api()

// async function list() {
//     try {
//         let pages = await api.listPages()
//         console.log(pages[0].getLinkByRel('self'))
//     }
//     catch (e) {
//         console.log('there was error an calling listPages');
//         console.log(e);
//     }
// }

// list()

api.test()


