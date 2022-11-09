import { Button, Divider, Image, message, Spin } from "antd";
import { FC, useCallback, useEffect, useRef, useState } from "react";
import { IconPro } from "../../components/Icons";
import {CloudDownloadOutlined,CloudUploadOutlined} from '@ant-design/icons';
import { getAppTypeList } from "../../serve";
import Detail from "./Detail";
import {store} from '@byzk/teamwork-sdk';
import {applications} from '@byzk/teamwork-inside-sdk';


interface IAppDetailProps {
  selectedId: any,
  firstId: any,
  appList: any
}


const AppDetail: FC<IAppDetailProps> = (props) => {

  const [list, setList] = useState<any[]>([])
  const [ifDetail, setIfDetail] = useState<boolean>(false)
  const fnsRef = useRef<any>()
  const [appName, setAppName] = useState<string>()
  const [showOne, setShowOne] = useState<boolean>(true)
  // const [showTwo, setShowTwo] = useState<boolean>(true)
  const [ifInstall, setIfInstall] = useState<any[]>([])
  const [loading, setLoading] = useState<boolean>(false)


  const getAppList = useCallback(async () => {
    if (!props.selectedId) { return }
    const list = await getAppTypeList(props.selectedId)
    if (list?.appList.length === 0) { return }
    setIfInstall(list?.appInstallationIdList)
    setList(list?.appList)
    setIfDetail(false)
    fnsRef.current.close()
  }, [props.selectedId])

  useEffect(() => { getAppList() }, [getAppList])

  const divList = list?.map((m, i) => {
    const stallId = ifInstall.filter(item => item === m.id)

    return <div className="small-div">
      <div className="list-div" onClick={() => {
        fnsRef.current.show(m)
        setIfDetail(true)
      }}>
        <div className="left"><Image width={40} preview={false} src={m.icon} /></div>
        <div className="content">
          <div className='title'>{m.name}</div>
          <div className='desc other-item' style={{ marginTop: 6 }} dangerouslySetInnerHTML={{ __html: m.desc }} />
        </div>
      </div>
      <div className="footer">
        {stallId && stallId[0] === m.id ? <Button type='primary' icon={<CloudDownloadOutlined />} danger onClick={async () => {
          setLoading(true)
          await applications.uninstall(m.id)
          getAppList()
          setLoading(false)
          message.success('卸载成功')
        }}>卸载</Button> :
          <Button type='primary' icon={<CloudUploadOutlined />} onClick={async () => {
            setLoading(true)
            await applications.install(m.id)
            getAppList()
            setLoading(false)
            message.success('安装成功')
          }}>安装</Button>
        }
      </div>
    </div>
  }, [getAppList, ifInstall])

  const getReord = useCallback(async () => {
    const appInfoList: any = await store.get('_content_menu_list')
    const newListInfo = appInfoList.filter((m: any) => m.id === props.selectedId)
    setAppName(newListInfo[0]?.name)
  }, [props.selectedId])

  useEffect(() => {
    getReord()
  }, [getReord])

  return (
    <>
      {!ifDetail && <div className="home">
        <Spin spinning={loading} tip='正在加载。。。'>
          {showOne && <>
            <div className='div-title'>
              <div style={{ fontSize: 18, fontWeight: 'bold', marginBottom: 10 }}>
                {appName}</div>
              {/* {props.selectedId === props.firstId &&
                <div>{showTwo ? <Button type='link' onClick={() => setShowTwo(false)}>查看全部</Button>
                  : <Button type='link' onClick={() => setShowTwo(true)} icon={<IconPro type='icon-fanhui' />}>返回</Button>}</div>
              } */}
            </div>
            <Divider />
            <div className="list">{divList}</div>
            <div style={{ height: 100 }}></div>
          </>}
          {/* {showTwo && props.selectedId === props.firstId ?
            <>
              <div className='div-title'>
                <div style={{ fontSize: 18, fontWeight: 'bold', marginBottom: 10 }}>下载更多</div>
                <div>{showOne ? <Button type='link' onClick={() => setShowOne(false)}>查看全部</Button>
                  : <Button type='link' onClick={() => setShowOne(true)}><IconPro type='icon-fanhui' />返回</Button>}</div>
              </div>
              <Divider />
              <div className="list" > {divList}</div>
            </>
            :
            <></>
          } */}
        </Spin>
      </div>}
      {ifDetail && <div className='callback' onClick={() => {
        setIfDetail(false)
        fnsRef.current.close()
      }}><IconPro style={{ fontSize: 26 }} type='icon-fanhui' /></div>}
      <Detail fns={fnsRef} finished={() => { }} />
    </>
  )
}

export default AppDetail