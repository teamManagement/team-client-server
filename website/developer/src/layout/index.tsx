import { Layout, Menu, message, Spin, Image, Button } from "antd";
import { PlusOutlined } from "@ant-design/icons";
import { useCallback, useEffect, useRef, useState } from "react";
import { MenuItemInfo } from "./types";
import { Routes, useLocation, useNavigate } from "react-router-dom";
import { LayoutContentEles } from "./LayoutItemMap";
import "./index.less";
import AddNewApp from "../views/Home/addNewApp";

const { Sider, Content } = Layout;

function uuid() {
  return 'xxxxxxxx-xxxx-4xxx-yxxx-xxxxxxxxxxxx'.replace(/[xy]/g, function (c) {
    var r = Math.random() * 16 | 0,
      v = c == 'x' ? r : (r & 0x3 | 0x8);
    return v.toString(16);
  });
}

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
  const [contentMenuList, setContentMenuList] = useState<any[]>([])
  const [image, setImage] = useState<any>()
  const fnsRef = useRef<any>()

  const queryMenuList = useCallback(async () => {
    try {
      setLoading(true);
      const a = [
        { id: '123', name: '/app', title: '首页', icon: 'icon-yunyingpan', sort: 1 },
      ]
      setMenuItems(a)
    } catch (e) {
      message.error(e as string)
    } finally {
      setLoading(false)
    }
  }, []);

  const queryContentMenuList = useCallback(async (menu?: MenuItemInfo) => {
    setLoading(true);
    try {
      setMenuContentItems([]);
      if (!menu) {
        return;
      }
      try {
        window.teamworkSDK.store.set(menuItems[0].id + "_content_menu_list", JSON.stringify(contentMenuList))
        setMenuContentItems(contentMenuList);
        if (contentMenuList.length > 0) {
          navigate(contentMenuList[0].name);
        }
      } catch (e) {
        message.error(e as string);
      }
    } finally {
      setLoading(false)
    }
  }, [navigate, contentMenuList]);

  const addNewApp = useCallback(() => {
    const newMenuList = JSON.parse(JSON.stringify(contentMenuList))
    newMenuList.push({ id: uuid(), name: `/app/${uuid()}/新应用`, title: '新应用', icon: 'https://127.0.0.1:65528/icons/undefined.png', sort: 1 })
    setContentMenuList(newMenuList)
    window.teamworkSDK.store.set(menuItems[0].id + "_content_menu_list", JSON.stringify(newMenuList))
    setMenuContentItems(newMenuList);
    if (newMenuList.length > 0) {
      navigate(newMenuList[newMenuList.length - 1].name);
    }
  }, [menuItems, contentMenuList])

  useEffect(() => { queryMenuList(); }, [queryMenuList]);

  const getStoreList = useCallback(async () => {
    setLoading(true)
    const data: any = await window.teamworkSDK.store.get<{ [key: string]: string }>(menuItems[0]?.id + "_content_menu_list")
    console.log("store: ", data)
    if (data.length > 0) {
      setContentMenuList(data)
    }
    setLoading(false)
    const list: any = await window.teamworkSDK.store.get<{ [key: string]: string }>('image-single')
    console.log(list);
    setImage(list)
  }, [menuItems])

  useEffect(() => {
    getStoreList()
  }, [getStoreList])

  useEffect(() => {
    if (menuItems.length === 0) { return; }
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

  const headerItem = useCallback((menuItemsInfo: MenuItemInfo[]) => {
    return menuItemsInfo.map((m) => {
      console.log(m);

      return (
        <Menu.Item
          icon={<Image style={{ marginLeft: 15, marginTop: -3 }} width={20} preview={false} src={m.icon} />}
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
  }, [navigate, image]);

  return (
    <Spin spinning={loading} tip="内容正在加载...">
      {menuContentItems.length > 0 ? <Layout className="project-layout">
        <Layout>
          {menuContentItems.length > 0 && (
            <Sider width={200} className="site-layout-background-sider" >
              <Menu
                mode="inline"
                selectedKeys={topMenuCurrentSelectedKey}
                defaultSelectedKeys={contentMenuList[contentMenuList.length - 1].name}
                style={{
                  borderRight: 0,
                  paddingTop: 19,
                  marginLeft: 0,
                }}
              >
                {headerItem(menuContentItems)}
              </Menu>
              <div className="addAppType">
                <Button type='link' className="type-btn" icon={<PlusOutlined />} onClick={() => addNewApp()}>添加应用</Button>
              </div>
            </Sider>
          )}
          <Layout style={{ padding: "0 24px 24px" }}>
            <Content
              className="site-layout-background"
              style={{
                height: "100%",
                padding: 24,
                margin: 0,
                marginTop: 24,
                overflow: 'auto',
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
        : <div className="add-btn">
          <Button type='primary' className="add-btn-item" onClick={() => fnsRef.current.show()}>新增应用</Button>
        </div>
      }
      <AddNewApp fns={fnsRef} finished={(name) => {
        setContentMenuList([{ id: uuid(), name: `/app/apps/${name}`, title: `${name}`, icon: image?image:'https://127.0.0.1:65528/icons/undefined.png', sort: 1 }])
      }} />
    </Spin >
  );
}
