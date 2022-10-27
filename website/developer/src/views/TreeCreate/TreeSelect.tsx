import { Divider, Input, message, Modal, Popconfirm, Spin, Tooltip, Tree } from "antd"
import { FC, useCallback, useEffect, useState } from "react"
import { isNull } from "../../components/utils";
import { v4 as uuid } from 'uuid'
import { CheckOutlined, CloseOutlined, EditOutlined, PlusOutlined } from '@ant-design/icons';
import './index.less';
import { addOrg, deleteOrg, getOrgList, upDateOrg } from "../../serve";

interface ITreeProps {
  firstList: any,
  getOrgId: (orgId: string) => void
}

const TreeSelectWar: FC<ITreeProps> = (props) => {
  const { firstList } = props
  const [treeData, setTreeData] = useState<any[]>([])
  const [loading, setLoading] = useState<boolean>(false)
  const [nowSelectedKey, setNowSelectedKey] = useState<string>()
  const [nowEditId, setNowEditId] = useState<string>()
  const [nowEditTxt, setNowEditTxt] = useState<string>()
  const [addNodeRecursionEnd, setAddNodeRecursionEnd] = useState<boolean>(false)
  const [deleteNodeRecursionEnd, setDeleteNodeRecursionEnd] = useState<boolean>(false)
  const [editRecursionEnd, setEditRecursionEnd] = useState<boolean>(false)
  const [virtualNodeKey, setVirtualNodeKey] = useState<string>()
  const [pid, setPid] = useState<string>()



  const loadTree = useCallback((firstList: any) => {
    if (firstList) {
      const newList = firstList.map((m: any) => {
        let children: any[] = []
        if (m.children) {
          children = m.children.map((m: any) => {
            return {
              key: m.id,
              title: m.name,
              pid: m.pid,
              children: m.children ? children : []
            }
          })
        }

        return {
          key: m.id,
          title: m.name,
          pid: m.pid,
          children: children
        }
      })
      setTreeData(newList)
    }
  }, [])

  const getTreeList = useCallback(async () => {
    const newlist = await getOrgList({})
    loadTree(newlist)
  }, [])

  useEffect(() => {
    loadTree(firstList)
  }, [loadTree])


  const hideEdit = useCallback(() => {
    setNowEditId(undefined)
    setNowEditTxt(undefined)
  }, [])

  const onEditOk = useCallback(async () => {
    if (isNull(nowEditTxt)) {
      Modal.error({ title: '请输入名称' });
      return;
    }
    if (!isNull(pid)) {
      await addOrg({ name: nowEditTxt, pid: pid })
      message.success('添加成功！')
    } else {
      await upDateOrg({ name: nowEditTxt, pid: pid })
      message.success('修改成功！')
    }
    getTreeList()
    hideEdit()
  }, [nowEditTxt, pid, hideEdit])

  const onEditCancel = useCallback(() => {
    console.log(nowEditId, virtualNodeKey)
    if (nowEditId === virtualNodeKey) {
      onDeleteNode(nowEditId, '');
    } else {
      hideEdit();
    }
  }, [nowEditId, virtualNodeKey])

  const addNode = useCallback((treeData: any[], key: any) => {
    for (let index = 0; index < treeData.length; index++) {
      if (addNodeRecursionEnd) {
        return;
      }
      const element = treeData[index];
      if (element.key === key) {
        if (isNull(element.children)) {
          element.children = []
        }
        var newKey = uuid();
        element.children.push({ key: newKey, title: '' });
        setNowEditId(newKey)
        setVirtualNodeKey(newKey)
        setNowEditTxt('')
        setPid(key)
        setAddNodeRecursionEnd(true)
        return;
      }
      if (!isNull(element.children)) {
        addNode(element.children, key);
      }
    }
  }, [])

  const onAddNode = useCallback((nodeId: any) => {
    let newTreeData = [...treeData]
    setAddNodeRecursionEnd(false)
    addNode(newTreeData, nodeId)
    setTreeData(newTreeData)
  }, [treeData])

  const editNode = useCallback((key: any, title: any) => {
    setNowEditId(key)
    setNowEditTxt(title)
  }, [])

  const deleteNode = useCallback((treeNode: any, key: any) => {
    if (deleteNodeRecursionEnd) {
      return;
    }
    if (isNull(treeNode.children) || treeNode.children!.length === 0) {
      return;
    }
    if (!isNull(treeNode.children!.find((m: any) => m.key === key))) {
      treeNode.children = [...treeNode.children!.filter((m: any) => m.key !== key)];
      setDeleteNodeRecursionEnd(true)
      return;
    }
    for (let index = 0; index < treeNode.children.length; index++) {
      deleteNode(treeNode.children[index], key);
    }
  }, [])

  const onDeleteNode = useCallback(async (nodeId: any, pid: any) => {
    if (virtualNodeKey === nodeId) {
      var newTreeData = [...treeData];
      setDeleteNodeRecursionEnd(false)
      deleteNode(newTreeData[0], nodeId);
      setTreeData(newTreeData)
      hideEdit()
    } else {
      await deleteOrg(pid, nodeId)
      getTreeList()
      message.success('删除成功！')
    }
  }, [virtualNodeKey])

  const onNodeClick = useCallback((nodeId: any) => {
    console.log(nodeId);
    props.getOrgId(nodeId)
  }, [])


  const editTitle = useCallback((treeData: any[], key: any, title: any) => {
    for (let index = 0; index < treeData.length; index++) {
      if (editRecursionEnd) {
        return;
      }
      const element = treeData[index];
      if (element.key === key) {
        element.title = title;
        setEditRecursionEnd(true)
        return;
      }
      if (!isNull(element.children)) {
        editTitle(element.children, key, title);
      }
    }
  }, [])

  return (
    <Spin spinning={loading} tip='内容正在加载...'>
      <div className="departmentTreeBox">
        {treeData?.length > 0 &&
          <Tree
            defaultExpandAll
            defaultExpandParent
            autoExpandParent
            selectable={false}
            treeData={treeData}
            titleRender={(node) =>
              <div
                className={(nowSelectedKey === node.key && nowEditId !== node.key) ? 'departmentTreeItem departmentTreeItemClicked' : nowEditId === node.key ? '' : 'departmentTreeItem'}
                onClick={(e) => { onNodeClick(node.key) }}>
                {
                  nowEditId === node.key ?
                    <span onClick={(event) => event.stopPropagation()}>
                      <Input
                        style={{ width: 120 }} size='small' value={nowEditTxt} maxLength={100}
                        onChange={(e) => setNowEditTxt(e.target.value)}
                        onPressEnter={() => onEditOk()}
                        onClick={(event) => event.stopPropagation()}
                        placeholder='输入名称'
                      />
                      <span className='departmentTreeEditingBtn' >
                        <CheckOutlined onClick={(event) => { onEditOk(); event.stopPropagation(); }} />
                        <Divider type='vertical' />
                        <CloseOutlined onClick={(event) => { onEditCancel(); event.stopPropagation(); }} />
                      </span>
                    </span>
                    :
                    <span className='departmentTreeNodeTitle' >
                      {node.title}
                    </span>
                }
                {
                  node.key === treeData[0].key &&
                  <span className='departmentTreeHoverBox' style={{ display: 'inline' }} >
                    <Tooltip title='增加子部门' >
                      <PlusOutlined onClick={(event) => { onAddNode(node.key); event.stopPropagation(); }} />
                    </Tooltip>
                  </span>
                }
                {
                  (isNull(nowEditId) && node.key !== treeData[0].key) &&
                  <span className='departmentTreeHoverBox' >
                    <Tooltip title='增加子部门' >
                      <PlusOutlined onClick={(event) => { onAddNode(node.key); event.stopPropagation(); }} />
                    </Tooltip>
                    <Divider type='vertical' />
                    <Tooltip title='编辑当前部门' >
                      <EditOutlined onClick={(event) => { editNode(node.key, node.title); event.stopPropagation(); }} />
                    </Tooltip>
                    <Divider type='vertical' />
                    <Popconfirm title='确定删除当前部门？'
                      onConfirm={(e) => { onDeleteNode(node.key, node.pid); e?.stopPropagation(); }}
                      onCancel={(e) => e?.stopPropagation()}
                      okText='确定'
                      cancelText='取消'
                    >
                      <Tooltip title='删除当前部门' placement='right' >
                        <CloseOutlined onClick={(event) => event.stopPropagation()} />
                      </Tooltip>
                    </Popconfirm>
                  </span>
                }
              </div>
            }
          />
        }
      </div>
    </Spin>
  );
}
export default TreeSelectWar