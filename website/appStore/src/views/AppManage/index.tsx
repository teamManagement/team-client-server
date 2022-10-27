import { Button } from "antd"
import { useState } from "react"
import { IconPro } from "../../components/Icons"
import { CustomTabs, TabsHeader } from "../../components/Tabs"
import CheckModal from "./CheckModal"
import WaitCheck from "./WaitCheck"

const headers: TabsHeader[] = [
  { key: 'waitCheck', title: '待审核应用' },
  { key: 'allApps', title: '所有应用' },
]

const AppManage: React.FC = () => {

  const [ifShowCheck, setIfShowCheck] = useState<boolean>(false)
  return (
    <>
      {!ifShowCheck && <CustomTabs headers={headers}>
        <WaitCheck key='waitCheck' tabType='waitCheck' checkFn={() => {
          setIfShowCheck(true)
        }} />
        <WaitCheck key='allApps' tabType='allApps' checkFn={() => { }} />
      </CustomTabs>}
      {ifShowCheck && <>
        <Button className='cab' type='default' onClick={() => setIfShowCheck(false)} icon={<IconPro style={{ fontSize: 26 }} type='icon-fanhui' />}>返回列表</Button>
        <CheckModal />
      </>}
    </>
  )
}

export default AppManage