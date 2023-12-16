const cacheName = 'CucinAssistant';
var filesToCache = [
    '/static/logo.png',
    '/static/sakura.css',
    '/static/style.css',
    'https://fonts.googleapis.com/css2?family=Inclusive+Sans&family=Satisfy&display=swap',
    'https://cdn.jsdelivr.net/npm/normalize.css@8.0.1/normalize.css',
    'https://code.jquery.com/jquery-3.7.1.slim.min.js'
];

filesToCache.map(function(u) {
    return new Request(u, {mode: 'no-cors'});
});


self.addEventListener('install', (e) => {
    e.waitUntil((async () => {
        const cache = await caches.open(cacheName);
        await cache.addAll(filesToCache);
    })());
});

self.addEventListener('activate', event => {
    event.waitUntil(self.clients.claim());
});

self.addEventListener('fetch', event => {
    event.respondWith(async () => {
        const cache = await caches.open(CACHE_NAME);
        const cachedResponse = await cache.match(event.request);

        if (cachedResponse !== undefined) {
            return cachedResponse;
        } else {
            return fetch(event.request)
        }
    });
});