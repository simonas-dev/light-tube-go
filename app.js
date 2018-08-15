const app = require('express')();
const http = require('http').Server(app);
const io = require('socket.io')(http);
const execSync = require('child_process').execSync;
const exec = require('child_process').exec;
const fs = require('fs');

var goProcessCode = -1

app.get('/', function(req, res){
  res.sendFile(__dirname + '/web/index.html');
});

io.on('connection', function(socket){
  console.log('a user connected');

  socket.on("*",function(event,data) {
    console.log(event);
    console.log(data);
  });

  socket.on('command', function(msg) {
    console.log("command" + msg)
    var map = {
      "start": function() {
        goProcessCode = exec("./run")
        console.log("start " + goProcessCode.pid)
      },
      "build": function() {
        goProcessCode = exec("./build")
        console.log("build " + goProcessCode.pid)
      },
      "kill": function() {
        if (goProcessCode != -1) {
          console.log("kill " + goProcessCode.pid)
          goProcessCode = execSync("kill " + goProcessCode.pid)
        }
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