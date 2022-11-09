import { Button, Divider, message } from "antd"
import { FC, useCallback, useEffect, useRef, useState } from "react"
import { PlusOutlined } from '@ant-design/icons'
import MonacoEditor from "react-monaco-editor"
import { contextmenu, store, hosts } from '@byzk/teamwork-sdk'
import CustomDef from "./CustomDef"
import './index.less'


const HostsDemo: FC = () => {
  const contextMenuRef = useRef<MenuList | null>()
  const [text, setText] = useState("");
  const fnsRef = useRef<any>()
  const [menuList, setMenuList] = useState<any[]>([])
  const [selectedId, setSelectId] = useState<any>('')
  const [rightId, setRightId] = useState<any>()
  const [isDisabled, setIsDisabled] = useState<boolean>(false)

  const clickMenu = useCallback(async (id: any) => {
    setSelectId(id)
    await store.set('menuId', id)
    console.log(id);
    const list = menuList.filter((m: any) => m.id === selectedId)
    const textChange = text
    console.log(textChange);

    if (id === 'lookall') {
      setText(`#--------- ${list[0].name} ------------\n\n` + textChange)
    } else {
      setText(textChange)
    }
  }, [text, menuList])


  useEffect(() => {
    const newMenuList = [
      { id: 123, name: '开发环境', text: '' },
      { id: 213, name: '测试环境', text: '' },
      { id: 321, name: '生产环境', text: '' }
    ]
    store.set('menuList', newMenuList)
    store.set('menuId', 123)
  }, [])

  const getMenu = useCallback(async () => {
    const list: any = await store.get('menuList')
    const selId: any = await store.get('menuId')
    setMenuList(list)
    setSelectId(selId)
    setRightId(selId)
    const newText = await hosts.export()
    console.log(newText.length, text.length);
    if (newText.length === text.length) {
      setIsDisabled(true)
    }
    setText(newText)
  }, [])

  useEffect(() => { getMenu() }, [getMenu])

  const rightClickMenu = useCallback(async () => {
    const rightClick = await store.get<string>('menuId')
    setRightId(rightClick)
    saveText()
  }, [])


  useEffect(() => {
    contextMenuRef.current = contextmenu.build([
      { label: '应用此方案', click: () => rightClickMenu() },
      { label: '编辑' },
      { label: '删除' },
    ])
    return () => {
      contextmenu.clearAll()
      contextMenuRef.current = null
    }
  }, [])


  const onContextMenuWrapper = useCallback(async (id: any) => {
    await store.set('menuId', id)
    if (contextMenuRef.current) {
      contextMenuRef.current.popup()
    }
  }, [contextMenuRef])


  const saveText = useCallback(async () => {
    await hosts.cover(text)
    const newText = await hosts.export()
    console.log(newText)
    setText(newText)
    message.success('设置成功!')
  }, [text, selectedId])

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
              onContextMenu={() => onContextMenuWrapper(m.id)}
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
              onChange={async (value) => {
                const newText = await hosts.export()
                console.log(newText.length, text.length);
                if (newText.length === text.length) {
                  setIsDisabled(true)
                }
                setIsDisabled(false)
                setText(value)
              }}
              options={{
                readOnly: selectedId === 'lookall' ? true : false,
                selectionHighlight: false,
                multiCursorModifier: 'ctrlCmd',
                automaticLayout: true,
                minimap: {
                  enabled: false,
                }
              }}
            />
          </div>
          <Divider />
          <div className="footer">
            {
              selectedId === 'lookall' ?
                <div className="onlyread">只读, 无法直接编辑</div>
                : <Button type='primary' disabled={isDisabled} onClick={() => saveText()}>保存(CTRL+S)</Button>
            }

          </div>
        </div>
      </div>
      <CustomDef fns={fnsRef} finished={() => { }} />
    </>
  )
}

export default HostsDemo