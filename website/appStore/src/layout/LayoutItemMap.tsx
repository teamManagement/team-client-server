import { Navigate, Route } from "react-router-dom";
import AppManage from "../views/AppManage";
import Home from "../views/Home";
import Manage from "../views/Manage";
import TestSwiper from "../views/TestSwiper";
import HomeType from "../views/Type";

/**
 * 菜单内容
 */
export const LayoutContentEles: React.ReactNode[] = [
  <Route path="/" key="other-route" element={<Navigate key="navigate-to-home" replace to={'/home/recommend'} />} />,
  <Route path="/home/recommend" key="/home/recommend" element={<Home />} />,
  <Route path='/home/type1' key='/home/type1' element={<HomeType />} />,
  <Route path='/manage/types' key='/manage/types' element={<Manage />} />,
  <Route path='/manage/apps' key='/manage/apps' element={<AppManage />} />,
  <Route path='/manage/test' key='/manage/test' element={<TestSwiper />} />,
];
