import { Button, Descriptions, Divider, Modal, Space, Tabs } from "antd"
import { useEffect, useRef, useState } from "react"
import { Swiper, SwiperSlide } from 'swiper/react';
import 'swiper/swiper-bundle.css';
import { Image } from 'antd';
import './index.less';
import SwiperCore, { Autoplay } from 'swiper';
import RejectModal from "./RejectModal";



SwiperCore.use([Autoplay]);

interface ICheckProps {
  fns: any,
  finished: () => void
}

const CheckModal: React.FC = () => {
  const { TabPane } = Tabs
  const [tabkey, setTabkey] = useState<any>('introduce')
  const fnsRef = useRef<any>()
  const [versionList, setVersionList] = useState<any[]>([1, 2, 4, 5, 6])
  const [authorList, setAuthorList] = useState<any[]>([1, 2, 4])

  const image = require('../../imgs/markdown.png')

  const list = versionList.map((m: any) => {
    return <div>
      <div>1.2.{m}</div>
      <div>1.描述</div>
      <div>2.描述</div>
      <Divider />
    </div>
  })

  const lists = authorList.map((m: any) => {
    return <div>
      <div>张三.{m}</div>
      <Divider />
    </div>
  })

  return (
    <>
      <div className='checkSign'>
        <Divider orientation='center' type='horizontal' >审核信息</Divider>
        <div className='title'>
          <div><Image width={100} src={image} /></div>
          <div className='right'>
            <div className='item'>名称：</div>
            <div className='item'>版本：</div>
            <div className='item'>短描述：</div>
          </div>
        </div>
        <div className="tabs-content">
          <Tabs activeKey={tabkey} onChange={(e) => setTabkey(e)}>
            <TabPane tab='详情介绍' key='introduce'>
              <div className="detailIntroduce">
                <div className="pics">
                  <Image width={100} src={image} />
                </div>
                <div className="longDesc">
                  markdown
                </div>
              </div>
            </TabPane>
            <TabPane tab='贡献者列表' key='contributors'>
              {lists}
            </TabPane>
            <TabPane tab='版本记录' key='version'>
              {list}
            </TabPane>
          </Tabs>
        </div>
        {/* <Descriptions column={3} style={{ marginBottom: 40 }}>
          <Descriptions.Item label='名称'>{'-'}</Descriptions.Item>
          <Descriptions.Item label='类别'>{'-'}</Descriptions.Item>
          <Descriptions.Item label='应用版本'>{'-'}</Descriptions.Item>
        </Descriptions>
        <Descriptions column={3} style={{ marginBottom: 40 }}>
          <Descriptions.Item label='贡献者列表'>{'-'}</Descriptions.Item>
           远程web(项目地址)/本地web(项目文件 文件hash 可下载)  
          <Descriptions.Item label='应用类型'>{'-'}</Descriptions.Item>
          <Descriptions.Item label='图标'>{'-'}</Descriptions.Item>
        </Descriptions>
        <Descriptions column={3} style={{ marginBottom: 40 }}>
          <Descriptions.Item label='所需权限'>{'-'}</Descriptions.Item>
          <Descriptions.Item label='短描述'>{'-'}</Descriptions.Item>
        </Descriptions>
        <div style={{ display: "flex", flexDirection: 'row' }}>
          <div>应用介绍：</div>
          <div className="war">
            <Swiper
              loop
              autoplay={{ delay: 2000 }}
              pagination={{ clickable: true }}
              preventClicks
              navigation
              onSlideChange={() => console.log('slide change')}
              onSwiper={(swiper) => console.log(swiper)}
            >
              {partnerLogo.map((value, index) => {
                return (
                  <SwiperSlide key={index}>
                    <Image height={100} src={value} />
                  </SwiperSlide>
                );
              })}
            </Swiper>
          </div>
        </div>
        <Descriptions column={3} style={{ marginBottom: 40, marginTop: 40 }}>
           否 版本历史 
          <Descriptions.Item label='是否为首次发布'>{'-'}</Descriptions.Item>
          <Descriptions.Item label='是否需要推荐位'>{'-'}</Descriptions.Item>
          <Descriptions.Item label='是否在应用商店内隐藏'>{'-'}</Descriptions.Item>
        </Descriptions>
        <Descriptions column={1} style={{ marginBottom: 40 }}>
           <Descriptions.Item label='长描述'><div dangerouslySetInnerHTML={{__html:'<h1>32432</h1>'}} style={{ height: 150, width: '70vw', border: '1px solid #999' }}>markdown</div></Descriptions.Item> 
          <Descriptions.Item label='长描述'><div style={{ height: 150, width: '70vw', border: '1px solid #999' }}>markdown</div></Descriptions.Item>
        </Descriptions> */}

        <Button.Group className="btns-group">
          <Space>
            <Button type='primary' danger onClick={() => fnsRef.current.show()}>审核拒绝</Button>
            <Button type='primary'>审核通过</Button>
          </Space>
        </Button.Group>
        <RejectModal fns={fnsRef} finished={() => { }} />
      </div>
    </>
  )
}

export default CheckModal