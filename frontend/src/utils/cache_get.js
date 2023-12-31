async function CachedFetch(url, options) {
    let expiry = 60 * 30;
    // let expiry = 1;
    if (typeof options === "number") {
        expiry = options;
        options = undefined;
    } else if (typeof options === "object") {
        expiry = options.seconds || expiry;
    }
    if (url.includes("get_all")) {
        expiry = 60 * 60 * 12;
        // expiry = 1;
    }
    let cacheKey = url;
    let cached = localStorage.getItem(cacheKey);
    let whenCached = localStorage.getItem(cacheKey + ":ts");
    if (cached !== null && whenCached !== null) {
        let age = (Date.now() - whenCached) / 1000;
        if (age < expiry) {
            let response = new Response(new Blob([cached]));
            return Promise.resolve(response);
        } else {
            localStorage.removeItem(cacheKey);
            localStorage.removeItem(cacheKey + ":ts");
        }
    }

    const response_1 = await fetch(url, options);
    if (response_1.status === 200) {
        let ct = response_1.headers.get("Content-Type");
        if (ct && (ct.match(/application\/json/i) || ct.match(/text\//i))) {
            response_1
                .clone()
                .text()
                .then((content) => {
                    localStorage.setItem(cacheKey, content);
                    localStorage.setItem(cacheKey + ":ts", Date.now());
                });
        }
    }

    return response_1;
}

export default CachedFetch;