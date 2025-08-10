//import {pageChanger} from "./gladsaxefolkekor";

//import PageChanger from "../components/resources/js/pagechanger/pagechanger.js";

if (!window.socket) {
    const wsProto = window.location.protocol === 'https:' ? 'wss' : 'ws';
    const wsHost = window.location.host; // includes hostname:port
    const wsPath = '/bryllup/websocket/hotreload/'; // server expects trailing slash
    window.socket = new WebSocket(`${wsProto}://${wsHost}${wsPath}`);

    window.socket.onmessage = function(event) {
        if (event.data) {
            location.reload();
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