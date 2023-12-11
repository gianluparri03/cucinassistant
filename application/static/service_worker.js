const cacheName = 'CucinAssistant';
const filesToCache = ['/', '/static/*'];


self.addEventListener('install', (e) => {
    e.waitUntil((async () => {
        const cache = await caches.open(cacheName)
        await cache.addAll(filesToCache);
    })());
});

self.addEventListener('activate', (e) => {
    e.waitUntil((async () => {
        const cacheKeys = await caches.keys();
        cacheKeys.map(async (key) => {
            if (key !== cacheName) { await caches.delete(key); }
        });
    })());

    return self.clients.claim();
});


self.addEventListener('fetch', (e) => {
    e.respondWith((async () => {
        const r = await caches.match(e.request);
        if (r) { return r; }

        try {
            const response = await fetch(e.request);
            return response;
        } catch(error) {
            const cache = await caches.open(cacheName);
            const cachedResponse = await cache.match('offline.html');
            return cachedResponse;
        }
    })());
});
