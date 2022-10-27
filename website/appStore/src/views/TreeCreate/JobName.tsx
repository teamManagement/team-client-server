import ProForm, { ProFormText } from "@ant-design/pro-form"
import { Button, Form, FormInstance, Input, message, Modal } from "antd"
import { useCallback, useEffect, useRef, useState } from "react"
import { addJob, addPost, apiPostRequest, updateJob, updatePost } from "../../serve"

interface ISignProps {
  fns: any,
  finished: any,
  type: 'job' | 'post'
}

const JobName: React.FC<ISignProps> = (props) => {
  const { type } = props
  const [visible, setVisible] = useState<boolean>(false)
  const [loading, setLoading] = useState<boolean>(false)
  const [orgId, setOrgId] = useState<any>('')
  const [info, setInfo] = useState<any>()

  const formRef = useRef<any>()
  useEffect(() => {
    props.fns.current = {
      show(orgId: any) {
        setVisible(true)
        setOrgId(orgId)
        if (orgId.name) {
          setInfo(orgId)
        }
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
    if (info) {
      console.log(orgId);
      if (type === 'job') {
        await updateJob({ id: info?.id, name: formValue.name, orgId: info.orgId })
      } else {
        await updatePost({ id: info?.id, name: formValue.name, orgId: info.orgId })
      }
      message.success('修改成功！')
    } else {
      if (type === 'job') {
        await addJob({ name: formValue.name, orgId: orgId })
      } else {
        await addPost({ name: formValue.name, orgId: orgId })
      }
      message.success('新增成功！')
    }
    setVisible(false)
    props.finished()
    setLoading(false)
  }, [props, orgId, type])

  return (
    <>
      <Modal onOk={onSave}
        okText='确定'
        cancelText='取消'
        title={type === 'job' ? '新增职位' : '新增岗位'}
        width='45vw'
        keyboard={false}
        className="sealModal"
        maskClosable={false}
        destroyOnClose
        open={visible}
        onCancel={() => {
          setVisible(false)
          setInfo([])
        }}
        footer={[
          <>
            <Button onClick={() => {
              setVisible(false)
              setInfo([])
            }}>取消</Button>
            <Button loading={loading} onClick={onSave} type='primary'>确定</Button>
          </>
        ]}
      >
        <Form ref={formRef} layout='vertical' initialValues={info}>
          <Form.Item label={type === 'job' ? '职位名称' : '岗位名称'} name='name'>
            <Input />
          </Form.Item>
        </Form>
      </Modal>
    </>
  )
}

export default JobName