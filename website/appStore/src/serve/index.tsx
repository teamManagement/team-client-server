import { Modal } from "antd";
import { api } from '@teamworktoolbox/inside-sdk'


export const getSignCert = async (errorCb?: any) => api.proxyHttpCoreServer('')

interface IHttpReq {
  url: string;
  data?: any;
  errorCb?: (error: string) => void;
  finalCb?: () => void;
}

export async function apiPostRequest(params: IHttpReq): Promise<any> {
  if (!params) {
    return Promise.reject('请求参数不能为空！')
  }
  try {
    const rst = await api.proxyHttpCoreServer(params.url, { jsonData: params.data })
    return Promise.resolve(rst)
  } catch (e: any) {
    Modal.error({ title: e.message, okText: '知道了' });
    params.errorCb?.(e.message);
    return
  } finally {
    params.finalCb?.()
  }
}

export async function apiLocalRequest(params: IHttpReq): Promise<any> {
  if (!params) {
    return Promise.reject('请求参数不能为空！')
  }
  try {
    const rst = await api.proxyHttpLocalServer(params.url, { jsonData: params.data, timeout: -1 })
    return Promise.resolve(rst)
  } catch (e: any) {
    Modal.error({ title: e.message, okText: '知道了' });
    params.errorCb?.(e.message);
    return
  } finally {
    params.finalCb?.()
  }
}

/**
 * 获取机构列表
 */
export async function getOrgList(data?: any): Promise<any> {
  return await apiPostRequest({ url: '/org/list', data })
}

/**
 * 添加机构
 */
export async function addOrg(data: any): Promise<any> {
  return await apiPostRequest({ url: '/org/add', data })
}

/**
 * 更新机构
 */
export async function upDateOrg(data: any,): Promise<any> {
  return await apiPostRequest({ url: '/org/update', data })
}

/**
 * 删除机构
 */
export async function deleteOrg(pid: any, id: any): Promise<any> {
  return await apiPostRequest({ url: `/org/del/${pid}/${id}` })
}



/**
 * 增加职位
 */
export async function addJob(data: any): Promise<any> {
  return await apiPostRequest({ url: `/org/job/add`, data })
}

/**
 * 职位列表
 */
export async function JobList(orgId: any): Promise<any> {
  return await apiPostRequest({ url: `/org/job/list/${orgId}` })
}

/**
 * 删除职位
 */
export async function deleteJob(orgId: any, jobId: any): Promise<any> {
  return await apiPostRequest({ url: `/org/job/del/${orgId}/${jobId}` })
}

/**
 * 更新职位
 */
export async function updateJob(data: any): Promise<any> {
  return await apiPostRequest({ url: "/org/job/update", data })
}


/**
 * 增加岗位
 */
export async function addPost(data: any): Promise<any> {
  return await apiPostRequest({ url: `/org/post/add`, data })
}

/**
 * 岗位列表
 */
export async function PostList(orgId: any): Promise<any> {
  return await apiPostRequest({ url: `/org/post/list/${orgId}` })
}

/**
 * 删除岗位
 */
export async function deletePost(orgId: any, jobId: any): Promise<any> {
  return await apiPostRequest({ url: `/org/post/del/${orgId}/${jobId}` })
}

/**
 * 更新岗位
 */
export async function updatePost(data: any): Promise<any> {
  return await apiPostRequest({ url: "/org/post/update", data })
}



/**
 * 获取类别列表
 * @param data 
 * @returns 
 */
export async function getTypeList(data: any): Promise<any> {
  return await apiPostRequest({ url: "app/category/list", data })
}



/**
 * 获取具体类别应用详情
 * @param data 
 * @returns 
 */
export async function getAppTypeList(appId: any): Promise<any> {
  return await apiPostRequest({ url: `app/store/${appId}` })
}


/**
* 获取具体类别应用列表
* @param data 
* @returns 
*/
export async function getAppUserList(appId: any): Promise<any> {
  return await apiPostRequest({ url: `/user/get/${appId}` })
}



//TODO 管理

/**
 * 管理员新增
 */
export async function addManageUsers(userId: any): Promise<any> {
  return await apiLocalRequest({ url: `/services/appstore/manager/add/${userId}` })
}

/**
 * 管理员删除
 */
export async function delManageUsers(userId: any): Promise<any> {
  return await apiLocalRequest({ url: `/services/appstore/manager/del/${userId}` })
}

/**
 * 管理员列表查询
 */
export async function reqManList(): Promise<any> {
  return await apiPostRequest({ url: `/appstore/manager/list` })
}


/**
 * 客户端版本列表查询
 */
export async function clientList(): Promise<any> {
  return await apiPostRequest({ url: `` })
}

/**
 * 客户端版本新增
 */
export async function addClientVersion(data: any): Promise<any> {
  return await apiLocalRequest({ url: ``, data })
}
