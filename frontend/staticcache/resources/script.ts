
import { DefaultApi } from "./api";
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

var container = document.getElementById('leftDrawer');
container.addEventListener('click', navlink_clicked, false);


function ajax(e) {
    console.log(e.target)
    console.log(e.currentTarget)
    console.log(e.target.className)
    // if ((e.target != e.currentTarget) && (e.target.className === "xhr")) {
        let linksSelf = api.findPageByKeyValue("prettyurl")
        console.log(linksSelf)
    //     // metatitle = mapmap(pagesmap,["Urlpath",e.target.href],"Metatitle")
        e.preventDefault();
    //     // if (e.target.href != location.href) {
    //     //   swapMainwindow(contenturl);
    //     //   swapTitle(metatitle);
    //     //   history.pushState(null, null, e.target.href);
    //     // }
        e.stopPropagation();
    // }
}

function navlink_clicked (e) {
    console.log('navlink_clicked')
  ajax(e);
//   na vdrawer_toggle(); // close navdrawer after click, in case it is open
}


