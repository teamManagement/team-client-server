import { Navigate, Route } from "react-router-dom";
import AppManage from "../views/AppManage";
import Home from "../views/Home";
import Manage from "../views/Manage";
import HomeType from "../views/Type";
import ManageUserMgr from "../views/ManageUserMgr";
import ClientManage from "../views/ClientMan";

/**
 * 菜单内容
 */
export const LayoutContentEles: React.ReactNode[] = [
  <Route path="/" key="other-route" element={<Navigate key="navigate-to-home" replace to={'/home/recommend'} />} />,
  <Route path="/home" key="/home" element={<Home />} />,
  <Route path='/home/type1' key='/home/type1' element={<HomeType />} />,
  <Route path='/manage/types' key='/manage/types' element={<Manage />} />,
  <Route path='/manage/apps' key='/manage/apps' element={<AppManage />} />,
  <Route path='/manage/userMgr' key='/manage/userMgr' element={<ManageUserMgr />} />,
  <Route path='/manage/clientMan' key='/manage/clientMan' element={<ClientManage />} />,

];
