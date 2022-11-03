import { Modal } from "antd";



export const getSignCert = async (errorCb?: any) => window.proxyApi.httpLocalServerProxy('')

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
    const rst = await window.proxyApi.httpWebServerProxy(params.url, { jsonData: params.data })
    return Promise.resolve(rst)
  } catch (e: any) {
    Modal.error({ title: e.message, okText: '知道了' })
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
