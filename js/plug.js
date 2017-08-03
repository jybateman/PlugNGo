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

function ChangeInput() {
    document.getElementById('NameText').style.display = "none"
    document.getElementById('NameInput').removeAttribute("style")
}

function StatusRange(id) {
    var req = []
    req.push("1")
    req.push(id)
    req.push(document.getElementById('StartDate').value)
    req.push(document.getElementById('EndDate').value)
    var sreq = implodeRequest(req)
    ws.send(sreq)
}

$('.datepicker').datepicker()
