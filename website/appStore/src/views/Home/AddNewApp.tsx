

import { message, Modal, Button, } from "antd";
import { FC, useCallback, useEffect, useRef, useState } from "react";
import { ProForm, ProFormText } from '@ant-design/pro-form';
import { isNull } from "../../components/utils";
import './index.less'

interface IAddAppProps {
  fns: any,
  finished: (name: any) => void
}


const AddNewApp: FC<IAddAppProps> = (props) => {
  const [visible, setVisible] = useState<boolean>(false)
  const [loading, setLoading] = useState<boolean>(false)

  const formRef = useRef<any>()
  useEffect(() => {
    props.fns.current = {
      show() {
        setVisible(true)
      },
      close() {
        setVisible(false)
      }
    }
  }, [])

  const onSave = useCallback(async () => {
    setLoading(true)
    formRef.current?.submit()
    const formValue = formRef.current?.getFieldsValue()
    console.log(formValue)
    if (isNull(formValue.name)) {
      setLoading(false)
      return;
    }
    setVisible(false)
    props.finished(formValue.name)
    setLoading(false)
  }, [props])
  return (
    <>
      <Modal onOk={onSave}
        okText='确定'
        cancelText='取消'
        title='新增应用'
        width='45vw'
        keyboard={false}
        className="appModal"
        maskClosable={false}
        destroyOnClose
        visible={visible}
        onCancel={() => setVisible(false)}
        footer={[
          <>
            <Button onClick={() => setVisible(false)}>取消</Button>
            <Button loading={loading} onClick={onSave} type='primary'>确定</Button>
          </>
        ]}
      >
        <ProForm formRef={formRef}>
          <ProFormText label='应用名称' name='name' rules={[{ required: true }]} />
        </ProForm>
      </Modal>
    </>
  )
}

export default AddNewApp