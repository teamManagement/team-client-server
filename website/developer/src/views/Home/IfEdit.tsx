import { Button, Divider, Input, Tooltip } from "antd";
import { FC, useEffect, useState } from "react";
import { EditOutlined, CheckOutlined, CloseOutlined } from '@ant-design/icons'

interface IEditProps {
  editText: any,
  type: 'input' | 'area',
  finished: (name: any) => void
}

const IfEditWar: FC<IEditProps> = (props) => {

  const [isEdit, setIsEdit] = useState<boolean>(false)
  const [nowEditText, setNowEditText] = useState<string>('')

  useEffect(() => {
    setNowEditText(props.editText)
  }, [props])


  return (
    <>
      {!isEdit && <div className='item'>
        {props.editText ? props.editText : nowEditText}
        <span style={{ marginLeft: 10 }}>
          <Tooltip title='编辑'>
            <Button type='link' onClick={() => setIsEdit(true)}><EditOutlined /></Button>
          </Tooltip>
        </span>
      </div>}
      {isEdit && <div className='item' style={{ marginTop: 20 }}>
        <span onClick={(event) => event.stopPropagation()}>
          {props.type === 'input' && <Input
            style={{ width: 120 }} size='small' defaultValue={'4234'} value={nowEditText} maxLength={20}
            onChange={(e) => setNowEditText(e.target.value)}
            onPressEnter={(event) => { setIsEdit(false); event.stopPropagation(); props.finished(nowEditText) }}
            onClick={(event) => event.stopPropagation()}
            placeholder='输入名称'
          />}
          {props.type === 'area' && <Input.TextArea
            style={{ width: 240 }} size='small' value={nowEditText} maxLength={100}
            onChange={(e) => setNowEditText(e.target.value)}
            onPressEnter={(event) => { setIsEdit(false); event.stopPropagation(); props.finished(nowEditText) }}
            onClick={(event) => event.stopPropagation()}
            placeholder='请输入短描述，最多100个字'
          />}
          <span style={{ margin: 10 }}>
            <a><CheckOutlined onClick={(event) => { setIsEdit(false); event.stopPropagation(); props.finished(nowEditText) }} /></a>
            <Divider type='vertical' />
            <a><CloseOutlined onClick={(event) => { setIsEdit(false); event.stopPropagation(); }} /></a>
          </span>
        </span>
      </div>}
    </>
  )
}

export default IfEditWar