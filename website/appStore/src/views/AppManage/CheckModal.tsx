import { Button, Descriptions, Divider, Modal, Space } from "antd"
import { useEffect, useRef, useState } from "react"
import { Swiper, SwiperSlide } from 'swiper/react';
import 'swiper/swiper-bundle.css';
import { Image } from 'antd'
import SwiperCore, { Autoplay } from 'swiper';
import RejectModal from "./RejectModal";
SwiperCore.use([Autoplay]);

interface ICheckProps {
  fns: any,
  finished: () => void
}

const CheckModal: React.FC = () => {

  const fnsRef = useRef<any>()
  const partnerLogo: Array<string> = [
    require('../../imgs/1.png'),
    require('../../imgs/2.png'),
    require('../../imgs/3.png'),
  ];

  return (
    <>
      <Divider orientation='left' type='horizontal' >审核信息</Divider>
      <Descriptions column={3} style={{ marginBottom: 40 }}>
        <Descriptions.Item label='名称'>{'-'}</Descriptions.Item>
        <Descriptions.Item label='类别'>{'-'}</Descriptions.Item>
        <Descriptions.Item label='应用版本'>{'-'}</Descriptions.Item>
      </Descriptions>
      <Descriptions column={3} style={{ marginBottom: 40 }}>
        <Descriptions.Item label='贡献者列表'>{'-'}</Descriptions.Item>
        {/* 远程web(项目地址)/本地web(项目文件 文件hash 可下载)  */}
        <Descriptions.Item label='应用类型'>{'-'}</Descriptions.Item>
        <Descriptions.Item label='图标'>{'-'}</Descriptions.Item>
      </Descriptions>
      <Descriptions column={3} style={{ marginBottom: 40 }}>
        <Descriptions.Item label='所需权限'>{'-'}</Descriptions.Item>
        <Descriptions.Item label='短描述'>{'-'}</Descriptions.Item>
      </Descriptions>
      <div style={{display:"flex",flexDirection:'row'}}>
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
        {/* 否 版本历史 */}
        <Descriptions.Item label='是否为首次发布'>{'-'}</Descriptions.Item>
        <Descriptions.Item label='是否需要推荐位'>{'-'}</Descriptions.Item>
        <Descriptions.Item label='是否在应用商店内隐藏'>{'-'}</Descriptions.Item>
      </Descriptions>
      <Descriptions column={1} style={{ marginBottom: 40 }}>
        {/* 转为markdown */}
        {/* <Descriptions.Item label='长描述'><div dangerouslySetInnerHTML={{__html:'<h1>32432</h1>'}} style={{ height: 150, width: '70vw', border: '1px solid #999' }}>markdown</div></Descriptions.Item> */}
        <Descriptions.Item label='长描述'><div  style={{ height: 150, width: '70vw', border: '1px solid #999' }}>markdown</div></Descriptions.Item>
      </Descriptions>

      <Button.Group className="btns-group">
        <Space>
          <Button type='primary' danger onClick={() => fnsRef.current.show()}>审核拒绝</Button>
          <Button type='primary'>审核通过</Button>
        </Space>
      </Button.Group>
      <RejectModal fns={fnsRef} finished={() => { }} />
    </>
  )
}

export default CheckModal