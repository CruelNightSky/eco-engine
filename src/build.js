const { exec } = require('child_process');
const fs = require('fs');

// Build go binary to .so file
exec('go build -buildmode=c-shared -o libecoengine.so main.shared.go', (err, stdout, stderr) => {
  if (err) {
    console.error(err);
    return;
  }
  console.log(stdout);
  console.error(stderr);
});

