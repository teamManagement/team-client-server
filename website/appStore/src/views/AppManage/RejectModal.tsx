import ProForm, { ProFormTextArea } from "@ant-design/pro-form"
import { Button, FormInstance, Modal } from "antd"
import { useCallback, useEffect, useRef, useState } from "react"
import { isNull } from "../../components/utils"
import './index.less'

interface ISignProps {
  fns: any,
  finished: () => void,
}


const RejectModal: React.FC<ISignProps> = (props) => {
  const [visible, setVisible] = useState<boolean>(false)
  const [loading, setLoading] = useState<boolean>(false)
  const formRef = useRef<FormInstance>()

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
    if (isNull(formValue.reject)) {
      setLoading(false)
      return
    }
    setLoading(false)
    Modal.success({ title: '审核成功', okText: '知道了' })
    props.finished()
  }, [props])

  return (
    <>
      <Modal onOk={onSave}
        okText='确定'
        cancelText='取消'
        title='拒绝原因'
        width='45vw'
        keyboard={false}
        className="modalPdf"
        maskClosable={false}
        destroyOnClose
        open={visible}
        onCancel={() => setVisible(false)}
        footer={[
          <Button onClick={() => setVisible(false)}>取消</Button>,
          <Button loading={loading} onClick={onSave} type='primary'>确定</Button>
        ]}
      >
        <ProForm formRef={formRef} layout='vertical' className="form-modal">
          <ProFormTextArea label='拒绝原因' placeholder='请输入拒绝原因' name='reject' rules={[{ required: true }]} />
        </ProForm>
      </Modal>
    </>
  )
}

export default RejectModal