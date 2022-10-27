import { Divider, Image, Tabs } from 'antd';
import { useState } from 'react'
import SwiperDeatil from './Swaiper';
import './index.less'

const Home: React.FC = () => {
  const { TabPane } = Tabs
  const [tabkey, setTabkey] = useState<any>('introduce')
  const [visible, setVisible] = useState<boolean>(false)
  const [versionList, setVersionList] = useState<any[]>([1, 2, 4, 5, 6, 7, 8, 9, 10, 11])


  const list = versionList.map((m: any) => {
    return <div>
      <div>1.2.{m}</div>
      <div>1.描述</div>
      <div>2.描述</div>
      <Divider />
    </div>
  })

  return (
    <>
      <div className='detail'>
        <div className='title'>
          <div><Image width={140} src='https://127.0.0.1:65528/icons/appstore.png' /></div>
          <div className='right'>
            <div className='item'>Markdown笔记  v1.8.6</div>
            <div className='item'>版本：1.8.6</div>
            <div className='item'>开发者:哈哈哈哈哈信息科技有限公司</div>
            <div className='item'>极佳的用户体验法大师傅士大夫描述</div>
          </div>
        </div>
        <div className='content-div'>
          <Tabs activeKey={tabkey} onChange={(e) => setTabkey(e)}>
            <TabPane tab='详情介绍' key='introduce'>
              <div className='tab content'><SwiperDeatil /></div>
            </TabPane>
            <TabPane tab='版本记录' key='version'>
              <div className='tab'>{list}</div>
            </TabPane>
          </Tabs>
        </div>
      </div>
    </>
  )
}

export default Home