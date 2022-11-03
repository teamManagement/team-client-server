import { Button, Divider, Image } from "antd";
import { FC, useCallback, useEffect, useRef, useState } from "react";
import { IconPro } from "../../components/Icons";
import Detail from "./Detail";


interface IAppDetailProps {
  selectedId: any,
  firstId: any,
}


const AppDetail: FC<IAppDetailProps> = (props) => {

  const [list, setList] = useState<any[]>([1, 2, 3, 4, 5, 6, 7, 8, 4, 3, 2, 2, 4])
  const [ifDetail, setIfDetail] = useState<boolean>(false)
  const fnsRef = useRef<any>()
  const [appName, setAppName] = useState<string>()
  const [showOne, setShowOne] = useState<boolean>(true)
  const [showTwo, setShowTwo] = useState<boolean>(true)

  const divList = list.map((m) => {
    return <div className="small-div">
      <div className="list-div" onClick={() => {
        fnsRef.current.show()
        setIfDetail(true)
      }}>
        <div className="left"><Image width={40} preview={false} src={'https://127.0.0.1:65528/icons/appstore.png'} /></div>
        <div className="content">
          <div className='title'>名称</div>
          <div className='desc'>描述哈哈哈哈哈哈但还是觉得辐射防护描述哈哈哈哈哈哈但还是觉得辐射防护描述哈哈哈哈哈哈但还是觉得辐射防护描述哈哈哈哈哈哈但还是觉得辐射防护描述哈哈哈哈哈哈但还是觉得辐射防护描述哈哈哈哈哈哈但还是觉得辐射防护描述哈哈哈哈哈哈但还是觉得辐射防护描述哈哈哈哈哈哈但还是觉得辐射防护描述哈哈哈哈哈哈但还是觉得辐射防护描述哈哈哈哈哈哈但还是觉得辐射防护</div>
        </div>
      </div>
      <div className="footer"><Button type='primary'>安装</Button></div>
    </div>
    // return <div className='bigDiv'>
    //   <div className='middleSDiv'>
    // <div className='smalldiv' onClick={() => {
    //   fnsRef.current.show()
    //   setIfDetail(true)
    // }}>
    //       <div className='iconStyle'><IconPro style={{ width: 50 }} type='icon-yunyingpan' /></div>
    //       <div className='content'>
    // <div className='title'>名称</div>
    // <div className='desc'>描述哈哈哈哈哈哈但还是觉得辐射防护描述哈哈哈哈哈哈但还是觉得辐射防护描述哈哈哈哈哈哈但还是觉得辐射防护描述哈哈哈哈哈哈但还是觉得辐射防护描述哈哈哈哈哈哈但还是觉得辐射防护描述哈哈哈哈哈哈但还是觉得辐射防护描述哈哈哈哈哈哈但还是觉得辐射防护描述哈哈哈哈哈哈但还是觉得辐射防护描述哈哈哈哈哈哈但还是觉得辐射防护描述哈哈哈哈哈哈但还是觉得辐射防护</div>
    //       </div>
    //       <div><Button>安装</Button></div>
    //     </div>
    //   </div>
    //   <div className='middleSDiv'>
    //     <div className='smalldiv' onClick={() => {
    //       fnsRef.current.show()
    //       setIfDetail(true)
    //     }}>
    //       <div className='iconStyle'><IconPro style={{ width: 50 }} type='icon-yunyingpan' /></div>
    //       <div className='content'>
    //         <div className='title'>名称</div>
    //         <div className='desc'>描述哈哈哈哈哈哈但还是觉得辐射防护描述哈哈哈哈哈哈但还是觉得辐射防护描述哈哈哈哈哈哈但还是觉得辐射防护描述哈哈哈哈哈哈但还是觉得辐射防护描述哈哈哈哈哈哈但还是觉得辐射防护描述哈哈哈哈哈哈但还是觉得辐射防护描述哈哈哈哈哈哈但还是觉得辐射防护描述哈哈哈哈哈哈但还是觉得辐射防护描述哈哈哈哈哈哈但还是觉得辐射防护描述哈哈哈哈哈哈但还是觉得辐射防护</div>
    //       </div>
    //       <div><Button>安装</Button></div>
    //     </div>
    //   </div>
    // </div>
  }, [])


  const getReord = useCallback(async () => {
    const appInfo: any = await window.teamworkSDK.store.get(props.selectedId)
    setAppName(appInfo.name)
  }, [props])

  useEffect(() => {
    getReord()
  }, [getReord])

  return (
    <>
      {!ifDetail && <div className="home">
        {showOne && <>
          <div className='div-title'>
            <div>{props.selectedId === props.firstId ? '最近更新' : appName}</div>
            {props.selectedId === props.firstId &&
              <div>{showTwo ? <Button type='link' onClick={() => setShowTwo(false)}>查看全部</Button>
                : <Button type='link' onClick={() => setShowTwo(true)} icon={<IconPro type='icon-fanhui' />}>返回</Button>}</div>
            }
          </div>
          <Divider />
          <div className="list" style={{ height: props.selectedId === props.firstId && showTwo ? '36vh' : "70vh" }}>{divList}</div>
          <div style={{ height: 100 }}></div>
        </>}
        {showTwo && props.selectedId === props.firstId ?
          <>
            <div className='div-title'>
              <div>下载最多</div>
              <div>{showOne ? <Button type='link' onClick={() => setShowOne(false)}>查看全部</Button>
                : <Button type='link' onClick={() => setShowOne(true)}><IconPro type='icon-fanhui' />返回</Button>}</div>
            </div>
            <Divider />
            <div className="list" style={{ height: showOne ? '36vh' : "70vh" }}> {divList}</div>
          </>
          :
          <></>
        }
      </div>}
      {ifDetail && <div className='callback' onClick={() => {
        setIfDetail(false)
        fnsRef.current.close()
      }}><IconPro style={{ fontSize: 26 }} type='icon-fanhui' /></div>}
      {/* <Detail fns={fnsRef} finished={() => { }} /> */}

    </>
  )
}

export default AppDetail