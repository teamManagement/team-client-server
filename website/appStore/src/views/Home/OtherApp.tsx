import { FC, useCallback, useEffect, useState } from "react";
import Detail from "./Detail";
import { Divider, Image, Tabs } from 'antd';
import './otherApp.less';
import SwiperDeatil from "./Swaiper";

interface IOtherProps {
  getId: any,
}

const OtherApp: FC<IOtherProps> = (props) => {
  const { TabPane } = Tabs
  const [tabkey, setTabkey] = useState<any>('introduce')
  const [visible, setVisible] = useState<boolean>(false)
  const [versionList, setVersionList] = useState<any[]>([1, 2, 4, 5, 6, 7, 8, 9, 10, 11])
  const [appName, setAppName] = useState<string>()

  const list = versionList.map((m: any) => {
    return <div>
      <div>1.2.{m}</div>
      <div>1.描述</div>
      <div>2.描述</div>
      <Divider />
    </div>
  })

  console.log(props.getId);
  const getReord = useCallback(async () => {
    const appInfo: any = await window.teamworkSDK.store.get(props.getId)
    setAppName(appInfo.title)
  }, [props])

  useEffect(() => {
    getReord()
  }, [getReord])

  return (
    <>
      <div className="otherApp">
        <div className='detail'>
          <div className='title'>
            <div><Image width={140} src={''} /></div>
            <div className='right'>
              <div className='item'>{appName}</div>
            </div>
          </div>
        </div>
      </div>
    </>
  )
}

export default OtherApp