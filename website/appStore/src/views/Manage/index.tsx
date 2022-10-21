import { ProColumns } from "@ant-design/pro-table"
import { Button, Popconfirm } from "antd"
import moment from "moment"
import { useCallback, useMemo, useRef } from "react"
import { AntProTableWrapper } from "../../components/Table"
import AddTypes from "./AddTypes"
import './index.less'

interface ManType {
  createTime: string;
}

const Manage: React.FC = () => {
  const actionRef = useRef<any>()
  const fnsRef = useRef<any>()

  const columns = useMemo<ProColumns<ManType>[]>(
    () => [
      {
        title: '序号', key: 'index', dataIndex: 'index', renderFormItem: () => { return null }, render: (node, row, index) => {
          return <>{index + 1}</>
        }
      },
      { title: '名称', key: 'name', dataIndex: 'name' },
      { title: '图标', key: 'icon', dataIndex: 'icon' },
      { title: '创建时间', key: 'createTime', dataIndex: 'createTime', render: (node, row) => <div>{moment(parseInt(row.createTime)).format('YYYY-MM-DD HH:mm:ss')}</div> },
      { title: '最后修改时间', key: 'lastEditTime', dataIndex: 'lastEditTime' },
      {
        title: '操作', key: 'option', dataIndex: 'option', render: (node, row) => {
          return <>
            <Button type='link' onClick={() => fnsRef.current.show(row)}>修改</Button>
            <Popconfirm
              title='是否确认删除？'
              onConfirm={() => { }}
            >
              <Button type='link' danger>删除</Button>
            </Popconfirm>
          </>
        }
      },
    ], [])

  const warpperReq = useCallback(() => {
    return [
      { name: '测试' }
    ]
  }, [])
  return (
    <>
      <Button className="addbtn" type='primary' onClick={() => fnsRef.current.show()}>新增类别</Button>
      <AntProTableWrapper
        rowKey='managType'
        actionRef={actionRef}
        wrapperSearchDisabled
        wrapperColumns={columns}
        wrapperRequest={warpperReq}
      />
      <AddTypes fns={fnsRef} finnished={() => { }} />
    </>
  )
}

export default Manage