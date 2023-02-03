/* eslint-disable jsx-a11y/anchor-is-valid */
/* eslint-disable jsx-a11y/anchor-has-content */
/* eslint-disable @typescript-eslint/no-unused-vars */
import { useState, useCallback, useEffect, useRef } from 'react'
import { Button, Form, Input, message, Popconfirm, Table } from 'antd'
import AddManageUserModal, { AddManageUserModalActionType } from './AddManageUserModal';
import { delManageUsers, reqManList } from '../../serve';
import { current } from '@byzk/teamwork-sdk';


const ManageUserMgr: React.FC = () => {

    const [fetchLoading, setFetchLoading] = useState<boolean>();
    const [page, setPage] = useState<number>(1);
    const [pageSize, setPageSize] = useState<number>(10);
    const [total, setTotal] = useState<number>(0);
    const [data, setData] = useState<any[]>([]);
    const loginId = useRef<string>()

    const [formRef] = Form.useForm();
    const addUserRef = useRef<AddManageUserModalActionType>();

    const geFirst = useCallback(() => {
        loginId.current = current.userInfo.id
    }, [])

    useEffect(() => {
        geFirst()
    }, [geFirst])

    const fetchData = useCallback(async (sPage?: number, sPageSize?: number) => {
        const page = sPage || 1;
        const pageSize = sPageSize || 10;
        const searchItems = formRef.getFieldsValue();
        var payload = { page, pageSize, ...searchItems };
        setFetchLoading(true);
        //todo
        const list = await reqManList()
        console.log(list);
        setFetchLoading(false);
        setData(list);
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


    console.log(loginId.current);


    return <div>
        <div style={{ marginBottom: 12, display: 'flex' }} >
            <div>
                <Form layout='inline' form={formRef} >
                    <Form.Item>
                        <Input placeholder='人员姓名' />
                    </Form.Item>
                </Form>
            </div>
            <div style={{ textAlign: 'right', width: '100%' }} >
                <Button type='primary' style={{ marginRight: 12 }} danger onClick={() => addUserRef.current?.show()} >新增管理员</Button>
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
                    title: '姓名', dataIndex: 'name',
                    render: (t, r, i) => r.user.name ?? '-'
                },
                {
                    title: '所在部门', dataIndex: 'department',
                    render: (t, r, i) => r.user.nowOrgInfo.org.name ?? '-'
                },
                {
                    title: '联系电话', dataIndex: 'phone',
                    render: (t, r, i) => r.user.phone ?? '-'
                },
                {
                    title: '备注', dataIndex: 'remark',
                    render: (t, r, i) => t ?? '-'
                },
                {
                    title: '操作',
                    render: (t, r, i) => <div>
                        {loginId.current && loginId.current !== r.userId && <Popconfirm title='确认要删除吗?' onConfirm={async () => {
                            await delManageUsers(r.userId)
                            message.success('删除成功！')
                        }} >
                            <Button type='link' danger>删除</Button>
                        </Popconfirm>}
                    </div>
                }
            ]}
        />
        <AddManageUserModal ref={addUserRef} onCompleted={() => refreshData()} />
    </div>
}

export default ManageUserMgr;