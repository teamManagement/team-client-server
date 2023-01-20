import { Layout, Menu, message, Spin } from "antd";
import { DownOutlined } from "@ant-design/icons";
import { useCallback, useEffect, useMemo, useState } from "react";
import { IconPro } from "../components/Icons";
import { MenuItemInfo, UserInfo } from "./types";
import { Routes, useLocation, useNavigate } from "react-router-dom";
import { LayoutContentEles } from "./LayoutItemMap";
// import { MessagePlugin } from 'tdesign-react'
import { api } from '@teamworktoolbox/inside-sdk'
import "./index.less";

const { Header, Sider, Content } = Layout;
const { SubMenu } = Menu;


/**
 * layout布局
 * @param props 菜单信息，默认选中的key
 * @returns layout布局
 */
export default function LayoutView(): React.ReactNode {
  const navigate = useNavigate();
  const location = useLocation();
  const [topMenuCurrentSelectedKey, setTopMenuCurrentSelectedKey] = useState<string[]>([]);
  const [menuItems, setMenuItems] = useState<MenuItemInfo[]>([]);
  const [menuContentItems, setMenuContentItems] = useState<MenuItemInfo[]>([]);
  const [loading, setLoading] = useState<boolean>(false);

  const [, setUserInfo] = useState<UserInfo | undefined>(undefined)




  const queryMenuList = useCallback(async () => {
    try {
      setLoading(true);
      const a = [
        { id: '123', name: '/home', title: '首页', icon: 'icon-yunyingpan', sort: 1 },
        { id: '456', name: '/manage', title: '管理', icon: 'icon-yingyong', sort: 1 },
      ]
      setMenuItems(a)
    } catch (e) {
      message.error(e as string)
    } finally {
      setLoading(false)
    }
  }, []);

  const queryContentMenuList = useCallback(
    async (menu?: MenuItemInfo) => {
      setLoading(true);
      try {
        setMenuContentItems([]);
        if (!menu) {
          return;
        }

        try {
          const list = sessionStorage.getItem(menu.id + "_content_menu_list");
          if (list) {
            const menuList: MenuItemInfo[] = JSON.parse(list);
            setMenuContentItems(menuList);
            if (menuList.length > 0) {
              navigate(menuList[0].name);
            }
            return;
          }
          let contentMenuList: any[] = []
          if (menu.id === '123') {
            contentMenuList = [
              // { id: '123456', name: '/home/recommend', title: '推荐', icon: '', sort: 1 },
              // { id: '12378934', name: '/home/type1', title: '类型1', icon: '', sort: 1 },
              // { id: '12378956', name: '/home/type2', title: '类型2', icon: '', sort: 1 },
              // { id: '12378987', name: '/home/type3', title: '类型3', icon: '', sort: 1 },
            ]
          }
          if (menu.id === '456') {
            contentMenuList = [
              { id: '456789', name: '/manage/types', title: '类别管理', icon: 'icon-yingyong', sort: 1 },
              { id: '456123', name: '/manage/apps', title: '应用管理', icon: 'icon-yingyong', sort: 1 },
              { id: '456321', name: '/manage/userMgr', title: '管理员列表', icon: 'icon-yingyong', sort: 1 },
            ]
          }
          sessionStorage.setItem(
            menu.id + "_content_menu_list",
            JSON.stringify(contentMenuList)
          );
          setMenuContentItems(contentMenuList);
          if (contentMenuList.length > 0) {
            navigate(contentMenuList[0].name);
          }
        } catch (e) {
          message.error(e as string);
        }
      } finally {
        setLoading(false);
      }
    },
    [navigate]
  );

  useEffect(() => {
    queryMenuList();
  }, [queryMenuList]);

  useEffect(() => {
    if (menuItems.length === 0) {
      return;
    }
    for (let m of menuItems) {
      if (m.name === location.pathname) {
        queryContentMenuList(m);
        break;
      }
    }
    const pathSplit = location.pathname.split("/");
    const result: string[] = [];
    for (let str of pathSplit) {
      if (str === "") {
        result.push("/");
        continue;
      }

      let prevStr = result[result.length - 1];
      if (!prevStr.endsWith("/")) {
        prevStr += "/";
      }
      result.push(prevStr + str);
    }

    setTopMenuCurrentSelectedKey(result);
  }, [location.pathname, menuItems, queryContentMenuList]);

  useEffect(() => {
    if (menuItems.length === 0) {
      sessionStorage.clear();
      return;
    }

    const cachePath = sessionStorage.getItem("c_p");
    if (cachePath) {
      if (location.pathname.startsWith(cachePath)) {
        return;
      }
      navigate(cachePath);
      return;
    }

    const targetPath = menuItems[0].name;
    sessionStorage.setItem("c_p", targetPath);
    navigate(targetPath);
  }, [menuItems, location.pathname, navigate]);

  const getUserInfo = useCallback(async () => {
    try {
      const users = await api.proxyHttpLocalServer<UserInfo>('/user/now')
      if (users.id !== '0') {
        setMenuItems([])
      }
      setUserInfo(users)
    } catch (e: any) {
      message.error('获取用户信息失败: ' + ((e as any).message || e))
    }
  }, [])

  useEffect(() => {
    getUserInfo()
  }, [getUserInfo])

  const headerItem = useCallback((menuItemsInfo: MenuItemInfo[]) => {
    return menuItemsInfo.map((m) => {
      if (m.children && m.children.length > 0) {
        return (
          <SubMenu
            key={m.name}
            icon={<IconPro type={m.icon} />}
            title={
              <>
                {m.title}&nbsp;&nbsp;
                <DownOutlined />
              </>
            }
          >
            {headerItem(m.children)}
          </SubMenu>
        );
      }
      return (
        <Menu.Item
          icon={<IconPro type={m.icon} />}
          onClick={async () => {
            sessionStorage.setItem("c_p", m.name);
            navigate(m.name);
          }}
          key={m.name}
        >
          {m.title}
        </Menu.Item>
      );
    });
  },
    [navigate]
  );

  const topMenuItemEles = useMemo(() => {
    return headerItem(menuItems);
  }, [headerItem, menuItems]);

  return (
    <Spin spinning={loading} tip="内容正在加载...">
      <Layout className="project-layout">
        {menuItems.length > 0 && <Header className="header" style={{ background: "#1345aa" }}>
          <div className="logo" />
          <Menu
            theme="light"
            mode="horizontal"
            selectedKeys={topMenuCurrentSelectedKey}
          >
            {topMenuItemEles}
          </Menu>
        </Header>}
        <Layout style={{
          margin: 20
        }} >
          {menuContentItems.length > 0 && (
            <Sider width={200} className="site-layout-background-sider" >
              <Menu
                mode="inline"
                selectedKeys={topMenuCurrentSelectedKey}
                style={{
                  borderRight: 0,
                  paddingTop: 19,
                  marginLeft: 0,
                }}
              >
                {headerItem(menuContentItems)}
              </Menu>
            </Sider>
          )}
          <Layout>
            <Content
              className="site-layout-background"
              style={{
                height: "100%",
                padding: menuContentItems.length > 0 ? '24px 18px' : 0,
                margin: 0,
                marginLeft: menuContentItems.length > 0 ? 20 : 0,
                overflow: menuContentItems.length > 0 ? 'auto' : 'hidden',
                background: "#fff",
              }}
            >
              <Routes>
                {LayoutContentEles}
              </Routes>
            </Content>
          </Layout>
        </Layout>
      </Layout>
    </Spin >
  );
}
