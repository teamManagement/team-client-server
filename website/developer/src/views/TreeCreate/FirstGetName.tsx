import ProForm, { ProFormText } from "@ant-design/pro-form"
import { Button, Form, FormInstance, Input, message, Modal } from "antd"
import { useCallback, useEffect, useRef, useState } from "react"

interface ISignProps {
  fns: any,
  finished: any
}

const FirstGetName: React.FC<ISignProps> = (props) => {
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
    try {
      await window.proxyApi.httpWebServerProxy('/org/add', { jsonData: formValue })
      setVisible(false)
      props.finished()
    } catch (e: any) {
      message.error(e)
    }
    setLoading(false)
  }, [props])

  return (
    <>
      <Modal onOk={onSave}
        okText='确定'
        cancelText='取消'
        title='新增人员管理'
        width='45vw'
        keyboard={false}
        className="sealModal"
        maskClosable={false}
        destroyOnClose
        open={visible}
        onCancel={() => setVisible(false)}
        footer={[
          <>
            <Button onClick={() => setVisible(false)}>取消</Button>
            <Button loading={loading} onClick={onSave} type='primary'>确定</Button>
          </>
        ]}
      >
        <Form ref={formRef} layout='vertical'>
          <Form.Item label='机构名称' name='name'>
            <Input />
          </Form.Item>
        </Form>
      </Modal>
    </>
  )
}

export default FirstGetName