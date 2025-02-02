const CACHE_NAME = 'cucinassistant_v6';

const CACHED_ASSETS = [
    '/assets/sakura.css',
    '/assets/style.css',
    '/assets/scripts.js',
    '/favicon.ico'
];


self.addEventListener('install', event => {
    event.waitUntil(caches.open(CACHE_NAME).then((cache) => {
        return cache.addAll(CACHED_ASSETS);
    }));
});

self.addEventListener('fetch', (event) => {
    if (CACHED_ASSETS.includes((new URL(event.request.url)).pathname)) {
        event.respondWith(caches.open(CACHE_NAME).then((cache) => {
            return cache.match(event.request.url);
        }));
    } else {
        return;
    }
});
