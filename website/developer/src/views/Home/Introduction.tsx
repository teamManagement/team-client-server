import { Button, Space, Tooltip } from 'antd';
import { FC, useState } from 'react';
import { EditOutlined } from '@ant-design/icons';
import BraftEditor from 'braft-editor';
import 'braft-editor/dist/index.css';
import 'braft-editor/dist/output.css'
import ImageBarWar from './ImageBar';

interface IImageProps{
  getFileMenuList:()=>void,
  getId:any
}

const Introduction: FC<IImageProps> = (props) => {
  const [ifEdit, setIfEdit] = useState<boolean>(true)
  const [braftValue, setBraftValue] = useState(BraftEditor.createEditorState(null))

  return (
    <>
      <div className="introduction">
        <div>
          <div style={{ marginBottom: 10, fontWeight: 'bold' }}>图片详情:</div>
          <ImageBarWar type='wall' getFileMenuList={props.getFileMenuList} getId={props.getId}/>
        </div>
        <div className='markDown'>
          <div style={{ marginBottom: 10, fontWeight: 'bold' }}>详细描述:
            <Tooltip title='编辑描述' >
              {!ifEdit && <Button type='link' className='edit-btn' onClick={() => {
                setBraftValue(BraftEditor.createEditorState(braftValue))
                setIfEdit(true)
              }}><EditOutlined /></Button>}
            </Tooltip>
          </div>
          <div className='markdown-content'>
            {!ifEdit && <div style={{ backgroundColor: "#f9f9f9", padding: '10px 10px' }} className="braft-output-content" dangerouslySetInnerHTML={{ __html: braftValue }}></div>}
            {ifEdit && <>
              <BraftEditor
                value={braftValue}
                onChange={(braftValue) => setBraftValue(braftValue.toHTML())}
              />
              <Button.Group style={{ position: 'absolute', right: 0, bottom: '10vh' }}>
                <Space><Button type='default' onClick={() => setIfEdit(false)}>取消</Button>
                  <Button type='primary' onClick={() => setIfEdit(false)}>确定</Button></Space>
              </Button.Group>
            </>}
          </div>
        </div>
      </div>
    </>
  )
}

export default Introduction