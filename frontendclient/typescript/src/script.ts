import { PagesContainer } from "./api";
import { PageVM } from "./pagevm";
import { PageVMElements } from "./pagevm";


import * as req from "./myrequest";

let host = '127.0.0.1';
let port = 8080;
let basepath = '/api';

let elements: PageVMElements = {
    linkcontainer: <HTMLElement>document.querySelector('.linkcontainer'),
    linklist: <HTMLElement>document.getElementById('linklist'),
    titledesktop: <HTMLElement>document.querySelector('.title-desktop'),
    titlemobile: <HTMLElement>document.querySelector('.title-mobile'),
    mainpanel: <HTMLElement>document.getElementById("mainPanel"),
    metatitle: <HTMLElement>document.getElementById("metatitle"),
    topbarmobile: <HTMLElement>document.querySelector(".top-bar-mobile"),
    menubutton: <HTMLElement>document.querySelector(".menubutton"),
}

let vm = new PageVM(elements)


window.onload = function () {
    window.addEventListener("popstate", doit, false);

    function doit() {
        let p = new PagesContainer(vm.pages)
        let page = p.findPageByKeyValue('prettyurl', '/' + location.href.split('/').pop())
        if (!page) {
            let err = new Error('Error in popstate event handler')
            throw err
        }
        vm.page = page
    }
}