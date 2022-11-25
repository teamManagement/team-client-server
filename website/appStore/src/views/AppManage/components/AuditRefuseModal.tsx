import { useState, useImperativeHandle, forwardRef, useCallback } from 'react';
import { Form, Input, Modal } from 'antd';

interface IAuditRefuseModalProps {
    onCompleted?: () => void;
}

export interface AuditRefuseModalActionType {
    show: (auditInfo: any) => void;
    hide: () => void;
}

const AuditRefuseModal = forwardRef((props: IAuditRefuseModalProps, ref) => {

    const [visible, setVisible] = useState<boolean>();
    const [auditInfo, setAduitInfo] = useState<any>();
    const [submitLoading, setSubmitLoading] = useState<boolean>();

    const [formRef] = Form.useForm();

    const show = useCallback((auditInfo: any) => {
        setVisible(true);
        setAduitInfo(auditInfo);
    }, []);

    const hide = useCallback(() => {
        setVisible(false);
    }, []);

    const onSubmit = useCallback(({ reason: string }: any) => {
        setSubmitLoading(true);

        var id = auditInfo?.id;
        //调用接口todo

        setSubmitLoading(false);

        props.onCompleted?.();
    }, [auditInfo?.id, props]);

    useImperativeHandle(ref, () => ({ show, hide }), [hide, show]);


    return (
        <Modal open={visible} title='输入拒绝理由' maskClosable={false} destroyOnClose onCancel={() => hide()} onOk={() => formRef.submit()} okButtonProps={{ loading: submitLoading }} >
            <Form form={formRef} onFinish={(e) => onSubmit(e)} >
                <Form.Item name='reason' rules={[{ required: true, message: '请输入拒绝理由' }]} >
                    <Input placeholder='请输入拒绝理由' />
                </Form.Item>
            </Form>
        </Modal>
    )
});

export default AuditRefuseModal;