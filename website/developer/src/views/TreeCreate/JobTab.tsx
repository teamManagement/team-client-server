import { ProColumns } from "@ant-design/pro-table";
import { Button, message, Popconfirm } from "antd";
import { FC, useCallback, useMemo, useRef } from "react";
import { AntProTableWrapper } from "../../components/Table";
import { deleteJob, deletePost, JobList, PostList } from "../../serve";
import JobName from "./JobName";
import { RFCToFormat } from "./utils";

interface TreeType {
  id: string,
  createdAt: string,
  updatedAt: string
}

const IndexEnum = {
  title: '序号', key: 'index', dataIndex: 'index', renderFormItem: () => { return null }, render: (node: any, row: any, index: any) => {
    return <>{index + 1}</>
  }
}

interface IJobProps {
  orgId: any,
  type: 'job' | 'post'
}

const JobTab: FC<IJobProps> = (props) => {
  const { orgId, type } = props
  const actionRef = useRef<any>()
  const fnsOrgRef = useRef<any>()

  const columns = useMemo<ProColumns<TreeType>[]>(
    () => [
      IndexEnum,
      { title: '名称', key: 'name', dataIndex: 'name' },
      { title: '创建时间', key: 'createdAt', dataIndex: 'createdAt', render: (node, row) => RFCToFormat(row.createdAt) },
      { title: '最后修改时间', key: 'updatedAt', dataIndex: 'updatedAt', render: (node, row) => RFCToFormat(row.updatedAt) },
      {
        title: '操作', key: 'option', dataIndex: 'option', render: (node, row) => {
          return <>
            <Button type='link' onClick={() => {
              fnsOrgRef.current.show(row)
            }}>修改</Button>
            <Popconfirm
              title='是否确认删除？'
              onConfirm={async () => {
                if (type === 'job') {
                  await deleteJob(orgId, row.id)
                } else {
                  await deletePost(orgId, row.id)
                }
                actionRef.current.reload()
                message.success('删除成功!', 3)
              }}
              okText='确定'
              cancelText='取消'
            >
              <Button type='link' danger>删除</Button>
            </Popconfirm>
          </>
        }
      },
    ], [orgId])

  const warpperReq = useCallback(async () => {
    if (orgId) {
      let jobList: any = []
      if (type === 'job') {
        jobList = await JobList(orgId)
      } else {
        jobList = await PostList(orgId)
      }
      return jobList
    }
  }, [orgId, type])

  return (
    <>
      <Button className="addNew" type='primary' onClick={() => fnsOrgRef.current.show(orgId)}>{type === 'job' ? '新增职位' : '新增岗位'}</Button>
      <AntProTableWrapper
        rowKey='treetype'
        actionRef={actionRef}
        wrapperSearchDisabled
        wrapperColumns={columns}
        wrapperRequest={warpperReq}
      />
      <JobName type={type} fns={fnsOrgRef} finished={() => actionRef.current.reload()} />
    </>
  )
}

export default JobTab