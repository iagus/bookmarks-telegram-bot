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

  if (stat.mtime == cache.mtime && stat.size == cache.size) {
    // will serve html from file system
  } else {
    // render data into temporary html to prevent serving half writes
    // we will switch it with the old version of the html once it's done
    const tmpHtmlPath = htmlPath + '.tmp'
    const html = fs.createWriteStream(tmpHtmlPath);
    let output = '<!DOCTYPE html><body>';
    let bookmarks = '';
    const readStream = fs.createReadStream(dataFile, { encoding: 'utf8' });
    const rl = readline.createInterface({
      input: readStream,
      crlfDelay: Infinity
    });

    rl.on('line', (line) => {
      bookmark = JSON.parse(line)
      output += '<div>\
        <a href=\"' + bookmark.Link + '\">' +
        bookmark?.data?.title + '</a></div>';
    });

    rl.on('close', function () {
      output += '</body></html>';
      html.write(output);
      html.end(() => fs.renameSync(tmpHtmlPath, htmlPath));
    });
  }

  // updates cache
  cache.mtime = stat.mtime;

  res.writeHead(200,{
    'Content-Type': 'text/html',
    'Content-Length': stat.size,
    'ETag': stat.mtime,
    'LastModified': stat.mtime.toUTCString()
  });
  res.createReadStream(htmlPath).pipe(res);
});

server.listen(port, () => {
  console.log(`Server running at http://localhost:${port}/`);
});

