import { Button, Divider, Menu } from "antd"
import { FC, useCallback, useEffect, useRef, useState } from "react"
import { PlusOutlined, SettingOutlined, QuestionCircleOutlined } from '@ant-design/icons'
import MonacoEditor from "react-monaco-editor";
import CustomDef from "./CustomDef";
import './index.less'

function uuid() {
  return '1xxxxxxx-xxxx-4xxx-yxxx-xxxxxxxxxxxx'.replace(/[xy]/g, function (c) {
    var r = Math.random() * 16 | 0,
      v = c == 'x' ? r : (r & 0x3 | 0x8);
    return v.toString(16);
  });
}

const leftContextMenu = window.electron?.ContextMenu.getById('leftContextMenu')
  ; (async () => {
    if (!leftContextMenu) {
      return
    }
    await leftContextMenu.clearItems()
    await leftContextMenu.appendMenuItem({ id: 'use', label: '应用此方案' })
    await leftContextMenu.appendMenuItem({ id: 'edit', label: '编辑' })
    await leftContextMenu.appendMenuItem({ id: 'delete', label: '删除' })
  })()

const HostsDemo: FC = () => {
  const [text, setText] = useState("");
  const fnsRef = useRef<any>()
  const [menuList, setMenuList] = useState<any[]>([])
  const [selectedId, setSelectId] = useState<string>('')
  const [rightId, setRightId] = useState<string>()

  const clickMenu = useCallback(async (id: any) => {
    console.log(id);
    setSelectId(id)
    await window.teamworkSDK.store.set('menuId', id)
  }, [])

  useEffect(() => {
    const newMenuList = [
      { id: uuid(), name: '开发环境' },
      { id: uuid(), name: '测试环境' },
      { id: uuid(), name: '生产环境' }
    ]
    window.teamworkSDK.store.set('menuList', newMenuList)
  }, [])

  const getMenu = useCallback(async () => {
    const list = await window.teamworkSDK.store.get<any[]>('menuList')
    const selId = await window.teamworkSDK.store.get<string>('menuId')
    setMenuList(list)
    setSelectId(selId)
  }, [])

  useEffect(() => { getMenu() }, [getMenu])


  useEffect(() => {
    leftContextMenu.registerItemClick('use', async () => {
      const rightClick = await window.teamworkSDK.store.get<string>('menuId')
      setRightId(rightClick)
    })
  }, [])

  const onContextMenuFn = useCallback(async (id: any) => {
    if (!leftContextMenu) { return }
    await window.teamworkSDK.store.set('menuId', id)
    leftContextMenu.popup()
  }, [leftContextMenu])

  const saveText = useCallback(async () => {
    // console.log(text);
    console.log(selectedId);
    console.log(menuList);

    const list = menuList.filter((m) => m.id === selectedId)

    console.log(list);

    // const coverText = "#--------- 开发环境 ------------\n\n\n" + text
    // await window.teamworkSDK.hosts.cover(coverText)
    // const newText = await window.teamworkSDK.hosts.export()
    // console.log(newText);
  }, [text])

  return (
    <>
      <div className="home">
        <div className="left">
          <div className={'lookall' === selectedId ? "left-menu-item select-item" : "left-menu-item"} id={'lookall'} onClick={() => clickMenu('lookall')}>查看系统hosts 文件内容</div>
          <div className="small">共用</div>
          <div className={'publicSet' === selectedId ? "left-inline-menu select-item" : "left-inline-menu"} id={'publicSet'} onClick={() => clickMenu('publicSet')}>
            <div className="radio">
              <div className="radio-content" />
            </div>
            公共配置</div>
          <div className="small">自定义</div>
          {menuList?.map((m: any) => {
            return <div
              className={m.id === selectedId ? "left-inline-menu select-item" : "left-inline-menu"}
              onContextMenu={() => onContextMenuFn(m.id)}
              id={m.id} onClick={() => clickMenu(m.id)}>
              <div className="radio">
                {m.id === rightId && <div className="radio-content" />}
              </div>
              {m.name}
            </div>
          })}
          <div className="left-footer">
            <div className="item" onClick={() => fnsRef.current.show()}><PlusOutlined /></div>
          </div>
        </div>
        <div className="right">
          <div className="content">
            <MonacoEditor
              theme="hc-light"
              value={text}
              onChange={(value) => { setText(value) }}
              options={{
                automaticLayout: true,
                minimap: {
                  enabled: false,
                }
              }}
            />
          </div>
          <Divider />
          <div className="footer">
            <Button type='primary' onClick={() => saveText()}>保存(CTRL+S)</Button>
          </div>
        </div>
      </div>
      <CustomDef fns={fnsRef} finished={() => { }} />
    </>
  )
}

export default HostsDemo