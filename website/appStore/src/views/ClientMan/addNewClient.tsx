
import { Form, FormInstance, Input, message, Modal } from "antd";
import { FC, useCallback, useEffect, useRef, useState } from "react";
import FileUpload from "../../components/FileUpload";
import { delEmptyQueryNodes, removeBase64Prefix } from "../../components/utils";

interface ICustomProps {
  fns: any,
  finished: () => void,
}

const AddNewClientModal: FC<ICustomProps> = ({ fns, finished }) => {
  const [open, setOpen] = useState<boolean>(false)
  const formRef = useRef<any>()

  useEffect(() => {
    fns.current = {
      show() {
        setOpen(true)
      },
      hide() {
        setOpen(false)
      }
    }
  }, [fns])


  const onSave = useCallback(async () => {
    const formValue = formRef.current?.getFieldsValue()
    formValue.file = removeBase64Prefix(formValue.file)
    const newValue = delEmptyQueryNodes(delEmptyQueryNodes(formValue))
    console.log(newValue);
    message.success('版本创建成功！')
    setOpen(false)
    finished()
  }, [])

  return (
    <>
      <Modal
        className="version-modal"
        title='创建版本'
        open={open}
        onCancel={() => setOpen(false)}
        onOk={onSave}
        destroyOnClose
        keyboard={false}
        maskClosable={false}
      >
        <Form ref={formRef} layout='vertical'>
          <Form.Item label='版本' name='version' rules={[{ required: true }]}>
            <Input />
          </Form.Item>
          <Form.Item label='描述' name='desc'>
            <Input />
          </Form.Item>
          <Form.Item label='版本文件' name='file' rules={[{ required: true }]}>
            <FileUpload tipTitle="选择文件" accept=".zip" placeHolder='选择zip文件,1M以内' maxSize={1 * 1024 * 1024} />
          </Form.Item>
        </Form>
      </Modal>
    </>
  )
}

export default AddNewClientModal