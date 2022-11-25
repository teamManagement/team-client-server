/* eslint-disable @typescript-eslint/no-unused-vars */
/* eslint-disable jsx-a11y/anchor-is-valid */
import { Button, Divider, Form, Input, Popconfirm, Table } from "antd"
import { useCallback, useRef, useState, useEffect } from "react"
import TypeEditor, { TypeEditorAction } from "./TypeEditor.tsx"
import { momentFormat } from "../../components/utils"
import './index.less';


const Manage: React.FC = () => {

  const [fetchLoading, setFetchLoading] = useState<boolean>();
  const [page, setPage] = useState<number>(1);
  const [pageSize, setPageSize] = useState<number>(10);
  const [total, setTotal] = useState<number>(0);
  const [data, setData] = useState<any[]>([]);

  const editorRef = useRef<TypeEditorAction>();
  const [formRef] = Form.useForm();

  const fetchData = useCallback(async (sPage?: number, sPageSize?: number) => {
    const page = sPage || 1;
    const pageSize = sPageSize || 10;
    const searchItems = formRef.getFieldsValue();
    var payload = { page, pageSize, ...searchItems };
    setFetchLoading(true);
    //todo

    setFetchLoading(false);
    setData([{ name: '啊哈哈哈' }]);
    setPage(page);
    setPageSize(pageSize);
    setTotal(1);
  }, [formRef]);

  const refreshData = useCallback(() => {
    fetchData(page, pageSize);
  }, [fetchData, page, pageSize]);

  useEffect(() => {
    fetchData();
    // eslint-disable-next-line react-hooks/exhaustive-deps
  }, []);

  return (
    <>
      <div style={{ marginBottom: 12, display: 'flex' }} >
        <div>
          <Form layout='inline' form={formRef} >
            <Form.Item  >
              <Input placeholder="类别名称" />
            </Form.Item>
          </Form>
        </div>
        <div style={{ textAlign: 'right', width: '100%' }} >
          <Button type='primary' danger style={{ marginRight: 12 }} onClick={() => editorRef.current?.show()}  >新增类别</Button>
          <Button type='primary' onClick={() => refreshData()} >查询</Button>
        </div>
      </div>
      <Table<any>
        bordered
        size='middle'
        loading={fetchLoading}
        pagination={{ current: page, pageSize, total, showTotal: () => `共${total}条数据` }}
        dataSource={data}
        rowKey='id'
        columns={[
          {
            title: '序号', align: 'center',
            render: (t, r, i) => i + 1
          },
          {
            title: '类别名称', dataIndex: 'name',
            render: (t) => t ?? '-'
          },
          {
            title: '图标', dataIndex: 'icon',
            render: (t) => t ?? '-'
          },
          {
            title: '创建时间', dataIndex: 'createTime',
            render: (t, r) => momentFormat(r.createTime)
          },
          {
            title: '最后修改时间', dataIndex: 'lastEditTime',
            render: (t) => momentFormat(t)
          },
          {
            title: '操作',
            render: (t, row) => <>
              <a onClick={() => editorRef.current?.show(row)} >修改</a>
              <Divider type='vertical' />
              <Popconfirm title='确认删除?' onConfirm={() => { }} >
                <a style={{ color: 'red' }} >删除</a>
              </Popconfirm>
            </>
          },
        ]}
      />
      <TypeEditor ref={editorRef} onFinish={() => refreshData()} />
    </>
  )
}

export default Manage