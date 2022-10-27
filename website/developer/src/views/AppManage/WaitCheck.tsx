import { ProColumns } from "@ant-design/pro-table"
import { Button, Popconfirm } from "antd"
import moment from "moment"
import { ProForm } from '@ant-design/pro-form'
import { useRef, useMemo, useCallback, useState } from "react"
import { AntProTableWrapper } from "../../components/Table"
import CheckModal from "./CheckModal"
import { IconPro } from "../../components/Icons"
import './index.less'

interface AppType {

}

interface IWaitCheckProps {
  tabType: 'waitCheck' | 'allApps',
  checkFn: () => void
}

const appStatusEnum: any = {
  1: { text: '下架', status: 'Error' },
  2: { text: '待审核', status: 'Processing' },
  3: { text: '审核通过', status: 'Error' },
  4: { text: '审核拒绝', status: 'Error' },
  5: { text: '正常', status: 'Success' },
}

const typesEnum: any = {
  1: '类型1',
  2: '类型2',
  3: '类型3'
}

const WaitCheck: React.FC<IWaitCheckProps> = (props) => {
  const actionRef = useRef<any>()
  const fnsRef = useRef<any>()
  const [ifCheck, setIfCheck] = useState<'check' | 'noCheck'>('noCheck')

  const columns = useMemo<ProColumns<AppType>[]>(
    () => [
      {
        title: '序号', key: 'index', align: 'center', dataIndex: 'index', renderFormItem: () => { return null }, render: (node, row, index) => {
          return <>{index + 1}</>
        }
      },
      { title: '应用名称', align: 'center', key: 'name', dataIndex: 'name' },
      { title: '类别名称', align: 'center', key: 'holders', dataIndex: 'holders', valueEnum: typesEnum },
      { title: '持有者', align: 'center', key: 'typeName', dataIndex: 'typeName' },
      { title: '应用类型', align: 'center', key: 'apptype', dataIndex: 'apptype' },
      { title: '图标', align: 'center', key: 'icon', dataIndex: 'icon', renderFormItem: () => { return null }, },
      { title: '短描述', align: 'center', key: 'shortDesc', dataIndex: 'shortDesc', renderFormItem: () => { return null }, },
      { title: '应用状态', align: 'center', key: 'appStatus', dataIndex: 'appStatus', renderFormItem: () => { return null }, valueEnum: appStatusEnum },
      { title: '是否上推荐', align: 'center', key: 'ifRecommed', dataIndex: 'ifRecommed', renderFormItem: () => { return null }, },
      { title: '是否在推荐页隐藏', align: 'center', key: 'ifNoDisplay', dataIndex: 'ifNoDisplay', renderFormItem: () => { return null }, },
      {
        title: '操作', width: '15%', align: 'center', key: 'option', dataIndex: 'option',
        renderFormItem: () => { return null },
        render: (node, row) => {
          return <>
            {props.tabType === 'waitCheck' && <>
              <Button type="link" onClick={() => {
                props.checkFn()
                // setI fCheck('check')
              }}>审核</Button>
            </>}
            {props.tabType === 'allApps' && <>
              <Popconfirm title='是否确定强制下架？'>
                <Button type="link">强制下架</Button>
              </Popconfirm>
              <Popconfirm title='是否确定强制整改'>
                <Button type="link">强制整改</Button>
              </Popconfirm>
            </>}
          </>
        }
      },
    ], [])

  const warpperReq = useCallback(() => {
    if (props.tabType === 'waitCheck') {
      return [
        { name: '测试', appStatus: '2' }
      ]
    }
    return [
      { name: '测试', appStatus: '3' },
      { name: '测试', appStatus: '5' },
    ]
  }, [])
  return (
    <>
      {ifCheck === 'noCheck' && <AntProTableWrapper
        rowKey='managType'
        actionRef={actionRef}
        wrapperColumns={columns}
        wrapperRequest={warpperReq}
      />}
      {ifCheck === 'check' && <>
        <Button className='cab' type='default' onClick={() => setIfCheck('noCheck')} icon={<IconPro style={{ fontSize: 26 }} type='icon-fanhui' />}>返回列表</Button>
        <CheckModal />
      </>}
    </>
  )
}

export default WaitCheck