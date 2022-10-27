import { Form, FormInstance, Modal } from "antd"
import { useEffect, useRef, useState } from "react"
import { ProForm, ProFormText } from '@ant-design/pro-form'
import './index.less'
import FileUpload from "../../components/FileUpload"

interface ITypeProps {
  fns: any,
  finnished: () => void,
}

const AddTypes: React.FC<ITypeProps> = (props) => {

  const [visible, setVisible] = useState<boolean>(false)
  const formRef = useRef<FormInstance>()

  useEffect(() => {
    props.fns.current = {
      show(info: any) {
        setVisible(true)
        console.log(info);
        formRef.current?.setFieldsValue(info)
      },
      close() { setVisible(false) }
    }
  }, [])

  return (
    <>
      <Modal
        title='新增类别'
        className="addTypes"
        width={'40vw'}
        open={visible}
        onCancel={() => setVisible(false)}
        destroyOnClose
      >
        <ProForm formRef={formRef} layout='horizontal'>
          <ProFormText label='名称' name='name' placeholder='请输入类别名称' />
          <Form.Item label='图标' name='icon'>
            <FileUpload valueMode='base64' showMode='picture' accept=".png" placeHolder='请选择png/jpeg文件' maxSize={1024 * 1024 * 1024} />
          </Form.Item>
        </ProForm>
      </Modal>
    </>
  )
}

export default AddTypes