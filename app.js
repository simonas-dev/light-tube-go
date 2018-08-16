const app = require('express')();
const http = require('http').Server(app);
const io = require('socket.io')(http);
const execSync = require('child_process').execSync;
const exec = require('child_process').exec;
const fs = require('fs');

var goProcessCode = null

function killProcess() {
  if (goProcessCode != null) {
    console.log("sudo pkill -TERM -P " + goProcessCode.pid)
    exec("sudo pkill -TERM -P " + goProcessCode.pid)
  }
}

app.get('/', function(req, res){
  res.sendFile(__dirname + '/web/index.html');
});

io.on('connection', function(socket){
  console.log('a user connected');

  socket.on('command', function(msg) {
    console.log("command" + msg)
    var map = {
      "start": function() {
        killProcess()
        goProcessCode = exec("sudo ./bin/app", (error, stdout, stderr) => {
          if (error) {
            console.error(`exec error: ${error}`);
            return;
          }
          console.log(`stdout: ${stdout}`);
          console.log(`stderr: ${stderr}`);
        });
        console.log("start " + goProcessCode.pid)
      },
      "build": function() {
        killProcess()
        execSync("go build -o ./bin/app cmd/app.go")
        goProcessCode = exec("sudo ./bin/app")
	console.log("build " + goProcessCode.pid)
      },
      "kill": function() {
        killProcess()
      }
    };
    map[msg]();
  });
  socket.on('config', function(msg) {
    console.log(msg)
    var json = JSON.parse(msg)
    fs.writeFile('./config.json', JSON.stringify(json, null, 2), function (err) {
        if (err) 
            return console.log(err);
    });
  });
});

http.listen(80, function(){
  console.log('listening on *:80');
});
