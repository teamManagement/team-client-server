import { useEffect, useRef, useState } from 'react'
import { Image, Tabs, Divider, Carousel, List, Avatar, Button, message, Modal } from 'antd'
import { CloudUploadOutlined, CloudDownloadOutlined } from '@ant-design/icons'
import './index.less'
import { applications } from '@byzk/teamwork-inside-sdk'


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


  const data = [
    { title: '苏林鑫' },
    { title: '柴敏' },
  ];

  useEffect(() => {
    if (!instalList || !appInfo) { return }
    const installs = instalList?.filter((m: any) => m === appInfo.id)
    setSelectId(installs[0])
  }, [])

  return (
    <>
      {visible && <div className='detail'>
        <div className='title'>
          <div><Image width={140} src={appInfo?.icon} /></div>
          <div className='right'>
            <div className='item' style={{ fontSize: 16, fontWeight: 'bold' }}>{appInfo?.name}</div>
            <div className='item'>版本: {appInfo?.version}</div>
            <div className='item'>开发者: 柴哈哈</div>
            <div className='item other-item' dangerouslySetInnerHTML={{ __html: '描述: ' + appInfo?.desc }} />
          </div>
          <div>
            <Button danger={appInfo.id === selectId ? true : false} icon={appInfo.id === selectId ? <CloudDownloadOutlined /> : <CloudUploadOutlined />} type='primary'
              onClick={async () => {
                try {
                  if (appInfo.id === selectId) {
                    await applications.uninstall(appInfo.id)
                  } else {
                    await applications.install(appInfo.id)
                  }
                  finished()
                  message.success(appInfo.id === selectId ? '卸载成功' : '安装成功')
                } catch (e: any) {
                  Modal.error({ title: e.message, okText: '知道了' })
                }
              }}>{selectId && appInfo.id === selectId ? '卸载' : '安装'}</Button>
          </div>
        </div>
        <div className='content-div'>
          <Tabs activeKey={tabkey} onChange={(e) => setTabkey(e)}>
            <TabPane tab='详情介绍' key='introduce'>
              <div className='tab content'>
                <div className='swiper-div'>
                  <Carousel className='swiper' autoplay dots ref={swiperRef}>
                    {JSON.parse(appInfo?.slideshow).map((value: any) =>
                      <Image height={300} src={value} />)}
                  </Carousel>
                  <div className="swiper-button-prev" onClick={() => swiperRef.current?.prev()} />
                  <div className="swiper-button-next" onClick={() => swiperRef.current?.next()} />
                </div>
                <div className='desc-div'>
                  <h1 style={{ fontSize: 16, fontWeight: 'bold' }}>长描述: </h1>
                  <div dangerouslySetInnerHTML={{ __html: appInfo?.desc }} />
                </div>
              </div>
            </TabPane>
            <TabPane tab='版本记录' key='version'>
              <div className='tab' >{list}</div>
            </TabPane>
            <TabPane tab='贡献人员' key='personList'>
              <div className='tab' >
                <List
                  itemLayout="horizontal"
                  dataSource={data}
                  renderItem={item => (
                    <List.Item>
                      <List.Item.Meta
                        avatar={<Avatar src={appInfo?.icon} />}
                        title={<Button type='link'>{item.title}</Button>}
                        description="主要开发人员"
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