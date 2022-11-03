import { Swiper, SwiperSlide } from 'swiper/react';
import { Image } from 'antd'
import 'swiper/swiper-bundle.css';
import SwiperCore, { Autoplay } from 'swiper';
import './index.less'
import 'swiper/css';
import React from 'react';
SwiperCore.use([Autoplay]);

const SwiperDeatil: React.FC = () => {
  // 轮播图最多九张 要是0张 不显示轮播图
  const partnerLogo: Array<string> = [
    require('../../imgs/1.png'),
    require('../../imgs/2.png'),
    require('../../imgs/3.png'),
    require('../../imgs/4.png'),
    require('../../imgs/5.png'),
    require('../../imgs/6.png'),
    require('../../imgs/7.png'),
    require('../../imgs/8.png'),
    require('../../imgs/9.png'),
  ];

  return (
    <>
      <div className='swiper'>
        <Swiper
          loop
          autoplay={{ delay: 2000 }}
          pagination
          preventClicks
          navigation
          onSlideChange={() => console.log('slide change')}
          onSwiper={(swiper) => console.log(swiper)}
        >
          {partnerLogo.map((value, index) => {
            return (
              <SwiperSlide key={index}>
                <Image height={400} src={value} />
              </SwiperSlide>
            );
          })}
          <div className='swiper-pagination'></div>
        </Swiper>
      </div>

      {/* 渲染富文本 */}
      <div className='desc-div'>
        <div>【Markdown 笔记特性】</div>
        <div style={{ marginTop: 20 }}>1. 实时预览、存储</div>
        <div>2. 与传统富文本编辑方式结合</div>
        <div>3. 直接粘贴截图或图片文件加入到笔记内容，图片存储在数据库文档中(云端同步不丢失)</div>
        <div>4. 代码块支持 168 种语言</div>
        <div>5. 支持 TODO 任务列表</div>
        <div>6. 支持数学公式</div>
        {/* <div>7. 可以导出 Markdown、Html、PDF、图片文件</div>
        <div>8. 分离多个窗口，同时编辑多个笔记</div>
        <div>9. 快速搜索笔记内容</div>
        <div>10. 笔记加入 uTools 搜索， 直接打开</div> */}
      </div>

      <div className='preClick' onClick={() => {
        console.log(121);
      }}></div>
      <div className='nextClick' onClick={() => {
        console.log(456);
      }}></div>
    </>
  )
}

export default SwiperDeatil

