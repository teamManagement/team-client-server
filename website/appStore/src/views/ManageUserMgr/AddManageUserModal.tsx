import { Modal, Row, Col, Select } from 'antd';
import { forwardRef, useImperativeHandle, useCallback, useState } from 'react'

export interface AddManageUserModalActionType {
    show: () => void;
    hide: () => void;
}

interface IProps {
    onCompleted?: () => void
}


const AddManageUserModal = forwardRef((props: IProps, ref) => {

    const [visible, setVisible] = useState<boolean>(false);
    const [loading, setLoading] = useState<boolean>(false);
    const [data, setData] = useState<any[]>();
    const [selected, setSelected] = useState<any>();
    const [submitLoading, setSubmitLoading] = useState<boolean>(false);

    const fetchData = useCallback(() => {
        setLoading(true);

        //todo

        setLoading(false);
        setData([
            { id: 1, name: '程亮', department: '电子签章业务部' },
            { id: 2, name: '张三', department: '测试部门' },
        ]);
    }, []);

    const show = useCallback(() => {
        setVisible(true);
        fetchData();
    }, [fetchData]);

    const hide = useCallback(() => {
        setVisible(false);
    }, []);


    useImperativeHandle(ref, () => ({ show, hide }), [show, hide]);

    const onComfirm = useCallback(() => {
        if (!selected) {
            Modal.error({ title: '请选择人员' });
            return;
        }
        setSubmitLoading(true);
        //todo
        setSubmitLoading(false);
        hide();
        props.onCompleted?.();
    }, [hide, props, selected]);



    return <Modal title='新增管理员' maskClosable={false} open={visible} onCancel={() => hide()} onOk={() => onComfirm()} okButtonProps={{ loading: submitLoading }} >
        <Row>
            <Col span={6} style={{ textAlign: 'right' }} >选择人员:</Col>
            <Col offset={1} style={{ marginBottom: 20 }}  >
                <Select
                    loading={loading}
                    showSearch
                    placeholder='输入想查询的人员姓名'
                    style={{ width: 250, marginTop: -5 }}
                    value={selected}
                    onChange={(e) => setSelected(e)}
                    options={(data || []).map((m) => { return { label: `${m.name}-${m.department}`, value: m.id } })}
                    filterOption={(input, option) => {
                        return (option?.label ?? '').includes(input)
                    }}
                />
            </Col>
        </Row>
    </Modal>
});


export default AddManageUserModal;