// "0" change/inform state on/off
// "1" Receive Status information

function implodeRequest(req) {
    var sreq = "";
    for (i = 0; i < req.length; i++) {
	sreq += req[i].length+":"+req[i]
    }
    // console.log(sreq)
    return sreq
}

function explodeRequest(sreq) {
    var req = [];
    for (i = 0; i < sreq.length; i++) {
	end = sreq.indexOf(":", i)
	len = parseInt(sreq.slice(i, end))
	req.push(sreq.slice(end+1, end+1+len))
	if (end > -1) {
	    i = end+len
	} else {
	    break
	}
    }
    return req
}

// SEND AND RECEIVE STATE REQUEST
// [0]: REQUEST ID
// [1]: ID
function SendChangeState(id) {
    var req = []
    req.push("0")
    req.push(id)
    sreq = implodeRequest(req)
    ws.send(sreq)
    document.getElementById("btn-"+id).disabled = true;
    setTimeout(function(){
	document.getElementById("btn-"+id).disabled = false;
    }, 2000);
}

// [0]: REQUEST ID
// [1]: ID
// [2]: State
function ChangeState(id, state) {
    if (parseInt(state) > 0) {
	document.getElementById(id).className = "Success"
	document.getElementById(id).title = "ON"
    } else {
	document.getElementById(id).className = "Danger"
	document.getElementById(id).title = "OFF"
    }
}

// SEND AND RECEIVE STATUS REQUEST
// [0]: REQUEST ID
function SendStatus(id) {
    var req = []
    req.push("1")
    req.push(id)
    sreq = implodeRequest(req)
    ws.send(sreq)
}

function Status(statInfo) {
    UpdateGraph(statInfo)
}


// WAIT FOR CONNECTION TO CONNECT
function waitForSocketConnection(socket, callback){
    setTimeout(
	function () {
	    if (socket.readyState === 1) {
		if(callback != null){
		    callback(path[2]);
		}
		return;
	    } else {
		waitForSocketConnection(socket, callback);
	    }
	}, 5); // wait 5 milisecond for the connection...
}

var ip = location.host;
var path = window.location.pathname;
var id = path.split("/");

var ws = new WebSocket("ws://"+ip+"/ws")
ws.onmessage = function (event) {
    if (event.data != "") {
	arr = explodeRequest(event.data)
	console.log(arr)
	switch (arr[0]) {
	case "0":
	    ChangeState(arr[1], arr[2])
	case "1":
	    Status(arr[1])
	}
    }
}
