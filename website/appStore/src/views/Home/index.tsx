import { Button, Image, Spin } from 'antd';
import { useCallback, useEffect, useRef, useState } from 'react'
import { PlusOutlined } from "@ant-design/icons";
import './index.less'
import AddNewApp from './AddNewApp';
import AppDetail from './AppDetail';
import OtherApp from './OtherApp';
import { IconPro } from '../../components/Icons';
import { getTypeList } from '../../serve';

function uuid() {
  return 'xxxxxxxx-xxxx-4xxx-yxxx-xxxxxxxxxxxx'.replace(/[xy]/g, function (c) {
    var r = Math.random() * 16 | 0,
      v = c == 'x' ? r : (r & 0x3 | 0x8);
    return v.toString(16);
  });
}

const initImage = 'https://127.0.0.1:65528/icons/undefined.png'

const Home: React.FC = () => {
  const fnsRef = useRef<any>()
  const [menuList, setMenuList] = useState<any[]>([{}])
  const [selectedId, setSelectId] = useState<any>()
  const [firstId, setFirstId] = useState<string>()
  const [loading, setLoading] = useState<boolean>(false)

  // const addnewAppFns = useCallback(async (newName: any) => {
  //   let newMenuList: any[] = [...menuList]
  //   const addInfo = { id: uuid(), title: `${newName}`, icon: '' }
  //   newMenuList.push(addInfo)
  //   await window.teamworkSDK.store.set("_content_menu_list", newMenuList)
  //   await window.teamworkSDK.store.set(addInfo.id, addInfo)
  //   const data: any = await window.teamworkSDK.store.get("_content_menu_list")
  //   console.log("store: ", data)
  //   setSelectId(addInfo.id)
  //   setMenuList(newMenuList);
  // }, [menuList])

  const setFirstList = useCallback(async () => {
    const addInfo = { id: uuid(), title: '首页', icon: '' }
    await window.teamworkSDK.store.set("_content_menu_list", [addInfo])
    await window.teamworkSDK.store.set(addInfo.id, addInfo)
  }, [])

  useEffect(() => {
    setFirstList()
  }, [setFirstList])

  const getFileList = useCallback(async () => {
    const data: any = await window.teamworkSDK.store.get("_content_menu_list")
    console.log("store: ", data)
    setSelectId(data[0]?.id)
    setFirstId(data[0]?.id)
    setMenuList(data)
    window.teamworkSDK.store.set('appId', data[0]?.id)

  }, [])

  // const flushFileList = useCallback(async () => {
  //   const data: any = await window.teamworkSDK.store.get("_content_menu_list")
  //   setMenuList(data)
  // }, [])

  useEffect(() => {
    getFileList()
  }, [getFileList])

  const getList = useCallback(async () => {
    setLoading(true)
    const list = await getTypeList({})
    console.log(list);
    setMenuList(list)
    setLoading(false)
  }, [])

  useEffect(() => {
    getList()
  }, [getList])

  return (
    <>
      <Spin tip='内容正在加载。。。' spinning={loading}>
        <div className='new-home'>
          {menuList && menuList.length > 0 ?
            <>
              <div className='new-home-left'>
                {menuList.map((m: any) => {
                  return <div>
                    <div key={m.id} className={selectedId === m.id ? "leftmenu-div selected" : "leftmenu-div"}
                      onClick={async () => {
                        setSelectId(m.id)
                        await window.teamworkSDK.store.set('appId', m.id)
                      }}
                    >
                      <div className='icon-left'>
                        <Image src={m.icon} preview={false} width={20} />
                      </div>
                      <div className='icon-name'>{m.name}</div>
                    </div>
                    <div className="addAppType">
                      {/* <Button type='link' className="type-btn-new" icon={<PlusOutlined />} onClick={() => fnsRef.current.show()}>添加应用</Button> */}
                    </div>
                  </div>
                })}
              </div>
              <AppDetail selectedId={selectedId} firstId={firstId} />
            </>
            :
            <div className="add-btn">
              {/* <Button type='primary' className="add-btn-item-new" icon={<PlusOutlined />} onClick={() => fnsRef.current.show()}>新增应用</Button> */}
            </div>
          }
          {/* <AddNewApp fns={fnsRef} finished={async (name) => addnewAppFns(name)} /> */}
        </div>
      </Spin>
    </>
  )
}

export default Home