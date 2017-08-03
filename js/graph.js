function CreateGraph(id) {
     graph = new Dygraph(
	document.getElementById("graph-"+id),
	"Date, Power, Voltage\n"
    );
}

function UpdateGraph(statInfo) {
    graph.updateOptions(
	{ 'file':
	  statInfo
	}
    );
}

var path = window.location.pathname.split("/")
if (path.length >= 3) {
    var graph
    // console.log(path[2])
    CreateGraph(path[2])
    waitForSocketConnection(ws, SendStatus)
}
