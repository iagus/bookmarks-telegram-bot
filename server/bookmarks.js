const http = require('http');
const path = require('path');
const fs = require('fs');
const readline = require('readline');
const port = 3000;
const htmlPath = path.join(__dirname, 'index.html');
const dataFile = process.env.BOOKMARKS_TELEGRAM_BOT_PATH;
const cache = { mtime: null }

const server = http.createServer((req, res) => {
  const stat = fs.statSync(dataFile, { throwIfNoEntry: false });

  if (!stat) {
    res.writeHead(503, { 'Content-Type': 'text/plain' });
    return res.end('Maybe try again later');
  }

  if (stat.mtime !== cache.mtime) {
    // render data into temporary html to prevent serving partial rewrites
    // we will replace the old version with it once it's done
    const tmpHtmlPath = htmlPath + '.tmp'
    const html = fs.createWriteStream(tmpHtmlPath);
    const readStream = fs.createReadStream(dataFile, { encoding: 'utf8' });
    const rl = readline.createInterface({
      input: readStream,
      crlfDelay: Infinity
    });

    html.write('<!DOCTYPE html><body>');

    rl.on('line', (line) => {
      bm = JSON.parse(line)
      html.write(`<div><a href="${bm.Link}">${bm?.data?.title}</a></div>`);
    });

    rl.on('close', function () {
      html.write('</body></html>');
      html.end(() => fs.renameSync(tmpHtmlPath, htmlPath));
    });

    // updates cache
    cache.mtime = stat.mtime;
  }

  res.writeHead(200,{
    'Content-Type': 'text/html',
    'Content-Length': stat.size,
    'ETag': stat.mtime,
    'LastModified': stat.mtime.toUTCString()
  });

  fs.createReadStream(htmlPath).pipe(res);
});

server.listen(port, () => {
  console.log(`Server running at http://localhost:${port}/`);
});

