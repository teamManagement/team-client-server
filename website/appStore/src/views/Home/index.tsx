import { Image, Modal, Spin } from 'antd';
import { useCallback, useEffect, useState } from 'react'
import AppDetail from './AppDetail';
import { getTypeList } from '../../serve';
import { store } from '@byzk/teamwork-sdk';
import { ToolOutlined } from '@ant-design/icons';
import './index.less'


const Home: React.FC = () => {
  const [menuList, setMenuList] = useState<any[]>([{}])
  const [selectedId, setSelectId] = useState<any>()
  const [firstId, setFirstId] = useState<string>()
  const [loading, setLoading] = useState<boolean>(false)
  const [appList, setAppList] = useState<any>()


  const getFileList = useCallback(async () => {
    const data: any = await store.get("_content_menu_list")
    console.log("store: ", data)
    if (data.length === 0) { return }
    setSelectId(data[0]?.id)
    setFirstId(data[0]?.id)
    setMenuList(data)
    store.set('appId', '1')
  }, [])


  useEffect(() => {
    getFileList()
  }, [getFileList])

  const getList = useCallback(async () => {
    try {
      setLoading(true)
      const list = await getTypeList({})
      await store.set("_content_menu_list", list)
      setMenuList(list)
    } catch (e: any) {
      Modal.error({ title: e.message })
    } finally {
      setLoading(false)
    }
  }, [])

  useEffect(() => { getList() }, [getList])

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
                        <ToolOutlined />
                      </div>
                      <div className='icon-name'>{m.name}</div>
                    </div>
                    <div className="addAppType">
                      {/* <Button type='link' className="type-btn-new" icon={<PlusOutlined />} onClick={() => fnsRef.current.show()}>添加应用</Button> */}
                    </div>
                  </div>
                })}
              </div>
              <AppDetail selectedId={selectedId} firstId={firstId} appList={appList} />
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