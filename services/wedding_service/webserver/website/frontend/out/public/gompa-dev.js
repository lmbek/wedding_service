//import {pageChanger} from "./gladsaxefolkekor";

//import PageChanger from "../components/resources/js/pagechanger/pagechanger.js";

if (!window.socket) {
    window.socket = new WebSocket('wss://localhost:8443/websocket/hotreload');

    window.socket.onmessage = function(event) {
        if (event.data) {
            location.reload()
            //let pageChanger = new PageChanger("body a[href]")
            //pageChanger.getPage(location.href, true);
        } else {
            // TODO: implement error here, giving error page or something...
            // insert error page
            //document.body.innerHTML = event.data;
        }
    };

    window.socket.onerror = function(error) {
        console.log('WebSocket Error:', error);
    };

}