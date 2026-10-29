// Service Worker for offline-first PWA behaviour
// This enables RakshaGrid to queue SOS requests even when the internet is down
// and flush them when connectivity returns via Background Sync API

const CACHE_NAME = "rakshagrid-v1";
const urlsToCache = [
  "/",
  "/index.html",
  "/manifest.json",
  "/src/main.jsx",
  "/src/App.jsx",
];

// Install event - cache app shell
self.addEventListener("install", (event) => {
  event.waitUntil(
    caches.open(CACHE_NAME).then((cache) => {
      return cache.addAll(urlsToCache);
    })
  );
});

// Activate event - clean up old caches
self.addEventListener("activate", (event) => {
  event.waitUntil(
    caches.keys().then((cacheNames) => {
      return Promise.all(
        cacheNames.map((cacheName) => {
          if (cacheName !== CACHE_NAME) {
            return caches.delete(cacheName);
          }
        })
      );
    })
  );
});

// Fetch event - serve from cache, fallback to network
self.addEventListener("fetch", (event) => {
  // Don't cache API requests
  if (event.request.url.includes("/v1/")) {
    event.respondWith(
      fetch(event.request).catch(() => {
        // If offline, queue the request via IndexedDB for later
        if (event.request.method === "POST") {
          queueRequest(event.request);
        }
        return new Response("Offline - request queued", {
          status: 202,
        });
      })
    );
  } else {
    // Cache-first strategy for static assets
    event.respondWith(
      caches.match(event.request).then((response) => {
        return response || fetch(event.request);
      })
    );
  }
});

// Queue requests in IndexedDB when offline
function queueRequest(request) {
  const db = indexedDB.open("RakshaGridQueue", 1);
  db.onsuccess = (event) => {
    const store = event.target.result
      .transaction("requests", "readwrite")
      .objectStore("requests");
    store.add({
      url: request.url,
      method: request.method,
      body: request.body,
      timestamp: Date.now(),
    });
  };
}

// Background Sync - flush queued requests when back online
self.addEventListener("sync", (event) => {
  if (event.tag === "sos-sync") {
    event.waitUntil(flushQueuedRequests());
  }
});

async function flushQueuedRequests() {
  const db = indexedDB.open("RakshaGridQueue", 1);
  db.onsuccess = (event) => {
    const store = event.target.result
      .transaction("requests", "readonly")
      .objectStore("requests");
    const requests = store.getAll();
    requests.onsuccess = async () => {
      for (const req of requests.result) {
        try {
          await fetch(req.url, {
            method: req.method,
            body: req.body,
          });
        } catch (err) {
          console.error("Failed to flush queued request:", err);
        }
      }
    };
  };
}
