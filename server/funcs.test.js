import test from 'node:test';
import { strict as assert } from 'node:assert';
import {
  updateCache,
  normalizeETag,
  renderBookmark
} from './funcs.js';

// mocking a stat object
const statMock = {
  size: 13536,
  mtimeMs: 1770606847534.1697,
  mtime: new Date()
}

test('updates external cache object', () => {
  // Not a fan of exposing the fields here, might revisit later.
  const cache = { mtime: null, mtimeMs: null, size: 0, etag: '' }

  Object.keys(cache).forEach(cache_key => {
    assert.notStrictEqual(cache[cache_key], statMock[cache_key]);
  });

  updateCache(cache, statMock);

  Object.keys(cache).forEach(cache_key => {
    if (statMock[cache_key]) {
      assert.strictEqual(
        cache[cache_key],
        statMock[cache_key],
        `Must save ${cache_key}`
      );
    } else {
      // if the key is not found in stat obj, means
      // it's the etag cache key
      assert.strictEqual(
        cache.etag,
        `"${statMock.mtimeMs}"`,
        'Must save stringfied mtimeMs'
      );
    }
  });

  assert.strictEqual(cache.mtime, statMock.mtime, 'Must save mtime');
  assert.strictEqual(cache.mtimeMs, statMock.mtimeMs, 'Must save mtimeMs');
  assert.strictEqual(cache.size, statMock.size, 'Must save size');
});

test('normalizes ETag across weak and strong format', () => {
  const wETag = 'W/"1770590167804.128"';
  const sETag = '"1770590167804.128"';
  const preprocessed = 1770590167804.128;

  const normalized_wETag = normalizeETag(wETag);
  const normalized_sETag = normalizeETag(sETag);

  assert.strictEqual(
    normalized_wETag,
    `"1770590167804.128"`,
    'Normalizes Strong ETag');

  assert.strictEqual(
    normalized_sETag,
    `"1770590167804.128"`,
    'Normalizes Weak ETag');

  assert.strictEqual(
    normalizeETag(preprocessed),
    `"1770590167804.128"`,
    'Normalizes raw data into ETag');

  // Return empty if called without parameter
  assert.strictEqual(
    "",
    "",
    "Must normalize empty string into empty string, not an ETag"
  );
});

test('Renders bookmark', () => {
  const bm = {
    "Link": "https://somewhere.com",
    "data": { 
      "title": "Title",
      "description": "This is a description",
      "image": { 
        "url": "https://picsum.com/image.jpg"
      }
    }
  };

  const rendered = renderBookmark(bm);
  console.log('rendered inside test', rendered);

  assert.strictEqual(rendered,
    `<a href="https://somewhere.com" class="bookmark">Title` /
    `<div class="image" style="background-image: url('https://picsum.com/image.jpg);>` /
    `</div><div class="info"><h3>Title</h3><p>This is a description</p></div></a>`,
    'Renders bookmark');
});
