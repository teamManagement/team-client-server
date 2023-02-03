import { Button, Col, Form, Input, Row, Table, Space, Popconfirm, message, Radio } from 'antd';
import { useAntdTable } from 'ahooks';
import { useRef, useState } from 'react';
import AddNewClientModal from './addNewClient';
import { clientList } from '../../serve';

interface Item {
  version: string,
  desc: string,
  createTime: any,
}

interface Result {
  total: number;
  list: Item[];
}

const ClientManage: React.FC = () => {
  const [form] = Form.useForm();
  const fnsRef = useRef<any>()
  const fnsDocRef = useRef<any>()
  const [status, setStatus] = useState<any>('all')
  const [formValue, setFormValue] = useState<any>()
  const getTableData = async (current?: any, pageSize?: any, formData?: any): Promise<Result> => {
    let query = `page=${current}&size=${pageSize}`;
    Object.entries(formData).forEach(([key, value]) => {
      if (value) {
        query += `&${key}=${value}`;
      }
    });

    // const clientInfo = await clientList()
    // console.log(clientInfo);

    return {
      total: 1,
      list: []
    }
  }

  const { tableProps, search, params } = useAntdTable(getTableData, {
    defaultPageSize: 5,
    form,
  });

  const { submit, reset, } = search;

  const columns = [
    { title: '版本', dataIndex: 'version', key: 'version' },
    { title: '描述', dataIndex: 'desc', key: 'desc' },
    { title: '创建时间', dataIndex: 'createTime', key: 'createTime' },
    { title: '是否为当前版本', dataIndex: 'ifMainVersion', key: 'ifMainVersion', },
    { title: '资源是否被清除', dataIndex: 'sourceIfClearn', key: 'sourceIfClearn', },
    {
      title: '操作', dataIndex: 'option', key: 'option', render: (node: any, row: any) => {
        return <Space>

        </Space>
      }
    },
  ];

  return (
    <>
      <div style={{ textAlign: "right" }}>
        <Button type='primary' onClick={() => fnsRef.current.show()}>创建版本</Button>
      </div>
      <Table style={{ marginTop: 20 }} columns={columns} rowKey="email" {...tableProps} />
      <AddNewClientModal fns={fnsRef} finished={() => { }} />
    </>
  )
}
export default ClientManage
