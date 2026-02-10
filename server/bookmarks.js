import { createServer } from "http";
import { createInterface } from "readline";
import { once } from "events";
import {
  statSync,
  createWriteStream,
  createReadStream,
  renameSync,
  unlink,
} from "fs";
import { updateCache, normalizeETag, renderBookmark } from "./funcs.js";

const port = 3001;
const htmlPath = "index.html";
const dataFile = process.env.BOOKMARKS_PATH;
const cache = { mtime: null, mtimeMs: null, size: 0, etag: "" };

async function writeLine(writer, line) {
  if (!writer.write(line)) {
    const stamp = `[bookmarks node][${new Date(Date.now()).toUTCString()}]`;
    console.log(`${stamp} Buffer full. Draining`);
    console.log(
      "buffer:",
      writer.writableLength,
      "/",
      writer.writableHighWaterMark,
    );
    await once(writer, "drain");
    console.log(`${stamp} Buffer drained`);
  }
}

const server = createServer(async (req, res) => {
  const stat = statSync(dataFile, { throwIfNoEntry: false });
  if (!stat) {
    res.writeHead(503, { "Content-Type": "text/plain" });
    res.end("Maybe try again later");
    return;
  }

  // We are using the data file's mtimeMs as the ETag header value.
  // To return 304, we need to normalize the format of the request
  // header value and the data source file's stat mtimeMs value.
  const normalized_header = normalizeETag(req.headers["if-none-match"]);
  const normalized_mtimeMs = normalizeETag(stat.mtimeMs);
  if (normalized_header === normalized_mtimeMs) {
    res.writeHead(304, {
      ETag: cache.etag,
      "Cache-Control": "public, max-age=3600",
    });
    res.end();
    return;
  }

  // We are using the data file's mtimeMs as identifier.
  // So it will also function as an ETag header. Here we are comparing the
  // current data source file state's mtimeMs with the cached mtimeMs from
  // the previous html build's data source file.
  // If those values are the same, we already have a fresh rendered html that
  // we proceed to send to the client.
  // Otherwise, the html is stale and must be re-rendered, then sent to client.
  if (stat.mtimeMs !== cache.mtimeMs) {
    // render data into temporary html to prevent serving partial rewrites
    // we will replace the old version with it once it's done
    const tmpHtmlPath = htmlPath + ".tmp";
    const html = createWriteStream(tmpHtmlPath);
    const readStream = createReadStream(dataFile, { encoding: "utf8" });
    const rl = createInterface({
      input: readStream,
      crlfDelay: Infinity,
    });

    try {
      await writeLine(html, "<!DOCTYPE html><body>");

      for await (const line of rl) {
        const bm = JSON.parse(line);
        const rendered = renderBookmark(bm);

        await writeLine(html, rendered);
      }

      await writeLine(html, "</body></html>");

      html.end();
      await once(html, "finish");

      renameSync(tmpHtmlPath, htmlPath);

      const htmlStat = statSync(htmlPath);
      res.writeHead(200, {
        "Content-Type": "text/html; charset=utf-8",
        "Content-Length": htmlStat.size,
        ETag: `"${stat.mtimeMs}"`, // keeping reference to data file
        "Last-Modified": htmlStat.mtime.toUTCString(),
        "Cache-Control": "public, max-age=3600",
      });

      // Save current build information for later checks.
      updateCache(cache, {
        mtime: stat.mtime,
        mtimeMs: stat.mtimeMs,
        size: htmlStat.size,
        etag: `"${stat.mtimeMs}"`,
      });

      createReadStream(htmlPath).pipe(res);
    } catch (err) {
      rl.close();
      readStream.destroy();
      html.destroy();
      unlink(tmpHtmlPath, () => {});

      throw err;
    }
  } else {
    res.writeHead(200, {
      "Content-Type": "text/html; charset=utf-8",
      "Content-Length": cache.size,
      ETag: cache.etag,
      "Last-Modified": cache.mtime.toUTCString(),
      "Cache-Control": "public, max-age=3600",
    });

    createReadStream(htmlPath).pipe(res);
  }
});

server.listen(port, () => {
  console.log(`Server running at http://localhost:${port}/`);
});
