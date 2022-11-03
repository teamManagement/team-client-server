import { Input, Modal } from "antd";
import { FC, useCallback, useEffect, useRef, useState } from "react";
import { ProForm, ProFormInstance, ProFormText } from '@ant-design/pro-form';
import './index.less';


interface ICustomProps {
  fns: any,
  finished: () => void
}

const CustomDef: FC<ICustomProps> = ({ fns, finished }) => {
  const [open, setOpen] = useState<boolean>(false)
  const formRef = useRef<ProFormInstance>()

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

  const onSave = useCallback(() => {
    formRef.current?.submit()
    const formValue = formRef.current?.getFieldsValue()
    if (!formValue.name) {
      return
    }
    finished()
  }, [finished])


  return (
    <>
      <Modal
        className="custom-modal"
        title='自定义方案'
        open={open}
        onCancel={() => setOpen(false)}
        destroyOnClose
        okText='确定'
        cancelText='取消'
        onOk={() => onSave()}
        keyboard={false}
        maskClosable={false}
      >
        <ProForm formRef={formRef} layout='vertical'>
          <ProFormText label='方案名称' name='name' placeholder="请输入自定义方案名称" rules={[{ required: true }]} />
        </ProForm>
      </Modal>
    </>
  )
}

export default CustomDef