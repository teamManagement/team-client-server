import { FC, useCallback, useEffect, useRef, useState } from "react";
import { Button, Divider, Tooltip } from 'antd';
import { EditOutlined, CheckOutlined, CloseOutlined } from '@ant-design/icons'
import { ProForm, ProFormSelect, ProFormSwitch, ProFormDependency, ProFormText } from '@ant-design/pro-form';
import IfEditWar from "./IfEdit";
import ImageBarWar from "./ImageBar";

interface IControlProps {
  getFileMenuList: any,
  getId: any
}

const ControlPanel: FC<IControlProps> = (props) => {
  const [ifEditSet, setIfEditSet] = useState<boolean>(false)
  const [nowSetEditTxt, setNowSetEditTxt] = useState<any>()
  const formRef = useRef<any>()
  const [appName, setAppName] = useState<any>('')

  const updateName = useCallback(async (newName: any) => {
    const appInfo: any = await window.teamworkSDK.store.get(props.getId)
    const newAppInfo = { ...appInfo, title: newName }
    await window.teamworkSDK.store.set(props.getId, newAppInfo)
    const appInfos: any = await window.teamworkSDK.store.get(props.getId)
    const allList: any = await window.teamworkSDK.store.get('_content_menu_list')
    const list = allList.map((m: any) => m.id === appInfos.id ? appInfos : m)
    await window.teamworkSDK.store.set('_content_menu_list', list)
    props.getFileMenuList()
  }, [props])

  const getReord = useCallback(async () => {
    const appInfo: any = await window.teamworkSDK.store.get(props.getId)
    setAppName(appInfo.title)
  }, [props])

  useEffect(() => {
    getReord()
  }, [getReord])

  const setOptions = {
    'long': {
      text: '远程地址',
      disabled: false
    },
    'local': {
      text: '本地目录 (暂不支持)',
      disabled: true,
    },
  }

  return (
    <>
      <div className="title">
        <div style={{ width: 120, marginTop: 20, height: 120 }} >
          <ImageBarWar type='single' getFileMenuList={props.getFileMenuList} getId={props.getId} />
        </div>
        <div className="right">
          <IfEditWar editText={appName} type='input' finished={(name) => updateName(name)} />
          <div className='item'>开发者: 超级管理员</div>
          <IfEditWar editText='极佳的用户体验法大师傅士大夫描述' type='area' finished={(name) => updateName(name)} />
        </div>
      </div>
      <div className="setting">
        <ProForm formRef={formRef}>
          <Divider orientation='left'>综合功能操作区</Divider>
          <div className="control-item">当前状态
            <ProFormDependency name={['ifDesk']}>
              {({ ifDesk }) => <div className="item-right" style={{ color: ifDesk ? 'green' : 'orange' }}> {ifDesk ? '调试中' : '未发布'} </div>}
            </ProFormDependency>
          </div>
          <div className="control-item">类别设置
            <div className="item-right">
              <ProFormSelect showSearch width={244} name='typeSet' options={[{ label: '测试', value: 'ceshi' }]} />
            </div>
          </div>
          <div className="control-item">在应用桌面中调试
            <div className="item-right">
              <ProFormSwitch name='ifDesk' initialValue='false' checkedChildren='启用' unCheckedChildren='禁用' />
            </div>
          </div>
          <div className="control-item">应用设置
            <div className="item-right">
              <Button type='primary' danger>删除/下架</Button>
            </div>
          </div>
          <ProFormDependency name={['ifDesk']}>
            {({ ifDesk }) => {
              if (ifDesk) {
                return <>
                  <Divider orientation='left'>调试操作区</Divider>
                  <div className="control-item">调试类别
                    <div className="item-right">
                      <ProFormSelect showSearch width={244} name='setting' initialValue='long' valueEnum={setOptions} />
                    </div>
                  </div>
                  <div className="control-item">
                    <div>
                      远程HTTP地址
                      <span style={{ marginLeft: 10 }}>
                        <Tooltip title='编辑'>
                          <Button type='link' onClick={() => setIfEditSet(true)}><EditOutlined /></Button>
                        </Tooltip>
                      </span>
                    </div>
                    {!ifEditSet && <Button type="link">{nowSetEditTxt ? nowSetEditTxt : 'http://apps.byzk.cn'}</Button>}
                    {ifEditSet && <div className="item-right" style={{ display: 'flex', alignItems: 'baseline' }}>
                      <ProFormText initialValue={nowSetEditTxt} width={178} name='address' placeholder='输入地址'
                        fieldProps={{
                          onChange: (e) => setNowSetEditTxt(e.target.value),
                          onClick: (event) => event.stopPropagation()
                        }}
                      />
                      <span style={{ margin: 10 }}>
                        <a><CheckOutlined onClick={(event) => { setIfEditSet(false); event.stopPropagation(); }} /></a>
                        <Divider type='vertical' />
                        <a><CloseOutlined onClick={(event) => { setIfEditSet(false); event.stopPropagation(); }} /></a>
                      </span>
                    </div>}
                  </div>
                </>
              }
            }}
          </ProFormDependency>
        </ProForm>
      </div>
    </>
  )
}


export default ControlPanel