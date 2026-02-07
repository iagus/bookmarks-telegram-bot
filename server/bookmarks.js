import { createServer } from 'http';
import { join } from 'path';
import { createInterface } from 'readline';
import { once } from 'events';
import {
  statSync,
  createWriteStream,
  createReadStream,
  renameSync
} from 'fs';

const port = 3001;
const htmlPath = 'index.html';
const dataFile = process.env.BOOKMARKS_PATH;
const cache = { mtimeMs: null }

async function writeLine(writer, line) {
  if (!writer.write(line)) {
    const stamp = `[bookmarks node][${new Date(Date.now()).toUTCString()}]`;
    console.log(`${stamp} Buffer full. Draining`);
    console.log('buffer:', writer.writableLength, '/', writer.writableHighWaterMark);
    await once(writer, 'drain');
    console.log(`${stamp} Buffer drained`);
  }
};

function writeHeaders(res, data) {
  res.writeHead(200, {
    'Content-Type': 'text/html',
    'Content-Length': data.size,
    'ETag': data.mtimeMs,
    'Last-Modified': data.mtime.toUTCString(),
    'Cache-Control': 'public, max-age=3600'
  });
}

function updateCache(stat) {
  cache.mtimeMs = stat.mtimeMs;
}

const server = createServer(async (req, res) => {
  if (req.headers['if-none-match'] === cache.mtimeMs) {
    res.writeHead(304, {
      'ETag': cache.mtimeMs,
      'Cache-Control': 'public, max-age=3600'
    });
    res.end();
    return;
  }

  const stat = statSync(dataFile, { throwIfNoEntry: false });
  if (!stat) {
    res.writeHead(503, { 'Content-Type': 'text/plain' });
    res.end('Maybe try again later');
    return;
  }

  if (stat.mtimeMs !== cache.mtimeMs) {
    // render data into temporary html to prevent serving partial rewrites
    // we will replace the old version with it once it's done
    const tmpHtmlPath = htmlPath + '.tmp'
    const html = createWriteStream(tmpHtmlPath);
    const readStream = createReadStream(dataFile, { encoding: 'utf8' });
    const rl = createInterface({
      input: readStream,
      crlfDelay: Infinity
    });

    try {
      await writeLine(html, '<!DOCTYPE html><body>');

      for await (const line of rl) {
        const bm = JSON.parse(line)
        const rendered = `<div><a href="${bm.Link}">${bm?.data?.title}</a></div>`;

        await writeLine(html, rendered);
      }

      await writeLine(html, '</body></html>');

      html.end();
      await once(html, 'finish');

      renameSync(tmpHtmlPath, htmlPath);

      updateCache(stat);

      const contentData = statSync(htmlPath);
      writeHeaders(res, contentData);

      // pipe to response
      createReadStream(htmlPath).pipe(res);

    } catch (err) {
      html.destroy();
      throw err;
    }
  } else {

    writeHeaders(res, stat);

    // pipe to response
    createReadStream(htmlPath).pipe(res);
  }
});

server.listen(port, () => {
  console.log(`Server running at http://localhost:${port}/`);
});

