// 0 change/inform state on/off

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

function ChangeState(id, state) {
    if (parseInt(state) > 0) {
	document.getElementById(id).className = "Success"
	document.getElementById(id).title = "ON"
    } else {
	document.getElementById(id).className = "Danger"
	document.getElementById(id).title = "OFF"
    }
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
	}
    }
}
