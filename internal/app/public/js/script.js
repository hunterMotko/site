const obs = new PerformanceObserver(list => {
 const entries = list.getEntries();
 console.log(entries);
});

obs.observe({ entryTypes: ["largest-contentful-paint", "layout-shift", "first-input", "paint"], buffered: true });
