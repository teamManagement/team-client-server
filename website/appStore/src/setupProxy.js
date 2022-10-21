const proxy = require("http-proxy-middleware").createProxyMiddleware;

const commonOptions = {
  target: "https://apps.byzk.cn",
  changeOrigin: true,
  secure: false,
  // agent: new ProxyAgent(proxyUri)
};

module.exports = function (app) {
  app.use(
    proxy("/appstore", {
      ...commonOptions,
      pathRewrite: {
        "^/appstore": "/appstore",
      },
    })
  );
};
