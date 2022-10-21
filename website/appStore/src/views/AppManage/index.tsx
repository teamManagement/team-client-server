import { CustomTabs, TabsHeader } from "../../components/Tabs"
import WaitCheck from "./WaitCheck"

const headers: TabsHeader[] = [
  { key: 'waitCheck', title: '待审核应用' },
  { key: 'allApps', title: '所有应用' },
]

const AppManage: React.FC = () => {
  return (
    <>
      <CustomTabs headers={headers}>
        <WaitCheck key='waitCheck' tabType='waitCheck' />
        <WaitCheck key='allApps' tabType='allApps' />
      </CustomTabs>
    </>
  )
}

export default AppManage