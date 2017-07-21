// 0 change/inform state on/off

function implodeRequest(req) {
    var sreq = "";
    for (i = 0; i < req.length; i++) {
	sreq += req[i].length+":"+req[i]
    }
    console.log(sreq)
    return sreq
}

function ChangeState(id) {
    var req = []
    req.push("0")
    req.push(id)
    sreq = implodeRequest(req)
    ws.send(sreq)
}


var ip = location.host;
var path = window.location.pathname;
var id = path.split("/");
var ws = new WebSocket("ws://"+ip+"/ws")
