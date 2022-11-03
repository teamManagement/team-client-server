import { FC, useCallback, useEffect, useState } from "react";
import { Divider, Tabs } from 'antd';
import Introduction from "./Introduction";
import ControlPanel from "./ControlPanel";

interface IRightProps {
  record: any,
  flushFileList: () => void,
}

const LeftMenuList: FC<IRightProps> = (props) => {

  const { TabPane } = Tabs
  const [tabkey, setTabkey] = useState<any>('mark')
  const [versionList, setVersionList] = useState<any[]>([1, 2, 4, 5, 6, 7, 8, 9, 10, 11])

  const list = versionList.map((m: any) => {
    return <div>
      <div>1.2.{m}</div>
      <div>1.描述</div>
      <div>2.描述</div>
      <Divider />
    </div>
  })


  const getFileMenuList = useCallback(async () => {
    console.log(props.record);
    props.flushFileList()
  }, [])

  useEffect(() => {
    getFileMenuList()
  }, [getFileMenuList])


  return (
    <>
      <div className='content-div'>
        <Tabs activeKey={tabkey} onChange={(e) => setTabkey(e)}>
          <TabPane tab='控制面板' key='mark'>
            <div className='tab'>
              <ControlPanel getId={props.record} getFileMenuList={props.flushFileList} />
            </div>
          </TabPane>
          <TabPane tab='详情介绍' key='introduce'>
            <div className='tab content'>
              <Introduction getFileMenuList={props.flushFileList} getId={props.record} />
            </div>
          </TabPane>
          <TabPane tab='版本记录' key='version'>
            <div className='tab'>{list}</div>
          </TabPane>
        </Tabs>
      </div>
    </>
  )
}

export default LeftMenuList