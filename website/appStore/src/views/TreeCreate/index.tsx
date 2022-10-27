import { Button, Divider, message, Tabs } from "antd";
import { FC, useCallback, useEffect, useRef, useState } from "react";
import './index.less';
import TreeSelectWar from "./TreeSelect";
import FirstGetName from "./FirstGetName";
import { getOrgList } from "../../serve";
import { CustomTabs, TabsHeader } from "../../components/Tabs";
import JobTab from "./JobTab";

const headers: TabsHeader[] = [
  { key: 'job', title: '职位' },
  { key: 'post', title: '岗位' },
]

const TreeCreate: FC = () => {
  const fnsRef = useRef<any>()
  const [treeList, setTreeList] = useState<any[]>([])
  const [getNewOrgId, setGetNewOrgId] = useState<string>(treeList[0]?.id)

  const loadTree = useCallback(async () => {
    try {
      const list: any = await getOrgList({})
      console.log(list);
      setTreeList(list)
    } catch (e: any) {
      message.error(e)
    }
  }, [])

  useEffect(() => {
    loadTree()
  }, [loadTree])
  
  return (
    <>
      {treeList.length === 0 && <div className="btn-add">
        <Button type='primary' onClick={() => fnsRef.current.show()}>新增机构</Button>
      </div>}
      {treeList.length > 0 && <div className="tree-create">
        <div className="left">
          <TreeSelectWar firstList={treeList} getOrgId={(orgId) => {
            console.log(orgId);
            setGetNewOrgId(orgId)
          }} />
        </div>
        <div className="line"><Divider type='vertical' /></div>
        <div className="right">
          <CustomTabs headers={headers}>
            <JobTab key='job' orgId={getNewOrgId} type='job' />
            <JobTab key='post' orgId={getNewOrgId} type='post' />
          </CustomTabs>
        </div>
      </div>}
      <FirstGetName fns={fnsRef} finished={() => { }} />
    </>
  )
}

export default TreeCreate