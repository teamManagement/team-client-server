/* eslint-disable jsx-a11y/anchor-is-valid */
/* eslint-disable @typescript-eslint/no-unused-vars */
import React, { useState, useCallback, useEffect } from "react";
import { Table, Tabs, Form, Button, Input, } from "antd";
import { LeftOutlined } from '@ant-design/icons'
import AuditPage from "./AuditPage";
import AppLeiBieSC from "./components/AppLeiBieSC";
import AppLeiXingSC from "./components/AppLeiXingSC";
import './index.less'

enum AppCheckStatus {
  waitCheck = 0,
  all = 10,
}

const AppManager: React.FC = () => {

  const [fetchLoading, setFetchLoading] = useState<boolean>();
  const [page, setPage] = useState<number>(1);
  const [pageSize, setPageSize] = useState<number>(10);
  const [total, setTotal] = useState<number>(0);
  const [data, setData] = useState<any[]>([]);

  const [checkStatus, setCheckStatus] = useState<AppCheckStatus>(AppCheckStatus.waitCheck);
  const [isAduit, setIsAduit] = useState<boolean>(false);
  const [nowAduit, setNowAduit] = useState<any>(false);

  const [formRef] = Form.useForm();

  const fetchData = useCallback(async (sPage?: number, sPageSize?: number) => {
    const page = sPage || 1;
    const pageSize = sPageSize || 10;
    const searchItems = formRef.getFieldsValue();
    var payload = { page, pageSize, ...searchItems, checkStatus };
    setFetchLoading(true);
    //todo

    setFetchLoading(false);
    setData([{ id: '1', mname: '测试' }]);
    setPage(page);
    setPageSize(pageSize);
    setTotal(1);
  }, [formRef, checkStatus]);

  const refreshData = useCallback(() => {
    fetchData(page, pageSize);
  }, [fetchData, page, pageSize]);

  useEffect(() => {
    if (!isAduit) {
      fetchData();
    }
    // eslint-disable-next-line react-hooks/exhaustive-deps
  }, [checkStatus, isAduit]);



  if (isAduit) {
    return <div>
      <div style={{ marginBottom: 24 }} >
        <Button icon={<LeftOutlined />} type='primary' ghost onClick={() => setIsAduit(false)} >返回列表</Button>
      </div>
      <AuditPage onCompleted={() => { setIsAduit(false) }} appId={nowAduit?.id} />
    </div>
  }

  return <div className="appManager" >
    <div>
      <Tabs
        activeKey={`${checkStatus}`}
        onChange={(e) => setCheckStatus(Number.parseInt(e))}
        items={[
          { label: '待审核应用', key: `${AppCheckStatus.waitCheck}` },
          { label: '所有应用', key: `${AppCheckStatus.all}` },
        ]} />
    </div>
    <div style={{ marginBottom: 12, display: 'flex' }} >
      <div style={{ width: '100%' }} >
        <Form layout='inline' form={formRef} >
          <Form.Item  >
            <Input placeholder="应用名称" />
          </Form.Item>
          <Form.Item >
            <Input placeholder="持有人" />
          </Form.Item>
          <Form.Item >
            <AppLeiBieSC />
          </Form.Item>
          <Form.Item  >
            <AppLeiXingSC />
          </Form.Item>
        </Form>
      </div>
      <div style={{ textAlign: 'right', }} >
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
          title: '应用名称', dataIndex: 'name',
          render: (t) => t ?? '-'
        },
        {
          title: '持有者', dataIndex: 'holder',
          render: (t) => t ?? '-'
        },
        {
          title: '应用类别', dataIndex: 'leibie',
          render: (t) => t ?? '-'
        },
        {
          title: '应用类型', dataIndex: 'leixing',
          render: (t) => t ?? '-'
        },
        {
          title: '图标', dataIndex: 'icon',
          render: (t) => t ?? '-'
        },
        {
          title: '短描述', dataIndex: 'showDesc',
          render: (t) => t ?? '-'
        },
        {
          title: '应用状态', dataIndex: 'auditStatus',
          render: (t) => t ?? '-'
        },
        {
          title: '是否上推荐', dataIndex: 'isTuiJian',
          render: (t) => t ?? '-'
        },
        {
          title: '是否在推荐页隐藏', dataIndex: 'isHideInTuiJian',
          render: (t) => t ?? '-'
        },
        {
          title: '操作',
          render: (t, r) => <>
            <a onClick={() => {
              setNowAduit(r);
              setIsAduit(true);
            }} >审核</a>
          </>
        }
      ]}
    />
  </div>
}

export default AppManager;