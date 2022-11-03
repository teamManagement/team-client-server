import { Button, Divider, Menu, Tabs, Image } from 'antd';
import { useCallback, useEffect, useRef, useState } from 'react'
import { useParams } from "react-router-dom"
import Introduction from './Introduction';
import ControlPanel from './ControlPanel';
import './index.less'
import LeftMenuList from './RightMenuList';
import AddNewApp from './addNewApp';
import { PlusOutlined } from "@ant-design/icons";
import RightMenuList from './RightMenuList';

function uuid() {
  return 'xxxxxxxx-xxxx-4xxx-yxxx-xxxxxxxxxxxx'.replace(/[xy]/g, function (c) {
    var r = Math.random() * 16 | 0,
      v = c == 'x' ? r : (r & 0x3 | 0x8);
    return v.toString(16);
  });
}

const initImage = 'https://127.0.0.1:65528/icons/undefined.png'

interface MenuInfo {
  id: string,
  title: string,
  icon: string,
}

const Home: React.FC = () => {
  const [menuList, setMenuList] = useState<MenuInfo[]>([])
  const fnsRef = useRef<any>()
  const [selectedId, setSelectId] = useState<any>()

  // window.teamworkSDK.store.set('params',params).then(async ()=>{
  //   const data =  await window.teamworkSDK.store.get<{[key:string]: string}>('params')
  //   console.log("store: ", data)
  // })


  const getFileList = useCallback(async () => {
    const data: any = await window.teamworkSDK.store.get("_content_menu_list")
    // const data: any = sessionStorage.getItem("_content_menu_list")
    console.log("store: ", data)
    // const list: any = await window.teamworkSDK.store.get<{ [key: string]: string }>('image-single')
    // console.log(list);
    setSelectId(data[0]?.id)
    setMenuList(data)
    window.teamworkSDK.store.set('appId', data[0]?.id)
    // setMenuList(JSON.parse(data))
  }, [])

  const flushFileList = useCallback(async () => {
    const data: any = await window.teamworkSDK.store.get("_content_menu_list")
    console.log(data);
    setMenuList(data)
  }, [])

  useEffect(() => {
    getFileList()
  }, [getFileList])


  const addNewApp = useCallback(async () => {
    let newMenuList = [...menuList]
    const addInfo = { id: uuid(), title: '新应用', icon: initImage, }
    newMenuList.push(addInfo)
    await window.teamworkSDK.store.set("_content_menu_list", newMenuList)
    await window.teamworkSDK.store.set(addInfo.id, addInfo)
    // sessionStorage.setItem("_content_menu_list", JSON.stringify(newMenuList))
    // sessionStorage.setItem(addInfo.id, JSON.stringify(addInfo))
    setSelectId(addInfo.id)
    setMenuList(newMenuList);
  }, [menuList])

  return (
    <>
      <div className='home'>
        {menuList && menuList.length > 0 ?
          <>
            <div className='left'>
              {menuList.map((m: any) => {
                return <div>
                  <div key={m.id} className={selectedId === m.id ? "leftmenu-div selected" : "leftmenu-div"}
                    onClick={async () => {
                      setSelectId(m.id)
                      await window.teamworkSDK.store.set('appId', m.id)
                    }}
                  >
                    <div className='icon-left'>
                      <Image src={m.icon} preview={false} width={50} />
                    </div>
                    <div className='icon-name'>{m.title}</div>
                  </div>
                  <div className="addAppType">
                    <Button type='link' className="type-btn-new" icon={<PlusOutlined />} onClick={() => fnsRef.current.show()}>添加应用</Button>
                  </div>
                </div>
              })}
            </div>
            <div className='detail'>
              <RightMenuList flushFileList={flushFileList} record={selectedId} />
            </div>
          </>
          : <div className="add-btn">
            <Button type='primary' className="add-btn-item-new" icon={<PlusOutlined />} onClick={() => fnsRef.current.show()}>新增应用</Button>
          </div>
        }
      </div>
      <AddNewApp fns={fnsRef} finished={async (name) => {
        let newMenuList = [...menuList]
        const addInfo = { id: uuid(), title: `${name}`, icon: initImage, }
        newMenuList.push(addInfo)
        // setMenuList([addInfo])
        // sessionStorage.setItem("_content_menu_list", JSON.stringify([{ id: uuid(), title: `${name}`, icon: initImage }]))
        // window.teamworkSDK.store.set("_content_menu_list", [{ id: uuid(), title: `${name}`, icon: initImage }])
        await window.teamworkSDK.store.set("_content_menu_list", newMenuList)
        await window.teamworkSDK.store.set(addInfo.id, addInfo)
        // getFileList()
        setSelectId(addInfo.id)
        setMenuList(newMenuList);
      }} />
    </>
  )
}

export default Home