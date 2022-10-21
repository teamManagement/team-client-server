import ReactDOM from "react-dom/client";
import { HashRouter } from "react-router-dom";
import Router from "./router";
import { ConfigProvider } from "antd";
import zhCN from "antd/lib/locale/zh_CN";
import "./index.less";

const root = ReactDOM.createRoot(
  document.getElementById("root") as HTMLElement
);
root.render(
  <HashRouter>
    <ConfigProvider locale={zhCN}>
      <Router />
    </ConfigProvider>
  </HashRouter>
);
