import { useEffect, useRef, useState } from 'react'
import { Image, Tabs, Divider, Carousel, List, Button, message, Modal } from 'antd'
import { applications } from '@byzk/teamwork-inside-sdk'
import { AppInfo } from '@byzk/teamwork-sdk'
import { CloudDownloadOutlined, CloudUploadOutlined } from '@ant-design/icons';
import './index.less'


interface IDeatilProps {
  fns: any,
  instalList: any,
  finished: () => void,
}

const Detail: React.FC<IDeatilProps> = ({ fns, finished, instalList }) => {
  const { TabPane } = Tabs
  const [tabkey, setTabkey] = useState<any>('introduce')
  const [visible, setVisible] = useState<boolean>(false)
  const [versionList, setVersionList] = useState<any[]>([1, 2, 4, 5, 6, 7, 8, 9, 10, 11])
  const [appInfo, setAppInfo] = useState<any>()
  const [selectId, setSelectId] = useState<string>('')
  const swiperRef = useRef<any>()


  useEffect(() => {
    fns.current = {
      show(info: any) {
        if (!info) { return }
        setVisible(true)

        setAppInfo(info)
        console.log(info);

      },
      close() {
        setVisible(false)
      }
    }
  }, [])

  const list = versionList.map((m: any) => {
    return <div>
      <div>1.2.{m}</div>
      <div>1.描述</div>
      <div>2.描述</div>
      <Divider />
    </div>
  })


  useEffect(() => {
    if (!instalList || !appInfo) { return }
    console.log(instalList);

    const installs = instalList?.filter((m: any) => m === appInfo.id)
    setSelectId(installs[0])
  }, [appInfo, instalList])

  return (
    <>
      {visible && <div className='detail'>
        <div className='title'>
          <div><Image width={140} src={appInfo?.icon} /></div>
          <div className='right'>
            <div className='item' style={{ fontSize: 16, fontWeight: 'bold' }}>{appInfo?.name}</div>
            <div className='item'>版本:&nbsp;&nbsp;{appInfo?.version}</div>
            <div className='item'>开发者:&nbsp;&nbsp;teamwork</div>
            <div className='item other-item' dangerouslySetInnerHTML={{ __html: '描述: ' + appInfo?.desc }} />
          </div>
          <div>
            {appInfo?.id === selectId ?
              <Button className='newinstall' icon={<CloudDownloadOutlined />} type='primary' danger onClick={async () => {
                try {
                  await applications.uninstall(appInfo?.id)
                  finished()
                  message.success('卸载成功')
                } catch (e: any) {
                  Modal.error({ title: e.message, okText: '知道了' })
                }
              }}>卸载</Button>
              : <Button icon={<CloudUploadOutlined />} className='newinstall' type='primary' onClick={async () => {
                try {
                  await applications.install(appInfo?.id as string)
                  finished()
                  message.success('安装成功')
                } catch (e: any) {
                  Modal.error({ title: e.message, okText: '知道了' })
                }
              }}>安装</Button>}
          </div>
        </div>
        <div className='content-div'>
          <Tabs activeKey={tabkey} onChange={(e) => setTabkey(e)}>
            <TabPane tab='详情介绍' key='introduce'>
              <div className='tab content'>
                <div className='swiper-div'>
                  <Carousel className='swiper' autoplay dots ref={swiperRef}>
                    {JSON.parse(appInfo?.slideshow as any).map((value: any) =>
                      <Image height={300} src={value} />)}
                  </Carousel>
                  <div className="swiper-button-prev" onClick={() => swiperRef.current?.prev()} />
                  <div className="swiper-button-next" onClick={() => swiperRef.current?.next()} />
                </div>
                <div className='desc-div'>
                  <h1 style={{ fontSize: 16, fontWeight: 'bold' }}>长描述: </h1>
                  <div dangerouslySetInnerHTML={{ __html: appInfo?.desc as string }} />
                </div>
              </div>
            </TabPane>
            {(appInfo?.type) && <TabPane tab='版本记录' key='version'>
              <div className='tab' >{list}</div>
            </TabPane>}
            <TabPane tab='贡献人员' key='personList'>
              <div className='tab' >
                <List
                  itemLayout="horizontal"
                  dataSource={[
                    // { title: appInfo ? 'teamwork(平台)' : '', desc: appInfo && appInfo?.authorInfo?.orgList[0].org.name },
                    { title: 'teamwork(平台)', desc: '超管' },
                  ]}
                  renderItem={item => (
                    <List.Item>
                      <List.Item.Meta
                        avatar={<div className='personStyle'>{item.title.slice(item.title.length - 3, item.title.length - 1)}</div>}
                        title={<a>{item.title}</a>}
                        // description={appInfo && appInfo?.authorInfo?.id === '0' ? '平台内置' : '辅助开发人员'}
                        description={'平台内置'}
                      />
                    </List.Item>
                  )}
                />
              </div>
            </TabPane>
          </Tabs>
        </div>
      </div>}
    </>
  )
}

export default Detail