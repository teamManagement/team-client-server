import { Form, Input, Modal } from "antd"
import { useState, useCallback, useImperativeHandle, forwardRef } from "react"
import FileUpload from "../../components/FileUpload"
import './index.less'

export interface TypeEditorProps {
  onFinish?: () => void;
}

export interface TypeEditorAction {
  show: (initValue?: any) => void;
  hide: () => void;
}

const TypeEditor = forwardRef((props: TypeEditorProps, ref) => {

  const [visible, setVisible] = useState<boolean>(false);
  const [editValue, setEditValue] = useState<any>(undefined);
  const [formRef] = Form.useForm();

  const show = useCallback((initValue?: any) => {
    setVisible(true);
    initValue && formRef.setFieldsValue({ ...initValue });
    setEditValue(initValue);
  }, [formRef]);

  const hide = useCallback(() => {
    setVisible(false);
    formRef.resetFields();
    setEditValue(undefined);
  }, [formRef]);

  useImperativeHandle(ref, () => ({ show, hide }), [show, hide]);

  const onSubmit = useCallback((formVal: any) => {
    console.info(formVal);
    if (editValue?.id) {

    }
  }, [editValue]);

  const formLabelCol = { span: 5 };
  const formWarpperCol = { span: 18 };

  return (
    <Modal title='新增类别' maskClosable={false} className="addTypes" visible={visible} onCancel={() => hide()} onOk={() => formRef.submit()}  >
      <Form form={formRef} layout='horizontal' labelCol={formLabelCol} wrapperCol={formWarpperCol} onFinish={(e) => onSubmit(e)} >
        <Form.Item label='类别名称' name='name' rules={[{ required: true, message: '请填写类别名称' }]} >
          <Input placeholder='请输入类别名称' />
        </Form.Item>
        <Form.Item label='图标' name='icon' >
          <FileUpload showMode='picture' accept='.png' placeHolder='请选择png文件' maxSize={3 * 1024 * 1024} />
        </Form.Item>
      </Form>
    </Modal>
  )

});


export default TypeEditor;