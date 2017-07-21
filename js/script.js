// 0 change/inform state on/off

function implodeRequest(req) {
    var sreq = "";
    for (i = 0; i < req.length; i++) {
	sreq += req[i].length+":"+req[i]
    }
    console.log(sreq)
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

function ChangeState(id) {
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


var ip = location.host;
var path = window.location.pathname;
var id = path.split("/");

var ws = new WebSocket("ws://"+ip+"/ws")
ws.onmessage = function (event) {
    if (event.data != "") {
	explodeRequest(event.data)
    }
}
