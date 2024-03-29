const CACHE_NAME = 'cucinassitant-2.0';
const CACHE_URLS = [
    '/static/css/base.css',
    '/static/css/landscape.css',
    '/static/css/portrait.css',
    '/static/css/sakura.css',

    '/static/img/logo.png',

    '/static/js/base.js',
    '/static/js/account.js',
    '/static/js/lists.js',
    '/static/js/menu.js',
    '/static/js/storage.js',

    '/favicon.ico'
];


self.addEventListener('install', event => {
    event.waitUntil(
        caches.open(CACHE_NAME)
            .then(cache => cache.addAll(CACHE_URLS))
            .then(self.skipWaiting())
    );
});

self.addEventListener('activate', event => {
    const currentCaches = [CACHE_NAME];

    event.waitUntil(
        caches.keys().then(cacheNames => {
            return cacheNames.filter(cacheName => cacheName != CACHE_NAME);
        }).then(cachesToDelete => {
            return Promise.all(cachesToDelete.map(cacheToDelete => {
                return caches.delete(cacheToDelete);
            }));
        }).then(() => self.clients.claim())
    );
});

self.addEventListener('fetch', event => {
    event.respondWith(
        caches.match(event.request).then(cachedResponse => {
            if (cachedResponse) {
                return cachedResponse;
            } else {
                return fetch(event.request);
            }
        })
    );
});
