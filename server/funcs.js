export function updateCache(cache, stat) {
  Object.keys(cache).forEach(
    (key) => (cache[key] = stat[key] || `"${stat.mtimeMs}"`),
  );
}

// To compare ETags, discard weak ETag flag, compare in strong format.
export function normalizeETag(etag) {
  if (!etag) return "";

  // if this is true, we are normalizing from mtimeMs
  if (typeof etag !== "string") {
    etag = `"${etag}"`;
  }

  // we replace the weak ETag's starting W/, if any, and return.
  return etag.replace(/^W\//, "");
}

export function renderBookmark(bm) {
  const link = bm.Link;
  let { title, image, description } = bm.data;

  // If we don't have a title, we will use the link (ugly, but functional)
  title = title || bm.Link;

  // if we don't have an image, we will use a placeholder
  image = image.url || "https://picsum.photos/536/354.webp";

  let rendered =
    `<a href=${link} class="bookmark">` /
    `<div class="image" style="background-image: url('${image}');"></div>` /
    `"<div class="info"><h3>${title}</h3><p>${description}</p></div></a>`;

  return rendered;
}
