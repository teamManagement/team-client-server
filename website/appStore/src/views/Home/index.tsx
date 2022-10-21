import { Divider, Button } from 'antd';
import { useState, useRef } from 'react'
import { IconPro } from '../../components/Icons';
import Detail from './Detail';
import './index.less'

const Home: React.FC = () => {
  const [list, setList] = useState<any[]>([1, 2, 3])
  const [ifDetail, setIfDetail] = useState<boolean>(false)
  const fnsRef = useRef<any>()

  const divList = list.map((m) => {
    return <div className='bigDiv'>
      <div className='middleSDiv'>
        <div className='smalldiv' onClick={() => {
          fnsRef.current.show()
          setIfDetail(true)
        }}>
          <div className='iconStyle'><IconPro type='icon-yunyingpan' /></div>
          <div className='content'>
            <div className='title'>名称</div>
            <div className='desc'>描述哈哈哈哈哈哈但还是觉得辐射防护描述哈哈哈哈哈哈但还是觉得辐射防护描述哈哈哈哈哈哈但还是觉得辐射防护描述哈哈哈哈哈哈但还是觉得辐射防护描述哈哈哈哈哈哈但还是觉得辐射防护描述哈哈哈哈哈哈但还是觉得辐射防护描述哈哈哈哈哈哈但还是觉得辐射防护描述哈哈哈哈哈哈但还是觉得辐射防护描述哈哈哈哈哈哈但还是觉得辐射防护描述哈哈哈哈哈哈但还是觉得辐射防护</div>
          </div>
        </div>
        <Button>安装</Button>
      </div>
      <div className='middleSDiv'>
        <div className='smalldiv' onClick={() => {
          fnsRef.current.show()
          setIfDetail(true)
        }}>
          <div className='iconStyle'><IconPro type='icon-yunyingpan' /></div>
          <div className='content'>
            <div className='title'>名称</div>
            <div className='desc'>描述哈哈哈哈哈哈但还是觉得辐射防护描述哈哈哈哈哈哈但还是觉得辐射防护描述哈哈哈哈哈哈但还是觉得辐射防护描述哈哈哈哈哈哈但还是觉得辐射防护描述哈哈哈哈哈哈但还是觉得辐射防护描述哈哈哈哈哈哈但还是觉得辐射防护描述哈哈哈哈哈哈但还是觉得辐射防护描述哈哈哈哈哈哈但还是觉得辐射防护描述哈哈哈哈哈哈但还是觉得辐射防护描述哈哈哈哈哈哈但还是觉得辐射防护</div>
          </div>
          <Button>安装</Button>
        </div>
      </div>
    </div>
  }, [])

  return (
    <>
      {!ifDetail && <div className="home">
        <div className='div-title'>
          <div>最近更新</div>
          <div><Button type='link'>查看全部</Button></div>
        </div>
        <Divider />
        <div>{divList}</div>

        <div className='div-title two'>
          <div>下载最多</div>
          <div><Button type='link'>查看全部</Button></div>
        </div>
        <Divider />
        <div>{divList}</div>
      </div>}
      {ifDetail && <div className='callback' onClick={() => {
        setIfDetail(false)
        fnsRef.current.close()
      }}><IconPro style={{ fontSize: 26 }} type='icon-fanhui' /></div>}
      <Detail fns={fnsRef} finished={() => { }} />
    </>
  )
}

export default Home