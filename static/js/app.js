window.onload = function () {

    var msg = document.getElementById("eventmsg");
    var log = document.getElementById("eventlog");

    function appendLog(item) {
        var doScroll = log.scrollTop > log.scrollHeight - log.clientHeight - 1;
        log.appendChild(item);
        if (doScroll) {
            log.scrollTop = log.scrollHeight - log.clientHeight;
        }
    }

    function setMsg(msf) {
        msg.innerHTML = msf;
    }

    if (log) {
        if (window["WebSocket"]) {

            console.log("Protocol: " + location.protocol);
            var wsURL = "ws://" + document.location.host + "/ws"

            // TODO: websocketUpgrade not set in Istio so errs on upgrade if WSS
            // if (location.protocol == 'https:') {
            //     wsURL = "wss://" + document.location.host + "/ws"
            // }
            console.log("WS URL: " + wsURL);
            var conn = new WebSocket(wsURL);

            conn.onopen = function () {
                setMsg("<b>Opening Connection</b>");
            };

            /*
                {
                    "specversion" : "0.2",
                    "type" : "com.github.pull.create",
                    "source" : "https://github.com/cloudevents/spec/pull/123",
                    "id" : "A234-1234-1234",
                    "time" : "2018-04-05T17:31:00Z",
                    "comexampleextension1" : "value",
                    "comexampleextension2" : {
                        "othervalue": 5
                    },
                    "contenttype" : "text/plain",
                    "data" : "message"
                }
            */

            conn.onmessage = function (evt) {
                var eventObj = JSON.parse(evt.data);
                console.log(eventObj);
                var item = document.createElement("div");
                item.className = "item";
                item.innerHTML = JSON.stringify(eventObj, null, 4);
                appendLog(item);
            };

            conn.onclose = function (evt) {
                setMsg("<b>Connection closed</b>");
            };


        } else {
            setMsg("<b class=error>Your browser does not support WebSockets</b>");
        }
    }
};