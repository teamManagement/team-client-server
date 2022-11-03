import { Navigate, Route } from "react-router-dom";
import Home from "../views/Home";
/**
 * 菜单内容
 */
export const LayoutContentEles: React.ReactNode[] = [
  <Route path="/" key="other-route" element={<Navigate key="navigate-to-home" replace to={'/app/:appId'} />} />,
  <Route path="/app/:appId/:appName" key="/app/:appId/:appName" element={<Home />} />,
];
